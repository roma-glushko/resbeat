//go:build linux
// +build linux

package gpu

import (
	"context"
	"fmt"
	"github.com/NVIDIA/go-nvml/pkg/nvml"
	"resbeat/pkg/resbeat/telemetry"
)

type GPUReader struct {
}

func (r *GPUReader) Init(ctx context.Context) error {
	logger := telemetry.FromContext(ctx)
	result := nvml.Init()

	if result != nvml.SUCCESS {
		return fmt.Errorf("Unable to initialize NVML: %v", nvml.ErrorString(result))
	}

	return nil
}

func (r *GPUReader) GPUStats(ctx context.Context) (*AllGPUStats, error) {
	logger := telemetry.FromContext(ctx)
	result := nvml.Init()

	if result != nvml.SUCCESS {
		return fmt.Errorf("Unable to initialize NVML: %v", nvml.ErrorString(result))
	}

	count, result := nvml.DeviceGetCount()

	if result != nvml.SUCCESS {
		return 0, fmt.Errorf("Unable to get device count: %v", nvml.ErrorString(result))
	}

	logger.Debug(fmt.Sprintf("Found %v GPU device(s)", count))

	stats := make(map[string]GPUStats, count)

	for i := 0; i < count; i++ {
		device, result := nvml.DeviceGetHandleByIndex(i)

		if result != nvml.SUCCESS {
			return nil, fmt.Errorf("Unable to get device at index %d: %v", i, nvml.ErrorString(result))
		}

		uuid, result := device.GetUUID()

		if result != nvml.SUCCESS {
			return nil, fmt.Errorf("Unable to get uuid of device at index %d: %v", i, nvml.ErrorString(result))
		}

		logger.Debug(fmt.Sprintf("GPU no %v - %v", i, uuid))

		utilization, result := device.GetUtilizationRates()

		if result != nvml.SUCCESS {
			return nil, fmt.Errorf("Unable to get GPU utilization of device %v: %v", uuid, nvml.ErrorString(result))
		}

		memory, result := device.GetMemoryInfo()

		if result != nvml.SUCCESS {
			return nil, fmt.Errorf("Unable to get memory of device %v: %v", uuid, nvml.ErrorString(result))
		}

		stats[uuid] = GPUStats{
			UsagePercentage:    utilization.Gpu,
			MemoryUsedInBytes:  memory.Used,
			TotalMemoryInBytes: memory.Total,
		}
	}

	return &stats, nil
}

func (r *GPUReader) GetGPUCount() (int, error) {
	count, result := nvml.DeviceGetCount()

	if result != nvml.SUCCESS {
		return 0, fmt.Errorf("Unable to get device count: %v", nvml.ErrorString(result))
	}

	return count, nil
}

func (r *GPUReader) Shutdown() error {
	result := nvml.Shutdown()

	if result != nvml.SUCCESS {
		return fmt.Errorf("Unable to shutdown NVML: %v", nvml.ErrorString(result))
	}

	return nil
}
