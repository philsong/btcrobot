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
	. "config"
	"entry"
	"fmt"
	"math/rand"
	"path/filepath"
	"process"
	"runtime"
	"time"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// 设置随机数种子
	rand.Seed(time.Now().Unix())
}

func main() {
	SavePid()

	printBanner()

	go startWEBserver()

	entry.RunRobot()
}

func printBanner() {
	version := "0.28"
	fmt.Println("[ ---------------------------------------------------------->>> ")
	fmt.Println(" BTC/LTC robot version ", version)
	fmt.Println(" *BTC/LTC操盘手自动化交易引擎*")
	fmt.Println(" btcrobot is a Bitcoin, Litecoin and Altcoin trading bot written in golang")
	fmt.Println(" it features multiple trading methods using technical analysis.")
	fmt.Println(" ")
	fmt.Println(" Disclaimer: USE AT YOUR OWN RISK!")
	fmt.Println(" ")
	fmt.Println(" 理性投机，风控第一.")
	fmt.Println(" ")
	fmt.Println(" *@feedback [btcrobot]http://weibo.com/bocaicfa")
	fmt.Println(" *@请在浏览器中打开 http://127.0.0.1:9090 配置相关参数")
	fmt.Println(" *@警告：API key和密码存放在conf/secret.json文件内，共享给他人前请务必删除，注意账号安全！！")
	fmt.Println(" <<<----------------------------------------------------------] ")
}

// 保存PID
func SavePid() {
	pidFile := Config["pid"]
	if !filepath.IsAbs(Config["pid"]) {
		pidFile = ROOT + "/" + pidFile
	}
	// TODO：错误不处理
	process.SavePidTo(pidFile)
}
