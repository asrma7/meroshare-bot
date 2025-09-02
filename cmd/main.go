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
	"github.com/robfig/cron/v3"
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
	shareRepo := repositories.NewShareRepository(db)

	authService := services.NewAuthService(cfg, &userRepo, redisClient)
	userService := services.NewUserService(db)
	accountService := services.NewAccountService(&accountRepo)
	shareService := services.NewShareService(&shareRepo)

	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	accountHandler := handlers.NewAccountHandler(accountService)
	shareHandler := handlers.NewShareHandler(shareService, accountService)

	routes.RegisterRoutes(r, authHandler, userHandler, accountHandler, shareHandler)

	c := cron.New()
	c.AddFunc("0 0 * * *", shareHandler.ApplyShare)
	c.Start()

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
