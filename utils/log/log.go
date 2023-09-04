package log

import (
	"go.uber.org/zap"
	//"gorm.io/gorm/logger"
	//"gorm.io/gorm/logger"
)

var logger *zap.Logger

func Init() {
	logger, _ = zap.NewProduction()

	logger.Info("log init success", zap.String("field", "zap log"))
}

func FieldLog(fieldName string, level string, msg string) {
	switch level {
	case "info":
		logger.Info(msg, zap.String("field", fieldName))
	case "error":
		logger.Error(msg, zap.String("field", fieldName))
	case "panic":
		logger.Panic(msg, zap.String("field", fieldName))
	}
}
