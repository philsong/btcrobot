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

package okcoin

import (
	. "config"
	"fmt"
	"logger"
	"net/http"
	"strconv"
)

type TradeAPI interface {
	AnalyzeKLine(peroid int) (ret bool)
	Buy(price, amount string) bool
	Sell(price, amount string) bool
	GetTradePrice(tradeDirection string) string
}

type Okcoin struct {
	client *http.Client

	Time   []string
	Price  []float64
	Volumn []float64
}

func NewOkcoin() *Okcoin {
	w := new(Okcoin)
	return w
}

func (w Okcoin) AnalyzeKLine(peroid int) (ret bool) {
	symbol := Option["symbol"]
	return w.AnalyzeKLinePeroid(symbol, peroid)
}

func (w Okcoin) Buy(tradePrice, tradeAmount string) bool {
	tradeAPI := NewOkcoinTrade(SecretOption["ok_partner"], SecretOption["ok_secret_key"])

	var buyId string
	symbol := Option["symbol"]
	if symbol == "btc_cny" {
		buyId = tradeAPI.BuyBTC(tradePrice, tradeAmount)
	} else if symbol == "ltc_cny" {
		buyId = tradeAPI.BuyLTC(tradePrice, tradeAmount)
	}

	if buyId != "0" {
		logger.Infoln("执行买入委托成功", tradePrice, tradeAmount)
		return true
	} else {
		logger.Infoln("执行买入委托失败", tradePrice, tradeAmount)
		return false
	}
}

func (w Okcoin) Sell(tradePrice, tradeAmount string) bool {
	tradeAPI := NewOkcoinTrade(SecretOption["ok_partner"], SecretOption["ok_secret_key"])

	var sellId string
	symbol := Option["symbol"]
	if symbol == "btc_cny" {
		sellId = tradeAPI.SellBTC(tradePrice, tradeAmount)
	} else if symbol == "ltc_cny" {
		sellId = tradeAPI.SellLTC(tradePrice, tradeAmount)
	}

	if sellId != "0" {
		logger.Infoln("执行卖出委托成功", tradePrice, tradeAmount)
		return true
	} else {
		logger.Infoln("执行卖出委托失败", tradePrice, tradeAmount)
		return false
	}
}

func (w Okcoin) GetTradePrice(tradeDirection string) string {
	if len(w.Price) == 0 {
		logger.Errorln("get price failed, array len=0")
		return "false"
	}

	slippage, err := strconv.ParseFloat(Option["slippage"], 64)
	if err != nil {
		logger.Debugln("config item slippage is not float")
		slippage = 0
	}

	var finalTradePrice float64
	if tradeDirection == "buy" {
		finalTradePrice = w.Price[len(w.Price)-1] + slippage*0.001
	} else if tradeDirection == "sell" {
		finalTradePrice = w.Price[len(w.Price)-1] - slippage*0.001
	} else {
		finalTradePrice = w.Price[len(w.Price)-1]
	}
	return fmt.Sprintf("%0.02f", finalTradePrice)
}
