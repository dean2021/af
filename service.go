// Copyright 2020 Dean.
// Authors: Dean <dean@csoio.com>
// Date: 2020/5/26 3:43 下午

// 服务管理

package af

import (
	"errors"
	"sync"
)

var services sync.Map

type Service interface{}

func AddService(name string, service Service) {
	services.Store(name, service)
}

func GetService(name string, service Service) error {
	load, o := services.Load(name)
	if !o {
		return errors.New("找不到" + name + "服务")
	}
	service = load.(Service)
	return nil
}
