package logger

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	once sync.Once
	cfg  zap.Config
)

func init() {
	once.Do(initLoggerConfig)
}

func initLoggerConfig() {
	level, err := zap.ParseAtomicLevel("debug")
	if err != nil {
		panic(fmt.Errorf("failed to parse log level: %v", err))
	}

	cfg = zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.Level = level
}

func NewLogger(name string, opts ...zap.Option) *zap.Logger {
	logger, err := cfg.Build(opts...)
	if err != nil {
		panic(fmt.Errorf("failed to build logger: %v", err))
	}

	if name == "" {
		return logger
	}

	return logger.Named(name)
}
