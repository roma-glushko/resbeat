package system

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCGroupV1_ReadValues(t *testing.T) {
	mounts := &CgroupMounts{
		subsystemMounts: map[string]string{
			"cpu":     "../../../../tests/fixtures/cgroupv1",
			"cpuacct": "../../../../tests/fixtures/cgroupv1",
			"memory":  "../../../../tests/fixtures/cgroupv1",
		},
	}

	reader := NewCGroupV1Reader(mounts)
	stat, err := reader.GetMemoryLimitInBytes()

	assert.Nil(t, err)
	assert.Equal(t, uint64(268435456), stat)

	stat, err = reader.GetMemoryUsageInBytes()

	assert.Nil(t, err)
	assert.Equal(t, uint64(1998848), stat)
}
