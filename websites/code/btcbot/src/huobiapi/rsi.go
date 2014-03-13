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

func getMomentum(closeToday, closeNDaysAgo float64) float64 {
	momentum := closeToday - closeNDaysAgo
	return momentum
}

func getROC(closeToday, closeNDaysAgo float64) float64 {
	momentum := closeToday - closeNDaysAgo
	roc := momentum / closeNDaysAgo
	return roc
}

func getUD(Price []float64) ([]float64, []float64) {
	length := len(Price)
	var u []float64 = make([]float64, length)
	var d []float64 = make([]float64, length)

	// Loop through the entire array.
	for i := 1; i < length; i++ {
		diff := Price[i] - Price[i-1]
		if diff > 0 {
			u[i] = diff
			d[i] = 0
		} else if diff < 0 {
			d[i] = -diff
			u[i] = 0
		}
	}

	return u, d
}

func getRSI(Price []float64, periods int) []float64 {
	var periodArr []float64
	length := len(Price)
	var rsi []float64 = make([]float64, length)

	// Loop through the entire array.
	for i := 0; i < length; i++ {
		// add points to the array.
		periodArr = append(periodArr, Price[i])

		// 1: Check if array is "filled" else create null point in line.
		// 2: Calculate average.
		// 3: Remove first value.
		if periods == len(periodArr) {
			u, d := getUD(periodArr)
			rs := arrayAvg(u) / arrayAvg(d)

			rsi[i] = 100 - 100.0/(1+rs)

			// remove first value in array.
			periodArr = periodArr[1:]
		} else {
			rsi[i] = 0
		}
	}

	return rsi
}

func rsi(Price []float64) {
	rsShort := getRSI(Price, 6)
	rsLong := getRSI(Price, 12)

	length := len(rsShort)

	// Loop through the entire array.
	for i := 0; i < length; i++ {
		logger.Infoln(rsShort[i], rsLong[i])
	}
}
