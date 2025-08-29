package routes

import (
	"github.com/asrma7/meroshare-bot/internal/handlers"
	"github.com/asrma7/meroshare-bot/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterAccountRoutes(router *gin.RouterGroup, authHandler handlers.AuthHandler, accountHandler handlers.AccountHandler) {
	router.Use(middlewares.AuthMiddleware(authHandler))
	router.POST("/accounts", accountHandler.CreateAccount)
}
