// +build windows
// Copyright 2020 Dean.
// Authors: Dean <dean@csoio.com>
// Date: 2020/5/19 6:17 下午

// 以服务形式启动插件

package main

import (
	"github.com/dean2021/af"
	"github.com/dean2021/af/_example/plugin"
	"github.com/kardianos/service"
	"log"
)

func main() {

	// 新建一个agent
	agent := af.NewAgent("hs-agent")

	// 限制100M内存
	agent.Config.Set("system.max_memory", "104857600")

	// 限制10% CPU使用率
	agent.Config.Set("system.max_cpu_quota", "10000")

	// 系统负载阈值, 超过此阈值测退出程序
	agent.Config.Set("system.max_load_limit", "0.7")

	// agent注册api
	agent.Config.Set("system.register.api", "http://soc.qa.csoio.com/api/hostsecurity/agent/register")

	// agent注册信息保存文件
	agent.Config.Set("system.register.save_file", "./data.toml")

	// 用户自定义配置
	agent.Config.Set("user.mysql", "mysql://127.0.0.1:3306")

	// 注册插件
	agent.Plugin(new(plugin.TestPlugin))
	agent.Plugin(new(plugin.TestPlugin2))

	config := &service.Config{
		Name:        "hs-agent",
		DisplayName: "hs-agent",
		Description: "security agent for hids service",
	}

	s, err := af.NewService(config, agent)
	if err != nil {
		log.Fatal(err)
	}

	err = s.Run()
	if err != nil {
		log.Fatal(err)
	}
}
