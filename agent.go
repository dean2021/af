// Copyright 2020 Dean.
// Authors: Dean <dean@csoio.com>
// Date: 2020/5/20 2:53 下午

// 文件介绍

package af

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Agent struct {

	// 唯一ID
	ID string `toml:"id"`

	// Agent名称
	Name string `toml:"name"`

	// 插件
	plugins map[string]Plugin `toml:"-"`

	// agent配置
	Config *Config `toml:"-"`

	// 插件关闭管理
	cancel context.CancelFunc `toml:"-"`
}

// 注册插件
func (a *Agent) Plugin(plugin Plugin) {
	if _, ok := a.plugins[plugin.Name()]; !ok {
		a.plugins[plugin.Name()] = plugin
	}
}

// 运行
func (a *Agent) Run() error {

	// 注册agent检查
	Register(a)

	// 启动agent
	err := a.Start()
	if err != nil {
		return err
	}

	return a.stopListen()
}

// 启动
func (a *Agent) Start() error {

	// 资源限制
	err := SystemResourceLimit(a)
	if err != nil {
		return errors.New("Unable to open system resource limit:" + err.Error())
	}

	// 负载监控, 超过阈值, 则agent自杀退出
	go SystemLoadMonitor(a)

	// TODO 性能指标上报，待实现

	// 启动所有插件
	a.StartPlugin()

	return nil
}

// 监听停止
func (a *Agent) stopListen() error {

	// 捕获结束信号
	var sigChan = make(chan os.Signal, 3)
	signal.Notify(sigChan, syscall.SIGTERM, os.Interrupt)
	<-sigChan

	return a.Stop()
}

// 停止
func (a *Agent) Stop() error {
	// TODO 停止agent前进行收尾，如记录日志
	log.Println("agent停止")

	os.Exit(0)
	return nil
}

// 启动所有插件
func (a *Agent) StartPlugin() {
	for _, plugin := range a.plugins {
		go func(p Plugin) {
			log.Println(p.Name() + "插件被启动")
			if err := p.Entry(a.Config); err != nil {
				log.Fatalf("start plugin [%s] error: %s", p.Name(), err.Error())
			}
			log.Println(p.Name() + "插件运行结束")
		}(plugin)
	}
}

// 添加默认配置
func setDefaultConfig(agent *Agent) {
	// 限制100M内存
	agent.Config.Set("system.max_memory", "104857600")
	// 限制10% CPU使用率
	agent.Config.Set("system.max_cpu_quota", "10000")
	// 系统负载阈值, 超过此阈值测退出程序
	agent.Config.Set("system.max_load_limit", "0.7")
	// agent注册完信息保存文件路径
	agent.Config.Set("system.register.save_file", "./data.toml")
}

// 初始化
func NewAgent(name string) *Agent {

	var agent = &Agent{
		Name:    name,
		plugins: make(map[string]Plugin),
		Config:  new(Config),
	}

	setDefaultConfig(agent)

	return agent
}
