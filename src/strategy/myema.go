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
	"logger"
	"strconv"
)

type MyEMAStrategy struct {
	PrevEMACross      string
	PrevEMAdif        float64
	LessBuyThreshold  bool
	LessSellThreshold bool
}

func init() {
	myEMAStrategy := new(MyEMAStrategy)
	myEMAStrategy.PrevEMACross = "unknown"
	Register("MYEMA", myEMAStrategy)
}

type EMARecord struct {
	EMAShort float64
	EMALong  float64
	DIF      float64
}

func CalcEMA(price []float64, periodShort int, periodLong int) (ema []EMARecord) {
	ema = make([]EMARecord, len(price))

	ema[0].EMAShort = price[0]
	ema[0].EMALong = price[0]
	ema[0].DIF = 0

	for i := 1; i < len(price); i++ {
		ema[i].EMAShort = ema[i-1].EMAShort*(float64(periodShort-1))/(float64(periodShort+1)) +
			price[i]*2.0/(float64(periodShort+1))
		ema[i].EMALong = ema[i-1].EMALong*(float64(periodLong-1))/(float64(periodLong+1)) +
			price[i]*2.0/(float64(periodLong+1))
		ema[i].DIF = ema[i].EMAShort - ema[i].EMALong
	}

	return
}

func (emaStrategy *MyEMAStrategy) checkThreshold(tradeType string, EMAdif float64) bool {
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

func (emaStrategy *MyEMAStrategy) is_upcross(emaLast, emaCurr float64) bool {
	if (emaLast <= 0 || emaStrategy.PrevEMACross == "down") && is_up(emaCurr) {
		return true
	}

	return false
}

func (emaStrategy *MyEMAStrategy) is_downcross(emaLast, emaCurr float64) bool {
	if (emaLast >= 0 || emaStrategy.PrevEMACross == "up") && is_down(emaCurr) {
		return true
	}

	return false
}

func (myEMAStrategy *MyEMAStrategy) Tick(records []Record) bool {
	shortEMA, _ := strconv.Atoi(Option["shortEMA"])
	longEMA, _ := strconv.Atoi(Option["longEMA"])

	var Price []float64
	for _, v := range records {
		Price = append(Price, v.Close)
	}

	length := len(records)
	ema := CalcEMA(Price, shortEMA, longEMA)
	emaLast := ema[length-2].DIF
	emaCurr := ema[length-1].DIF

	if myEMAStrategy.PrevEMACross == "unknown" {
		if is_up(emaLast) {
			myEMAStrategy.PrevEMACross = "up"
		} else if is_down(emaLast) {
			myEMAStrategy.PrevEMACross = "down"
		} else {
			myEMAStrategy.PrevEMACross = "unknown"
		}

		logger.Infoln("prev cross is", myEMAStrategy.PrevEMACross)
	}

	if emaCurr != myEMAStrategy.PrevEMAdif {
		myEMAStrategy.PrevEMAdif = emaCurr
		logger.Infof("EMA [%0.02f, %0.02f, %0.02f] Diff:%0.04f \t%0.04f\n", records[length-1].Close, ema[length-1].EMAShort, ema[length-1].EMALong, emaLast, emaCurr)
	}

	// reset LessBuyThreshold LessSellThreshold flag when (^ or V) happen
	if myEMAStrategy.LessBuyThreshold && is_down(emaCurr) {
		myEMAStrategy.LessBuyThreshold = false
		myEMAStrategy.PrevEMACross = "down" // reset
		logger.Infoln("down->up(EMA diff < buy threshold)->down ^")
	}
	if myEMAStrategy.LessSellThreshold && is_up(emaCurr) {
		myEMAStrategy.LessSellThreshold = false
		myEMAStrategy.PrevEMACross = "up" // reset
		logger.Infoln("up->down(EMA diff > sell threshold)->up V")
	}

	// do buy when cross up
	if myEMAStrategy.is_upcross(emaLast, emaCurr) || myEMAStrategy.LessBuyThreshold {
		if Option["enable_trading"] == "1" && PrevTrade != "buy" {
			myEMAStrategy.PrevEMACross = "up"
			if myEMAStrategy.checkThreshold("buy", emaCurr) {
				Buy()
			}
		}
	}

	// do sell when cross down
	if myEMAStrategy.is_downcross(emaLast, emaCurr) || myEMAStrategy.LessSellThreshold {
		myEMAStrategy.PrevEMACross = "down"
		if Option["enable_trading"] == "1" && PrevTrade != "sell" {
			if myEMAStrategy.checkThreshold("sell", emaCurr) {
				Sell()
			}
		}
	}

	if PrevBuyPirce < records[length-1].High {
		PrevBuyPirce = records[length-1].High
	}

	processStoploss(lastPrice)

	return true
}
