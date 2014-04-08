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

package strategy

import (
	. "common"
	. "config"
	"email"
	"logger"
)

type KDJexStrategy struct {
	PrevKDJTrade string
	PrevTime     string
	PrevPrice    float64
}

func init() {
	kdjexStrategy := new(KDJexStrategy)
	kdjexStrategy.PrevKDJTrade = "init"
	Register("KDJ-EX", kdjexStrategy)
}

//KDJ-EX strategy
func (kdjexStrategy *KDJexStrategy) Perform(tradeAPI TradeAPI, records []Record) bool {
	tradeAmount := Option["tradeAmount"]

	var Time []string
	var Price []float64
	var Volumn []float64
	for _, v := range records {
		Time = append(Time, v.TimeStr)
		Price = append(Price, v.Close)
		Volumn = append(Volumn, v.Volumn)
		//Price = append(Price, (v.Close+v.Open+v.High+v.Low)/4.0)
		//Price = append(Price, v.Low)
	}

	length := len(records)

	if kdjexStrategy.PrevTime == records[length-1].TimeStr &&
		kdjexStrategy.PrevPrice == records[length-1].Close {
		return false
	}

	//K线为白，D线为黄，J线为红,K in middle
	k, d, j := getKDJ(records)

	if kdjexStrategy.PrevTime != records[length-1].TimeStr ||
		kdjexStrategy.PrevPrice != records[length-1].Close {
		kdjexStrategy.PrevTime = records[length-1].TimeStr
		kdjexStrategy.PrevPrice = records[length-1].Close

		logger.Infoln(records[length-1].TimeStr, records[length-1].Close)
		logger.Infof("d(黄线）%0.0f\tk(白线）%0.0f\tj(红线）%0.0f\n", d[length-2], k[length-2], j[length-2])
		logger.Infof("d(黄线）%0.0f\tk(白线）%0.0f\tj(红线）%0.0f\n", d[length-1], k[length-1], j[length-1])
	}

	if ((j[length-2] < k[length-2] && k[length-2] < d[length-2]) || kdjexStrategy.PrevKDJTrade == "sell") &&
		(j[length-1] > k[length-1] && k[length-1] > d[length-1]) {
		logger.Infoln("KDJ up cross")
		if (kdjexStrategy.PrevKDJTrade == "init" && d[length-2] <= 30) || kdjexStrategy.PrevKDJTrade == "sell" {
			//do buy
			warning := "KDJ up cross, 买入buy In<----市价" + tradeAPI.GetTradePrice("", Price[length-1]) +
				",委托价" + tradeAPI.GetTradePrice("buy", Price[length-1])
			logger.Infoln(warning)
			if tradeAPI.Buy(tradeAPI.GetTradePrice("buy", Price[length-1]), tradeAmount) {
				warning += "[委托成功]"
			} else {
				warning += "[委托失败]"
			}

			kdjexStrategy.PrevKDJTrade = "buy"

			go email.TriggerTrender(warning)
		}

	}

	if ((j[length-2] > k[length-2] && k[length-2] > d[length-2]) || kdjexStrategy.PrevKDJTrade == "buy") &&
		(j[length-1] < k[length-1] && k[length-1] < d[length-1]) {

		logger.Infoln("KDJ down cross")
		if (kdjexStrategy.PrevKDJTrade == "init" && d[length-2] >= 70) || kdjexStrategy.PrevKDJTrade == "buy" {
			//do sell
			warning := "KDJ down cross, 卖出Sell Out---->市价" + tradeAPI.GetTradePrice("", Price[length-1]) +
				",委托价" + tradeAPI.GetTradePrice("sell", Price[length-1])
			logger.Infoln(warning)
			if tradeAPI.Sell(tradeAPI.GetTradePrice("sell", Price[length-1]), tradeAmount) {
				warning += "[委托成功]"
			} else {
				warning += "[委托失败]"
			}

			kdjexStrategy.PrevKDJTrade = "sell"

			go email.TriggerTrender(warning)
		}

	}

	return true
}
