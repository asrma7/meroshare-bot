package redis

import (
	"context"

	"github.com/redis/go-redis/v9"

	"github.com/asrma7/meroshare-bot/pkg/config"
	"github.com/asrma7/meroshare-bot/pkg/logs"
)

var ctx = context.Background()

func InitRedisClient(cfg *config.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		logs.Error("Failed to connect to Redis", map[string]interface{}{"error": err})
		return nil
	}

	logs.Info("Connected to Redis successfully", nil)
	return rdb
}
