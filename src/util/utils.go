// Copyright 2014 The btcbot Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://github.com/philsong/
// Author：Phil	78623269@qq.com

package util

import (
	"compress/gzip"
	"fmt"
	"io"
	"logger"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func AddRecord(record, filename string) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return
	}
	defer file.Close()

	file.WriteString(fmt.Sprintf("%s\n", record))
}

func InterfaceToFloat64(iv interface{}) (retV float64) {
	switch ivTo := iv.(type) {
	case float64:
		retV = ivTo
	case string:
		{
			var err error
			retV, err = strconv.ParseFloat(ivTo, 64)
			if err != nil {
				logger.Fatalln("convert failed, crash!")
				return 0
			}
		}
	default:
		logger.Fatalln(ivTo)
		logger.Fatalln("don't know the type, crash!")
		return 0
	}

	return retV
}

func DumpGZIP(r io.Reader) string {
	var body string
	reader, _ := gzip.NewReader(r)
	for {
		buf := make([]byte, 1024)
		n, err := reader.Read(buf)

		if err != nil && err != io.EOF {
			panic(err)
		}

		if n == 0 {
			break
		}
		body += string(buf)
	}
	return body
}

func RandomString(l int) string {
	rand.Seed(time.Now().UnixNano())
	var result string

	for i := 0; i < l; i++ {
		result += string(randdigit())
	}
	return result
}

func randdigit() uint8 {
	answers := "0123456789"

	return answers[rand.Intn(len(answers))]
}

func IntegerToString(value int64) (s string) {
	s = strconv.FormatInt(value, 10)
	return
}

func StringToInteger(s string) (value int64) {
	value, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		value = 0
	}
	return
}

func FloatToString(value float64) (s string) {
	s = strconv.FormatFloat(value, 'f', -1, 64)
	return
}

func StringToFloat(in string) float64 {
	out, err := strconv.ParseFloat(in, 64)
	if err != nil {
		logger.Fatalln("don't know the type, crash!")
		logger.Traceln("config item amount is not number, get from buy price")
	}

	return out
}

func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func DeleteFile(filepath string) {
	if Exist(filepath) {
		os.Remove(filepath)
	}
}
