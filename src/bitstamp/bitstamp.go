package bitstamp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	API_URL = "https://www.bitstamp.net/api/"
)

type Ticker struct {
	Last float64 `json:",string"`
	High float64 `json:",string"`
	Low  float64 `json:",string"`
	Ask  float64 `json:",string"`
	Bid  float64 `json:",string"`
}

type OrderBook struct {
	Time   time.Time
	Orders []Order
}

type Order struct {
	Type   string // "ask" or "bid"
	Price  float64
	Amount float64
}

type Trade struct {
	Time   time.Time
	Id     string
	Price  float64
	Amount float64
}

type Api struct {
	User     string
	Password string
}

// Creates a new api object given a config file. The config file must
// be json formated to inlude User and Password
func NewFromConfig(cfgfile string) (api *Api, err error) {
	api = new(Api)
	return api, nil

}

func New(user, password string) (api *Api, err error) {

	api = &Api{
		User:     user,
		Password: password,
	}
	return api, nil
}

func (api *Api) get(url string) (body []byte, err error) {
	resp, err := http.Get(fmt.Sprint(API_URL, url))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	return
}

func (api *Api) GetTicker() (ticker *Ticker, err error) {
	body, err := api.get("ticker/")
	if err != nil {
		return
	}
	ticker = new(Ticker)
	err = json.Unmarshal(body, ticker)
	if err != nil {
		return
	}
	return

}

func (api *Api) GetOrderBook() (orderbook *OrderBook, err error) {
	body, err := api.get("order_book/")
	if err != nil {
		return
	}
	orderbook = new(OrderBook)
	defaultstruct := make(map[string]interface{})
	err = json.Unmarshal(body, &defaultstruct)
	if err != nil {
		return
	}
	timestamp_str := defaultstruct["timestamp"].(string)
	timestamp, err := strconv.ParseInt(timestamp_str, 10, 64)
	if err != nil {
		return
	}
	orderbook.Time = time.Unix(timestamp, 0)

	bids := defaultstruct["bids"].([]interface{})
	asks := defaultstruct["asks"].([]interface{})
	total := len(bids) + len(asks)
	orderbook.Orders = make([]Order, total)
	for i, bid := range bids {
		_bid := bid.([]interface{})
		price, err := strconv.ParseFloat(_bid[0].(string), 64)
		if err != nil {
			return orderbook, err
		}
		amount, err := strconv.ParseFloat(_bid[1].(string), 64)
		if err != nil {
			return orderbook, err
		}
		order := Order{
			Type:   "bid",
			Price:  price,
			Amount: amount,
		}
		orderbook.Orders[i] = order
	}
	for i, ask := range asks {
		_ask := ask.([]interface{})
		price, err := strconv.ParseFloat(_ask[0].(string), 64)
		if err != nil {
			return orderbook, err
		}
		amount, err := strconv.ParseFloat(_ask[1].(string), 64)
		if err != nil {
			return orderbook, err
		}
		order := Order{
			Type:   "ask",
			Price:  price,
			Amount: amount,
		}
		orderbook.Orders[i+len(bids)] = order
	}

	return
}

// Get the list of last trades with default parameters
func (api *Api) GetTrades() (trades []Trade, err error) {
	body, err := api.get("transactions/")
	return formatTrades(body)
}

// Specify wich trades you want.
func (api *Api) GetTradesParams(offset, limit int64, sort string) (trades []Trade, err error) {
	values := url.Values{}
	values.Add("offset", strconv.FormatInt(offset, 10))
	values.Add("limit", strconv.FormatInt(limit, 10))
	values.Add("sort", sort)
	body, err := api.get("transactions/?" + values.Encode())
	if err != nil {
		return
	}
	return formatTrades(body)
}

func formatTrades(body []byte) (trades []Trade, err error) {
	var defaultstruct []interface{}
	err = json.Unmarshal(body, &defaultstruct)
	if err != nil {
		return
	}
	trades = make([]Trade, len(defaultstruct))
	for i, _trade := range defaultstruct {
		_t := _trade.(map[string]interface{})
		price, err := strconv.ParseFloat(_t["price"].(string), 64)
		if err != nil {
			return trades, err
		}
		amount, err := strconv.ParseFloat(_t["amount"].(string), 64)
		if err != nil {
			return trades, err
		}
		timestamp_str := _t["date"].(string)
		timestamp, err := strconv.ParseInt(timestamp_str, 10, 64)
		if err != nil {
			return trades, err
		}
		time := time.Unix(timestamp, 0)
		trade := Trade{
			Time:   time,
			Id:     strconv.FormatInt(int64(_t["tid"].(float64)), 10),
			Price:  price,
			Amount: amount,
		}
		trades[i] = trade
	}
	return
}
