// Copyright 2020 Dean.
// Authors: Dean <dean@csoio.com>
// Date: 2020/5/19 6:28 下午

// 插件接口定义

package af

type Plugin interface {

	// 插件名称
	Name() string

	// 收到指令被动调用
	Command(name string, body string)

	// 插件入口函数
	Entry(config *Config, notify *Notify, logger Logger) error
}
