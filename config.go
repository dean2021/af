// Copyright 2020 Dean.
// Authors: Dean <dean@csoio.com>
// Date: 2020/5/22 3:01 下午

// 配置管理

package af

import (
	"sync"
)

type Config struct {
	sync.Map
}

// 设置配置
func (c *Config) Set(key, value string) {
	c.Store(key, value)
}

// 获取配置
func (c *Config) Get(key string) string {
	v, ok := c.Load(key)
	if !ok {
		return ""
	}
	return v.(string)
}
