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
	"email"
	"fmt"
	"logger"
	"strconv"
)

type EMAMACDStrategy struct {
	PrevMACDTrade string
	PrevMACDdif   float64

	PrevEMATrade      string
	PrevEMACross      string
	PrevEMAdif        float64
	PrevBuyPirce      float64
	LessBuyThreshold  bool
	LessSellThreshold bool
}

func init() {
	emamacdStrategy := new(EMAMACDStrategy)
	emamacdStrategy.PrevEMACross = "unknown"
	Register("EMAMACD", emamacdStrategy)
}

func (emamacdStrategy *EMAMACDStrategy) checkThreshold(direction string, EMAdif float64) bool {
	if direction == "buy" {
		buyThreshold, err := strconv.ParseFloat(Option["buyThreshold"], 64)
		if err != nil {
			logger.Errorln("config item buyThreshold is not float")
			return false
		}

		if EMAdif > buyThreshold {
			logger.Infof("EMAdif(%0.03f) > buyThreshold(%0.03f), trigger to buy\n", EMAdif, buyThreshold)
			emamacdStrategy.LessBuyThreshold = false
			return true
		} else {
			if emamacdStrategy.LessBuyThreshold == false {
				logger.Infof("cross up, but EMAdif(%0.03f) <= buyThreshold(%0.03f)\n", EMAdif, buyThreshold)
				emamacdStrategy.LessBuyThreshold = true
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
			emamacdStrategy.LessSellThreshold = false
			return true
		} else {
			if emamacdStrategy.LessSellThreshold == false {
				logger.Infof("cross down, but EMAdif(%0.03f) >= sellThreshold(%0.03f)\n", EMAdif, sellThreshold)
				emamacdStrategy.LessSellThreshold = true
			}
		}
	}

	return false
}

func (emamacdStrategy *EMAMACDStrategy) is_upcross(prevema, ema float64) bool {
	if is_uptrend(ema) {
		if prevema <= 0 || emamacdStrategy.PrevEMACross == "down" {
			return true
		}
	}

	return false
}

func (emamacdStrategy *EMAMACDStrategy) is_downcross(prevema, ema float64) bool {
	if is_downtrend(ema) {
		if prevema >= 0 || emamacdStrategy.PrevEMACross == "up" {
			return true
		}
	}

	return false
}

//EMA strategy
func (emamacdStrategy *EMAMACDStrategy) Tick(records []Record) bool {
	//read config
	shortEMA, _ := strconv.Atoi(Option["shortEMA"])
	longEMA, _ := strconv.Atoi(Option["longEMA"])
	signalPeriod, _ := strconv.Atoi(Option["signalPeriod"])

	tradeAmount := Option["tradeAmount"]

	stoploss, err := strconv.ParseFloat(Option["stoploss"], 64)
	if err != nil {
		logger.Errorln("config item stoploss is not float")
		return false
	}

	MACDsellThreshold, err := strconv.ParseFloat(Option["MACDsellThreshold"], 64)
	if err != nil {
		logger.Errorln("config item MACDsellThreshold is not float")
		return false
	}

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

	MACDdif := getMACDdif(emaShort, emaLong)
	MACDSignal := getMACDSignal(MACDdif, signalPeriod)
	MACDHistogram := getMACDHistogram(MACDdif, MACDSignal)

	length := len(Price)
	if emamacdStrategy.PrevEMACross == "unknown" {
		if is_uptrend(EMAdif[length-3]) {
			emamacdStrategy.PrevEMACross = "up"
		} else if is_downtrend(EMAdif[length-3]) {
			emamacdStrategy.PrevEMACross = "down"
		} else {
			emamacdStrategy.PrevEMACross = "unknown"
		}
		logger.Infoln("prev cross is", emamacdStrategy.PrevEMACross)
		if is_uptrend(EMAdif[length-3]) {
			logger.Infoln("上一个趋势是上涨，等待卖出点触发")
		} else if is_downtrend(EMAdif[length-3]) {
			logger.Infoln("上一个趋势是下跌，等待买入点触发")
		} else {
			logger.Infoln("上一个趋势是unknown。。。")
		}
	}

	//go TriggerPrice(Price[length-1])
	if EMAdif[length-1] != emamacdStrategy.PrevEMAdif {
		emamacdStrategy.PrevEMAdif = EMAdif[length-1]
		logger.Infof("EMA [%0.02f,%0.02f,%0.02f] Diff:%0.03f\t%0.03f\n", Price[length-1], emaShort[length-1], emaLong[length-1], EMAdif[length-2], EMAdif[length-1])
	}

	if MACDdif[length-1] != emamacdStrategy.PrevMACDdif {
		emamacdStrategy.PrevMACDdif = MACDdif[length-1]
		logger.Infof("MACD:d=%5.03f\ts=%5.03f\th=%5.03f\tpre-h=%5.03f\n", MACDdif[length-1], MACDSignal[length-1], MACDHistogram[length-1], MACDHistogram[length-2])
	}

	//reset LessBuyThreshold LessSellThreshold flag when (^ or V) happen
	if emamacdStrategy.LessBuyThreshold && is_downtrend(EMAdif[length-1]) {
		emamacdStrategy.LessBuyThreshold = false
		emamacdStrategy.PrevEMACross = "down" //reset
		logger.Infoln("down->up(EMA diff < buy threshold)->down ^")

	}
	if emamacdStrategy.LessSellThreshold && is_uptrend(EMAdif[length-1]) {
		emamacdStrategy.LessSellThreshold = false
		emamacdStrategy.PrevEMACross = "up" //reset
		logger.Infoln("up->down(EMA diff > sell threshold)->up V")
	}

	//EMA cross
	if emamacdStrategy.is_upcross(EMAdif[length-2], EMAdif[length-1]) || emamacdStrategy.LessBuyThreshold { //up cross

		//do buy when cross up
		if emamacdStrategy.is_upcross(EMAdif[length-2], EMAdif[length-1]) || emamacdStrategy.LessBuyThreshold {
			if Option["enable_trading"] == "1" && emamacdStrategy.PrevEMATrade != "buy" {

				emamacdStrategy.PrevEMACross = "up"

				if emamacdStrategy.checkThreshold("buy", EMAdif[length-1]) {
					emamacdStrategy.PrevEMATrade = "buy"
					emamacdStrategy.PrevMACDTrade = "init"
					diff := fmt.Sprintf("%0.03f", EMAdif[length-1])
					warning := "EMA up cross, 买入buy In<----市价" + getTradePrice("", Price[length-1]) +
						",委托价" + getTradePrice("buy", Price[length-1]) + ",diff" + diff
					logger.Infoln(warning)
					if Buy(getTradePrice("buy", Price[length-1]), tradeAmount) != "0" {
						emamacdStrategy.PrevBuyPirce = Price[length-1]
						warning += "[委托成功]"
					} else {
						warning += "[委托失败]"
					}

					go email.TriggerTrender(warning)
				}
			}
		}

		//backup the kline data for analyze
		if Config["env"] == "dev" {
			backup(Time[length-1])
		}
	}

	//macd cross
	if EMAdif[length-1] > 0 || emamacdStrategy.PrevEMATrade == "buy" {
		if (MACDHistogram[length-2] > 0.000001 && MACDHistogram[length-1] < MACDsellThreshold) &&
			emamacdStrategy.PrevMACDTrade != "sell" {
			if Option["enable_trading"] == "1" && emamacdStrategy.PrevMACDTrade != "sell" {
				emamacdStrategy.PrevMACDTrade = "sell"
				emamacdStrategy.PrevBuyPirce = 0
				emamacdStrategy.PrevEMATrade = "sell"

				histogram := fmt.Sprintf("%0.03f", MACDHistogram[length-1])
				warning := "MACD down cross, 卖出Sell Out---->市价" + getTradePrice("", Price[length-1]) +
					",委托价" + getTradePrice("sell", Price[length-1]) + ",histogram" + histogram
				logger.Infoln(warning)
				if Sell(getTradePrice("sell", Price[length-1]), tradeAmount) != "0" {
					warning += "[委托成功]"
				} else {
					warning += "[委托失败]"
				}

				go email.TriggerTrender(warning)
			}
		}
	}

	//do sell when price is below stoploss point
	if Price[length-1] < emamacdStrategy.PrevBuyPirce*(1-stoploss*0.01) {
		if Option["enable_trading"] == "1" && emamacdStrategy.PrevEMATrade != "sell" {

			warning := "stop loss, 卖出Sell Out---->市价" + getTradePrice("", Price[length-1]) + ",委托价" + getTradePrice("sell", Price[length-1])
			logger.Infoln(warning)

			if Sell(getTradePrice("sell", Price[length-1]), tradeAmount) != "0" {
				warning += "[委托成功]"
			} else {
				warning += "[委托失败]"
			}

			go email.TriggerTrender(warning)

			emamacdStrategy.PrevEMATrade = "sell"
			emamacdStrategy.PrevMACDTrade = "init"
			emamacdStrategy.PrevBuyPirce = 0
		}
	}

	return true
}
