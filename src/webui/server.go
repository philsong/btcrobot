package webui

import (
	. "common"
	"config"
	"email"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/martini-contrib/auth"
	"github.com/go-martini/martini"
	"huobi"
	"logger"
	"net/http"
	"okcoin"
)

func webui() {
	m := martini.Classic()
	m.Get("/secret", func() string {
		// show something
		err := config.LoadSecretOption()
		if err != nil {
			logger.Errorln(err)
			return `{"errno": 1, "msg":", "读取秘钥配置数据失败"}`
		}
		Option_json, err := json.Marshal(config.SecretOption)
		if err != nil {
			logger.Errorln(err)
			return `{"errno": 1, "msg":", "解析秘钥配置数据失败"}`
		}

		return string(Option_json)
	})

	m.Post("/secret", func(res http.ResponseWriter, req *http.Request) { // res and req are injected by Martini
		config.SecretOption["username"] = req.FormValue("username")
		config.SecretOption["password"] = req.FormValue("password")

		config.SecretOption["huobi_access_key"] = req.FormValue("huobi_access_key")
		config.SecretOption["huobi_secret_key"] = req.FormValue("huobi_secret_key")

		config.SecretOption["ok_partner"] = req.FormValue("ok_partner")
		config.SecretOption["ok_secret_key"] = req.FormValue("ok_secret_key")

		config.SecretOption["smtp_username"] = req.FormValue("smtp_username")
		config.SecretOption["smtp_password"] = req.FormValue("smtp_password")
		config.SecretOption["smtp_host"] = req.FormValue("smtp_host")
		config.SecretOption["smtp_addr"] = req.FormValue("smtp_addr")

		// 更新个人信息
		err := config.SaveSecretOption()
		if err != nil {
			fmt.Fprint(res, "写入秘钥配置数据失败")
			return
		}

		fmt.Fprint(res, "更新秘钥配置成功!")

		go email.TriggerTrender("btcrobot测试邮件，您能收到这封邮件说明您的SMTP配置成功，来自星星的机器人")

		return
	})

	m.Get("/engine", func() string {
		// show something
		err := config.LoadOption()
		if err != nil {
			logger.Errorln(err)
			return `{"errno": 1, "msg":", "读取引擎配置数据失败"}`
		}
		Option_json, err := json.Marshal(config.Option)
		if err != nil {
			logger.Errorln(err)
			return `{"errno": 1, "msg":", "解析引擎配置数据失败"}`
		}

		return string(Option_json)
	})

	m.Post("/engine", func(res http.ResponseWriter, req *http.Request) {
		if req.FormValue("enable_trading") == "on" {
			config.Option["enable_trading"] = "1"
		} else {
			config.Option["enable_trading"] = "0"
		}

		// open传递过来的是“on”或没传递
		if req.FormValue("enable_email") == "on" {
			config.Option["enable_email"] = "1"
		} else {
			config.Option["enable_email"] = "0"
		}

		config.Option["to_email"] = req.FormValue("to_email")

		config.Option["tick_interval"] = req.FormValue("tick_interval")
		config.Option["datacenter"] = req.FormValue("datacenter")
		config.Option["tradecenter"] = req.FormValue("tradecenter")
		config.Option["symbol"] = req.FormValue("symbol")
		config.Option["strategy"] = req.FormValue("strategy")
		config.Option["shortEMA"] = req.FormValue("shortEMA")
		config.Option["longEMA"] = req.FormValue("longEMA")
		config.Option["signalPeriod"] = req.FormValue("signalPeriod")

		config.Option["tradeAmount"] = req.FormValue("tradeAmount")
		config.Option["slippage"] = req.FormValue("slippage")

		config.Option["totalHour"] = req.FormValue("totalHour")
		config.Option["buyThreshold"] = req.FormValue("buyThreshold")
		config.Option["sellThreshold"] = req.FormValue("sellThreshold")
		config.Option["MACDbuyThreshold"] = req.FormValue("MACDbuyThreshold")
		config.Option["MACDsellThreshold"] = req.FormValue("MACDsellThreshold")

		config.Option["stoploss"] = req.FormValue("stoploss")

		// 更新个人信息
		err := config.SaveOption()
		if err != nil {
			fmt.Fprint(res, "写入引擎配置数据失败")
			return
		}

		fmt.Fprint(res, "更新引擎配置成功!")
	})

	m.Get("/trade", func() string {
		// show something
		config.LoadTrade()

		err := config.LoadTrade()
		if err != nil {
			logger.Errorln(err)
			return `{"errno": 1, "msg":", "读取Trade配置数据失败"}`
		}

		Option_json, err := json.Marshal(config.TradeOption)
		if err != nil {
			logger.Errorln(err)
			return `{"errno": 1, "msg":", "解析引擎配置数据失败"}`
		}

		return string(Option_json)
	})

	m.Post("/trade", func(res http.ResponseWriter, req *http.Request) {
		msgtype := req.FormValue("msgtype")
		if msgtype == "dobuy" {
			config.TradeOption["buyprice"] = req.FormValue("buyprice")
			config.TradeOption["buyamount"] = req.FormValue("buyamount")
		} else if msgtype == "dosell" {
			config.TradeOption["sellprice"] = req.FormValue("sellprice")
			config.TradeOption["sellamount"] = req.FormValue("sellamount")
		} else {
			fmt.Fprint(res, "无效的POST请求")
		}

		// 更新个人信息
		err := config.SaveTrade()
		if err != nil {
			fmt.Fprint(res, "写入Trade配置数据失败")
		}

		var tradeAPI TradeAPI
		if config.Option["tradecenter"] == "huobi" {
			tradeAPI = huobi.NewHuobi()
		} else if config.Option["tradecenter"] == "okcoin" {
			tradeAPI = okcoin.NewOkcoin()
		} else {
			fmt.Fprint(res, "没有选择交易所名称")
		}

		var ret string
		if msgtype == "dobuy" {
			ret = tradeAPI.Buy(config.TradeOption["buyprice"], config.TradeOption["buyamount"])
		} else if msgtype == "dosell" {
			ret = tradeAPI.Sell(config.TradeOption["sellprice"], config.TradeOption["sellamount"])
		}

		if ret != "0" {
			fmt.Fprint(res, "交易委托成功")
		} else {
			fmt.Fprint(res, "交易委托失败")
		}
	})

	m.Use(auth.Basic(config.SecretOption["username"], config.SecretOption["password"]))
	m.Use(martini.Static("./static"))

	logger.Infoln(http.ListenAndServe(config.Config["host"], m))

	m.Run()

	fmt.Println("[ ---------------------------------------------------------->>> ")
	fmt.Printf("start web server failed, please check if %s is already used.", config.Config["host"])
	fmt.Println(" <<<----------------------------------------------------------] ")
}

func RunServer() {
	webui()
}
