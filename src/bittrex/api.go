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
	bittrex   *bittrex.Bittrex
	market    string
}

var (
	once   sync.Once
	manage *Bittrex
)

func Manager() (m *Bittrex) {
	if manage == nil {
		once.Do(func() {
			manage = new(Bittrex)
			// Bittrex client
			manage.market = "BTC-BOOM"
			manage.bittrex = bittrex.New(SecretOption["bittrex_access_key"], SecretOption["bittrex_secret_key"])
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
	bittrex := w.bittrex
	err := bittrex.CancelOrder(order_id)
	if err != nil {
		ret = false
		return
	}
	return true
}

func (w Bittrex) GetOrderBook() (ret bool, orderBook OrderBook) {
	bittrex := w.bittrex
	bittreOrderBook, err := bittrex.GetOrderBook(w.market, "both", 100)
	if err != nil {
		ret = false
		return
	}
	ret = true

	for i := 0; i < 10; i++ {
		orderBook.Asks[i].Price = bittreOrderBook.Sell[len(bittreOrderBook.Sell)-10+i].Rate
		orderBook.Asks[i].Amount = bittreOrderBook.Sell[len(bittreOrderBook.Sell)-10+i].Quantity
		orderBook.Bids[i].Price = bittreOrderBook.Buy[i].Rate
		orderBook.Bids[i].Amount = bittreOrderBook.Buy[i].Quantity
	}

	return
}

func (w Bittrex) GetOrder(order_id string) (ret bool, order Order) {
	bittrex := w.bittrex
	openOrders, err := bittrex.GetOpenOrders(w.market)
	if err != nil {
		ret = false
		return
	}
	fmt.Println(openOrders)
	ret = true
	return
}

func (w Bittrex) GetAccount() (account Account, ret bool) {
	return
}

func (w Bittrex) Buy(tradePrice, tradeAmount string) (buyId string) {
	bittrex := w.bittrex
	uuid, err := bittrex.BuyLimit(w.market, toFloat(tradeAmount), toFloat(tradePrice))
	if err != nil {
		buyId = "0"
	} else {
		buyId = uuid
	}

	return buyId
}

func (w Bittrex) Sell(tradePrice, tradeAmount string) (sellId string) {
	bittrex := w.bittrex
	uuid, err := bittrex.SellLimit(w.market, toFloat(tradeAmount), toFloat(tradePrice))
	if err != nil {
		sellId = "0"
	} else {
		sellId = uuid
	}

	return sellId
}

func (w *Bittrex) GetKLine(peroid int) (ret bool, records []Record) {
	//symbol := Option["symbol"]
	bittrex := w.bittrex

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
			fmt.Println(TimeStr, candle.Open, candle.Close, candle.High, candle.Low, candle.Volume)
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
			fmt.Println(TimeStr, candle.Open, candle.Close, candle.High, candle.Low, candle.Volume)
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
		ret, _ = strconv.ParseFloat(v, 64)
	}
	return ret
}

func float2str(i float64) string {
	return strconv.FormatFloat(i, 'f', -1, 64)
}
