package telemetry

import (
	"context"
	"fmt"
	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
)

type loggerCtxKey struct {
}

type LogFormats = string

const (
	PlainFormat LogFormats = "text"
	JsonFormat  LogFormats = "json"
)

var logger *zap.Logger

func SetupLogger(ctx context.Context, format LogFormats, level string) (context.Context, *zap.Logger, error) {
	if logger != nil {
		return ctx, logger, nil
	}

	logLevel := zap.InfoLevel

	if level != "" {
		parsedLevel, err := zapcore.ParseLevel(level)

		if err != nil {
			log.Println(
				fmt.Errorf("invalid level, defaulting to INFO: %w", err),
			)
		}

		logLevel = parsedLevel
	}

	var core zapcore.Core

	// TODO: validate the format

	switch format {
	case JsonFormat:
		config := ecszap.NewDefaultEncoderConfig()
		core = ecszap.NewCore(config, os.Stdout, logLevel)
	case PlainFormat:
		config := zap.NewDevelopmentEncoderConfig()
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder := zapcore.NewConsoleEncoder(config)

		core = zapcore.NewCore(encoder, os.Stdout, logLevel)
	}

	logger = zap.New(core, zap.AddCaller())

	return WithContext(ctx, logger), logger, nil
}

func UnsetupLogger(ctx context.Context) context.Context {
	cLogger := FromContext(ctx)

	if cLogger != nil {
		ctx = context.WithValue(ctx, loggerCtxKey{}, nil)
	}

	logger = nil

	return ctx
}

func WithContext(ctx context.Context, logger *zap.Logger) context.Context {
	if lp, ok := ctx.Value(loggerCtxKey{}).(*zap.Logger); ok {
		if lp == logger {
			// Do not store same logger.
			return ctx
		}
	}

	return context.WithValue(ctx, loggerCtxKey{}, logger)
}

func FromContext(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(loggerCtxKey{}).(*zap.Logger); ok {
		return l
	}

	if l := logger; l != nil {
		return l
	}

	return zap.NewNop()
}
