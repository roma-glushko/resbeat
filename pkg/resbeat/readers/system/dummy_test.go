package system

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDummySystemReader_GetMemoryStats(t *testing.T) {
	memUsage := uint64(1000000)
	memLim := uint64(2500000)

	reader := DummyStatsReader{
		memoryUsageInBytes: memUsage,
		memoryLimitInBytes: memLim,
	}

	stat, err := reader.MemoryUsageInBytes()

	assert.Nil(t, err)
	assert.Equal(t, memUsage, stat)

	stat, err = reader.MemoryLimitInBytes()

	assert.Nil(t, err)
	assert.Equal(t, memLim, stat)
}

func TestDummySystemReader_GetCPUStats(t *testing.T) {
	usage := uint64(100000000)
	lim := 2.0

	reader := DummyStatsReader{
		cpuUsageInNano:  usage,
		cpuLimitInCores: lim,
	}

	stat, err := reader.CPUUsageInNanos()

	assert.Nil(t, err)
	assert.Equal(t, usage, stat)

	limStat, err := reader.CPUUsageLimitInCores()

	assert.Nil(t, err)
	assert.Equal(t, lim, limStat)
}
