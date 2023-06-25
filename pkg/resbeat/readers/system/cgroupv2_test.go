package system

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCGroupV2_ReadMemoryUsage(t *testing.T) {
	reader := NewCGroupV2Reader("../../../../tests/fixtures/cgroupv2")

	stat, err := reader.GetMemoryUsageInBytes()

	assert.Nil(t, err)
	assert.Equal(t, uint64(15708160), stat)
}

func TestCGroupV2_ReadMemoryLimit(t *testing.T) {
	reader := NewCGroupV2Reader("../../../../tests/fixtures/cgroupv2")

	stat, err := reader.GetMemoryLimitInBytes()

	assert.Nil(t, err)
	assert.Equal(t, uint64(15728640), stat)
}
