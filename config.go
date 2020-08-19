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
	remoteConfig *RemoteConfig
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

func (c *Config) Watch(key string, callback func(k string, v string)) {
	c.remoteConfig.WatchChange(key, callback)
}

func (c *Config) GetRemoteConfig(key string) ([]KeyValue, error) {
	return c.remoteConfig.GetConfig(key)
}

func (c *Config) SetRemoteConfig(key string, val string) error {
	return c.remoteConfig.SetConfig(key, val)
}
