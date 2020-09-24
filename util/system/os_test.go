// Copyright 2020 Dean.
// Authors: Dean <dean@csoio.com>
// Date: 2020/9/24 12:08 下午

// 文件介绍

package system

import "testing"

func TestGetSystemInfo(t *testing.T) {
	info := GetSystemInfo()
	t.Log(info)
}
