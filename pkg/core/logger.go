package core

import (
	"go.uber.org/zap"
)

var logger, _ = zap.NewDevelopment()

var Log = NewLogger("Core")

func NewLogger(name string) *zap.Logger {
	return logger.Named(name)
}
