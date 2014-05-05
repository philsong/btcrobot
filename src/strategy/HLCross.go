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
	"strconv"
)

type HLCrossStrategy struct {
	PrevHLCrossTrade string
	PrevBuyPirce     float64
	PrevClosePrice   float64
	PrevHighPrice    float64
	PrevLowPrice     float64
}

func init() {
	HLCross := new(HLCrossStrategy)
	Register("HLCross", HLCross)
}

//HLCross strategy
func (HLCross *HLCrossStrategy) Tick(records []Record) bool {
	//read config

	tradeAmount := Option["tradeAmount"]
	stoploss, err := strconv.ParseFloat(Option["stoploss"], 64)
	if err != nil {
		logger.Errorln("config item stoploss is not float")
		return false
	}

	var Time []string
	var Price []float64
	var Volumn []float64
	for _, v := range records {
		Time = append(Time, v.TimeStr)
		Price = append(Price, v.Close)
		Volumn = append(Volumn, v.Volumn)
	}

	length := len(Price)

	if HLCross.PrevClosePrice != records[length-1].Close ||
		HLCross.PrevHighPrice != records[length-2].High ||
		HLCross.PrevLowPrice != records[length-2].Low {
		HLCross.PrevClosePrice = records[length-1].Close
		HLCross.PrevHighPrice = records[length-2].High
		HLCross.PrevLowPrice = records[length-2].Low

		logger.Infof("nowClose %0.02f prevHigh %0.02f prevLow %0.02f\n", records[length-1].Close, records[length-2].High, records[length-2].Low)
	}

	//HLCross cross
	if records[length-1].Close > records[length-2].High {
		if Option["enable_trading"] == "1" && HLCross.PrevHLCrossTrade != "buy" {
			warning := "HLCross up, 买入buy In<----市价" + getTradePrice("", Price[length-1]) +
				",委托价" + getTradePrice("buy", Price[length-1])
			logger.Infoln(warning)
			if Buy(getTradePrice("buy", Price[length-1]), tradeAmount) != "0" {
				HLCross.PrevBuyPirce = Price[length-1]
				warning += "[委托成功]"
				HLCross.PrevHLCrossTrade = "buy"
			} else {
				warning += "[委托失败]"
			}

			go email.TriggerTrender(warning)
		}
	} else if records[length-1].Close < records[length-2].Low {
		if Option["enable_trading"] == "1" && HLCross.PrevHLCrossTrade != "sell" {
			warning := "HLCross down, 卖出Sell Out---->市价" + getTradePrice("", Price[length-1]) +
				",委托价" + getTradePrice("sell", Price[length-1])
			logger.Infoln(warning)
			if Sell(getTradePrice("sell", Price[length-1]), tradeAmount) != "0" {
				warning += "[委托成功]"
				HLCross.PrevHLCrossTrade = "sell"
				HLCross.PrevBuyPirce = 0
			} else {
				warning += "[委托失败]"
			}

			go email.TriggerTrender(warning)
		}
	}

	//do sell when price is below stoploss point
	if Price[length-1] < HLCross.PrevBuyPirce*(1-stoploss*0.01) {
		if Option["enable_trading"] == "1" && HLCross.PrevHLCrossTrade != "sell" {
			warning := "stop loss, 卖出Sell Out---->市价" + getTradePrice("", Price[length-1]) + ",委托价" + getTradePrice("sell", Price[length-1])
			logger.Infoln(warning)
			if Sell(getTradePrice("sell", Price[length-1]), tradeAmount) != "0" {
				warning += "[委托成功]"
				HLCross.PrevHLCrossTrade = "sell"
				HLCross.PrevBuyPirce = 0
			} else {
				warning += "[委托失败]"
			}

			go email.TriggerTrender(warning)
		}
	}

	return true
}
