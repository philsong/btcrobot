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
	. "config"
	. "util"
)

type SimOrder struct {
	OrderTime string
	OrderType string
	Symbol    string
	Price     float64
	Amount    float64
}

var SimOrders []SimOrder

type SimulateTrade struct {
	errno int64
}

func NewSimulateTrade() *SimulateTrade {
	w := new(SimulateTrade)
	return w
}

func (w *SimulateTrade) Cancel_order(symbol, order_id string) bool {
	return false
}

func (w *SimulateTrade) Cancel_BTCorder(order_id string) (ret bool) {
	return w.Cancel_order("btc_cny", order_id)
}

func (w *SimulateTrade) Cancel_LTCorder(order_id string) (ret bool) {
	return w.Cancel_order("ltc_cny", order_id)
}

func (w *SimulateTrade) BuyBTC(price, amount string) string {
	LoadSimulate()

	cny := StringToFloat(SimAccount["CNY"])
	cnyAmount := StringToFloat(price) * StringToFloat(amount)

	if cnyAmount > cny {
		cnyAmount = cny
	}

	btcAmount := cnyAmount / StringToFloat(price)
	if btcAmount >= 0.01 {
		SimAccount["CNY"] = FloatToString(cny - cnyAmount)
		SimAccount["BTC"] = FloatToString(btcAmount + StringToFloat(SimAccount["BTC"]))
		SaveSimulate()
		return RandomString(6)
	} else {
		return "0"
	}
}

func (w *SimulateTrade) SellBTC(price, amount string) string {
	LoadSimulate()

	btc := StringToFloat(SimAccount["BTC"])

	if btc > StringToFloat(amount) {
		btc = StringToFloat(amount)
	}

	if btc >= 0.01 {
		cny := StringToFloat(price) * btc
		SimAccount["CNY"] = FloatToString(StringToFloat(SimAccount["CNY"]) + cny)
		SimAccount["BTC"] = FloatToString(StringToFloat(SimAccount["BTC"]) - btc)
		SaveSimulate()
		return RandomString(6)
	} else {
		return "0"
	}
}

func (w *SimulateTrade) BuyLTC(price, amount string) string {
	LoadSimulate()

	cny := StringToFloat(SimAccount["CNY"])
	cnyAmount := StringToFloat(price) * StringToFloat(amount)

	if cnyAmount > cny {
		cnyAmount = cny
	}

	ltcAmount := cnyAmount / StringToFloat(price)
	if ltcAmount >= 0.1 {
		SimAccount["CNY"] = FloatToString(cny - cnyAmount)
		SimAccount["LTC"] = FloatToString(ltcAmount + StringToFloat(SimAccount["LTC"]))
		SaveSimulate()
		return RandomString(6)
	} else {
		return "0"
	}
}

func (w *SimulateTrade) SellLTC(price, amount string) string {
	LoadSimulate()

	ltc := StringToFloat(SimAccount["LTC"])

	if ltc > StringToFloat(amount) {
		ltc = StringToFloat(amount)
	}

	if ltc >= 0.1 {
		cny := StringToFloat(price) * ltc
		SimAccount["CNY"] = FloatToString(StringToFloat(SimAccount["CNY"]) + cny)
		SimAccount["LTC"] = FloatToString(StringToFloat(SimAccount["LTC"]) - ltc)
		SaveSimulate()
		return RandomString(6)
	} else {
		return "0"
	}
}
