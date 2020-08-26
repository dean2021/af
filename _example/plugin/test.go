// Copyright 2020 Dean.
// Authors: Dean <dean@csoio.com>
// Date: 2020/5/20 4:29 下午

// 文件介绍

package plugin

import (
	"fmt"
	"github.com/dean2021/af"
	"time"
)

type TestPlugin struct{}

func (tp *TestPlugin) Name() string {
	return "testplugin"
}

func (tp *TestPlugin) Entry(config *af.Config, notify *af.Notify, logger af.Logger) error {

	// 获取指令
	go notify.WatchCommand(tp.Name(), func(command af.Command) {
		fmt.Println("[COMMAND] 指令名:", command.Name, "指令内容:", command.Body)
	})

	// 监听配置变更
	go notify.WatchConfig(tp.Name(), func(value string) {
		fmt.Println("[CONFIG] ", value)
	})

	// 直接获取配置
	keyVal, err := notify.GetConfig(tp.Name())
	if err == nil {
		fmt.Println("获取到配置:", keyVal)
	}

	for {
		logger.Printf("[%s]插件运行中...", tp.Name())
		time.Sleep(time.Second * 5)
	}
}
