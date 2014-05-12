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
	"time"
)

type HLCrossStrategy struct {
	PrevClosePrice float64
	PrevHighPrice  float64
	PrevLowPrice   float64
}

func init() {
	HLCross := new(HLCrossStrategy)
	Register("HLCross", HLCross)
}

//HLCross strategy
func (HLCross *HLCrossStrategy) Tick(records []Record) bool {
	//read config

	tradeAmount := Option["tradeAmount"]

	numTradeAmount, err := strconv.ParseFloat(Option["tradeAmount"], 64)
	if err != nil {
		logger.Errorln("config item tradeAmount is not float")
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
	if records[length-2].Close > records[length-2].Open &&
		records[length-3].Close > records[length-3].Open &&
		records[length-4].Close > records[length-4].Open {
		if Option["enable_trading"] == "1" && PrevTrade != "buy" {
			if GetAvailable_cny() < numTradeAmount {
				warning = "HLCross up, but 没有足够的法币可买"
				PrevTrade = "buy"
			} else {
				var tradePrice string
				if true {
					ret, orderBook := GetOrderBook()
					if !ret {
						logger.Infoln("get orderBook failed 1")
						ret, orderBook = GetOrderBook() //try again
						if !ret {
							logger.Infoln("get orderBook failed 2")
							return false
						}
					}

					logger.Infoln("卖一", (orderBook.Asks[len(orderBook.Asks)-1]))
					logger.Infoln("买一", orderBook.Bids[0])

					tradePrice = fmt.Sprintf("%f", orderBook.Bids[0].Price+0.01)
					warning += "---->限价单" + tradePrice
				} else {
					tradePrice = getTradePrice("buy", Price[length-1])
					warning += "---->市价单" + tradePrice
				}

				buyID := Buy(tradePrice, tradeAmount)
				if buyID != "0" {
					warning += "[委托成功]" + buyID

					buyOrders[time.Now()] = buyID
					PrevTrade = "buy"
				} else {
					warning += "[委托失败]"
				}
			}

			logger.Infoln(warning)
			go email.TriggerTrender(warning)
		}
	} else if records[length-2].Close < records[length-2].Open &&
		records[length-3].Close < records[length-3].Open &&
		records[length-4].Close < records[length-4].Open {
		if Option["enable_trading"] == "1" && PrevTrade != "sell" {
			if GetAvailable_coin() < numTradeAmount {
				warning = "HLCross down, but 没有足够的币可卖"
				PrevTrade = "sell"
				PrevBuyPirce = 0
			} else {
				warning = "HLCross down, 卖出Sell Out---->市价" + getTradePrice("", Price[length-1])
				var tradePrice string
				if true {
					ret, orderBook := GetOrderBook()
					if !ret {
						logger.Infoln("get orderBook failed 1")
						ret, orderBook = GetOrderBook() //try again
						if !ret {
							logger.Infoln("get orderBook failed 2")
							return false
						}
					}

					logger.Infoln("卖一", (orderBook.Asks[len(orderBook.Asks)-1]))
					logger.Infoln("买一", orderBook.Bids[0])

					tradePrice = fmt.Sprintf("%f", orderBook.Asks[len(orderBook.Asks)-1].Price-0.01)
					warning += "---->限价单" + tradePrice
				} else {
					tradePrice = getTradePrice("sell", Price[length-1])
					warning += "---->市价单" + tradePrice
				}

				sellID := Sell(tradePrice, tradeAmount)
				if sellID != "0" {
					warning += "[委托成功]"
					sellOrders[time.Now()] = sellID
					PrevTrade = "sell"
					PrevBuyPirce = 0
				} else {
					warning += "[委托失败]"
				}
			}

			logger.Infoln(warning)

			go email.TriggerTrender(warning)
		}
	}

	//do sell when price is below stoploss point
	processStoploss(Price)

	processTimeout()

	return true
}
