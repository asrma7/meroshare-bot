package database

import (
	"log"
	"os"
	"time"

	"github.com/asrma7/meroshare-bot/internal/models"
	"github.com/asrma7/meroshare-bot/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDB(cfg *config.Config) (*gorm.DB, error) {
	logLevel := logger.Info
	ignoreRNF := false
	if cfg.Environment == "prod" {
		logLevel = logger.Warn
		ignoreRNF = true
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: ignoreRNF,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(postgres.Open(cfg.DBConnString), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Account{},
		&models.AppliedShare{},
		&models.AppliedShareError{},
	); err != nil {
		return nil, err
	}
	return db, nil
}
