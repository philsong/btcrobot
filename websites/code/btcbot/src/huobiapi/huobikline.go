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
	"logger"
)

/*
import (
	"io"
	"logger"
	"strings"
)
*/

type PeroidRecord struct {
	Date   string
	Time   string
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volumn float64
	Amount float64
}

type MinuteRecord struct {
	Time   string
	Price  float64
	Volumn float64
	Amount float64
}

func (w *Huobi) AnalyzePeroidLine(filename string, content string) bool {
	//logger.Infoln(content)
	//logger.Infoln(filename)
	PeroidRecords := ParsePeroidCSV(filename)

	var Time []string
	var Price []float64
	var Volumn []float64
	for _, v := range PeroidRecords {
		Time = append(Time, v.Date+" "+v.Time)
		Price = append(Price, v.Close)
		Volumn = append(Volumn, v.Volumn)
		//Price = append(Price, (v.Close+v.Open+v.High+v.Low)/4.0)
		//Price = append(Price, v.Low)
	}
	w.Time = Time
	w.Price = Price
	w.Volumn = Volumn

	//rsi(Price)
	if Config["env"] == "test" {
		w.do2Percent(Time, Price)
		return true

		k, d, j := doKDJ(PeroidRecords)
		length := len(k)
		// Loop through the entire array.
		for i := 0; i < length; i++ {
			logger.Infof("[%s-%s]%d/%d/%d\n", PeroidRecords[i].Date, PeroidRecords[i].Time, int(k[i]), int(d[i]), int(j[i]))
		}
	} else {
		w.doEMA(Time, Price, Volumn)
		return true
	}

	return true
}

func (w *Huobi) AnalyzeMinuteLine(filename string, content string) bool {
	//logger.Infoln(content)
	logger.Debugln(filename)
	MinuteRecords := ParseMinuteCSV(filename)
	var Time []string
	var Price []float64
	var Volumn []float64
	for _, v := range MinuteRecords {
		Time = append(Time, v.Time)
		Price = append(Price, v.Price)
		Volumn = append(Volumn, v.Volumn)
	}

	w.Time = Time
	w.Price = Price
	w.Volumn = Volumn

	if Config["env"] == "test" {
		w.do2Percent(Time, Price)
		return true
	} else {
		w.doEMA(Time, Price, Volumn)
		return true
	}
}

func (w *Huobi) getTradePrice(tradeDirection string) string {
	if len(w.Price) == 0 {
		logger.Errorln("get price failed, array len=0")
		return "false"
	}
	var finalTradePrice float64
	if tradeDirection == "buy" {
		finalTradePrice = w.Price[len(w.Price)-1] + w.Slippage
	} else if tradeDirection == "sell" {
		finalTradePrice = w.Price[len(w.Price)-1] - w.Slippage
	} else {
		finalTradePrice = w.Price[len(w.Price)-1]
	}
	return fmt.Sprintf("%0.02f", finalTradePrice)
}
