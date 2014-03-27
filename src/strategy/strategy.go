package strategy

import (
	"fmt"
	"logger"
)

type TradeAPI interface {
	Buy(price, amount string) bool
	Sell(price, amount string) bool
	GetTradePrice(tradeDirection string) string
}

// Strategy is the interface that must be implemented by a strategy driver.
type Strategy interface {
	Perform(tradeAPI TradeAPI, Time []string, Price []float64, Volumn []float64) bool
}

var strategys = make(map[string]Strategy)

// Register makes a strategy available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(strageteyName string, strategy Strategy) {
	if strategy == nil {
		panic("sql: Register strategy is nil")
	}
	if _, dup := strategys[strageteyName]; dup {
		panic("sql: Register called twice for strategy " + strageteyName)
	}
	strategys[strageteyName] = strategy
}

//entry call
func Perform(strageteyName string, tradeAPI TradeAPI, Time []string, Price []float64, Volumn []float64) bool {
	strategy, ok := strategys[strageteyName]
	if !ok {
		fmt.Errorf("sql: unknown strategy %q (forgotten import?)", strageteyName)
		return false
	}

	//
	if len(Time) == 0 || len(Price) == 0 || len(Volumn) == 0 {
		logger.Errorln("warning:detect exception data", len(Time), len(Price), len(Volumn))
		return false
	}

	length := len(Price)

	//check exception data in trade center
	if checkException(Price[length-2], Price[length-1], Volumn[length-1]) == false {
		logger.Errorln("detect exception data of trade center", Price[length-2], Price[length-1], Volumn[length-1])
		return false
	}

	return strategy.Perform(tradeAPI, Time, Price, Volumn)
}

//check exception data in trade center
func checkException(yPrevPrice, Price, Volumn float64) bool {
	if Price > yPrevPrice+10 && Volumn < 1 {
		return false
	}

	if Price < yPrevPrice-10 && Volumn < 1 {
		return false
	}

	return true
}
