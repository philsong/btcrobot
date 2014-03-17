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

package controller

import (
	"config"
	"encoding/json"
	"filter"
	"fmt"
	"github.com/studygolang/mux"
	"huobiapi"
	"logger"
	"net/http"
)

// 用户个人首页
// URI: /trade/{username}
func WelcomeHandler(rw http.ResponseWriter, req *http.Request) {
	//util.Redirect(rw, req, "/static/trade")
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/trade/indictor.html")
}

func TradeHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	msgtype := vars["msgtype"]

	if req.Method != "POST" || msgtype == "" {
		// 获取用户信息
		err := config.LoadTrade()
		if err != nil {
			logger.Errorln(err)
			fmt.Fprint(rw, `{"errno": 1, "msg":"`, "读取Trade配置数据失败", `"}`)
			return
		}
		// 设置模板数据
		filter.SetData(req, map[string]interface{}{"trade": config.TradeOption})
		req.Form.Set(filter.CONTENT_TPL_KEY, "/template/trade/trade.html")
		return
	} else if req.Method == "POST" {
		if msgtype == "dobuy" {
			config.TradeOption["buyprice"] = req.FormValue("buyprice")
			config.TradeOption["buyamount"] = req.FormValue("buyamount")
		} else if msgtype == "dosell" {
			config.TradeOption["sellprice"] = req.FormValue("sellprice")
			config.TradeOption["sellamount"] = req.FormValue("sellamount")
		} else {
			fmt.Fprint(rw, `{"errno": 1, "msg":"`, "无效的POST请求", `"}`)
			return
		}

		// 更新个人信息
		err := config.SaveTrade()
		if err != nil {
			fmt.Fprint(rw, `{"errno": 1, "msg":"`, "写入Trade配置数据失败", `"}`)
			return
		}

		huobi := huobiapi.NewHuobi()
		var ret bool
		if msgtype == "dobuy" {
			ret = huobi.BuyIn(config.TradeOption["buyprice"], config.TradeOption["buyamount"])
		} else if msgtype == "dosell" {
			ret = huobi.SellOut(config.TradeOption["sellprice"], config.TradeOption["sellamount"])
		}

		if ret != true {
			fmt.Fprint(rw, `{"errno": 1, "msg":"`, "交易委托失败", `"}`)
			return
		} else {
			fmt.Fprint(rw, `{"errno": 1, "msg":"`, "交易委托成功", `"}`)
			return
		}
	}
}

func EngineHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	msgtype := vars["msgtype"]

	if req.Method != "POST" && msgtype == "" {
		// 获取用户信息
		err := config.LoadOption()
		if err != nil {
			logger.Errorln(err)
			fmt.Fprint(rw, `{"errno": 1, "msg":"`, "读取引擎配置数据失败", `"}`)
			return
		}
		// 设置模板数据
		filter.SetData(req, map[string]interface{}{"config": config.Option})
		req.Form.Set(filter.CONTENT_TPL_KEY, "/template/trade/engine.html")
		return
	} else if req.Method != "POST" && msgtype == "ajax" {
		Option_json, err := json.Marshal(config.Option)
		if err != nil {
			logger.Errorln(err)
			fmt.Fprint(rw, `{"errno": 1, "msg":"`, "读取引擎配置数据失败", `"}`)
		} else {
			fmt.Fprint(rw, string(Option_json))
		}
		return
	} else {
		logger.Debugln("===[", req.FormValue("disable_trading"), "]")
		if req.FormValue("disable_trading") == "on" {
			config.Option["disable_trading"] = "1"
		} else {
			config.Option["disable_trading"] = "0"
		}

		logger.Debugln("===[", req.FormValue("disable_backtesting"), "]")
		if req.FormValue("disable_backtesting") == "on" {
			config.Option["disable_backtesting"] = "1"
		} else {
			config.Option["disable_backtesting"] = "0"
		}

		config.Option["tick_interval"] = req.FormValue("tick_interval")
		config.Option["shortEMA"] = req.FormValue("shortEMA")
		config.Option["longEMA"] = req.FormValue("longEMA")

		config.Option["tradeAmount"] = req.FormValue("tradeAmount")

		config.Option["totalHour"] = req.FormValue("totalHour")

		// 更新个人信息
		err := config.SaveOption()
		if err != nil {
			fmt.Fprint(rw, `{"errno": 1, "msg":"`, "写入引擎配置数据失败", `"}`)
			return
		}

		fmt.Fprint(rw, `{"errno": 0, "msg":"更新引擎配置成功!"}`)
	}
}

func RemindHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	msgtype := vars["msgtype"]
	if req.Method != "POST" && msgtype == "" {
		// 获取用户信息
		err := config.LoadRemind()
		if err != nil {
			fmt.Fprint(rw, `{"errno": 1, "error":"`, "读取提醒配置数据失败", `"}`)
			return
		}
		// 设置模板数据
		filter.SetData(req, map[string]interface{}{"config": config.Remind})
		req.Form.Set(filter.CONTENT_TPL_KEY, "/template/trade/remind.html")
		return
	} else if req.Method != "POST" && msgtype == "ajax" {
		Remind_json, err := json.Marshal(config.Remind)
		if err != nil {
			fmt.Fprint(rw, `{"errno": 1, "msg":"`, "读取提醒配置数据失败", `"}`)
		} else {
			fmt.Fprint(rw, string(Remind_json))
		}
		return
	} else {
		// open传递过来的是“on”或没传递
		if req.FormValue("disable_email") == "on" {
			config.Remind["disable_email"] = "1"
		} else {
			config.Remind["disable_email"] = "0"
		}
		config.Remind["lowest_price"] = req.FormValue("lowest_price")
		config.Remind["highest_price"] = req.FormValue("highest_price")
		config.Remind["to_email"] = req.FormValue("to_email")

		// 更新个人信息
		err := config.SaveRemind()
		if err != nil {
			fmt.Fprint(rw, `{"errno": 1, "error":"`, "写入提醒配置数据失败", `"}`)
			return
		}

		fmt.Fprint(rw, `{"errno": 0, "msg":"更新提醒配置成功!"}`)
	}
}

func SecretHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	msgtype := vars["msgtype"]
	if req.Method != "POST" && msgtype == "" {
		// 获取用户信息
		err := config.LoadSecretOption()
		if err != nil {
			fmt.Fprint(rw, `{"errno": 1, "error":"`, "读取秘钥配置数据失败", `"}`)
			return
		}
		// 设置模板数据
		filter.SetData(req, map[string]interface{}{"config": config.SecretOption})
		req.Form.Set(filter.CONTENT_TPL_KEY, "/template/trade/secret.html")
		return
	} else {
		config.SecretOption["access_key"] = req.FormValue("access_key")
		config.SecretOption["secret_key"] = req.FormValue("secret_key")

		config.SecretOption["email"] = req.FormValue("email")
		config.SecretOption["password"] = req.FormValue("password")

		config.SecretOption["smtp_username"] = req.FormValue("smtp_username")
		config.SecretOption["smtp_password"] = req.FormValue("smtp_password")
		config.SecretOption["smtp_host"] = req.FormValue("smtp_host")
		config.SecretOption["smtp_addr"] = req.FormValue("smtp_addr")

		// 更新个人信息
		err := config.SaveSecretOption()
		if err != nil {
			fmt.Fprint(rw, `{"errno": 1, "error":"`, "写入秘钥配置数据失败", `"}`)
			return
		}

		fmt.Fprint(rw, `{"errno": 0, "msg":"更新秘钥配置成功!"}`)
	}
}
