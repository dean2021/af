// Copyright 2020 Dean.
// Authors: Dean <dean@csoio.com>
// Date: 2020/5/20 4:29 下午

// 文件介绍

package plugin

import (
	"github.com/dean2021/af"
	"strconv"
	"time"
)

type TestPlugin struct{}

func (tp *TestPlugin) Name() string {
	return "testplugin"
}

func (tp *TestPlugin) Entry(config *af.Config, logger af.Logger) error {
	i := 0
	for {
		i++
		config.Set("xxx", strconv.Itoa(i))
		logger.Printf("[%s]插件运行中...", tp.Name())
		time.Sleep(time.Second)
	}
	return nil
}
