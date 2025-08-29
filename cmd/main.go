package main

import (
	"github.com/asrma7/meroshare-bot/internal/handlers"
	"github.com/asrma7/meroshare-bot/internal/repositories"
	"github.com/asrma7/meroshare-bot/internal/routes"
	"github.com/asrma7/meroshare-bot/internal/services"
	"github.com/asrma7/meroshare-bot/pkg/config"
	"github.com/asrma7/meroshare-bot/pkg/database"
	"github.com/asrma7/meroshare-bot/pkg/logs"
	"github.com/asrma7/meroshare-bot/pkg/redis"
	"github.com/asrma7/meroshare-bot/pkg/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	logs.InitLogger()

	cfg := config.LoadConfig()

	db, err := database.ConnectDB(cfg)
	if err != nil {
		logs.Error("Failed to connect to database", map[string]any{"error": err})
		return
	}

	redisClient := redis.InitRedisClient(cfg)
	if redisClient == nil {
		logs.Error("Failed to connect to Redis", nil)
		return
	}

	if cfg.Environment == "prod" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.Default()
	r.Use(utils.NewCors())

	userRepo := repositories.NewUserRepository(db)
	accountRepo := repositories.NewAccountRepository(db)

	authHandler := handlers.NewAuthHandler(services.NewAuthService(cfg, &userRepo, redisClient))
	accountHandler := handlers.NewAccountHandler(services.NewAccountService(&accountRepo))

	routes.RegisterRoutes(r, authHandler, accountHandler)

	logs.Info("Starting server", map[string]any{
		"port": cfg.Port,
		"env":  cfg.Environment,
	})

	if err := r.Run(":" + cfg.Port); err != nil {
		logs.Logger.Fatal("Failed to start server", map[string]any{
			"error": err,
		})
	}
}
