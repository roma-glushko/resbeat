package system

import (
	"context"
	"errors"
	"fmt"
	"resbeat/pkg/resbeat/telemetry"
)

var ErrNotAContainer = errors.New("resbeat collects system stats inside a container only")

// SystemStatsReader represents components that reads resource stats from different resource controllers
type SystemStatsReader interface {
	MemoryUsageInBytes() (uint64, error)
	MemoryLimitInBytes() (uint64, error)
	CPUUsageLimitInCores() (float64, error)
	CPUUsageInNanos() (uint64, error)
}

// NewSystemReader
func NewSystemReader(ctx context.Context) (SystemStatsReader, error) {
	logger := telemetry.FromContext(ctx)
	cgroupType, mounts, err := getCGroupMounts(procMountsPath)

	if cgroupType == CGroupUnknown {
		return nil, ErrNotAContainer
	}

	logger.Info(fmt.Sprintf("found system controller: %s", cgroupType))

	if err != nil {
		return nil, fmt.Errorf("failed to init %s controller: %v", cgroupType, err)
	}

	switch cgroupType {
	case CGroupV2:
		return NewCGroupV2Reader(mounts.GetRootDir()), nil
	case CGroupV1:
		return NewCGroupV1Reader(mounts), nil
	}

	panic("cgroup controller should have been processed")
}
