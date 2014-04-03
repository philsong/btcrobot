package strategy

import (
	. "common"
	"logger"
)

// Strategy is the interface that must be implemented by a strategy driver.
type Strategy interface {
	Perform(tradeAPI TradeAPI, records []Record) bool
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
func Perform(strageteyName string, tradeAPI TradeAPI, records []Record) bool {
	strategy, ok := strategys[strageteyName]
	if !ok {
		logger.Errorf("sql: unknown strategy %q (forgotten import?)", strageteyName)
		return false
	}

	//
	if len(records) == 0 {
		logger.Errorln("warning:detect exception data", len(records))
		return false
	}

	var Price []float64
	var Volumn []float64
	for _, v := range records {
		Price = append(Price, v.Close)
		Volumn = append(Volumn, v.Volumn)
	}

	length := len(Price)

	//check exception data in trade center
	if checkException(records[length-2], records[length-1]) == false {
		logger.Errorln("detect exception data of trade center", Price[length-2], Price[length-1], Volumn[length-1])
		return false
	}

	return strategy.Perform(tradeAPI, records)
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
