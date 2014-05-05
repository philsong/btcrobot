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
	. "common"
	. "config"
	"fmt"
	"huobi"
	"logger"
	"okcoin"
	"strategy"
	"strconv"
	"time"
)

/*
func backtesting() {
	fmt.Println("back testing begin...")
	huobi := huobi.NewHuobi()
	huobi.Disable_trading = 1

	peroids := []int{1, 5, 15, 30, 60, 100}
	for _, peroid := range peroids {
		if huobi.AnalyzeKLine(peroid) == true {
		} else {
			logger.Errorln("TradeKLine failed.")
		}
	}

	fmt.Println("生成 1/5/15/30/60分钟及1天 周期的后向测试报告于log/reportxxx.log文件中,请查看")

	fmt.Println("back testing end ...")
}
*/

func marketAPI() (marketAPI MarketAPI) {
	if Option["datacenter"] == "huobi" {
		marketAPI = huobi.NewHuobi()
	} else if Option["datacenter"] == "okcoin" {
		marketAPI = okcoin.NewOkcoin()
	} else {
		logger.Fatalln("Please config the datacenter firstly...")

	}
	return
}

func tradeAPI() (tradeAPI TradeAPI) {
	if Option["tradecenter"] == "huobi" {
		tradeAPI = huobi.NewHuobi()
	} else if Option["tradecenter"] == "okcoin" {
		tradeAPI = okcoin.NewOkcoin()
	} else {
		logger.Fatalln("Please config the tradecenter firstly...")

	}
	return
}

func RobotWorker() {
	fmt.Println("env", Config["env"])
	if DebugEnv || Config["env"] == "dev" {
		fmt.Println("test working...")

		var tradeAPI TradeAPI
		tradeAPI = okcoin.NewOkcoin()
		tradeAPI.GetAccount()
		tradeAPI.GetOrderBook()

		tradeAPI = huobi.NewHuobi()
		tradeAPI.GetAccount()
		ret, orderbook := tradeAPI.GetOrderBook()
		fmt.Println(ret, orderbook)

		//testHuobiAPI()
		//testOkcoinLTCAPI()
		return
	}

	ticker := time.NewTicker(1 * time.Second) //2s
	defer ticker.Stop()

	totalHour, _ := strconv.ParseInt(Option["totalHour"], 0, 64)
	if totalHour < 1 {
		totalHour = 1
	}

	fmt.Println("robot working...")

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
	normal bool                   //true means exit normal, otherwise
	state  map[string]interface{} //goroutine state
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

	/*
		now := time.Now()
		seed := now.UnixNano()
		rand.Seed(seed)
		num := rand.Int63()
		fmt.Println(num)
		if num%2 != 0 {
			fmt.Println("1")
			panic("not evening")
		} else {
			fmt.Println("0")
			runtime.Goexit()
		}
	*/
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

	//	fmt.Println(tradeAPI.GetAccount())
	if false {
		buyId := tradeAPI.BuyBTC("1000", "0.001")
		sellId := tradeAPI.SellBTC("10000", "0.001")

		//fmt.Println(tradeAPI.Get_delegations())
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

	//orderTable, ret := tradeAPI.Get_LTCorder("-1")
	//fmt.Println(ret, orderTable)

	//ret = tradeAPI.Cancel_LTCorder("100253")
	//fmt.Println(ret)
}
