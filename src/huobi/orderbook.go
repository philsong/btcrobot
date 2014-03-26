/*
  btcrobot is a Bitcoin, Litecoin and Altcoin trading bot written in golang,
  it features multiple trading methods using technical analysis.

  Disclaimer:

  USE AT YOUR OWN RISK!

  The author of this project is NOT responsible for any damage or loss caused
  by this software. There can be bugs and the bot may not perform as expected
  or specified. Please consider testing it first with paper trading /
  backtesting on historical data. Also look at the code to see what how
  it's working.

  Weibo:http://weibo.com/bocaicfa
*/

package huobi

import (
	. "config"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"logger"
	"net/http"
	"strings"
	"time"
)

func (w *Huobi) getOrderBook(symbol string) (ret bool) {
	var huobisymbol string
	if symbol == "btc_cny" {
		huobisymbol = "huobibtccny"
	} else {
		huobisymbol = "huobiltccny"
		logger.Fatal("huobi does not support LTC by now, wait for huobi provide it.", huobisymbol)
		return false
	}

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

				ioutil.WriteFile("cache/OrderBook.json", bodyByte, 0644)
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
				return w.OrderBook(body)
			}
		}
	} else {
		logger.Errorf("HTTP returned status %v", resp)
	}

	logger.Errorln("why in here?")
	return false
}

type SellBuy struct {
	Price  string
	Level  int     ////涨幅
	Amount float64 //成交量
}

type Trade struct {
	Time   string
	Price  float64
	Amount float64
	Type   string
}

type Top_buy_sell struct {
	Price  string
	Level  int
	Amount float64
	Accu   float64
}

type OrderBook struct {
	Sells    [10]SellBuy
	Buys     [10]SellBuy
	Trades   [15]Trade
	P_new    float64
	Level    float64
	Amount   float64
	Total    float64
	Amp      int
	P_open   float64
	P_high   float64
	P_low    float64
	P_last   float64
	Top_sell [5]Top_buy_sell
	Top_buy  [5]Top_buy_sell
}

func (w *Huobi) OrderBook(body string) bool {
	body = strings.TrimPrefix(body, "view_detail(")
	body = strings.TrimSuffix(body, ")")

	ioutil.WriteFile("cache/OrderBook.json", []byte(body), 0644)

	var orderBook OrderBook
	if err := json.Unmarshal([]byte(body), &orderBook); err != nil {
		logger.Infoln(err)
		return false
	}
	logger.Infoln(orderBook)
	return true
}
