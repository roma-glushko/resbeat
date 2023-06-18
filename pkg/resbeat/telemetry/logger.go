package telemetry

import (
	"context"
	"go.uber.org/zap"
)

type contextKey string

const loggerKey = contextKey("telemetry")

var rootLogger *zap.Logger

func SetupLogger(ctx context.Context) (*zap.Logger, error) {
	config := zap.NewDevelopmentConfig()
	config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	config.OutputPaths = []string{"stdout"}

	logger, err := config.Build(
		zap.AddStacktrace(zap.ErrorLevel),
	)

	if err != nil {
		return nil, err
	}

	logger = logger.Named("resbeat")
	rootLogger = logger

	return logger, nil
}

func WithContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func FromContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return rootLogger
	}

	if logger, ok := ctx.Value(loggerKey).(*zap.Logger); ok {
		return logger
	}

	return rootLogger
}
