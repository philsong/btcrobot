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
	. "common"
	. "config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"logger"
	"net/http"
	"strings"
	"time"
	"util"
)

func (w *Huobi) getOrderBook(symbol string) (ret bool, hbOrderBook OrderBook) {
	// init to false
	ret = false
	var huobisymbol string
	if symbol == "btc_cny" {
		huobisymbol = "huobibtccny"
	} else {
		huobisymbol = "huobiltccny"
		logger.Fatal("huobi does not support LTC by now, wait for huobi provide it.", huobisymbol)
		return
	}

	rnd := util.RandomString(20)

	now := time.Now().UnixNano() / 1000000

	req, err := http.NewRequest("GET", fmt.Sprintf(Config["hb_trade_detail"], rnd, now, now), nil)
	if err != nil {
		logger.Fatal(err)
		return
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
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		logger.Errorf("HTTP returned status %v", resp)
		return
	}

	var body string
	contentEncoding := resp.Header.Get("Content-Encoding")
	logger.Tracef("HTTP returned Content-Encoding %s", contentEncoding)
	switch contentEncoding {
	case "gzip":
		body = util.DumpGZIP(resp.Body)
	default:
		bodyByte, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Errorln("read the http stream failed")
			return
		} else {
			body = string(bodyByte)
			ioutil.WriteFile("cache/OrderBook.json", bodyByte, 0644)
		}
	}

	return w.analyzeOrderBook(body)
}

type SellBuy struct {
	Price  float64
	Level  float64 // 涨幅
	Amount float64 // 成交量
}

type Trade struct {
	Time   string
	Price  float64
	Amount float64
	Type   string
}

type Top_buy_sell struct {
	Price  float64
	Level  float64
	Amount float64
	Accu   float64
}

type HBOrderBook struct {
	Sells    [10]SellBuy
	Buys     [10]SellBuy
	Trades   [60]Trade
	P_new    float64
	Level    float64
	Amount   float64
	Total    float64
	Amp      float64
	P_open   float64
	P_high   float64
	P_low    float64
	P_last   float64
	Top_sell [5]Top_buy_sell
	Top_buy  [5]Top_buy_sell
}

func (w *Huobi) analyzeOrderBook2(body string) (ret bool, orderBook OrderBook) {
	// init to false
	ret = false
	body = strings.TrimPrefix(body, "view_detail(")
	body = strings.TrimSuffix(body, ")")

	ioutil.WriteFile("cache/OrderBook.json", []byte(body), 0644)

	var hbOrderBook HBOrderBook
	if err := json.Unmarshal([]byte(body), &hbOrderBook); err != nil {
		logger.Infoln(err)
		return
	}
	logger.Infoln(hbOrderBook)

	ret = true
	return
}

func (w *Huobi) analyzeOrderBook(body string) (ret bool, orderBook OrderBook) {
	logger.Debugln("analyzeOrderBook start....")
	// init to false
	ret = false

	ioutil.WriteFile("cache/view_detail.jsonp", []byte(body), 0644)

	body = strings.TrimPrefix(body, "view_detail(")
	body = strings.TrimSuffix(body, ")")

	ioutil.WriteFile("cache/analyzeOrderBook.json", []byte(body), 0644)

	var hbOrderBook HBOrderBook
	var view_detail map[string]interface{}

	logger.Debugln("Unmarshal begin")
	if err := json.Unmarshal([]byte(body), &view_detail); err != nil {
		logger.Errorln("analyzeOrderBook Unmarshal failed")
		return
	}
	logger.Debugln("Unmarshal end")

	p_new := view_detail["p_new"].(float64)
	level := view_detail["level"].(float64)
	amount := view_detail["amount"].(float64)
	total := view_detail["total"].(float64)
	amp := view_detail["amp"].(float64)
	p_open := view_detail["p_open"].(float64)
	p_high := view_detail["p_high"].(float64)
	p_low := view_detail["p_low"].(float64)
	p_last := view_detail["p_last"].(float64)

	hbOrderBook.P_new = p_new
	hbOrderBook.Level = level
	hbOrderBook.Amount = amount
	hbOrderBook.Total = total
	hbOrderBook.Amp = amp
	hbOrderBook.P_open = p_open
	hbOrderBook.P_high = p_high
	hbOrderBook.P_low = p_low
	hbOrderBook.P_last = p_last

	logger.Debugln("analyzeOrderBook sells....")
	sells := view_detail["sells"].([]interface{})
	ret = parse_buy_sell(sells, &hbOrderBook.Sells)
	if !ret {
		return
	}

	logger.Debugln("analyzeOrderBook buys....")
	buys := view_detail["buys"].([]interface{})
	ret = parse_buy_sell(buys, &hbOrderBook.Buys)
	if !ret {
		return
	}

	logger.Debugln("analyzeOrderBook trades....")
	trades := view_detail["trades"].([]interface{})
	ret = parse_trade(trades, &hbOrderBook.Trades)
	if !ret {
		return
	}

	logger.Debugln("analyzeOrderBook top_buy....")
	top_buys := view_detail["top_buy"].([]interface{})
	ret = parse_topbuy(top_buys, &hbOrderBook.Top_buy)
	if !ret {
		return
	}

	logger.Debugln("analyzeOrderBook top_sell....")
	top_sells := view_detail["top_sell"].([]interface{})
	ret = parse_topbuy(top_sells, &hbOrderBook.Top_sell)
	if !ret {
		return
	}

	logger.Debugln("analyzeOrderBook end....")
	logger.Debugln(hbOrderBook)

	for i := 0; i < 10; i++ {
		orderBook.Asks[i].Price = hbOrderBook.Sells[9-i].Price
		orderBook.Asks[i].Amount = hbOrderBook.Sells[9-i].Amount
		orderBook.Bids[i].Price = hbOrderBook.Buys[i].Price
		orderBook.Bids[i].Amount = hbOrderBook.Buys[i].Amount
	}

	ret = true
	return
}

func parse_trade(trades []interface{}, trades_data *[60]Trade) bool {
	for k, v := range trades {
		switch vt := v.(type) {
		case map[string]interface{}:
			logger.Debugln(k, " is a map:")
			logger.Debugf("trades[%d]\n", k)

			for ik, iv := range vt {
				switch ik {
				case "time":
					trades_data[k].Time = iv.(string)
				case "price":
					trades_data[k].Price = util.InterfaceToFloat64(iv)
				case "amount":
					trades_data[k].Amount = util.InterfaceToFloat64(iv)
				case "type":
					trades_data[k].Type = iv.(string)
				}
			}
		default:
			logger.Errorln(k, v)
			logger.Fatalln("don't know the type, crash!")
			return false
		}
	}
	return true
}

func parse_buy_sell(sells_buys []interface{}, sells_buys_data *[10]SellBuy) bool {
	for k, v := range sells_buys {
		switch vt := v.(type) {
		case map[string]interface{}:
			logger.Debugln(k, " is a map:")
			logger.Debugf("sells/buys[%d]\n", k)

			for ik, iv := range vt {
				switch ik {
				case "price":
					sells_buys_data[k].Price = util.InterfaceToFloat64(iv)
				case "level":
					sells_buys_data[k].Level = util.InterfaceToFloat64(iv)
				case "amount":
					sells_buys_data[k].Amount = util.InterfaceToFloat64(iv)
				}
			}
		default:
			logger.Errorln(k, v)
			logger.Fatalln("don't know the type, crash!")
			return false
		}
	}

	return true
}

func parse_topbuy(topbuys []interface{}, topbuys_data *[5]Top_buy_sell) bool {
	for k, v := range topbuys {
		switch vt := v.(type) {
		case map[string]interface{}:
			logger.Debugln(k, " is a map:")
			logger.Debugf("topbuys[%d]\n", k)

			for ik, iv := range vt {
				switch ik {
				case "price":
					topbuys_data[k].Price = util.InterfaceToFloat64(iv)
				case "amount":
					topbuys_data[k].Amount = util.InterfaceToFloat64(iv)
				case "level":
					topbuys_data[k].Level = util.InterfaceToFloat64(iv)
				case "accu":
					topbuys_data[k].Accu = util.InterfaceToFloat64(iv)
				}
			}
		default:
			logger.Errorln(k, v)
			logger.Fatalln("don't know the type, crash!")
			return false
		}
	}
	return true
}

// just fuck the huobi old shit bug
func parse_topsell(topsells map[string]interface{}, topsells_data *[5]Top_buy_sell) bool {
	index := 4
	for k, v := range topsells {
		switch vt := v.(type) {
		case map[string]interface{}:
			logger.Debugln(k, " is a map:")
			logger.Debugf("topsells[%s]\n", k)

			for ik, iv := range vt {
				logger.Debugln(ik, iv)
				switch ik {
				case "price":
					topsells_data[index].Price = util.InterfaceToFloat64(iv)
				case "amount":
					topsells_data[index].Amount = util.InterfaceToFloat64(iv)
				case "level":
					topsells_data[index].Level = util.InterfaceToFloat64(iv)
				case "accu":
					topsells_data[index].Accu = util.InterfaceToFloat64(iv)
				}
			}

			index--

		default:
			logger.Errorln(k, v)
			logger.Fatalln("don't know the type, crash!")
			return false
		}
	}
	return true
}
