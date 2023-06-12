package resbeat

import (
	"context"
	"fmt"
	"resbeat/pkg/resbeat/readers/system"
	"resbeat/pkg/resbeat/telemetry"
	"sync"
	"time"
)

type Monitor struct {
	reader    *system.SystemStatsReader
	mtx       *sync.RWMutex
	prevUsage *Usage
	usage     *Usage
}

func NewMonitor(reader *system.SystemStatsReader) *Monitor {
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
	currentUsage := Usage{
		CollectedAt: time.Now().UTC(),
		System:      m.collectSystemUsage(ctx),
	}

	return &currentUsage
}

func (m *Monitor) collectSystemUsage(ctx context.Context) *SystemStats {
	logger := telemetry.FromContext(ctx)
	systemReader := m.reader

	if systemReader == nil {
		// the system reader was not inited successfully
		return nil
	}

	cpuUsage, err := m.collectCPUUsage()

	if err != nil {
		logger.Error(fmt.Sprintf("error during collecting CPU stats: %v (usage data will be skipped)", err))
	}

	memoryUsage, err := m.collectMemoryUsage()

	if err != nil {
		logger.Error(fmt.Sprintf("error during collecting memory stats: %v (usage data will be skipped)", err))
	}

	return &SystemStats{
		CPU:    cpuUsage,
		Memory: memoryUsage,
	}
}

func (m *Monitor) collectMemoryUsage() (*MemoryStats, error) {
	if m.reader == nil {
		return nil, nil
	}

	systemReader := *m.reader

	memoryUsageInBytes, err := systemReader.GetMemoryUsageInBytes()

	if err != nil {
		return nil, err
	}

	memoryLimitInBytes, err := systemReader.GetMemoryLimitInBytes()

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
	if m.reader == nil {
		return nil, nil
	}

	systemReader := *m.reader
	prevCPUUsage := m.prevUsage

	var usagePercentage float64
	var usageDelta uint64

	limitInCores, err := systemReader.GetCPUUsageLimitInCores()

	if err != nil {
		return nil, err
	}

	usageInNanos, err := systemReader.GetCPUUsageInNanos()

	if err != nil {
		return nil, err
	}

	if prevCPUUsage == nil {
		usagePercentage = 0.0
		usageDelta = usageInNanos
	} else {
		usageDelta = usageInNanos - prevCPUUsage.System.CPU.UsageInNanos
		timeDelta := time.Now().UTC().Nanosecond() - prevCPUUsage.CollectedAt.Nanosecond()

		usagePercentage = float64(usageDelta) / float64(timeDelta) / limitInCores / 100.0
	}

	return &CPUStats{
		LimitInCores:    limitInCores,
		UsageInNanos:    usageDelta,
		UsagePercentage: usagePercentage,
	}, nil
}
