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
	"fmt"
	"logger"
	"os"
	"strconv"
)

type TradeAPI interface {
	BuyIn(price, amount string) bool
	SellOut(price, amount string) bool
	GetTradePrice(tradeDirection string) string
	GetPrevTrend() string
	SetPrevTrend(trend string)
}

//EMA strategy
func PerformEMA(tradeAPI TradeAPI, Time []string, Price []float64, Volumn []float64) {

	//
	if len(Time) == 0 || len(Price) == 0 || len(Volumn) == 0 {
		logger.Errorln("detect exception data")
		return
	}

	//read config
	shortEMA, _ := strconv.Atoi(Option["shortEMA"])
	longEMA, _ := strconv.Atoi(Option["longEMA"])

	_, err := strconv.ParseFloat(Option["tradeAmount"], 64)
	if err != nil {
		logger.Debugln("config item tradeAmount is not float")
		return
	}
	tradeAmount := Option["tradeAmount"]

	//compute the indictor
	emaShort := EMA(Price, shortEMA)
	emaLong := EMA(Price, longEMA)
	EMAdif := getEMAdif(emaShort, emaLong)

	length := len(Price)

	//EMA cross
	if (EMAdif[length-2] < 0 && EMAdif[length-1] > 0) || (EMAdif[length-2] > 0 && EMAdif[length-1] < 0) { //up cross

		//check exception data in trade center
		if checkException(Price[length-2], Price[length-1], Volumn[length-1]) == false {
			logger.Errorln("detect exception data of trade center", Price[length-2], Price[length-1], Volumn[length-1])
			return
		}

		//do buy when cross up
		if EMAdif[length-2] < 0 && EMAdif[length-1] > 0 {
			if Option["disable_trading"] != "1" && tradeAPI.GetPrevTrend() != "up" {
				tradeAPI.SetPrevTrend("up")
				logger.Infoln("EMA up cross, 买入buy In", tradeAPI.GetTradePrice(""))
				tradeAPI.BuyIn(tradeAPI.GetTradePrice("buy"), tradeAmount)
			}
		}

		//do sell when cross down
		if EMAdif[length-2] > 0 && EMAdif[length-1] < 0 {
			if Option["disable_trading"] != "1" && tradeAPI.GetPrevTrend() != "down" {
				tradeAPI.SetPrevTrend("down")
				logger.Infoln("EMA down cross, 卖出Sell Out", tradeAPI.GetTradePrice(""))
				tradeAPI.SellOut(tradeAPI.GetTradePrice("sell"), tradeAmount)
			}
		}

		//backup the kline data for analyze
		backup(Time[length-1])
	}
}

//for backup the kline file to detect the huobi bug
func backup(Time string) {

	oldFile := "cache/TradeKLine_minute.data"
	newFile := fmt.Sprintf("%s_%s", oldFile, Time)
	os.Rename(oldFile, newFile)
}

//check exception data in trade center
func checkException(yPrevPrice, Price, Volumn float64) bool {
	if Price > yPrevPrice+10 && Volumn < 1 {
		return false
	}

	if Price < yPrevPrice-10 && Volumn < 1 {
		return false
	}

	return true
}

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
