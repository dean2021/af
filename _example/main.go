// Copyright 2020 Dean.
// Authors: Dean <dean@csoio.com>
// Date: 2020/5/19 6:17 下午

package main

import (
	"github.com/dean2021/af"
	"github.com/dean2021/af/_example/logger"
	"github.com/dean2021/af/_example/plugin"
	"github.com/sirupsen/logrus"
	"log"
	"path"
	"time"
)

func main() {

	// 新建一个agent
	agent := af.NewAgent("hs-agent")

	// 是否启用cgroup, 默认启用, 参数 on/off
	agent.Config.Set("system.cgroup_enable", "on")

	// 限制100M内存
	agent.Config.Set("system.max_memory", "104857600")

	// 限制10% CPU使用率
	agent.Config.Set("system.max_cpu_quota", "10000")

	// 系统负载阈值, 超过此阈值则退出程序
	agent.Config.Set("system.max_load_limit", "0.7")

	// 系统cpu使用率阈值，超过此阈值则自杀
	agent.Config.Set("system.max_cpu_usage_limit", "80")

	// agent注册api
	agent.Config.Set("system.register.api", "http://www.baidu.com/api/hostsecurity/agent/register")

	// agent注册信息保存文件
	agent.Config.Set("system.register.save_file", "./data.toml")

	// 用户自定义配置
	agent.Config.Set("user.mysql", "mysql://127.0.0.1:3306")
	agent.Config.Set("service.grpc.addr", "localhost:50001")

	// 添加服务
	//af.AddService("grpc", new(service.DataService))

	// 添加插件
	agent.Plugin(new(plugin.TestPlugin))
	agent.Plugin(new(plugin.TestPlugin2))
	//agent.Plugin(new(plugin.TestService))

	// 替换日志组件
	l := logrus.New()
	// 设置日志格式
	l.SetFormatter(&logger.JSONFormatter{})
	// 添加log http hook
	l.AddHook(logger.NewHttpHook(logrus.AllLevels, &logger.JSONFormatter{}, "http://www.baidu.com/api/hostsecurity/agentLog/add"))
	// 添加log滚动文件切割hook
	hook, err := logger.NewRotateHook(path.Join("./", "logs"), "debug.log", time.Hour*24, time.Second*60)
	if err != nil {
		panic(err)
	}
	l.AddHook(hook)
	agent.SetLogger(l)

	// 运行agent
	err = agent.Run()
	if err != nil {
		log.Fatal(err)
	}
}
