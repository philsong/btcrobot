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

func (emaStrategy *EMAStrategy) checkThreshold(direction string, EMAdif float64) bool {
	if direction == "buy" {
		buyThreshold, err := strconv.ParseFloat(Option["buyThreshold"], 64)
		if err != nil {
			logger.Errorln("config item buyThreshold is not float")
			return false
		}

		if EMAdif > buyThreshold {
			logger.Infof("EMAdif(%0.03f) > buyThreshold(%0.03f), trigger to buy\n", EMAdif, buyThreshold)
			emaStrategy.LessBuyThreshold = false
			return true
		} else {
			if emaStrategy.LessBuyThreshold == false {
				logger.Infof("cross up, but EMAdif(%0.03f) <= buyThreshold(%0.03f)\n", EMAdif, buyThreshold)
				emaStrategy.LessBuyThreshold = true
			}
		}
	} else {
		sellThreshold, err := strconv.ParseFloat(Option["sellThreshold"], 64)
		if err != nil {
			logger.Errorln("config item sellThreshold is not float")
			return false
		}

		if sellThreshold > 0 {
			sellThreshold = -sellThreshold
		}

		if EMAdif < sellThreshold {
			logger.Infof("EMAdif(%0.03f) <  sellThreshold(%0.03f), trigger to sell\n", EMAdif, sellThreshold)
			emaStrategy.LessSellThreshold = false
			return true
		} else {
			if emaStrategy.LessSellThreshold == false {
				logger.Infof("cross down, but EMAdif(%0.03f) >= sellThreshold(%0.03f)\n", EMAdif, sellThreshold)
				emaStrategy.LessSellThreshold = true
			}
		}
	}

	return false
}

/* test cases
2014/03/23 07:17:58 EMA Diff:-0.002	-0.000	Price:96.41
2014/03/23 07:18:06 EMA Diff:-0.002	0.003	Price:96.44
2014/03/23 07:18:12 EMA Diff:-0.002	-0.000	Price:96.41
2014/03/23 07:18:56 EMA Diff:-0.000	0.014	Price:96.52
*/

func is_uptrend(ema float64) bool {
	if ema > 0.000001 {
		return true
	} else {
		return false
	}
}

func is_downtrend(ema float64) bool {
	if ema < -0.000001 {
		return true
	} else {
		return false
	}
}

func (emaStrategy *EMAStrategy) is_upcross(prevema, ema float64) bool {
	if is_uptrend(ema) {
		if prevema <= 0 || emaStrategy.PrevEMACross == "down" {
			return true
		}
	}

	return false
}

func (emaStrategy *EMAStrategy) is_downcross(prevema, ema float64) bool {
	if is_downtrend(ema) {
		if prevema >= 0 || emaStrategy.PrevEMACross == "up" {
			return true
		}
	}

	return false
}

//EMA strategy
func (emaStrategy *EMAStrategy) Tick(records []Record) bool {
	//read config
	shortEMA, _ := strconv.Atoi(Option["shortEMA"])
	longEMA, _ := strconv.Atoi(Option["longEMA"])

	var Time []string
	var Price []float64
	var Volumn []float64
	for _, v := range records {
		Time = append(Time, v.TimeStr)
		Price = append(Price, v.Close)
		Volumn = append(Volumn, v.Volumn)
		//Price = append(Price, (v.Close+v.Open+v.High+v.Low)/4.0)
		//Price = append(Price, v.Low)
	}

	//compute the indictor
	emaShort := EMA(Price, shortEMA)
	emaLong := EMA(Price, longEMA)
	EMAdif := getMACDdif(emaShort, emaLong)

	length := len(Price)
	if emaStrategy.PrevEMACross == "unknown" {
		if is_uptrend(EMAdif[length-3]) {
			emaStrategy.PrevEMACross = "up"
		} else if is_downtrend(EMAdif[length-3]) {
			emaStrategy.PrevEMACross = "down"
		} else {
			emaStrategy.PrevEMACross = "unknown"
		}
		logger.Infoln("prev cross is", emaStrategy.PrevEMACross)
		if is_uptrend(EMAdif[length-3]) {
			logger.Infoln("上一个趋势是上涨，等待卖出点触发")
		} else if is_downtrend(EMAdif[length-3]) {
			logger.Infoln("上一个趋势是下跌，等待买入点触发")
		} else {
			logger.Infoln("上一个趋势是unknown。。。")
		}
	}

	//go TriggerPrice(Price[length-1])
	if EMAdif[length-1] != emaStrategy.PrevEMAdif {
		emaStrategy.PrevEMAdif = EMAdif[length-1]
		logger.Infof("EMA [%0.02f,%0.02f,%0.02f] Diff:%0.03f\t%0.03f\n", Price[length-1], emaShort[length-1], emaLong[length-1], EMAdif[length-2], EMAdif[length-1])
	}

	//reset LessBuyThreshold LessSellThreshold flag when (^ or V) happen
	if emaStrategy.LessBuyThreshold && is_downtrend(EMAdif[length-1]) {
		emaStrategy.LessBuyThreshold = false
		emaStrategy.PrevEMACross = "down" //reset
		logger.Infoln("down->up(EMA diff < buy threshold)->down ^")

	}
	if emaStrategy.LessSellThreshold && is_uptrend(EMAdif[length-1]) {
		emaStrategy.LessSellThreshold = false
		emaStrategy.PrevEMACross = "up" //reset
		logger.Infoln("up->down(EMA diff > sell threshold)->up V")
	}

	//EMA cross
	if (emaStrategy.is_upcross(EMAdif[length-2], EMAdif[length-1]) || emaStrategy.LessBuyThreshold) ||
		(emaStrategy.is_downcross(EMAdif[length-2], EMAdif[length-1]) || emaStrategy.LessSellThreshold) { //up cross

		//do buy when cross up
		if emaStrategy.is_upcross(EMAdif[length-2], EMAdif[length-1]) || emaStrategy.LessBuyThreshold {
			if Option["enable_trading"] == "1" && PrevTrade != "buy" {
				emaStrategy.PrevEMACross = "up"
				if emaStrategy.checkThreshold("buy", EMAdif[length-1]) {
					Buy()
				}
			}
		}

		//do sell when cross down
		if emaStrategy.is_downcross(EMAdif[length-2], EMAdif[length-1]) || emaStrategy.LessSellThreshold {
			emaStrategy.PrevEMACross = "down"
			if Option["enable_trading"] == "1" && PrevTrade != "sell" {
				if emaStrategy.checkThreshold("sell", EMAdif[length-1]) {
					Sell()
				}
			}
		}

		//backup the kline data for analyze
		if Config["env"] == "dev" {
			backup(Time[length-1])
		}
	}

	//do sell when price is below stoploss point
	processStoploss(Price)

	processTimeout()

	return true
}

//for backup the kline file to detect the huobi bug
func backup(Time string) {

	oldFile := "cache/TradeKLine_minute.data"
	newFile := fmt.Sprintf("%s_%s", oldFile, Time)
	os.Rename(oldFile, newFile)
}
