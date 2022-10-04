package zaplog

import (
	"fmt"
	"go.uber.org/zap"
	"lambda-search-nir/service/application/logger"
)

type Logger struct {
	Logger *zap.Logger
}

func NewZapLogger() logger.Logger {
	logger, _ := zap.NewProduction()
	return &Logger{
		Logger: logger,
	}
}

func (l *Logger) Info(message ...interface{}) {
	str := fmt.Sprintf("%v", message)
	l.Logger.Info(str)
}

func (l *Logger) Error(message ...interface{}) {
	str := fmt.Sprintf("%v", message)
	l.Logger.Error(str)
}

func (l *Logger) Fatal(message ...interface{}) {
	str := fmt.Sprintf("%v", message)
	l.Logger.Fatal(str)
}
