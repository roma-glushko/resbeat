package system

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

const (
	memorySubsystem        string = "memory"
	cpuSubsystem           string = "cpu"
	cpuAccountingSubsystem string = "cpuacct"
)

var requiredSubsystems = [...]string{
	memorySubsystem,
	cpuSubsystem,
	cpuAccountingSubsystem,
}

var mountSplitter = regexp.MustCompile(`\s+`)

func (s *CgroupMounts) GetMemoryPath() string {
	return s.subsystemMounts[memorySubsystem]
}

func (s *CgroupMounts) GetCPUPath() string {
	return s.subsystemMounts[cpuSubsystem]
}

func (s *CgroupMounts) GetCPUAccountingPath() string {
	return s.subsystemMounts[cpuAccountingSubsystem]
}

type CGroupV1Reader struct {
	subsystemMounts *CgroupMounts
}

func (r *CGroupV1Reader) getStat(statFilePath string) (uint64, error) {
	// TODO: handle no limit case e.g -1 negative value
	statFile, err := os.Open(statFilePath)

	if err != nil {
		return 0, err
	}

	var statRaw []byte

	statRaw, err = io.ReadAll(statFile)

	if err != nil {
		_ = statFile.Close()
		return 0, err
	}

	var statValue uint64

	statValue, err = strconv.ParseUint(string(bytes.TrimSpace(statRaw)), 10, 64)

	if err != nil {
		_ = statFile.Close()
		return 0, err
	}

	return statValue, statFile.Close()
}

func (r *CGroupV1Reader) MemoryUsageInBytes() (uint64, error) {
	return r.getStat(filepath.Join(r.subsystemMounts.GetMemoryPath(), "memory.usage_in_bytes"))
}

func (r *CGroupV1Reader) MemoryLimitInBytes() (uint64, error) {
	return r.getStat(filepath.Join(r.subsystemMounts.GetMemoryPath(), "memory.limit_in_bytes"))
}

func (r *CGroupV1Reader) CPUUsageLimitInCores() (float64, error) {
	cpuQuota, err := r.getCPUQuotaInMicros()

	if err != nil {
		return 0, err
	}

	cpuPeriod, err := r.getCPUPeriodInMicros()

	if err != nil {
		return 0, err
	}

	return float64(cpuQuota) / float64(cpuPeriod), nil
}

func (r *CGroupV1Reader) CPUUsageInNanos() (uint64, error) {
	return r.getStat(filepath.Join(r.subsystemMounts.GetCPUAccountingPath(), "cpuacct.usage"))
}

func (r *CGroupV1Reader) getCPUQuotaInMicros() (uint64, error) {
	return r.getStat(filepath.Join(r.subsystemMounts.GetCPUPath(), "cpu.cfs_quota_us"))
}

func (r *CGroupV1Reader) getCPUPeriodInMicros() (uint64, error) {
	return r.getStat(filepath.Join(r.subsystemMounts.GetCPUPath(), "cpu.cfs_period_us"))
}

func NewCGroupV1Reader(cgroupMounts *CgroupMounts) *CGroupV1Reader {
	return &CGroupV1Reader{
		subsystemMounts: cgroupMounts,
	}
}
