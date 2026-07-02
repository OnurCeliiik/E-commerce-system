package service

import (
	"context"
	"fmt"
	"log"
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

type ProcessedNotificationRepository interface {
	TryClaim(ctx context.Context, orderID uuid.UUID, eventType string) (bool, error)
}

type notificationService struct {
	emailSender   EmailSender
	processedRepo ProcessedNotificationRepository
}

func NewNotificationService(
	emailSender EmailSender,
	processedRepo ProcessedNotificationRepository,
) *notificationService {
	return &notificationService{
		emailSender:   emailSender,
		processedRepo: processedRepo,
	}
}

func (s *notificationService) ProcessInventoryReserved(ctx context.Context, event dto.InventoryReservedEvent) error {
	const eventType = "inventory_reserved"
	claimed, err := s.processedRepo.TryClaim(ctx, event.OrderID, eventType)
	if err != nil {
		metrics.RecordNotificationEvent(eventType, "error")
		return err
	}
	if !claimed {
		log.Printf("skip duplicate %s order_id=%s", eventType, event.OrderID)
		metrics.RecordNotificationEvent(eventType, "duplicate")
		return nil
	}

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
	const eventType = "inventory_reservation_failed"
	claimed, err := s.processedRepo.TryClaim(ctx, event.OrderID, eventType)
	if err != nil {
		metrics.RecordNotificationEvent(eventType, "error")
		return err
	}
	if !claimed {
		log.Printf("skip duplicate %s order_id=%s", eventType, event.OrderID)
		metrics.RecordNotificationEvent(eventType, "duplicate")
		return nil
	}

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
