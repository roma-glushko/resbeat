//go:build linux
// +build linux

package gpu

import (
	"context"
	"fmt"
	"github.com/NVIDIA/go-nvml/pkg/nvml"
)

type GPUReader struct {
}

func (r *GPUReader) Init(ctx context.Context) error {
	logger := telemetry.FromContext(ctx)
	result := nvml.Init()

	if result != nvml.SUCCESS {
		return fmt.Errorf("Unable to initialize NVML: %v", nvml.ErrorString(result))
	}

	count, err := r.GetGPUCount()

	if err != nil {
		return err
	}

	logger.Debug(fmt.Sprintf("Found %v GPU device(s)", count))

	return nil
}

func (r *GPUReader) GPUStats() (*AllGPUStats, error) {
	stats := make(map[string]GPUStats, count)

	count, err := r.GetGPUCount()

	if err != nil {
		return nil, err
	}

	for i := 0; i < count; i++ {
		device, result := nvml.DeviceGetHandleByIndex(i)

		if result != nvml.SUCCESS {
			return nil, fmt.Errorf("Unable to get device at index %d: %v", i, nvml.ErrorString(result))
		}

		uuid, result := device.GetUUID()

		if result != nvml.SUCCESS {
			return nil, fmt.Errorf("Unable to get uuid of device at index %d: %v", i, nvml.ErrorString(result))
		}

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

func (r *GPUReader) GetGPUCount() (int, err) {
	count, result := nvml.DeviceGetCount()

	if result != nvml.SUCCESS {
		return nil, fmt.Errorf("Unable to get device count: %v", nvml.ErrorString(result))
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
