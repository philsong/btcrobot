// Copyright 2014 The btcbot Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://github.com/philsong/
// Author：Phil	78623269@qq.com

package logger

import (
	. "config"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var (
	// 日志文件
	trade_file    = ROOT + "/log/trade.log"
	override_file = ROOT + "/log/report"
	info_file     = ROOT + "/log/info.log"
	debug_file    = ROOT + "/log/debug.log"
	trace_file    = ROOT + "/log/trace.log"
	error_file    = ROOT + "/log/error.log"
	fatal_file    = ROOT + "/log/fatal.log"
)

func init() {
	os.Mkdir(ROOT+"/log/", 0777)
	if Config["env"] == "test" {
		override_file += "_test"
	}
}

type logger struct {
	*log.Logger
}

func New(out io.Writer) *logger {
	return &logger{
		Logger: log.New(out, "", log.LstdFlags),
	}
}

func NewReport(out io.Writer) *logger {
	return &logger{
		Logger: log.New(out, "", log.LstdFlags),
	}
}

func Tradef(format string, args ...interface{}) {
	file, err := os.OpenFile(trade_file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return
	}
	defer file.Close()
	New(file).Printf(format, args...)
	if Config["infoconsole"] == "1" {
		log.Printf(format, args...)
	}
}

func Tradeln(args ...interface{}) {
	file, err := os.OpenFile(trade_file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return
	}
	defer file.Close()
	New(file).Println(args...)
	if Config["infoconsole"] == "1" {
		log.Println(args...)
	}
}

func Infof(format string, args ...interface{}) {
	file, err := os.OpenFile(info_file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return
	}
	defer file.Close()
	New(file).Printf(format, args...)
	if Config["infoconsole"] == "1" {
		log.Printf(format, args...)
	}
}

func Infoln(args ...interface{}) {
	file, err := os.OpenFile(info_file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return
	}
	defer file.Close()
	New(file).Println(args...)
	if Config["infoconsole"] == "1" {
		log.Println(args...)
	}
}

func Errorf(format string, args ...interface{}) {
	file, err := os.OpenFile(error_file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return
	}
	defer file.Close()
	New(file).Printf(format, args...)
	if Config["errorconsole"] == "1" {
		log.Printf(format, args...)
	}
}

func Errorln(args ...interface{}) {
	file, err := os.OpenFile(error_file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return
	}
	defer file.Close()
	// 加上文件调用和行号
	_, callerFile, line, ok := runtime.Caller(1)
	if ok {
		args = append([]interface{}{"[", filepath.Base(callerFile), "]", line}, args...)
	}
	New(file).Println(args...)
	if Config["errorconsole"] == "1" {
		log.Println(args...)
	}
}

func Fatalf(format string, args ...interface{}) {
	file, err := os.OpenFile(fatal_file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return
	}
	defer file.Close()
	New(file).Printf(format, args...)
	if Config["fatalconsole"] == "1" {
		log.Printf(format, args...)
	}
}

func Fatalln(args ...interface{}) {
	file, err := os.OpenFile(fatal_file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return
	}
	defer file.Close()
	// 加上文件调用和行号
	_, callerFile, line, ok := runtime.Caller(1)
	if ok {
		args = append([]interface{}{"[", filepath.Base(callerFile), "]", line}, args...)
	}
	New(file).Println(args...)
	if Config["fatalconsole"] == "1" {
		log.Println(args...)
	}
}

func Fatal(args ...interface{}) {
	file, err := os.OpenFile(fatal_file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return
	}
	defer file.Close()
	// 加上文件调用和行号
	_, callerFile, line, ok := runtime.Caller(1)
	if ok {
		args = append([]interface{}{"[", filepath.Base(callerFile), "]", line}, args...)
	}
	New(file).Println(args...)
	if Config["fatalconsole"] == "1" {
		log.Println(args...)
	}
}

func Debugf(format string, args ...interface{}) {
	if Config["debug"] == "1" {
		file, err := os.OpenFile(debug_file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
		if err != nil {
			return
		}
		defer file.Close()
		New(file).Printf(format, args...)

		if Config["debugconsole"] == "1" {
			log.Printf(format, args...)
		}
	}

}

func Debugln(args ...interface{}) {
	if Config["debug"] == "1" {
		file, err := os.OpenFile(debug_file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
		if err != nil {
			return
		}
		defer file.Close()
		// 加上文件调用和行号
		_, callerFile, line, ok := runtime.Caller(1)
		if ok {
			args = append([]interface{}{"[", filepath.Base(callerFile), "]", line}, args...)
		}
		New(file).Println(args...)
		if Config["debugconsole"] == "1" {
			log.Println(args...)
		}
	}
}

func Tracef(format string, args ...interface{}) {

	file, err := os.OpenFile(trace_file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return
	}
	defer file.Close()
	New(file).Printf(format, args...)

}

func Traceln(args ...interface{}) {

	file, err := os.OpenFile(trace_file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return
	}
	defer file.Close()
	// 加上文件调用和行号
	_, callerFile, line, ok := runtime.Caller(1)
	if ok {
		args = append([]interface{}{"[", filepath.Base(callerFile), "]", line}, args...)
	}
	New(file).Println(args...)
}

var what string

func OverrideStart(Peroid int) {
	what = fmt.Sprintf("%03d", Peroid)
	file, err := os.OpenFile(override_file+what+".log", os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return
	}
	defer file.Close()
}

func Overridef(format string, args ...interface{}) {
	file, err := os.OpenFile(override_file+what+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return
	}
	defer file.Close()
	NewReport(file).Printf(format, args...)
}

func Overrideln(args ...interface{}) {
	file, err := os.OpenFile(override_file+what+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return
	}
	defer file.Close()

	NewReport(file).Println(args...)
}
