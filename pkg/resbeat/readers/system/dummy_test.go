package system

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDummySystemReader_GetMemoryStats(t *testing.T) {
	memUsage := uint64(1000000)
	memLim := uint64(2500000)

	reader := DummyStatsReader{
		MemoryUsageInBytes: memUsage,
		MemoryLimitInBytes: memLim,
	}

	stat, err := reader.GetMemoryUsageInBytes()

	assert.Nil(t, err)
	assert.Equal(t, memUsage, stat)

	stat, err = reader.GetMemoryLimitInBytes()

	assert.Nil(t, err)
	assert.Equal(t, memLim, stat)
}

func TestDummySystemReader_GetCPUStats(t *testing.T) {
	usage := uint64(100000000)
	lim := 2.0

	reader := DummyStatsReader{
		CPUUsageInNano:  usage,
		CPULimitInCores: lim,
	}

	stat, err := reader.GetCPUUsageInNanos()

	assert.Nil(t, err)
	assert.Equal(t, usage, stat)

	limStat, err := reader.GetCPUUsageLimitInCores()

	assert.Nil(t, err)
	assert.Equal(t, lim, limStat)
}
