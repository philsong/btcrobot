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

package huobi

import (
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
	"strconv"
	"strings"
	"time"
	"util"
)

/*
	http://www.huobi.com/help/index.php?a=api_help
*/
type HuobiTrade struct {
	access_key string
	secret_key string
}

func NewHuobiTrade(access_key, secret_key string) *HuobiTrade {
	w := new(HuobiTrade)
	w.access_key = access_key
	w.secret_key = secret_key
	return w
}

func (w *HuobiTrade) createSign(pParams map[string]string) string {
	pParams["secret_key"] = w.secret_key

	ms := util.NewMapSorter(pParams)
	sort.Sort(ms)

	v := url.Values{}
	for _, item := range ms {
		v.Add(item.Key, item.Val)
	}

	h := md5.New()

	io.WriteString(h, v.Encode())
	sign := fmt.Sprintf("%x", h.Sum(nil))

	return sign
}

func (w *HuobiTrade) httpRequest(pParams map[string]string) (string, error) {
	v := url.Values{}
	for key, val := range pParams {
		v.Add(key, val)
	}

	req, err := http.NewRequest("POST", Config["hb_api_url"], strings.NewReader(v.Encode()))
	if err != nil {
		logger.Fatal(err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", Config["hb_base_url"])
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36")
	logger.Traceln(req)

	c := util.NewTimeoutClient()

	logger.Tracef("HTTP req begin HuobiTrade")
	resp, err := c.Do(req)
	logger.Tracef("HTTP req end HuobiTrade")
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
			ioutil.WriteFile("cache/api_url.json", bodyByte, 0644)
		}

		logger.Traceln(body)

		return body, nil

	} else {
		logger.Tracef("resp %v", resp)
	}

	return "", nil
}

type ErrorMsg struct {
	Code int
	Msg  string
	Time int
}

type HBOrderItem struct {
	Id               int
	Type             int
	order_price      float64
	order_amount     float64
	processed_amount float64
	order_time       int
}

type HBOrder struct {
	Id               int
	Type             int
	Order_price      string
	Order_amount     string
	Processed_amount string
	Processed_price  string
	Total            string
	Fee              string
	Vot              string
	Status           int
}

func (w *HuobiTrade) check_json_result(body string) (errorMsg ErrorMsg, ret bool) {
	if strings.Contains(body, "code") != true {
		ret = true
		return
	}

	doc := json.NewDecoder(strings.NewReader(body))

	if err := doc.Decode(&errorMsg); err == io.EOF {
		logger.Traceln(err)
	} else if err != nil {
		logger.Fatal(err)
	}

	if errorMsg.Code != 0 {
		logger.Errorln(errorMsg)
		ret = false
		return
	}

	ret = true
	return
}

type Account_info struct {
	Total                 string
	Net_asset             string
	Available_cny_display string
	Available_btc_display string
	Frozen_cny_display    string
	Frozen_btc_display    string
	Loan_cny_display      string
	Loan_btc_display      string
}

func (w *HuobiTrade) GetAccount() (account_info Account_info, ret bool) {
	pParams := make(map[string]string)
	pParams["method"] = "get_account_info"
	pParams["access_key"] = w.access_key
	now := time.Now().Unix()
	pParams["created"] = strconv.FormatInt(now, 10)
	pParams["sign"] = w.createSign(pParams)

	ret = true

	body, err := w.httpRequest(pParams)
	if err != nil {
		ret = false
		return
	}

	_, ret = w.check_json_result(body)
	if ret == false {
		return
	}

	doc := json.NewDecoder(strings.NewReader(body))

	if err := doc.Decode(&account_info); err == io.EOF {
		logger.Fatal(err)
	} else if err != nil {
		logger.Fatal(err)
	}

	logger.Traceln(account_info)

	return
}

func (w *HuobiTrade) Get_orders() (ret bool, m []HBOrderItem) {
	pParams := make(map[string]string)
	pParams["method"] = "get_orders"
	pParams["access_key"] = w.access_key
	now := time.Now().Unix()
	pParams["created"] = strconv.FormatInt(now, 10)
	pParams["sign"] = w.createSign(pParams)

	ret = true

	body, err := w.httpRequest(pParams)
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
		logger.Fatal(err)
	}

	logger.Traceln(m)

	return
}

func (w *HuobiTrade) Get_order(id string) (ret bool, m HBOrder) {
	pParams := make(map[string]string)
	pParams["method"] = "order_info"
	pParams["access_key"] = w.access_key
	pParams["id"] = id
	now := time.Now().Unix()
	pParams["created"] = strconv.FormatInt(now, 10)
	pParams["sign"] = w.createSign(pParams)

	ret = true

	body, err := w.httpRequest(pParams)
	if err != nil {
		ret = false
		logger.Infoln(err)
		return
	}

	_, ret = w.check_json_result(body)
	if ret == false {
		logger.Infoln(body)
		return
	}

	doc := json.NewDecoder(strings.NewReader(body))
	if err := doc.Decode(&m); err == io.EOF {
		ret = false
		logger.Infoln(err)
	} else if err != nil {
		ret = false
		logger.Infoln(err)
	}

	return
}

func (w *HuobiTrade) Cancel_order(id string) bool {
	pParams := make(map[string]string)
	pParams["method"] = "cancel_order"
	pParams["access_key"] = w.access_key
	pParams["id"] = id
	now := time.Now().Unix()
	pParams["created"] = strconv.FormatInt(now, 10)
	pParams["sign"] = w.createSign(pParams)

	body, err := w.httpRequest(pParams)
	if err != nil {
		return false
	}
	_, ret := w.check_json_result(body)
	if ret == false {
		return false
	}

	doc := json.NewDecoder(strings.NewReader(body))

	type Msg struct {
		Result string
	}

	var m Msg
	if err := doc.Decode(&m); err == io.EOF {
		logger.Traceln(err)
	} else if err != nil {
		logger.Fatal(err)
	}

	logger.Traceln(m)

	if m.Result == "success" {
		return true
	} else {
		return false
	}
}

func (w *HuobiTrade) doTrade(method, price, amount string) int {
	pParams := make(map[string]string)
	pParams["method"] = method
	pParams["access_key"] = w.access_key
	pParams["price"] = price
	pParams["amount"] = amount
	now := time.Now().Unix()
	pParams["created"] = strconv.FormatInt(now, 10)
	pParams["sign"] = w.createSign(pParams)

	body, err := w.httpRequest(pParams)
	if err != nil {
		return 0
	}
	_, ret := w.check_json_result(body)
	if ret == false {
		return 0
	}

	doc := json.NewDecoder(strings.NewReader(body))

	type Msg struct {
		Result string
		Id     int
	}

	var m Msg
	if err := doc.Decode(&m); err == io.EOF {
		logger.Traceln(err)
	} else if err != nil {
		logger.Fatal(err)
	}

	logger.Infoln(m)

	if m.Result == "success" {
		return m.Id
	} else {
		return 0
	}
}

func (w *HuobiTrade) BuyBTC(price, amount string) string {
	buyId := w.doTrade("buy", price, amount)
	return (fmt.Sprintf("%d", buyId))
}

func (w *HuobiTrade) SellBTC(price, amount string) string {
	sellId := w.doTrade("sell", price, amount)
	return (fmt.Sprintf("%d", sellId))
}

func (w *HuobiTrade) BuyLTC(price, amount string) string {
	// todo
	buyId := 0
	return (fmt.Sprintf("%d", buyId))
}

func (w *HuobiTrade) SellLTC(price, amount string) string {
	// todo
	sellId := 0
	return (fmt.Sprintf("%d", sellId))
}
