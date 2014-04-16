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
	. "common"
	. "config"
	"logger"
	"time"
)

type Huobi struct {
}

func NewHuobi() *Huobi {
	w := new(Huobi)
	return w
}

func (w Huobi) CancelOrder(order_id string) (ret bool) {
	tradeAPI := NewHuobiTrade(SecretOption["huobi_access_key"], SecretOption["huobi_secret_key"])
	symbol := Option["symbol"]
	if symbol == "btc_cny" {
		return tradeAPI.Cancel_order(order_id)
	} else if symbol == "ltc_cny" {
		return false
	}

	return false
}

func (w Huobi) GetOrderBook() (ret bool, orderBook OrderBook) {
	symbol := Option["symbol"]
	return w.getOrderBook(symbol)
}

func (w Huobi) GetKLine(peroid int) (ret bool, records []Record) {
	symbol := Option["symbol"]
	return w.AnalyzeKLinePeroid(symbol, peroid)
}

func (w Huobi) GetAccountInfo() (AccountInfo AccountInfo, ret bool) {
	tradeAPI := NewHuobiTrade(SecretOption["huobi_access_key"], SecretOption["huobi_secret_key"])

	userInfo, ret := tradeAPI.GetAccountInfo()

	if !ret {
		logger.Traceln("Huobi GetAccountInfo failed")

		return
	} else {
		AccountInfo.Available_cny = userInfo.Available_cny_display
		AccountInfo.Available_btc = userInfo.Available_btc_display
		AccountInfo.Available_ltc = "N/A"

		AccountInfo.Frozen_cny = userInfo.Frozen_cny_display
		AccountInfo.Frozen_btc = userInfo.Frozen_btc_display
		AccountInfo.Frozen_ltc = "N/A"

		logger.Infof("Huobi资产: \n 可用cny:%-10s \tbtc:%-10s \tltc:%-10s \n 冻结cny:%-10s \tbtc:%-10s \tltc:%-10s\n",
			AccountInfo.Available_cny,
			AccountInfo.Available_btc,
			AccountInfo.Available_ltc,
			AccountInfo.Frozen_cny,
			AccountInfo.Frozen_btc,
			AccountInfo.Frozen_ltc)
		return
	}
}

func (w Huobi) Buy(tradePrice, tradeAmount string) (buyId string) {
	tradeAPI := NewHuobiTrade(SecretOption["huobi_access_key"], SecretOption["huobi_secret_key"])

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

	time.Sleep(3 * time.Second)
	_, ret := w.GetAccountInfo()
	if !ret {
		logger.Infoln("GetAccountInfo failed")
	}

	return buyId
}

func (w Huobi) Sell(tradePrice, tradeAmount string) (sellId string) {
	tradeAPI := NewHuobiTrade(SecretOption["huobi_access_key"], SecretOption["huobi_secret_key"])

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

	time.Sleep(3 * time.Second)
	_, ret := w.GetAccountInfo()
	if !ret {
		logger.Infoln("GetAccountInfo failed")
	}

	return sellId
}
