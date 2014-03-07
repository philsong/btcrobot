// Copyright 2014 The btcbot Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://github.com/philsong/
// Author：Phil	78623269@qq.com

package huobiapi

import (
	"logger"
)

func getMomentum(closeToday, closeNDaysAgo float64) float64 {
	momentum := closeToday - closeNDaysAgo
	return momentum
}

func getROC(closeToday, closeNDaysAgo float64) float64 {
	momentum := closeToday - closeNDaysAgo
	roc := momentum / closeNDaysAgo
	return roc
}

func getUD(yData []float64) ([]float64, []float64) {
	length := len(yData)
	var u []float64 = make([]float64, length)
	var d []float64 = make([]float64, length)

	// Loop through the entire array.
	for i := 1; i < length; i++ {
		diff := yData[i] - yData[i-1]
		if diff > 0 {
			u[i] = diff
			d[i] = 0
		} else if diff < 0 {
			d[i] = -diff
			u[i] = 0
		}
	}

	return u, d
}

func getRSI(yData []float64, periods int) []float64 {
	var periodArr []float64
	length := len(yData)
	var rsi []float64 = make([]float64, length)

	// Loop through the entire array.
	for i := 0; i < length; i++ {
		// add points to the array.
		periodArr = append(periodArr, yData[i])

		// 1: Check if array is "filled" else create null point in line.
		// 2: Calculate average.
		// 3: Remove first value.
		if periods == len(periodArr) {
			u, d := getUD(periodArr)
			rs := arrayAvg(u) / arrayAvg(d)

			rsi[i] = 100 - 100.0/(1+rs)

			// remove first value in array.
			periodArr = periodArr[1:]
		} else {
			rsi[i] = 0
		}
	}

	return rsi
}

func rsi(yData []float64) {
	rsShort := getRSI(yData, 6)
	rsLong := getRSI(yData, 12)

	length := len(rsShort)

	// Loop through the entire array.
	for i := 0; i < length; i++ {
		logger.Infoln(rsShort[i], rsLong[i])
	}
}
