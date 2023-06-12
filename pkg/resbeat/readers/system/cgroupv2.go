package system

import (
	"os"
)

type CGroupV2Reader struct {
}

func (r *CGroupV2Reader) GetMemoryUsageInBytes() (uint64, error) {

}

func (r *CGroupV2Reader) GetMemoryLimitInBytes() (uint64, error) {

}

func (r *CGroupV2Reader) GetCPUUsageLimitInCores() (float64, error) {

}

func (r *CGroupV2Reader) GetCPUUsageInNanos() (uint64, error) {

}

func (r *CGroupV2Reader) GetCPUQuotaInMicros() (uint64, error) {

}

func supportCGroupV2() (bool, error) {
	_, err := os.Stat("/sys/fs/cgroup/cgroup.controllers")

	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func NewCGroupV2Reader() (*CGroupV2Reader, error) {

}
