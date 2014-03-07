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

func (w *Huobi) do2Percent(xData []string, yData []float64) {
	lowpoint := yData[0]
	highpoint := yData[0]
	lastTrade := "init"

	factor := 0.001
	logger.OverrideStart(w.Peroid)
	logger.Overrideln("start")
	for i := 0; i < len(yData); i++ {
		//logger.Overrideln("", i, lastTrade, xData[i], yData[i])
		if lastTrade != "buy" && yData[i]-lowpoint > 0 && yData[i]-lowpoint > factor*lowpoint {
			logger.Overrideln("++", i, xData[i], yData[i], lowpoint, highpoint, yData[i]-lowpoint, factor*lowpoint)

			highpoint = yData[i]
			lowpoint = yData[i]
			lastTrade = "buy"
		} else if lastTrade != "sell" && yData[i]-highpoint < 0 && yData[i]-highpoint < -factor*highpoint {
			logger.Overrideln("--", i, xData[i], yData[i], lowpoint, highpoint, yData[i]-highpoint, -factor*highpoint)

			highpoint = yData[i]
			lowpoint = yData[i]
			lastTrade = "sell"
		} else {
			if highpoint < yData[i] {
				highpoint = yData[i]
			} else if lowpoint > yData[i] {
				lowpoint = yData[i]
			}
		}
	}
	logger.Overrideln("end")

	return
}
