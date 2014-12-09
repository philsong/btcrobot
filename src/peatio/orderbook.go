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

package peatio

import (
	. "common"
	. "config"
	"encoding/json"
	"io/ioutil"
	"logger"
	"net/http"
	"util"
)

func (w *Peatio) getOrderBook(symbol string) (ret bool, orderBook OrderBook) {
	// init to false
	ret = false
	req, err := http.NewRequest("GET", Config["peatio_depth_url"], nil)
	if err != nil {
		logger.Fatal(err)
		return
	}

	req.Header.Set("Referer", Config["peatio_base_url"])
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

			ioutil.WriteFile("cache/peatiodepth.json", bodyByte, 0644)
		}
	}

	return w.analyzeOrderBook(body)
}

type OKMarketOrder struct {
	Price  float64 // 价格
	Amount float64 // 委单量
}

type _PeatioOrderBook struct {
	Asks [10]interface{}
	Bids [10]interface{}
}

type PeatioOrderBook struct {
	Asks [10]OKMarketOrder
	Bids [10]OKMarketOrder
}

func convert2struct(_peatioOrderBook _PeatioOrderBook) (peatioOrderBook PeatioOrderBook) {
	for k, v := range _peatioOrderBook.Asks {
		switch vt := v.(type) {
		case []interface{}:
			for ik, iv := range vt {
				switch ik {
				case 0:
					peatioOrderBook.Asks[k].Price = util.InterfaceToFloat64(iv)
				case 1:
					peatioOrderBook.Asks[k].Amount = util.InterfaceToFloat64(iv)
				}
			}
		}
	}

	for k, v := range _peatioOrderBook.Bids {
		switch vt := v.(type) {
		case []interface{}:
			for ik, iv := range vt {
				switch ik {
				case 0:
					peatioOrderBook.Bids[k].Price = util.InterfaceToFloat64(iv)
				case 1:
					peatioOrderBook.Bids[k].Amount = util.InterfaceToFloat64(iv)
				}
			}
		}
	}
	return
}

func (w *Peatio) analyzeOrderBook(content string) (ret bool, orderBook OrderBook) {
	// init to false
	ret = false
	var _peatioOrderBook _PeatioOrderBook
	if err := json.Unmarshal([]byte(content), &_peatioOrderBook); err != nil {
		logger.Infoln(err)
		return
	}

	peatioOrderBook := convert2struct(_peatioOrderBook)

	for i := 0; i < 10; i++ {
		orderBook.Asks[i].Price = peatioOrderBook.Asks[len(_peatioOrderBook.Asks)-10+i].Price
		orderBook.Asks[i].Amount = peatioOrderBook.Asks[len(_peatioOrderBook.Asks)-10+i].Amount
		orderBook.Bids[i].Price = peatioOrderBook.Bids[i].Price
		orderBook.Bids[i].Amount = peatioOrderBook.Bids[i].Amount
	}

	//logger.Infoln(orderBook)
	ret = true
	return
}
