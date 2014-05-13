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

type OOStrategy struct {
}

func init() {
	oo := new(OOStrategy)
	Register("OPENORDER", oo)
}

func (oo *OOStrategy) Tick(records []Record) bool {

	const btcslap = 0.2
	const ltcslap = 0.05
	const ordercount = 1

	ret, orderbook := GetOrderBook()
	if !ret {
		logger.Infoln("get orderbook failed 1")
		ret, orderbook = GetOrderBook() //try again
		if !ret {
			logger.Infoln("get orderbook failed 2")
			return false
		}
	}

	logger.Infoln("卖一", orderbook.Asks[len(orderbook.Asks)-1])
	logger.Infoln("买一", orderbook.Bids[0])

	diff := 1.0

	if orderbook.Bids[0].Price+diff <= orderbook.Asks[len(orderbook.Asks)-1].Price {
		for i := 1; i <= ordercount; i++ {
			buyID := Buy()
			if buyID != "0" {
				Sell()
			}
		}
	}

	processTimeout()

	return true
}
