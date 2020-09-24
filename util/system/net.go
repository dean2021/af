// Copyright 2020 Dean.
// Authors: Dean <dean@csoio.com>
// Date: 2020/9/24 10:41 上午

package system

import (
	"net"
)

type IP struct {
	IPv4 []string
	IPv6 []string
}

// 获取真实IP
func GetRealIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "Unknown"
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

// 获取所有IP
func GetIP() IP {
	var ip = IP{}
	interfaceAddr, err := net.InterfaceAddrs()
	if err != nil {
		return ip
	}
	for _, address := range interfaceAddr {
		ipNet, isValidIpNet := address.(*net.IPNet)
		if isValidIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ip.IPv4 = append(ip.IPv4, ipNet.IP.String())
			} else {
				ip.IPv6 = append(ip.IPv6, ipNet.IP.String())
			}
		}
	}
	return ip
}

// 获取本机mac地址
func GetMacAddrs() (macAddrs []string) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return
	}
	for _, netInterface := range netInterfaces {
		macAddr := netInterface.HardwareAddr.String()
		if len(macAddr) == 0 {
			continue
		}
		macAddrs = append(macAddrs, macAddr)
	}
	return macAddrs
}
