package strategy

import (
	"fmt"
)

type TradeAPI interface {
	Buy(price, amount string) bool
	Sell(price, amount string) bool
	GetTradePrice(tradeDirection string) string
}

// Stragetey is the interface that must be implemented by a stragetey driver.
type Stragetey interface {
	Perform(tradeAPI TradeAPI, Time []string, Price []float64, Volumn []float64) bool
}

var strageteys = make(map[string]Stragetey)

// Register makes a stragetey available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(strageteyName string, stragetey Stragetey) {
	if stragetey == nil {
		panic("sql: Register stragetey is nil")
	}
	if _, dup := strageteys[strageteyName]; dup {
		panic("sql: Register called twice for stragetey " + strageteyName)
	}
	strageteys[strageteyName] = stragetey
}

//entry call
func Perform(strageteyName string, tradeAPI TradeAPI, Time []string, Price []float64, Volumn []float64) bool {
	stragetey, ok := strageteys[strageteyName]
	if !ok {
		fmt.Errorf("sql: unknown stragetey %q (forgotten import?)", strageteyName)
		return false
	}

	return stragetey.Perform(tradeAPI, Time, Price, Volumn)
}
