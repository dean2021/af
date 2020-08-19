// Copyright 2020 Dean.
// Authors: Dean <dean@csoio.com>
// Date: 2020/5/20 2:53 下午

// agent框架主入口

package af

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const Version = "2.0.1"

type Agent struct {

	// 唯一ID
	ID string `toml:"id"`

	// Agent名称
	Name string `toml:"name"`

	// 插件
	plugins map[string]Plugin `toml:"-"`

	// agent配置
	Config *Config `toml:"-"`

	// 日志
	logger Logger `toml:"-"`

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
	err := Register(a)
	if err != nil {
		return fmt.Errorf("failed to register agent: %v", err)
	}

	// 注入系统配置
	a.Config.Set("system.agent.id", a.ID)
	a.Config.Set("system.agent.name", a.Name)

	// 启用etcd作为动态配置组件
	if a.Config.Get("system.etcd.enable") == "on" {
		remoteConfig, err := NewRemoteConfig(a)
		if err != nil {
			return err
		}
		a.Config.remoteConfig = remoteConfig
	}

	// 启动agent
	err = a.Start()
	if err != nil {
		return err
	}

	return a.stopListen()
}

// 启动
func (a *Agent) Start() error {

	// 负载监控, 超过阈值, 则agent自杀退出,
	SystemLoadCheck(a)
	go func(agent *Agent) {
		sleepTime := time.Minute
		for {
			SystemLoadCheck(a)
			time.Sleep(sleepTime)
		}
	}(a)

	// 监控cpu使用率
	go func(agent *Agent) {
		sleepTime := time.Second * 5
		for {
			time.Sleep(sleepTime)
			SystemCpuUsageCheck(a)
		}
	}(a)

	// 资源限制
	err := SystemResourceLimit(a)
	if err != nil {
		return errors.New("Unable to open system resource limit:" + err.Error())
	}

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
	// 停止agent前进行收尾，如记录日志
	a.logger.Println("agent停止")
	os.Exit(0)
	return nil
}

// 启动所有插件
func (a *Agent) StartPlugin() {
	for _, plugin := range a.plugins {
		go func(p Plugin) {
			a.logger.Println(p.Name() + "插件被启动")
			if err := p.Entry(a.Config, a.logger); err != nil {
				a.logger.Fatalf("[%s]插件返回错误: %s", p.Name(), err.Error())
			}
			a.logger.Println(p.Name() + "插件运行结束")
		}(plugin)
	}
}
func (a *Agent) SetLogger(logger Logger) {
	a.logger = logger
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
	// 系统cpu使用率阈值，超过此阈值则
	agent.Config.Set("system.max_cpu_usage_limit", "80")
	// 是否启用cgroup, 默认启用, 参数 on/off
	agent.Config.Set("system.cgroup_enable", "on")

	// 是否启用etcd, 默认不启用, 参数 on/off
	agent.Config.Set("system.etcd.enable", "off")
	// etcd节点地址，多个用逗号隔开
	agent.Config.Set("system.etcd.endpoints", "127.0.0.1:2379,127.0.0.1:2378")
	// etcd连接超时时间
	agent.Config.Set("system.etcd.dial_timeout", "30")
	// AutoSyncInterval is the interval to update endpoints with its latest members.
	// 0 disables auto-sync. By default auto-sync is disabled.
	agent.Config.Set("system.etcd.auto_sync_interval", "0")
	// etcd连接账号
	agent.Config.Set("system.etcd.username", "")
	// etcd连接密码
	agent.Config.Set("system.etcd.password", "")
	// 命名空间, 防止和其他应用的path冲突
	agent.Config.Set("system.etcd.namespace", "af")

}

// 初始化
func NewAgent(name string) *Agent {

	var agent = &Agent{
		Name:    name,
		plugins: make(map[string]Plugin),
		Config:  new(Config),
	}

	// 默认日志
	agent.SetLogger(log.New(os.Stdout, "", log.LstdFlags))

	setDefaultConfig(agent)

	return agent
}
