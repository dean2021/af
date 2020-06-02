// Copyright 2019 Dean.
// Authors: Dean.lu <dean@csoio.com>
// Date: 2019/11/15 4:22 下午

// Cache

package cache

import (
	"github.com/dean2021/af/util/file"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
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

func GetFloat64(key string) (float64, bool, error) {
	value, has, err := Get(key)
	if err != nil {
		return 0, false, err
	}
	var v float64
	if has {
		v, err = strconv.ParseFloat(value, 64)
		if err != nil {
			return 0, false, err
		}
	}
	return v, true, nil
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
	ok := file.PathExists(cacheFileDir)
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
	ok := file.PathExists(cacheFilePath)
	if !ok {
		return "", false, nil
	}
	b, err := ioutil.ReadFile(cacheFilePath)
	if err != nil {
		return "", false, err
	}
	return string(b), true, nil
}
