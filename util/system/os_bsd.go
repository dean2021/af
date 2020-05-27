// +build darwin netbsd freebsd openbsd dragonfly

// Copyright 2019 Dean.
// Authors: Dean.lu <dean@csoio.com>
// Date: 2019-10-30 12:48

// 文件介绍

package system

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
	"syscall"
)

// Get Mac kernel version.
func KernelVersion() (string, error) {
	cmd := exec.Command("uname", "-r")
	var o, e bytes.Buffer
	cmd.Stdout = &o
	cmd.Stderr = &e
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	err := cmd.Start()
	if err != nil {
		return "", err
	}
	err = cmd.Wait()
	if err != nil {
		return "", err
	}
	version := strings.Split(o.String(), "-")
	if len(version) == 0 {
		return "", errors.New("unable to get kernel version number")
	}

	return strings.TrimSpace(version[0]), err
}

func FileStatSys(sys interface{}) *syscall.Stat_t {
	return sys.(*syscall.Stat_t)
}
