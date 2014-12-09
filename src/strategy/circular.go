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
	"logger"
	"strconv"
)

type circularStrategy struct {
	PrevClosePrice float64
}

func init() {
	circular := new(circularStrategy)
	Register("circular", circular)
}

// circular strategy
func (circular *circularStrategy) Tick(records []Record) bool {
	if circular.PrevClosePrice == lastPrice {
		return false
	}

	basePrice, err := strconv.ParseFloat(Option["basePrice"], 64)
	if err != nil {
		logger.Debugln("config item basePrice is not float")
		return false
	}

	fluctuation, err := strconv.ParseFloat(Option["fluctuation"], 64)
	if err != nil {
		logger.Debugln("config item fluctuation is not float")
		return false
	}

	circular.PrevClosePrice = lastPrice

	logger.Infof("lastPrice %0.02f\n", lastPrice)
	if lastPrice >= basePrice+fluctuation {
		Sell()
	} else if lastPrice <= basePrice-fluctuation {
		Buy()
	}

	// do sell when price is below stoploss point
	processStoploss(lastPrice)

	processTimeout()

	return true
}
