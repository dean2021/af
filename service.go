// Copyright 2020 Dean.
// Authors: Dean <dean@csoio.com>
// Date: 2020/5/20 5:04 下午

// 服务管理

package af

import (
	"github.com/kardianos/service"
	"log"
)

type program struct {
	agent *Agent
}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}
func (p *program) run() {
	err := p.agent.Run()
	if err != nil {
		log.Fatal(err)
	}
}
func (p *program) Stop(s service.Service) error {
	return p.agent.Stop()
}

// 初始化一个服务
func NewService(config *service.Config, agent *Agent) (service.Service, error) {
	prg := &program{
		agent: agent,
	}
	s, err := service.New(prg, config)
	if err != nil {
		return s, nil
	}
	return s, nil
}
