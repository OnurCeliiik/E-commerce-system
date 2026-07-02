package email

import (
	"context"
	"fmt"
	"net/smtp"
)

type SMTPConfig struct {
	Host        string
	Port        int
	Username    string
	Password    string
	FromAddress string
	FromName    string
}

type SMTPSender struct {
	config SMTPConfig
}

func NewSMTPSender(config SMTPConfig) *SMTPSender {
	if config.FromAddress == "" {
		config.FromAddress = "noreply@ecommerce.local"
	}
	if config.FromName == "" {
		config.FromName = "E-commerce"
	}
	if config.Port == 0 {
		config.Port = 1025
	}
	return &SMTPSender{config: config}
}

func (s *SMTPSender) Send(_ context.Context, msg Message) error {
	from := fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromAddress)
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	body := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=utf-8\r\n\r\n%s",
		from,
		msg.To,
		msg.Subject,
		msg.Body,
	)

	var auth smtp.Auth
	if s.config.Username != "" {
		auth = smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)
	}

	return smtp.SendMail(addr, auth, s.config.FromAddress, []string{msg.To}, []byte(body))
}
