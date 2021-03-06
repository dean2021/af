// Copyright 2020 Dean.
// Authors: Dean <dean@csoio.com>
// Date: 2020/5/20 5:22 下午

// 系统管理

package af

import (
	"github.com/kardianos/service"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
	"runtime"
	"strconv"
)

type program struct {
	agent  *Agent
	logger Logger
}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}
func (p *program) run() {
	err := p.agent.Run()
	if err != nil {
		p.logger.Fatal(err)
	}
}
func (p *program) Stop(s service.Service) error {
	return p.agent.Stop()
}

// 初始化一个服务
func NewSystemService(config *service.Config, agent *Agent) (service.Service, error) {
	prg := &program{
		agent:  agent,
		logger: agent.logger,
	}
	s, err := service.New(prg, config)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// 系统负载检查
// 超过设定好的阈值, agent则自动停止
func SystemLoadCheck(agent *Agent) {
	maxLoadLimit, err := strconv.ParseFloat(agent.Config.Get("system.max_load_limit"), 64)
	if err != nil {
		agent.logger.Fatal(err)
	}

	// windows采用cpu使用率监控
	if runtime.GOOS == "windows" {
		return
	}

	avg, err := load.Avg()
	if err != nil {
		agent.logger.Fatal(err)
	}
	// cpu每个核平均负载(5分钟内)
	avgCoreLoad := avg.Load5 / float64(runtime.NumCPU())
	// 系统负载超过阈值, 则agent退出
	if avgCoreLoad > maxLoadLimit {
		agent.logger.Printf("系统负载过高(%v),已超过设定阈值(%v), agent退出", avgCoreLoad, maxLoadLimit)
		agent.Stop()
	}
}

// 系统cpu使用率
// 超过系统设置好的阈值, agent则自动停止
func SystemCpuUsageCheck(agent *Agent) {
	totalAvg, _ := cpu.Percent(0, false)
	maxCpuUsageLimit, err := strconv.ParseFloat(agent.Config.Get("system.max_cpu_usage_limit"), 64)
	if err != nil {
		agent.logger.Fatal(err)
	}
	if totalAvg[0] > maxCpuUsageLimit {
		agent.logger.Printf("系统cpu使用率过高(%v),已超过设定阈值(%v), agent退出", totalAvg[0], maxCpuUsageLimit)
		agent.Stop()
	}
}
