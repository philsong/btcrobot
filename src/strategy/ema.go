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
	"logger"
	"os"
	"strconv"
)

type EMAStrategy struct {
	PrevEMACross      string
	PrevEMAdif        float64
	LessBuyThreshold  bool
	LessSellThreshold bool
}

func init() {
	emaStrategy := new(EMAStrategy)
	emaStrategy.PrevEMACross = "unknown"
	Register("EMA", emaStrategy)
}

func (emaStrategy *EMAStrategy) checkThreshold(tradeType string, EMAdif float64) bool {
	if tradeType == "buy" {
		buyThreshold := toFloat(Option["buyThreshold"])

		if EMAdif > buyThreshold {
			logger.Infof("EMAdif(%0.04f) > buyThreshold(%0.04f), trigger to buy\n", EMAdif, buyThreshold)
			emaStrategy.LessBuyThreshold = false
			return true
		} else {
			if emaStrategy.LessBuyThreshold == false {
				logger.Infof("cross up, but EMAdif(%0.04f) <= buyThreshold(%0.04f)\n", EMAdif, buyThreshold)
				emaStrategy.LessBuyThreshold = true
			}
		}
	} else {
		sellThreshold := toFloat(Option["sellThreshold"])

		if sellThreshold > 0 {
			sellThreshold = -sellThreshold
		}

		if EMAdif < sellThreshold {
			logger.Infof("EMAdif(%0.04f) <  sellThreshold(%0.04f), trigger to sell\n", EMAdif, sellThreshold)
			emaStrategy.LessSellThreshold = false
			return true
		} else {
			if emaStrategy.LessSellThreshold == false {
				logger.Infof("cross down, but EMAdif(%0.04f) >= sellThreshold(%0.04f)\n", EMAdif, sellThreshold)
				emaStrategy.LessSellThreshold = true
			}
		}
	}

	return false
}

func is_up(ema float64) bool {
	if ema > 0.000001 {
		return true
	} else {
		return false
	}
}

func is_down(ema float64) bool {
	if ema < -0.000001 {
		return true
	} else {
		return false
	}
}

func (emaStrategy *EMAStrategy) is_upcross(emaLast, emaCurr float64) bool {
	if (emaLast <= 0 || emaStrategy.PrevEMACross == "down") && is_up(emaCurr) {
		return true
	}

	return false
}

func (emaStrategy *EMAStrategy) is_downcross(emaLast, emaCurr float64) bool {
	if (emaLast >= 0 || emaStrategy.PrevEMACross == "up") && is_down(emaCurr) {
		return true
	}

	return false
}

// EMA strategy
func (emaStrategy *EMAStrategy) Tick(records []Record) bool {
	// read config
	shortEMA, _ := strconv.Atoi(Option["shortEMA"])
	longEMA, _ := strconv.Atoi(Option["longEMA"])

	var Price []float64
	for _, v := range records {
		Price = append(Price, v.Close)
	}

	// compute the indictor
	emaShort := EMA(Price, shortEMA)
	emaLong := EMA(Price, longEMA)
	EMAdif := getMACDdif(emaShort, emaLong)
	emaLast := EMAdif[length-3]
	emaCurr := EMAdif[length-2]

	if emaStrategy.PrevEMACross == "unknown" {
		if is_up(emaLast) {
			emaStrategy.PrevEMACross = "up"
		} else if is_down(EMAdif[length-3]) {
			emaStrategy.PrevEMACross = "down"
		} else {
			emaStrategy.PrevEMACross = "unknown"
		}

		logger.Infoln("prev cross is", emaStrategy.PrevEMACross)
	}

	//go TriggerPrice(Price[length-1])
	if emaCurr != emaStrategy.PrevEMAdif {
		emaStrategy.PrevEMAdif = emaCurr
		logger.Infof("EMA [%0.02f,%0.02f,%0.02f] Diff:%0.04f\t%0.04f\n", lastPrice, emaShort[length-1], emaLong[length-1], emaLast, emaCurr)
	}

	ret, orderBook := GetOrderBook()
	if !ret {
		logger.Infoln("get order book failed")
	} else {
		askstotal := 0.0
		for i := 0; i < len(orderBook.Asks); i++ {
			askstotal += orderBook.Asks[i].Amount
		}
		bidstotal := 0.0
		for i := 0; i < len(orderBook.Bids); i++ {
			bidstotal += orderBook.Bids[i].Amount
		}

		logger.Infoln("sell:buy", askstotal, bidstotal)
	}

	// reset LessBuyThreshold LessSellThreshold flag when (^ or V) happen
	if emaStrategy.LessBuyThreshold && is_down(emaCurr) {
		emaStrategy.LessBuyThreshold = false
		emaStrategy.PrevEMACross = "down" // reset
		logger.Infoln("down->up(EMA diff < buy threshold)->down ^")
	}
	if emaStrategy.LessSellThreshold && is_up(emaCurr) {
		emaStrategy.LessSellThreshold = false
		emaStrategy.PrevEMACross = "up" // reset
		logger.Infoln("up->down(EMA diff > sell threshold)->up V")
	}

	// do buy when cross up
	if emaStrategy.is_upcross(emaLast, emaCurr) || emaStrategy.LessBuyThreshold {
		if Option["enable_trading"] == "1" && PrevTrade != "buy" {
			emaStrategy.PrevEMACross = "up"
			if emaStrategy.checkThreshold("buy", emaCurr) {
				Buy()
			}
		}
	}

	// do sell when cross down
	if emaStrategy.is_downcross(emaLast, emaCurr) || emaStrategy.LessSellThreshold {
		emaStrategy.PrevEMACross = "down"
		if Option["enable_trading"] == "1" && PrevTrade != "sell" {
			if emaStrategy.checkThreshold("sell", emaCurr) {
				Sell()
			}
		}
	}

	// EMA cross
	if (emaStrategy.is_upcross(emaLast, emaCurr) || emaStrategy.LessBuyThreshold) ||
		(emaStrategy.is_downcross(emaLast, emaCurr) || emaStrategy.LessSellThreshold) {
		// backup the kline data for analyze
		if Config["env"] == "dev" {
			backup(records[length-1].TimeStr)
		}
	}

	// do sell when price is below stoploss point
	processStoploss(lastPrice)

	processTimeout()

	return true
}

// for backup the kline file to detect the huobi bug
func backup(Time string) {
	oldFile := "cache/TradeKLine_minute.data"
	newFile := fmt.Sprintf("%s_%s", oldFile, Time)
	os.Rename(oldFile, newFile)
}
