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
	"fmt"
	"io/ioutil"
	"logger"
	"math/rand"
	"net/http"
	"os"
	"strings"
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
			return w.analyzePeroidLine(fmt.Sprintf("cache/TradeKLine_%03d.data", peroid), body)
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
			return w.analyzeMinuteLine(fmt.Sprintf("cache/TradeKLine_minute.data"), body)
		}

	} else {
		logger.Tracef("HTTP returned status %v", resp)
	}

	return false
}
