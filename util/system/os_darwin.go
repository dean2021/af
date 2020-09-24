// Copyright 2020 Dean.
// Authors: Dean <dean@csoio.com>
// Date: 2020/9/22 5:28 下午

// 文件介绍

package system

import (
	"bytes"
	"golang.org/x/sys/unix"
	"runtime"
)

func GetPlatform() string {
	return runtime.GOOS
}

func GetVersion() string {
	var utsName unix.Utsname
	err := unix.Uname(&utsName)
	if err != nil {
		return "unknown"
	}
	return string(utsName.Release[:bytes.IndexByte(utsName.Release[:], 0)])
}
