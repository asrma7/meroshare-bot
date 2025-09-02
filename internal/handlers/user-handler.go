package handlers

import (
	"net/http"

	"github.com/asrma7/meroshare-bot/internal/services"
	"github.com/asrma7/meroshare-bot/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler interface {
	GetUserDashboard(c *gin.Context)
	ResetUserLogs(c *gin.Context)
}

type userHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) UserHandler {
	return &userHandler{
		userService: userService,
	}
}

func (h *userHandler) GetUserDashboard(c *gin.Context) {
	userID := c.GetString("userID")
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

	dashboardData, err := h.userService.GetUserDashboard(userIDParsed)
	if err != nil {
		errorResp, statusCode := errors.GetErrorResponse(err)
		c.JSON(statusCode, errorResp)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "dashboard": dashboardData})
}

func (h *userHandler) ResetUserLogs(c *gin.Context) {
	userID := c.GetString("userID")
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

	err = h.userService.ResetUserLogs(userIDParsed)
	if err != nil {
		errorResp, statusCode := errors.GetErrorResponse(err)
		c.JSON(statusCode, errorResp)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "User logs reset successfully"})
}
