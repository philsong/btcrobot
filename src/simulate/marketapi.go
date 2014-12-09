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

package simulate

import (
	. "common"
	. "config"
	"fmt"
	"io/ioutil"
	"logger"
	"net/http"
	"strconv"
	"strings"
	"time"
	"util"
)

func (w *Simulate) AnalyzeKLinePeroid(symbol string, peroid int) (ret bool, records []Record) {
	var oksymbol string
	if symbol == "btc_cny" {
		oksymbol = "okcoinbtccny"
	} else {
		oksymbol = "okcoinltccny"
	}

	ret = false
	now := time.Now().UnixNano() / 1000000

	req, err := http.NewRequest("GET", fmt.Sprintf(Config["ok_kline_url"], 60*peroid, oksymbol, now), nil)
	if err != nil {
		logger.Fatal(err)
		return
	}

	req.Header.Set("Referer", Config["ok_base_url"])
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36")

	logger.Traceln(req)

	c := util.NewTimeoutClient()
	logger.Tracef("okHTTP req begin")
	resp, err := c.Do(req)
	logger.Tracef("okHTTP req end")
	if err != nil {
		logger.Traceln(err)
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		logger.Tracef("HTTP returned status %v", resp)
		return
	}

	var body string
	contentEncoding := resp.Header.Get("Content-Encoding")

	logger.Tracef("HTTP returned Content-Encoding %s", contentEncoding)
	logger.Traceln(resp.Header.Get("Content-Type"))

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

			ioutil.WriteFile(fmt.Sprintf("cache/okTradeKLine_%03d.data", peroid), bodyByte, 0644)
		}
	}

	return analyzePeroidLine(body)
}

func analyzePeroidLine(content string) (ret bool, records []Record) {
	logger.Traceln("Okcoin analyzePeroidLine begin....")
	content = strings.TrimPrefix(content, "[[")
	content = strings.TrimSuffix(content, "]]")

	ret = false
	for _, value := range strings.Split(content, `],[`) {
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

		records = append(records, record)
	}

	logger.Traceln("Okcoin parsePeroidArray end....")
	ret = true
	return
}
