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

package huobi

import (
	. "config"
	"fmt"
	"logger"
	"net/http"
	"strconv"
)

type MarketAPI interface {
	AnalyzeKLine(peroid int) (ret bool)
}

type TradeAPI interface {
	BuyIn(price, amount string) bool
	SellOut(price, amount string) bool
	GetTradePrice(tradeDirection string) string
	GetPrevTrend() string
	SetPrevTrend(trend string)
}

type Huobi struct {
	client       *http.Client
	tradeAPI     *HuobiTrade
	prevEMATrend string

	Time   []string
	Price  []float64
	Volumn []float64
}

func NewHuobi() *Huobi {
	w := new(Huobi)
	return w
}

func (w Huobi) AnalyzeKLine(peroid int) (ret bool) {
	if peroid == 1 {
		return w.AnalyzeKLineMinute()
	} else {
		return w.AnalyzeKLinePeroid(peroid)
	}
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

func (w Huobi) GetTradePrice(tradeDirection string) string {
	if len(w.Price) == 0 {
		logger.Errorln("get price failed, array len=0")
		return "false"
	}

	slippage, err := strconv.ParseFloat(Config["slippage"], 64)
	if err != nil {
		logger.Debugln("config item slippage is not float")
		slippage = 0
	}

	var finalTradePrice float64
	if tradeDirection == "buy" {
		finalTradePrice = w.Price[len(w.Price)-1] + slippage
	} else if tradeDirection == "sell" {
		finalTradePrice = w.Price[len(w.Price)-1] - slippage
	} else {
		finalTradePrice = w.Price[len(w.Price)-1]
	}
	return fmt.Sprintf("%0.02f", finalTradePrice)
}

func (w Huobi) SetPrevTrend(trend string) {
	w.prevEMATrend = trend
}

func (w Huobi) GetPrevTrend() string {
	return w.prevEMATrend
}
