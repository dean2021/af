## Agent Framework


### 系统配置

1. system.max_memory ： 设置最大占用内容, 默认100m
2. system.max_cpu_quota : 设置最大cpu使用率， 默认10%
3. system.max_load_limit : 设置系统负载自杀阈值，默认0.7
4. system.register.api ： agent注册api地址, 必填
5. system.register.save_file : agent注册完信息保存文件路径, 默认：存放到当前目录下data.toml文件中
6. system.max_cpu_usage_limit : 设置cpu使用率自杀阈值，默认80(%)
7. system.cgroup_enable ： 是否启用cgroup, 默认启用, 参数 on/off

### 稳定性

| 控制        | windows   |  linux  |
| --------   | -----:  | :----:  |
|负载监控      | 无   |   有     |
| 资源限制（cpu/memory）         |  无   |   有   |
| CPU使用率        |    有    |  有  |

### TODO

- [x] 资源限制(内存/CPU)
- [x] 资源监控/自杀
- [x] 统一日志
- [x] 插件扩展
- [x] 统一注册
- [x] 统一日志
- [x] 服务注册
- [ ] 进程守护
- [ ] 升级卸载
- [ ] 进程守护
- [x] 配置变更
- [x] 指令通讯


### 框架使用

main.go 
```go


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

	// 限制100M内存
	agent.Config.Set("system.max_memory", "104857600")

	// 限制10% CPU使用率
	agent.Config.Set("system.max_cpu_quota", "10000")

	// 系统负载阈值, 超过此阈值则退出程序
	agent.Config.Set("system.max_load_limit", "0.7")

	// 系统cpu使用率阈值，超过此阈值则
	agent.Config.Set("system.max_cpu_usage_limit", "80")

	// agent注册api
	agent.Config.Set("system.register.api", "http://soc.qa.csoio.com/api/hostsecurity/agent/register")

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
	l.AddHook(logger.NewHttpHook(logrus.AllLevels, &logger.JSONFormatter{}, "http://www.baidu.com/logserver"))
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



```
