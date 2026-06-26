package email

import (
	"context"
	"log"
)

// LogSender prints emails to the application log (development / learning stub).
type LogSender struct{}

func NewLogSender() *LogSender {
	return &LogSender{}
}

func (LogSender) Send(_ context.Context, msg Message) error {
	log.Printf("[EMAIL] to=%s subject=%q\n%s", msg.To, msg.Subject, msg.Body)
	return nil
}
