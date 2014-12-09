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

package email

import (
	. "config"
	"encoding/base64"
	"fmt"
	"logger"
	"net/mail"
	"net/smtp"
	"strconv"
	"strings"
)

func encodeRFC2047(String string) string {
	// use mail's rfc2047 to encode any string
	addr := mail.Address{String, ""}
	return strings.Trim(addr.String(), " <>")
}

/*
 *	user : example@example.com login smtp server user
 *	password: xxxxx login smtp server password
 *	host: smtp.example.com:port   smtp.163.com:25
 *	to: example@example.com;example1@163.com;example2@sina.com.cn;...
 *  subject:The subject of mail
 *  body: The content of mail
 *  mailtyoe: mail type html or text
 */

// 发送电子邮件功能
func SendMail(subject, content string, tos []string) error {
	message := `From: btc-robot
				To: ` + strings.Join(tos, ",") + `
				Subject: ` + subject + `
				Content-Type: text/html;charset=UTF-8

				` + content

	auth := smtp.PlainAuth("", SecretOption["smtp_username"], SecretOption["smtp_password"], SecretOption["smtp_host"])
	err := smtp.SendMail(SecretOption["smtp_addr"], auth, Option["from_email"], tos, []byte(message))
	if err != nil {
		logger.Infoln("Send Mail to", strings.Join(tos, ","), "error:", err)
		return err
	}
	logger.Infoln("Send Mail to", strings.Join(tos, ","), "Successfully")
	return nil
}

func SendAlertEmail(receiver, alert string) error {
	// Set up authentication information.
	auth := smtp.PlainAuth("", SecretOption["smtp_username"], SecretOption["smtp_password"], SecretOption["smtp_host"])
	from := mail.Address{"BTCRobot监控中心", SecretOption["smtp_username"]}
	to := mail.Address{"收件人", receiver}
	title := "BTCRobot来电--->" + alert

	body := `
	<html>
	<body>
	<h3>
	%s
	</h3>
	<p>
	捐助BTC，支持开发<span style="font-size: 80%"><a href="bitcoin:1NDnnWCUu926z4wxA3sNBGYWNQD3mKyes8">1NDnnWCUu926z4wxA3sNBGYWNQD3mKyes8</a></span>
	</p>
	</body>
	</html>
	`
	body = fmt.Sprintf(body, alert)

	header := make(map[string]string)
	header["From"] = from.String()
	header["To"] = to.String()
	header["Subject"] = encodeRFC2047(title)
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	logger.Debugln("Try sending Mail to", to)
	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err := smtp.SendMail(
		SecretOption["smtp_addr"],
		auth,
		from.Address,
		[]string{to.Address},
		[]byte(message))
	if err != nil {
		logger.Infoln("Send Mail to", to, "error:", err)
		return err
	}
	logger.Debugln("Send Mail to", to, "Successfully")
	return nil
}

func NoticeEmailV2() {
	tos := []string{"78623269@qq.com", "songbohr@gmail.com"}

	subject := "Test send email"

	body := `
	<html>
	<body>
	<h3>
	"Test send email by btcrobot"
	</h3>
	</body>
	</html>
	`
	fmt.Println("send email")
	err := SendMail(subject, body, tos)
	if err != nil {
		fmt.Println("send mail error!")
		fmt.Println(err)
	} else {
		fmt.Println("send mail success!")
	}
}

func TriggerTrender(alert string) error {

	if Option["enable_email"] == "1" {
		if alert != "" {
			SendAlertEmail(Option["to_email"], alert)
		}
	}

	return nil
}

func TriggerPrice(price float64) error {
	lowest_price, err := strconv.ParseFloat(Option["lowest_price"], 64)
	if err != nil {
		logger.Debugln("config item lowest_price is not float")
		return err
	}
	highest_price, err := strconv.ParseFloat(Option["highest_price"], 64)
	if err != nil {
		logger.Debugln("config item highest_price is not float")
		return err
	}

	var alert string
	if Option["enable_email"] == "1" {
		if price < lowest_price {
			alert = fmt.Sprintf("价格 %f 低于设定的阀值 %f", price, Option["lowest_price"])
		} else if price > highest_price {
			alert = fmt.Sprintf("价格 %f 超过设定的阀值 %f", price, Option["highest_price"])
		}

		if alert != "" {
			SendAlertEmail(Option["to_email"], alert)
		}
	}

	return nil
}
