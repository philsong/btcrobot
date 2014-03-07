package email

import (
	. "config"
	"encoding/base64"
	"fmt"
	"logger"
	"net/mail"
	"net/smtp"
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
	err := smtp.SendMail(SecretOption["smtp_addr"], auth, Remind["from_email"], tos, []byte(message))
	if err != nil {
		logger.Errorln("Send Mail to", strings.Join(tos, ","), "error:", err)
		return err
	}
	logger.Debugln("Send Mail to", strings.Join(tos, ","), "Successfully")
	return nil
}

func SendAlertEmail(receiver, alert string) error {
	// Set up authentication information.

	auth := smtp.PlainAuth("", SecretOption["smtp_username"], SecretOption["smtp_password"], SecretOption["smtp_host"])

	from := mail.Address{"BTCRobot监控中心", "78623269@qq.com"}
	to := mail.Address{"收件人", receiver}
	title := "BTC价格预警" + alert

	body := `
	<html>
	<body>
	<h3>
	"%s"
	</h3>
	捐助BTC，支持开发<span style="font-size: 80%"><a href="bitcoin:1NDnnWCUu926z4wxA3sNBGYWNQD3mKyes8">1NDnnWCUu926z4wxA3sNBGYWNQD3mKyes8</a></span>

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

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err := smtp.SendMail(
		SecretOption["smtp_addr"],
		auth,
		from.Address,
		[]string{to.Address},
		[]byte(message),
		//[]byte("This is the email body."),
	)
	if err != nil {
		logger.Errorln("Send Mail to", to, "error:", err)
		return err
	}
	logger.Debugln("Send Mail to", to, "Successfully")
	return nil
}

func NoticeEmail2() {
	tos := []string{"78623269@qq.com", "philsong@techtrex.com"}

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
