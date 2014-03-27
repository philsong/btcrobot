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
	. "controller"
	"filter"
	"fmt"
	"github.com/studygolang/mux"
	"logger"
	"net/http"
)

func startWEBserver() {
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
	fmt.Printf("start web server failed, please check if %s is already used.", Config["host"])
	fmt.Println(" <<<----------------------------------------------------------] ")
}

func initRouter() *mux.Router {

	router := mux.NewRouter()

	// 大部分handler都需要页面展示
	frontViewFilter := filter.NewViewFilter()
	// 表单校验过滤器（配置了验证规则就会执行）
	formValidateFilter := new(filter.FormValidateFilter)
	router.FilterChain(mux.NewFilterChain([]mux.Filter{formValidateFilter, frontViewFilter}...))

	router.HandleFunc("/", IndictorHandler)
	router.HandleFunc("/secret{json:(|.json)}", SecretHandler)
	router.HandleFunc("/engine{msgtype:(|ajax|get|post)}{json:(|.json)}", EngineHandler)
	router.HandleFunc("/trade{msgtype:(|ajax|dobuy|dosell)}{json:(|.json)}", TradeHandler)

	// 404页面
	router.HandleFunc("/{*}", NotFoundHandler)

	return router
}
