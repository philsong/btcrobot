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
	. "config"
	"logger"
	"math"
	"service"
)

func getTrendAtIndex(emaShort, emaLong []float64, EMAMinThreshold float64, TresholdLevel int, i int) int {
	// This function return the calculated trend at index i, with respect to EMA-values,
	// thresholds and no of samples before triggering.
	// Return values:
	// 0		= no trend
	// 1/-1	= weak trend up/down (below thresholds)
	// 2/-2	= strong trend up/down (above thresholds)
	// 3/-3	= strong trend up/down and enough samples has passed (according to settings "Buy/Sell after X samples")

	if i < 5 {
		// All data not available
		return 0
	}

	var trend = 0
	var dif1 = getEMAdifAt(emaShort, emaLong, i)
	if dif1 > 0 {
		trend = 1
		if dif1 > EMAMinThreshold {
			trend = 2
			var dif2 = getEMAdifAt(emaShort, emaLong, i-1)
			var dif3 = getEMAdifAt(emaShort, emaLong, i-2)
			var dif4 = getEMAdifAt(emaShort, emaLong, i-3)
			var dif5 = getEMAdifAt(emaShort, emaLong, i-4)
			if (TresholdLevel == 1) ||
				(TresholdLevel == 2 && (dif2 > EMAMinThreshold)) ||
				(TresholdLevel == 3 && (dif2 > EMAMinThreshold) && (dif3 > EMAMinThreshold)) ||
				(TresholdLevel == 4 && (dif2 > EMAMinThreshold) && (dif3 > EMAMinThreshold) && (dif4 > EMAMinThreshold)) ||
				(TresholdLevel == 5 && (dif2 > EMAMinThreshold) && (dif3 > EMAMinThreshold) && (dif4 > EMAMinThreshold) && (dif5 > EMAMinThreshold)) {
				trend = 3
			}
		}
	} else if dif1 < 0 {
		trend = -1
		if dif1 < -EMAMinThreshold {
			trend = -2
			var dif2 = getEMAdifAt(emaShort, emaLong, i-1)
			var dif3 = getEMAdifAt(emaShort, emaLong, i-2)
			var dif4 = getEMAdifAt(emaShort, emaLong, i-3)
			var dif5 = getEMAdifAt(emaShort, emaLong, i-4)
			if (TresholdLevel == 1) ||
				(TresholdLevel == 2 && (dif2 < -EMAMinThreshold)) ||
				(TresholdLevel == 3 && (dif2 < -EMAMinThreshold) && (dif3 < -EMAMinThreshold)) ||
				(TresholdLevel == 4 && (dif2 < -EMAMinThreshold) && (dif3 < -EMAMinThreshold) && (dif4 < -EMAMinThreshold)) ||
				(TresholdLevel == 5 && (dif2 < -EMAMinThreshold) && (dif3 < -EMAMinThreshold) && (dif4 < -EMAMinThreshold) && (dif5 < -EMAMinThreshold)) {
				trend = -3
			}
		}
	}
	return trend
}

func (w *Huobi) findLatestSolidTrend(emaShort, emaLong []float64, EMAMinThreshold float64,
	TresholdLevel int, length int) {
	for i := length - 2; i >= 4; i-- {
		trend := getTrendAtIndex(emaShort, emaLong, EMAMinThreshold, TresholdLevel, i)
		if math.Abs(float64(trend)) == 3 {
			w.latestSolidTrend = trend
			break
		}
	}

	logger.Infoln("Latest solid trend: ", w.latestSolidTrend)
}

func (w *Huobi) trade(emaShort, emaLong []float64, EMAMinThreshold float64,
	TresholdLevel int, length int, tradeOnlyAfterSwitch int, tradeAmount string) {
	currentTrend := getTrendAtIndex(emaShort, emaLong, EMAMinThreshold, TresholdLevel, length-1)
	logger.Debugln("currentTrend is ", currentTrend)

	if currentTrend > 1 {
		// Trend is up
		if currentTrend == 3 {
			// Trend is up, also according to the "Buy after X samples"-setting
			if (tradeOnlyAfterSwitch == 1) && (w.latestSolidTrend == 3) {
				// tradeOnlyAfterSwitch==true but the trend has not switched: Don't trade
				logger.Debugln("Trend has not switched (still up). The setting \"tradeOnlyAfterSwitch==true\", so do not trade...")
				return
			}
			w.latestSolidTrend = 3

			if Option["disable_trading"] == "1" {
				logger.Debugln("Simulted BUY (Simulation only: no trade was made)")
			} else {
				logger.Infoln("Trend has switched, 探测到买入点")
				go service.TriggerTrender("探测到买入点")

				w.Do_buy(w.getTradePrice("buy"), tradeAmount)
			}
			//logger.Infoln("Trend is up, but no " + currency + " to spend...");
		} else {
			logger.Debugf("Trend is up, but not for long enough (needs to be \"up\" for at least %d samples)\n", TresholdLevel)
		}
	} else if currentTrend < -1 {
		// Trend is down
		if currentTrend == -3 {
			// Trend is down, also according to the "Sell after X samples"-setting
			if (tradeOnlyAfterSwitch == 1) && (w.latestSolidTrend == -3) {
				// tradeOnlyAfterSwitch==true but the trend has not switched: Don't trade
				logger.Debugln("Trend has not switched (still down). The setting \"tradeOnlyAfterSwitch==true\", so do not trade...")
				return
			}
			w.latestSolidTrend = -3

			if Option["disable_trading"] == "1" {
				logger.Infoln("Simulted SELL (Simulation only: no trade was made)")
			} else {
				logger.Infoln("Trend has switched, 探测到卖出点")
				go service.TriggerTrender("探测到卖出点")

				w.Do_sell(w.getTradePrice("sell"), tradeAmount)
			}
			//logger.Infoln("Trend is down, but no BTC to sell...");
		} else {
			logger.Debugf("Trend is down, but not for long enough (needs to be \"down\" for at least t %d samples)\n", TresholdLevel)
		}
	} else {
		logger.Debugln("Trend is undefined/weak")
	}
}
