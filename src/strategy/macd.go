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
	"logger"
	"strconv"
)

type MACDStrategy struct {
	PrevMACDTrade     string
	PrevEMACross      string
	PrevMACDdif       float64
	PrevBuyPirce      float64
	LessBuyThreshold  bool
	LessSellThreshold bool
}

func init() {
	macdStrategy := new(MACDStrategy)
	macdStrategy.PrevEMACross = "unknown"
	Register("MACD", macdStrategy)
}

func getMACDHistogram(MACDdif, MACDSignal []float64) []float64 {
	var MACDHistogram []float64
	length := len(MACDSignal)
	for i := 0; i < length; i++ {
		MACDHistogramAt := getMACDHistogramAt(MACDdif, MACDSignal, i)
		MACDHistogram = append(MACDHistogram, MACDHistogramAt)
	}
	return MACDHistogram
}

//MACD strategy
func (macdStrategy *MACDStrategy) Perform(tradeAPI TradeAPI, Time []string, Price []float64, Volumn []float64) bool {
	//read config
	shortEMA, _ := strconv.Atoi(Option["shortEMA"])
	longEMA, _ := strconv.Atoi(Option["longEMA"])

	signalPeriod, _ := strconv.Atoi(Option["signalPeriod"])
	/*
		MACDMinThreshold, err := strconv.ParseFloat(Option["MACDMinThreshold"], 64)
		if err != nil {
			logger.Debugln("config item MACDMinThreshold is not float")
			return false
		}
	*/

	tradeAmount := Option["tradeAmount"]
	stoploss, err := strconv.ParseFloat(Option["stoploss"], 64)
	if err != nil {
		logger.Errorln("config item stoploss is not float")
		return false
	}

	//compute the indictor
	emaShort := EMA(Price, shortEMA)
	emaLong := EMA(Price, longEMA)
	MACDdif := getMACDdif(emaShort, emaLong)
	MACDSignal := getMACDSignal(MACDdif, signalPeriod)
	MACDHistogram := getMACDHistogram(MACDdif, MACDSignal)

	length := len(Price)

	if MACDdif[length-1] != macdStrategy.PrevMACDdif {
		macdStrategy.PrevMACDdif = MACDdif[length-1]
		logger.Infof("MACD:d%5.03f\ts%5.03f\th%5.03f\tPrice:%5.02f\n", MACDdif[length-1], MACDSignal[length-1], MACDHistogram[length-1], Price[length-1])
	}

	//macd cross
	for i := 1; i < length; i++ {
		if MACDHistogram[i-1] < 0 && MACDHistogram[i] > 0 {
			if i == length-1 && macdStrategy.PrevMACDTrade != "buy" {
				macdStrategy.PrevMACDTrade = "buy"
				warning := "MACD up cross, 买入buy In<----市价" + tradeAPI.GetTradePrice("") + ",委托价" + tradeAPI.GetTradePrice("buy")
				logger.Infoln(warning)
				if tradeAPI.Buy(tradeAPI.GetTradePrice("buy"), tradeAmount) {
					macdStrategy.PrevBuyPirce = Price[length-1]
					warning += "[委托成功]"
				} else {
					warning += "[委托失败]"
				}

				go email.TriggerTrender(warning)
			}
		} else if MACDHistogram[i-1] > 0 && MACDHistogram[i] < 0 {

			if i == length-1 && macdStrategy.PrevMACDTrade != "sell" {
				macdStrategy.PrevMACDTrade = "sell"
				warning := "MACD down cross, 卖出Sell Out---->市价" + tradeAPI.GetTradePrice("") + ",委托价" + tradeAPI.GetTradePrice("sell")
				logger.Infoln(warning)
				if tradeAPI.Sell(tradeAPI.GetTradePrice("sell"), tradeAmount) {
					warning += "[委托成功]"
				} else {
					warning += "[委托失败]"
				}

				go email.TriggerTrender(warning)
			}
		}
	}

	//do sell when price is below stoploss point
	if Price[length-1] < macdStrategy.PrevBuyPirce*(1-stoploss*0.01) {
		if Option["disable_trading"] != "1" && macdStrategy.PrevMACDTrade != "sell" {
			macdStrategy.PrevMACDTrade = "sell"
			warning := "stop loss, 卖出Sell Out---->市价" + tradeAPI.GetTradePrice("") + ",委托价" + tradeAPI.GetTradePrice("sell")
			logger.Infoln(warning)
			if tradeAPI.Sell(tradeAPI.GetTradePrice("sell"), tradeAmount) {
				warning += "[委托成功]"
			} else {
				warning += "[委托失败]"
			}

			go email.TriggerTrender(warning)
		}
	}

	return true
}

func getMACDdifAt(emaShort, emaLong []float64, idx int) float64 {
	var ces = emaShort[idx]
	var cel = emaLong[idx]
	if cel == 0 {
		return 0
	} else {
		return (ces - cel)
	}
}

func getMACDdif(emaShort, emaLong []float64) []float64 {
	// loop through data
	var MACDdif []float64
	length := len(emaShort)
	for i := 0; i < length; i++ {
		MACDdifAt := getMACDdifAt(emaShort, emaLong, i)
		MACDdif = append(MACDdif, MACDdifAt)
	}

	return MACDdif
}

func getMACDSignal(MACDdif []float64, signalPeriod int) []float64 {
	signal := EMA(MACDdif, signalPeriod)
	return signal
}

func getMACDHistogramAt(MACDdif, MACDSignal []float64, idx int) float64 {
	var dif = MACDdif[idx]
	var signal = MACDSignal[idx]
	if signal == 0 {
		return 0
	} else {
		return dif - signal
	}
}
