package readers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCGroupV1_MountsFileParsing(t *testing.T) {
	expectedSubsystem := map[string]string{
		"cpu":     "/sys/fs/cgroup/cpu",
		"cpuacct": "/sys/fs/cgroup/cpuacct",
		"memory":  "/sys/fs/cgroup/memory",
	}

	foundSubsystems, err := getSubsystemsMounts("../../../tests/fixtures/cgroupv1/mounts")

	assert.Nil(t, err)
	assert.NotNil(t, subsystems)
	assert.Equal(t, expectedSubsystem, foundSubsystems.subsystemMounts)
}

func TestCGroupV1_MountsWrongCGroup(t *testing.T) {
	subsystems, err := getSubsystemsMounts("../../../tests/fixtures/cgroupv2/mounts")

	assert.ErrorIs(t, err, CGroupV1NotSupported)
	assert.Nil(t, subsystems)
}

func TestCGroupV1_ReadValues(t *testing.T) {
	mounts := &SubsystemMounts{
		subsystemMounts: map[string]string{
			"cpu":     "../../../tests/fixtures/cgroupv1",
			"cpuacct": "../../../tests/fixtures/cgroupv1",
			"memory":  "../../../tests/fixtures/cgroupv1",
		},
	}

	reader := CGroupV1Reader{subsystemMounts: mounts}
}
