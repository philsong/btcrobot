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

package peatio

import (
	. "common"
	. "config"
	"fmt"
	"logger"
	"strconv"
	"time"
)

type Peatio struct {
}

func NewPeatio() *Peatio {
	w := new(Peatio)
	return w
}

func (w Peatio) GetKLine(peroid int) (ret bool, records []Record) {
	symbol := Option["symbol"]
	return w.AnalyzeKLinePeroid(symbol, peroid)
}

func (w Peatio) GetOrderBook() (ret bool, orderBook OrderBook) {
	symbol := Option["symbol"]
	return w.getOrderBook(symbol)
}

func (w Peatio) Buy(tradePrice, tradeAmount string) (buyId string) {
	tradeAPI, _ := newPeatio(SecretOption["peatio_access_key"], SecretOption["peatio_secret_key"], "btc", 0)
	ibuyId, _ := tradeAPI.Buy(toFloat(tradePrice), toFloat(tradeAmount))
	buyId = fmt.Sprintf("%d", ibuyId)

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

func (w Peatio) Sell(tradePrice, tradeAmount string) (sellId string) {
	tradeAPI, _ := newPeatio(SecretOption["peatio_access_key"], SecretOption["peatio_secret_key"], "btc", 0)
	isellId, _ := tradeAPI.Buy(toFloat(tradePrice), toFloat(tradeAmount))
	sellId = fmt.Sprintf("%d", isellId)

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

func (w Peatio) GetOrder(order_id string) (ret bool, order Order) {
	symbol := Option["symbol"]
	if symbol == "ltc_cny" {
		ret = false
		return
	}

	tradeAPI, _ := newPeatio(SecretOption["peatio_access_key"], SecretOption["peatio_secret_key"], "btc", 0)
	iorder_id, _ := strconv.Atoi(order_id)
	pOrder, err := tradeAPI.GetOrder(int64(iorder_id))
	if err != nil {
		ret = false
		return
	}

	order.Id = int(pOrder.Id)
	order.Price = pOrder.Price
	order.Amount = pOrder.Amount
	order.Deal_amount = pOrder.DealAmount

	ret = true
	return
}

func (w Peatio) CancelOrder(order_id string) (ret bool) {
	tradeAPI, _ := newPeatio(SecretOption["peatio_access_key"], SecretOption["peatio_secret_key"], "btc", 0)
	iorder_id, _ := strconv.Atoi(order_id)

	ret, _ = tradeAPI.CancelOrder(int64(iorder_id))
	return
}

func (w Peatio) GetAccount() (account Account, ret bool) {
	tradeAPI, _ := newPeatio(SecretOption["peatio_access_key"], SecretOption["peatio_secret_key"], "btc", 0)

	userInfo, err := tradeAPI.GetAccount()

	if err != nil {
		logger.Traceln("Peatio GetAccount failed")
		ret = false
		return
	} else {
		ret = true
		account.Available_cny = float2str(userInfo.Balance)
		account.Available_btc = float2str(userInfo.Stocks)
		account.Available_ltc = "N/A"

		account.Frozen_cny = float2str(userInfo.FrozenBalance)
		account.Frozen_btc = float2str(userInfo.FrozenStocks)
		account.Frozen_ltc = "N/A"

		logger.Infof("Peatio资产: \n 可用cny:%-10s \tbtc:%-10s \tltc:%-10s \n 冻结cny:%-10s \tbtc:%-10s \tltc:%-10s\n",
			account.Available_cny,
			account.Available_btc,
			account.Available_ltc,
			account.Frozen_cny,
			account.Frozen_btc,
			account.Frozen_ltc)
		return
	}
}
