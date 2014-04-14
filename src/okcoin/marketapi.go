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

package okcoin

import (
	. "common"
	. "config"
	"fmt"
	"io/ioutil"
	"logger"
	"net/http"
	"strategy"
	"strconv"
	"strings"
	"time"
	"util"
)

/*
	SEE DOC:
	TRADE API
	https://www.okcoin.com/t-1000097.html

	行情API
	https://www.okcoin.com/shequ/themeview.do?tid=1000052&currentPage=1

	//non-official API :P
	K线数据step单位为second
	https://www.okcoin.com/kline/period.do?step=60&symbol=okcoinbtccny&nonce=1394955131098

	https://www.okcoin.com/kline/trades.do?since=10625682&symbol=okcoinbtccny&nonce=1394955760557

	https://www.okcoin.com/kline/depth.do?symbol=okcoinbtccny&nonce=1394955767484

	https://www.okcoin.com/real/ticker.do?symbol=0&random=61

	//old kline for btc
	日数据
	https://www.okcoin.com/klineData.do?type=3&marketFrom=0
	5分钟数据
	https://www.okcoin.com/klineData.do?type=1&marketFrom=0

	/for ltc
	https://www.okcoin.com/klineData.do?type=3&marketFrom=3
*/
func (w *Okcoin) AnalyzeKLinePeroid(symbol string, peroid int) (ret bool) {
	var oksymbol string
	if symbol == "btc_cny" {
		oksymbol = "okcoinbtccny"
	} else {
		oksymbol = "okcoinltccny"
	}

	now := time.Now().UnixNano() / 1000000

	req, err := http.NewRequest("GET", fmt.Sprintf(Config["ok_kline_url"], 60*peroid, oksymbol, now), nil)
	if err != nil {
		logger.Fatal(err)
		return false
	}

	req.Header.Set("Referer", Config["ok_base_url"])
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36")

	logger.Traceln(req)

	c := util.NewTimeoutClient()
	logger.Tracef("HTTP req begin AnalyzeKLinePeroid")
	resp, err := c.Do(req)
	logger.Tracef("HTTP req end AnalyzeKLinePeroid")
	if err != nil {
		logger.Traceln(err)
		return false
	}

	defer resp.Body.Close()
	if resp.StatusCode == 200 {
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
				return false
			} else {
				body = string(bodyByte)

				ioutil.WriteFile(fmt.Sprintf("cache/okTradeKLine_%03d.data", peroid), bodyByte, 0644)
			}
		}

		logger.Traceln(resp.Header.Get("Content-Type"))

		return w.analyzePeroidLine(fmt.Sprintf("cache/okTradeKLine_%03d.data", peroid), body)
	} else {
		logger.Tracef("HTTP returned status %v", resp)
	}

	return false
}

func parsePeroidArray(content string) (ret bool, records []Record) {
	logger.Traceln("Okcoin parsePeroidArray begin....")
	content = strings.TrimPrefix(content, "[[")
	content = strings.TrimSuffix(content, "]]")

	ret = false
	for _, value := range strings.Split(content, `],[`) {
		//logger.Traceln(value)
		v := strings.Split(value, ",")
		if len(v) < 8 {
			logger.Debugln("wrong data")
			return
		}

		var record Record
		Time, err := strconv.ParseInt(v[0], 0, 64)
		if err != nil {
			logger.Debugln("config item is not float")
			return
		}

		_, err = strconv.ParseInt(v[1], 0, 64)
		if err != nil {
			logger.Debugln("config item is not float")
			return
		}
		_, err = strconv.ParseInt(v[2], 0, 64)
		if err != nil {
			logger.Debugln("config item is not float")
			return
		}

		Open, err := strconv.ParseFloat(v[3], 64)
		if err != nil {
			logger.Debugln("config item is not float")
			return
		}

		Close, err := strconv.ParseFloat(v[4], 64)
		if err != nil {
			logger.Debugln("config item is not float")
			return
		}

		High, err := strconv.ParseFloat(v[5], 64)
		if err != nil {
			logger.Debugln("config item is not float")
			return
		}

		Low, err := strconv.ParseFloat(v[6], 64)
		if err != nil {
			logger.Debugln("config item is not float")
			return
		}

		Volumn, err := strconv.ParseFloat(v[7], 64)
		if err != nil {
			logger.Debugln("config item is not float")
			return
		}

		const layout = "2006-01-02 15:04:05"
		t := time.Unix(Time, 0)
		record.TimeStr = t.Format(layout)
		record.Time = Time
		record.Open = Open
		record.High = High
		record.Low = Low
		record.Close = Close
		record.Volumn = Volumn

		//logger.Traceln(records)

		records = append(records, record)
	}

	logger.Traceln("Okcoin parsePeroidArray end....")
	ret = true
	return
}

func (w *Okcoin) analyzePeroidLine(filename string, content string) bool {
	//logger.Infoln(content)
	//logger.Infoln(filename)
	ret, records := parsePeroidArray(content)
	if ret == false {
		logger.Errorln("Okcoin parsePeroidArray failed....")
		return false
	}

	strategyName := Option["strategy"]
	strategy.Perform(strategyName, *w, records)

	return true
}
