package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/asrma7/meroshare-bot/pkg/logs"
	"github.com/joho/godotenv"
)

type Config struct {
	Environment   string
	Port          string
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	DBConnString  string
	AccessSecret  string
	RefreshSecret string
	TokenExpiry   time.Duration
	RefreshExpiry time.Duration
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		logs.Info("No .env file found, using environment variables", nil)
		logs.Debug("Error loading env file", map[string]interface{}{"error": err})
	}

	return &Config{
		Environment:   getEnv("ENVIRONMENT", "development"),
		Port:          getEnv("PORT", "8080"),
		RedisAddr:     getRedisAddr(),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvInt("REDIS_DB", 0),
		DBConnString:  getDBConnString(),
		AccessSecret:  getEnv("ACCESS_SECRET", "your_access_secret"),
		RefreshSecret: getEnv("REFRESH_SECRET", "your_refresh_secret"),
		TokenExpiry:   time.Minute * 15,
		RefreshExpiry: time.Hour * 24 * 7,
	}
}

func getDBConnString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=prefer",
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "password"),
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_NAME", "meroshare"),
	)
}

func getRedisAddr() string {
	return fmt.Sprintf(
		"%s:%s",
		getEnv("REDIS_HOST", "localhost"),
		getEnv("REDIS_PORT", "6379"),
	)
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
