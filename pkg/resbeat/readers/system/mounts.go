package system

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type CGroup string

const (
	procMountsPath string = "/proc/mounts"
	CGroupUnknown  CGroup = "unknown"
	CGroupV1       CGroup = "cgroupv1"
	CGroupV2       CGroup = "cgroupv2"
)

type CgroupMounts struct {
	subsystemMounts map[string]string
}

func getCGroupMounts(mountsPath string) (CGroup, *CgroupMounts, error) {
	procMount, err := os.Open(mountsPath)

	if err != nil {
		return CGroupUnknown, nil, fmt.Errorf("failed to read %s: %v", mountsPath, err)
	}

	scanner := bufio.NewScanner(procMount)

	cgroupMounts := make(map[string]string, 10)

	for scanner.Scan() {
		mountInfo := mountSplitter.Split(scanner.Text(), -1)

		if len(mountInfo) < 6 {
			// a broken line, skipping it
			continue
		}

		fsType, mountPath := mountInfo[2], mountInfo[1]

		if !strings.Contains(fsType, "cgroup") {
			continue
		}

		if strings.Contains(fsType, "cgroup2") {
			// we are dealing with newer version of cgroups
			cgroupMounts[cgroupV2RootDir] = mountPath

			return CGroupV2, &CgroupMounts{
				subsystemMounts: cgroupMounts,
			}, procMount.Close()
		}

		pathParts := strings.Split(mountPath, "/")

		subsystemMountParts := pathParts[:len(pathParts)-1]
		subsystemNames := strings.Split(pathParts[len(pathParts)-1], ",")
		mountPartsLen := len(subsystemMountParts)

		for _, subsystem := range subsystemNames {
			for _, requiredSubsystem := range requiredSubsystems {
				if subsystem == requiredSubsystem {
					subsystemPath := make([]string, mountPartsLen, mountPartsLen+1)
					copy(subsystemPath, subsystemMountParts)
					subsystemPath = append(subsystemPath, subsystem)
					subsystemPath[0] = string(filepath.Separator)

					cgroupMounts[subsystem] = filepath.Join(subsystemPath...)

					break
				}
			}
		}

	}

	missedSubsystems := make([]string, 0, len(requiredSubsystems))

	for _, reqSubsystem := range requiredSubsystems {
		if _, found := cgroupMounts[reqSubsystem]; !found {
			missedSubsystems = append(missedSubsystems, reqSubsystem)
		}
	}

	if len(missedSubsystems) == len(requiredSubsystems) {
		return CGroupUnknown, nil, fmt.Errorf("could not found any cgroup mounts")
	}

	if len(missedSubsystems) > 0 {
		return CGroupV1, nil, fmt.Errorf("missing some of the required cgroupv1 subsystems: %q", missedSubsystems)
	}

	return CGroupV1, &CgroupMounts{
		subsystemMounts: cgroupMounts,
	}, procMount.Close()
}
