package mail

import (
	"bytes"
	"echotalk/internal/config"
	"fmt"

	"gopkg.in/gomail.v2"
)

type Email struct {
	dialer *gomail.Dialer
}

func NewEmail() *Email {
	cfg := config.Config
	dialer := gomail.NewDialer(cfg.SmtpHost, int(cfg.SmtpPort), cfg.MailId, cfg.MailPassword)

	return &Email{
		dialer: dialer,
	}
}

func (email *Email) SendVerifyEmail(to string, verifyLink string, location string) error {
	tmpl, err := loadVerifyEmailTemplate()
	if err != nil {
		return fmt.Errorf("failed to parse verify email template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, map[string]string{
		"VerifyLink": config.Config.WebURL + location + "?verify=" + verifyLink,
	}); err != nil {
		return fmt.Errorf("failed to execute verify email template: %w", err)
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", config.Config.MailFrom)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", "[echotalk] 이메일 인증 메일입니다.")
	msg.SetBody("text/html", buf.String())

	if err := email.dialer.DialAndSend(msg); err != nil {
		return fmt.Errorf("failed to send verify email: %w", err)
	}

	return nil
}
