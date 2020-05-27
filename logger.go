// Copyright 2020 Dean.
// Authors: Dean <dean@csoio.com>
// Date: 2020/5/27 3:16 下午

// 文件介绍

package af

type Logger interface {
	Printf(format string, v ...interface{})
	Print(v ...interface{})
	Println(v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatalln(v ...interface{})
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
	Panicln(v ...interface{})
}
