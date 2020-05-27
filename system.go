// Copyright 2020 Dean.
// Authors: Dean <dean@csoio.com>
// Date: 2020/5/20 5:22 下午

// 系统管理

package af

import (
	"github.com/kardianos/service"
	"github.com/shirou/gopsutil/load"
	"log"
	"runtime"
	"strconv"
	"time"
)

type program struct {
	agent *Agent
}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}
func (p *program) run() {
	err := p.agent.Run()
	if err != nil {
		log.Fatal(err)
	}
}
func (p *program) Stop(s service.Service) error {
	return p.agent.Stop()
}

// 初始化一个服务
func NewSystemService(config *service.Config, agent *Agent) (service.Service, error) {
	prg := &program{
		agent: agent,
	}
	s, err := service.New(prg, config)
	if err != nil {
		return s, nil
	}
	return s, nil
}

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

	sleepTime := time.Minute
	for {
		avg, err := load.Avg()
		if err != nil {
			log.Fatal(err)
		}
		// cpu每个核平均负载(5分钟内)
		avgCoreLoad := avg.Load5 / float64(runtime.NumCPU())
		// 系统负载超过阈值, 则agent退出
		if avgCoreLoad > maxLoadLimit {
			log.Printf("系统负载过高(%v),已超过设定阈值(%v), agent退出", avgCoreLoad, maxLoadLimit)
			agent.Stop()
		}
		time.Sleep(sleepTime)
	}
}
