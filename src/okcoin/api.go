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

var ErrorCodeMap = map[int64]string{
	10000: "必选参数不能为空",
	10001: "用户请求过于频繁",
	10002: "系统错误",
	10003: "未在请求限制列表中,稍后请重试",
	10004: "IP限制不能请求该资源",
	10005: "密钥不存在",
	10006: "用户不存在",
	10007: "签名不匹配",
	10008: "非法参数",
	10009: "订单不存在",
	10010: "余额不足",
	10011: "买卖的数量小于BTC/LTC最小买卖额度",
	10012: "当前网站暂时只支持btc_cny ltc_cny",
	10013: "此接口只支持https请求",
	10014: "下单价格不得≤0或≥1000000",
	10015: "下单价格与最新成交价偏差过大",
}

type Okcoin struct {
	tradeAPI *OkcoinTrade
}

func NewOkcoin() *Okcoin {
	w := new(Okcoin)
	w.tradeAPI = NewOkcoinTrade(SecretOption["ok_partner"], SecretOption["ok_secret_key"])
	return w
}

func (w Okcoin) CancelOrder(order_id string) (ret bool) {
	tradeAPI := w.tradeAPI
	symbol := Option["symbol"]
	return tradeAPI.Cancel_order(symbol, order_id)
}

func (w Okcoin) GetOrderBook() (ret bool, orderBook OrderBook) {
	symbol := Option["symbol"]
	return w.getOrderBook(symbol)
}

func (w Okcoin) GetOrder(order_id string) (ret bool, order Order) {
	symbol := Option["symbol"]
	tradeAPI := w.tradeAPI

	ret, ok_orderTable := tradeAPI.Get_order(symbol, order_id)
	if ret == false {
		return
	}

	order.Id = ok_orderTable.Orders[0].Orders_id
	order.Price = ok_orderTable.Orders[0].Avg_rate
	order.Amount = ok_orderTable.Orders[0].Amount
	order.Deal_amount = ok_orderTable.Orders[0].Deal_amount

	return
}

func (w Okcoin) GetKLine(peroid int) (ret bool, records []Record) {
	symbol := Option["symbol"]
	return w.AnalyzeKLinePeroid(symbol, peroid)
}

func (w Okcoin) GetAccount() (account Account, ret bool) {
	tradeAPI := w.tradeAPI

	userInfo, ret := tradeAPI.GetAccount()

	//logger.Infoln("account:", userInfo)
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
	tradeAPI := w.tradeAPI

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

	time.Sleep(1 * time.Second)

	return buyId
}

func (w Okcoin) Sell(tradePrice, tradeAmount string) (sellId string) {
	tradeAPI := w.tradeAPI
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

	time.Sleep(1 * time.Second)
	_, ret := w.GetAccount()

	if !ret {
		logger.Infoln("GetAccount failed")
	}

	return sellId
}
