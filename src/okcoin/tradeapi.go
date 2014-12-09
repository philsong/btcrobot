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

package okcoin

import (
	. "common"
	. "config"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"logger"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"util"
)

/*
	SEE DOC:
	TRADE API
	https://www.okcoin.com/t-1000097.html

	行情API
	https://www.okcoin.com/shequ/themeview.do?tid=1000052&currentPage=1

	// non-official API :P
	K线数据step单位为second
	https://www.okcoin.com/kline/period.do?step=60&symbol=okcoinbtccny&nonce=1394955131098

	https://www.okcoin.com/kline/trades.do?since=10625682&symbol=okcoinbtccny&nonce=1394955760557

	https://www.okcoin.com/kline/depth.do?symbol=okcoinbtccny&nonce=1394955767484

	https://www.okcoin.com/real/ticker.do?symbol=0&random=61

	// old kline for btc
	日数据
	https://www.okcoin.com/klineData.do?type=3&marketFrom=0
	5分钟数据
	https://www.okcoin.com/klineData.do?type=1&marketFrom=0

	// for ltc
	https://www.okcoin.com/klineData.do?type=3&marketFrom=3
*/

type OkcoinTrade struct {
	partner    string
	secret_key string
	errno      int64
}

func NewOkcoinTrade(partner, secret_key string) *OkcoinTrade {
	w := new(OkcoinTrade)
	w.partner = partner
	w.secret_key = secret_key
	return w
}

func (w *OkcoinTrade) createSign(pParams map[string]string) string {
	ms := util.NewMapSorter(pParams)
	sort.Sort(ms)

	v := url.Values{}
	for _, item := range ms {
		v.Add(item.Key, item.Val)
	}

	h := md5.New()

	io.WriteString(h, v.Encode()+w.secret_key)
	sign := fmt.Sprintf("%X", h.Sum(nil))

	return sign
}

func (w *OkcoinTrade) httpRequest(api_url string, pParams map[string]string) (string, error) {
	pParams["sign"] = w.createSign(pParams)

	v := url.Values{}
	for key, val := range pParams {
		v.Add(key, val)
	}

	req, err := http.NewRequest("POST", api_url, strings.NewReader(v.Encode()))
	if err != nil {
		logger.Fatal(err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "https://www.okcoin.cn/")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36")
	logger.Traceln(req)

	c := util.NewTimeoutClient()

	logger.Tracef("HTTP req begin OkcoinTrade")
	resp, err := c.Do(req)
	logger.Tracef("HTTP req end OkcoinTrade")
	if err != nil {
		logger.Fatal(err)
		return "", err
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
			ioutil.WriteFile("cache/okapi_url.json", bodyByte, 0644)
		}

		return body, nil

	} else {
		logger.Tracef("resp %v", resp)
	}

	return "", nil
}

type ErrorMsg struct {
	Result    bool
	ErrorCode int
}

func (w *OkcoinTrade) check_json_result(body string) (errorMsg ErrorMsg, ret bool) {
	if strings.Contains(body, "result") != true {
		ret = false
		return
	}

	doc := json.NewDecoder(strings.NewReader(body))

	if err := doc.Decode(&errorMsg); err == io.EOF {
		logger.Traceln(err)
	} else if err != nil {
		logger.Fatal(err)
		logger.Fatalln(body)
	}

	if errorMsg.Result != true {
		logger.Errorln(errorMsg)
		SetLastError(errorMsg.ErrorCode)
		ret = false
		return
	}

	ret = true
	return
}

/*
　｛
　　　 　 "result":true,
　　　 　 "info":{
　　　 　　 "funds":{
　　　　　 　"free":{
　　　　　 　 　"cny":1000,
　　　 　　 　　 "btc":10,
　　　 　　 　　"ltc":0
　　　　　　　 },
　　　　　　　 "freezed":{
　　　　　　　　 "cny":1000,
　　　　 　　　　 "btc":10,
　　　　 　 　　　"ltc":0
　　　　　　　 }
　　　 　　 }
　　　　 }
　　｝
*/

type Asset struct {
	Net   string
	Total string
}
type UnionFund struct {
	BTC string
	LTC string
}

type Money struct {
	BTC string
	CNY string
	LTC string
}

type Funds struct {
	Free      Money
	Freezed   Money
	Borrow    Money
	Asset     Asset
	UnionFund UnionFund
}

type Info struct {
	Funds Funds
}

type UserInfo struct {
	Result bool
	Info   Info
}

func (w *OkcoinTrade) GetAccount() (userInfo UserInfo, ret bool) {
	pParams := make(map[string]string)
	pParams["partner"] = w.partner

	ret = true

	body, err := w.httpRequest(Config["ok_api_userinfo"], pParams)
	if err != nil {
		ret = false
		return
	}

	_, ret = w.check_json_result(body)
	if ret == false {
		return
	}

	logger.Traceln(body)
	doc := json.NewDecoder(strings.NewReader(body))

	if err := doc.Decode(&userInfo); err == io.EOF {
		ret = false
		logger.Traceln(err)
	} else if err != nil {
		ret = false
		logger.Fatal(err)
	}

	logger.Traceln(userInfo)

	return
}

type OKOrder struct {
	Orders_id   int
	Status      int
	Symbol      string
	Type        string
	Rate        float64
	Amount      float64
	Deal_amount float64
	Avg_rate    float64
}

type OKOrderTable struct {
	Result bool
	Orders []OKOrder
}

func (w *OkcoinTrade) Get_order(symbol, order_id string) (ret bool, m OKOrderTable) {
	pParams := make(map[string]string)
	pParams["partner"] = w.partner
	pParams["symbol"] = symbol
	pParams["order_id"] = order_id

	ret = true

	body, err := w.httpRequest(Config["ok_api_getorder"], pParams)
	if err != nil {
		ret = false
		return
	}

	_, ret = w.check_json_result(body)
	if ret == false {
		return
	}

	doc := json.NewDecoder(strings.NewReader(body))

	if err := doc.Decode(&m); err == io.EOF {
		logger.Traceln(err)
	} else if err != nil {
		logger.Errorln(err)
		logger.Errorln(body)
		logger.Errorln(m)
	}

	return
}

func (w *OkcoinTrade) Get_BTCorder(order_id string) (ret bool, m OKOrderTable) {
	return w.Get_order("btc_cny", order_id)
}

func (w *OkcoinTrade) Get_LTCorder(order_id string) (ret bool, m OKOrderTable) {
	return w.Get_order("ltc_cny", order_id)
}

func (w *OkcoinTrade) Cancel_order(symbol, order_id string) bool {
	pParams := make(map[string]string)
	pParams["partner"] = w.partner
	pParams["symbol"] = symbol
	pParams["order_id"] = order_id

	body, err := w.httpRequest(Config["ok_api_cancelorder"], pParams)
	if err != nil {
		return false
	}
	_, ret := w.check_json_result(body)
	if ret == false {
		return false
	}

	doc := json.NewDecoder(strings.NewReader(body))

	type Msg struct {
		Result   bool
		Order_id int
	}

	var m Msg
	if err := doc.Decode(&m); err == io.EOF {
		logger.Traceln(err)
	} else if err != nil {
		logger.Fatal(err)
	}

	logger.Traceln(m)

	if m.Result == true {
		return true
	} else {
		logger.Infoln(m)

		return false
	}
}

func (w *OkcoinTrade) Cancel_BTCorder(order_id string) (ret bool) {
	return w.Cancel_order("btc_cny", order_id)
}

func (w *OkcoinTrade) Cancel_LTCorder(order_id string) (ret bool) {
	return w.Cancel_order("ltc_cny", order_id)
}

func (w *OkcoinTrade) doTrade(symbol, method, rate, amount string) int {
	pParams := make(map[string]string)
	pParams["partner"] = w.partner
	pParams["symbol"] = symbol
	pParams["type"] = method
	pParams["rate"] = rate
	pParams["amount"] = amount

	body, err := w.httpRequest(Config["ok_api_trade"], pParams)
	if err != nil {
		return 0
	}
	_, ret := w.check_json_result(body)
	if ret == false {
		return 0
	}

	doc := json.NewDecoder(strings.NewReader(body))

	type Msg struct {
		Result   bool
		Order_id int
	}

	var m Msg
	if err := doc.Decode(&m); err == io.EOF {
		logger.Traceln(err)
	} else if err != nil {
		logger.Fatal(err)
	}

	logger.Traceln(m)

	if m.Result == true {
		return m.Order_id
	} else {
		return 0
	}
}

func (w *OkcoinTrade) BuyBTC(price, amount string) string {
	buyId := w.doTrade("btc_cny", "buy", price, amount)
	return (fmt.Sprintf("%d", buyId))
}

func (w *OkcoinTrade) SellBTC(price, amount string) string {
	sellId := w.doTrade("btc_cny", "sell", price, amount)
	return (fmt.Sprintf("%d", sellId))
}

func (w *OkcoinTrade) BuyLTC(price, amount string) string {
	buyId := w.doTrade("ltc_cny", "buy", price, amount)
	return (fmt.Sprintf("%d", buyId))
}

func (w *OkcoinTrade) SellLTC(price, amount string) string {
	sellId := w.doTrade("ltc_cny", "sell", price, amount)
	return (fmt.Sprintf("%d", sellId))
}

func (w *OkcoinTrade) BuyMarketBTC(price, amount string) string {
	buyId := w.doTrade("btc_cny", "buy_market", price, amount)
	return (fmt.Sprintf("%d", buyId))
}

func (w *OkcoinTrade) SellMarketBTC(price, amount string) string {
	sellId := w.doTrade("btc_cny", "sell_market", price, amount)
	return (fmt.Sprintf("%d", sellId))
}

func (w *OkcoinTrade) BuyMarketLTC(price, amount string) string {
	buyId := w.doTrade("ltc_cny", "buy_market", price, amount)
	return (fmt.Sprintf("%d", buyId))
}

func (w *OkcoinTrade) SellMarketLTC(price, amount string) string {
	sellId := w.doTrade("ltc_cny", "sell_market", price, amount)
	return (fmt.Sprintf("%d", sellId))
}
