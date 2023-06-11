package resbeat

import (
	"context"
	"fmt"
	"math"
	"resbeat/pkg/resbeat/telemetry"
	"sync"
	"time"
)

// StatsReader represents components that reads resource stats from different resource controllers
type StatsReader interface {
	GetMemoryUsageInBytes() (uint64, error)
	GetMemoryLimitInBytes() (uint64, error)
	GetCPUUsageLimitInCores() (float64, error)
	GetCPUUsageInNanos() (uint64, error)
}

type Monitor struct {
	reader    StatsReader
	mtx       *sync.RWMutex
	prevUsage *Usage
	usage     *Usage
}

func NewMonitor(reader StatsReader) *Monitor {
	return &Monitor{
		reader:    reader,
		mtx:       &sync.RWMutex{},
		prevUsage: nil,
		usage:     nil,
	}
}

func (w *Monitor) Usage() *Usage {
	w.mtx.RLock()
	defer w.mtx.RUnlock()

	return w.usage
}

func (m *Monitor) Run(ctx context.Context, frequency time.Duration) <-chan bool {
	logger := telemetry.FromContext(ctx)
	beat := make(chan bool)

	go func() {
		timer := time.NewTicker(frequency)

		defer func() {
			logger.Info("resource monitor is shutting down")

			timer.Stop()
			close(beat)
		}()

		// init the usage stats on resbeat's startup
		m.usage = m.collectCurrentUsage(ctx)

		for {
			select {
			case <-ctx.Done():
				return
			case <-timer.C:
				usage := m.collectCurrentUsage(ctx)

				m.mtx.Lock()
				m.usage, m.prevUsage = usage, m.usage
				m.mtx.Unlock()

				beat <- true
			}
		}
	}()

	return beat
}

func (m *Monitor) collectCurrentUsage(ctx context.Context) *Usage {
	logger := telemetry.FromContext(ctx)
	cpuUsage, err := m.collectCPUUsage()

	if err != nil {
		logger.Error(fmt.Sprintf("error during collecting CPU stats: %v (usage data will be skipped)", err))
	}

	memoryUsage, err := m.collectMemoryUsage()

	if err != nil {
		logger.Error(fmt.Sprintf("error during collecting memory stats: %v (usage data will be skipped)", err))
	}

	currentUsage := Usage{
		CollectedAt: time.Now().UTC(),
		CPU:         cpuUsage,
		Memory:      memoryUsage,
	}

	return &currentUsage
}

func (m *Monitor) collectMemoryUsage() (*MemoryStats, error) {
	memoryUsageInBytes, err := m.reader.GetMemoryUsageInBytes()

	if err != nil {
		return nil, err
	}

	memoryLimitInBytes, err := m.reader.GetMemoryLimitInBytes()

	if err != nil {
		return nil, err
	}

	return &MemoryStats{
		UsageInBytes:    memoryUsageInBytes,
		LimitInBytes:    memoryLimitInBytes,
		UsagePercentage: float64(memoryUsageInBytes) / float64(memoryLimitInBytes),
	}, nil
}

func (m *Monitor) collectCPUUsage() (*CPUStats, error) {
	prevCPUUsage := m.prevUsage

	var usagePercentage float64
	var usageDelta uint64

	limitInCores, err := m.reader.GetCPUUsageLimitInCores()

	if err != nil {
		return nil, err
	}

	usageInNanos, err := m.reader.GetCPUUsageInNanos()

	if err != nil {
		return nil, err
	}

	if prevCPUUsage == nil {
		usagePercentage = 0.0
		usageDelta = usageInNanos
	} else {
		usageDelta = usageInNanos - prevCPUUsage.CPU.UsageInNanos
		timeDelta := time.Now().UTC().Nanosecond() - prevCPUUsage.CollectedAt.Nanosecond()

		usagePercentage = math.Abs(float64(usageDelta) / float64(timeDelta) / limitInCores)
	}

	return &CPUStats{
		LimitInCors:     limitInCores,
		UsageInNanos:    usageDelta,
		UsagePercentage: usagePercentage,
	}, nil
}
