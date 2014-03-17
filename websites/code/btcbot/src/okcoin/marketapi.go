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
	. "config"
	"fmt"
	"io/ioutil"
	"logger"
	"net/http"
	"os"
	"strategy"
	"strings"
	"time"
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

	req, err := http.NewRequest("GET", fmt.Sprintf("https://www.okcoin.com/kline/period.do?step=%d&symbol=%s&nonce=%d", 60*peroid, oksymbol, now), nil)
	if err != nil {
		logger.Fatal(err)
		return false
	}

	req.Header.Set("Referer", "https://www.okcoin.com/")
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

				ioutil.WriteFile(fmt.Sprintf("cache/okTradeKLine_%03d.data", peroid), bodyByte, os.ModeAppend)
			}
		}

		logger.Traceln(resp.Header.Get("Content-Type"))

		ret := strings.Contains(body, "您需要登录才能继续")
		if ret {
			logger.Traceln("您需要登录才能继续")
			return false
		} else {
			return w.analyzePeroidLine(fmt.Sprintf("cache/okTradeKLine_%03d.data", peroid), body)
		}

	} else {
		logger.Tracef("HTTP returned status %v", resp)
	}

	return false
}

type PeroidRecord struct {
	Time   int
	zero1  int
	zero2  int
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volumn float64
}

func (w *Okcoin) analyzePeroidLine(filename string, content string) bool {
	//logger.Infoln(content)
	//logger.Infoln(filename)
	//PeroidRecords := parsePeroidCSV(filename)

	var PeroidRecords []PeroidRecord
	logger.Traceln("Okcoin analyzePeroidLine begin....")
	content = strings.TrimPrefix(content, "[[")
	content = strings.TrimSuffix(content, "]]")

	for _, value := range strings.Split(content, `],[`) {

		//logger.Traceln(value)
		v := strings.Split(value, ",")
		if len(valueItem) != 8 {
			logger.Traceln("wrong data")
			break
		}

		var peroidRecord PeroidRecord
		_, err := strconv.ParseInt(v[0], 64)
		if err != nil {
			logger.Debugln("config item tradeAmount is not float")
			return false
		}

		peroidRecord.Time = v[0]
		peroidRecord.zero1 = v[1]
		peroidRecord.zero2 = v[2]
		_, err := strconv.ParseFloat(v[0], 64)
		if err != nil {
			logger.Debugln("config item tradeAmount is not float")
			return false
		}

		peroidRecord.Open = v[3]
		peroidRecord.High = v[4]
		peroidRecord.Low = v[5]
		peroidRecord.Close = v[6]
		peroidRecord.Volumn = v[7].(float64)

		logger.Traceln(peroidRecord)

		//ids = append(ids, valueItem)
	}
	logger.Traceln("Okcoin analyzePeroidLine end....")
	return false

	fmt.Println(PeroidRecords)
	return true
	var Time []string
	var Price []float64
	var Volumn []float64
	for _, v := range PeroidRecords {
		//Time = append(Time, v.Time)
		Price = append(Price, v.Close)
		Volumn = append(Volumn, v.Volumn)
		//Price = append(Price, (v.Close+v.Open+v.High+v.Low)/4.0)
		//Price = append(Price, v.Low)
	}
	w.Time = Time
	w.Price = Price
	w.Volumn = Volumn
	strategyName := Option["strategy"]
	strategy.Perform(strategyName, *w, Time, Price, Volumn)

	return true
}
