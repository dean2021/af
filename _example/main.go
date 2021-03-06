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

	// 启动etcd
	agent.Config.Set("system.etcd.enable", "on")

	// etcd连接设置
	agent.Config.Set("system.etcd.endpoints", "127.0.0.1:2379")
	agent.Config.Set("system.etcd.dial_timeout", "5s")
	// AutoSyncInterval is the interval to update endpoints with its latest members.
	// 0 disables auto-sync. By default auto-sync is disabled.
	agent.Config.Set("system.etcd.auto_sync_interval", "5m0s")
	// etcd连接账号
	agent.Config.Set("system.etcd.username", "")
	// etcd连接密码
	agent.Config.Set("system.etcd.password", "")
	// 命名空间, 防止和其他应用的path冲突
	agent.Config.Set("system.etcd.namespace", "af")

	// 添加服务
	//af.AddService("grpc", new(service.DataService))

	// 添加插件
	agent.Plugin(new(plugin.TestPlugin2))
	//agent.Plugin(new(plugin.TestService))

	// 替换日志组件
	l := logrus.New()
	// 设置日志格式
	l.SetFormatter(&logger.JSONFormatter{})
	// 添加log http hook
	l.AddHook(logger.NewHttpHook(logrus.AllLevels, &logger.JSONFormatter{}, "http://127.0.0.1:2333/api/hostsecurity/agentLog/add"))
	agent.SetLogger(l)

	// 运行agent
	err := agent.Run()
	if err != nil {
		log.Fatal(err)
	}
}
