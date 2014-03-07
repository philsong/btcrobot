// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

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

func SendWelcomeMail(emailaddr []string) {
	content := `Welcome to BTC预警网.<br><br>
				欢迎您，成功注册成为 BTC预警网 会员<br><br>
				前往 <a href="http://pingliwang.com:8080">BTC预警网</a><br>
				<div style="text-align:right;">&copy;2013 BTC预警网</div>`
	email.SendMail("BTC预警网 注册成功通知", content, emailaddr)
}

// 发重置密码邮件
func SendResetpwdMail(emailaddr, uuid string) {
	content := `您好，` + emailaddr + `,<br/><br/>
&nbsp;&nbsp;&nbsp;&nbsp;我们的系统收到一个请求，说您希望通过电子邮件重新设置您在 <a href="http://` + config.Config["domain"] + `">微言微语</a> 的密码。您可以点击下面的链接重设密码：<br/><br/>

&nbsp;&nbsp;&nbsp;&nbsp;http://` + config.Config["domain"] + `/account/resetpwd?code=` + uuid + ` <br/><br/>

如果这个请求不是由您发起的，那没问题，您不用担心，您可以安全地忽略这封邮件。<br/><br/>

如果您有任何疑问，可以回复这封邮件向我们提问。谢谢！<br/><br/>`

	email.SendMail("重设密码 ", content, []string{emailaddr})
}
