package mail

import (
	"bytes"
	"echotalk/internal/config"
	"fmt"
	"html/template"

	"gopkg.in/gomail.v2"
)

func SendVerifyEmail(to string, verifyLink string, location string) error {
	tmpl, err := template.ParseFiles("template.html")
	if err != nil {
		return fmt.Errorf("failed to parse verify email template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, map[string]string{
		"VerifyLink": config.Config.WebURL + "/" + location + "?verify=" + verifyLink,
	}); err != nil {
		return fmt.Errorf("failed to execute verify email template: %w", err)
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", config.Config.MailFrom)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", "[echotalk] 이메일 인증 메일입니다.")
	msg.SetBody("text/html", buf.String())

	if err := gomail.Send(Sender, msg); err != nil {
		return fmt.Errorf("failed to send verify email: %w", err)
	}

	return nil
}
