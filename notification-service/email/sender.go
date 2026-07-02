package email

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// Sender delivers notification emails.
type Sender interface {
	Send(ctx context.Context, msg Message) error
}

func NewFromEnv() (Sender, error) {
	host := strings.TrimSpace(os.Getenv("SMTP_HOST"))
	if host == "" {
		log.Println("SMTP_HOST not set, using LogSender")
		return NewLogSender(), nil
	}

	port := 1025
	if portStr := strings.TrimSpace(os.Getenv("SMTP_PORT")); portStr != "" {
		parsed, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("invalid SMTP_PORT %q: %w", portStr, err)
		}
		port = parsed
	}

	log.Printf("using SMTP sender at %s:%d", host, port)
	return NewSMTPSender(SMTPConfig{
		Host:        host,
		Port:        port,
		Username:    os.Getenv("SMTP_USERNAME"),
		Password:    os.Getenv("SMTP_PASSWORD"),
		FromAddress: os.Getenv("SMTP_FROM_ADDRESS"),
		FromName:    os.Getenv("SMTP_FROM_NAME"),
	}), nil
}
