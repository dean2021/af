// Copyright 2019 Dean.
// Authors: Dean.lu <dean@csoio.com>
// Date: 2019-10-30 16:57

// 日志切割

package logger

import (
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

func NewRotateHook(logPath string, logFileName string, maxAge time.Duration, rotationTime time.Duration) (*lfshook.LfsHook, error) {
	err := os.MkdirAll(logPath, os.ModePerm)
	if err != nil {
		return nil, errors.New("config local file system logger error:" + err.Error())
	}
	baseLogPath := path.Join(logPath, logFileName)
	writer, err := rotatelogs.New(
		baseLogPath+".%Y-%m-%d",
		rotatelogs.WithLinkName(baseLogPath),      // 生成软链，指向最新日志文
		rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)
	if err != nil {
		return nil, errors.New("config local file system logger error:" + err.Error())
	}

	return lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer, // 为不同级别设置不同的输出目的
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &logrus.TextFormatter{DisableColors: true, TimestampFormat: "2006-01-02 15:04:05.000"}), nil
}
