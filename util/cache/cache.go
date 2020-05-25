// Copyright 2019 同程艺龙, Inc.
// Authors: Dean.lu <dean.lu@ly.com>
// Date: 2019/11/15 4:22 下午

// Cache

package cache

import (
	"SecurityAgent/pkg/util"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

var (
	MemoryCache sync.Map
)

func SetCache(key string, value string) error {
	err := SetFileCache(key, value)
	if err != nil {
		return err
	}
	SetMemoryCache(key, value)
	return nil
}

func Get(key string) (string, bool, error) {
	value := GetMemoryCache(key)
	if value != "" {
		return value, true, nil
	}
	value, has, err := GetFileCache(key)
	return value, has, err
}

// Write cache connect to memory
func SetMemoryCache(key string, value string) {
	MemoryCache.Store(key, value)
}

func DelMemoryCache(key string) {
	MemoryCache.Delete(key)
}

// Read cache content from memory
func GetMemoryCache(key string) string {
	value, ok := MemoryCache.Load(key)
	if !ok {
		return ""
	}
	return value.(string)
}

// Write cache content to file
func SetFileCache(cacheFilePath string, value string) error {
	cacheFilePath, err := filepath.Abs(cacheFilePath)
	if err != nil {
		return err
	}
	cacheFileDir := filepath.Dir(cacheFilePath)
	ok := util.PathExists(cacheFileDir)
	if !ok {
		err = os.MkdirAll(cacheFileDir, 0666)
		if err != nil {
			return err
		}
	}
	return ioutil.WriteFile(cacheFilePath, []byte(value), 0666)
}

// Read cache content from file
func GetFileCache(cacheFilePath string) (string, bool, error) {
	ok := util.PathExists(cacheFilePath)
	if !ok {
		return "", false, nil
	}
	b, err := ioutil.ReadFile(cacheFilePath)
	if err != nil {
		return "", false, err
	}
	return string(b), true, nil
}
