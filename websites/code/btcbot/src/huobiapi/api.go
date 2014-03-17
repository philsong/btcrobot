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

package huobiapi

import (
	. "config"
	"fmt"
	"logger"
	"net/http"
)

type Huobi struct {
	client          *http.Client
	prevEMATrend    string
	Disable_trading int

	Peroid   int
	Slippage float64
	Time     []string
	Price    []float64
	Volumn   []float64
}

func NewHuobi() *Huobi {
	w := new(Huobi)
	return w
}

func (w Huobi) BuyIn(tradePrice, tradeAmount string) bool {
	tradeAPI := NewHuobiTrade(SecretOption["access_key"], SecretOption["secret_key"])
	buyId := tradeAPI.Buy(tradePrice, tradeAmount)
	logger.Infoln("buyId", buyId)
	if buyId != "0" {
		logger.Infoln("执行买入委托成功", tradePrice, tradeAmount)

		return true
	} else {
		logger.Infoln("执行买入委托失败", tradePrice, tradeAmount)
		return false
	}
}

func (w Huobi) SellOut(tradePrice, tradeAmount string) bool {
	tradeAPI := NewHuobiTrade(SecretOption["access_key"], SecretOption["secret_key"])
	sellId := tradeAPI.Sell(tradePrice, tradeAmount)
	logger.Infoln("sellId", sellId)
	if sellId != "0" {
		logger.Infoln("执行卖出委托成功", tradePrice, tradeAmount)
		w.SetPrevTrend("down")
		return true
	} else {
		logger.Infoln("执行卖出委托失败", tradePrice, tradeAmount)
		return false
	}
}

func (w Huobi) SetPrevTrend(trend string) {
	w.prevEMATrend = trend
}

func (w Huobi) GetPrevTrend() string {
	return w.prevEMATrend
}

func (w Huobi) GetTradePrice(tradeDirection string) string {
	if len(w.Price) == 0 {
		logger.Errorln("get price failed, array len=0")
		return "false"
	}
	var finalTradePrice float64
	if tradeDirection == "buy" {
		finalTradePrice = w.Price[len(w.Price)-1] + w.Slippage
	} else if tradeDirection == "sell" {
		finalTradePrice = w.Price[len(w.Price)-1] - w.Slippage
	} else {
		finalTradePrice = w.Price[len(w.Price)-1]
	}
	return fmt.Sprintf("%0.02f", finalTradePrice)
}
