package resbeat

import (
	"context"
	"resbeat/pkg/resbeat/readers"
	"resbeat/pkg/resbeat/telemetry"
	"sync"
	"time"
)

type Monitor struct {
	reader    readers.StatsReader
	mtx       *sync.RWMutex
	prevUsage *Usage
	usage     *Usage
}

func NewMonitor(reader readers.StatsReader) *Monitor {
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
			logger.Info("monitor is shutting down")

			timer.Stop()
			close(beat)
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case <-timer.C:
				m.mtx.Lock()
				m.usage, m.prevUsage = m.collectCurrentUsage(), m.usage
				m.mtx.Unlock()

				beat <- true
			}
		}
	}()

	return beat
}

func (m *Monitor) collectCurrentUsage() *Usage {
	memoryUsageInBytes, _ := m.reader.GetMemoryUsageInBytes() // TODO: handle errors
	memoryLimitInBytes, _ := m.reader.GetMemoryLimitInBytes()

	currentUsage := Usage{
		CPU: &CPUStats{},
		Memory: &MemoryStats{
			UsageInBytes:    memoryUsageInBytes,
			LimitInBytes:    memoryLimitInBytes,
			UsagePercentage: float64(memoryUsageInBytes) / float64(memoryLimitInBytes),
		},
	}

	return &currentUsage
}
