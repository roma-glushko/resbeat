package system

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCGroupV2_ReadMemoryUsage(t *testing.T) {
	reader := NewCGroupV2Reader("../../../../tests/fixtures/cgroupv2")

	stat, err := reader.MemoryUsageInBytes()

	assert.Nil(t, err)
	assert.Equal(t, uint64(15708160), stat)
}

func TestCGroupV2_ReadMemoryLimit(t *testing.T) {
	reader := NewCGroupV2Reader("../../../../tests/fixtures/cgroupv2")

	stat, err := reader.MemoryLimitInBytes()

	assert.Nil(t, err)
	assert.Equal(t, uint64(15728640), stat)
}

func TestCGroupV2_ReadCPUCores(t *testing.T) {
	reader := NewCGroupV2Reader("../../../../tests/fixtures/cgroupv2")

	stat, err := reader.CPUUsageLimitInCores()

	assert.Nil(t, err)
	assert.Equal(t, 0.01, stat)
}

func TestCGroupV2_ReadCPUUsage(t *testing.T) {
	reader := NewCGroupV2Reader("../../../../tests/fixtures/cgroupv2")

	stat, err := reader.CPUUsageInNanos()

	assert.Nil(t, err)
	assert.Equal(t, uint64(1084617), stat)
}
