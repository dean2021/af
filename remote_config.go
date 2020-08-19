// Copyright 2020 Dean.
// Authors: Dean <dean@csoio.com>
// Date: 2020/5/20 4:00 下午

// 基于ETcd实现的指令通讯

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

type RemoteConfig struct {

	// etcd根路径
	//使用命名空间方式进行防止etcd路径冲突
	namespace string

	// ETcd 连接客户端
	client *clientv3.Client

	// ETcd 上下文
	context context.Context

	// 唯一di
	uuid string
}

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// 监听配置变更
func (rc *RemoteConfig) WatchChange(key string, EventFunc func(k string, v string)) {
	path := fmt.Sprintf("/%s/%s/config/%s/", rc.namespace, rc.uuid, key)
	log.Println("监听:", path)
	for {
		v, err := rc.getRevision()
		if err == nil {
			rch := rc.client.Watch(rc.context, path, clientv3.WithPrefix(), clientv3.WithRev(v))
			for wResp := range rch {
				for _, ev := range wResp.Events {
					err := rc.setRevision(ev.Kv.ModRevision + 1)
					if err != nil {
						log.Println(err)
					}
					if ev.Type == clientv3.EventTypePut {
						EventFunc(string(ev.Kv.Value), string(ev.Kv.Key))
					}
				}
			}
		} else {
			log.Println(err)
		}
		time.Sleep(time.Second)
	}
}

// 获取配置信息
func (rc *RemoteConfig) SetConfig(key string, value string) error {
	path := fmt.Sprintf("/%s/%s/config/%s/", rc.namespace, rc.uuid, key)
	_, err := rc.client.Put(rc.context, path, value)
	return err
}

// 获取配置信息
func (rc *RemoteConfig) GetConfig(key string) ([]KeyValue, error) {
	var keyValue []KeyValue
	path := fmt.Sprintf("/%s/%s/config/%s/", rc.namespace, rc.uuid, key)
	resp, err := rc.client.Get(rc.context, path, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	for _, v := range resp.Kvs {
		keyValue = append(keyValue, KeyValue{
			Key:   string(v.Key),
			Value: string(v.Value),
		})
	}
	return keyValue, err
}

// 保存当前命令版本
func (rc *RemoteConfig) setRevision(revision int64) error {
	path := fmt.Sprintf(
		"/%s/%s/revision", rc.namespace, rc.uuid)
	_, err := rc.client.Put(rc.context, path, strconv.FormatInt(revision, 10))
	return err
}

// 获取命令版本
func (rc *RemoteConfig) getRevision() (int64, error) {
	var curRevision int64
	path := fmt.Sprintf("/%s/%s/revision", rc.namespace, rc.uuid)
	resp, err := rc.client.Get(context.Background(), path)
	if err != nil {
		return 0, err
	}
	// 当保存的revision为空情下,用最新revision
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

func NewRemoteConfig(agent *Agent) (*RemoteConfig, error) {

	namespace := agent.Config.Get("system.etcd.namespace")
	endpoints := agent.Config.Get("system.etcd.endpoints")
	username := agent.Config.Get("system.etcd.username")
	password := agent.Config.Get("system.etcd.password")

	dialTimeout, err := time.ParseDuration(agent.Config.Get("system.etcd.dial_timeout"))
	if err != nil {
		return nil, err
	}
	autoSyncInterval, err := time.ParseDuration(agent.Config.Get("system.etcd.auto_sync_interval"))
	if err != nil {
		return nil, err
	}

	cfg := clientv3.Config{
		Endpoints:        strings.Split(endpoints, ","),
		DialTimeout:      dialTimeout,
		AutoSyncInterval: autoSyncInterval,
		Username:         username,
		Password:         password,
	}

	// 配置etcd为远程配置管理
	client, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}

	return &RemoteConfig{
		namespace: namespace,
		uuid:      agent.ID,
		context:   context.Background(),
		client:    client,
	}, nil
}
