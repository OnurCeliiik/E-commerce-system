package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/OnurCeliiik/ecommerce/services/notification/dto"
	"github.com/OnurCeliiik/ecommerce/services/notification/email"
	"github.com/google/uuid"
)

// EmailSender delivers notification emails. Implementations live in the email package.
type EmailSender interface {
	Send(ctx context.Context, msg email.Message) error
}

type notificationService struct {
	emailSender EmailSender
}

func NewNotificationService(emailSender EmailSender) *notificationService {
	return &notificationService{emailSender: emailSender}
}

func (s *notificationService) ProcessOrderCreated(ctx context.Context, event dto.OrderCreatedEvent) error {
	msg := email.Message{
		To:      recipientForUser(event.UserID),
		Subject: fmt.Sprintf("Order confirmed — %s", event.OrderID),
		Body:    buildOrderConfirmationBody(event),
	}
	return s.emailSender.Send(ctx, msg)
}

func recipientForUser(userID uuid.UUID) string {
	// Slice 1: we only have user_id on the event, not the email address.
	// Replace with a user-service lookup or customer_email on the event later.
	return fmt.Sprintf("user-%s@stub.local", userID)
}

func buildOrderConfirmationBody(event dto.OrderCreatedEvent) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Your order %s has been placed.\n", event.OrderID)
	fmt.Fprintf(&b, "Status: %s\n", event.Status)
	fmt.Fprintf(&b, "Total: %.2f\n\n", event.Total)
	fmt.Fprintf(&b, "Items:\n")
	for _, item := range event.Items {
		fmt.Fprintf(&b, "- %s x%d @ %.2f\n", item.ProductID, item.Quantity, item.UnitPrice)
	}
	return b.String()
}
