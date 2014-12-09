package peatio

// todo: 懒得搞了，你们玩呗

// https://peatio.com/documents/api_v2 http://demo.peat.io/documents/websocket_api

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	. "github.com/bitly/go-simplejson"
	"hash"
	"io/ioutil"
	"logger"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"util"
)

const (
	_ = iota // ignore first value by assigning to blank identifier
	ORDER_TYPE_BUY
	ORDER_TYPE_SELL
	ORDER_STATE_PARTIAL
	ORDER_STATE_PENDING
	ORDER_STATE_CLOSED
	ORDER_STATE_CANCELED
)

type iPeatio struct {
	accessKey string
	secretKey string
	currency  string
	symbol    string
	step      int64
	timeout   time.Duration
}

type (
	/*
		Record struct {
			Time   int64
			Open   float64
			High   float64
			Low    float64
			Close  float64
			Volume float64
		}
	*/

	MarketOrder struct {
		Price  float64
		Amount float64
	}

	Fee struct {
		Sell float64
		Buy  float64
	}

	Ticker struct {
		High   float64
		Low    float64
		Sell   float64
		Buy    float64
		Last   float64
		Volume float64
	}

	Trade struct {
		Id     int64
		Time   int64
		Amount float64
		Price  float64
		Type   int64
	}

	POrder struct {
		Id         int64
		Amount     float64
		DealAmount float64
		Price      float64
		Status     int64
		Type       int64
	}

	Depth struct {
		Asks []MarketOrder
		Bids []MarketOrder
	}

	PAccount struct {
		Stocks        float64
		FrozenStocks  float64
		Balance       float64
		FrozenBalance float64
	}
)

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

func newPeatio(accessKey, secretKey, currency string, peroid int64) (*iPeatio, error) {
	currency = strings.ToLower(currency)
	if currency != "btc" && currency != "pts" && currency != "dog" {
		return nil, errors.New("Currency not support " + currency)
	}

	s := new(iPeatio)
	s.accessKey = accessKey
	s.secretKey = secretKey
	s.currency = currency
	s.step = peroid
	s.symbol = currency + "cny"
	s.timeout = 20 * time.Second
	return s, nil
}

type Request struct {
	http.Request
	EncodeParams string
	Method       string
	Uri          string
	Timeout      time.Duration
}

func (p *Request) DoJSON() (*Json, error) {
	j, _ := NewJson([]byte(""))
	var req *http.Request
	var err error
	if p.Method == "POST" {
		req, err = http.NewRequest(p.Method, p.Uri, strings.NewReader(p.EncodeParams))
	} else {
		req, err = http.NewRequest(p.Method, p.Uri+"?"+p.EncodeParams, nil)
	}

	if err != nil {
		logger.Fatal(err)
		return j, err
	}

	if p.Method == "POST" {
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36")
	logger.Traceln(req)

	c := util.NewTimeoutClient()

	logger.Tracef("HTTP req begin HuobiTrade")
	resp, err := c.Do(req)
	logger.Tracef("HTTP req end HuobiTrade")
	if err != nil {
		logger.Fatal(err)
		return j, err
	}
	defer resp.Body.Close()

	logger.Tracef("api_url resp StatusCode=%v", resp.StatusCode)
	logger.Tracef("api_url resp=%v", resp)
	if resp.StatusCode == 200 {
		var body string

		contentEncoding := resp.Header.Get("Content-Encoding")
		logger.Tracef("HTTP returned Content-Encoding %s", contentEncoding)
		logger.Traceln(resp.Header.Get("Content-Type"))

		switch contentEncoding {
		case "gzip":
			body = util.DumpGZIP(resp.Body)

		default:
			bodyByte, _ := ioutil.ReadAll(resp.Body)
			body = string(bodyByte)
			ioutil.WriteFile("cache/api_url.json", bodyByte, 0644)
		}

		logger.Traceln(body)

		return NewJson([]byte(body))

	} else {
		logger.Tracef("resp %v", resp)
	}

	return j, nil
}

// CheckMAC returns true if messageMAC is a valid HMAC tag for message.
func HMACEncrypt(h func() hash.Hash, message, key string) string {
	mac := hmac.New(h, []byte(key))
	mac.Write([]byte(message))
	messageMAC := mac.Sum(nil)
	return string(messageMAC)
}

func encodeParams(params map[string]string) string {
	v := url.Values{}
	for key, val := range params {
		v.Add(key, val)
	}

	return v.Encode()
}

func (p *iPeatio) tapiCall(httpMethod, method string, params map[string]string) (js *Json, err error) {
	params["access_key"] = p.accessKey
	params["tonce"] = strconv.FormatInt(time.Now().Unix(), 10)
	params["signature"] = HMACEncrypt(sha256.New, httpMethod+"|/api/v2/"+method+".json|"+encodeParams(params), p.secretKey)
	jsonUri := fmt.Sprintf("https://peatio.com/api/v2/%s.json", method)
	req := Request{
		Method:  httpMethod,
		Timeout: p.timeout,
	}
	req.Uri = jsonUri
	req.EncodeParams = encodeParams(params)

	js, err = req.DoJSON()
	if err != nil {
		return
	}

	if obj, ok := js.CheckGet("error"); ok {
		return nil, errors.New(fmt.Sprintf("%+v", obj))
	}
	return
}

func (p *iPeatio) GetAccount() (account PAccount, err error) {
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

func (p *iPeatio) GetOrders() (orders []POrder, err error) {
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
		var order POrder
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

func (p *iPeatio) GetOrder(orderId int64) (order POrder, err error) {
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
	switch mp["state"].(string) {
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
