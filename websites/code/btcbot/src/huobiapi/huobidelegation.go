// Copyright 2014 The btcbot Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// https://github.com/philsong/
// Author：Phil	78623269@qq.com

package huobiapi

import (
	"logger"
	"strings"
)

func (w *Huobi) TradeCancel1stPage(ids []string) bool {
	logger.Debugln("TradeCancel1stPage start....")

	for _, id := range ids {
		w.TradeCancel(id)
	}

	logger.Debugln("TradeCancel1stPage end-----")
	return true
}

func (w *Huobi) TradeDelegationAnalyze(body string) bool {
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
