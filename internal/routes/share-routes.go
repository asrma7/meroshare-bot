package routes

import (
	"github.com/asrma7/meroshare-bot/internal/handlers"
	"github.com/asrma7/meroshare-bot/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterShareRoutes(r *gin.RouterGroup, authHandler handlers.AuthHandler, shareHandler handlers.ShareHandler) {
	r.Use(middlewares.AuthMiddleware(authHandler))
	r.GET("/shares/applied", shareHandler.GetAppliedShares)
	r.GET("/shares/errors", shareHandler.GetAppliedShareErrors)
}
