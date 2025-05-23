package mailer

import (
	"context"
	"crypto/tls"
	"gopkg.in/gomail.v2"
)

type GOMailer struct {
	dialer *gomail.Dialer
}

func NewGoMailClient(user, pass, host string, port int) GOMailer {
	d := gomail.NewDialer(host, port, user, pass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return GOMailer{dialer: d}
}

func (c GOMailer) SendEmail(ctx context.Context, data SendEmail) error {
	var recipients []string
	if len(data.ToListEmail) > 1 {
		recipients = data.ToListEmail
	} else {
		recipients = append(recipients, data.ToEmail)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", "noreply@system.com")
	m.SetHeader("To", recipients...)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", data.Body)

	return nil
}
