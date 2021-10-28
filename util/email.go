package util

import (
	"crypto/tls"
	"gopkg.in/gomail.v2"
)

func SendMail(to string, subject string, body string) error {
	cfg := LoadEmailCfg()

	d := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Passwd)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	m := gomail.NewMessage()
	m.SetHeader("From", cfg.Username)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	return d.DialAndSend(m)
}
