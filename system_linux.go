// Copyright 2020 Dean.
// Authors: Dean <dean@csoio.com>
// Date: 2020/5/20 5:22 下午

// 进程资源限制, linux系统下采用cgroup实现

package af

import (
	"fmt"
	"github.com/dean2021/af/util/file"
	"github.com/containerd/cgroups"
	"github.com/opencontainers/runtime-spec/specs-go"
	"os"
	"path"
	"strconv"
	"syscall"
)

const (
	CgroupRootPath = "/cgroup"
)

func mount(source string, fstype string, path string, option string) error {
	ok := file.PathExists(path)
	if !ok {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
	}
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	return syscall.Mount(source, path, fstype, uintptr(defaultMountFlags), option)
}

func mountCgroupFileSystem() error {
	/*
		mkdir /tmp/cgroup
		mount -t tmpfs cgroup_root /tmp/cgroup
		mkdir /tmp/cgroup/cpuset
		mount -t cgroup -ocpuset cpuset /tmp/cgroup/cpuset/
		mkdir /tmp/cgroup/cpu
		mount -t cgroup -ocpu cpu /tmp/cgroup/cpu/
		mkdir /tmp/cgroup/memory
		mount -t cgroup -omemory memory /tmp/cgroup/memory/
	*/
	err := mount("cgroup", "tmpfs", CgroupRootPath, "")
	if err != nil {
		return err
	}
	subSystems := []string{"cpuset", "cpu", "memory", "blkio"}
	for _, subSystem := range subSystems {
		err := mount("cgroup", "cgroup", path.Join(CgroupRootPath, subSystem), subSystem)
		if err != nil {
			return err
		}
	}
	return nil
}

// 根据进程添加资源限制
func SystemResourceLimit(agent *Agent) error {

	maxCPUQuota, err := strconv.ParseInt(agent.Config.Get("system.max_cpu_quota"), 10, 64)
	if err != nil {
		return err
	}
	maxMemory, err := strconv.ParseInt(agent.Config.Get("system.max_memory"), 10, 64)
	if err != nil {
		return err
	}

	// 检查cgroup是否挂载
	_, err = cgroups.V1()
	if err == cgroups.ErrMountPointNotExist {
		err = mountCgroupFileSystem()
		if err != nil {
			return fmt.Errorf("mount cgroup: %v", err)
		}
	}

	// 将进程加入cgroup
	control, err := cgroups.New(cgroups.V1, cgroups.StaticPath("/cgroup"), &specs.LinuxResources{
		CPU: &specs.LinuxCPU{
			Quota: &maxCPUQuota,
		},
		Memory: &specs.LinuxMemory{
			Limit: &maxMemory,
		},
		//BlockIO: &specs.LinuxBlockIO{
		//	Weight: &configs.CgroupBlockIOWeight,
		//},
	})
	if err != nil {
		return err
	}
	err = control.Add(cgroups.Process{
		Pid: os.Getpid(),
	})
	return err
}
