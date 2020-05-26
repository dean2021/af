// Copyright 2019 同程艺龙, Inc.
// Authors: Dean.lu <dean.lu@ly.com>
// Date: 2019-10-30 12:49

// 文件介绍

package system

import (
	"errors"
	"github.com/dean2021/af/util/array"
	"strings"
	"syscall"
)

// Get Linux kernel version.
func KernelVersion() (string, error) {
	var uName syscall.Utsname
	if err := syscall.Uname(&uName); err != nil {
		return "", errors.New("unable to get kernel version:" + err.Error())
	}
	release := array.ToString(uName.Release)
	versions := strings.Split(release, "-")
	if len(versions) == 0 {
		return release, nil
	}
	return versions[0], nil
}

func FileStatSys(sys interface{}) *syscall.Stat_t {
	return sys.(*syscall.Stat_t)
}
