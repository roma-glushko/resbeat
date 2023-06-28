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
	cgroupType, foundMounts, err := getCGroupMounts("../../../../tests/fixtures/cgroupv1/mounts")

	assert.Nil(t, err)
	assert.Equal(t, CGroupV1, cgroupType)
	assert.NotNil(t, foundMounts)
	assert.Equal(t, expectedSubsystem, foundMounts.subsystemMounts)
}

func TestCGroupV1_ParseMountsWithCommaSepSubsystems(t *testing.T) {
	cgroupType, foundMounts, err := getCGroupMounts("../../../../tests/fixtures/cgroupv1/mounts.2")

	assert.Nil(t, err)
	assert.Equal(t, CGroupV1, cgroupType)
	assert.NotNil(t, requiredSubsystems)
	assert.Equal(t, expectedSubsystem, foundMounts.subsystemMounts)
}

func TestCGroupV1_MountsWrongCGroup(t *testing.T) {
	cgroupType, foundMounts, err := getCGroupMounts("../../../../tests/fixtures/cgroupv2/mounts")

	assert.Equal(t, CGroupV2, cgroupType)
	assert.Equal(t, "/sys/fs/cgroup", foundMounts.GetRootDir())
	assert.Nil(t, err)
}

func TestCGroupV1_MissingRequiredSubsystems(t *testing.T) {
	cgroupType, foundMounts, err := getCGroupMounts("../../../../tests/fixtures/cgroupv1/mounts.missedsubsys")

	assert.Equal(t, CGroupV1, cgroupType)
	assert.Nil(t, foundMounts)
	assert.ErrorContains(t, err, "missing some of the required cgroupv1 subsystems")
	assert.ErrorContains(t, err, "memory")
}

func TestCGroupV1_MountsFileNotFound(t *testing.T) {
	cgroupType, foundMounts, err := getCGroupMounts("../../../../tests/fixtures/cgroupv3.notfound/mounts")

	assert.Equal(t, CGroupUnknown, cgroupType)
	assert.Nil(t, foundMounts)
	assert.ErrorContains(t, err, "failed to read")
}
