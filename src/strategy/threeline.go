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
	"fmt"
	"logger"
	"time"
	. "util"
)

type ThreelineStrategy struct {
	PrevTrade string
	PrevTime  string
	BuyPrice  float64
	SellHalf  bool
}

func init() {
	threelineStrategy := new(ThreelineStrategy)
	Register("THREELINE", threelineStrategy)
}

func coinName() (coin string) {
	if Option["symbol"] == "btc_cny" {
		coin = "比特币"
	} else {
		coin = "莱特币"
	}
	return
}

// 三均线交易策略
func (threelineStrategy *ThreelineStrategy) Buy(price float64) {
	price = price + StringToFloat(Option["slippage"])
	threelineStrategy.BuyPrice = price
	priceStr := fmt.Sprintf("%f", price)
	amount := Option["tradeAmount"]
	buy(priceStr, amount)
	threelineStrategy.PrevTrade = "buy"
	threelineStrategy.SellHalf = false
	t := time.Unix(GetBtTime(), 0)
	logger.Backtestf("%s 在simulate，根据策略%s周期%s，以价格%f买入%s个%s\n",
		t.Format("2006-01-02 15:04:05"), Option["strategy"], Option["tick_interval"],
		price, amount, coinName())
}

func (threelineStrategy *ThreelineStrategy) Sell(price float64) {
	price = price - StringToFloat(Option["slippage"])
	threelineStrategy.BuyPrice = 0
	priceStr := fmt.Sprintf("%f", price)
	amount := Option["tradeAmount"]
	sell(priceStr, amount)
	threelineStrategy.PrevTrade = "sell"
	t := time.Unix(GetBtTime(), 0)
	logger.Backtestf("%s 在simulate，根据策略%s周期%s，以价格%f卖出%s个%s\n",
		t.Format("2006-01-02 15:04:05"), Option["strategy"], Option["tick_interval"],
		price, amount, coinName())
}

func (threelineStrategy *ThreelineStrategy) Tick(records []Record) bool {
	var Price []float64
	for _, v := range records {
		Price = append(Price, v.Close)
	}

	length := len(records)
	ema4 := EMA(Price, 13)
	ema9 := EMA(Price, 34)
	ema18 := EMA(Price, 55)

	if records[length-1].Close > ema4[length-1] &&
		ema4[length-1] > ema9[length-1] &&
		ema9[length-1] > ema18[length-1] {
		if threelineStrategy.PrevTrade != "buy" {
			threelineStrategy.Buy(records[length-1].Close)
		}
	}

	if records[length-1].Close < ema4[length-1] &&
		ema4[length-1] < ema9[length-1] &&
		ema9[length-1] < ema18[length-1] {
		if threelineStrategy.PrevTrade != "sell" {
			threelineStrategy.Sell(records[length-1].Close)
		}
	}

	return true
}
