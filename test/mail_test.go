package test

import (
	"cloud-disk/core/define"
	"crypto/tls"
	"github.com/jordan-wright/email"
	"net/smtp"
	"testing"
)

func TestSendMall(t *testing.T) {
	e := email.NewEmail()
	e.From = "Get <z15874610586@163.com>"
	e.To = []string{"3215881731@qq.com"}

	e.Subject = "验证码发送测试"
	//e.Text = []byte("Text Body is, of course, supported!")
	e.HTML = []byte("你的验证码：<h1>123</h1>")
	//e.Send("smtp.gmail.com:587", smtp.PlainAuth("", "test@gmail.com", "password123", "smtp.gmail.com"))
	err := e.SendWithTLS("smtp.163.com:465", smtp.PlainAuth("", "z15874610586@163.com", define.MailPassword, "smtp.163.com"),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.163.com"})
	if err != nil {
		t.Fatal(err)
	}
}
