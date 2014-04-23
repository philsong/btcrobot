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

package webui

import (
	. "common"
	"config"
	"crypto/sha256"
	"crypto/subtle"
	"email"
	"encoding/base64"
	"encoding/json"
	"filter"
	"fmt"
	"github.com/studygolang/mux"
	"huobi"
	"logger"
	"net/http"
	"okcoin"
)

func NotFoundHandler(rw http.ResponseWriter, req *http.Request) {
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/404.html")
}

// SecureCompare performs a constant time compare of two strings to limit timing attacks.
func SecureCompare(given string, actual string) bool {
	givenSha := sha256.Sum256([]byte(given))
	actualSha := sha256.Sum256([]byte(actual))

	return subtle.ConstantTimeCompare(givenSha[:], actualSha[:]) == 1
}

// Basic returns a Handler that authenticates via Basic Auth. Writes a http.StatusUnauthorized
// if authentication fails
func Basic(rw http.ResponseWriter, req *http.Request) bool {
	var siteAuth = base64.StdEncoding.EncodeToString([]byte(config.SecretOption["username"] + ":" + config.SecretOption["password"]))
	auth := req.Header.Get("Authorization")
	if !SecureCompare(auth, "Basic "+siteAuth) {
		logger.Infoln("auth failed")
		rw.Header().Set("WWW-Authenticate", "Basic realm=\"Authorization Required\"")
		http.Error(rw, "请输入用户名密码：默认为admin/123456，登陆后请尽快修改默认密码。", http.StatusUnauthorized)
		return false
	}

	return true
}

func GuideHandler(rw http.ResponseWriter, req *http.Request) {
	if !Basic(rw, req) {
		return
	}
	//util.Redirect(rw, req, "/static/trade")
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/trade/guide.html")
}

// 用户个人首页
// URI: /trade/{username}
func HuobiIndictorHandler(rw http.ResponseWriter, req *http.Request) {
	if !Basic(rw, req) {
		return
	}
	//util.Redirect(rw, req, "/static/trade")
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/trade/indictor_huobi.html")
}

func OkcoinIndictorHandler(rw http.ResponseWriter, req *http.Request) {
	if !Basic(rw, req) {
		return
	}
	//util.Redirect(rw, req, "/static/trade")
	req.Form.Set(filter.CONTENT_TPL_KEY, "/template/trade/indictor_okcoin.html")
}

func TradeHandler(rw http.ResponseWriter, req *http.Request) {
	if !Basic(rw, req) {
		return
	}
	vars := mux.Vars(req)

	msgtype := vars["msgtype"]

	if req.Method != "POST" && msgtype == "" {
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
	} else if req.Method != "POST" && msgtype == "ajax" {
		rw.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域

		Option_json, err := json.Marshal(config.TradeOption)
		if err != nil {
			logger.Errorln(err)
			fmt.Fprint(rw, `{"errno": 1, "msg":"`, "读取引擎配置数据失败", `"}`)
		} else {
			fmt.Fprint(rw, string(Option_json))
		}
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

		var tradeAPI TradeAPI
		if config.Option["tradecenter"] == "huobi" {
			tradeAPI = huobi.NewHuobi()
		} else if config.Option["tradecenter"] == "okcoin" {
			tradeAPI = okcoin.NewOkcoin()
		} else {
			fmt.Fprint(rw, `{"errno": 1, "msg":"`, "没有选择交易所名称", `"}`)
			return
		}

		var ret string
		if msgtype == "dobuy" {
			ret = tradeAPI.Buy(config.TradeOption["buyprice"], config.TradeOption["buyamount"])
		} else if msgtype == "dosell" {
			ret = tradeAPI.Sell(config.TradeOption["sellprice"], config.TradeOption["sellamount"])
		}

		if ret != "0" {
			fmt.Fprint(rw, `{"errno": 1, "msg":"`, "交易委托成功", `"}`)
			return
		} else {
			fmt.Fprint(rw, `{"errno": 1, "msg":"`, "交易委托失败", `"}`)
			return
		}
	}
}

func EngineHandler(rw http.ResponseWriter, req *http.Request) {
	if !Basic(rw, req) {
		return
	}
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
		rw.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域

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

		// open传递过来的是“on”或没传递
		if req.FormValue("disable_email") == "on" {
			config.Option["disable_email"] = "1"
		} else {
			config.Option["disable_email"] = "0"
		}

		config.Option["to_email"] = req.FormValue("to_email")

		logger.Debugln("===[", req.FormValue("disable_backtesting"), "]")
		if req.FormValue("disable_backtesting") == "on" {
			config.Option["disable_backtesting"] = "1"
		} else {
			config.Option["disable_backtesting"] = "0"
		}

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
			fmt.Fprint(rw, `{"errno": 1, "msg":"`, "写入引擎配置数据失败", `"}`)
			return
		}

		fmt.Fprint(rw, `{"errno": 0, "msg":"更新引擎配置成功!"}`)
	}
}

func SecretHandler(rw http.ResponseWriter, req *http.Request) {
	if !Basic(rw, req) {
		return
	}
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
	} else if req.Method != "POST" && msgtype == "ajax" {
		rw.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域

		Option_json, err := json.Marshal(config.SecretOption)
		if err != nil {
			logger.Errorln(err)
			fmt.Fprint(rw, `{"errno": 1, "msg":"`, "读取引擎配置数据失败", `"}`)
		} else {
			fmt.Fprint(rw, string(Option_json))
		}
		return
	} else {
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
			fmt.Fprint(rw, `{"errno": 1, "error":"`, "写入秘钥配置数据失败", `"}`)
			return
		}

		fmt.Fprint(rw, `{"errno": 0, "msg":"更新秘钥配置成功!"}`)

		go email.TriggerTrender("btcrobot测试邮件，您能收到这封邮件说明您的SMTP配置成功，来自星星的机器人")
	}
}
