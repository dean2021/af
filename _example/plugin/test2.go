// Copyright 2020 Dean.
// Authors: Dean <dean@csoio.com>
// Date: 2020/5/20 4:29 下午

// 文件介绍

package plugin

import (
	"github.com/dean2021/af"
	"time"
)

type TestPlugin2 struct {
	logger af.Logger
}

func (tp *TestPlugin2) Name() string {
	return "testplugin2"
}

func (tp *TestPlugin2) WatchConfig(v string, k string) {
	tp.logger.Println("发现配置变更", v, k)
}

func (tp *TestPlugin2) Entry(config *af.Config, logger af.Logger) error {

	tp.logger = logger

	go config.Watch(tp.Name()+"/script", tp.WatchConfig)

	keyValue, err := config.GetRemoteConfig(tp.Name() + "/script")
	if err != nil {
		panic(err)
	}

	for _, kv := range keyValue {
		logger.Println("获取配置信息", kv)
	}

	logger.Println(tp.Name() + "读取配置" + config.Get("user.mysql"))
	for {
		//fmt.Println(config.Get("xxx"))
		logger.Printf("[%s]插件运行中...", tp.Name())
		time.Sleep(time.Second)
	}
	return nil
}
