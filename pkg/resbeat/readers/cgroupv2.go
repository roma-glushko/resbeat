package readers

import (
	"errors"
	"os"
)

var CGroupV2NotSupported = errors.New("cgroup v2 is not supported")

type CGroupV2Controller struct {
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

//func NewCGroupV2Controller() (*CGroupV2Controller, error) {
//
//}
