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

func (tp *TestPlugin) Command(name string, body string) {
	fmt.Println("[COMMAND] 指令名:", name, "指令内容:", body)
}

func (tp *TestPlugin) Config(body string) {
	fmt.Println("[CONFIG] ", body)
}

func (tp *TestPlugin) Entry(config *af.Config, notify *af.Notify, logger af.Logger) error {
	// 直接获取配置
	cfg, err := notify.GetConfig(tp.Name())
	if err == nil {
		fmt.Println("获取到配置:", cfg)
	}
	for {
		logger.Printf("[%s]插件运行中...", tp.Name())
		time.Sleep(time.Second * 5)
	}
}
