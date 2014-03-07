// Copyright 2014 The btcbot Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://github.com/philsong/
// Author：Phil	78623269@qq.com

package huobiapi

import (
//"logger"
)

/*
wisdom The formula is

RSV = (CLOSE-LLV(LOW,9))/(HHV(HIGH,9)-LLV(LOW,9))*100;
K = SMA(RSV,3,1);
D = SMA(K,3,1);
J = 3*K-2*D;

LLV means lowest value in period
HHV means highest value in period
SMA is Simple Moving Average
*/
func doKDJ(records []PeroidRecord) ([]float64, []float64, []float64) {
	periods := 9
	k, d := kd(records, periods)
	j := j(k, d)
	return k, d, j
}

func j(k, d []float64) []float64 {
	length := len(k)
	var j []float64 = make([]float64, length)

	// Loop through the entire array.
	for i := 0; i < length; i++ {
		j[i] = 3*k[i] - 2*d[i]

	}

	return j
}

func kd(records []PeroidRecord, periods int) ([]float64, []float64) {
	var periodLowArr, periodHighArr []float64
	length := len(records)
	var rsv []float64 = make([]float64, length)
	var k []float64 = make([]float64, length)
	var d []float64 = make([]float64, length)

	// Loop through the entire array.
	for i := 0; i < length; i++ {
		// add points to the array.
		periodLowArr = append(periodLowArr, records[i].Low)
		periodHighArr = append(periodHighArr, records[i].High)

		// 1: Check if array is "filled" else create null point in line.
		// 2: Calculate average.
		// 3: Remove first value.
		if periods == len(periodLowArr) {
			lowest := arrayLowest(periodLowArr)
			highest := arrayHighest(periodHighArr)
			//logger.Infoln(i, records[i].Close, lowest, highest)
			rsv[i] = 100 * (records[i].Close - lowest) / (highest - lowest)
			k[i] = (2.0/3)*k[i-1] + 1.0/3*rsv[i]
			d[i] = (2.0/3)*d[i-1] + 1.0/3*k[i]
			// remove first value in array.
			periodLowArr = periodLowArr[1:]
			periodHighArr = periodHighArr[1:]
		} else {
			k[i] = 50
			d[i] = 50
			rsv[i] = 0
		}
	}

	return k, d
}

func Highest(yData []float64, periods int) []float64 {
	var periodArr []float64
	length := len(yData)
	var HighestLine []float64 = make([]float64, length)

	// Loop through the entire array.
	for i := 0; i < length; i++ {
		// add points to the array.
		periodArr = append(periodArr, yData[i])
		// 1: Check if array is "filled" else create null point in line.
		// 2: Calculate average.
		// 3: Remove first value.
		if periods == len(periodArr) {
			HighestLine[i] = arrayHighest(periodArr)

			// remove first value in array.
			periodArr = periodArr[1:]
		} else {
			HighestLine[i] = 0
		}
	}

	return HighestLine
}

func Lowest(yData []float64, periods int) []float64 {
	var periodArr []float64
	length := len(yData)
	var LowestLine []float64 = make([]float64, length)

	// Loop through the entire array.
	for i := 0; i < length; i++ {
		// add points to the array.
		periodArr = append(periodArr, yData[i])
		// 1: Check if array is "filled" else create null point in line.
		// 2: Calculate average.
		// 3: Remove first value.
		if periods == len(periodArr) {
			LowestLine[i] = arrayLowest(periodArr)

			// remove first value in array.
			periodArr = periodArr[1:]
		} else {
			LowestLine[i] = 0
		}
	}

	return LowestLine
}

/* Function based on the idea of a simple moving average.
 * @param yData : array of y variables.
 * @param periods : The amount of "days" to average from.
 * @return an array containing the SMA.
**/
func SMA(yData []float64, periods int) []float64 {
	var periodArr []float64
	length := len(yData)
	var smLine []float64 = make([]float64, length)

	// Loop through the entire array.
	for i := 0; i < length; i++ {
		// add points to the array.
		periodArr = append(periodArr, yData[i])

		// 1: Check if array is "filled" else create null point in line.
		// 2: Calculate average.
		// 3: Remove first value.
		if periods == len(periodArr) {
			smLine[i] = arrayAvg(periodArr)

			// remove first value in array.
			periodArr = periodArr[1:]
		} else {
			smLine[i] = 0
		}
	}

	return smLine
}

func arrayLowest(yData []float64) float64 {
	length := len(yData)
	var lowest = yData[0]

	// Loop through the entire array.
	for i := 1; i < length; i++ {
		if yData[i] < lowest {
			lowest = yData[i]
		}
	}

	return lowest
}

func arrayHighest(yData []float64) float64 {
	length := len(yData)
	var highest = yData[0]

	// Loop through the entire array.
	for i := 1; i < length; i++ {
		if yData[i] > highest {
			highest = yData[i]
		}
	}

	return highest
}
