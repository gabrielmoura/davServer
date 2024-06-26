package log

import "go.uber.org/zap"

var Logger *zap.Logger

func InitLogger() {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	Logger = logger
}
