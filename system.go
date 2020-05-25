// Copyright 2020 Dean.
// Authors: Dean <dean@csoio.com>
// Date: 2020/5/20 5:22 下午

// 系统管理

package af

import (
	"github.com/shirou/gopsutil/load"
	"log"
	"runtime"
	"strconv"
	"time"
)

// 负载监控
func SystemLoadMonitor(agent *Agent) {

	maxLoadLimit, err := strconv.ParseFloat(agent.Config.Get("system.max_load_limit"), 64)
	if err != nil {
		log.Fatal(err)
	}

	// TODO windows负载监控暂未实现
	if runtime.GOOS == "windows" {
		return
	}

	sleepTime := time.Minute * 15
	for {
		avg, err := load.Avg()
		if err != nil {
			log.Fatal(err)
		}
		// cpu每个核平均负载(15分钟内)
		avgCoreLoad := avg.Load15 / float64(runtime.NumCPU())
		// 系统负载超过阈值, 则agent退出
		if avgCoreLoad > maxLoadLimit {
			log.Printf("系统负载过高(%v),已超过设定阈值(%v), agent退出", avgCoreLoad, maxLoadLimit)
			agent.Stop()
		}
		time.Sleep(sleepTime)
	}
}
