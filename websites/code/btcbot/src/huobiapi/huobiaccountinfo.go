// Copyright 2014 The btcbot Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://github.com/philsong/
// Author：Phil	78623269@qq.com

package huobiapi

import (
	"encoding/json"
	"logger"
)

/*
import (
	"io"
	"logger"
	"strings"
)
*/

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

func (w *Huobi) TradeMyTradeInfoAnalyze(body string) bool {

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
