// Copyright 2020 Dean.
// Authors: Dean <dean@csoio.com>
// Date: 2020/9/24 12:28 下午

// 文件介绍

package system

import (
	"os"
	"path/filepath"
	"strings"
)

func GetCurrentDirectory() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return strings.Replace(dir, "\\", "/", -1)
}
