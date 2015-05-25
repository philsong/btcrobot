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

package main

import (
	"config"
	"fmt"
	"io/ioutil"
	"math/rand"
	"monitor"
	"os"
	"path/filepath"
	"runtime"
	// "runtime"
	"strconv"
	"time"
	"webui"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// 设置随机数种子
	rand.Seed(time.Now().Unix())
}

func main() {
	SavePid()

	printBanner()

	go webui.RunServer()

	monitor.RunRobot()
}

func printBanner() {
	version := "V0.5"
	fmt.Println("[ ---------------------------------------------------------->>> ")
	fmt.Println(" BTC/LTC自动化算法交易引擎", version)
	fmt.Println(" btcrobot is a Bitcoin, Litecoin and Altcoin trading bot written in golang")
	fmt.Println(" it features multiple trading methods using technical analysis.")
	fmt.Println(" ")
	fmt.Println(" *@请在浏览器中打开 http://127.0.0.1:9090 配置相关参数")
	fmt.Println(" *@警告：API key和密码存放在conf/secret.json文件内，共享给他人前请务必删除，注意账号安全！！")
	fmt.Println(" <<<----------------------------------------------------------] ")
}

// 保存PID
func SavePid() {
	pidFile := config.Config["pid"]
	if pidFile == "" {
		pidFile = "pid/btcrobot.pid"
	}

	if !filepath.IsAbs(config.Config["pid"]) {
		pidFile = config.ROOT + "/" + pidFile
	}

	// 保存pid
	pidPath := filepath.Dir(pidFile)
	if err := os.MkdirAll(pidPath, 0777); err != nil {
		return
	}
	ioutil.WriteFile(pidFile, []byte(strconv.Itoa(os.Getpid())), 0777)
}
