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

type HLCrossStrategy struct {
	PrevClosePrice float64
	PrevHighPrice  float64
	PrevLowPrice   float64
}

func init() {
	HLCross := new(HLCrossStrategy)
	Register("HLCross", HLCross)
}

// HLCross strategy
func (HLCross *HLCrossStrategy) Tick(records []Record) bool {
	// read config
	shortEMA, _ := strconv.Atoi(Option["shortEMA"])

	var Price []float64
	for _, v := range records {
		Price = append(Price, v.Close)
	}

	// compute the indictor
	emaShort := EMA(Price, shortEMA)

	if HLCross.PrevClosePrice != records[length-1].Close ||
		HLCross.PrevHighPrice != records[length-2].High ||
		HLCross.PrevLowPrice != records[length-2].Low {
		HLCross.PrevClosePrice = records[length-1].Close
		HLCross.PrevHighPrice = records[length-2].High
		HLCross.PrevLowPrice = records[length-2].Low

		logger.Infof("lastPrice %0.02f prevHigh %0.02f prevLow %0.02f\n",
			lastPrice, records[length-2].High, records[length-2].Low)
	}

	// HLCross cross
	if Price[length-2] > emaShort[length-2] &&
		records[length-2].Volumn > 500 &&
		records[length-2].High > records[length-3].High &&
		records[length-2].Low > records[length-3].Low {
		Buy()
	} else if Price[length-2] < emaShort[length-2] &&
		(records[length-2].High < records[length-3].High ||
			records[length-2].Low < records[length-3].Low) {
		Sell()
	}

	// do sell when price is below stoploss point
	processStoploss(lastPrice)

	processTimeout()

	return true
}
