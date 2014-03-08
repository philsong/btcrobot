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

package service

import (
	"config"
	"email"
	"fmt"
	"logger"
	"strconv"
)

func TriggerTrender(alert string) error {

	if config.Remind["disable_email"] != "1" {
		if alert != "" {
			email.SendAlertEmail(config.Remind["to_email"], alert)
		}
	}

	return nil
}

func TriggerPrice(price float64) error {
	lowest_price, err := strconv.ParseFloat(config.Remind["lowest_price"], 64)
	if err != nil {
		logger.Debugln("config item lowest_price is not float")
		return err
	}
	highest_price, err := strconv.ParseFloat(config.Remind["highest_price"], 64)
	if err != nil {
		logger.Debugln("config item highest_price is not float")
		return err
	}

	var alert string
	if config.Remind["disable_email"] != "1" {
		if price < lowest_price {
			alert = fmt.Sprintf("价格 %f 低于设定的阀值 %f", price, config.Remind["lowest_price"])
		} else if price > highest_price {
			alert = fmt.Sprintf("价格 %f 超过设定的阀值 %f", price, config.Remind["highest_price"])
		}

		if alert != "" {
			email.SendAlertEmail(config.Remind["to_email"], alert)
		}
	}

	return nil
}
