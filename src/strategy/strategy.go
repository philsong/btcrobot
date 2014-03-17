package strategy

import (
	"fmt"
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

	return strategy.Perform(tradeAPI, Time, Price, Volumn)
}
