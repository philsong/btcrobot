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
)

func getEMAdifAt(emaShort, emaLong []float64, idx int) float64 {
	var cel = emaLong[idx]
	var ces = emaShort[idx]
	if cel == 0 {
		return 0
	} else {
		return 100 * (ces - cel) / ((ces + cel) / 2)
	}
}

func getEMAdif(emaShort, emaLong []float64) []float64 {
	// loop through data
	var EMAdifs []float64
	length := len(emaShort)
	for i := 0; i < length; i++ {
		EMAdifAt := getEMAdifAt(emaShort, emaLong, i)
		EMAdifs = append(EMAdifs, EMAdifAt)
	}

	return EMAdifs
}

/* Function based on the idea of an exponential moving average.
 *
 * Formula: EMA = Price(t) * k + EMA(y) * (1 - k)
 * t = today y = yesterday N = number of days in EMA k = 2/(2N+1)
 *
 * @param Price : array of y variables.
 * @param periods : The amount of "days" to average from.
 * @return an array containing the EMA.
**/
func EMA(Price []float64, periods int) []float64 {
	var t float64
	y := 0.0
	n := float64(periods)
	var k float64
	k = 2 / (n + 1)
	var ema float64 // exponential moving average.

	var periodArr []float64
	var startpos int
	length := len(Price)
	var emaLine []float64 = make([]float64, length)

	// loop through data
	for i := 0; i < length; i++ {
		if Price[i] != 0 {
			startpos = i + 1
			break
		} else {
			emaLine[i] = 0
		}
	}

	for i := startpos; i < length; i++ {
		periodArr = append(periodArr, Price[i])

		// 0: runs if the periodArr has enough points.
		// 1: set currentvalue (today).
		// 2: set last value. either by past avg or yesterdays ema.
		// 3: calculate todays ema.
		if periods == len(periodArr) {

			t = Price[i]

			if y == 0 {
				y = arrayAvg(periodArr)
			} else {
				ema = (t * k) + (y * (1 - k))
				y = ema
			}

			emaLine[i] = y

			// remove first value in array.
			periodArr = periodArr[1:]

		} else {

			emaLine[i] = 0
		}

	}

	return emaLine
}

/* Function that returns average of an array's values.
 *
**/
func arrayAvg(arr []float64) float64 {
	sum := 0.0

	for i := 0; i < len(arr); i++ {
		sum = sum + arr[i]
	}

	return (sum / (float64)(len(arr)))
}

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
func getKDJ(records []Record) ([]float64, []float64, []float64) {
	period := 9
	k, d := kd(records, period)
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

func kd(records []Record, period int) ([]float64, []float64) {
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
		if period == len(periodLowArr) {
			lowest := arrayLowest(periodLowArr)
			highest := arrayHighest(periodHighArr)
			//logger.Infoln(i, records[i].Close, lowest, highest)
			if highest-lowest < 0.000001 {
				rsv[i] = 100
			} else {
				rsv[i] = 100 * (records[i].Close - lowest) / (highest - lowest)
			}

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

func Highest(Price []float64, periods int) []float64 {
	var periodArr []float64
	length := len(Price)
	var HighestLine []float64 = make([]float64, length)

	// Loop through the entire array.
	for i := 0; i < length; i++ {
		// add points to the array.
		periodArr = append(periodArr, Price[i])
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

func Lowest(Price []float64, periods int) []float64 {
	var periodArr []float64
	length := len(Price)
	var LowestLine []float64 = make([]float64, length)

	// Loop through the entire array.
	for i := 0; i < length; i++ {
		// add points to the array.
		periodArr = append(periodArr, Price[i])
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
 * @param Price : array of y variables.
 * @param periods : The amount of "days" to average from.
 * @return an array containing the SMA.
**/
func SMA(Price []float64, periods int) []float64 {
	var periodArr []float64
	length := len(Price)
	var smLine []float64 = make([]float64, length)

	// Loop through the entire array.
	for i := 0; i < length; i++ {
		// add points to the array.
		periodArr = append(periodArr, Price[i])

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

func arrayLowest(Price []float64) float64 {
	length := len(Price)
	var lowest = Price[0]

	// Loop through the entire array.
	for i := 1; i < length; i++ {
		if Price[i] < lowest {
			lowest = Price[i]
		}
	}

	return lowest
}

func arrayHighest(Price []float64) float64 {
	length := len(Price)
	var highest = Price[0]

	// Loop through the entire array.
	for i := 1; i < length; i++ {
		if Price[i] > highest {
			highest = Price[i]
		}
	}

	return highest
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
