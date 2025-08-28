package middlewares

import (
	"net/http"
	"strings"

	"github.com/asrma7/meroshare-bot/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AuthMiddleware(authHandler handlers.AuthHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := authHandler.ValidateToken(tokenString)
		if err != nil {
			if strings.Contains(err.Error(), "token expired") {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			}
			c.Abort()
			return
		}

		if claims.UserID == uuid.Nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID.String())
		c.Next()
	}
}
