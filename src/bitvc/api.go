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

package bitvc

import (
	. "common"
	. "config"
	"huobi"
	"logger"
	"strconv"
	"time"
)

type Bitvc struct {
}

func NewBitvc() *Bitvc {
	w := new(Bitvc)
	return w
}

func (w Bitvc) CancelOrder(order_id string) (ret bool) {
	tradeAPI := NewBitvcTrade(SecretOption["huobi_access_key"], SecretOption["huobi_secret_key"])
	symbol := Option["symbol"]
	if symbol == "btc_cny" {
		return tradeAPI.Cancel_order(order_id)
	} else if symbol == "ltc_cny" {
		return false
	}

	return false
}

func (w Bitvc) GetOrderBook() (ret bool, orderBook OrderBook) {
	tradeAPI := huobi.NewHuobi()
	return tradeAPI.GetOrderBook()
}

func (w Bitvc) GetOrder(order_id string) (ret bool, order Order) {
	tradeAPI := NewBitvcTrade(SecretOption["huobi_access_key"], SecretOption["huobi_secret_key"])

	symbol := Option["symbol"]
	if symbol == "ltc_cny" {
		ret = false
		return
	}

	ret, hbOrder := tradeAPI.Get_order(order_id)
	if ret == false {
		ret = false
		return
	}

	order.Id = hbOrder.Id

	Price, err := strconv.ParseFloat(hbOrder.Order_price, 64)
	if err != nil {
		logger.Errorln("config item order_price is not float")
		ret = false
		return
	}
	order.Price = Price

	Amount, err := strconv.ParseFloat(hbOrder.Order_amount, 64)
	if err != nil {
		logger.Errorln("config item order_amount is not float")
		ret = false
		return
	}
	order.Amount = Amount

	Deal_amount, err := strconv.ParseFloat(hbOrder.Processed_amount, 64)
	if err != nil {
		logger.Errorln("config item processed_amount is not float")
		ret = false
		return
	}

	order.Deal_amount = Deal_amount

	ret = true
	return
}

func (w Bitvc) GetAccount() (account Account, ret bool) {
	tradeAPI := NewBitvcTrade(SecretOption["huobi_access_key"], SecretOption["huobi_secret_key"])

	userInfo, ret := tradeAPI.GetAccount()

	if !ret {
		logger.Traceln("Bitvc GetAccount failed")

		return
	} else {
		account.Available_cny = userInfo.Available_cny_display
		account.Available_btc = userInfo.Available_btc_display
		account.Available_ltc = "N/A"

		account.Frozen_cny = userInfo.Frozen_cny_display
		account.Frozen_btc = userInfo.Frozen_btc_display
		account.Frozen_ltc = "N/A"

		logger.Infof("Bitvc资产: \n 可用cny:%-10s \tbtc:%-10s \tltc:%-10s \n 冻结cny:%-10s \tbtc:%-10s \tltc:%-10s\n",
			account.Available_cny,
			account.Available_btc,
			account.Available_ltc,
			account.Frozen_cny,
			account.Frozen_btc,
			account.Frozen_ltc)
		return
	}
}

func (w Bitvc) Buy(tradePrice, tradeAmount string) (buyId string) {
	tradeAPI := NewBitvcTrade(SecretOption["huobi_access_key"], SecretOption["huobi_secret_key"])

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
	_, ret := w.GetAccount()
	if !ret {
		logger.Infoln("GetAccount failed")
	}

	return buyId
}

func (w Bitvc) Sell(tradePrice, tradeAmount string) (sellId string) {
	tradeAPI := NewBitvcTrade(SecretOption["huobi_access_key"], SecretOption["huobi_secret_key"])

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
	_, ret := w.GetAccount()
	if !ret {
		logger.Infoln("GetAccount failed")
	}

	return sellId
}
