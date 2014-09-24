package webui

import (
	. "common"
	. "config"
	"email"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/martini-contrib/auth"
	"github.com/philsong/martini"
	"huobi"
	"logger"
	"net/http"
	"okcoin"
	"strconv"
)

func webui() {
	m := martini.Classic()

	m.Get("/", func(res http.ResponseWriter, req *http.Request) { // res and req are injected by Martini
		http.Redirect(res, req, "/index.html", http.StatusFound)
	})

	m.Get("/secret", func() string {
		// show something
		err := LoadSecretOption()
		if err != nil {
			logger.Errorln(err)
			return `{"errno": 1, "msg":", "读取秘钥配置数据失败"}`
		}
		Option_json, err := json.Marshal(SecretOption)
		if err != nil {
			logger.Errorln(err)
			return `{"errno": 1, "msg":", "解析秘钥配置数据失败"}`
		}

		return string(Option_json)
	})

	m.Post("/secret", func(res http.ResponseWriter, req *http.Request) { // res and req are injected by Martini
		SecretOption["username"] = req.FormValue("username")
		SecretOption["password"] = req.FormValue("password")

		SecretOption["bitvc_email"] = req.FormValue("bitvc_email")
		SecretOption["bitvc_password"] = req.FormValue("bitvc_password")

		SecretOption["huobi_access_key"] = req.FormValue("huobi_access_key")
		SecretOption["huobi_secret_key"] = req.FormValue("huobi_secret_key")

		SecretOption["smtp_username"] = req.FormValue("smtp_username")
		SecretOption["smtp_password"] = req.FormValue("smtp_password")
		SecretOption["smtp_host"] = req.FormValue("smtp_host")
		SecretOption["smtp_addr"] = req.FormValue("smtp_addr")

		SecretOption["OKCoinAPIkey"] = req.FormValue("OKCoinAPIkey")

		SecretOption["ok_partner"] = req.FormValue("ok_partner")
		SecretOption["ok_secret_key"] = req.FormValue("ok_secret_key")

		SecretOption["ok_partner1"] = req.FormValue("ok_partner1")
		SecretOption["ok_secret_key1"] = req.FormValue("ok_secret_key1")

		SecretOption["ok_partner2"] = req.FormValue("ok_partner2")
		SecretOption["ok_secret_key2"] = req.FormValue("ok_secret_key2")

		SecretOption["ok_partner3"] = req.FormValue("ok_partner3")
		SecretOption["ok_secret_key3"] = req.FormValue("ok_secret_key3")

		SecretOption["ok_partner4"] = req.FormValue("ok_partner4")
		SecretOption["ok_secret_key4"] = req.FormValue("ok_secret_key4")

		SecretOption["ok_partner5"] = req.FormValue("ok_partner5")
		SecretOption["ok_secret_key5"] = req.FormValue("ok_secret_key5")

		SecretOption["ok_partner6"] = req.FormValue("ok_partner6")
		SecretOption["ok_secret_key6"] = req.FormValue("ok_secret_key6")

		SecretOption["ok_partner7"] = req.FormValue("ok_partner7")
		SecretOption["ok_secret_key7"] = req.FormValue("ok_secret_key7")

		SecretOption["ok_partner8"] = req.FormValue("ok_partner8")
		SecretOption["ok_secret_key8"] = req.FormValue("ok_secret_key8")

		SecretOption["ok_partner9"] = req.FormValue("ok_partner9")
		SecretOption["ok_secret_key9"] = req.FormValue("ok_secret_key9")

		SecretOption["ok_partner10"] = req.FormValue("ok_partner10")
		SecretOption["ok_secret_key10"] = req.FormValue("ok_secret_key10")

		// 更新个人信息
		err := SaveSecretOption()
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
		err := LoadOption()
		if err != nil {
			logger.Errorln(err)
			return `{"errno": 1, "msg":", "读取引擎配置数据失败"}`
		}
		Option_json, err := json.Marshal(Option)
		if err != nil {
			logger.Errorln(err)
			return `{"errno": 1, "msg":", "解析引擎配置数据失败"}`
		}

		return string(Option_json)
	})

	m.Post("/engine", func(res http.ResponseWriter, req *http.Request) {
		if req.FormValue("enable_trading") == "on" {
			Option["enable_trading"] = "1"
		} else {
			Option["enable_trading"] = "0"
		}

		if req.FormValue("discipleMode") == "on" {
			Option["discipleMode"] = "1"
		} else {
			Option["discipleMode"] = "0"
		}

		Option["discipleValue"] = req.FormValue("discipleValue")

		// open传递过来的是“on”或没传递
		if req.FormValue("enable_email") == "on" {
			Option["enable_email"] = "1"
		} else {
			Option["enable_email"] = "0"
		}

		Option["to_email"] = req.FormValue("to_email")

		Option["tick_interval"] = req.FormValue("tick_interval")
		Option["datacenter"] = req.FormValue("datacenter")
		Option["tradecenter"] = req.FormValue("tradecenter")
		Option["symbol"] = req.FormValue("symbol")
		Option["strategy"] = req.FormValue("strategy")
		Option["shortEMA"] = req.FormValue("shortEMA")
		Option["longEMA"] = req.FormValue("longEMA")
		Option["signalPeriod"] = req.FormValue("signalPeriod")

		Option["tradeAmount"] = req.FormValue("tradeAmount")
		Option["slippage"] = req.FormValue("slippage")

		Option["totalHour"] = req.FormValue("totalHour")
		Option["buyThreshold"] = req.FormValue("buyThreshold")
		Option["sellThreshold"] = req.FormValue("sellThreshold")
		Option["MACDbuyThreshold"] = req.FormValue("MACDbuyThreshold")
		Option["MACDsellThreshold"] = req.FormValue("MACDsellThreshold")

		Option["basePrice"] = req.FormValue("basePrice")
		Option["fluctuation"] = req.FormValue("fluctuation")

		Option["stoploss"] = req.FormValue("stoploss")

		// 更新个人信息
		err := SaveOption()
		if err != nil {
			fmt.Fprint(res, "写入引擎配置数据失败")
			return
		}

		fmt.Fprint(res, "更新引擎配置成功!")
	})

	m.Get("/trade", func() string {
		// show something
		LoadTrade()

		err := LoadTrade()
		if err != nil {
			logger.Errorln(err)
			return `{"errno": 1, "msg":", "读取Trade配置数据失败"}`
		}

		Option_json, err := json.Marshal(TradeOption)
		if err != nil {
			logger.Errorln(err)
			return `{"errno": 1, "msg":", "解析引擎配置数据失败"}`
		}

		return string(Option_json)
	})

	m.Post("/trade", func(res http.ResponseWriter, req *http.Request) {
		msgtype := req.FormValue("msgtype")
		var nbuyprice, nbuytotalamount, nbuyinterval, nmaxbuyamountratio float64
		var nbuytimes int
		var nsellprice, nselltotalamount, nsellinterval, nmaxsellamountratio float64
		var nselltimes int

		var err error
		if msgtype == "dobuy" {
			TradeOption["buyprice"] = req.FormValue("buyprice")
			TradeOption["buytotalamount"] = req.FormValue("buytotalamount")
			TradeOption["buyinterval"] = req.FormValue("buyinterval")
			TradeOption["buytimes"] = req.FormValue("buytimes")
			TradeOption["maxbuyamountratio"] = req.FormValue("maxbuyamountratio")
			nbuyprice, err = strconv.ParseFloat(Option["buyprice"], 64)
			if err != nil {
				logger.Errorln("config item buyprice is not float")
				return
			}
			nbuytotalamount, err = strconv.ParseFloat(Option["buytotalamount"], 64)
			if err != nil {
				logger.Errorln("config item numbuytotalamount is not float")
				return
			}
			nbuyinterval, err = strconv.ParseFloat(Option["buyinterval"], 64)
			if err != nil {
				logger.Errorln("config item buyinterval is not float")
				return
			}
			nbuytimes, err = strconv.Atoi(Option["buytimes"])
			if err != nil {
				logger.Errorln("config item numbuytimes is not float")
				return
			}

			nmaxbuyamountratio, err = strconv.ParseFloat(Option["maxbuyamountratio"], 64)
			if err != nil {
				logger.Errorln("config item numbuytotalamount is not float")
				return
			}
		} else if msgtype == "dosell" {
			TradeOption["sellprice"] = req.FormValue("sellprice")
			TradeOption["selltotalamount"] = req.FormValue("selltotalamount")
			TradeOption["sellinterval"] = req.FormValue("sellinterval")
			TradeOption["selltimes"] = req.FormValue("selltimes")
			TradeOption["maxsellamountratio"] = req.FormValue("maxsellamountratio")

			nsellprice, err = strconv.ParseFloat(Option["sellprice"], 64)
			if err != nil {
				logger.Errorln("config item sellprice is not float")
				return
			}
			nselltotalamount, err = strconv.ParseFloat(Option["selltotalamount"], 64)
			if err != nil {
				logger.Errorln("config item selltotalamount is not float")
				return
			}
			nsellinterval, err = strconv.ParseFloat(Option["v"], 64)
			if err != nil {
				logger.Errorln("config item sellinterval is not float")
				return
			}
			nselltimes, err = strconv.Atoi(Option["selltimes"])
			if err != nil {
				logger.Errorln("config item selltimes is not float")
				return
			}

			nmaxsellamountratio, err = strconv.ParseFloat(Option["maxsellamountratio"], 64)
			if err != nil {
				logger.Errorln("config item maxsellamountratio is not float")
				return
			}

		} else {
			fmt.Fprint(res, "无效的POST请求")
		}

		// 更新个人信息
		err = SaveTrade()
		if err != nil {
			fmt.Fprint(res, "写入Trade配置数据失败")
		}

		var tradeAPI TradeAPI
		if Option["tradecenter"] == "huobi" {
			tradeAPI = huobi.NewHuobi()
		} else if Option["tradecenter"] == "okcoin" {
			tradeAPI = okcoin.NewOkcoin()
		} else {
			fmt.Fprint(res, "没有选择交易所名称")
		}

		var ret string

		var arrbuyTradePrice []float64
		splitnbuyinterval := nbuyinterval / float64(nbuytimes)
		splitbuyTradeAmount := nbuytotalamount / float64(nbuytimes)
		if splitbuyTradeAmount/nbuytotalamount > nmaxbuyamountratio {
			return
		}

		if msgtype == "dobuy" {
			for i := 0; i < nbuytimes; i++ {
				warning := "oo, 买入buy In<----限价单"

				if i < nbuytimes/2 {
					arrbuyTradePrice[i] = nbuyprice - float64(i)*splitnbuyinterval
				} else {
					arrbuyTradePrice[i] = nbuyprice + float64((i-nbuytimes/2))*splitnbuyinterval
				}

				tradePrice := fmt.Sprintf("%f", arrbuyTradePrice[i])
				strsplitTradeAmount := fmt.Sprintf("%f", splitbuyTradeAmount)

				ret = tradeAPI.Buy(tradePrice, strsplitTradeAmount)

				if ret != "0" {
					fmt.Fprint(res, "交易委托成功")
				} else {
					fmt.Fprint(res, "交易委托失败")
				}
				logger.Infoln(warning)
			}
		}

		var arrsellTradePrice []float64
		splitnsellinterval := nsellinterval / float64(nselltimes)
		splitsellTradeAmount := nselltotalamount / float64(nselltimes)
		if splitsellTradeAmount/nselltotalamount > nmaxsellamountratio {
			return
		}
		if msgtype == "dosell" {
			for i := 0; i < nselltimes; i++ {
				warning := "oo, 卖出buy In<----限价单"

				if i < nselltimes/2 {
					arrsellTradePrice[i] = nsellprice - float64(i)*splitnsellinterval
				} else {
					arrsellTradePrice[i] = nsellprice + float64((i-nselltimes/2))*splitnsellinterval
				}

				tradePrice := fmt.Sprintf("%f", arrsellTradePrice[i])
				strsplitTradeAmount := fmt.Sprintf("%f", splitbuyTradeAmount)

				ret = tradeAPI.Sell(tradePrice, strsplitTradeAmount)

				if ret != "0" {
					fmt.Fprint(res, "交易委托成功")
				} else {
					fmt.Fprint(res, "交易委托失败")
				}
				logger.Infoln(warning)
			}
		}
	})

	m.Use(auth.Basic(SecretOption["username"], SecretOption["password"]))
	m.Use(martini.Static("./static"))
	m.Use(martini.Static("../static"))

	logger.Infoln(http.ListenAndServe(Config["host"], m))

	m.Run()

	fmt.Println("[ ---------------------------------------------------------->>> ")
	fmt.Printf("start web server failed, please check if %s is already used.", Config["host"])
	fmt.Println(" <<<----------------------------------------------------------] ")
}

func RunServer() {
	webui()
}
