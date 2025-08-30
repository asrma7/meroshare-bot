package routes

import (
	"github.com/asrma7/meroshare-bot/internal/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, authHandler handlers.AuthHandler, accountHandler handlers.AccountHandler, shareHandler handlers.ShareHandler) {
	api := router.Group("/api/v1")

	RegisterAuthRoutes(api, authHandler)
	RegisterAccountRoutes(api, authHandler, accountHandler)
	RegisterShareRoutes(api, authHandler, shareHandler)
}
