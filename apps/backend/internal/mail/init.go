package mail

import (
	"echotalk/internal/config"

	"gopkg.in/gomail.v2"
)

var Sender gomail.SendCloser

func init() {
	cfg := config.Config
	dialer := gomail.NewDialer(cfg.SmtpHost, int(cfg.SmtpPort), cfg.MailId, cfg.MailPassword)

	sender, err := dialer.Dial()
	if err != nil {
		panic("failed to connect to SMTP server: " + err.Error())
	}

	Sender = sender
}
