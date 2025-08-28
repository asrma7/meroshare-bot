package routes

import (
	"github.com/asrma7/meroshare-bot/internal/handlers"
	"github.com/asrma7/meroshare-bot/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(router *gin.RouterGroup, authHandler handlers.AuthHandler) {
	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)
	router.GET("/profile", middlewares.AuthMiddleware(authHandler), authHandler.GetProfile)
	router.POST("/refresh", authHandler.RefreshToken)
}
