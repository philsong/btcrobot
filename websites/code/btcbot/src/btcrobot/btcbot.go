// Copyright 2014 The btcbot Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://github.com/philsong/
// Author：Phil	78623269@qq.com

package main

import (
	. "config"
	"fmt"
	"huobiapi"
	"logger"
	"strconv"
	"time"
)

func doTradeDelegation() {
	huobi := huobiapi.NewHuobi()
	logger.Infoln("doTradeDelegation start....")
	if huobi.Login() == true {
		logger.Debugln("Login successfully.")

		huobi.TradeDelegation()
	} else {
		logger.Debugln("Login failed.")
	}

	logger.Infoln("doTradeDelegation end-----")
}

func backtesting(done chan bool) {
	fmt.Println("back testing begin...")
	huobi := huobiapi.NewHuobi()

	peroids := []int{1, 5, 15, 30, 60, 100}
	for _, v := range peroids {
		huobi.Peroid = v
		if huobi.TradeKLinePeroid(huobi.Peroid) == true {

		} else {
			logger.Errorln("TradeKLine failed.")
		}
	}
	fmt.Println("back testing end ...")
	done <- true
}

func testKLineAPI(done chan bool) {
	ticker := time.NewTicker(time.Millisecond * 2000)

	huobi := huobiapi.NewHuobi()
	huobi.Peroid, _ = strconv.Atoi(Option["tick_interval"])

	slippage, err := strconv.ParseFloat(Config["slippage"], 64)
	if err != nil {
		logger.Debugln("config item slippage is not float")
		slippage = 0
	}
	huobi.Slippage = slippage

	go func() {
		for _ = range ticker.C {
			if huobi.Peroid == 1 {
				huobi.TradeKLineMinute()
			} else {
				huobi.TradeKLinePeroid(huobi.Peroid)
			}
		}
	}()
	time.Sleep(time.Millisecond * 24 * 60 * 60 * 1000)

	ticker.Stop()
	fmt.Println("Ticker stopped")
	done <- true
}

func TestTradeAPI() {
	tradeAPI := huobiapi.NewHuobiTrade(SecretOption["access_key"], SecretOption["secret_key"])
	//	fmt.Println(tradeAPI.Get_account_info())
	if false {
		buyId := tradeAPI.Buy("1000", "0.001")
		sellId := tradeAPI.Sell("10000", "0.001")

		//fmt.Println(tradeAPI.Get_delegations())
		if tradeAPI.Cancel_delegation(buyId) {
			fmt.Printf("cancel %s success \n", buyId)
		} else {
			fmt.Printf("cancel %s falied \n", buyId)
		}

		if tradeAPI.Cancel_delegation(sellId) {
			fmt.Printf("cancel %s success \n", sellId)
		} else {
			fmt.Printf("cancel %s falied \n", sellId)
		}
	}

	fmt.Println(tradeAPI.Get_delegations())
}

func tradeService() {

	done := make(chan bool, 1)
	fmt.Println("working...")

	if Config["backtesting"] == "true" {
		go backtesting(done)
		<-done
	} else {
		fmt.Println("trade monitor...")
		go testKLineAPI(done)
		<-done
	}

	fmt.Println("done")

	return
	//doTradeDelegation()
}
