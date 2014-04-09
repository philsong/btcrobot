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
	. "common"
	. "config"
	"fmt"
	"logger"
	"strconv"
	"time"
)

type Okcoin struct {
}

func NewOkcoin() *Okcoin {
	w := new(Okcoin)
	return w
}

func (w Okcoin) CancelOrder(order_id string) (ret bool) {
	tradeAPI := NewOkcoinTrade(SecretOption["ok_partner"], SecretOption["ok_secret_key"])
	symbol := Option["symbol"]
	return tradeAPI.Cancel_order(symbol, order_id)
}

func (w Okcoin) GetOrderBook() (ret bool, orderBook OrderBook) {
	symbol := Option["symbol"]
	return w.getOrderBook(symbol)
}

func (w Okcoin) AnalyzeKLine(peroid int) (ret bool) {
	symbol := Option["symbol"]
	return w.AnalyzeKLinePeroid(symbol, peroid)
}

func (w Okcoin) Get_account_info() (userMoney UserMoney, ret bool) {
	tradeAPI := NewOkcoinTrade(SecretOption["ok_partner"], SecretOption["ok_secret_key"])

	userInfo, ret := tradeAPI.Get_account_info()

	if !ret {
		logger.Traceln("okcoin Get_account_info failed")
		return
	} else {
		logger.Traceln(userInfo)

		userMoney.Available_cny = userInfo.Info.Funds.Free.CNY
		userMoney.Available_btc = userInfo.Info.Funds.Free.BTC
		userMoney.Available_ltc = userInfo.Info.Funds.Free.LTC

		userMoney.Frozen_cny = userInfo.Info.Funds.Freezed.CNY
		userMoney.Frozen_btc = userInfo.Info.Funds.Freezed.BTC
		userMoney.Frozen_ltc = userInfo.Info.Funds.Freezed.LTC

		logger.Infof("okcoin资产: \n 可用cny:%-10s \tbtc:%-10s \tltc:%-10s \n 冻结cny:%-10s \tbtc:%-10s \tltc:%-10s\n",
			userMoney.Available_cny,
			userMoney.Available_btc,
			userMoney.Available_ltc,
			userMoney.Frozen_cny,
			userMoney.Frozen_btc,
			userMoney.Frozen_ltc)
		//logger.Infoln(userMoney)
		return
	}
}

func (w Okcoin) Buy(tradePrice, tradeAmount string) (buyId string) {
	tradeAPI := NewOkcoinTrade(SecretOption["ok_partner"], SecretOption["ok_secret_key"])

	symbol := Option["symbol"]
	if symbol == "btc_cny" {
		buyId = tradeAPI.BuyBTC(tradePrice, tradeAmount)
	} else if symbol == "ltc_cny" {
		buyId = tradeAPI.BuyLTC(tradePrice, tradeAmount)
	}

	if buyId != "0" {
		logger.Infoln("执行买入委托成功", tradePrice, tradeAmount)
	} else {
		logger.Infoln("执行买入委托失败", tradePrice, tradeAmount)
	}

	time.Sleep(3 * time.Second)
	_, ret := w.Get_account_info()

	if !ret {
		logger.Infoln("Get_account_info failed")
	}

	return buyId
}

func (w Okcoin) Sell(tradePrice, tradeAmount string) (sellId string) {
	tradeAPI := NewOkcoinTrade(SecretOption["ok_partner"], SecretOption["ok_secret_key"])

	symbol := Option["symbol"]
	if symbol == "btc_cny" {
		sellId = tradeAPI.SellBTC(tradePrice, tradeAmount)
	} else if symbol == "ltc_cny" {
		sellId = tradeAPI.SellLTC(tradePrice, tradeAmount)
	}
	if sellId != "0" {
		logger.Infoln("执行卖出委托成功", tradePrice, tradeAmount)
	} else {
		logger.Infoln("执行卖出委托失败", tradePrice, tradeAmount)
	}

	time.Sleep(3 * time.Second)
	_, ret := w.Get_account_info()

	if !ret {
		logger.Infoln("Get_account_info failed")
	}

	return sellId
}

func (w Okcoin) GetTradePrice(tradeDirection string, price float64) string {
	slippage, err := strconv.ParseFloat(Option["slippage"], 64)
	if err != nil {
		logger.Debugln("config item slippage is not float")
		slippage = 0
	}

	var finalTradePrice float64
	if tradeDirection == "buy" {
		finalTradePrice = price * (1 + slippage*0.001)
	} else if tradeDirection == "sell" {
		finalTradePrice = price * (1 - slippage*0.001)
	} else {
		finalTradePrice = price
	}

	return fmt.Sprintf("%0.02f", finalTradePrice)
}
