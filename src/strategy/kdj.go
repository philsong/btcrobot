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
	"logger"
)

type KDJStrategy struct {
	PrevKDJTrade string
	PrevTime     string
	PrevPrice    float64
}

func init() {
	kdjStrategy := new(KDJStrategy)
	kdjStrategy.PrevKDJTrade = "init"
	Register("KDJ", kdjStrategy)
}

//xxx strategy
func (kdjStrategy *KDJStrategy) Tick(records []Record) bool {
	//实现自己的策略

	tradeAmount := Option["tradeAmount"]

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

	length := len(records)

	if kdjStrategy.PrevTime == records[length-1].TimeStr &&
		kdjStrategy.PrevPrice == records[length-1].Close {
		return false
	}

	//K线为白，D线为黄，J线为红,K in middle
	k, d, j := getKDJ(records)

	if kdjStrategy.PrevTime != records[length-1].TimeStr ||
		kdjStrategy.PrevPrice != records[length-1].Close {
		kdjStrategy.PrevTime = records[length-1].TimeStr
		kdjStrategy.PrevPrice = records[length-1].Close

		logger.Infoln(records[length-1].TimeStr, records[length-1].Close)
		logger.Infof("d(黄线）%0.0f\tk(白线）%0.0f\tj(红线）%0.0f\n", d[length-2], k[length-2], j[length-2])
		logger.Infof("d(黄线）%0.0f\tk(白线）%0.0f\tj(红线）%0.0f\n", d[length-1], k[length-1], j[length-1])
	}

	if ((j[length-2] < k[length-2] && k[length-2] < d[length-2]) || kdjStrategy.PrevKDJTrade == "sell") &&
		(j[length-1] > k[length-1] && k[length-1] > d[length-1]) {
		logger.Infoln("KDJ up cross")
		if (kdjStrategy.PrevKDJTrade == "init" && d[length-2] <= 30) || kdjStrategy.PrevKDJTrade == "sell" {
			//do buy
			warning := "KDJ up cross, 买入buy In<----市价" + getTradePrice("", Price[length-1]) +
				",委托价" + getTradePrice("buy", Price[length-1])
			logger.Infoln(warning)
			if Buy(getTradePrice("buy", Price[length-1]), tradeAmount) != "0" {
				warning += "[委托成功]"
			} else {
				warning += "[委托失败]"
			}

			kdjStrategy.PrevKDJTrade = "buy"

			go email.TriggerTrender(warning)
		}

	}

	if ((j[length-2] > k[length-2] && k[length-2] > d[length-2]) || kdjStrategy.PrevKDJTrade == "buy") &&
		(j[length-1] < k[length-1] && k[length-1] < d[length-1]) {

		logger.Infoln("KDJ down cross")
		if (kdjStrategy.PrevKDJTrade == "init" && d[length-2] >= 70) || kdjStrategy.PrevKDJTrade == "buy" {
			//do sell
			warning := "KDJ down cross, 卖出Sell Out---->市价" + getTradePrice("", Price[length-1]) +
				",委托价" + getTradePrice("sell", Price[length-1])
			logger.Infoln(warning)
			if Sell(getTradePrice("sell", Price[length-1]), tradeAmount) != "0" {
				warning += "[委托成功]"
			} else {
				warning += "[委托失败]"
			}

			kdjStrategy.PrevKDJTrade = "sell"

			go email.TriggerTrender(warning)
		}

	}

	return true
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
