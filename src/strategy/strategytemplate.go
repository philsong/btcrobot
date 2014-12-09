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
	"fmt"
)

type XXXStrategy struct{}

func init() {
	xxxStrategy := new(XXXStrategy)
	Register("xxx", xxxStrategy)
}

// xxx strategy
func (xxxStrategy *XXXStrategy) Tick(records []Record) bool {
	fmt.Println("empty strgatey template, you can realize your own trade strategy in here")
	// 实现自己的策略
	return false
}
