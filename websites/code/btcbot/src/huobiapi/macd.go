// Copyright 2014 The btcbot Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://github.com/philsong/
// Author：Phil	78623269@qq.com

package huobiapi

import (
	. "config"
	"fmt"
	"logger"
	"service"
	"strconv"
)

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

func getMACDHistogram(MACDdif, MACDSignal []float64) []float64 {
	var MACDHistogram []float64
	length := len(MACDSignal)
	for i := 0; i < length; i++ {
		MACDHistogramAt := getMACDHistogramAt(MACDdif, MACDSignal, i)
		MACDHistogram = append(MACDHistogram, MACDHistogramAt)
	}
	return MACDHistogram
}

func (w *Huobi) doMACD(xData []string, yData []float64) {
	if len(yData) == 0 {
		return
	}
	//read config
	shortEMA, _ := strconv.Atoi(Config["shortEMA"])
	middleEMA, _ := strconv.Atoi(Config["middleEMA"])
	longEMA, _ := strconv.Atoi(Config["longEMA"])

	signalPeriod, _ := strconv.Atoi(Config["signalPeriod"])
	MACDMinThreshold, err := strconv.ParseFloat(Config["MACDMinThreshold"], 64)
	if err != nil {
		logger.Debugln("config item MACDMinThreshold is not float")
		return
	}

	tradeAmount := Config["MACDtradeAmount"]

	//compute the indictor
	emaShort := EMA(yData, shortEMA)
	emaMiddle := EMA(yData, middleEMA)
	emaLong := EMA(yData, longEMA)
	MACDdif := getMACDdif(emaShort, emaLong)
	MACDSignal := getMACDSignal(MACDdif, signalPeriod)
	MACDHistogram := getMACDHistogram(MACDdif, MACDSignal)

	length := len(yData)

	logger.OverrideStart(w.Peroid)

	logger.Overridef("macd cross分析[%d/%d]\n", shortEMA, longEMA)
	//macd cross
	for i := 1; i < length; i++ {
		if MACDHistogram[i-1] < 0 && MACDHistogram[i] > 0 {
			if MACDHistogram[i]-MACDHistogram[i-1] >= 0 && yData[i] > emaMiddle[i] {
				logger.Overrideln("++", i, xData[i], yData[i], fmt.Sprintf("%0.02f", MACDSignal[i]))
			} else {
				logger.Overrideln(" +", i, xData[i], yData[i], fmt.Sprintf("%0.02f", MACDSignal[i]))
			}
			if i == length-1 && w.latestMACDTrend != 1 && MACDSignal[i] < -MACDMinThreshold {
				w.latestMACDTrend = 1
				logger.Infoln("MACD has switched, 探测到买入点", w.getTradePrice("buy"))
				go service.TriggerTrender("MACD has switched, 探测到买入点")

				w.Do_buy(w.getTradePrice("buy"), tradeAmount)
			}
		} else if MACDHistogram[i-1] > 0 && MACDHistogram[i] < 0 {
			if MACDHistogram[i]-MACDHistogram[i-1] <= 0 && yData[i] < emaMiddle[i] {
				logger.Overrideln("--", i, xData[i], yData[i], fmt.Sprintf("%0.02f", MACDSignal[i]))
			} else {
				logger.Overrideln(" -", i, xData[i], yData[i], fmt.Sprintf("%0.02f", MACDSignal[i]))
			}
			if i == length-1 && w.latestMACDTrend != -1 && MACDSignal[i] > MACDMinThreshold {
				w.latestMACDTrend = -1
				logger.Infoln("MACD has switched, 探测到卖出点", w.getTradePrice("sell"))
				go service.TriggerTrender("MACD has switched, 探测到卖出点")

				w.Do_sell(w.getTradePrice("sell"), tradeAmount)
			}
		}
	}
}
