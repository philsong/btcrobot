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
	"logger"
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

func (w Okcoin) GetOrder(order_id string) (ret bool, order Order) {
	symbol := Option["symbol"]
	tradeAPI := NewOkcoinTrade(SecretOption["ok_partner"], SecretOption["ok_secret_key"])

	ret, ok_orderTable := tradeAPI.Get_order(symbol, order_id)
	if ret == false {
		return
	}

	order.Id = ok_orderTable.Orders[0].Orders_id

	order.Amount = ok_orderTable.Orders[0].Amount

	order.Deal_amount = ok_orderTable.Orders[0].Deal_amount

	return
}

func (w Okcoin) GetKLine(peroid int) (ret bool, records []Record) {
	symbol := Option["symbol"]
	return w.AnalyzeKLinePeroid(symbol, peroid)
}

func (w Okcoin) GetAccount() (account Account, ret bool) {
	tradeAPI := NewOkcoinTrade(SecretOption["ok_partner"], SecretOption["ok_secret_key"])

	userInfo, ret := tradeAPI.GetAccount()

	if !ret {
		logger.Traceln("okcoin GetAccount failed")
		return
	} else {
		logger.Traceln(userInfo)

		account.Available_cny = userInfo.Info.Funds.Free.CNY
		account.Available_btc = userInfo.Info.Funds.Free.BTC
		account.Available_ltc = userInfo.Info.Funds.Free.LTC

		account.Frozen_cny = userInfo.Info.Funds.Freezed.CNY
		account.Frozen_btc = userInfo.Info.Funds.Freezed.BTC
		account.Frozen_ltc = userInfo.Info.Funds.Freezed.LTC

		logger.Infof("okcoin资产: \n 可用cny:%-10s \tbtc:%-10s \tltc:%-10s \n 冻结cny:%-10s \tbtc:%-10s \tltc:%-10s\n",
			account.Available_cny,
			account.Available_btc,
			account.Available_ltc,
			account.Frozen_cny,
			account.Frozen_btc,
			account.Frozen_ltc)
		//logger.Infoln(Account)
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

	time.Sleep(2 * time.Second)

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
	_, ret := w.GetAccount()

	if !ret {
		logger.Infoln("GetAccount failed")
	}

	return sellId
}
