/*
  btcbot is a Bitcoin trading bot for HUOBI.com written
  in golang, it features multiple trading methods using
  technical analysis.

  Disclaimer:

  USE AT YOUR OWN RISK!

  The author of this project is NOT responsible for any damage or loss caused
  by this software. There can be bugs and the bot may not perform as expected
  or specified. Please consider testing it first with paper trading /
  backtesting on historical data. Also look at the code to see what how
  it's working.

  Weibo:http://weibo.com/bocaicfa
*/

package process

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

// 保存pid
func SavePidTo(pidFile string) error {
	pidPath := filepath.Dir(pidFile)
	if err := os.MkdirAll(pidPath, 0777); err != nil {
		return err
	}
	return ioutil.WriteFile(pidFile, []byte(strconv.Itoa(os.Getpid())), 0777)
}

// 获得可执行程序所在目录
func ExecutableDir() (string, error) {
	pathAbs, err := filepath.Abs(os.Args[0])
	if err != nil {
		return "", err
	}
	return filepath.Dir(pathAbs), nil
}
