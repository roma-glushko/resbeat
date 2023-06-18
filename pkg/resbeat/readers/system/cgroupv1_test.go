package system

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var expectedSubsystem = map[string]string{
	"cpu":     "/sys/fs/cgroup/cpu",
	"cpuacct": "/sys/fs/cgroup/cpuacct",
	"memory":  "/sys/fs/cgroup/memory",
}

func TestCGroupV1_ParseMounts(t *testing.T) {
	foundSubsystems, err := getSubsystemsMounts("../../../../tests/fixtures/cgroupv1/mounts")

	assert.Nil(t, err)
	assert.NotNil(t, requiredSubsystems)
	assert.Equal(t, expectedSubsystem, foundSubsystems.subsystemMounts)
}

func TestCGroupV1_ParseMountsWithCommaSepSubsystems(t *testing.T) {
	foundSubsystems, err := getSubsystemsMounts("../../../../tests/fixtures/cgroupv1/mounts.2")

	assert.Nil(t, err)
	assert.NotNil(t, requiredSubsystems)
	assert.Equal(t, expectedSubsystem, foundSubsystems.subsystemMounts)
}

func TestCGroupV1_MountsWrongCGroup(t *testing.T) {
	subsystems, err := getSubsystemsMounts("../../../../tests/fixtures/cgroupv2/mounts")

	assert.ErrorIs(t, err, ErrCGroupNotSupported)
	assert.Nil(t, subsystems)
}

func TestCGroupV1_MissingRequiredSubsystems(t *testing.T) {
	subsystems, err := getSubsystemsMounts("../../../../tests/fixtures/cgroupv1/mounts.missedsubsys")

	assert.Nil(t, subsystems)
	assert.ErrorContains(t, err, "missing some of the required subsystems")
	assert.ErrorContains(t, err, "memory")
}

//func TestCGroupV1_ReadValues(t *testing.T) {
//	mounts := &SubsystemMounts{
//		subsystemMounts: map[string]string{
//			"cpu":     "../../../tests/fixtures/cgroupv1",
//			"cpuacct": "../../../tests/fixtures/cgroupv1",
//			"memory":  "../../../tests/fixtures/cgroupv1",
//		},
//	}
//
//	reader := CGroupV1Reader{subsystemMounts: mounts}
//}
