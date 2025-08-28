package handlers

import (
	"net/http"

	"github.com/asrma7/meroshare-bot/internal/requests"
	"github.com/asrma7/meroshare-bot/internal/services"
	"github.com/asrma7/meroshare-bot/pkg/errors"
	"github.com/asrma7/meroshare-bot/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	RefreshToken(c *gin.Context)
	GetProfile(c *gin.Context)
	ValidateToken(token string) (*services.CustomClaims, error)
}

type authHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) AuthHandler {
	return &authHandler{
		authService: authService,
	}
}

func (h *authHandler) Register(c *gin.Context) {
	var req requests.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errResp := errors.ErrorResponse{
			Type:    "VALIDATION_ERROR",
			Message: "Invalid request data",
			Details: map[string]string{"error": err.Error()},
		}
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	userID, err := h.authService.RegisterUser(req.Username, req.Password, req.Email, utils.CapitalizeFirstLetter(req.FirstName), utils.CapitalizeFirstLetter(req.LastName))
	if err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			errorResp := errors.ErrorResponse{
				Type:    appErr.Type,
				Message: appErr.Message,
			}
			if appErr.Details != nil {
				errorResp.Details = appErr.Details
			}
			statusCode := errors.GetErrorStatusCode(err)
			c.JSON(statusCode, errorResp)
		} else {
			errResp := errors.ErrorResponse{
				Type:    "INTERNAL_ERROR",
				Message: "An unexpected error occurred",
			}
			c.JSON(http.StatusInternalServerError, errResp)
		}
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "success", "user_id": userID})
}

func (h *authHandler) Login(c *gin.Context) {
	var req requests.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errResp := errors.ErrorResponse{
			Type:    "VALIDATION_ERROR",
			Message: "Invalid request data",
			Details: map[string]string{"error": err.Error()},
		}
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	accessToken, refreshToken, err := h.authService.LoginUser(req.Identifier, req.Password)
	if err != nil {
		errorResp, statusCode := errors.GetErrorResponse(err)
		c.JSON(statusCode, errorResp)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":        "success",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *authHandler) RefreshToken(c *gin.Context) {
	var req requests.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errResp := errors.ErrorResponse{
			Type:    "VALIDATION_ERROR",
			Message: "Invalid request data",
			Details: map[string]string{"error": err.Error()},
		}
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	accessToken, refreshToken, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		errorResp, statusCode := errors.GetErrorResponse(err)
		c.JSON(statusCode, errorResp)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":        "success",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *authHandler) GetProfile(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	if userID == "" {
		errResp := errors.ErrorResponse{
			Type:    "UNAUTHORIZED",
			Message: "User ID not found in context",
		}
		c.JSON(http.StatusUnauthorized, errResp)
		return
	}
	userIDParsed, err := uuid.Parse(userID)
	if err != nil {
		errResp := errors.ErrorResponse{
			Type:    "VALIDATION_ERROR",
			Message: "Invalid user ID format",
		}
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	user, err := h.authService.GetProfile(userIDParsed)
	if err != nil {
		errorResp, statusCode := errors.GetErrorResponse(err)
		c.JSON(statusCode, errorResp)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "profile": gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
	}})
}

func (h *authHandler) ValidateToken(token string) (*services.CustomClaims, error) {
	userID, err := h.authService.ValidateToken(token)
	if err != nil {
		return nil, err
	}
	return userID, nil
}
