package readers

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	procMountsPath         string = "/proc/mounts"
	memorySubsystem        string = "memory"
	cpuSubsystem           string = "cpu"
	cpuAccountingSubsystem string = "cpuacct"
)

var subsystems = [...]string{
	memorySubsystem,
	cpuSubsystem,
	cpuAccountingSubsystem,
}

var CGroupV1NotSupported = errors.New("cgroup v1 is not supported")

var mountSplitter = regexp.MustCompile("\\s+")

type SubsystemMounts struct {
	subsystemMounts map[string]string
}

func (s *SubsystemMounts) GetMemoryPath() string {
	return s.subsystemMounts[memorySubsystem]
}

func (s *SubsystemMounts) GetCPUPath() string {
	return s.subsystemMounts[cpuSubsystem]
}

func (s *SubsystemMounts) GetCPUAccountingPath() string {
	return s.subsystemMounts[cpuAccountingSubsystem]
}

func getSubsystemsMounts(mountsPath string) (*SubsystemMounts, error) {
	procMount, err := os.Open(mountsPath)

	if err != nil {
		return nil, err
	}

	defer procMount.Close()

	scanner := bufio.NewScanner(procMount)

	subsystemMounts := map[string]string{}

	for scanner.Scan() {
		mountInfo := mountSplitter.Split(scanner.Text(), -1)
		fsType, mountPath := mountInfo[2], mountInfo[1]

		if !strings.Contains(fsType, "cgroup") {
			continue
		}

		if strings.Contains(fsType, "cgroup2") {
			// we are dealing with newer version of cgroups
			return nil, CGroupV1NotSupported
		}

		pathParts := strings.Split(mountPath, "/")
		subsystem := pathParts[len(pathParts)-1]

		for _, requiredSubsystem := range subsystems {
			if subsystem == requiredSubsystem {
				subsystemMounts[subsystem] = mountPath
				break
			}
		}

	}

	return &SubsystemMounts{
		subsystemMounts: subsystemMounts,
	}, nil
}

type CGroupV1Reader struct {
	subsystemMounts *SubsystemMounts
}

func (r *CGroupV1Reader) getStat(statFilePath string) (uint64, error) {
	// TODO: handle no limit case e.g -1 negative value
	statFile, err := os.Open(statFilePath)

	if err != nil {
		return 0, err
	}

	defer statFile.Close()

	var statRaw []byte

	statRaw, err = io.ReadAll(statFile)

	if err != nil {
		return 0, err
	}

	var statValue uint64

	statValue, err = strconv.ParseUint(string(bytes.TrimSpace(statRaw)), 10, 64)

	if err != nil {
		return 0, err
	}

	return statValue, nil
}

func (r *CGroupV1Reader) GetMemoryUsageInBytes() (uint64, error) {
	return r.getStat(filepath.Join(r.subsystemMounts.GetMemoryPath(), "memory.usage_in_bytes"))
}

func (r *CGroupV1Reader) GetMemoryLimitInBytes() (uint64, error) {
	return r.getStat(filepath.Join(r.subsystemMounts.GetMemoryPath(), "memory.limit_in_bytes"))
}

func (r *CGroupV1Reader) GetCPUUsageInNanos() (uint64, error) {
	return r.getStat(filepath.Join(r.subsystemMounts.GetCPUAccountingPath(), "cpuacct.usage"))
}

func (r *CGroupV1Reader) GetCPUQuotaInMicros() (uint64, error) {
	return r.getStat(filepath.Join(r.subsystemMounts.GetCPUPath(), "cpu.cfs_quota_us"))
}

func (r *CGroupV1Reader) GetCPUPeriodInMicros() (uint64, error) {
	return r.getStat(filepath.Join(r.subsystemMounts.GetCPUPath(), "cpu.cfs_period_us"))
}

func NewCGroupV1Reader() (*CGroupV1Reader, error) {
	subsystemMounts, err := getSubsystemsMounts(procMountsPath)

	if err != nil {
		return nil, err
	}

	return &CGroupV1Reader{
		subsystemMounts: subsystemMounts,
	}, nil
}
