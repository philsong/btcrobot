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

type HLCrossStrategy struct {
	PrevClosePrice float64
	PrevHighPrice  float64
	PrevLowPrice   float64
}

func init() {
	HLCross := new(HLCrossStrategy)
	Register("HLCross", HLCross)
}

//HLCross strategy
func (HLCross *HLCrossStrategy) Tick(records []Record) bool {

	if HLCross.PrevClosePrice != records[length-1].Close ||
		HLCross.PrevHighPrice != records[length-2].High ||
		HLCross.PrevLowPrice != records[length-2].Low {
		HLCross.PrevClosePrice = records[length-1].Close
		HLCross.PrevHighPrice = records[length-2].High
		HLCross.PrevLowPrice = records[length-2].Low

		logger.Infof("nowClose %0.02f prevHigh %0.02f prevLow %0.02f\n", records[length-1].Close, records[length-2].High, records[length-2].Low)
	}

	//HLCross cross
	if records[length-1].Close > records[length-2].High {
		logger.Infoln("HLCross up")
		Buy()
	} else if records[length-1].Close < records[length-2].Low {
		Sell()
	}

	//do sell when price is below stoploss point
	processStoploss(lastPrice)

	processTimeout()

	return true
}
