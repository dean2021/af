## Agent Framework


### 系统配置

1. system.max_memory ： 设置最大占用内容, 默认100m
2. system.max_cpu_quota : 设置最大cpu使用率， 默认10%
3. system.max_load_limit : 设置系统负载自杀阈值，默认0.7
4. system.register.api ： agent注册api地址, 必填
5. system.register.save_file : agent注册完信息保存文件路径, 默认：存放到当前目录下data.toml文件中


### 框架使用

main.go 
```go

package main

import (
	"github.com/dean2021/af"
	"github.com/dean2021/af/_example/plugin"
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

	// 添加插件
	agent.Plugin(new(plugin.TestPlugin))
	agent.Plugin(new(plugin.TestPlugin2))

	// 运行agent
	err := agent.Run()
	if err != nil {
		log.Fatal(err)
	}
}
```
