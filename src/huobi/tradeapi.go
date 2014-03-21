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
	client     *http.Client
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

	req, err := http.NewRequest("POST", Config["api_url"], strings.NewReader(v.Encode()))
	if err != nil {
		logger.Fatal(err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "https://www.huobi.com/")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36")
	logger.Traceln(req)

	if w.client == nil {
		w.client = &http.Client{nil, nil, nil}
	}

	resp, err := w.client.Do(req)
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
			body = DumpGZIP(resp.Body)

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

type DelegationsMsg struct {
	Id               int
	Type             int
	order_price      string
	order_amount     string
	processed_amount string
	order_time       int
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

func (w *HuobiTrade) Get_account_info() (string, error) {
	pParams := make(map[string]string)
	pParams["method"] = "get_account_info"
	pParams["access_key"] = w.access_key
	now := time.Now().Unix()
	pParams["created"] = strconv.FormatInt(now, 10)
	pParams["sign"] = w.createSign(pParams)

	return w.httpRequest(pParams)
}

func (w *HuobiTrade) Get_delegations() (m []DelegationsMsg, ret bool) {
	pParams := make(map[string]string)
	pParams["method"] = "get_delegations"
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

func (w *HuobiTrade) Get_delegation_info(id string) (string, error) {
	pParams := make(map[string]string)
	pParams["method"] = "delegation_info"
	pParams["access_key"] = w.access_key
	pParams["id"] = id
	now := time.Now().Unix()
	pParams["created"] = strconv.FormatInt(now, 10)
	pParams["sign"] = w.createSign(pParams)

	return w.httpRequest(pParams)
}

func (w *HuobiTrade) Cancel_delegation(id string) bool {
	pParams := make(map[string]string)
	pParams["method"] = "cancel_delegation"
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
	//todo
	buyId := 0
	return (fmt.Sprintf("%d", buyId))
}

func (w *HuobiTrade) SellLTC(price, amount string) string {
	//todo
	sellId := 0
	return (fmt.Sprintf("%d", sellId))
}
