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
	"fmt"
)

type OOStrategy struct {
}

func init() {
	oo := new(OOStrategy)
	Register("OPENORDER", oo)
}

func (oo *OOStrategy) Tick(records []Record) bool {
	const btcslap = 0.5
	const ltcslap = 0.8

	sell1, buy1, ret := getOrderPrice()
	if !ret {
		return false
	}

	diff := btcslap
	fmt.Println(buy1, diff, sell1)
	if buy1+diff <= sell1 {
		amount := Option["tradeAmount"]

		buyID := buy(toString(buy1+0.01), amount)
		if buyID != "0" {
			sellID := sell(toString(sell1-0.01), amount)
			for {
				fmt.Println(sellID, GetLastError(), "retry sell")
				if sellID == "0" && GetLastError() == 10001 {
					sellID = sell(toString(sell1-0.01), amount)
				} else {
					break
				}
			}
		}
	}

	processTimeout()

	return true
}
