package resbeat

import (
	"context"
	"github.com/stretchr/testify/assert"
	"resbeat/pkg/resbeat/readers/system"
	"testing"
	"time"
)

func TestMonitor_ExitOnCanceledCtx(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	reader := system.NewDummyStatsReader(nil, nil, nil, nil)
	monitor := NewMonitor(reader)

	updateC := monitor.Run(ctx, 10*time.Millisecond)
	<-updateC
	cancel()

	monitor.wg.Wait()
	assert.NotNil(t, monitor.Usage())
}

func TestMonitor_CPUUsageReport(t *testing.T) {
	ctx := context.Background()
	usageInNano1, limInCores := uint64(149537069), 2.0

	reader := system.NewDummyStatsReader(
		nil,
		nil,
		&usageInNano1,
		&limInCores,
	)
	monitor := NewMonitor(reader)

	monitor.collectUsageOnTick(ctx)
	cpuUsage := monitor.Usage().System.CPU

	assert.NotNil(t, cpuUsage)
	assert.NotNil(t, cpuUsage.collectedAt)
	assert.Equal(t, usageInNano1, cpuUsage.AccumulatedUsageInNanos())
	assert.Equal(t, 0.0, cpuUsage.UsagePercentage)
	assert.Equal(t, limInCores, cpuUsage.LimitInCores)

	usageInNano2 := uint64(149547069)
	reader.SetCPUUsageInNanos(usageInNano2)

	// recollect stats again
	monitor.collectUsageOnTick(ctx)
	prevUsage := cpuUsage
	cpuUsage = monitor.Usage().System.CPU

	assert.NotNil(t, cpuUsage)
	assert.NotNil(t, cpuUsage.collectedAt)
	assert.Greater(t, cpuUsage.collectedAt, prevUsage.collectedAt)
	assert.Equal(t, usageInNano2, cpuUsage.AccumulatedUsageInNanos())
	assert.Equal(t, limInCores, cpuUsage.LimitInCores)
	assert.Equal(t, usageInNano2-usageInNano1, cpuUsage.UsageInNanos)

	// assert.(t, , cpuUsage.UsagePercentage)
}

func TestMonitor_ErrorToInitSystemReader(t *testing.T) {
	ctx := context.Background()
	var noReader *system.DummyStatsReader // we would get nil on system reader init error
	monitor := NewMonitor(noReader)

	monitor.collectUsageOnTick(ctx)
	usage := monitor.Usage()

	assert.NotNil(t, usage.CollectedAt)
	assert.Nil(t, usage.System)
}

func BenchmarkMonitor_CGroupV2UsageCollection(b *testing.B) {
	ctx := context.Background()
	reader := system.NewCGroupV2Reader("../../../../tests/fixtures/cgroupv2")
	monitor := NewMonitor(reader)

	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		monitor.collectUsageOnTick(ctx)
	}
}
