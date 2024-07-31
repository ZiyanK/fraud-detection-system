package logger

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// CreateLogger created a logger instance
func CreateLogger() *zap.Logger {
	logger := zap.Must(zap.NewProduction())
	if viper.Get("mode") == "dev" {
		logger = zap.Must(zap.NewDevelopment())
	}
	defer logger.Sync()

	return logger
}
