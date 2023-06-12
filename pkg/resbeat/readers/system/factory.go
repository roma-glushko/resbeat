package system

import (
	"errors"
	"fmt"
)

var ErrNotAContainer = errors.New("resbeat collects system stats inside a container only")

// SystemStatsReader represents components that reads resource stats from different resource controllers
type SystemStatsReader interface {
	GetMemoryUsageInBytes() (uint64, error)
	GetMemoryLimitInBytes() (uint64, error)
	GetCPUUsageLimitInCores() (float64, error)
	GetCPUUsageInNanos() (uint64, error)
}

// NewSystemReader
func NewSystemReader() (SystemStatsReader, error) {
	var reader SystemStatsReader

	reader, err := NewCGroupV2Reader()

	if err == nil {
		return reader, nil
	}

	if err != nil && !errors.Is(err, ErrCGroupNotSupported) {
		return nil, fmt.Errorf("faild to read cgroupv2 controller: %w", err)
	}

	reader, err = NewCGroupV1Reader()

	if err == nil {
		return reader, nil
	}

	if err != nil && errors.Is(err, ErrCGroupNotSupported) {
		return nil, ErrNotAContainer
	}

	return nil, fmt.Errorf("faild to read cgroupv1 controller: %w", err)
}
