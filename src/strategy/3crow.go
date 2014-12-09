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

type the3crowStrategy struct {
	PrevClosePrice float64
}

func init() {
	the3crow := new(the3crowStrategy)
	Register("the3crow", the3crow)
}

// the3crow strategy
func (the3crow *the3crowStrategy) Tick(records []Record) bool {
	if the3crow.PrevClosePrice == lastPrice {
		return false
	}

	the3crow.PrevClosePrice = lastPrice

	logger.Infof("lastPrice %0.02f\n", lastPrice)
	logger.Infof("3 open %0.02f close %0.02f\n", records[length-2].Open, records[length-2].Close)
	logger.Infof("2 open %0.02f close %0.02f\n", records[length-3].Open, records[length-3].Close)
	logger.Infof("1 open %0.02f close %0.02f\n", records[length-4].Open, records[length-4].Close)

	if records[length-2].Close > records[length-2].Open {
		logger.Infof("3阳")
	} else {
		logger.Infof("3阴")
	}

	if records[length-3].Close > records[length-3].Open {
		logger.Infof("2阳")
	} else {
		logger.Infof("2阴")
	}
	if records[length-4].Close > records[length-4].Open {
		logger.Infof("1阳")
	} else {
		logger.Infof("1阴")
	}
	logger.Infoln("---------")

	// the3crow cross
	if records[length-2].Close > records[length-2].Open &&
		records[length-3].Close > records[length-3].Open &&
		records[length-4].Close > records[length-4].Open {
		Buy()
	} else if records[length-2].Close < records[length-2].Open &&
		records[length-3].Close < records[length-3].Open &&
		records[length-4].Close < records[length-4].Open {
		Sell()
	}

	// do sell when price is below stoploss point
	processStoploss(lastPrice)

	processTimeout()

	return true
}
