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
	"fmt"
	"logger"
	"strconv"
	"time"
)

type OOStrategy struct {
	PrevTime     string
	PrevPrice    float64
	PrevBuyPirce float64
	BuyId        []string
	SellId       []string
	BuyBegin     time.Time
	SellBegin    time.Time
}

func init() {
	oo := new(OOStrategy)
	Register("OPENORDER", oo)
}

//KDJ-EX strategy
func (oo *OOStrategy) Tick(records []Record) bool {

	const btcslap = 0.2
	const ltcslap = 0.05
	const ordercount = 1

	numTradeAmount, err := strconv.ParseFloat(Option["tradeAmount"], 64)
	if err != nil {
		logger.Errorln("config item tradeAmount is not float")
		return false
	}

	nSplitTradeAmount := numTradeAmount / float64(ordercount)
	splitTradeAmount := fmt.Sprintf("%f", nSplitTradeAmount)

	ret, orderbook := GetOrderBook()
	if !ret {
		logger.Infoln("get orderbook failed 1")
		ret, orderbook = GetOrderBook() //try again
		if !ret {
			logger.Infoln("get orderbook failed 2")
			return false
		}
	}

	logger.Infoln("卖一", orderbook.Asks[len(orderbook.Asks)-1])
	logger.Infoln("买一", orderbook.Bids[0])

	diff := 0.05

	if orderbook.Bids[0].Price+diff <= orderbook.Asks[len(orderbook.Asks)-1].Price {
		for i := 1; i <= ordercount; i++ {
			ret := true
			warning := "oo, 买入buy In<----限价单"
			tradePrice := fmt.Sprintf("%f", orderbook.Bids[0].Price+0.01)
			buyID := Buy(tradePrice, splitTradeAmount)
			if buyID != "0" {
				warning += "[委托成功]"
				ret = true
			} else {
				warning += "[委托失败]"
				ret = false
			}

			logger.Infoln(warning)

			if ret {
				warning = "oo, 卖出Sell Out---->限价单"
				tradePrice = fmt.Sprintf("%f", orderbook.Asks[len(orderbook.Asks)-1].Price-0.01)
				sellID := Sell(tradePrice, splitTradeAmount)
				if sellID != "0" {
					warning += "[委托成功]"
				} else {
					warning += "[委托失败]"
				}

				logger.Infoln(warning)
			}
		}
	}

	processTimeout()

	return true
}
