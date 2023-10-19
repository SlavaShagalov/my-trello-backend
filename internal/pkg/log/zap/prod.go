package zap

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func ProdConfig() zapcore.EncoderConfig {
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

func NewProdLogger(filePath string) (*zap.Logger, *os.File, error) {
	cfg := ProdConfig()

	fileEncoder := zapcore.NewConsoleEncoder(cfg)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, nil, err
	}

	fileCore := zapcore.NewCore(fileEncoder, zapcore.Lock(zapcore.AddSync(file)), zapcore.DebugLevel)

	logger := zap.New(fileCore, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return logger, file, nil
}
