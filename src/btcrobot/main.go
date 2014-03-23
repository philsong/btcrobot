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
	"logger"
	"math/rand"
	"net/http"
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

	version := "0.23"
	fmt.Println("[ ---------------------------------------------------------->>> ")
	fmt.Println(" BTC robot version ", version)
	fmt.Println(" *BTC操盘手自动化交易引擎*")
	fmt.Println(" btcrobot is a Bitcoin, Litecoin and Altcoin trading bot written in golang")
	fmt.Println(" it features multiple trading methods using technical analysis.")
	fmt.Println(" ")
	fmt.Println(" Disclaimer:")
	fmt.Println(" ")
	fmt.Println(" USE AT YOUR OWN RISK!")
	fmt.Println("")
	fmt.Println(" The author of this project is NOT responsible for any damage or loss caused")
	fmt.Println(" by this software. There can be bugs and the bot may not perform as expected")
	fmt.Println(" or specified. Please consider testing it first with paper trading /")
	fmt.Println(" backtesting on historical data. Also understand what how it's working.")
	fmt.Println("")
	fmt.Println(" *@author [btcrobot]")
	fmt.Println(" *@feedback http://weibo.com/bocaicfa")

	fmt.Println(" *@Open http://127.0.0.1:9090 in browser to config the robot")
	fmt.Println(" *@请在浏览器中打开 http://127.0.0.1:9090 配置相关参数")
	fmt.Println(" *@警告：API key和密码存放在conf/secret.json文件内，共享给他人前请务必删除，注意账号安全！！")
	fmt.Println(" <<<----------------------------------------------------------] ")
	SavePid()

	go entry.RunRobot()

	// 服务静态文件
	http.Handle("/static/", http.FileServer(http.Dir(ROOT)))

	router := initRouter()
	http.Handle("/", router)
	if Config["env"] == "test" {
		logger.Infoln(http.ListenAndServe("0.0.0.0:9090", nil))
	} else {
		logger.Infoln(http.ListenAndServe(Config["host"], nil))
	}

	fmt.Println("[ ---------------------------------------------------------->>> ")
	fmt.Println("start web server failed, please check if the port 9090 is already used.")
	fmt.Println(" <<<----------------------------------------------------------] ")

	oneHour := 60 * 60 * 1000 * time.Millisecond

	time.Sleep(24 * oneHour)
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
