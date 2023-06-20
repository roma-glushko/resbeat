package telemetry

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
	"testing"
)

func TestLogging_CreateLogger(t *testing.T) {
	tests := map[string]struct {
		format LogFormats
		level  string
	}{
		"plain logger with debug level": {
			format: PlainFormat,
			level:  "debug",
		},
		"plain logger with info level": {
			format: PlainFormat,
			level:  "info",
		},
		"plain logger with warning level": {
			format: PlainFormat,
			level:  "warning",
		},
		"plain logger with error level": {
			format: PlainFormat,
			level:  "error",
		},
		"json logger with debug level": {
			format: JsonFormat,
			level:  "debug",
		},
		"json logger with info level": {
			format: JsonFormat,
			level:  "info",
		},
		"json logger with warning level": {
			format: JsonFormat,
			level:  "warning",
		},
		"json logger with error level": {
			format: JsonFormat,
			level:  "error",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			ctx, rootLogger, err := SetupLogger(ctx, test.format, test.level)

			assert.Nil(t, err)
			assert.NotNil(t, rootLogger)
			assert.Equal(t, rootLogger, FromContext(ctx))
			assert.Equal(t, rootLogger, logger)

			UnsetupLogger(ctx)
		})
	}
}

func TestLogging_UnrecognizedLevel(t *testing.T) {
	ctx := context.Background()
	ctx, rootLogger, err := SetupLogger(ctx, PlainFormat, "supercritical")

	assert.Nil(t, err)
	assert.NotNil(t, rootLogger)
	assert.Equal(t, zapcore.InfoLevel, rootLogger.Level())

	UnsetupLogger(ctx)
}

// TODO: cover this case
//func TestLogging_UnrecognizedFormat(t *testing.T) {
//	ctx := context.Background()
//	ctx, rootLogger, err := SetupLogger(ctx, "yaml", "info")
//
//	assert.Nil(t, err)
//	assert.NotNil(t, rootLogger)
//
//	UnsetupLogger(ctx)
//}
