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
	. "common"
	. "config"
	"email"
	"fmt"
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
func (macdStrategy *MACDStrategy) Perform(tradeAPI TradeAPI, records []Record) bool {
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

	//compute the indictor
	emaShort := EMA(Price, shortEMA)
	emaLong := EMA(Price, longEMA)
	MACDdif := getMACDdif(emaShort, emaLong)
	MACDSignal := getMACDSignal(MACDdif, signalPeriod)
	MACDHistogram := getMACDHistogram(MACDdif, MACDSignal)

	length := len(Price)

	if MACDdif[length-1] != macdStrategy.PrevMACDdif {
		macdStrategy.PrevMACDdif = MACDdif[length-1]
		logger.Infof("MACD:d%5.03f\ts%5.03f\tph%5.03f\th%5.03f\tPrice:%5.02f\n", MACDdif[length-1], MACDSignal[length-1], MACDHistogram[length-2], MACDHistogram[length-1], Price[length-1])
	}

	//macd cross
	if (MACDHistogram[length-2] < -0.000001 && MACDHistogram[length-1] > MACDbuyThreshold) ||
		(macdStrategy.PrevMACDTrade == "sell" && MACDHistogram[length-2] > 0.000001 && MACDHistogram[length-1] > MACDbuyThreshold) {
		if Option["disable_trading"] != "1" && macdStrategy.PrevMACDTrade != "buy" {
			macdStrategy.PrevMACDTrade = "buy"

			histogram := fmt.Sprintf("%0.03f", MACDHistogram[length-1])
			warning := "MACD up cross, 买入buy In<----市价" + tradeAPI.GetTradePrice("", Price[length-1]) +
				",委托价" + tradeAPI.GetTradePrice("buy", Price[length-1]) + ",histogram" + histogram
			logger.Infoln(warning)
			if tradeAPI.Buy(tradeAPI.GetTradePrice("buy", Price[length-1]), tradeAmount) {
				macdStrategy.PrevBuyPirce = Price[length-1]
				warning += "[委托成功]"
			} else {
				warning += "[委托失败]"
			}

			go email.TriggerTrender(warning)
		}
	} else if (MACDHistogram[length-2] > 0.000001 && MACDHistogram[length-1] < MACDsellThreshold) ||
		(macdStrategy.PrevMACDTrade == "buy" && MACDHistogram[length-2] < -0.000001 && MACDHistogram[length-1] < MACDsellThreshold) {
		if Option["disable_trading"] != "1" && macdStrategy.PrevMACDTrade != "sell" {
			macdStrategy.PrevMACDTrade = "sell"

			histogram := fmt.Sprintf("%0.03f", MACDHistogram[length-1])
			warning := "MACD down cross, 卖出Sell Out---->市价" + tradeAPI.GetTradePrice("", Price[length-1]) +
				",委托价" + tradeAPI.GetTradePrice("sell", Price[length-1]) + ",histogram" + histogram
			logger.Infoln(warning)
			if tradeAPI.Sell(tradeAPI.GetTradePrice("sell", Price[length-1]), tradeAmount) {
				warning += "[委托成功]"
			} else {
				warning += "[委托失败]"
			}

			go email.TriggerTrender(warning)
		}
	}

	//do sell when price is below stoploss point
	if Price[length-1] < macdStrategy.PrevBuyPirce*(1-stoploss*0.01) {
		if Option["disable_trading"] != "1" && macdStrategy.PrevMACDTrade != "sell" {

			warning := "stop loss, 卖出Sell Out---->市价" + tradeAPI.GetTradePrice("", Price[length-1]) + ",委托价" + tradeAPI.GetTradePrice("sell", Price[length-1])
			logger.Infoln(warning)
			if tradeAPI.Sell(tradeAPI.GetTradePrice("sell", Price[length-1]), tradeAmount) {
				warning += "[委托成功]"
			} else {
				warning += "[委托失败]"
			}

			go email.TriggerTrender(warning)

			macdStrategy.PrevMACDTrade = "sell"
			macdStrategy.PrevBuyPirce = 0
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
