// Copyright 2020 Dean.
// Authors: Dean <dean@csoio.com>
// Date: 2020/9/22 5:28 下午

// 文件介绍

package system

import (
	"bytes"
	"golang.org/x/sys/unix"
	"io/ioutil"
	"runtime"
	"strings"
)

func GetPlatform() string {
	b, err := ioutil.ReadFile("/etc/redhat-release")
	if err != nil {
		b, err = ioutil.ReadFile("/etc/issue")
		if err != nil {
			return runtime.GOOS
		}
	}
	platform := strings.Replace(string(b), "\\n \\l", "", -1)
	platform = strings.TrimSpace(platform)
	return platform
}

func GetVersion() string {
	var utsName unix.Utsname
	err := unix.Uname(&utsName)
	if err != nil {
		return "unknown"
	}
	return string(utsName.Release[:bytes.IndexByte(utsName.Release[:], 0)])
}
