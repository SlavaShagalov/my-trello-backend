package zap

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func TestConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func NewTestLogger() *zap.Logger {
	consoleCfg := TestConfig()

	consoleEncoder := zapcore.NewConsoleEncoder(consoleCfg)
	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel)

	logger := zap.New(consoleCore, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return logger
}
