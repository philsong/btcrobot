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

package Bittrex

import (
	. "common"
	. "config"
	"fmt"
	"github.com/philsong/go-bittrex"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Bittrex struct {
	init      bool
	records   []Record
	LastEpoch string
}

var (
	once   sync.Once
	manage *Bittrex
)

func Manager() (m *Bittrex) {
	if manage == nil {
		once.Do(func() {
			manage = new(Bittrex)
		})
	}
	m = manage
	return
}

func NewBittrex() *Bittrex {
	w := new(Bittrex)
	return w
}

func (w Bittrex) CancelOrder(order_id string) (ret bool) {
	return
}

func (w Bittrex) GetOrderBook() (ret bool, orderBook OrderBook) {
	return
}

func (w Bittrex) GetOrder(order_id string) (ret bool, order Order) {
	return
}

func (w *Bittrex) GetKLine(peroid int) (ret bool, records []Record) {
	//symbol := Option["symbol"]

	// Bittrex client
	bittrex := bittrex.New(SecretOption["bittrex_access_key"], SecretOption["bittrex_secret_key"])

	fmt.Println(w.init, len(w.records), w.LastEpoch)
	if w.init == false {
		candles, err := bittrex.GetHisCandles("BTC-BC")
		if err != nil {
			ret = false
			fmt.Println(err)
			return
		} else {
			ret = true
		}
		fmt.Println(err, len(candles))

		var record Record
		for _, candle := range candles {
			TimeStr := strings.TrimPrefix(candle.TimeStamp, "/Date(")
			TimeStr = strings.TrimSuffix(TimeStr, "000)/")
			//secs, _ := strconv.Atoi(TimeStr)
			//fmt.Println(time.Unix(int64(secs), 0))
			//fmt.Println(TimeStr, candle.TimeStamp, candle.Open, candle.Close, candle.High, candle.Low, candle.Volume)
			//return

			record.TimeStr = TimeStr
			record.Open = candle.Open
			record.Close = candle.Close
			record.High = candle.High
			record.Low = candle.Low
			record.Volumn = candle.Volume
			records = append(records, record)
			//fmt.Println(time.Unix(int64(secs), 0))
			//fmt.Println(TimeStr, candle.Open, candle.Close, candle.High, candle.Low, candle.Volume)
		}

		secs, _ := strconv.Atoi(records[len(records)-1].TimeStr)
		fmt.Println(time.Unix(int64(secs), 0))

		if w.init == false && len(records) > 0 {
			w.init = true
			w.records = records
			w.LastEpoch = records[len(records)-1].TimeStr + "000"
			fmt.Println(w.LastEpoch)
		}

	} else {
		records = w.records
		candles, err := bittrex.GetNewCandles("BTC-BC", w.LastEpoch)
		if err != nil {
			ret = false
			fmt.Println(err)
			return
		} else {
			ret = true
		}

		fmt.Println(err, len(candles))

		var record Record
		for _, candle := range candles {
			TimeStr := strings.TrimPrefix(candle.TimeStamp, "/Date(")
			TimeStr = strings.TrimSuffix(TimeStr, "000)/")
			secs, _ := strconv.Atoi(TimeStr)
			//fmt.Println(time.Unix(int64(secs), 0))
			//fmt.Println(TimeStr, candle.TimeStamp, candle.Open, candle.Close, candle.High, candle.Low, candle.Volume)
			//return

			record.TimeStr = TimeStr
			record.Open = candle.Open
			record.Close = candle.Close
			record.High = candle.High
			record.Low = candle.Low
			record.Volumn = candle.Volume
			records = append(records, record)
			fmt.Println(time.Unix(int64(secs), 0))
			//fmt.Println(TimeStr, candle.Open, candle.Close, candle.High, candle.Low, candle.Volume)
		}

		if len(records) > 0 {
			w.records = records
			w.LastEpoch = records[len(records)-1].TimeStr + "000"
		}
		secs, _ := strconv.Atoi(records[len(records)-1].TimeStr)
		fmt.Println(len(records), time.Unix(int64(secs), 0))
	}

	fmt.Println(time.Now())

	//fmt.Println(records)
	return
}

func (w Bittrex) GetAccount() (account Account, ret bool) {
	return
}

func (w Bittrex) Buy(tradePrice, tradeAmount string) (buyId string) {
	return
}

func (w Bittrex) Sell(tradePrice, tradeAmount string) (sellId string) {
	return
}
