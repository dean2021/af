// Copyright 2019 Dean.
// Authors: Dean.lu <dean@csoio.com>
// Date: 2019-10-30 16:57

// HTTP日志hook

package logger

import (
	"bytes"
	"errors"
	"github.com/dean2021/af/util/retry"
	"log"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type HttpHook struct {
	levels    []logrus.Level
	formatter logrus.Formatter
	client    *http.Client
	uri       string
}

type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"result"`
}

// Create a new log HttpHook.
func NewHttpHook(levels []logrus.Level, formatter logrus.Formatter, uri string) *HttpHook {
	client := &http.Client{}
	hook := &HttpHook{
		levels:    levels,
		formatter: formatter,
		client:    client,
		uri:       uri,
	}
	return hook
}

func (hh *HttpHook) Levels() []logrus.Level {
	return hh.levels
}

func (hh *HttpHook) Fire(entry *logrus.Entry) error {
	b, err := hh.formatter.Format(entry)
	if err != nil {
		return err
	}
	err = retry.Do(
		func() error {
			req, err := http.NewRequest("POST", hh.uri, bytes.NewBuffer(b))
			if err != nil {
				return err
			}
			req.Header.Set("Content-Type", "application/json")
			resp, err := hh.client.Do(req)
			if err != nil {
				return err
			}

			if resp.Body == nil {
				return errors.New("接口返回内容为空")
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				return errors.New("接口返回状态码非200: " + resp.Status)
			}
			return nil
		},
		// 重试三次
		retry.Attempts(3),
		retry.Delay(time.Second),
		retry.LastErrorOnly(false),
		retry.OnRetry(func(n uint, err error) {
			log.Printf("Failed to send log#%d: %s\n", n, err)
		}),
	)

	return err
}
