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
func (emamacdStrategy *EMAMACDStrategy) Perform(tradeAPI TradeAPI, Time []string, Price []float64, Volumn []float64) bool {
	//read config
	shortEMA, _ := strconv.Atoi(Option["shortEMA"])
	longEMA, _ := strconv.Atoi(Option["longEMA"])
	signalPeriod, _ := strconv.Atoi(Option["signalPeriod"])

	nTradeAmount, err := strconv.ParseFloat(Option["tradeAmount"], 64)
	if err != nil {
		logger.Debugln("config item tradeAmount is not float")
		return false
	}

	MacdTradeAmount := fmt.Sprintf("%0.02f", 0.5*nTradeAmount)

	tradeAmount := Option["tradeAmount"]

	stoploss, err := strconv.ParseFloat(Option["stoploss"], 64)
	if err != nil {
		logger.Errorln("config item stoploss is not float")
		return false
	}

	//compute the indictor
	emaShort := EMA(Price, shortEMA)
	emaLong := EMA(Price, longEMA)
	EMAdif := getEMAdif(emaShort, emaLong)

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

	if MACDdif[length-1] != emamacdStrategy.PrevMACDdif {
		emamacdStrategy.PrevMACDdif = MACDdif[length-1]
		logger.Infof("MACD:d%5.03f\ts%5.03f\th%5.03f\tPrice:%5.02f\n", MACDdif[length-1], MACDSignal[length-1], MACDHistogram[length-1], Price[length-1])
	}

	//go TriggerPrice(Price[length-1])
	if EMAdif[length-1] != emamacdStrategy.PrevEMAdif {
		emamacdStrategy.PrevEMAdif = EMAdif[length-1]
		logger.Infof("EMA Diff:%0.03f\t%0.03f\tPrice:%0.02f\n", EMAdif[length-2], EMAdif[length-1], Price[length-1])
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
	if (emamacdStrategy.is_upcross(EMAdif[length-2], EMAdif[length-1]) || emamacdStrategy.LessBuyThreshold) ||
		(emamacdStrategy.is_downcross(EMAdif[length-2], EMAdif[length-1]) || emamacdStrategy.LessSellThreshold) { //up cross

		//do buy when cross up
		if emamacdStrategy.is_upcross(EMAdif[length-2], EMAdif[length-1]) || emamacdStrategy.LessBuyThreshold {
			if Option["disable_trading"] != "1" && emamacdStrategy.PrevEMATrade != "buy" {

				emamacdStrategy.PrevEMACross = "up"

				if emamacdStrategy.checkThreshold("buy", EMAdif[length-1]) {
					emamacdStrategy.PrevEMATrade = "buy"
					warning := "EMA up cross, 买入buy In<----市价" + tradeAPI.GetTradePrice("") + ",委托价" + tradeAPI.GetTradePrice("buy")
					logger.Infoln(warning)
					if tradeAPI.Buy(tradeAPI.GetTradePrice("buy"), tradeAmount) {
						emamacdStrategy.PrevBuyPirce = Price[length-1]
						warning += "[委托成功]"
					} else {
						warning += "[委托失败]"
					}

					go email.TriggerTrender(warning)
				}
			}
		}

		//do sell when cross down
		if emamacdStrategy.is_downcross(EMAdif[length-2], EMAdif[length-1]) || emamacdStrategy.LessSellThreshold {
			emamacdStrategy.PrevEMACross = "down"
			if Option["disable_trading"] != "1" && emamacdStrategy.PrevEMATrade != "sell" {
				if emamacdStrategy.checkThreshold("sell", EMAdif[length-1]) {
					emamacdStrategy.PrevEMATrade = "sell"
					warning := "EMA down cross, 卖出Sell Out---->市价" + tradeAPI.GetTradePrice("") + ",委托价" + tradeAPI.GetTradePrice("sell")
					logger.Infoln(warning)

					var ematradeAmount string
					if emamacdStrategy.PrevMACDTrade == "sell" {
						ematradeAmount = MacdTradeAmount
						emamacdStrategy.PrevMACDTrade = "init"
					} else {
						ematradeAmount = tradeAmount
					}

					if tradeAPI.Sell(tradeAPI.GetTradePrice("sell"), ematradeAmount) {
						warning += "[委托成功]"
					} else {
						warning += "[委托失败]"
					}

					go email.TriggerTrender(warning)
				}
			}
		}

		//macd cross
		if emamacdStrategy.PrevEMATrade == "buy" {
			if MACDHistogram[length-2] < 0 && MACDHistogram[length-1] > 0 {
				if Option["disable_trading"] != "1" && emamacdStrategy.PrevMACDTrade == "sell" {
					emamacdStrategy.PrevMACDTrade = "buy"
					warning := "MACD up cross, 买入buy In<----市价" + tradeAPI.GetTradePrice("") + ",委托价" + tradeAPI.GetTradePrice("buy")
					logger.Infoln(warning)
					if tradeAPI.Buy(tradeAPI.GetTradePrice("buy"), MacdTradeAmount) {
						emamacdStrategy.PrevBuyPirce = Price[length-1]
						warning += "[委托成功]"
					} else {
						warning += "[委托失败]"
					}

					go email.TriggerTrender(warning)
				}
			} else if MACDHistogram[length-2] > 0 && MACDHistogram[length-1] < 0 {
				if Option["disable_trading"] != "1" && emamacdStrategy.PrevMACDTrade != "sell" {
					emamacdStrategy.PrevMACDTrade = "sell"
					warning := "MACD down cross, 卖出Sell Out---->市价" + tradeAPI.GetTradePrice("") + ",委托价" + tradeAPI.GetTradePrice("sell")
					logger.Infoln(warning)
					if tradeAPI.Sell(tradeAPI.GetTradePrice("sell"), MacdTradeAmount) {
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

	//do sell when price is below stoploss point
	if Price[length-1] < emamacdStrategy.PrevBuyPirce*(1-stoploss*0.01) {
		if Option["disable_trading"] != "1" && emamacdStrategy.PrevEMATrade != "sell" {
			emamacdStrategy.PrevEMATrade = "sell"
			warning := "stop loss, 卖出Sell Out---->市价" + tradeAPI.GetTradePrice("") + ",委托价" + tradeAPI.GetTradePrice("sell")
			logger.Infoln(warning)

			var ematradeAmount string
			if emamacdStrategy.PrevMACDTrade == "sell" {
				ematradeAmount = MacdTradeAmount
			} else {
				ematradeAmount = tradeAmount
			}

			if tradeAPI.Sell(tradeAPI.GetTradePrice("sell"), ematradeAmount) {
				warning += "[委托成功]"
			} else {
				warning += "[委托失败]"
			}

			go email.TriggerTrender(warning)
		}
	}

	return true
}
