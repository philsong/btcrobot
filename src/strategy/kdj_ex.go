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
	"email"
	"fmt"
	"logger"
	"strconv"
	"time"
)

type KDJexStrategy struct {
	PrevKDJTrade string
	PrevTime     string
	PrevPrice    float64
	PrevBuyPirce float64
	BuyId        []string
	SellId       []string
	BuyBegin     time.Time
	SellBegin    time.Time
}

func init() {
	kdjex := new(KDJexStrategy)
	kdjex.PrevKDJTrade = "init"

	Register("KDJ-EX", kdjex)
}

func SendEmail(warning string) {
	if !GetBacktest() {
		go email.TriggerTrender(warning)
	}
}

// KDJ-EX strategy
func (kdjex *KDJexStrategy) Tick(records []Record) bool {
	const btcslap = 0.2
	const ltcslap = 0.01
	const timeout = 300 // 秒
	const ordercount = 5

	tradeAmount := Option["tradeAmount"]

	numTradeAmount, err := strconv.ParseFloat(Option["tradeAmount"], 64)
	if err != nil {
		logger.Errorln("config item tradeAmount is not float")
		return false
	}

	var slappage float64
	symbol := Option["symbol"]
	if symbol == "btc_cny" {
		slappage = btcslap
	} else {
		slappage = ltcslap
	}

	var coin string
	if Option["symbol"] == "btc_cny" {
		coin = "比特币"
	} else {
		coin = "莱特币"
	}

	nSplitTradeAmount := numTradeAmount / float64(ordercount)
	splitTradeAmount := fmt.Sprintf("%f", nSplitTradeAmount)

	stoploss, err := strconv.ParseFloat(Option["stoploss"], 64)
	if err != nil {
		logger.Errorln("config item stoploss is not float")
		return false
	}

	var Time []string
	var Price []float64
	var Volumn []float64
	for _, v := range records {
		Time = append(Time, v.TimeStr)
		Price = append(Price, v.Close)
		Volumn = append(Volumn, v.Volumn)
	}

	length := len(records)

	if kdjex.PrevTime == records[length-1].TimeStr &&
		kdjex.PrevPrice == records[length-1].Close {
		return false
	}

	// K线为白，D线为黄，J线为红，K in middle
	k, d, j := getKDJ(records)

	if kdjex.PrevTime != records[length-1].TimeStr ||
		kdjex.PrevPrice != records[length-1].Close {
		kdjex.PrevTime = records[length-1].TimeStr
		kdjex.PrevPrice = records[length-1].Close

		logger.Infoln(records[length-1].TimeStr, records[length-1].Close)
		logger.Infof("d(黄线）%0.0f\tk(白线）%0.0f\tj(红线）%0.0f\n", d[length-2], k[length-2], j[length-2])
		logger.Infof("d(黄线）%0.0f\tk(白线）%0.0f\tj(红线）%0.0f\n", d[length-1], k[length-1], j[length-1])
		if j[length-1] > k[length-1] {
			logger.Infoln("KDJ up trend")
		} else {
			logger.Infoln("KDJ down trend")
		}
	}

	if (j[length-2] < k[length-2] && k[length-2] < d[length-2]) &&
		(j[length-1] > k[length-1] && k[length-1] > d[length-1]) {

		logger.Infoln("----------------->KDJ up cross", kdjex.PrevKDJTrade, d[length-2])

		if kdjex.PrevKDJTrade != "buy" && j[length-2] <= 20 {
			// do buy
			ret, orderbook := GetOrderBook()
			if !ret {
				logger.Infoln("get orderbook failed 1")
				ret, orderbook = GetOrderBook() // try again
				if !ret {
					logger.Infoln("get orderbook failed 2")
					return false
				}
			}

			logger.Infoln("卖一", (orderbook.Asks[len(orderbook.Asks)-1]))
			logger.Infoln("买一", orderbook.Bids[0])

			logger.Infoln("X 两根K线最低价", records[length-2].Low, records[length-1].Low)
			logger.Infoln("X 两根K线最高价", records[length-2].High, records[length-1].High)

			avgLow := (records[length-2].Close + records[length-1].Low) / 2.0
			logger.Infoln("X 两根K线的最低平均价", avgLow)

			warning := "KDJ up cross, 买入buy In<----限价单"
			for i := 1; i <= ordercount; i++ {
				warning := "KDJ up cross, 买入buy In<----限价单"
				tradePrice := fmt.Sprintf("%f", avgLow+slappage*float64(i))
				buyID := buy(tradePrice, splitTradeAmount)
				if buyID != "0" {
					warning += "[委托成功]"
					kdjex.BuyId = append(kdjex.BuyId, buyID)
					if !GetBacktest() {
						logger.Tradef("在%s，根据策略%s周期%s，以价格%s买入%s个%s\n", Option["tradecenter"], Option["strategy"], Option["tick_interval"], tradePrice, splitTradeAmount, coin)
					} else {
						t := time.Unix(GetBtTime(), 0)
						logger.Backtestf("%s 在simulate，根据策略%s周期%s，以价格%s买入%s个%s\n", t.Format("2006-01-02 15:04:05"), Option["strategy"], Option["tick_interval"], tradePrice, splitTradeAmount, coin)
					}
				} else {
					warning += "[委托失败]"
				}

				logger.Infoln(warning)
			}

			kdjex.BuyBegin = time.Now()
			kdjex.PrevKDJTrade = "buy"

			kdjex.PrevBuyPirce = avgLow

			logger.Infoln("------------>>>stoploss price", kdjex.PrevBuyPirce*(1-stoploss*0.01))

			_, ret = GetAccount()
			if !ret {
				logger.Infoln("GetAccount failed")
			}

			SendEmail(warning)
		}
	}

	if (j[length-2] > k[length-2] && k[length-2] > d[length-2]) &&
		(j[length-1] < k[length-1] && k[length-1] < d[length-1]) {

		logger.Infoln("<----------------------KDJ down cross", kdjex.PrevKDJTrade, d[length-2])

		if kdjex.PrevKDJTrade != "sell" && j[length-2] >= 80 {
			// do sell
			ret, orderbook := GetOrderBook()
			if !ret {
				logger.Infoln("get orderbook failed 1")
				ret, orderbook = GetOrderBook() // try again
				if !ret {
					logger.Infoln("get orderbook failed 2")
					return false
				}
			}

			logger.Infoln("卖一", (orderbook.Asks[len(orderbook.Asks)-1]))
			logger.Infoln("买一", orderbook.Bids[0])

			logger.Infoln("X 两根K线最低价", records[length-2].Low, records[length-1].Low)
			logger.Infoln("X 两根K线最高价", records[length-2].High, records[length-1].High)

			avgHigh := (records[length-2].Close + records[length-1].High) / 2.0
			logger.Infoln("X 两根K线的最高平均价", avgHigh)
			warning := "KDJ down cross, 卖出Sell Out---->限价单"

			for i := 1; i <= ordercount; i++ {
				warning := "KDJ down cross, 卖出Sell Out---->限价单"
				tradePrice := fmt.Sprintf("%f", avgHigh-slappage*float64(i))
				sellID := sell(tradePrice, splitTradeAmount)
				if sellID != "0" {
					warning += "[委托成功]"
					kdjex.SellId = append(kdjex.SellId, sellID)
					if !GetBacktest() {
						logger.Tradef("在%s，根据策略%s周期%s，以价格%s卖出%s个%s\n", Option["tradecenter"], Option["strategy"], Option["tick_interval"], tradePrice, splitTradeAmount, coin)
					} else {
						t := time.Unix(GetBtTime(), 0)
						logger.Backtestf("%s 在simulate，根据策略%s周期%s，以价格%s卖出%s个%s\n", t.Format("2006-01-02 15:04:05"), Option["strategy"], Option["tick_interval"], tradePrice, splitTradeAmount, coin)
					}
				} else {
					warning += "[委托失败]"
				}

				logger.Infoln(warning)
			}

			kdjex.SellBegin = time.Now()
			kdjex.PrevKDJTrade = "sell"

			_, ret = GetAccount()
			if !ret {
				logger.Infoln("GetAccount failed")
			}
			SendEmail(warning)
		}
	}

	// do sell when price is below stoploss point

	if Price[length-1] < kdjex.PrevBuyPirce*(1-stoploss*0.01) {
		if Option["disable_trading"] != "1" && kdjex.PrevKDJTrade != "sell" {
			kdjex.PrevKDJTrade = "sell"
			kdjex.PrevBuyPirce = 0
			warning := "!<------------------stop loss, 卖出Sell Out---->市价" + getTradePrice("", Price[length-1]) + ",委托价" + getTradePrice("sell", Price[length-1])
			logger.Infoln(warning)
			tradePrice := getTradePrice("sell", Price[length-1])
			if sell(tradePrice, tradeAmount) != "0" {
				warning += "[委托成功]"
				if !GetBacktest() {
					logger.Tradef("在%s，根据策略%s周期%s，以价格%s卖出%s个%s\n", Option["tradecenter"], Option["strategy"], Option["tick_interval"], tradePrice, tradeAmount, coin)
				} else {
					t := time.Unix(GetBtTime(), 0)
					logger.Backtestf("%s 在simulate，根据策略%s周期%s，以价格%s卖出%s个%s\n", t.Format("2006-01-02 15:04:05"), Option["strategy"], Option["tick_interval"], tradePrice, tradeAmount, coin)
				}
			} else {
				warning += "[委托失败]"
				for i := 1; i <= ordercount; i++ {
					warning := "stop loss, 卖出Sell Out---->限价单"
					tradePrice := getTradePrice("sell", Price[length-1])
					sellID := sell(tradePrice, splitTradeAmount)
					if sellID != "0" {
						warning += "[委托成功]"
						kdjex.SellId = append(kdjex.SellId, sellID)
						if !GetBacktest() {
							logger.Tradef("在%s，根据策略%s周期%s，以价格%s卖出%s个%s\n", Option["tradecenter"], Option["strategy"], Option["tick_interval"], tradePrice, splitTradeAmount, coin)
						} else {
							t := time.Unix(GetBtTime(), 0)
							logger.Backtestf("%s 在simulate，根据策略%s周期%s，以价格%s卖出%s个%s\n", t.Format("2006-01-02 15:04:05"), Option["strategy"], Option["tick_interval"], tradePrice, splitTradeAmount, coin)
						}
					} else {
						warning += "[委托失败]"
					}

					logger.Infoln(warning)
				}
			}

			kdjex.SellBegin = time.Now()
			kdjex.PrevKDJTrade = "sell"

			_, ret := GetAccount()
			if !ret {
				logger.Infoln("GetAccount failed")
			}
			SendEmail(warning)
		}
	}

	// check timeout trade
	now := time.Now()

	logger.Infoln("time go ", int64(now.Sub(kdjex.BuyBegin)/time.Second))
	logger.Infoln("BuyId len", len(kdjex.BuyId), cap(kdjex.BuyId))
	logger.Infoln("SellId len", len(kdjex.SellId), cap(kdjex.SellId))

	if len(kdjex.BuyId) != 0 &&
		int64(now.Sub(kdjex.BuyBegin)/time.Second) > timeout {
		// todo
		for _, BuyId := range kdjex.BuyId {
			warning := "<--------------buy order timeout, cancel-------------->" + BuyId
			if CancelOrder(BuyId) {
				warning += "[Cancel委托成功]"
			} else {
				warning += "[Cancel委托失败]"
			}
			logger.Infoln(warning)
			time.Sleep(1 * time.Second)
			time.Sleep(500 * time.Microsecond)
		}
		kdjex.BuyId = kdjex.BuyId[:0]
	}

	if len(kdjex.SellId) != 0 &&
		int64(now.Sub(kdjex.SellBegin)/time.Second) > timeout {
		// todo
		for _, SellId := range kdjex.SellId {
			warning := "<--------------sell order timeout, cancel------------->" + SellId
			if CancelOrder(SellId) {
				warning += "[Cancel委托成功]"
			} else {
				warning += "[Cancel委托失败]"
			}
			logger.Infoln(warning)
			time.Sleep(1 * time.Second)
			time.Sleep(500 * time.Microsecond)
		}
		kdjex.SellId = kdjex.SellId[:0]
	}

	return true
}
