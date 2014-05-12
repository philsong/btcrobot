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

type MACDStrategy struct {
	PrevEMACross      string
	PrevMACDdif       float64
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
func (macdStrategy *MACDStrategy) Tick(records []Record) bool {
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
		(PrevTrade == "sell" && MACDHistogram[length-2] > 0.000001 && MACDHistogram[length-1] > MACDbuyThreshold) {
		if Option["enable_trading"] == "1" && PrevTrade != "buy" {
			PrevTrade = "buy"

			histogram := fmt.Sprintf("%0.03f", MACDHistogram[length-1])
			warning := "MACD up cross, 买入buy In<----市价" + getTradePrice("", Price[length-1]) +
				",委托价" + getTradePrice("buy", Price[length-1]) + ",histogram" + histogram
			logger.Infoln(warning)
			if Buy(getTradePrice("buy", Price[length-1]), tradeAmount) != "0" {
				PrevBuyPirce = Price[length-1]
				warning += "[委托成功]"
			} else {
				warning += "[委托失败]"
			}

			go email.TriggerTrender(warning)
		}
	} else if (MACDHistogram[length-2] > 0.000001 && MACDHistogram[length-1] < MACDsellThreshold) ||
		(PrevTrade == "buy" && MACDHistogram[length-2] < -0.000001 && MACDHistogram[length-1] < MACDsellThreshold) {
		if Option["enable_trading"] == "1" && PrevTrade != "sell" {
			PrevTrade = "sell"

			histogram := fmt.Sprintf("%0.03f", MACDHistogram[length-1])
			warning := "MACD down cross, 卖出Sell Out---->市价" + getTradePrice("", Price[length-1]) +
				",委托价" + getTradePrice("sell", Price[length-1]) + ",histogram" + histogram
			logger.Infoln(warning)
			if Sell(getTradePrice("sell", Price[length-1]), tradeAmount) != "0" {
				warning += "[委托成功]"
				PrevBuyPirce = 0
			} else {
				warning += "[委托失败]"
			}

			go email.TriggerTrender(warning)
		}
	}

	//do sell when price is below stoploss point
	processStoploss(Price)

	return true
}
