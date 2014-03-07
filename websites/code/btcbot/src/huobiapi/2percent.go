// Copyright 2014 The btcbot Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://github.com/philsong/
// Author：Phil	78623269@qq.com

package huobiapi

import (
	"logger"
)

func (w *Huobi) do2Percent(xData []string, yData []float64) {
	lowpoint := yData[0]
	highpoint := yData[0]
	lastTrade := "init"

	factor := 0.001
	logger.OverrideStart(w.Peroid)
	logger.Overrideln("start")
	for i := 0; i < len(yData); i++ {
		//logger.Overrideln("", i, lastTrade, xData[i], yData[i])
		if lastTrade != "buy" && yData[i]-lowpoint > 0 && yData[i]-lowpoint > factor*lowpoint {
			logger.Overrideln("++", i, xData[i], yData[i], lowpoint, highpoint, yData[i]-lowpoint, factor*lowpoint)

			highpoint = yData[i]
			lowpoint = yData[i]
			lastTrade = "buy"
		} else if lastTrade != "sell" && yData[i]-highpoint < 0 && yData[i]-highpoint < -factor*highpoint {
			logger.Overrideln("--", i, xData[i], yData[i], lowpoint, highpoint, yData[i]-highpoint, -factor*highpoint)

			highpoint = yData[i]
			lowpoint = yData[i]
			lastTrade = "sell"
		} else {
			if highpoint < yData[i] {
				highpoint = yData[i]
			} else if lowpoint > yData[i] {
				lowpoint = yData[i]
			}
		}
	}
	logger.Overrideln("end")

	return
}
