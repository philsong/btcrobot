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

//kdjStrategy strategy
func (kdjStrategy *KDJStrategy) Tick(records []Record) bool {
	if kdjStrategy.PrevTime == records[length-1].TimeStr &&
		kdjStrategy.PrevPrice == lastPrice {
		return false
	}

	//K线为白，D线为黄，J线为红,K in middle
	k, d, j := getKDJ(records)

	if kdjStrategy.PrevTime != records[length-1].TimeStr ||
		kdjStrategy.PrevPrice != records[length-1].Close {
		kdjStrategy.PrevTime = records[length-1].TimeStr
		kdjStrategy.PrevPrice = lastPrice

		logger.Infoln(records[length-1].TimeStr, records[length-1].Close)
		logger.Infof("d(黄线）%0.0f\tk(白线）%0.0f\tj(红线）%0.0f\n", d[length-2], k[length-2], j[length-2])
		logger.Infof("d(黄线）%0.0f\tk(白线）%0.0f\tj(红线）%0.0f\n", d[length-1], k[length-1], j[length-1])
	}

	if ((j[length-2] < k[length-2] && k[length-2] < d[length-2]) || PrevTrade == "sell") &&
		(j[length-1] > k[length-1] && k[length-1] > d[length-1]) {
		logger.Infoln("KDJ up cross")
		if d[length-2] <= 30 {
			Buy()
		}
	}

	if ((j[length-2] > k[length-2] && k[length-2] > d[length-2]) || PrevTrade == "buy") &&
		(j[length-1] < k[length-1] && k[length-1] < d[length-1]) {
		logger.Infoln("KDJ down cross")
		if d[length-2] >= 70 {
			Sell()
		}
	}

	//do sell when price is below stoploss point
	processStoploss(lastPrice)
	processTimeout()

	return true
}
