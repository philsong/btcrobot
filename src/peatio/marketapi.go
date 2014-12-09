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
	"fmt"
	"io/ioutil"
	"logger"
	"net/http"
	"strconv"
	"strings"
	"time"
	"util"
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

func (w *Peatio) AnalyzeKLinePeroid(symbol string, peroid int) (ret bool, records []Record) {
	ret = false
	if symbol != "btc_cny" {
		logger.Fatal("I only add btccny for peatio by now.")
		return
	}

	req, err := http.NewRequest("GET", fmt.Sprintf(Config["peatio_kline_url"], peroid), nil)
	if err != nil {
		logger.Fatal(err)
		return
	}

	req.Header.Set("Referer", Config["peatio_base_url"])
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36")

	logger.Traceln(req)

	c := util.NewTimeoutClient()

	logger.Tracef("HTTP req begin AnalyzeKLinePeroid")
	resp, err := c.Do(req)
	logger.Tracef("HTTP req end AnalyzeKLinePeroid")
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
		}
	}

	ioutil.WriteFile(fmt.Sprintf("cache/peatioKLine_%03d.data", peroid), []byte(body), 0644)

	return analyzePeroidLine(body)
}

func analyzePeroidLine(content string) (ret bool, records []Record) {
	logger.Traceln("peatioKLine analyzePeroidLine begin....")
	content = strings.TrimPrefix(content, "[[")
	content = strings.TrimSuffix(content, "]]")

	ret = false
	for _, value := range strings.Split(content, `],[`) {
		//logger.Traceln(value)
		v := strings.Split(value, ",")
		if len(v) < 6 {
			logger.Debugln("wrong data")
			return
		}

		var record Record
		Time, err := strconv.ParseInt(v[0], 0, 64)
		if err != nil {
			logger.Debugln("config item is not float")
			return
		}

		Open, err := strconv.ParseFloat(v[1], 64)
		if err != nil {
			logger.Debugln("config item is not float")
			return
		}

		Close, err := strconv.ParseFloat(v[2], 64)
		if err != nil {
			logger.Debugln("config item is not float")
			return
		}

		High, err := strconv.ParseFloat(v[3], 64)
		if err != nil {
			logger.Debugln("config item is not float")
			return
		}

		Low, err := strconv.ParseFloat(v[4], 64)
		if err != nil {
			logger.Debugln("config item is not float")
			return
		}

		Volumn, err := strconv.ParseFloat(v[5], 64)
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

	logger.Traceln("peatioKLine parsePeroidArray end....")
	ret = true
	return
}
