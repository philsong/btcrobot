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

package simulate

import (
	. "common"
	. "config"
	"logger"
	. "util"
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

type Simulate struct {
	tradeAPI *SimulateTrade
}

func NewSimulate() *Simulate {
	w := new(Simulate)
	w.tradeAPI = NewSimulateTrade()
	return w
}

func (w Simulate) CancelOrder(order_id string) (ret bool) {
	return true
}

func (w Simulate) GetOrderBook() (ret bool, orderBook OrderBook) {
	symbol := Option["symbol"]
	return w.getOrderBook(symbol)
}

func (w Simulate) GetOrder(order_id string) (ret bool, order Order) {
	order.Id = int(StringToInteger(order_id))
	order.Price = 1000
	order.Amount = 1.0
	order.Deal_amount = 1.0

	ret = true

	return
}

func (w Simulate) GetKLine(peroid int) (ret bool, records []Record) {
	symbol := Option["symbol"]
	return w.AnalyzeKLinePeroid(symbol, peroid)
}

func (w Simulate) GetAccount() (account Account, ret bool) {
	LoadSimulate()

	account.Available_cny = SimAccount["CNY"]
	account.Available_btc = SimAccount["BTC"]
	account.Available_ltc = SimAccount["LTC"]

	account.Frozen_cny = "0.0"
	account.Frozen_btc = "0.0"
	account.Frozen_ltc = "0.0"

	ret = true

	return
}

func (w Simulate) Buy(tradePrice, tradeAmount string) (buyId string) {
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

	_, ret := w.GetAccount()
	if !ret {
		logger.Infoln("GetAccount failed")
	}

	return buyId
}

func (w Simulate) Sell(tradePrice, tradeAmount string) (sellId string) {
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

	_, ret := w.GetAccount()
	if !ret {
		logger.Infoln("GetAccount failed")
	}

	return sellId
}
