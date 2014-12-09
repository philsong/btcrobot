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

package bitvc

import (
	"compress/gzip"
	. "config"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"logger"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"util"
)

/*
	http://www.huobi.com/help/index.php?a=api_help
*/
type BitvcTrade struct {
	isLogin     bool
	client      *http.Client
	MyTradeInfo TradeInfo
}

func NewBitvcTrade(access_key, secret_key string) *BitvcTrade {
	w := new(BitvcTrade)
	return w
}

func (w *BitvcTrade) httpRequest(reqURL string, pParams map[string]string) (string, error) {
	v := url.Values{}
	for key, val := range pParams {
		v.Add(key, val)
	}

	req, err := http.NewRequest("POST", reqURL, strings.NewReader(v.Encode()))
	if err != nil {
		logger.Fatal(err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", Config["bitvc_base_url"])
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36")
	logger.Traceln(req)

	c := util.NewTimeoutClient()

	logger.Tracef("HTTP req begin BitvcTrade")
	resp, err := c.Do(req)
	logger.Tracef("HTTP req end BitvcTrade")
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

func (w *BitvcTrade) check_json_result(body string) (errorMsg ErrorMsg, ret bool) {
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

/*
{"code":0,"msg":"","ext":{"cny_balance":{"CNY":{"total":340062540230000,"available":340062540230000,"frozen":0},"BTC":{"total":177247000000,"available":177247000000,"frozen":0},"LTC":{"total":0,"available":0,"frozen":0},"now_btc_price":"38380100000000","now_ltc_price":"471500000000","LOAN_CNY":0,"LOAN_BTC":177247000000,"LOAN_LTC":0,"net_asset":340062540230000,"total":1020338298700000},"coin_saving_balance":{"btc_total":"0","ltc_total":"0"}}}
*/
func (w *BitvcTrade) GetAccount() (account_info Account_info, ret bool) {

	w.Login()

	pParams := make(map[string]string)

	ret = true

	body, err := w.httpRequest("https://www.bitvc.com/ajax/user_balance", pParams)
	if err != nil {
		ret = false
		return
	}

	fmt.Println(body)
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

func (w *BitvcTrade) Get_orders() (ret bool, m []HBOrderItem) {
	pParams := make(map[string]string)

	ret = true

	body, err := w.httpRequest(Config["bitvc_login_url"], pParams)
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

func (w *BitvcTrade) Get_order(id string) (ret bool, m HBOrder) {
	pParams := make(map[string]string)

	ret = true

	body, err := w.httpRequest(Config["bitvc_login_url"], pParams)
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

func (w *BitvcTrade) Cancel_order(id string) bool {
	pParams := make(map[string]string)

	body, err := w.httpRequest(Config["bitvc_login_url"], pParams)
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

func (w *BitvcTrade) doTrade(method, price, amount string) int {
	pParams := make(map[string]string)

	body, err := w.httpRequest(Config["bitvc_login_url"], pParams)
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

func (w *BitvcTrade) BuyBTC(price, amount string) string {
	buyId := w.doTrade("buy", price, amount)
	return (fmt.Sprintf("%d", buyId))
}

func (w *BitvcTrade) SellBTC(price, amount string) string {
	sellId := w.doTrade("sell", price, amount)
	return (fmt.Sprintf("%d", sellId))
}

func (w *BitvcTrade) BuyLTC(price, amount string) string {
	// todo
	buyId := 0
	return (fmt.Sprintf("%d", buyId))
}

func (w *BitvcTrade) SellLTC(price, amount string) string {
	// todo
	sellId := 0
	return (fmt.Sprintf("%d", sellId))
}

/*
		function calc_password_security_score(t) {
	    var e = 0;
	    return t.length < 4 ? e :
	     (t.length >= 8 && e++, t.length >= 10 && e++, /[a-z]/.test(t) && /[A-Z]/.test(t) && e++, /[0-9]/.test(t) && e++, /.[!,@,#,$,%,^,&,*,?,_,~, -,£,(,)]/.test(t) && e++, e)
		}
*/
func getPSS(clear_password string) int {
	pwd_security_score := 0
	fmt.Println(clear_password)
	if len(clear_password) < 4 {
		pwd_security_score = 0
	} else {
		if len(clear_password) >= 8 {
			pwd_security_score++
		}

		if len(clear_password) >= 10 {
			pwd_security_score++
		}

		regLower := regexp.MustCompile(`[[:lower:]]`)
		fmt.Printf("%q\n", regLower.FindAllString(clear_password, -1))
		regUpper := regexp.MustCompile(`[[:upper:]]`)
		fmt.Printf("%q\n", regUpper.FindAllString(clear_password, -1))

		matchedLower, err := regexp.MatchString(`[[:lower:]]`, clear_password)
		fmt.Println(matchedLower, err)
		matchedUpper, err := regexp.MatchString(`[[:upper:]]`, clear_password)
		fmt.Println(matchedUpper, err)

		if matchedLower && matchedUpper {
			pwd_security_score++
		}

		regDigit := regexp.MustCompile(`[[:digit:]]`)
		fmt.Printf("%q\n", regDigit.FindAllString(clear_password, -1))

		matchedDigit, err := regexp.MatchString(`[[:digit:]]`, clear_password)
		fmt.Println(matchedDigit, err)

		if matchedDigit {
			pwd_security_score++
		}

		regSpecial := regexp.MustCompile(`[!|@|#|$|%|\\^|&|\\*|\\?|_|~|-|£|\\(|\\)]`)
		fmt.Printf("%q\n", regSpecial.FindAllString(clear_password, -1))

		matchedSpecial, err := regexp.MatchString(`[!|@|#|$|%|\\^|&|\\*|\\?|_|~|-|£|\\(|\\)]`, clear_password)
		fmt.Println(matchedSpecial, err)

		if matchedSpecial {
			pwd_security_score++
		}
	}

	fmt.Println(pwd_security_score)
	return pwd_security_score
}

/*
	txcode := map[int]string{
		0:  `买单已委托，<a href="/trade/index.php?a=delegation">查看结果</a>`,
		2:  `没有足够的人民币`,
		10:	`没有足够的比特币`,
		16: `您需要登录才能继续`,
		17: `没有权限`,
		42:	`该委托已经取消，不能修改`,
		44:	`交易价钱太低`,
		56:`卖出价格不能低于限价的95%`}

	logger.Traceln(txcode[m.Code])
*/
func (w *BitvcTrade) Login() bool {
	fmt.Println("login....")
	login_url := Config["bitvc_login_url"]
	email := SecretOption["bitvc_email"]
	clear_password := SecretOption["bitvc_password"]
	password := util.Md5(clear_password + "hi,pwd")

	pwd_security_score := getPSS(clear_password)
	str_pwd_security_score := fmt.Sprintf("%d", pwd_security_score)
	post_arg := url.Values{"email": {email}, "password": {password}, "backurl": {"/index/index"}, "pwd_security_score": {str_pwd_security_score}}

	//logger.Traceln(strings.NewReader(post_arg.Encode()))
	req, err := http.NewRequest("POST", login_url, strings.NewReader(post_arg.Encode()))
	if err != nil {
		logger.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", Config["bitvc_base_url"])
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36")
	logger.Infoln(req)

	//jar := NewJar()
	/*
		// how to do compatible like c define?
		jar, _ := cookiejar.New(nil)
		fmt.Println("version:", runtime.Version())
		if runtime.Version() != "go1.3" {
			w.client = &http.Client{nil, nil, jar}
		} else {
			w.client = &http.Client{nil, nil, jar, 10 * time.Second}
		}
	*/
	//w.client = new(http.Client)

	resp, err := w.client.Do(req)
	if err != nil {
		logger.Fatal(err)
	}
	defer resp.Body.Close()
	logger.Infof("Login resp StatusCode=%v", resp.StatusCode)
	logger.Infof("Login resp=%v", resp)
	if resp.StatusCode == 200 {
		var body string

		contentEncoding := resp.Header.Get("Content-Encoding")
		logger.Infof("HTTP returned Content-Encoding %s", contentEncoding)
		switch contentEncoding {
		case "gzip":
			body = DumpGZIP(resp.Body)

		default:
			bodyByte, _ := ioutil.ReadAll(resp.Body)
			body = string(bodyByte)
			ioutil.WriteFile("login.html", bodyByte, os.ModeAppend)
		}

		logger.Traceln(resp.Header.Get("Content-Type"))
		ret := strings.Contains(body, "用户名或者密码错误")
		if ret {
			logger.Traceln("用户名或者密码错误")
			return false
		}

		w.isLogin = true
		return true
	} else if resp.StatusCode == 500 {
		w.isLogin = true
		return true
	} else {
		logger.Infof("resp %v", resp)
	}

	return false
}

func (w *BitvcTrade) TradeAdd(a, price, amount string) bool {

	if w.isLogin == false {
		if w.Login() == false {
			return false
		}
	}
	/*
		if w.checkAccount(a, price, amount) == false {
			return false
		}
	*/
	post_arg := url.Values{
		"a":      {a},
		"price":  {price},
		"amount": {amount},
	}

	req, err := http.NewRequest("POST", Config["trade_add_url"], strings.NewReader(post_arg.Encode()))
	if err != nil {
		logger.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "https://www.huobi.com/")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	logger.Traceln(req)

	resp, err := w.client.Do(req)
	if err != nil {
		logger.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		var body string

		contentEncoding := resp.Header.Get("Content-Encoding")
		logger.Tracef("HTTP returned Content-Encoding %s", contentEncoding)
		switch contentEncoding {
		case "gzip":
			body = DumpGZIP(resp.Body)

		default:
			bodyByte, _ := ioutil.ReadAll(resp.Body)
			body = string(bodyByte)
			ioutil.WriteFile("TradeAdd.json", bodyByte, os.ModeAppend)
		}

		logger.Traceln(resp.Header.Get("Content-Type"))

		if resp.Header.Get("Content-Type") == "application/json" {
			doc := json.NewDecoder(strings.NewReader(body))

			type Msg struct {
				Code int
				Msg  string
			}

			var m Msg
			if err := doc.Decode(&m); err == io.EOF {
				logger.Traceln(err)
			} else if err != nil {
				logger.Fatal(err)
			}

			logger.Infoln(m)

			if m.Code == 0 {
				return true
			} else {
				return false
			}
		} else {
			ret := strings.Contains(body, "您需要登录才能继续")
			if ret {
				logger.Errorln("您需要登录才能继续")
				w.isLogin = false
				return false
			} else {
				return true
			}
		}

	} else {
		logger.Tracef("HTTP returned status %v", resp)
	}

	return false

}

func (w *BitvcTrade) checkAccount(a, price, amount string) bool {
	btc, cny := w.get_account_info()

	FPrice, err := strconv.ParseFloat(price, 64)
	if err != nil {
		logger.Debugln("price is not float")
		return false
	}
	FAmount, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		logger.Debugln("amount is not float")
		return false
	}
	if a == "do_buy" {
		if float64(cny) < FPrice*FAmount {
			return false
		}
	} else {
		if float64(btc) < FAmount {
			return false
		}
	}

	return true
}

func (w *BitvcTrade) Do_buy(tradePrice, tradeAmount string) bool {
	if w.Login() == true {
		if w.TradeAdd("do_buy", tradePrice, tradeAmount) == true {
			logger.Infoln("执行买入委托成功", tradePrice, tradeAmount)
			return true
		} else {
			logger.Infoln("执行买入失败，原因：买入操作被服务器拒绝", tradePrice, tradeAmount)
		}
	} else {
		logger.Infoln("执行买入失败，原因：登陆失败", tradePrice, tradeAmount)

	}

	return false
}

func (w *BitvcTrade) Do_sell(tradePrice, tradeAmount string) bool {
	if w.Login() == true {
		if w.TradeAdd("do_sell", tradePrice, tradeAmount) == true {
			logger.Infoln("执行卖出委托成功", tradePrice, tradeAmount)
			return true
		} else {
			logger.Infoln("执行卖出失败，原因：卖出操作被服务器拒绝", tradePrice, tradeAmount)
		}
	} else {
		logger.Infoln("执行卖出失败，原因：登陆失败", tradePrice, tradeAmount)
	}

	return false
}

func (w *BitvcTrade) TradeUpdate(id, price, amount string) bool {
	post_arg := url.Values{
		"id":     {id},
		"price":  {price},
		"amount": {amount},
	}

	req, err := http.NewRequest("POST", Config["trade_update_url"], strings.NewReader(post_arg.Encode()))
	if err != nil {
		logger.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "https://www.huobi.com/")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	logger.Traceln(req)

	resp, err := w.client.Do(req)
	if err != nil {
		logger.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		var body string

		contentEncoding := resp.Header.Get("Content-Encoding")
		logger.Tracef("HTTP returned Content-Encoding %s", contentEncoding)
		switch contentEncoding {
		case "gzip":
			body = DumpGZIP(resp.Body)

		default:
			bodyByte, _ := ioutil.ReadAll(resp.Body)
			body = string(bodyByte)
			ioutil.WriteFile("TradeUpdate.json", bodyByte, os.ModeAppend)
		}

		logger.Traceln(resp.Header.Get("Content-Type"))

		if resp.Header.Get("Content-Type") == "application/json" {
			doc := json.NewDecoder(strings.NewReader(body))

			type Msg struct {
				Code int
				Msg  string
			}

			var m Msg
			if err := doc.Decode(&m); err == io.EOF {
				logger.Traceln(err)
			} else if err != nil {
				logger.Fatal(err)
			}

			logger.Traceln(m)

			if m.Code == 0 {
				return true
			} else {
				return false
			}
		} else {
			ret := strings.Contains(body, "您需要登录才能继续")
			if ret {
				logger.Traceln("您需要登录才能继续")
				w.isLogin = false
				return false
			} else {
				return true
			}
		}

	} else {
		logger.Tracef("HTTP returned status %v", resp)
	}

	return false
}

func (w *BitvcTrade) QueryMyTradeInfo() bool {
	fmt.Println(w.isLogin)
	if w.isLogin == false {
		if w.Login() == false {
			return false
		}
	}

	req, err := http.NewRequest("GET", fmt.Sprintf(Config["my_trade_info"], rand.Float64()), nil)
	if err != nil {
		logger.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", Config["trade_query_url"])
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	logger.Traceln(req)

	resp, err := w.client.Do(req)
	if err != nil {
		logger.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		var body string

		contentEncoding := resp.Header.Get("Content-Encoding")
		logger.Tracef("HTTP returned Content-Encoding %s", contentEncoding)
		switch contentEncoding {
		case "gzip":
			body = DumpGZIP(resp.Body)

		default:
			bodyByte, _ := ioutil.ReadAll(resp.Body)
			body = string(bodyByte)
			ioutil.WriteFile("MyTradeInfo.json", bodyByte, os.ModeAppend)
		}

		logger.Traceln(resp.Header.Get("Content-Type"))

		if resp.Header.Get("Content-Type") == "application/json" {
			w.TradeMyTradeInfoAnalyze(body)
			if w.MyTradeInfo.Code == 0 {
				return true
			} else {
				return false
			}
		} else {
			ret := strings.Contains(body, "您需要登录才能继续")
			if ret {
				logger.Traceln("您需要登录才能继续")
				w.isLogin = false
				return false
			} else {
				return true
			}
		}

	} else {
		logger.Tracef("HTTP returned status %v", resp)
	}

	return false
}

func (w *BitvcTrade) TradeCancel(id string) bool {

	req, err := http.NewRequest("GET", fmt.Sprintf(Config["trade_cancel_url"], id), nil)
	if err != nil {
		logger.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", Config["trade_delegation"])
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	logger.Traceln(req)

	resp, err := w.client.Do(req)
	if err != nil {
		logger.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		var body string

		contentEncoding := resp.Header.Get("Content-Encoding")
		logger.Tracef("HTTP returned Content-Encoding %s", contentEncoding)
		switch contentEncoding {
		case "gzip":
			body = DumpGZIP(resp.Body)

		default:
			bodyByte, _ := ioutil.ReadAll(resp.Body)
			body = string(bodyByte)
			ioutil.WriteFile("TradeCancel.json", bodyByte, os.ModeAppend)
		}

		logger.Traceln(resp.Header.Get("Content-Type"))

		if resp.Header.Get("Content-Type") == "application/json" {
			doc := json.NewDecoder(strings.NewReader(body))

			type Msg struct {
				Code int
				Msg  string
			}

			var m Msg
			if err := doc.Decode(&m); err == io.EOF {
				logger.Traceln(err)
			} else if err != nil {
				logger.Fatal(err)
			}
			logger.Traceln(m)

			if m.Code == 0 {
				return true
			} else {
				return false
			}
		} else {
			ret := strings.Contains(body, "您需要登录才能继续")
			if ret {
				logger.Traceln("您需要登录才能继续")
				w.isLogin = false
				return false
			} else {
				return true
			}
		}

	} else {
		logger.Tracef("HTTP returned status %v", resp)
	}

	return false
}

func (w *BitvcTrade) TradeDelegation() bool {
	fmt.Println(w.isLogin)
	if w.isLogin == false {
		if w.Login() == false {
			return false
		}
	}

	req, err := http.NewRequest("GET", Config["trade_delegation"], nil)
	if err != nil {
		logger.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", Config["trade_add_url"])
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	logger.Traceln(req)

	resp, err := w.client.Do(req)
	if err != nil {
		logger.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		var body string

		contentEncoding := resp.Header.Get("Content-Encoding")
		logger.Tracef("HTTP returned Content-Encoding %s", contentEncoding)
		switch contentEncoding {
		case "gzip":
			body = DumpGZIP(resp.Body)

		default:
			bodyByte, _ := ioutil.ReadAll(resp.Body)
			body = string(bodyByte)
			ioutil.WriteFile("TradeDelegation.html", bodyByte, os.ModeAppend)
		}

		logger.Traceln(resp.Header.Get("Content-Type"))

		if resp.Header.Get("Content-Type") == "application/json" {
			doc := json.NewDecoder(strings.NewReader(body))

			type Msg struct {
				Code int
				Msg  string
			}

			for {
				var m Msg
				if err := doc.Decode(&m); err == io.EOF {
					logger.Traceln(err)
					break
				} else if err != nil {
					logger.Fatal(err)
				}
				logger.Traceln(m)

				if m.Code == 0 {
					return true
				} else {
					return false
				}
			}

			return false
		} else {
			ret := strings.Contains(body, "您需要登录才能继续")
			if ret {
				logger.Traceln("您需要登录才能继续")
				w.isLogin = false
				return false
			} else {
				return w.TradeDelegationAnalyze(body)
			}
		}

	} else {
		logger.Tracef("HTTP returned status %v", resp)
	}

	return false
}

func (w *BitvcTrade) get_account_info() (btc, cny int) {
	if w.QueryMyTradeInfo() == true {
		logger.Infoln("BTC：", w.MyTradeInfo.Extra.Balance.Available_btc_display)
		logger.Infoln("CNY：", w.MyTradeInfo.Extra.Balance.Available_cny_display)
		logger.Infoln("Frozen BTC：", w.MyTradeInfo.Extra.Balance.Frozen_btc_display)
		logger.Infoln("Frozen CNY: ", w.MyTradeInfo.Extra.Balance.Frozen_cny_display)
		logger.Infoln("CNY：", w.MyTradeInfo.Extra.Balance.Available_cny_display)
		logger.Infoln("total", w.MyTradeInfo.Extra.Balance.Total)
		logger.Infoln("(BTC,CNY)", w.MyTradeInfo.Extra.Balance.Available_btc, w.MyTradeInfo.Extra.Balance.Available_cny)
		return w.MyTradeInfo.Extra.Balance.Available_btc, w.MyTradeInfo.Extra.Balance.Available_cny
	} else {
		logger.Errorln("get_account_info failed.")
		return 0, 0
	}
}

func DumpGZIP(r io.Reader) string {
	var body string
	reader, _ := gzip.NewReader(r)
	for {
		buf := make([]byte, 1024)
		n, err := reader.Read(buf)

		if err != nil && err != io.EOF {
			panic(err)
		}

		if n == 0 {
			break
		}
		body += string(buf)
	}
	return body
}

type Sell struct {
	Price         string `json:"price"`
	Amount        string `json:"amount"`
	Available_btc string `json:"available_btc"`
}

type Buy struct {
	Price         string `json:"price"`
	Amount        string `json:"amount"`
	Available_cny string `json:"available_cny"`
}

type Balance struct {
	Id                    int    `json:"id"`
	Initialized_cny       int    `json:"initialized_cny"`
	Initialized_btc       int    `json:"initialized_btc"`
	User_id               int    `json:"user_id"`
	Available_cny         int    `json:"available_cny"`
	Available_btc         int    `json:"available_btc"`
	Available_usd         int    `json:"available_usd"`
	Frozen_cny            int    `json:"frozen_cny"`
	Frozen_btc            int    `json:"frozen_btc"`
	Frozen_usd            int    `json:"frozen_usd"`
	Debt_bitcoin          int    `json:"debt_bitcoin"`
	Debt_rmb              int    `json:"debt_rmb"`
	Total                 string `json:"total"`
	Loan_total            string `json:"loan_total"`
	Net_asset             string `json:"net_asset"`
	Loan_cny_display      string `json:"loan_cny_display"`
	Loan_btc_display      string `json:"loan_btc_display"`
	Available_btc_display string `json:"available_btc_display"`
	Available_cny_display string `json:"available_cny_display"`
	Frozen_btc_display    string `json:"frozen_btc_display"`
	Frozen_cny_display    string `json:"frozen_cny_display"`
}
type Extra struct {
	Sell    Sell    `json:"sell"`
	Buy     Buy     `json:"buy"`
	Balance Balance `json:"balance"`
}

type TradeInfo struct {
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
	Extra Extra  `json:"extra"`
}

func (w *BitvcTrade) TradeMyTradeInfoAnalyze(body string) bool {

	if err := json.Unmarshal([]byte(body), &w.MyTradeInfo); err != nil {
		logger.Debugln("error:", err)
		logger.Debugln("MyTradeInfo json....panic!!!")
		logger.Debugln(body)
		logger.Debugln("MyTradeInfo---------------------------panic!!!")
		return false
	}
	/*
		doc := json.NewDecoder(strings.NewReader(body))
		if err := doc.Decode(&w.MyTradeInfo); err == io.EOF {
			logger.Debugln(err)
		} else if err != nil {
			logger.Fatal(err)
		}
	*/
	logger.Debugln(w.MyTradeInfo)
	return true
}

func (w *BitvcTrade) TradeCancel1stPage(ids []string) bool {
	logger.Debugln("TradeCancel1stPage start....")

	for _, id := range ids {
		w.TradeCancel(id)
	}

	logger.Debugln("TradeCancel1stPage end-----")
	return true
}

func (w *BitvcTrade) TradeDelegationAnalyze(body string) bool {
	logger.Debugln("TradeDelegationAnalyze start....")

	var ids []string
	for _, pat := range strings.Split(body, `">撤单`) {
		if len(pat) == 0 {
			// Empty strings such as from a trailing comma can be ignored.
			continue
		}
		patLev := strings.Split(pat, "a=cancel&id=")
		if len(patLev) != 2 || len(patLev[0]) == 0 || len(patLev[1]) == 0 {
			logger.Debugln("parse end")
			break
		}

		logger.Debugln(patLev[1])

		ids = append(ids, patLev[1])
	}

	logger.Debugln(ids)

	w.TradeCancel1stPage(ids)

	logger.Debugln("TradeDelegationAnalyze end-----")
	return true
}
