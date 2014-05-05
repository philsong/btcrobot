package strategy

import (
	. "common"
	. "config"
	"fmt"
	"logger"
	"strconv"
)

// Strategy is the interface that must be implemented by a strategy driver.
type Strategy interface {
	Tick(records []Record) bool
}

var strategys = make(map[string]Strategy)
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
