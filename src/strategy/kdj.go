/*
  btcrobot is a Bitcoin, Litecoin and Altcoin trading bot written in golang,
  it features multiple trading methods using technical analysis.

  Disclaimer:

  USE AT YOUR OWN RISK!

  The author of this project is NOT responsible for any damage or loss caused
  by this software. There can be bugs and the bot may not Tick as expected
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

type KDJStrategy struct {
	PrevTime  string
	PrevPrice float64
}

func init() {
	kdjStrategy := new(KDJStrategy)
	PrevTrade = "init"
	Register("KDJ", kdjStrategy)
}

//xxx strategy
func (kdjStrategy *KDJStrategy) Tick(records []Record) bool {
	//实现自己的策略

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

	if kdjStrategy.PrevTime == records[length-1].TimeStr &&
		kdjStrategy.PrevPrice == records[length-1].Close {
		return false
	}

	//K线为白，D线为黄，J线为红,K in middle
	k, d, j := getKDJ(records)

	if kdjStrategy.PrevTime != records[length-1].TimeStr ||
		kdjStrategy.PrevPrice != records[length-1].Close {
		kdjStrategy.PrevTime = records[length-1].TimeStr
		kdjStrategy.PrevPrice = records[length-1].Close

		logger.Infoln(records[length-1].TimeStr, records[length-1].Close)
		logger.Infof("d(黄线）%0.0f\tk(白线）%0.0f\tj(红线）%0.0f\n", d[length-2], k[length-2], j[length-2])
		logger.Infof("d(黄线）%0.0f\tk(白线）%0.0f\tj(红线）%0.0f\n", d[length-1], k[length-1], j[length-1])
	}

	if ((j[length-2] < k[length-2] && k[length-2] < d[length-2]) || PrevTrade == "sell") &&
		(j[length-1] > k[length-1] && k[length-1] > d[length-1]) {
		logger.Infoln("KDJ up cross")
		if (PrevTrade == "init" && d[length-2] <= 30) || PrevTrade == "sell" {
			//do buy
			warning := "KDJ up cross, 买入buy In<----市价" + getTradePrice("", Price[length-1]) +
				",委托价" + getTradePrice("buy", Price[length-1])
			logger.Infoln(warning)
			if Buy(getTradePrice("buy", Price[length-1]), tradeAmount) != "0" {
				warning += "[委托成功]"
				PrevBuyPirce = Price[length-1]
				PrevTrade = "buy"
			} else {
				warning += "[委托失败]"
			}

			go email.TriggerTrender(warning)
		}

	}

	if ((j[length-2] > k[length-2] && k[length-2] > d[length-2]) || PrevTrade == "buy") &&
		(j[length-1] < k[length-1] && k[length-1] < d[length-1]) {

		logger.Infoln("KDJ down cross")
		if (PrevTrade == "init" && d[length-2] >= 70) || PrevTrade == "buy" {
			//do sell
			warning := "KDJ down cross, 卖出Sell Out---->市价" + getTradePrice("", Price[length-1]) +
				",委托价" + getTradePrice("sell", Price[length-1])
			logger.Infoln(warning)
			if Sell(getTradePrice("sell", Price[length-1]), tradeAmount) != "0" {
				warning += "[委托成功]"
				PrevTrade = "sell"
				PrevBuyPirce = 0
			} else {
				warning += "[委托失败]"
			}

			go email.TriggerTrender(warning)
		}

	}

	//do sell when price is below stoploss point
	processStoploss(Price)

	return true
}
