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

package simulate

import (
	. "common"
	. "config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"logger"
	"net/http"
	"util"
)

func (w *Simulate) getOrderBook(symbol string) (ret bool, orderBook OrderBook) {
	// init to false
	ret = false
	req, err := http.NewRequest("GET", fmt.Sprintf(Config["ok_depth_url"], symbol), nil)
	if err != nil {
		logger.Fatal(err)
		return
	}

	req.Header.Set("Referer", Config["ok_base_url"])
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36")

	logger.Traceln(req)

	c := util.NewTimeoutClient()
	logger.Tracef("HTTP req begin getOrderBook")
	resp, err := c.Do(req)
	logger.Tracef("HTTP req end getOrderBook")
	if err != nil {
		logger.Traceln(err)
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

			ioutil.WriteFile("cache/okdepth.json", bodyByte, 0644)
		}
	}

	return w.analyzeOrderBook(body)
}

type OKMarketOrder struct {
	Price  float64 // 价格
	Amount float64 // 委单量
}

type _OKOrderBook struct {
	Asks [200]interface{}
	Bids [200]interface{}
}

type OKOrderBook struct {
	Asks [200]OKMarketOrder
	Bids [200]OKMarketOrder
}

func convert2struct(_okOrderBook _OKOrderBook) (okOrderBook OKOrderBook) {
	for k, v := range _okOrderBook.Asks {
		switch vt := v.(type) {
		case []interface{}:
			for ik, iv := range vt {
				switch ik {
				case 0:
					okOrderBook.Asks[k].Price = util.InterfaceToFloat64(iv)
				case 1:
					okOrderBook.Asks[k].Amount = util.InterfaceToFloat64(iv)
				}
			}
		}
	}

	for k, v := range _okOrderBook.Bids {
		switch vt := v.(type) {
		case []interface{}:
			for ik, iv := range vt {
				switch ik {
				case 0:
					okOrderBook.Bids[k].Price = util.InterfaceToFloat64(iv)
				case 1:
					okOrderBook.Bids[k].Amount = util.InterfaceToFloat64(iv)
				}
			}
		}
	}
	return
}

func (w *Simulate) analyzeOrderBook(content string) (ret bool, orderBook OrderBook) {
	// init to false
	ret = false
	var _okOrderBook _OKOrderBook
	if err := json.Unmarshal([]byte(content), &_okOrderBook); err != nil {
		logger.Infoln(err)
		return
	}

	okOrderBook := convert2struct(_okOrderBook)

	for i := 0; i < 10; i++ {
		orderBook.Asks[i].Price = okOrderBook.Asks[len(_okOrderBook.Asks)-10+i].Price
		orderBook.Asks[i].Amount = okOrderBook.Asks[len(_okOrderBook.Asks)-10+i].Amount
		orderBook.Bids[i].Price = okOrderBook.Bids[i].Price
		orderBook.Bids[i].Amount = okOrderBook.Bids[i].Amount
	}

	ret = true
	return
}
