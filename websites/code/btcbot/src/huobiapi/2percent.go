/*
  btcbot is a Bitcoin trading bot for HUOBI.com written
  in golang, it features multiple trading methods using
  technical analysis.

  Disclaimer:

  USE AT YOUR OWN RISK!

  The author of this project is NOT responsible for any damage or loss caused
  by this software. There can be bugs and the bot may not perform as expected
  or specified. Please consider testing it first with paper trading /
  backtesting on historical data. Also look at the code to see what how
  it's working.

  Weibo:http://weibo.com/bocaicfa
*/

package huobiapi

import (
	"logger"
)

func (w *Huobi) do2Percent(Time []string, Price []float64) {
	lowpoint := Price[0]
	highpoint := Price[0]
	lastTrade := "init"

	factor := 0.001
	logger.OverrideStart(w.Peroid)
	logger.Overrideln("start")
	for i := 0; i < len(Price); i++ {
		//logger.Overrideln("", i, lastTrade, Time[i], Price[i])
		if lastTrade != "buy" && Price[i]-lowpoint > 0 && Price[i]-lowpoint > factor*lowpoint {
			logger.Overrideln("++", i, Time[i], Price[i], lowpoint, highpoint, Price[i]-lowpoint, factor*lowpoint)

			highpoint = Price[i]
			lowpoint = Price[i]
			lastTrade = "buy"
		} else if lastTrade != "sell" && Price[i]-highpoint < 0 && Price[i]-highpoint < -factor*highpoint {
			logger.Overrideln("--", i, Time[i], Price[i], lowpoint, highpoint, Price[i]-highpoint, -factor*highpoint)

			highpoint = Price[i]
			lowpoint = Price[i]
			lastTrade = "sell"
		} else {
			if highpoint < Price[i] {
				highpoint = Price[i]
			} else if lowpoint > Price[i] {
				lowpoint = Price[i]
			}
		}
	}
	logger.Overrideln("end")

	return
}
