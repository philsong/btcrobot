/*
  btcrobot is a Bitcoin, Litecoin and Altcoin trading bot written in golang,
  it features multiple trading methods using technical analysis.

  Disclaimer:

  USE AT YOUR OWN RISK!

  The author of this project is NOT responsible for any damage or loss caused
  by this software. There can be bugs and the bot may not perform as expected
  or specified. Please consider testing it first with paper trading /
  backtesting on historical data. Also look at the code to see what how
  it's working.

  Weibo:http://weibo.com/bocaicfa
*/

package strategy

import (
	"fmt"
)

type EMAMACDStrategy struct{}

func init() {
	emamacdStrategy := EMAMACDStrategy{}
	Register("emamacd", emamacdStrategy)
}

//xxx strategy
func (emamacdStrategy EMAMACDStrategy) Perform(tradeAPI TradeAPI, Time []string, Price []float64, Volumn []float64) bool {
	fmt.Println("哥哥，还没实现呢。。。正在搞。。")
	//实现自己的策略
	return false
}
