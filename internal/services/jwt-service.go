package services

import (
	"context"
	"fmt"
	"time"

	"github.com/asrma7/meroshare-bot/pkg/config"
	"github.com/asrma7/meroshare-bot/pkg/errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type CustomClaims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

type JWTService interface {
	GenerateToken(userID uuid.UUID) (accessToken string, refreshToken string, err error)
	ValidateToken(token string) (*CustomClaims, error)
	RefreshToken(refreshToken string) (newAccessToken string, newRefreshToken string, err error)
	ExtractClaims(token string) (jwt.MapClaims, error)
	ParseToken(token string) (*jwt.Token, error)
}

type jwtService struct {
	accessSecret  string
	refreshSecret string
	accessTTL     time.Duration
	refreshTTL    time.Duration
	redisClient   *redis.Client
}

func NewJWTService(cfg *config.Config, redisClient *redis.Client) JWTService {
	return &jwtService{
		accessSecret:  cfg.AccessSecret,
		refreshSecret: cfg.RefreshSecret,
		accessTTL:     cfg.TokenExpiry,
		refreshTTL:    cfg.RefreshExpiry,
		redisClient:   redisClient,
	}
}

func (s *jwtService) GenerateToken(userID uuid.UUID) (string, string, error) {
	now := time.Now()

	accessClaims := &CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTTL)),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "meroshare-bot",
		},
	}

	refreshJTI := uuid.NewString()
	refreshClaims := &CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        refreshJTI,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshTTL)),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "meroshare-bot",
		},
	}

	accessToken, err := s.signToken(accessClaims, s.accessSecret)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.signToken(refreshClaims, s.refreshSecret)
	if err != nil {
		return "", "", err
	}

	key := "refresh:" + refreshJTI
	err = s.redisClient.Set(context.Background(), key, userID.String(), s.refreshTTL).Err()
	if err != nil {
		return "", "", fmt.Errorf("failed to store refresh token in redis: %w", err)
	}

	return accessToken, refreshToken, nil
}

func (s *jwtService) ValidateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.accessSecret), nil
	}, jwt.WithLeeway(0))

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.NewUnauthorizedError("token expired")
		}
		return nil, fmt.Errorf("invalid token %w", err)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

func (s *jwtService) RefreshToken(refreshToken string) (string, string, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &CustomClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.refreshSecret), nil
	}, jwt.WithLeeway(0))

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", "", errors.NewUnauthorizedError("token expired")
		}
		return "", "", fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return "", "", errors.NewUnauthorizedError("invalid token claims")
	}

	jti := claims.ID
	if jti == "" {
		return "", "", errors.NewUnauthorizedError("invalid token: missing JTI")
	}

	key := "refresh:" + jti
	ctx := context.Background()

	storedUserID, err := s.redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", "", errors.NewForbiddenError("Refresh token has already been used or is invalid")
	} else if err != nil {
		return "", "", fmt.Errorf("failed to verify refresh token: %w", err)
	}

	if storedUserID != claims.UserID.String() {
		return "", "", errors.NewForbiddenError("Refresh token user mismatch")
	}

	if err := s.redisClient.Del(ctx, key).Err(); err != nil {
		return "", "", fmt.Errorf("failed to invalidate refresh token: %w", err)
	}

	return s.GenerateToken(claims.UserID)
}

func (s *jwtService) ExtractClaims(tokenString string) (jwt.MapClaims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

func (s *jwtService) ParseToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (any, error) {
		return []byte(s.accessSecret), nil
	})
}

func (s *jwtService) signToken(claims *CustomClaims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return signedToken, nil
}
