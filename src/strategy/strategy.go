package strategy

import (
	. "common"
	. "config"
	"email"
	"fmt"
	"logger"
	"strconv"
)

// Strategy is the interface that must be implemented by a strategy driver.
type Strategy interface {
	Tick(records []Record) bool
}

var strategys = make(map[string]Strategy)
var PrevTrade string
var PrevBuyPirce float64
var warning string

var gTradeAPI TradeAPI

// Register makes a strategy available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(strategyName string, strategy Strategy) {
	if strategy == nil {
		panic("sql: Register strategy is nil")
	}
	if _, dup := strategys[strategyName]; dup {
		panic("sql: Register called twice for strategy " + strategyName)
	}
	strategys[strategyName] = strategy
}

//entry call
func Tick(tradeAPI TradeAPI, records []Record) bool {

	strategyName := Option["strategy"]
	strategy, ok := strategys[strategyName]
	if !ok {
		logger.Errorf("sql: unknown strategy %q (forgotten import? private strategy?)", strategyName)
		return false
	}

	if strategyName != "OPENORDER" {
		length := len(records)
		//
		if length == 0 {
			logger.Errorln("warning:detect exception data", len(records))
			return false
		}

		//check exception data in trade center
		if checkException(records[length-2], records[length-1]) == false {
			logger.Errorln("detect exception data of trade center",
				records[length-2].Close, records[length-1].Close, records[length-1].Volumn)
			return false
		}
	}

	gTradeAPI = tradeAPI
	return strategy.Tick(records)
}

//check exception data in trade center
func checkException(recordPrev, recordNow Record) bool {
	if recordNow.Close > recordPrev.Close+10 && recordNow.Volumn < 1 {
		return false
	}

	if recordNow.Close < recordPrev.Close-10 && recordNow.Volumn < 1 {
		return false
	}

	return true
}

func getTradePrice(tradeDirection string, price float64) string {
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

func Buy(price, amount string) string {
	return gTradeAPI.Buy(price, amount)
}
func Sell(price, amount string) string {
	return gTradeAPI.Sell(price, amount)
}
func CancelOrder(order_id string) bool {
	return gTradeAPI.CancelOrder(order_id)
}
func GetAccount() (Account, bool) {
	return gTradeAPI.GetAccount()
}
func GetOrderBook() (ret bool, orderBook OrderBook) {
	return gTradeAPI.GetOrderBook()
}

func GetOrder(order_id string) (ret bool, order Order) {
	return gTradeAPI.GetOrder(order_id)
}

func GetAvailable_cny() float64 {
	account, ret := GetAccount()
	if !ret {
		logger.Errorln("GetAccount failed")
		return 0
	}

	numAvailable_cny, err := strconv.ParseFloat(account.Available_cny, 64)
	if err != nil {
		logger.Errorln("tradeAmount is not float")
		return 0
	}
	//balance > 0
	return numAvailable_cny
}

func GetAvailable_btc() float64 {
	account, ret := GetAccount()
	if !ret {
		logger.Errorln("GetAccount failed")
		return 0
	}

	numAvailable_btc, err := strconv.ParseFloat(account.Available_btc, 64)
	if err != nil {
		logger.Errorln("Available_btc is not float")
		return 0
	}
	//nCoins > 0
	return numAvailable_btc
}

func GetAvailable_ltc() float64 {
	account, ret := GetAccount()
	if !ret {
		logger.Errorln("GetAccount failed")
		return 0
	}

	numAvailable_ltc, err := strconv.ParseFloat(account.Available_ltc, 64)
	if err != nil {
		logger.Errorln("Available_ltc is not float")
		return 0
	}
	//nCoins > 0
	return numAvailable_ltc
}

func GetAvailable_coin() float64 {
	symbol := Option["symbol"]
	if symbol == "btc_cny" {
		return GetAvailable_btc()
	} else {
		return GetAvailable_ltc()
	}
}

//////////////////////////////////
//common stop loss function
//////////////////////////////////

func stop_loss_detect(Price []float64) bool {
	length := len(Price)

	stoploss, err := strconv.ParseFloat(Option["stoploss"], 64)
	if err != nil {
		logger.Errorln("config item stoploss is not float")
		return false
	}

	//do sell when price is below stoploss point
	stoplossPrice := PrevBuyPirce * (1 - stoploss*0.01)
	if Price[length-1] <= stoplossPrice {
		if Option["enable_trading"] == "1" && PrevTrade != "sell" {
			var tradePrice string
			if Option["discipleMode"] == "1" {
				if Price[length-1] > stoplossPrice {
					tradePrice = getTradePrice("sell", Price[length-1])
				} else {
					discipleValue, err := strconv.ParseFloat(Option["discipleValue"], 64)
					if err != nil {
						logger.Errorln("config item discipleValue is not float")
						return false
					}

					tradePrice = fmt.Sprintf("%f", PrevBuyPirce+discipleValue)
				}
			} else {
				tradePrice = getTradePrice("sell", Price[length-1])
			}

			warning := "stop loss, 卖出Sell Out---->市价" + getTradePrice("", Price[length-1]) + ",委托价" + tradePrice
			logger.Infoln(warning)

			tradeAmount := Option["tradeAmount"]
			if Sell(tradePrice, tradeAmount) != "0" {
				warning += "[委托成功]"
				PrevTrade = "sell"
				PrevBuyPirce = 0
			} else {
				warning += "[委托失败]"
			}

			go email.TriggerTrender(warning)
		}
	}

	return true
}
