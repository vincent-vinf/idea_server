package utils

import (
	"crypto/tls"
	"gopkg.in/gomail.v2"
	"idea_server/global"
)

func SendMail(subject, body string, tos []string) error {
	emailCfg := global.IDEA_CONFIG.Email
	d := gomail.NewDialer(emailCfg.Host, emailCfg.Port, emailCfg.From, emailCfg.Secret)
	if emailCfg.IsSSL {
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	m := gomail.NewMessage()
	m.SetHeader("From", emailCfg.From)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	m.SetHeader("To", "dparticle@163.com")
	for _, to := range tos {
		m.SetHeader("To", to)
		if err := d.DialAndSend(m) ; err != nil {
			return err
		}
	}
	return nil
}
