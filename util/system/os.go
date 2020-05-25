package system

import (
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

// 获取主机名
func HostName() string {
	hostName, err := os.Hostname()
	if err != nil {
		logrus.Fatal(err)
		return "Unknown"
	}
	return hostName
}

func GetCurrentDirectory() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return strings.Replace(dir, "\\", "/", -1)
}
