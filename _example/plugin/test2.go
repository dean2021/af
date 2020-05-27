// Copyright 2020 Dean.
// Authors: Dean <dean@csoio.com>
// Date: 2020/5/20 4:29 下午

// 文件介绍

package plugin

import (
	"fmt"
	"github.com/dean2021/af"
	"log"
	"time"
)

type TestPlugin2 struct{}

func (tp *TestPlugin2) Name() string {
	return "testplugin2"
}

func (tp *TestPlugin2) Entry(config *af.Config, logger af.Logger) error {
	fmt.Println(tp.Name() + "读取配置" + config.Get("user.mysql"))
	for {
		//fmt.Println(config.Get("xxx"))
		log.Printf("[%s]插件运行中...", tp.Name())
		time.Sleep(time.Second)
	}
	return nil
}
