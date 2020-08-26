// Copyright 2020 Dean.
// Authors: Dean <dean@csoio.com>
// Date: 2020/8/25 6:26 下午

// 文件介绍

package af

import (
	"context"
	"fmt"
	strings2 "github.com/dean2021/af/util/stringsplus"
	"go.etcd.io/etcd/clientv3"
	"log"
	"strconv"
	"strings"
	"time"
)

type Notify struct {
	// ETcd 命名空间
	namespace string
	// ETcd 连接客户端
	client *clientv3.Client
	// ETcd 上下文
	context context.Context
	// 唯一id
	uuid string
}

type Command struct {
	Name string
	Body string
}
type CommandHandleFunc func(command Command)
type ConfigHandleFunc func(value string)

// 监听指令
// 路径设计： /hs/uuid-xxx/plugin-xxx/command/指令名
func (n *Notify) WatchCommand(pluginName string, commandHandleFunc CommandHandleFunc) {
	path := fmt.Sprintf("/%s/%s/%s/command/", n.namespace, n.uuid, pluginName)
	log.Println("监听:", path)
	for {
		v, err := n.getRevision()
		if err == nil {
			rch := n.client.Watch(n.context, path, clientv3.WithPrefix(), clientv3.WithRev(v))
			for wResp := range rch {
				for _, ev := range wResp.Events {
					err := n.setRevision(ev.Kv.ModRevision + 1)
					if err != nil {
						log.Println(err)
						continue
					}
					if ev.Type == clientv3.EventTypePut {
						log.Println("发现命令")
						p := strings.Split(string(ev.Kv.Key), "/")
						if len(p) >= 4 {
							commandHandleFunc(Command{
								Name: p[5],
								Body: string(ev.Kv.Value),
							})
						}
					}
				}
			}
		} else {
			log.Println(err)
		}
		// TODO 重试时间待优化, 改成backoff
		time.Sleep(time.Second)
	}
}

// 监听配置变更
// 路径设计： /hs/uuid-xxx/plugin-xxx/config
func (n *Notify) WatchConfig(pluginName string, configHandleFunc ConfigHandleFunc) {
	path := fmt.Sprintf("/%s/%s/%s/config", n.namespace, n.uuid, pluginName)
	log.Println("监听:", path)
	for {
		version, err := n.getRevision()
		if err == nil {
			rch := n.client.Watch(n.context, path, clientv3.WithPrefix(), clientv3.WithRev(version))
			for wResp := range rch {
				for _, ev := range wResp.Events {
					err := n.setRevision(ev.Kv.ModRevision + 1)
					if err != nil {
						log.Println(err)
						continue
					}
					if ev.Type == clientv3.EventTypePut {
						configHandleFunc(string(ev.Kv.Value))
					}
				}
			}
		} else {
			log.Println(err)
		}
		// TODO 重试时间待优化, 改成backoff
		time.Sleep(time.Second)
	}
}

// 获取配置信息
// 路径设计： /hs/uuid-xxx/plugin-xxx/config
func (n *Notify) GetConfig(pluginName string) (string, error) {
	var value string
	path := fmt.Sprintf("/%s/%s/%s/config", n.namespace, n.uuid, pluginName)
	resp, err := n.client.Get(n.context, path, clientv3.WithPrefix())
	if err != nil {
		return "", err
	}
	if resp.Count > 0 {
		value = string(resp.Kvs[0].Value)
	}
	return value, err
}

// 保存当前命令版本
func (n *Notify) setRevision(revision int64) error {
	path := fmt.Sprintf("/%s/%s/revision", n.namespace, n.uuid)
	_, err := n.client.Put(n.context, path, strconv.FormatInt(revision, 10))
	return err
}

// 获取命令版本
func (n *Notify) getRevision() (int64, error) {
	var curRevision int64
	path := fmt.Sprintf("/%s/%s/revision", n.namespace, n.uuid)
	resp, err := n.client.Get(context.Background(), path)
	if err != nil {
		return 0, err
	}
	// 当保存的revision为空情况下,用最新revision
	if len(resp.Kvs) == 0 {
		return resp.Header.Revision + 1, nil
	}
	for _, ev := range resp.Kvs {
		version, err := strings2.ByteToInt64(ev.Value)
		if err != nil {
			continue
		}
		curRevision = version
	}
	return curRevision + 1, err
}

func NewNotify(namespace string, uuid string, client *clientv3.Client) *Notify {
	return &Notify{
		namespace: namespace,
		context:   context.Background(),
		client:    client,
		uuid:      uuid,
	}
}
