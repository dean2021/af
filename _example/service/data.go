// Copyright 2020 Dean.
// Authors: Dean <dean@csoio.com>
// Date: 2020/5/26 4:01 下午

// grpc数据通信

package service

import (
	"context"
	"github.com/dean2021/af"
	pb "github.com/dean2021/af/_example/service/helloworld"
	"google.golang.org/grpc"
)

type DataService struct {
	conn *grpc.ClientConn
}

func (ds *DataService) InitRPCService(config *af.Config) error {
	if ds.conn != nil {
		return nil
	}
	var err error
	ds.conn, err = grpc.Dial(config.Get("service.grpc.addr"), grpc.WithInsecure())
	return err
}

func (ds *DataService) SendMsg(data string, logger af.Logger) error {
	logger.Println("发送数据" + data)
	c := pb.NewGreeterClient(ds.conn)
	r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: data})
	if err != nil {
		logger.Fatalf("could not greet: %v", err)
	}
	logger.Println(r.Message)
	return nil
}
