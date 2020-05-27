// Copyright 2020 Dean.
// Authors: Dean <dean@csoio.com>
// Date: 2020/5/20 4:29 下午

// 文件介绍

package plugin

import (
	"github.com/dean2021/af"
	"github.com/dean2021/af/_example/service"
	"log"
	"time"
)

type TestService struct{}

func (tp *TestService) Name() string {
	return "TestService"
}

func (tp *TestService) Entry(config *af.Config, logger af.Logger) error {

	var s service.DataService
	err := af.GetService("grpc", &s)
	if err != nil {
		log.Fatalf("加载服务失败")
	}

	err = s.InitRPCService(config)
	if err != nil {
		log.Fatalf("服务初始化失败")
	}

	for {
		_ = s.SendMsg("1111111")
		log.Printf("[%s]插件运行中...", tp.Name())
		time.Sleep(time.Second * 2)
	}
	return nil
}
