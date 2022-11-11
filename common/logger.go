package common

import "go.uber.org/zap"

var Logger *zap.Logger

func InitLogger() {
	Logger, _ = zap.NewProduction()
}

func CleanLogger() {
	Logger.Sync()
}
