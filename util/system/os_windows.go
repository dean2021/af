// Copyright 2020 Dean.
// Authors: Dean <dean@csoio.com>
// Date: 2020/9/22 5:28 下午

package system

import (
	"errors"
	"github.com/StackExchange/wmi"
	"runtime"
)

type operatingSystem struct {
	BuildNumber             string
	Caption                 string
	CSDVersion              string
	OperatingSystemSKU      uint32
	ServicePackMajorVersion uint16
	ServicePackMinorVersion uint16
	Version                 string
	OSArchitecture          string
}

// 获取操作系统信息
// 参考：https://docs.microsoft.com/en-us/windows/win32/cimwin32prov/win32-operatingsystem
func getOperatingSystemInfo() (*operatingSystem, error) {
	var OperatingSystems []operatingSystem
	err := wmi.Query("SELECT * FROM Win32_OperatingSystem", &OperatingSystems)
	if err != nil {
		return nil, err
	}
	if len(OperatingSystems) == 0 {
		return nil, errors.New("failed to get system information")
	}
	return &OperatingSystems[0], nil
}

func GetPlatform() string {
	osInfo, err := getOperatingSystemInfo()
	if err != nil {
		return runtime.GOOS
	}
	return osInfo.Caption
}

func GetVersion() string {
	osInfo, err := getOperatingSystemInfo()
	if err != nil {
		return "Unknown"
	}
	return osInfo.Version
}
