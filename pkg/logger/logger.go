package logger

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

func init() {
	// 初始化 Zap 日志
	var err error
	Logger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
}

// Sync 刷新日志缓冲区
func Sync() {
	_ = Logger.Sync()
}