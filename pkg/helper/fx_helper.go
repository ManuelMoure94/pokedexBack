package fxhelper

import (
	"time"

	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Logger() fxevent.Logger {
	appName := "pokedex_backend_go"
	logger := zap.L().WithOptions(zap.IncreaseLevel(zapcore.WarnLevel), zap.WithCaller(false)).Named(appName)
	return &fxevent.ZapLogger{
		Logger: logger,
	}
}

func Timeout() time.Duration {
	return time.Minute
}
