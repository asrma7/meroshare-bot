package logs

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func InitLogger() {
	Logger = logrus.New()
	Logger.SetLevel(logrus.DebugLevel)
	Logger.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: true,
	})
	Logger.SetOutput(os.Stdout)
}

func Info(message string, fields map[string]interface{}) {
	Logger.WithFields(fields).Info(message)
}

func Debug(message string, fields map[string]interface{}) {
	Logger.WithFields(fields).Debug(message)
}

func Warn(message string, fields map[string]interface{}) {
	Logger.WithFields(fields).Warn(message)
}

func Error(message string, fields map[string]interface{}) {
	Logger.WithFields(fields).Error(message)
}
