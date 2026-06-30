package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/OnurCeliiik/ecommerce/services/notification/dto"
	"github.com/OnurCeliiik/ecommerce/services/notification/email"
	"github.com/OnurCeliiik/ecommerce/services/notification/metrics"
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

func (s *notificationService) ProcessInventoryReserved(ctx context.Context, event dto.InventoryReservedEvent) error {

	msg := email.Message{
		To:      recipientEmail(event.CustomerEmail, event.UserID),
		Subject: fmt.Sprintf("Order confirmed — %s", event.OrderID),
		Body:    buildOrderConfirmedBody(event),
	}
	if err := s.emailSender.Send(ctx, msg); err != nil {
		metrics.RecordNotificationEvent("inventory_reserved", "error")
		return err
	}
	metrics.RecordNotificationEvent("inventory_reserved", "success")
	return nil
}

func (s *notificationService) ProcessInventoryReservationFailed(ctx context.Context, event dto.InventoryReservationFailedEvent) error {
	msg := email.Message{
		To:      recipientEmail(event.CustomerEmail, event.UserID),
		Subject: fmt.Sprintf("Order could not be fulfilled — %s", event.OrderID),
		Body:    buildOrderFailedBody(event),
	}
	if err := s.emailSender.Send(ctx, msg); err != nil {
		metrics.RecordNotificationEvent("inventory_reservation_failed", "error")
		return err
	}
	metrics.RecordNotificationEvent("inventory_reservation_failed", "success")
	return nil
}

func recipientEmail(customerEmail string, userID uuid.UUID) string {
	if customerEmail != "" {
		return customerEmail
	}
	return fmt.Sprintf("user-%s@stub.local", userID)
}

func buildOrderConfirmedBody(event dto.InventoryReservedEvent) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Good news — your order %s is confirmed.\n", event.OrderID)
	fmt.Fprintf(&b, "Total: %.2f\n\n", event.Total)
	fmt.Fprintf(&b, "Items:\n")
	for _, item := range event.Items {
		fmt.Fprintf(&b, "- %s x%d @ %.2f\n", item.ProductID, item.Quantity, item.UnitPrice)
	}
	return b.String()
}

func buildOrderFailedBody(event dto.InventoryReservationFailedEvent) string {
	return fmt.Sprintf(
		"We're sorry — we could not reserve stock for order %s.\nReason: %s\n",
		event.OrderID,
		event.Reason,
	)
}
