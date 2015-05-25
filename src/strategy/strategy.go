/*
  btcrobot is a Bitcoin, Litecoin and Altcoin trading bot written in golang,
  it features multiple trading methods using technical analysis.

  Disclaimer:

  USE AT YOUR OWN RISK!

  The author of this project is NOT responsible for any damage or loss caused
  by this software. There can be bugs and the bot may not Tick as expected
  or specified. Please consider testing it first with paper trading /
  backtesting on historical data. Also look at the code to see what how
  it's working.

  Weibo:http://weibo.com/bocaicfa
*/

package strategy

import (
	. "common"
	. "config"
	"db"
	"email"
	"fmt"
	"logger"
	"strconv"
	"time"
)

// Strategy is the interface that must be implemented by a strategy driver.
type Strategy interface {
	Tick(records []Record) bool
}

const timeout = 10 // minute

var magic int64
var strategys = make(map[string]Strategy)
var PrevTrade string
var PrevBuyPirce float64
var warning string

var buyOrders map[time.Time]string
var dealOrders map[time.Time]Order
var sellOrders map[time.Time]string
var recancelbuyOrders map[time.Time]string
var resellOrders map[time.Time]string

var buy_average float64
var buy_amount float64
var lastPrice float64
var length int
var isStoploss bool

var gTradeAPI TradeAPI

func init() {
	buyOrders = make(map[time.Time]string)
	dealOrders = make(map[time.Time]Order)
	sellOrders = make(map[time.Time]string)
	recancelbuyOrders = make(map[time.Time]string)
	resellOrders = make(map[time.Time]string)

	buy_average = 0

	PrevTrade = "init"
}

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

// entry call
func Tick(tradeAPI TradeAPI, records []Record) bool {
	strategyName := Option["strategy"]
	strategy, ok := strategys[strategyName]
	if !ok {
		logger.Errorf("sql: unknown strategy %q (forgotten import? private strategy?)", strategyName)
		return false
	}

	logger.Infoln("strategyName", strategyName)

	if strategyName != "OPENORDER" {
		length = len(records)
		if length == 0 {
			logger.Errorln("warning:detect exception data", len(records))
			return false
		}

		fmt.Println(length, records)
		if length < 2 {
			logger.Errorln("warning:detect exception data", len(records))
			return false
		}

		if Option["tick_interval"] != "1" {
			// check exception data in trade center
			if CheckException(records[length-2], records[length-1]) == false {
				logger.Errorln("detect exception data of trade center",
					records[length-2].Close, records[length-1].Close, records[length-1].Volumn)
				return false
			}
		}

		lastPrice = records[length-1].Close
	}

	gTradeAPI = tradeAPI

	return strategy.Tick(records)
}

// check exception data in trade center
func CheckException(recordPrev, recordNow Record) bool {
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
		finalTradePrice = price * (1 + slippage*0.01)
	} else if tradeDirection == "sell" {
		finalTradePrice = price * (1 - slippage*0.01)
	} else {
		finalTradePrice = price
	}

	return fmt.Sprintf("%0.02f", finalTradePrice)
}

func getOrderPrice() (sell1, buy1 float64, ret bool) {
	if GetBacktest() {
		sell1 = GetBtPrice()
		buy1 = GetBtPrice()
		ret = true
		return
	}

	ret, orderBook := GetOrderBook()
	if !ret {
		logger.Infoln("get orderBook failed 1")
		ret, orderBook = GetOrderBook() // try again
		if !ret {
			logger.Infoln("get orderBook failed 2")
			return 0, 0, ret
		}
	}

	logger.Infoln("卖一", orderBook.Asks[len(orderBook.Asks)-1])
	logger.Infoln("买一", orderBook.Bids[0])

	return orderBook.Asks[len(orderBook.Asks)-1].Price, orderBook.Bids[0].Price, ret
}

func getBuyPrice() (price string, nPrice float64, warning string) {
	// compute the price
	slippage, err := strconv.ParseFloat(Option["slippage"], 64)
	if err != nil {
		logger.Debugln("config item slippage is not float")
		slippage = 0.01
	}

	_, buy1, ret := getOrderPrice()
	if !ret {
		return
	}

	nPrice = buy1 * (1 + slippage*0.01)
	price = fmt.Sprintf("%f", nPrice)
	warning += "---->限价单" + price

	return
}

func getSellPrice() (price string, nPrice float64, warning string) {
	// compute the price
	slippage, err := strconv.ParseFloat(Option["slippage"], 64)
	if err != nil {
		logger.Debugln("config item slippage is not float")
		slippage = 0.01
	}

	sell1, _, ret := getOrderPrice()
	if !ret {
		return
	}
	nPrice = sell1 * (1 - slippage*0.01)

	if !isStoploss && Option["discipleMode"] == "1" {
		if nPrice < PrevBuyPirce {
			discipleValue, err := strconv.ParseFloat(Option["discipleValue"], 64)
			if err != nil {
				logger.Errorln("config item discipleValue is not float")
				discipleValue = 0.01
			}

			nPrice = PrevBuyPirce + discipleValue
		}
	}

	price = fmt.Sprintf("%f", nPrice)
	warning += "---->限价单" + price

	return
}

func buy(price, amount string) string {
	nPrice, err := strconv.ParseFloat(price, 64)
	if err != nil {
		logger.Debugln("price is not float")
		return "0"
	}

	nAmount, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		logger.Infoln("amount is not float")
		return "0"
	}

	Available_cny := GetAvailable_cny()
	if Available_cny < nPrice*nAmount {
		var nMinTradeAmount float64
		nAmount = Available_cny / nPrice
		symbol := Option["symbol"]
		if symbol == "btc_cny" {
			nMinTradeAmount = 0.1
		} else {
			nMinTradeAmount = 0.01
		}
		if nAmount < nMinTradeAmount {
			warning += "没有足够的法币余额"
			logger.Infoln(warning)
			PrevTrade = "buy"
			PrevBuyPirce = nPrice
			return "0"
		}
	}

	var buyID string
	if Option["enable_trading"] == "1" {
		buyID = gTradeAPI.Buy(price, amount)
	} else {
		buyID = "-1"
	}

	if buyID != "0" {
		timestamp := time.Now()
		magic += 1
		if !GetBacktest() {
			err := db.SetTx("buy", buyID, timestamp.Unix(), amount, price, magic)
			if err != nil {
				fmt.Println("SetTx", err)
			}
		}
		buyOrders[timestamp] = buyID
		PrevTrade = "buy"
		PrevBuyPirce = nPrice
	}

	if !GetBacktest() {
		cmd, id, timestamp, amount, price, magic, err := db.GetTx(buyID)
		if err != nil {
			fmt.Println("GetTx", err)
		} else {
			fmt.Println("buy items:", cmd, id, timestamp, amount, price, magic)
		}
	}

	return buyID
}

func sell(price, amount string) string {
	Available_coin := GetAvailable_coin()
	if Available_coin < 0.01 {
		warning = "没有足够的币"
		logger.Infoln(warning)
		PrevTrade = "sell"
		PrevBuyPirce = 0
		return "0"
	}

	nAmount, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		logger.Infoln("amount is not float")
		return "0"
	}

	if nAmount > Available_coin {
		nAmount = Available_coin
		amount = fmt.Sprintf("%02f", nAmount)
	}

	var sellID string
	if Option["enable_trading"] == "1" {
		sellID = gTradeAPI.Sell(price, amount)
	} else {
		sellID = "-1"
	}

	if sellID != "0" {
		timestamp := time.Now()
		magic += 1
		if !GetBacktest() {
			err := db.SetTx("sell", sellID, timestamp.Unix(), price, amount, magic)
			if err != nil {
				fmt.Println("SetTx", err)
			}
		}
		sellOrders[timestamp] = sellID
		PrevTrade = "sell"
		PrevBuyPirce = 0
	}

	if !GetBacktest() {
		cmd, id, timestamp, amount, price, magic, err := db.GetTx(sellID)
		if err != nil {
			fmt.Println("GetTx", err)
		} else {
			fmt.Println("sell items:", cmd, id, timestamp, amount, price, magic)
		}
	}

	return sellID
}

func Buy() string {
	if PrevTrade == "buy" {
		return "0"
	}

	// init
	isStoploss = false

	// compute the price
	price, nPrice, warning := getBuyPrice()

	// compute the amount
	amount := Option["tradeAmount"]
	nAmount, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		logger.Infoln("amount is not float")
		return "0"
	}

	Available_cny := GetAvailable_cny()
	if Available_cny < nPrice*nAmount {
		var nMinTradeAmount float64
		nAmount = Available_cny / nPrice
		symbol := Option["symbol"]
		if symbol == "btc_cny" {
			nMinTradeAmount = 0.1
		} else {
			nMinTradeAmount = 0.01
		}
		if nAmount < nMinTradeAmount {
			warning += "没有足够的法币余额"
			logger.Infoln(warning)
			PrevTrade = "buy"
			PrevBuyPirce = nPrice
			return "0"
		}

		amount = fmt.Sprintf("%02f", nAmount)
	}

	warning += "---->数量" + amount

	buyID := buy(price, amount)
	if buyID == "-1" {
		warning += " [模拟]"
	} else if buyID == "0" {
		warning += "[委托失败]"
	} else {
		warning += "[委托成功]" + buyID
	}

	logger.Infoln(warning)

	var coin string
	if Option["symbol"] == "btc_cny" {
		coin = "比特币"
	} else {
		coin = "莱特币"
	}

	if buyID != "0" {
		if !GetBacktest() {
			logger.Tradef("在%s，根据策略%s周期%s，以价格%s买入%s个%s\n", Option["tradecenter"], Option["strategy"], Option["tick_interval"], price, amount, coin)
			go email.TriggerTrender(warning)
		} else {
			t := time.Unix(GetBtTime(), 0)
			logger.Backtestf("%s 在simulate，根据策略%s周期%s，以价格%s买入%s个%s\n", t.Format("2006-01-02 15:04:05"), Option["strategy"], Option["tick_interval"], price, amount, coin)
		}
	}

	return buyID
}

func Sell() string {
	if PrevTrade == "sell" {
		return "0"
	}

	// compute the price
	price, _, warning := getSellPrice()

	// compute the amount
	Available_coin := GetAvailable_coin()
	if Available_coin < 0.01 {
		warning = "没有足够的币"
		logger.Infoln(warning)
		PrevTrade = "sell"
		PrevBuyPirce = 0
		return "0"
	}

	amount := Option["tradeAmount"]
	nAmount, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		logger.Infoln("amount is not float")
		return "0"
	}

	if nAmount > Available_coin {
		nAmount = Available_coin
		amount = fmt.Sprintf("%02f", nAmount)
	}

	sellID := sell(price, amount)
	if sellID == "-1" {
		warning += " [模拟]"
	} else if sellID == "0" {
		warning += "[委托失败]"
	} else {
		warning += "[委托成功]" + sellID
	}

	logger.Infoln(warning)

	var coin string
	if Option["symbol"] == "btc_cny" {
		coin = "比特币"
	} else {
		coin = "莱特币"
	}

	if sellID != "0" {
		if !GetBacktest() {
			logger.Tradef("在%s，根据策略%s周期%s，以价格%s卖出%s个%s\n", Option["tradecenter"], Option["strategy"], Option["tick_interval"], price, amount, coin)
			go email.TriggerTrender(warning)
		} else {
			t := time.Unix(GetBtTime(), 0)
			logger.Backtestf("%s 在simulate，根据策略%s周期%s，以价格%s卖出%s个%s\n", t.Format("2006-01-02 15:04:05"), Option["strategy"], Option["tick_interval"], price, amount, coin)
		}
	}

	return sellID
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
		return -1
	}

	numAvailable_cny, err := strconv.ParseFloat(account.Available_cny, 64)
	if err != nil {
		logger.Errorln("Available_cny is not float")
		return -1
	}
	return numAvailable_cny
}

func GetAvailable_btc() float64 {
	account, ret := GetAccount()
	if !ret {
		logger.Errorln("GetAccount failed")
		return -1
	}

	numAvailable_btc, err := strconv.ParseFloat(account.Available_btc, 64)
	if err != nil {
		logger.Errorln("Available_btc is not float")
		return -1
	}
	return numAvailable_btc
}

func GetAvailable_ltc() float64 {
	account, ret := GetAccount()
	if !ret {
		logger.Errorln("GetAccount failed")
		return -1
	}

	numAvailable_ltc, err := strconv.ParseFloat(account.Available_ltc, 64)
	if err != nil {
		logger.Errorln("Available_ltc is not float")
		return -1
	}
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
// common stop loss function
//////////////////////////////////

func processStoploss(Price float64) bool {
	stoploss, err := strconv.ParseFloat(Option["stoploss"], 64)
	if err != nil {
		logger.Errorln("config item stoploss is not float")
		return false
	}

	// do sell when price is below stoploss point
	stoplossPrice := PrevBuyPirce * (1 - stoploss*0.01)
	if Price <= stoplossPrice {
		warning := "stop loss, 卖出Sell Out---->"
		logger.Infoln(warning)
		logger.Infoln(Price, stoplossPrice, PrevBuyPirce, stoploss)

		isStoploss = true

		Sell()
	}

	return true
}

// check timeout trade
func processTimeout() bool {
	// last cancel failed, recancel
	for tm, id := range recancelbuyOrders {
		warning := fmt.Sprintf("<-----re-cancel %s-------------->", id)
		if CancelOrder(id) {
			warning += "[Cancel委托成功]"
			delete(recancelbuyOrders, tm)
		} else {
			warning += "[Cancel委托失败]"
			errno := GetLastError()
			if errno == 10009 {
				logger.Infoln(errno)
				delete(recancelbuyOrders, tm)
			}
		}

		logger.Infoln(warning)
		time.Sleep(1 * time.Second)
		time.Sleep(500 * time.Microsecond)
	}

	for tm, tradeAmount := range resellOrders {
		warning := fmt.Sprintf("<-----re-sell %f-------------->", tradeAmount)
		logger.Infoln(warning)
		sellID := Sell()
		if sellID != "0" {
			warning += "[re-sell委托成功]"
			delete(resellOrders, tm)
			sellOrders[time.Now()] = sellID //append or just update "set"
		} else {
			warning += "[re-sell委托失败]"
		}

		logger.Infoln(warning)
		time.Sleep(1 * time.Second)
		time.Sleep(500 * time.Microsecond)
	}

	now := time.Now()
	if len(buyOrders) != 0 {
		// todo
		logger.Infoln("BuyId len", len(buyOrders))
		for tm, id := range buyOrders {
			ret, order := GetOrder(id)
			if ret == false {
				continue
			}
			if order.Amount == order.Deal_amount {
				buy_average = (buy_amount*buy_average + order.Deal_amount*order.Price) / (buy_amount + order.Deal_amount)
				logger.Infof("buy_average full=%0.02f,%0.02f,%0.02f\n", buy_average, order.Deal_amount, buy_amount)

				dealOrders[tm] = order
				buy_amount += order.Deal_amount
				delete(buyOrders, tm)
			} else {
				if int64(now.Sub(tm)/time.Minute) <= timeout {
					continue
				}

				if order.Deal_amount > 0.0001 { // 部分成交的买卖单
					buy_average = (buy_amount*buy_average + order.Deal_amount*order.Price) / (buy_amount + order.Deal_amount)
					logger.Infof("buy_average partial=%0.02f,%0.02f,%0.02f\n", buy_average, order.Deal_amount, buy_amount)
					dealOrders[tm] = order
					buy_amount += order.Deal_amount
				}

				warning := fmt.Sprintf("<-----buy Delegation timeout, cancel %s[deal:%f]-------------->", id, order.Deal_amount)
				logger.Infoln(order)
				if CancelOrder(id) {
					warning += "[buy Cancel委托成功]"
				} else {
					warning += "[buy Cancel委托失败]"
					recancelbuyOrders[time.Now()] = id
				}

				delete(buyOrders, tm)
				logger.Infoln(warning)
				time.Sleep(1 * time.Second)
				time.Sleep(500 * time.Microsecond)
			}
		}
	}

	if len(sellOrders) != 0 {
		// todo
		logger.Infoln("SellId len", len(sellOrders))
		for tm, id := range sellOrders {
			if int64(now.Sub(tm)/time.Second) <= timeout {
				continue
			}

			ret, order := GetOrder(id)
			if ret == false {
				continue
			}

			if order.Amount == order.Deal_amount {
				delete(sellOrders, tm)
			} else {
				if int64(now.Sub(tm)/time.Minute) <= timeout {
					continue
				}
				if order.Deal_amount < order.Amount {
					ret, orderBook := GetOrderBook()
					if !ret {
						logger.Infoln("get orderBook failed 1")
						ret, orderBook = GetOrderBook() // try again
						if !ret {
							logger.Infoln("get orderBook failed 2")
							return false
						}
					}

					warning := fmt.Sprintf("<-----sell Delegation timeout, cancel %s[deal:%f]-------------->", id, order.Deal_amount)
					if CancelOrder(id) {
						warning += "[sell Cancel委托成功]"

						delete(sellOrders, tm)
						// update to delete, start a new order for sell in below

						sell_amount := order.Amount - order.Deal_amount

						logger.Infoln("卖一", (orderBook.Asks[len(orderBook.Asks)-1]))
						logger.Infoln("买一", orderBook.Bids[0])

						warning := "timeout, resell 卖出Sell Out---->限价单"
						tradePrice := fmt.Sprintf("%f", orderBook.Asks[len(orderBook.Asks)-1].Price-0.01)
						tradeAmount := fmt.Sprintf("%f", sell_amount)
						sellID := sell(tradePrice, tradeAmount)
						if sellID != "0" {
							warning += "[resell 委托成功]"
							sellOrders[time.Now()] = sellID //append or just update "set"
						} else {
							warning += "[resell 委托失败]"
							resellOrders[time.Now()] = tradeAmount
						}
						logger.Infoln(warning)
					} else {
						warning += "[sell Cancel委托失败]"

						errno := GetLastError()
						if errno == 10009 {
							logger.Infoln(errno)
							delete(sellOrders, tm)
						}
					}
					logger.Infoln(warning)
					time.Sleep(1 * time.Second)
					time.Sleep(500 * time.Microsecond)
				}
			}
		}
	}

	return true
}

// todo:need to think about edge issue carefully
// compute any period k-line base on 1 minute kline
func getKLine(records []Record, periods int) (recordsN []Record) {
	length := len(records)
	lengthN := length / periods

	// Loop through the entire array.
	for i := 0; i < periods*lengthN; i = i + periods {
		var recordN Record
		recordN.TimeStr = records[i].TimeStr
		recordN.Time = records[i].Time

		recordN.Open = records[i].Open
		recordN.Close = records[i+periods-1].Close

		var LowPrice []float64
		var HignPrice []float64
		for j := 0; j < periods; j++ {
			LowPrice = append(LowPrice, records[i+j].Low)
			HignPrice = append(HignPrice, records[i+j].High)
			recordN.Volumn += records[i+j].Volumn
		}
		recordN.Low = arrayLowest(LowPrice)
		recordN.High = arrayHighest(HignPrice)

		// add points to the array.
		recordsN = append(recordsN, recordN)
	}

	return recordsN
}

func toString(s interface{}) string {
	if v, ok := s.(string); ok {
		return v
	}
	return fmt.Sprintf("%v", s)
}

func toFloat(s interface{}) float64 {
	var ret float64
	switch v := s.(type) {
	case float64:
		ret = v
	case int64:
		ret = float64(v)
	case string:
		ret, err := strconv.ParseFloat(v, 64)
		if err != nil {
			logger.Errorln("convert ", s, " to float failed")
			return ret
		}
	}
	return ret
}

func float2str(i float64) string {
	return strconv.FormatFloat(i, 'f', -1, 64)
}
