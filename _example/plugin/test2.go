// Copyright 2020 Dean.
// Authors: Dean <dean@csoio.com>
// Date: 2020/5/20 4:29 下午

// 文件介绍

package plugin

import (
	"github.com/dean2021/af"
	"time"
)

type TestPlugin2 struct{}

func (tp *TestPlugin2) Name() string {
	return "testplugin2"
}

func (tp *TestPlugin2) Entry(config *af.Config, logger af.Logger) error {
	logger.Println(tp.Name() + "读取配置" + config.Get("user.mysql"))
	for {
		//fmt.Println(config.Get("xxx"))
		logger.Printf("[%s]插件运行中...", tp.Name())
		time.Sleep(time.Second)
	}
	return nil
}
