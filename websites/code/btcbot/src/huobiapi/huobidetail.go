/*
  btcbot is a Bitcoin trading bot for HUOBI.com written
  in golang, it features multiple trading methods using
  technical analysis.

  Disclaimer:

  USE AT YOUR OWN RISK!

  The author of this project is NOT responsible for any damage or loss caused
  by this software. There can be bugs and the bot may not perform as expected
  or specified. Please consider testing it first with paper trading /
  backtesting on historical data. Also look at the code to see what how
  it's working.

  Weibo:http://weibo.com/bocaicfa
*/

package huobiapi

import (
	"encoding/json"
	"io/ioutil"
	"logger"
	"os"
	"strings"
)

type SellBuy struct {
	price  float64
	level  float64 ////涨幅
	amount float64 //成交量
}

type Trade struct {
	time   string
	price  float64
	amount float64
	type_  string
}

type Top_buy_sell struct {
	Price  float64
	level  float64
	amount float64
	accu   float64
}

type View_detail struct {
	sells      [10]SellBuy
	buys       [10]SellBuy
	trades     [15]Trade
	Vp_new     float64
	level      float64
	amount     float64
	total      float64
	amp        float64
	Vp_open    float64
	Vp_high    float64
	Vp_low     float64
	Vp_last    float64
	Vtop_sells [5]Top_buy_sell
	Vtop_buys  [5]Top_buy_sell
}

func (w *Huobi) TradeDetailAnalyze(body string) bool {
	logger.Traceln("TradeDetailAnalyze start....")

	logger.Traceln("TradeDetailAnalyze body....")
	logger.Traceln(body)
	logger.Traceln("---------------------------")

	body = strings.TrimPrefix(body, "view_detail(")
	body = strings.TrimSuffix(body, ")")

	logger.Traceln("TradeDetailAnalyze json....")
	logger.Traceln(body)
	ioutil.WriteFile("cache/TradeDetailAnalyze.json", []byte(body), os.ModeAppend)
	logger.Traceln("---------------------------")

	/*
		var view_detail View_detail
		if err := json.Unmarshal([]byte(body), &view_detail); err != nil {
			logger.Traceln("TradeDetailAnalyze json....panic!!!")
			logger.Traceln(body)
			logger.Traceln("---------------------------panic!!!")
			panic(err)
		}

		fmt.Println(view_detail)
		return true
	*/

	detail_data := &w.Detail_data
	var view_detail map[string]interface{}

	if err := json.Unmarshal([]byte(body), &view_detail); err != nil {
		logger.Traceln("TradeDetailAnalyze json....panic!!!")
		logger.Traceln(body)
		logger.Traceln("---------------------------panic!!!")
		return false
	}

	p_new := view_detail["p_new"].(float64)
	level := view_detail["level"].(float64)
	amount := view_detail["amount"].(float64)
	total := view_detail["total"].(float64)
	amp := view_detail["amp"].(float64)
	p_open := view_detail["p_open"].(float64)
	p_high := view_detail["p_high"].(float64)
	p_low := view_detail["p_low"].(float64)
	p_last := view_detail["p_last"].(float64)

	detail_data.Vp_new = p_new
	detail_data.level = level
	detail_data.amount = amount
	detail_data.total = total
	detail_data.amp = amp
	detail_data.Vp_open = p_open
	detail_data.Vp_high = p_high
	detail_data.Vp_low = p_low
	detail_data.Vp_last = p_last

	sells := view_detail["sells"].([]interface{})
	parse_buy_sell(sells, &detail_data.sells)

	buys := view_detail["buys"].([]interface{})
	parse_buy_sell(buys, &detail_data.buys)

	trades := view_detail["trades"].([]interface{})
	parse_trade(trades, &detail_data.trades)

	top_buys := view_detail["top_buy"].([]interface{})
	parse_topbuy(top_buys, &detail_data.Vtop_buys)

	top_sells := view_detail["top_sell"].(map[string]interface{})
	parse_topsell(top_sells, &detail_data.Vtop_sells)

	logger.Traceln(detail_data)
	return true
	/*

		doc := json.NewDecoder(strings.NewReader(body))

		//var view_detail View_detail

		if err := doc.Decode(&view_detail); err == io.EOF {
			logger.Traceln(err)
		} else if err != nil {
			logger.Fatal(err)
		}

		logger.Infoln(view_detail)

		logger.Traceln("TradeDetailAnalyze end-----")
		return true
	*/
}

func parse_trade(trades []interface{}, trades_data *[15]Trade) {
	for k, v := range trades {
		switch vt := v.(type) {
		case map[string]interface{}:
			logger.Traceln(k, " is a map:")
			logger.Tracef("trades[%d]\n", k)

			for ik, iv := range vt {
				switch ik {
				case "time":
					trades_data[k].time = iv.(string)
				case "price":
					trades_data[k].price = InterfaceToFloat64(iv)
				case "amount":
					trades_data[k].amount = InterfaceToFloat64(iv)
				case "type":
					trades_data[k].type_ = iv.(string)
				}
			}
		default:
			logger.Errorln(k, v)
			logger.Fatalln("don't know the type, crash!")
		}
	}
}

func parse_buy_sell(sells_buys []interface{}, sells_buys_data *[10]SellBuy) {
	for k, v := range sells_buys {
		switch vt := v.(type) {
		case map[string]interface{}:
			logger.Traceln(k, " is a map:")
			logger.Tracef("sells/buys[%d]\n", k)

			for ik, iv := range vt {
				switch ik {
				case "price":
					sells_buys_data[k].price = InterfaceToFloat64(iv)
				case "level":
					sells_buys_data[k].level = InterfaceToFloat64(iv)
				case "amount":
					sells_buys_data[k].amount = InterfaceToFloat64(iv)
				}
			}
		default:
			logger.Errorln(k, v)
			logger.Fatalln("don't know the type, crash!")
		}
	}
}

func parse_topbuy(topbuys []interface{}, topbuys_data *[5]Top_buy_sell) {
	for k, v := range topbuys {
		switch vt := v.(type) {
		case map[string]interface{}:
			logger.Traceln(k, " is a map:")
			logger.Tracef("topbuys[%d]\n", k)

			for ik, iv := range vt {
				switch ik {
				case "price":
					topbuys_data[k].Price = InterfaceToFloat64(iv)
				case "amount":
					topbuys_data[k].amount = InterfaceToFloat64(iv)
				case "level":
					topbuys_data[k].level = InterfaceToFloat64(iv)
				case "accu":
					topbuys_data[k].accu = InterfaceToFloat64(iv)
				}
			}
		default:
			logger.Errorln(k, v)
			logger.Fatalln("don't know the type, crash!")
		}
	}
}

func parse_topsell(topsells map[string]interface{}, topsells_data *[5]Top_buy_sell) {
	index := 4
	for k, v := range topsells {
		switch vt := v.(type) {
		case map[string]interface{}:
			logger.Traceln(k, " is a map:")
			logger.Tracef("topsells[%s]\n", k)

			for ik, iv := range vt {
				logger.Traceln(ik, iv)
				switch ik {
				case "price":
					topsells_data[index].Price = InterfaceToFloat64(iv)
				case "amount":
					topsells_data[index].amount = InterfaceToFloat64(iv)
				case "level":
					topsells_data[index].level = InterfaceToFloat64(iv)
				case "accu":
					topsells_data[index].accu = InterfaceToFloat64(iv)
				}
			}

			index--

		default:
			logger.Errorln(k, v)
			logger.Fatalln("don't know the type, crash!")
		}
	}
}
