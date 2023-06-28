package resbeat

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"math"
	"reflect"
	"resbeat/pkg/resbeat/readers/system"
	"resbeat/pkg/resbeat/telemetry"
	"sync"
	"time"
)

type Monitor struct {
	reader    system.SystemStatsReader
	mu        *sync.RWMutex
	prevUsage *Usage
	usage     *Usage
	wg        *sync.WaitGroup
}

func NewMonitor(reader system.SystemStatsReader) *Monitor {
	return &Monitor{
		reader:    reader,
		mu:        &sync.RWMutex{},
		prevUsage: nil,
		usage:     nil,
		wg:        &sync.WaitGroup{},
	}
}

func (w *Monitor) Usage() *Usage {
	w.mu.RLock()
	defer w.mu.RUnlock()

	return w.usage
}

func (m *Monitor) Run(ctx context.Context, frequency time.Duration) <-chan bool {
	logger := telemetry.FromContext(ctx)
	beat := make(chan bool)

	m.wg.Add(1)

	go func() {
		timer := time.NewTicker(frequency)

		defer func() {
			logger.Info("resource monitor is shutting down")

			timer.Stop()
			close(beat)
			m.wg.Done()
		}()

		// init the usage stats on resbeat's startup
		m.usage = m.collectCurrentUsage(ctx)

		for {
			select {
			case <-ctx.Done():
				return
			case <-timer.C:
				m.collectUsageOnTick(ctx)
				beat <- true
			}
		}
	}()

	return beat
}

func (m *Monitor) collectUsageOnTick(ctx context.Context) {
	m.prevUsage = m.usage
	usage := m.collectCurrentUsage(ctx)

	m.mu.Lock()
	m.usage = usage
	m.mu.Unlock()
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

	if systemReader == nil || reflect.ValueOf(systemReader).IsNil() {
		// the system reader was not inited successfully
		return nil
	}

	cpuUsage, err := m.collectCPUUsage(ctx)

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

func (m *Monitor) clampPercentage(value float64) float64 {
	minRange, maxRange := 0.0, 1.0

	return math.Max(minRange, math.Min(value, maxRange))
}

func (m *Monitor) collectMemoryUsage() (*MemoryStats, error) {
	if m.reader == nil {
		return nil, nil
	}

	systemReader := m.reader

	memoryUsageInBytes, err := systemReader.MemoryUsageInBytes()

	if err != nil {
		return nil, err
	}

	memoryLimitInBytes, err := systemReader.MemoryLimitInBytes()

	if err != nil {
		return nil, err
	}

	return &MemoryStats{
		UsageInBytes:    memoryUsageInBytes,
		LimitInBytes:    memoryLimitInBytes,
		UsagePercentage: m.clampPercentage(float64(memoryUsageInBytes) / float64(memoryLimitInBytes)),
	}, nil
}

func (m *Monitor) collectCPUUsage(ctx context.Context) (*CPUStats, error) {
	logger := telemetry.FromContext(ctx)

	if m.reader == nil {
		return nil, nil
	}

	systemReader := m.reader
	prevUsage := m.prevUsage

	var usagePercentage float64
	var usageDelta uint64
	var timeDelta int64

	limitInCores, err := systemReader.CPUUsageLimitInCores()

	if err != nil {
		return nil, err
	}

	collectedAt := time.Now().UTC()
	accumulatedUsageInNanos, err := systemReader.CPUUsageInNanos()

	if err != nil {
		return nil, err
	}

	if prevUsage == nil {
		logger.Debug("no previous CPU usage report")
		usagePercentage = 0.0
		usageDelta = accumulatedUsageInNanos
	} else {
		prevCPUUsage := prevUsage.System.CPU
		usageDelta = accumulatedUsageInNanos - prevCPUUsage.AccumulatedUsageInNanos()
		timeDelta = collectedAt.Sub(prevCPUUsage.CollectedAt()).Nanoseconds()

		usagePercentage = float64(usageDelta) / float64(timeDelta) / limitInCores
		// 3-06-28T13:37:31.770Z        DEBUG   resbeat/monitor.go:191  CPU report      {"limitInCores": 0.5, "accumulatedUsage": 29133460000, "usageDelta": 1503878000, "timeDelta": 3000259608, "usagePercentage": 1.0024985811161178}
	}

	logger.Debug(
		"CPU report",
		zap.Float64("limitInCores", limitInCores),
		zap.Uint64("accumulatedUsage", accumulatedUsageInNanos),
		zap.Uint64("usageDelta", usageDelta),
		zap.Int64("timeDelta", timeDelta),
		zap.Float64("usagePercentage", usagePercentage),
	)

	return &CPUStats{
		collectedAt:             collectedAt,
		accumulatedUsageInNanos: accumulatedUsageInNanos,
		LimitInCores:            limitInCores,
		UsageInNanos:            usageDelta,
		UsagePercentage:         m.clampPercentage(usagePercentage),
	}, nil
}
