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

type EMAMACDEMAStrategy struct {
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
	emamacdemaStrategy := new(EMAMACDEMAStrategy)
	emamacdemaStrategy.PrevEMACross = "unknown"
	Register("EMAMACDEMA", emamacdemaStrategy)
}

func (emamacdemaStrategy *EMAMACDEMAStrategy) checkThreshold(direction string, EMAdif float64) bool {
	if direction == "buy" {
		buyThreshold, err := strconv.ParseFloat(Option["buyThreshold"], 64)
		if err != nil {
			logger.Errorln("config item buyThreshold is not float")
			return false
		}

		if EMAdif > buyThreshold {
			logger.Infof("EMAdif(%0.03f) > buyThreshold(%0.03f), trigger to buy\n", EMAdif, buyThreshold)
			emamacdemaStrategy.LessBuyThreshold = false
			return true
		} else {
			if emamacdemaStrategy.LessBuyThreshold == false {
				logger.Infof("cross up, but EMAdif(%0.03f) <= buyThreshold(%0.03f)\n", EMAdif, buyThreshold)
				emamacdemaStrategy.LessBuyThreshold = true
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
			emamacdemaStrategy.LessSellThreshold = false
			return true
		} else {
			if emamacdemaStrategy.LessSellThreshold == false {
				logger.Infof("cross down, but EMAdif(%0.03f) >= sellThreshold(%0.03f)\n", EMAdif, sellThreshold)
				emamacdemaStrategy.LessSellThreshold = true
			}
		}
	}

	return false
}

func (emamacdemaStrategy *EMAMACDEMAStrategy) is_upcross(prevema, ema float64) bool {
	if is_up(ema) {
		if prevema <= 0 || emamacdemaStrategy.PrevEMACross == "down" {
			return true
		}
	}

	return false
}

func (emamacdemaStrategy *EMAMACDEMAStrategy) is_downcross(prevema, ema float64) bool {
	if is_down(ema) {
		if prevema >= 0 || emamacdemaStrategy.PrevEMACross == "up" {
			return true
		}
	}

	return false
}

// EMA strategy
func (emamacdemaStrategy *EMAMACDEMAStrategy) Tick(records []Record) bool {
	// read config
	shortEMA, _ := strconv.Atoi(Option["shortEMA"])
	longEMA, _ := strconv.Atoi(Option["longEMA"])
	signalPeriod, _ := strconv.Atoi(Option["signalPeriod"])

	nTradeAmount, err := strconv.ParseFloat(Option["tradeAmount"], 64)
	if err != nil {
		logger.Errorln("config item tradeAmount is not float")
		return false
	}

	MacdTradeAmount := fmt.Sprintf("%0.02f", 0.5*nTradeAmount)

	tradeAmount := Option["tradeAmount"]

	stoploss, err := strconv.ParseFloat(Option["stoploss"], 64)
	if err != nil {
		logger.Errorln("config item stoploss is not float")
		return false
	}

	MACDbuyThreshold, err := strconv.ParseFloat(Option["MACDbuyThreshold"], 64)
	if err != nil {
		logger.Errorln("config item MACDbuyThreshold is not float")
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

	// compute the indictor
	emaShort := EMA(Price, shortEMA)
	emaLong := EMA(Price, longEMA)
	EMAdif := getMACDdif(emaShort, emaLong)

	MACDdif := getMACDdif(emaShort, emaLong)
	MACDSignal := getMACDSignal(MACDdif, signalPeriod)
	MACDHistogram := getMACDHistogram(MACDdif, MACDSignal)

	length := len(Price)
	if emamacdemaStrategy.PrevEMACross == "unknown" {
		if is_up(EMAdif[length-3]) {
			emamacdemaStrategy.PrevEMACross = "up"
		} else if is_down(EMAdif[length-3]) {
			emamacdemaStrategy.PrevEMACross = "down"
		} else {
			emamacdemaStrategy.PrevEMACross = "unknown"
		}
		logger.Infoln("prev cross is", emamacdemaStrategy.PrevEMACross)
		if is_up(EMAdif[length-3]) {
			logger.Infoln("上一个趋势是上涨，等待卖出点触发")
		} else if is_down(EMAdif[length-3]) {
			logger.Infoln("上一个趋势是下跌，等待买入点触发")
		} else {
			logger.Infoln("上一个趋势是unknown。。。")
		}
	}

	if EMAdif[length-1] != emamacdemaStrategy.PrevEMAdif {
		emamacdemaStrategy.PrevEMAdif = EMAdif[length-1]
		logger.Infof("EMA [%0.02f,%0.02f,%0.02f] Diff:%0.03f\t%0.03f\n", Price[length-1], emaShort[length-1], emaLong[length-1], EMAdif[length-2], EMAdif[length-1])
	}

	if MACDdif[length-1] != emamacdemaStrategy.PrevMACDdif {
		emamacdemaStrategy.PrevMACDdif = MACDdif[length-1]
		logger.Infof("MACD:d=%5.03f\ts=%5.03f\th=%5.03f\tpre-h=%5.03f\n", MACDdif[length-1], MACDSignal[length-1], MACDHistogram[length-1], MACDHistogram[length-2])
	}

	// reset LessBuyThreshold LessSellThreshold flag when (^ or V) happen
	if emamacdemaStrategy.LessBuyThreshold && is_down(EMAdif[length-1]) {
		emamacdemaStrategy.LessBuyThreshold = false
		emamacdemaStrategy.PrevEMACross = "down" //reset
		logger.Infoln("down->up(EMA diff < buy threshold)->down ^")

	}
	if emamacdemaStrategy.LessSellThreshold && is_up(EMAdif[length-1]) {
		emamacdemaStrategy.LessSellThreshold = false
		emamacdemaStrategy.PrevEMACross = "up" //reset
		logger.Infoln("up->down(EMA diff > sell threshold)->up V")
	}

	// EMA cross
	if (emamacdemaStrategy.is_upcross(EMAdif[length-2], EMAdif[length-1]) || emamacdemaStrategy.LessBuyThreshold) ||
		(emamacdemaStrategy.is_downcross(EMAdif[length-2], EMAdif[length-1]) || emamacdemaStrategy.LessSellThreshold) { //up cross

		// do buy when cross up
		if emamacdemaStrategy.is_upcross(EMAdif[length-2], EMAdif[length-1]) || emamacdemaStrategy.LessBuyThreshold {
			if Option["enable_trading"] == "1" && emamacdemaStrategy.PrevEMATrade != "buy" {

				emamacdemaStrategy.PrevEMACross = "up"

				if emamacdemaStrategy.checkThreshold("buy", EMAdif[length-1]) {
					emamacdemaStrategy.PrevEMATrade = "buy"
					diff := fmt.Sprintf("%0.03f", EMAdif[length-1])
					warning := "EMA up cross, 买入buy In<----市价" + getTradePrice("", Price[length-1]) +
						",委托价" + getTradePrice("buy", Price[length-1]) + ",diff" + diff
					logger.Infoln(warning)
					if buy(getTradePrice("buy", Price[length-1]), tradeAmount) != "0" {
						emamacdemaStrategy.PrevBuyPirce = Price[length-1]
						warning += "[委托成功]"
					} else {
						warning += "[委托失败]"
					}

					go email.TriggerTrender(warning)
				}
			}
		}

		// do sell when cross down
		if emamacdemaStrategy.is_downcross(EMAdif[length-2], EMAdif[length-1]) || emamacdemaStrategy.LessSellThreshold {
			emamacdemaStrategy.PrevEMACross = "down"
			if Option["enable_trading"] == "1" && emamacdemaStrategy.PrevEMATrade != "sell" {
				if emamacdemaStrategy.checkThreshold("sell", EMAdif[length-1]) {
					emamacdemaStrategy.PrevEMATrade = "sell"
					diff := fmt.Sprintf("%0.03f", EMAdif[length-1])
					warning := "EMA down cross, 卖出Sell Out---->市价" + getTradePrice("", Price[length-1]) +
						",委托价" + getTradePrice("sell", Price[length-1]) + ",diff" + diff
					logger.Infoln(warning)

					var ematradeAmount string
					if emamacdemaStrategy.PrevMACDTrade == "sell" {
						ematradeAmount = MacdTradeAmount
						emamacdemaStrategy.PrevMACDTrade = "init"
					} else {
						ematradeAmount = tradeAmount
					}

					if sell(getTradePrice("sell", Price[length-1]), ematradeAmount) != "0" {
						warning += "[委托成功]"
					} else {
						warning += "[委托失败]"
					}

					go email.TriggerTrender(warning)
				}
			}
		}

		// backup the kline data for analyze
		if Config["env"] == "dev" {
			backup(Time[length-1])
		}
	}

	// macd cross
	if MACDdif[length-1] > 0 {
		if (MACDHistogram[length-2] < -0.000001 && MACDHistogram[length-1] > MACDbuyThreshold) ||
			(emamacdemaStrategy.PrevMACDTrade == "sell" && MACDHistogram[length-2] > 0.000001 && MACDHistogram[length-1] > MACDbuyThreshold) {
			if Option["enable_trading"] == "1" && emamacdemaStrategy.PrevMACDTrade == "sell" {
				emamacdemaStrategy.PrevMACDTrade = "buy"
				histogram := fmt.Sprintf("%0.03f", MACDHistogram[length-1])
				warning := "MACD up cross, 买入buy In<----市价" + getTradePrice("", Price[length-1]) +
					",委托价" + getTradePrice("buy", Price[length-1]) + ",histogram" + histogram
				logger.Infoln(warning)
				if buy(getTradePrice("buy", Price[length-1]), MacdTradeAmount) != "0" {
					emamacdemaStrategy.PrevBuyPirce = Price[length-1]
					warning += "[委托成功]"
				} else {
					warning += "[委托失败]"
				}

				go email.TriggerTrender(warning)
			}
		} else if (Price[length-1] < emaLong[length-1]) &&
			((MACDHistogram[length-2] > 0.000001 && MACDHistogram[length-1] < MACDsellThreshold) ||
				(emamacdemaStrategy.PrevMACDTrade == "buy" && MACDHistogram[length-2] < -0.000001 && MACDHistogram[length-1] < MACDsellThreshold)) {
			if Option["enable_trading"] == "1" && emamacdemaStrategy.PrevMACDTrade != "sell" {
				emamacdemaStrategy.PrevMACDTrade = "sell"
				histogram := fmt.Sprintf("%0.03f", MACDHistogram[length-1])
				warning := "MACD down cross, 卖出Sell Out---->市价" + getTradePrice("", Price[length-1]) +
					",委托价" + getTradePrice("sell", Price[length-1]) + ",histogram" + histogram
				logger.Infoln(warning)
				if sell(getTradePrice("sell", Price[length-1]), MacdTradeAmount) != "0" {
					warning += "[委托成功]"
				} else {
					warning += "[委托失败]"
				}

				go email.TriggerTrender(warning)
			}
		}
	}

	// do sell when price is below stoploss point
	if Price[length-1] < emamacdemaStrategy.PrevBuyPirce*(1-stoploss*0.01) {
		if Option["enable_trading"] == "1" && emamacdemaStrategy.PrevEMATrade != "sell" {

			warning := "stop loss, 卖出Sell Out---->市价" + getTradePrice("", Price[length-1]) + ",委托价" + getTradePrice("sell", Price[length-1])
			logger.Infoln(warning)

			var ematradeAmount string
			if emamacdemaStrategy.PrevMACDTrade == "sell" {
				ematradeAmount = MacdTradeAmount
			} else {
				ematradeAmount = tradeAmount
			}

			if sell(getTradePrice("sell", Price[length-1]), ematradeAmount) != "0" {
				warning += "[委托成功]"
			} else {
				warning += "[委托失败]"
			}

			go email.TriggerTrender(warning)

			emamacdemaStrategy.PrevEMATrade = "sell"
			emamacdemaStrategy.PrevMACDTrade = "init"
			emamacdemaStrategy.PrevBuyPirce = 0
		}
	}

	return true
}
