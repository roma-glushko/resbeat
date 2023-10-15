package gpu

import (
	"context"
	"resbeat/pkg/resbeat/telemetry"
)

func NewGPUReader(ctx context.Context) (*GPUReader, error) {
	logger := telemetry.FromContext(ctx)
	var reader GPUReader

	if err := reader.Init(ctx); err != nil {
		return nil, err
	}

	logger.Info("NVML is initialized")

	return &reader, nil
}
