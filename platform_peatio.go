package platform

// https://peatio.com/documents/api_v2 http://demo.peat.io/documents/websocket_api
import (
	"crypto/sha256"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type iPeatio struct {
	accessKey string
	secretKey string
	currency  string
	symbol    string
	step      int64
	timeout   time.Duration
}

func newPeatio(accessKey, secretKey, currency string, peroid int64) (IExchange, error) {
	currency = strings.ToLower(currency)
	if currency != "btc" && currency != "pts" && currency != "dog" {
		return nil, errors.New("Currency not support " + currency)
	}

	step, ok := PeriodsMap[peroid]
	if !ok {
		return nil, errors.New("Peroid not support")
	}

	s := new(iPeatio)
	s.accessKey = accessKey
	s.secretKey = secretKey
	s.currency = currency
	s.step = step
	s.symbol = currency + "cny"
	s.timeout = 20 * time.Second
	return s, nil
}

func (p *iPeatio) GetName() string {
	return "Peatio"
}

func (p *iPeatio) GetRate() (rate float64) {
    return 1.0
}

func (p *iPeatio) GetCurrency() string {
	return strings.ToUpper(p.currency)
}

func (p *iPeatio) GetTicker() (ticker Ticker, err error) {
	js, err := Request{Uri: fmt.Sprintf("https://peatio.com/api/v2/tickers/%s.json", p.symbol), Timeout: p.timeout}.DoJSON()
	if err != nil {
		return
	}
	dic := js.Get("ticker")
	ticker.Buy = dic.Get("buy").MustFloat64()
	ticker.Sell = dic.Get("sell").MustFloat64()
	ticker.Last = dic.Get("last").MustFloat64()
	ticker.High = dic.Get("high").MustFloat64()
	ticker.Low = dic.Get("low").MustFloat64()
	ticker.Volume = dic.Get("vol").MustFloat64()
	return
}

func (p *iPeatio) GetDepth() (depth Depth, err error) {
	// &asks_limit=20&bids_limit=20
	js, err := Request{Uri: fmt.Sprintf("https://peatio.com/api/v2/order_book.json?market=%s", p.symbol), Timeout: p.timeout}.DoJSON()
	if err != nil {
		return
	}

	for _, mp := range js.Get("asks").MustArray() {
		if dic, ok := mp.(map[string]interface{}); ok {
			depth.Asks = append(depth.Asks, MarketOrder{toFloat(dic["price"]), toFloat(dic["volume"])})
		}
	}

	for _, mp := range js.Get("bids").MustArray() {
		if dic, ok := mp.(map[string]interface{}); ok {
			depth.Bids = append(depth.Bids, MarketOrder{toFloat(dic["price"]), toFloat(dic["volume"])})
		}
	}
	return
}

func (p *iPeatio) GetTrades() (trades []Trade, err error) {
    return
}

func (p *iPeatio) GetRecords() (records []Record, err error) {
	panic("Peatio not support GetRecords")
	return
}

func (p *iPeatio) tapiCall(httpMethod, method string, params map[string]string) (js *Json, err error) {
	params["access_key"] = p.accessKey
	params["tonce"] = strconv.FormatInt(now(), 10)
	params["signature"] = HMACEncrypt(sha256.New, httpMethod+"|/api/v2/"+method+".json|"+encodeParams(params), p.secretKey)
	jsonUri := fmt.Sprintf("https://peatio.com/api/v2/%s.json", method)
	req := Request{
		Method:  httpMethod,
		Timeout: p.timeout,
	}
	req.Uri = jsonUri
	qs := encodeParams(params)
	if httpMethod == "POST" {
		req.AddHeader("Accept", "application/json")
		req.AddHeader("Content-Type", "application/x-www-form-urlencoded")
		req.Body = qs
	} else if len(qs) > 0 {
		req.Uri = jsonUri + "?" + qs
	}

	js, err = req.DoJSON()
	if err != nil {
		return
	}

	if obj, ok := js.CheckGet("error"); ok {
		return nil, errors.New(fmt.Sprintf("%+v", obj))
	}
	return
}

func (p *iPeatio) GetAccount() (account Account, err error) {
	js, err := p.tapiCall("GET", "members/me", map[string]string{})
	if err != nil {
		return
	}
	for _, item := range js.Get("accounts").MustArray() {
		mp := item.(map[string]interface{})
		switch toString(mp["currency"]) {
		case p.currency:
			account.FrozenStocks = toFloat(mp["locked"])
			account.Stocks = toFloat(mp["balance"])
		case "cny":
			account.FrozenBalance = toFloat(mp["locked"])
			account.Balance = toFloat(mp["balance"])
		}
	}
	return
}

func (p *iPeatio) __trade(tradeType string, price, amount float64) (int64, error) {
	js, err := p.tapiCall("POST", "orders", map[string]string{
		"market": p.symbol,
		"side":   tradeType,
		"price":  float2str(price),
		"volume": float2str(amount),
	})
	if err != nil {
		return 0, err
	}

	orderId := js.Get("id").MustInt64()
	return orderId, nil
}
func (p *iPeatio) Buy(price, amount float64) (int64, error) {
	return p.__trade("buy", price, amount)
}

func (p *iPeatio) Sell(price, amount float64) (int64, error) {
	return p.__trade("sell", price, amount)
}

func (p *iPeatio) GetOrders() (orders []Order, err error) {
	js, err := p.tapiCall("GET", "orders", map[string]string{
		"state":  "wait",
		"market": p.symbol,
		"limit":  "100",
	})
    if err != nil {
		return nil, err
	}
	for _, item := range js.MustArray() {
		mp := item.(map[string]interface{})
		var order Order
		order.Id = int64(toFloat(mp["id"]))
		order.Amount = toFloat(mp["volume"])
		order.Price = toFloat(mp["price"])
		order.DealAmount = toFloat(mp["executed_volume"])
		if mp["side"].(string) == "buy" {
			order.Type = ORDER_TYPE_BUY
		} else {
			order.Type = ORDER_TYPE_SELL
		}
		if order.DealAmount > 0 {
			order.Status = ORDER_STATE_PARTIAL
		} else {
			order.Status = ORDER_STATE_PENDING
		}
		orders = append(orders, order)
	}
	return
}

func (p *iPeatio) GetOrder(orderId int64) (order Order, err error) {
    var js *Json
	js, err = p.tapiCall("GET", "order", map[string]string{
		"id": strconv.FormatInt(orderId, 10),
	})
    if err != nil {
		return
	}

	mp := js.MustMap()
	order.Id = int64(toFloat(mp["id"]))
	order.Amount = toFloat(mp["volume"])
	order.Price = toFloat(mp["price"])
	order.DealAmount = toFloat(mp["executed_volume"])
	if mp["side"].(string) == "buy" {
		order.Type = ORDER_TYPE_BUY
	} else {
		order.Type = ORDER_TYPE_SELL
	}
	switch mp["side"].(string) {
	case "wait":
		if order.DealAmount > 0 {
			order.Status = ORDER_STATE_PARTIAL
		} else {
			order.Status = ORDER_STATE_PENDING
		}
    case "done":
        order.Status = ORDER_STATE_CLOSED
    case "cancel":
        order.Status = ORDER_STATE_CANCELED
	}
    return
}

func (p *iPeatio) CancelOrder(orderId int64) (ret bool, err error) {
    _, err = p.tapiCall("POST", "order/delete", map[string]string{
        "id": strconv.FormatInt(orderId, 10),
    })
    if err != nil {
        return
    }
    ret = true
    return
}

func (p *iPeatio) GetMinStock() float64 {
	if p.currency == "btc" {
		return 0.01
	}
	return 0.1
}

func (p *iPeatio) GetFee() (fee Fee, err error) {
	fee.Buy = 0.0
	fee.Sell = 0.0
	return
}
