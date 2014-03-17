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

type TradeAPI interface {
	AnalyzeKLine(peroid int) (ret bool)
	Buy(price, amount string) bool
	Sell(price, amount string) bool
	GetTradePrice(tradeDirection string) string
}

type Huobi struct {
	client   *http.Client
	tradeAPI *TradeAPI

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

func (w Huobi) Buy(tradePrice, tradeAmount string) bool {
	tradeAPI := NewHuobiTrade(SecretOption["access_key"], SecretOption["secret_key"])

	var buyId string
	if Option["symbol"] == "btc_cny" {
		buyId = tradeAPI.BuyBTC(tradePrice, tradeAmount)
	} else if Option["symbol"] == "ltc_cny" {
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

func (w Huobi) Sell(tradePrice, tradeAmount string) bool {
	tradeAPI := NewHuobiTrade(SecretOption["access_key"], SecretOption["secret_key"])

	var sellId string
	if Option["symbol"] == "btc_cny" {
		sellId = tradeAPI.SellBTC(tradePrice, tradeAmount)
	} else if Option["symbol"] == "ltc_cny" {
		sellId = tradeAPI.SellLTC(tradePrice, tradeAmount)
	}
	logger.Infoln("sellId", sellId)
	if sellId != "0" {
		logger.Infoln("执行卖出委托成功", tradePrice, tradeAmount)
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
