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
	"common"
	. "config"
	"fmt"
	"logger"
	"net/http"
	"strconv"
)

type Huobi struct {
	client *http.Client

	Time   []string
	Price  []float64
	Volumn []float64
}

func NewHuobi() *Huobi {
	w := new(Huobi)
	return w
}

func (w Huobi) AnalyzeKLine(peroid int) (ret bool) {
	symbol := Option["symbol"]
	if peroid == 1 {
		return w.AnalyzeKLineMinute(symbol)
	} else {
		return w.AnalyzeKLinePeroid(symbol, peroid)
	}
}

func (w Huobi) Get_account_info() (userMoney common.UserMoney, ret bool) {
	tradeAPI := NewHuobiTrade(SecretOption["huobi_access_key"], SecretOption["huobi_secret_key"])

	userInfo, ret := tradeAPI.Get_account_info()

	if !ret {
		logger.Traceln("Huobi Get_account_info failed")

		return
	} else {
		userMoney.Available_cny = userInfo.Available_cny_display
		userMoney.Available_btc = userInfo.Available_btc_display
		userMoney.Available_ltc = "N/A"

		userMoney.Frozen_cny = userInfo.Frozen_cny_display
		userMoney.Frozen_btc = userInfo.Frozen_btc_display
		userMoney.Frozen_ltc = "N/A"

		logger.Infof("Huobi资产: \n 可用cny:%-10s \tbtc:%-10s \tltc:%-10s \n 冻结cny:%-10s \tbtc:%-10s \tltc:%-10s\n",
			userMoney.Available_cny,
			userMoney.Available_btc,
			userMoney.Available_ltc,
			userMoney.Frozen_cny,
			userMoney.Frozen_btc,
			userMoney.Frozen_ltc)
		return
	}
}

func (w Huobi) Buy(tradePrice, tradeAmount string) bool {
	tradeAPI := NewHuobiTrade(SecretOption["huobi_access_key"], SecretOption["huobi_secret_key"])

	var buyId string
	if Option["symbol"] == "btc_cny" {
		buyId = tradeAPI.BuyBTC(tradePrice, tradeAmount)
	} else if Option["symbol"] == "ltc_cny" {
		buyId = tradeAPI.BuyLTC(tradePrice, tradeAmount)
	}

	if buyId != "0" {
		logger.Infoln("执行买入委托成功", tradePrice, tradeAmount)
	} else {
		logger.Infoln("执行买入委托失败", tradePrice, tradeAmount)
	}

	_, ret := w.Get_account_info()

	if !ret {
		logger.Infoln("Get_account_info failed")
	}

	if buyId != "0" {
		return true
	} else {
		return false
	}
}

func (w Huobi) Sell(tradePrice, tradeAmount string) bool {
	tradeAPI := NewHuobiTrade(SecretOption["huobi_access_key"], SecretOption["huobi_secret_key"])

	var sellId string
	if Option["symbol"] == "btc_cny" {
		sellId = tradeAPI.SellBTC(tradePrice, tradeAmount)
	} else if Option["symbol"] == "ltc_cny" {
		sellId = tradeAPI.SellLTC(tradePrice, tradeAmount)
	}

	if sellId != "0" {
		logger.Infoln("执行卖出委托成功", tradePrice, tradeAmount)
	} else {
		logger.Infoln("执行卖出委托失败", tradePrice, tradeAmount)
	}

	_, ret := w.Get_account_info()

	if !ret {
		logger.Infoln("Get_account_info failed")
	}

	if sellId != "0" {
		return true
	} else {
		return false
	}
}

func (w Huobi) GetTradePrice(tradeDirection string) string {
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
		finalTradePrice = w.Price[len(w.Price)-1] * (1 + slippage*0.001)
	} else if tradeDirection == "sell" {
		finalTradePrice = w.Price[len(w.Price)-1] * (1 - slippage*0.001)
	} else {
		finalTradePrice = w.Price[len(w.Price)-1]
	}

	return fmt.Sprintf("%0.02f", finalTradePrice)
}
