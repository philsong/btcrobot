/*
  btcbot is a Bitcoin trading bot for HUOBI.com written
  in golang, it features multiple trading methods using
  technical analysis.

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
	"fmt"
)

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

func (w *Okcoin) AnalyzeKLinePeroid(peroid int) (ret bool) {
	fmt.Println("todo..")
	return false
}
