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

type MyMACDStrategy struct {
	PrevEMACross      string
	PrevMACDdif       float64
	LessBuyThreshold  bool
	LessSellThreshold bool
}

func init() {
	macdStrategy := new(MyMACDStrategy)
	macdStrategy.PrevEMACross = "unknown"
	Register("MYMACD", macdStrategy)
}

type MACD struct {
	EMAShort float64
	EMALong  float64
	DIF      float64
	DEA      float64
	BAR      float64
}

func CalcMACD(price []float64, periodShort int, periodLong int, periodDIF int) (macd []MACD) {
	macd = make([]MACD, len(price))

	macd[0].EMAShort = price[0]
	macd[0].EMALong = price[0]
	macd[0].DIF = 0
	macd[0].DEA = 0
	macd[0].BAR = 0

	for i := 1; i < len(price); i++ {
		macd[i].EMAShort = macd[i-1].EMAShort*(float64(periodShort-1))/(float64(periodShort+1)) +
			price[i]*2.0/(float64(periodShort+1))
		macd[i].EMALong = macd[i-1].EMALong*(float64(periodLong-1))/(float64(periodLong+1)) +
			price[i]*2.0/(float64(periodLong+1))
		macd[i].DIF = macd[i].EMAShort - macd[i].EMALong
		macd[i].DEA = macd[i-1].DEA*(float64(periodDIF-1))/(float64(periodDIF+1)) +
			macd[i].DIF*2.0/(float64(periodDIF+1))
		macd[i].BAR = macd[i].DIF - macd[i].DEA
	}

	return
}

func (macdStrategy *MyMACDStrategy) Tick(records []Record) bool {
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
	if MACDsellThreshold > 0 {
		MACDsellThreshold = 0 - MACDsellThreshold
	}

	shortEMA, _ := strconv.Atoi(Option["shortEMA"])
	longEMA, _ := strconv.Atoi(Option["longEMA"])
	signalPeriod, _ := strconv.Atoi(Option["signalPeriod"])

	var Price []float64
	for _, v := range records {
		Price = append(Price, v.Close)
	}

	length := len(records)
	macd := CalcMACD(Price, shortEMA, longEMA, signalPeriod)
	barLast := macd[length-2].BAR
	barCurr := macd[length-1].BAR

	// macd cross
	if (barLast < -0.000001 && barCurr > MACDbuyThreshold) ||
		(PrevTrade == "sell" && barLast > 0.000001 && barCurr > MACDbuyThreshold) {
		Buy()
	} else if (barLast > 0.000001 && barCurr < MACDsellThreshold) ||
		(PrevTrade == "buy" && barLast < -0.000001 && barCurr < MACDsellThreshold) {
		Sell()
	}

	// do sell when price is below stoploss point
	processStoploss(lastPrice)
	processTimeout()

	return true
}
