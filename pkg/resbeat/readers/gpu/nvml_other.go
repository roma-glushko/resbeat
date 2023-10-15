//go:build !linux
// +build !linux

package gpu

import (
	"context"
	"errors"
)

var ERR_GPU_NOT_SUPPORTED = errors.New("GPU monitoring is supported only on Linux (with the NVML toolchain installed)")

type GPUReader struct {
}

func (*GPUReader) Init(ctx context.Context) error {
	return ERR_GPU_NOT_SUPPORTED
}

func (*GPUReader) GPUStats() (*AllGPUStats, error) {
	return nil, ERR_GPU_NOT_SUPPORTED
}

func (*GPUReader) Shutdown() error {
	return ERR_GPU_NOT_SUPPORTED
}
