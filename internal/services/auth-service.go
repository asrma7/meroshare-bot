package services

import (
	"fmt"
	"strings"

	"github.com/asrma7/meroshare-bot/internal/models"
	"github.com/asrma7/meroshare-bot/internal/repositories"
	"github.com/asrma7/meroshare-bot/pkg/config"
	"github.com/asrma7/meroshare-bot/pkg/errors"
	"github.com/asrma7/meroshare-bot/pkg/utils"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type AuthService interface {
	RegisterUser(username, password, email, firstname, lastname string) (uuid.UUID, error)
	LoginUser(username, password string) (string, string, error)
	RefreshToken(refreshToken string) (newAccessToken string, newRefreshToken string, err error)
	GetProfile(userID uuid.UUID) (*models.User, error)
	ValidateToken(token string) (*CustomClaims, error)
}

type authService struct {
	cfg        *config.Config
	userRepo   repositories.UserRepository
	jwtService JWTService
}

func NewAuthService(cfg *config.Config, userRepo *repositories.UserRepository, redis *redis.Client) AuthService {
	return &authService{
		cfg:        cfg,
		userRepo:   *userRepo,
		jwtService: NewJWTService(cfg, redis),
	}
}

func (s *authService) RegisterUser(username, password, email, firstname, lastname string) (uuid.UUID, error) {
	_, err := s.userRepo.GetUserByUsername(username)
	if err == nil {
		return uuid.Nil, errors.NewConflictError("username already exists")
	}
	_, err = s.userRepo.GetUserByEmail(email)
	if err == nil {
		return uuid.Nil, errors.NewConflictError("email already exists")
	}

	// Further validation
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return uuid.Nil, errors.NewInternalError(err)
	}

	user := &models.User{
		Username:  username,
		Email:     email,
		Password:  hashedPassword,
		FirstName: firstname,
		LastName:  lastname,
	}

	userID, err := s.userRepo.CreateUser(user)
	if err != nil {
		return uuid.Nil, errors.NewInternalError(err)
	}

	return userID, nil
}

func (s *authService) LoginUser(identifier, password string) (string, string, error) {
	user, err := s.userRepo.GetUserByUsername(identifier)
	if err != nil {
		user, err = s.userRepo.GetUserByEmail(identifier)
		if err != nil {
			return "", "", errors.NewUnauthorizedError("invalid username/email or password")
		}
	}
	if err := utils.CheckPasswordHash(password, user.Password); err != nil {
		return "", "", errors.NewUnauthorizedError("invalid username/email or password")
	}
	accessToken, refreshToken, err := s.jwtService.GenerateToken(user.ID)
	if err != nil {
		return "", "", errors.NewInternalError(err)
	}
	return accessToken, refreshToken, nil
}

func (s *authService) RefreshToken(refreshToken string) (string, string, error) {
	claims, err := s.jwtService.ExtractClaims(refreshToken)
	if err != nil {
		return "", "", errors.NewUnauthorizedError("invalid refresh token")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", "", errors.NewUnauthorizedError("invalid refresh token claims")
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		return "", "", errors.NewUnauthorizedError("invalid user ID in token")
	}

	user, err := s.userRepo.GetUserByID(uid)
	if err != nil || user == nil {
		return "", "", errors.NewUnauthorizedError("user not found")
	}

	accessToken, newRefreshToken, err := s.jwtService.RefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}

func (s *authService) GetProfile(userID uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, errors.NewNotFoundError(fmt.Sprintf("user with ID %s not found", userID))
		}
		return nil, errors.NewInternalError(err)
	}
	return user, nil
}

func (s *authService) ValidateToken(token string) (*CustomClaims, error) {
	claims, err := s.jwtService.ValidateToken(token)
	if err != nil {
		return nil, err
	}
	return claims, nil
}
