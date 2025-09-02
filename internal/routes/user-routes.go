package routes

import (
	"github.com/asrma7/meroshare-bot/internal/handlers"
	"github.com/asrma7/meroshare-bot/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.RouterGroup, authHandler handlers.AuthHandler, userHandler handlers.UserHandler) {
	r.Use(middlewares.AuthMiddleware(authHandler))
	r.GET("/dashboard", userHandler.GetUserDashboard)
	r.POST("/reset-logs", userHandler.ResetUserLogs)
}
