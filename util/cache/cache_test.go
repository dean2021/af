// Copyright 2019 Dean.
// Authors: Dean.lu <dean@csoio.com>
// Date: 2019/11/15 4:41 下午

// 文件介绍

package cache

import (
	"testing"
)

func TestGetCache(t *testing.T) {
	key := "/tmp/xxxxx"
	value, has, err := Get(key)
	if err != nil {
		t.Fatal(err)
	}
	if !has {
		t.Log("缓存内容为空")
	}
	t.Log(value)
}

func TestSetCache(t *testing.T) {
	key := "/tmp/xxxxx"
	err := SetCache(key, "0.7")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetFloat64(t *testing.T) {
	key := "/tmp/xxxxx"
	value, has, err := GetFloat64(key)
	if err != nil {
		t.Fatal(err)
	}
	if !has {
		t.Log("缓存内容为空")
	}
	t.Log(value)
}
