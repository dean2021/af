// Copyright 2020 Dean.
// Authors: Dean <dean@csoio.com>
// Date: 2020/5/19 6:28 下午

// 插件接口定义

package af

type Plugin interface {

	// 插件名称
	Name() string

	// 插件入口函数
	Entry(config *Config, logger Logger) error
}
