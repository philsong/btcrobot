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

package huobiapi

import (
	. "config"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"logger"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

/*
	txcode := map[int]string{
		0:  `买单已委托，<a href="/trade/index.php?a=delegation">查看结果</a>`,
		2:  `没有足够的人民币`,
		10:	`没有足够的比特币`,
		16: `您需要登录才能继续`,
		17: `没有权限`,
		42:	`该委托已经取消，不能修改`,
		44:	`交易价钱太低`,
		56:`卖出价格不能低于限价的95%`}

	logger.Traceln(txcode[m.Code])
*/
type Huobi struct {
	isLogin          bool
	client           *http.Client
	Detail_data      View_detail
	MyTradeInfo      TradeInfo
	latestSolidTrend int
	latestMACDTrend  int
	latestEMATrend   int
	lastLowestprice  float64
	lastBuyprice     float64

	lastAction string

	Disable_trading int

	Peroid   int
	Slippage float64
	xData    []string
	yData    []float64
	Volumn   []float64
}

func NewHuobi() *Huobi {
	w := new(Huobi)
	return w
}

func (w *Huobi) Login() bool {
	login_url := Config["login_url"]
	email := Config["email"]
	password := Config["password"]

	post_arg := url.Values{"email": {email}, "password": {password}}

	//logger.Traceln(strings.NewReader(post_arg.Encode()))
	req, err := http.NewRequest("POST", login_url, strings.NewReader(post_arg.Encode()))
	if err != nil {
		logger.Fatal(err)
		return false
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "https://www.huobi.com/")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36")
	logger.Traceln(req)

	//jar := NewJar()
	jar, _ := cookiejar.New(nil)
	w.client = &http.Client{nil, nil, jar}
	//w.client = new(http.Client)
	resp, err := w.client.Do(req)
	if err != nil {
		logger.Fatal(err)
		return false
	}
	defer resp.Body.Close()
	logger.Tracef("Login resp StatusCode=%v", resp.StatusCode)
	logger.Tracef("Login resp=%v", resp)
	if resp.StatusCode == 200 {
		var body string

		contentEncoding := resp.Header.Get("Content-Encoding")
		logger.Tracef("HTTP returned Content-Encoding %s", contentEncoding)
		switch contentEncoding {
		case "gzip":
			body = DumpGZIP(resp.Body)

		default:
			bodyByte, _ := ioutil.ReadAll(resp.Body)
			body = string(bodyByte)
			ioutil.WriteFile("cache/login.html", bodyByte, os.ModeAppend)
		}

		logger.Traceln(resp.Header.Get("Content-Type"))
		ret := strings.Contains(body, "用户名或者密码错误")
		if ret {
			logger.Traceln("用户名或者密码错误")
			return false
		}

		w.isLogin = true
		return true
	} else if resp.StatusCode == 500 {
		w.isLogin = true
		return true
	} else {
		logger.Tracef("resp %v", resp)
	}

	return false
}

func (w *Huobi) TradeAdd(a, price, amount string) bool {

	if w.isLogin == false {
		if w.Login() == false {
			return false
		}
	}
	/*
		if w.checkAccount(a, price, amount) == false {
			return false
		}
	*/
	post_arg := url.Values{
		"a":      {a},
		"price":  {price},
		"amount": {amount},
	}

	req, err := http.NewRequest("POST", Config["trade_add_url"], strings.NewReader(post_arg.Encode()))
	if err != nil {
		logger.Fatal(err)
		return false
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "https://www.huobi.com/")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	logger.Traceln(req)

	resp, err := w.client.Do(req)
	if err != nil {
		logger.Fatal(err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		var body string

		contentEncoding := resp.Header.Get("Content-Encoding")
		logger.Tracef("HTTP returned Content-Encoding %s", contentEncoding)
		switch contentEncoding {
		case "gzip":
			body = DumpGZIP(resp.Body)

		default:
			bodyByte, _ := ioutil.ReadAll(resp.Body)
			body = string(bodyByte)
			ioutil.WriteFile("cache/TradeAdd.json", bodyByte, os.ModeAppend)
		}

		logger.Traceln(resp.Header.Get("Content-Type"))

		if resp.Header.Get("Content-Type") == "application/json" {
			doc := json.NewDecoder(strings.NewReader(body))

			type Msg struct {
				Code int
				Msg  string
			}

			var m Msg
			if err := doc.Decode(&m); err == io.EOF {
				logger.Traceln(err)
			} else if err != nil {
				logger.Fatal(err)
			}

			logger.Infoln(m)

			if m.Code == 0 {
				return true
			} else {
				return false
			}
		} else {
			ret := strings.Contains(body, "您需要登录才能继续")
			if ret {
				logger.Errorln("您需要登录才能继续")
				w.isLogin = false
				return false
			} else {
				return true
			}
		}

	} else {
		logger.Tracef("HTTP returned status %v", resp)
	}

	return false

}

func (w *Huobi) checkAccount(a, price, amount string) bool {
	btc, cny := w.get_account_info()

	FPrice, err := strconv.ParseFloat(price, 64)
	if err != nil {
		logger.Debugln("price is not float")
		return false
	}
	FAmount, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		logger.Debugln("amount is not float")
		return false
	}
	if a == "do_buy" {
		if float64(cny) < FPrice*FAmount {
			return false
		}
	} else {
		if float64(btc) < FAmount {
			return false
		}
	}

	return true
}

func (w *Huobi) Do_buy(tradePrice, tradeAmount string) bool {

	tradeAPI := NewHuobiTrade(SecretOption["access_key"], SecretOption["secret_key"])
	buyId := tradeAPI.Buy(tradePrice, tradeAmount)
	logger.Infoln("buyId", buyId)
	if buyId != "0" {
		return true
	}

	if w.Login() == true {
		if w.TradeAdd("do_buy", tradePrice, tradeAmount) == true {
			logger.Infoln("执行买入委托成功", tradePrice, tradeAmount)
			return true
		} else {
			logger.Infoln("执行买入失败，原因：买入操作被服务器拒绝", tradePrice, tradeAmount)
		}
	} else {
		logger.Infoln("执行买入失败，原因：登陆失败", tradePrice, tradeAmount)

	}

	return false
}

func (w *Huobi) Do_sell(tradePrice, tradeAmount string) bool {
	tradeAPI := NewHuobiTrade(SecretOption["access_key"], SecretOption["secret_key"])
	sellId := tradeAPI.Sell(tradePrice, tradeAmount)
	logger.Infoln("sellId", sellId)
	if sellId != "0" {
		return true
	}

	if w.Login() == true {
		if w.TradeAdd("do_sell", tradePrice, tradeAmount) == true {
			logger.Infoln("执行卖出委托成功", tradePrice, tradeAmount)
			return true
		} else {
			logger.Infoln("执行卖出失败，原因：卖出操作被服务器拒绝", tradePrice, tradeAmount)
		}
	} else {
		logger.Infoln("执行卖出失败，原因：登陆失败", tradePrice, tradeAmount)
	}

	return false
}

func (w *Huobi) TradeUpdate(id, price, amount string) bool {
	post_arg := url.Values{
		"id":     {id},
		"price":  {price},
		"amount": {amount},
	}

	req, err := http.NewRequest("POST", Config["trade_update_url"], strings.NewReader(post_arg.Encode()))
	if err != nil {
		logger.Fatal(err)
		return false
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "https://www.huobi.com/")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	logger.Traceln(req)

	resp, err := w.client.Do(req)
	if err != nil {
		logger.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		var body string

		contentEncoding := resp.Header.Get("Content-Encoding")
		logger.Tracef("HTTP returned Content-Encoding %s", contentEncoding)
		switch contentEncoding {
		case "gzip":
			body = DumpGZIP(resp.Body)

		default:
			bodyByte, _ := ioutil.ReadAll(resp.Body)
			body = string(bodyByte)
			ioutil.WriteFile("cache/TradeUpdate.json", bodyByte, os.ModeAppend)
		}

		logger.Traceln(resp.Header.Get("Content-Type"))

		if resp.Header.Get("Content-Type") == "application/json" {
			doc := json.NewDecoder(strings.NewReader(body))

			type Msg struct {
				Code int
				Msg  string
			}

			var m Msg
			if err := doc.Decode(&m); err == io.EOF {
				logger.Traceln(err)
			} else if err != nil {
				logger.Fatal(err)
				return false
			}

			logger.Traceln(m)

			if m.Code == 0 {
				return true
			} else {
				return false
			}
		} else {
			ret := strings.Contains(body, "您需要登录才能继续")
			if ret {
				logger.Traceln("您需要登录才能继续")
				w.isLogin = false
				return false
			} else {
				return true
			}
		}

	} else {
		logger.Tracef("HTTP returned status %v", resp)
	}

	return false
}

func (w *Huobi) QueryMyTradeInfo() bool {
	fmt.Println(w.isLogin)
	if w.isLogin == false {
		if w.Login() == false {
			return false
		}
	}

	req, err := http.NewRequest("GET", fmt.Sprintf(Config["my_trade_info"], rand.Float64()), nil)
	if err != nil {
		logger.Fatal(err)
		return false
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", Config["trade_query_url"])
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	logger.Traceln(req)

	resp, err := w.client.Do(req)
	if err != nil {
		logger.Fatal(err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		var body string

		contentEncoding := resp.Header.Get("Content-Encoding")
		logger.Tracef("HTTP returned Content-Encoding %s", contentEncoding)
		switch contentEncoding {
		case "gzip":
			body = DumpGZIP(resp.Body)

		default:
			bodyByte, _ := ioutil.ReadAll(resp.Body)
			body = string(bodyByte)
			ioutil.WriteFile("cache/MyTradeInfo.json", bodyByte, os.ModeAppend)
		}

		logger.Traceln(resp.Header.Get("Content-Type"))

		if resp.Header.Get("Content-Type") == "application/json" {
			w.TradeMyTradeInfoAnalyze(body)
			if w.MyTradeInfo.Code == 0 {
				return true
			} else {
				return false
			}
		} else {
			ret := strings.Contains(body, "您需要登录才能继续")
			if ret {
				logger.Traceln("您需要登录才能继续")
				w.isLogin = false
				return false
			} else {
				return true
			}
		}

	} else {
		logger.Tracef("HTTP returned status %v", resp)
	}

	return false
}

func (w *Huobi) TradeCancel(id string) bool {

	req, err := http.NewRequest("GET", fmt.Sprintf(Config["trade_cancel_url"], id), nil)
	if err != nil {
		logger.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", Config["trade_delegation"])
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	logger.Traceln(req)

	resp, err := w.client.Do(req)
	if err != nil {
		logger.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		var body string

		contentEncoding := resp.Header.Get("Content-Encoding")
		logger.Tracef("HTTP returned Content-Encoding %s", contentEncoding)
		switch contentEncoding {
		case "gzip":
			body = DumpGZIP(resp.Body)

		default:
			bodyByte, _ := ioutil.ReadAll(resp.Body)
			body = string(bodyByte)
			ioutil.WriteFile("cache/TradeCancel.json", bodyByte, os.ModeAppend)
		}

		logger.Traceln(resp.Header.Get("Content-Type"))

		if resp.Header.Get("Content-Type") == "application/json" {
			doc := json.NewDecoder(strings.NewReader(body))

			type Msg struct {
				Code int
				Msg  string
			}

			var m Msg
			if err := doc.Decode(&m); err == io.EOF {
				logger.Traceln(err)
			} else if err != nil {
				logger.Fatal(err)
			}
			logger.Traceln(m)

			if m.Code == 0 {
				return true
			} else {
				return false
			}
		} else {
			ret := strings.Contains(body, "您需要登录才能继续")
			if ret {
				logger.Traceln("您需要登录才能继续")
				w.isLogin = false
				return false
			} else {
				return true
			}
		}

	} else {
		logger.Tracef("HTTP returned status %v", resp)
	}

	return false
}

func (w *Huobi) TradeDelegation() bool {
	fmt.Println(w.isLogin)
	if w.isLogin == false {
		if w.Login() == false {
			return false
		}
	}

	req, err := http.NewRequest("GET", Config["trade_delegation"], nil)
	if err != nil {
		logger.Fatal(err)
		return false
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", Config["trade_add_url"])
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	logger.Traceln(req)

	resp, err := w.client.Do(req)
	if err != nil {
		logger.Fatal(err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		var body string

		contentEncoding := resp.Header.Get("Content-Encoding")
		logger.Tracef("HTTP returned Content-Encoding %s", contentEncoding)
		switch contentEncoding {
		case "gzip":
			body = DumpGZIP(resp.Body)

		default:
			bodyByte, _ := ioutil.ReadAll(resp.Body)
			body = string(bodyByte)
			ioutil.WriteFile("cache/TradeDelegation.html", bodyByte, os.ModeAppend)
		}

		logger.Traceln(resp.Header.Get("Content-Type"))

		if resp.Header.Get("Content-Type") == "application/json" {
			doc := json.NewDecoder(strings.NewReader(body))

			type Msg struct {
				Code int
				Msg  string
			}

			for {
				var m Msg
				if err := doc.Decode(&m); err == io.EOF {
					logger.Traceln(err)
					break
				} else if err != nil {
					logger.Fatal(err)
				}
				logger.Traceln(m)

				if m.Code == 0 {
					return true
				} else {
					return false
				}
			}

			return false
		} else {
			ret := strings.Contains(body, "您需要登录才能继续")
			if ret {
				logger.Traceln("您需要登录才能继续")
				w.isLogin = false
				return false
			} else {
				return w.TradeDelegationAnalyze(body)
			}
		}

	} else {
		logger.Tracef("HTTP returned status %v", resp)
	}

	return false
}

func (w *Huobi) getNewPrice() (float64, bool) {
	if w.TradeDetail() == true {
		logger.Traceln("new：", w.Detail_data.Vp_new)
		logger.Traceln("last：", w.Detail_data.Vp_last)
		logger.Traceln("high：", w.Detail_data.Vp_high)
		logger.Traceln("low：", w.Detail_data.Vp_low)
		logger.Traceln("sell：", w.Detail_data.Vtop_sells[0])
		logger.Traceln("buy：", w.Detail_data.Vtop_buys[0])
		return w.Detail_data.Vp_new, true
	} else {
		logger.Errorln("getNewPrice failed.")
		return 0, false
	}
}

func (w *Huobi) get_account_info() (btc, cny int) {
	if w.QueryMyTradeInfo() == true {
		logger.Infoln("BTC：", w.MyTradeInfo.Extra.Balance.Available_btc_display)
		logger.Infoln("CNY：", w.MyTradeInfo.Extra.Balance.Available_cny_display)
		logger.Infoln("Frozen BTC：", w.MyTradeInfo.Extra.Balance.Frozen_btc_display)
		logger.Infoln("Frozen CNY: ", w.MyTradeInfo.Extra.Balance.Frozen_cny_display)
		logger.Infoln("CNY：", w.MyTradeInfo.Extra.Balance.Available_cny_display)
		logger.Infoln("total", w.MyTradeInfo.Extra.Balance.Total)
		logger.Infoln("(BTC,CNY)", w.MyTradeInfo.Extra.Balance.Available_btc, w.MyTradeInfo.Extra.Balance.Available_cny)
		return w.MyTradeInfo.Extra.Balance.Available_btc, w.MyTradeInfo.Extra.Balance.Available_cny
	} else {
		logger.Errorln("get_account_info failed.")
		return 0, 0
	}
}

func (w *Huobi) TradeDetail() (ret bool) {
	rnd := RandomString(20)

	now := time.Now().UnixNano() / 1000000

	req, err := http.NewRequest("GET", fmt.Sprintf(Config["trade_detail"], rnd, now, now), nil)
	if err != nil {
		logger.Fatal(err)
	}

	req.Header.Set("Referer", Config["base_url"])
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36")
	req.Header.Add("Accept-Encoding", "identity")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Traceln(err)
		logger.Traceln(req)
		logger.Traceln(resp)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		var body string

		contentEncoding := resp.Header.Get("Content-Encoding")
		logger.Tracef("HTTP returned Content-Encoding %s", contentEncoding)
		switch contentEncoding {
		case "gzip":
			body = DumpGZIP(resp.Body)

		default:
			bodyByte, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				logger.Errorln("read the http stream failed")
				return false
			} else {
				body = string(bodyByte)

				ioutil.WriteFile("cache/TradeDetail.json", bodyByte, os.ModeAppend)
			}
		}

		logger.Traceln(resp.Header.Get("Content-Type"))

		if resp.Header.Get("Content-Type") == "application/json" {
			doc := json.NewDecoder(strings.NewReader(body))

			type Msg struct {
				Code int
				Msg  string
			}

			var m Msg
			if err := doc.Decode(&m); err == io.EOF {
				logger.Traceln(err)
			} else if err != nil {
				logger.Fatal(err)
				return false
			}
			logger.Traceln(m)

			if m.Code == 0 {
				return true
			} else {
				logger.Errorln(m)
				return false
			}
		} else {
			ret := strings.Contains(body, "您需要登录才能继续")
			if ret {
				logger.Errorln("您需要登录才能继续")
				logger.Errorln(body)
				return false
			} else {
				return w.TradeDetailAnalyze(body)
			}
		}
	} else {
		logger.Errorf("HTTP returned status %v", resp)
	}

	logger.Errorln("why in here?")
	return false
}

func (w *Huobi) TradeKLinePeroid(peroid int) (ret bool) {
	req, err := http.NewRequest("GET", fmt.Sprintf(Config["trade_kline_url"], peroid, rand.Float64()), nil)
	if err != nil {
		logger.Fatal(err)
		return false
	}

	req.Header.Set("Referer", Config["trade_flash_url"])
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36")

	logger.Traceln(req)

	if w.client == nil {
		w.client = &http.Client{nil, nil, nil}
	}
	resp, err := w.client.Do(req)
	if err != nil {
		logger.Errorln(err)
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		var body string

		contentEncoding := resp.Header.Get("Content-Encoding")
		logger.Tracef("HTTP returned Content-Encoding %s", contentEncoding)
		switch contentEncoding {
		case "gzip":
			body = DumpGZIP(resp.Body)

		default:
			bodyByte, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				logger.Errorln("read the http stream failed")
				return false
			} else {
				body = string(bodyByte)

				ioutil.WriteFile(fmt.Sprintf("cache/TradeKLine_%03d.data", peroid), bodyByte, os.ModeAppend)
			}
		}

		logger.Traceln(resp.Header.Get("Content-Type"))

		ret := strings.Contains(body, "您需要登录才能继续")
		if ret {
			logger.Traceln("您需要登录才能继续")
			return false
		} else {
			return w.AnalyzePeroidLine(fmt.Sprintf("cache/TradeKLine_%03d.data", peroid), body)
		}

	} else {
		logger.Tracef("HTTP returned status %v", resp)
	}

	return false
}

func (w *Huobi) TradeKLineMinute() (ret bool) {
	req, err := http.NewRequest("GET", Config["trade_fenshi"], nil)
	if err != nil {
		logger.Fatal(err)
	}

	req.Header.Set("Referer", Config["trade_flash_url"])
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36")

	logger.Traceln(req)

	if w.client == nil {
		w.client = &http.Client{nil, nil, nil}
	}
	resp, err := w.client.Do(req)
	if err != nil {
		logger.Errorln(err)
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		var body string

		contentEncoding := resp.Header.Get("Content-Encoding")
		logger.Tracef("HTTP returned Content-Encoding %s", contentEncoding)
		switch contentEncoding {
		case "gzip":
			body = DumpGZIP(resp.Body)

		default:
			bodyByte, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				logger.Errorln("read the http stream failed")
				return false
			} else {
				body = string(bodyByte)

				ioutil.WriteFile(fmt.Sprintf("cache/TradeKLine_minute.data"), bodyByte, os.ModeAppend)
			}
		}

		logger.Traceln(resp.Header.Get("Content-Type"))

		ret := strings.Contains(body, "您需要登录才能继续")
		if ret {
			logger.Traceln("您需要登录才能继续")
			return false
		} else {
			return w.AnalyzeMinuteLine(fmt.Sprintf("cache/TradeKLine_minute.data"), body)
		}

	} else {
		logger.Tracef("HTTP returned status %v", resp)
	}

	return false
}
