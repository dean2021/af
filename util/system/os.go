// Copyright 2020 Dean.
// Authors: Dean <dean@csoio.com>
// Date: 2020/9/22 5:28 下午

package system

import (
	"encoding/json"
	"os"
	"runtime"
	"strings"
)

type OS struct {
	// 主机名
	HostName string
	// 真实使用的IP 地址
	RealIP string
	// MAC地址
	MAC string
	// IPv4 IP地址
	IPv4 string
	// IPv6 ip地址
	IPv6 string
	// 系统内核, linux / windows
	Kernel string
	// 系统位数 32/64
	KernelBit int
	// 系统版本,  如: 3.3.0-3.58.24-Rel-2019-05-16_11-57-14-122 / 3.3.0-3.58.22-WIN-Rel-2019-05-15_09-53-18-102
	KernelVersion string
	// 平台名，如: CentOS release 6.4 (Final) / Windows Server 2008 R2 Enterprise Service Pack 1 (build 7601)
	Platform string
}

func (os *OS) String() string {
	s, _ := json.Marshal(os)
	return string(s)
}

func GetSystemInfo() *OS {
	hostName, err := os.Hostname()
	if err != nil {
		hostName = "Unknown"
	}
	ip := GetIP()
	return &OS{
		HostName:      hostName,
		RealIP:        GetRealIP(),
		MAC:           strings.Join(GetMacAddrs(), ","),
		IPv4:          strings.Join(ip.IPv4, ","),
		IPv6:          strings.Join(ip.IPv6, ","),
		Kernel:        runtime.GOOS,
		KernelBit:     32 << (^uint(0) >> 63),
		KernelVersion: GetVersion(),
		Platform:      GetPlatform(),
	}
}
