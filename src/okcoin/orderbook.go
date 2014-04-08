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

package okcoin

import (
	. "common"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"logger"
	"net/http"
	//"strconv"
	//"strings"
	"reflect"
	"util"
)

func (w *Okcoin) getOrderBook(symbol string) (ret bool, orderBook OrderBook) {
	//init to false
	ret = false
	req, err := http.NewRequest("GET", fmt.Sprintf("https://www.okcoin.com/api/depth.do?symbol=%s", symbol), nil)
	if err != nil {
		logger.Fatal(err)
		return
	}

	req.Header.Set("Referer", "https://www.okcoin.com/")
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
		body = DumpGZIP(resp.Body)

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
	Price  float64 //价格
	Amount float64 //委单量
}

type OKOrderBook struct {
	Asks [60]interface{}
	Bids [60]interface{}
}

func (w *Okcoin) analyzeOrderBook(content string) (ret bool, orderBook OrderBook) {
	//init to false
	ret = false
	var okOrderBook OKOrderBook
	if err := json.Unmarshal([]byte(content), &okOrderBook); err != nil {
		logger.Infoln(err)
		return
	}

	//logger.Infoln(orderBook.Asks)
	logger.Infoln((okOrderBook.Asks[len(okOrderBook.Asks)-1]))
	logger.Infoln(okOrderBook.Bids[0])
	/*
		i := 0
		for _, value := range okOrderBook.Bids {

			value = strings.TrimPrefix(value.(type)(string), "[")
			value = strings.TrimSuffix(value, "]")
			//logger.Traceln(value)
			v := strings.Split(value, ",")
			if len(v) < 2 {
				logger.Debugln("wrong data")
				return
			}

			orderBook.Bids[i].Price, err = strconv.ParseFloat(v[0], 64)
			if err != nil {
				logger.Debugln("config item is not float")
				return
			}
			orderBook.Bids[i].Amount, err = strconv.ParseFloat(v[1], 64)
			if err != nil {
				logger.Debugln("config item is not float")
				return
			}
			i++
			if i == 9 {
				break
			}
		}

		i = 0
		for _, value := range okOrderBook.Asks {
			if i < 51 {
				continue
			}
			value = strings.TrimPrefix(value, "[")
			value = strings.TrimSuffix(value, "]")
			//logger.Traceln(value)
			v := strings.Split(value, ",")
			if len(v) < 2 {
				logger.Debugln("wrong data")
				return
			}

			orderBook.Asks[i].Price, err = strconv.ParseFloat(v[0], 64)
			if err != nil {
				logger.Debugln("config item is not float")
				return
			}
			orderBook.Asks[i].Amount, err = strconv.ParseFloat(v[1], 64)
			if err != nil {
				logger.Debugln("config item is not float")
				return
			}
			i++
		}
	*/
	var sells_buys_data [60]OKMarketOrder
	parse_array(okOrderBook.Bids, &sells_buys_data)
	/*
		for i := 0; i < 10; i++ {
			orderBook.Asks[i].Price = okOrderBook.Asks[i].Price
			orderBook.Asks[i].Amount = okOrderBook.Asks[i].Amount
			orderBook.Bids[i].Price = okOrderBook.Bids[i].Price
			orderBook.Bids[i].Amount = okOrderBook.Bids[i].Amount
		}
	*/
	//OrderBook
	//logger.Infoln(orderBook)
	return
}

func parse_array(sells_buys [60]interface{}, sells_buys_data *[60]OKMarketOrder) bool {
	for k, v := range sells_buys {
		switch vt := v.(type) {
		case map[string]interface{}:
			logger.Debugln(k, " is a map:")
			logger.Debugf("sells/buys[%d]\n", k)

			for ik, iv := range vt {
				switch ik {
				case "price":
					sells_buys_data[k].Price = util.InterfaceToFloat64(iv)
				case "amount":
					sells_buys_data[k].Amount = util.InterfaceToFloat64(iv)
				}
			}
		default:
			logger.Fatalln(k, v)

			logger.Fatalln(vt)

			switch reflect.TypeOf(vt).Kind() {
			case reflect.Slice:
				s := reflect.ValueOf(vt)
				fmt.Println(s)

				for i := 0; i < s.Len(); i++ {
					fmt.Println(s.Index(i))
					fmt.Println(util.InterfaceToFloat64(s.Index(i)))
				}
			}
			return true
		}
	}

	return false
}
