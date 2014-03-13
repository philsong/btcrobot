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
	. "controller"
	"filter"
	"github.com/studygolang/mux"
)

func initRouter() *mux.Router {

	router := mux.NewRouter()

	// 大部分handler都需要页面展示
	frontViewFilter := filter.NewViewFilter()
	// 表单校验过滤器（配置了验证规则就会执行）
	formValidateFilter := new(filter.FormValidateFilter)
	router.FilterChain(mux.NewFilterChain([]mux.Filter{formValidateFilter, frontViewFilter}...))

	router.HandleFunc("/", WelcomeHandler)
	router.HandleFunc("/secret{json:(|.json)}", SecretHandler)
	router.HandleFunc("/engine{msgtype:(|ajax|get|post)}{json:(|.json)}", EngineHandler)
	router.HandleFunc("/trade{msgtype:(|ajax|dobuy|dosell)}{json:(|.json)}", TradeHandler)
	router.HandleFunc("/remind{msgtype:(|ajax|get|post)}{json:(|.json)}", RemindHandler)

	// 错误处理handler
	router.HandleFunc("/noauthorize", NoAuthorizeHandler) // 无权限handler
	// 404页面
	router.HandleFunc("/{*}", NotFoundHandler)

	return router
}
