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
)

type arbitrageStrategy struct {
	PrevClosePrice float64
}

func init() {
	arbitrage := new(arbitrageStrategy)
	Register("arbitrage", arbitrage)
}

// arbitrage strategy
func (arbitrage *arbitrageStrategy) Tick(records []Record) bool {
	const btcslap = 1.8
	const ltcslap = 0.8

	sell1, buy1, ret := getOrderPrice()
	if !ret {
		return false
	}
	/*
		sell2, buy2, ret := getOrderPrice()
		if !ret {
			return false
		}
	*/
	diff := btcslap

	if buy1+diff <= sell1 {
		buyID := Buy()
		if buyID != "0" {
			Sell()
		}
	}

	processTimeout()

	return true
}
