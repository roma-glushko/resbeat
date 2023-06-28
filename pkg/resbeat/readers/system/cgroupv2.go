package system

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

const (
	cgroupV2RootDir    = "cgroupv2.rootDir"
	memoryUsageInBytes = "memory.current"
	memoryLimitInBytes = "memory.max"

	cpuUsage  = "cpu.stat"
	cpuLimits = "cpu.max"
)

func (s *CgroupMounts) GetRootDir() string {
	return s.subsystemMounts[cgroupV2RootDir]
}

type CGroupV2Reader struct {
	rootDir string
}

func (r *CGroupV2Reader) getStat(statFilePath string) (uint64, error) {
	file, err := os.Open(statFilePath)

	if err != nil {
		return 0, fmt.Errorf("failed to read stat from %s: %v", statFilePath, err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	present := scanner.Scan()

	if !present {
		return 0, fmt.Errorf("stat file %s is empty", statFilePath)
	}

	valueStr := string(bytes.TrimSpace(scanner.Bytes()))

	if valueStr == "max" {
		return 0, nil
	}

	value, err := strconv.ParseUint(valueStr, 10, 64)

	if err != nil {
		return 0, fmt.Errorf("failed to parse stat value from %s: %v", statFilePath, err)
	}

	return value, nil
}

func (r *CGroupV2Reader) getStatMap(statFilePath string) (statMap map[string]uint64, err error) {
	statFile, err := os.Open(statFilePath)

	if err != nil {
		return nil, err
	}

	defer func() {
		fileErr := statFile.Close()

		if fileErr != nil && err == nil {
			err = fileErr
		}
	}()

	statMap = make(map[string]uint64, 10)
	scanner := bufio.NewScanner(statFile)

	for scanner.Scan() {
		fields := bytes.Fields(scanner.Bytes())

		if len(fields) != 2 {
			continue
		}

		name := string(fields[0])
		valueStr := string(bytes.TrimSpace(fields[1]))

		value, parseErr := strconv.ParseUint(valueStr, 10, 64)

		if parseErr != nil {
			return nil, parseErr
		}

		statMap[name] = value
	}

	return statMap, nil
}

func (r *CGroupV2Reader) getCPULimits(statFilePath string) (quota, period uint64, err error) {
	limitFile, err := os.Open(statFilePath)

	if err != nil {
		return 0, 0, errors.Wrapf(err, "failed to read %s file: %v", statFilePath, err)
	}

	content, err := io.ReadAll(limitFile)

	if err != nil {
		return 0, 0, errors.Wrapf(err, "failed to read %s file: %v", statFilePath, err)
	}

	fields := strings.Fields(string(content))

	if len(fields) > 2 || len(fields) == 0 {
		return 0, 0, fmt.Errorf("unexpected format when reading %s file: %s", statFilePath, content)
	}

	if fields[0] == "max" {
		quota = 0
	} else {
		quota, err = strconv.ParseUint(fields[0], 10, 64)

		if err != nil {
			return 0, 0, fmt.Errorf("failed to read CPU quota from %s: %v", statFilePath, err)
		}
	}

	if len(fields) == 2 {
		period, err = strconv.ParseUint(fields[1], 10, 64)

		if err != nil {
			return 0, 0, fmt.Errorf("failed to read CPU period from %s: %v", statFilePath, err)
		}
	}

	return quota, period, nil
}

func (r *CGroupV2Reader) MemoryUsageInBytes() (uint64, error) {
	statFilePath := filepath.Join(r.rootDir, memoryUsageInBytes)

	return r.getStat(statFilePath)
}

func (r *CGroupV2Reader) MemoryLimitInBytes() (uint64, error) {
	statFilePath := filepath.Join(r.rootDir, memoryLimitInBytes)

	return r.getStat(statFilePath)
}

func (r *CGroupV2Reader) CPUUsageLimitInCores() (usage float64, err error) {
	statFilePath := filepath.Join(r.rootDir, cpuLimits)
	quota, period, err := r.getCPULimits(statFilePath)

	if err != nil {
		return 0, err
	}

	if quota == 0 || period == 0 {
		return float64(runtime.NumCPU()), nil
	}

	return float64(quota) / float64(period), nil
}

func (r *CGroupV2Reader) CPUUsageInNanos() (uint64, error) {
	micro_to_nano := uint64(1000)
	statFilePath := filepath.Join(r.rootDir, cpuUsage)
	stats, err := r.getStatMap(statFilePath)

	if err != nil {
		return 0, fmt.Errorf("failed to read CPU usage: %v", err)
	}

	usage, found := stats["usage_usec"] // microseconds

	if !found {
		return 0, fmt.Errorf("did not found CPU usage (usage_usec) in %v", stats)
	}

	return usage * micro_to_nano, nil
}

func NewCGroupV2Reader(rootDir string) *CGroupV2Reader {
	return &CGroupV2Reader{
		rootDir: rootDir,
	}
}
