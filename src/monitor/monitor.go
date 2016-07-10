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

package monitor

import (
	"bittrex"
	"bitvc"
	. "common"
	. "config"
	"fmt"
	"huobi"
	"logger"
	"okcoin"
	"peatio"
	"simulate"
	"strategy"
	"strconv"
	"time"
)

func marketAPI() (marketAPI MarketAPI) {
	logger.Infof(Option["datacenter"])
	if Option["datacenter"] == "huobi" {
		marketAPI = huobi.NewHuobi()
	} else if Option["datacenter"] == "okcoin" {
		marketAPI = okcoin.NewOkcoin()
	} else if Option["datacenter"] == "peatio" {
		marketAPI = peatio.NewPeatio()
	} else if Option["datacenter"] == "bittrex" {
		marketAPI = Bittrex.Manager()
	} else {
		logger.Fatalln("Please config the market center...")
		panic(-1)
	}
	return
}

func tradeAPI() (tradeAPI TradeAPI) {
	if Option["tradecenter"] == "huobi" {
		tradeAPI = huobi.NewHuobi()
	} else if Option["tradecenter"] == "okcoin" {
		tradeAPI = okcoin.NewOkcoin()
	} else if Option["tradecenter"] == "bitvc" {
		tradeAPI = bitvc.NewBitvc()
	} else if Option["tradecenter"] == "peatio" {
		tradeAPI = peatio.NewPeatio()
	} else if Option["tradecenter"] == "bittrex" {
		tradeAPI = Bittrex.Manager()
	} else if Option["tradecenter"] == "simulate" {
		tradeAPI = simulate.NewSimulate()
	} else {
		logger.Fatalln("Please config the exchange center...")
		panic(0)
	}
	return
}

func RobotWorker() {
	ticker := time.NewTicker(1 * time.Second) // one second
	defer ticker.Stop()

	totalHour, _ := strconv.ParseInt(Option["totalHour"], 0, 64)
	if totalHour < 1 {
		totalHour = 1
	}

	fmt.Println("trade robot start working...")

	go func() {
		for _ = range ticker.C {
			peroid, _ := strconv.Atoi(Option["tick_interval"])
			strategyName := Option["strategy"]
			ret := true
			var records []Record
			if strategyName != "OPENORDER" {
				ret, records = marketAPI().GetKLine(peroid)
			}

			if ret != false {
				strategy.Tick(tradeAPI(), records)
			}
		}
	}()

	logger.Infof("程序将持续运行%d小时后停止", time.Duration(totalHour))

	time.Sleep(time.Duration(totalHour) * time.Hour)

	logger.Infof("程序到达设定时长%d小时，停止运行。", time.Duration(totalHour))
}

const worker_number = 1

type message struct {
	normal bool                   // true means exit normal, otherwise
	state  map[string]interface{} // goroutine state
}

func worker(mess chan message) {
	defer func() {
		exit_message := message{state: make(map[string]interface{})}
		i := recover()
		if i != nil {
			exit_message.normal = false
		} else {
			exit_message.normal = true
		}
		mess <- exit_message
	}()

	RobotWorker()
}

func supervisor(mess chan message) {
	for i := 0; i < worker_number; i++ {
		m := <-mess
		switch m.normal {
		case true:
			logger.Infoln("exit normal, nothing serious!")
		case false:
			logger.Infoln("exit abnormal, something went wrong")
		}
	}
}

func RunRobot() {
	mess := make(chan message, 10)
	for i := 0; i < worker_number; i++ {
		go worker(mess)
	}

	supervisor(mess)
}

func testHuobiAPI() {
	tradeAPI := huobi.NewHuobiTrade(SecretOption["huobi_access_key"], SecretOption["huobi_secret_key"])
	accout_info, _ := tradeAPI.GetAccount()
	fmt.Println(accout_info)

	fmt.Println(tradeAPI.GetAccount())
	if false {
		buyId := tradeAPI.BuyBTC("1000", "0.001")
		sellId := tradeAPI.SellBTC("10000", "0.001")

		if tradeAPI.Cancel_order(buyId) {
			fmt.Printf("cancel %s success \n", buyId)
		} else {
			fmt.Printf("cancel %s falied \n", buyId)
		}

		if tradeAPI.Cancel_order(sellId) {
			fmt.Printf("cancel %s success \n", sellId)
		} else {
			fmt.Printf("cancel %s falied \n", sellId)
		}
	}

	fmt.Println(tradeAPI.Get_orders())
}

func testBitVCAPI() {
	tradeAPI := bitvc.NewBitvc()
	accout_info, _ := tradeAPI.GetAccount()
	fmt.Println(accout_info)
	/*
		fmt.Println(tradeAPI.GetAccount())
		if false {
			buyId := tradeAPI.BuyBTC("1000", "0.001")
			sellId := tradeAPI.SellBTC("10000", "0.001")

			if tradeAPI.Cancel_order(buyId) {
				fmt.Printf("cancel %s success \n", buyId)
			} else {
				fmt.Printf("cancel %s falied \n", buyId)
			}

			if tradeAPI.Cancel_order(sellId) {
				fmt.Printf("cancel %s success \n", sellId)
			} else {
				fmt.Printf("cancel %s falied \n", sellId)
			}
		}

		fmt.Println(tradeAPI.Get_orders())
	*/
}

func testOkcoinBTCAPI() {
	tradeAPI := okcoin.NewOkcoinTrade(SecretOption["ok_partner"], SecretOption["ok_secret_key"])
	accout_info, _ := tradeAPI.GetAccount()
	fmt.Println(accout_info)

	buyret := tradeAPI.BuyBTC("1000", "0.01")
	fmt.Println(buyret)
	sellret := tradeAPI.SellBTC("10000", "0.01")
	fmt.Println(sellret)

	var orderTable okcoin.OKOrderTable
	ret, orderTable := tradeAPI.Get_BTCorder("-1")
	fmt.Println(ret, orderTable)

	time.Sleep(2000 * time.Millisecond)

	ret, orderTable = tradeAPI.Get_LTCorder("-1")
	fmt.Println(ret, orderTable)

	ret = tradeAPI.Cancel_BTCorder("-1")
	fmt.Println(ret)

	time.Sleep(2000 * time.Millisecond)

	ret = tradeAPI.Cancel_LTCorder("-1")
	fmt.Println(ret)
}

func testOkcoinLTCAPI() {
	tradeAPI := okcoin.NewOkcoinTrade(SecretOption["ok_partner"], SecretOption["ok_secret_key"])

	buyret := tradeAPI.BuyMarketLTC("100", "0.1")
	fmt.Println(buyret)

	time.Sleep(2000 * time.Millisecond)

	sellret := tradeAPI.SellMarketLTC("150", "0.1")
	fmt.Println(sellret)

	time.Sleep(2000 * time.Millisecond)

	buyret = tradeAPI.BuyLTC("100", "0.1")
	fmt.Println(buyret)

	time.Sleep(2000 * time.Millisecond)

	sellret = tradeAPI.SellLTC("150", "0.1")
	fmt.Println(sellret)

	ret, orderTable := tradeAPI.Get_LTCorder("-1")
	fmt.Println(ret, orderTable)

	ret = tradeAPI.Cancel_LTCorder("100253")
	fmt.Println(ret)
}
