package service_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/OnurCeliiik/ecommerce/services/notification/dto"
	"github.com/OnurCeliiik/ecommerce/services/notification/email"
	"github.com/OnurCeliiik/ecommerce/services/notification/service"
	"github.com/google/uuid"
)

type mockEmailSender struct {
	send func(ctx context.Context, msg email.Message) error
}

func (m *mockEmailSender) Send(ctx context.Context, msg email.Message) error {
	return m.send(ctx, msg)
}

type mockProcessedNotificationRepo struct {
	tryClaim func(ctx context.Context, orderID uuid.UUID, eventType string) (bool, error)
}

func (m *mockProcessedNotificationRepo) TryClaim(ctx context.Context, orderID uuid.UUID, eventType string) (bool, error) {
	return m.tryClaim(ctx, orderID, eventType)
}

func TestProcessInventoryReserved_Success(t *testing.T) {
	orderID := uuid.New()
	userID := uuid.New()

	var claimedOrderID uuid.UUID
	var claimedEventType string
	var sentMsg email.Message

	repo := &mockProcessedNotificationRepo{
		tryClaim: func(_ context.Context, id uuid.UUID, eventType string) (bool, error) {
			claimedOrderID = id
			claimedEventType = eventType
			return true, nil
		},
	}
	sender := &mockEmailSender{
		send: func(_ context.Context, msg email.Message) error {
			sentMsg = msg
			return nil
		},
	}

	svc := service.NewNotificationService(sender, repo)
	event := dto.InventoryReservedEvent{
		OrderID:       orderID,
		UserID:        userID,
		CustomerEmail: "test@example.com",
		Total:         100,
		Items: []dto.OrderLineItem{
			{ProductID: uuid.New(), Quantity: 1, UnitPrice: 100},
		},
	}

	err := svc.ProcessInventoryReserved(context.Background(), event)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if claimedOrderID != orderID {
		t.Fatalf("expected TryClaim order %s, got %s", orderID, claimedOrderID)
	}
	if claimedEventType != "inventory_reserved" {
		t.Fatalf("expected event type inventory_reserved, got %s", claimedEventType)
	}
	if sentMsg.To != "test@example.com" {
		t.Fatalf("expected email to test@example.com, got %s", sentMsg.To)
	}
	if !strings.Contains(sentMsg.Subject, orderID.String()) {
		t.Fatalf("expected subject to contain order id, got %q", sentMsg.Subject)
	}
	if !strings.Contains(sentMsg.Body, "confirmed") {
		t.Fatalf("expected confirmation body, got %q", sentMsg.Body)
	}
}

func TestProcessInventoryReserved_Duplicate(t *testing.T) {
	sendCalled := false

	repo := &mockProcessedNotificationRepo{
		tryClaim: func(_ context.Context, _ uuid.UUID, _ string) (bool, error) {
			return false, nil
		},
	}
	sender := &mockEmailSender{
		send: func(_ context.Context, _ email.Message) error {
			sendCalled = true
			return nil
		},
	}

	svc := service.NewNotificationService(sender, repo)
	err := svc.ProcessInventoryReserved(context.Background(), dto.InventoryReservedEvent{
		OrderID: uuid.New(),
		UserID:  uuid.New(),
	})
	if err != nil {
		t.Fatalf("expected no error on duplicate, got %v", err)
	}
	if sendCalled {
		t.Fatal("expected Send not to be called for duplicate event")
	}
}

func TestProcessInventoryReserved_ClaimError(t *testing.T) {
	wantErr := errors.New("db down")

	repo := &mockProcessedNotificationRepo{
		tryClaim: func(_ context.Context, _ uuid.UUID, _ string) (bool, error) {
			return false, wantErr
		},
	}
	sender := &mockEmailSender{
		send: func(_ context.Context, _ email.Message) error {
			t.Fatal("Send should not be called when TryClaim fails")
			return nil
		},
	}

	svc := service.NewNotificationService(sender, repo)
	err := svc.ProcessInventoryReserved(context.Background(), dto.InventoryReservedEvent{
		OrderID: uuid.New(),
		UserID:  uuid.New(),
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected claim error, got %v", err)
	}
}

func TestProcessInventoryReserved_SendError(t *testing.T) {
	wantErr := errors.New("smtp failed")

	repo := &mockProcessedNotificationRepo{
		tryClaim: func(_ context.Context, _ uuid.UUID, _ string) (bool, error) {
			return true, nil
		},
	}
	sender := &mockEmailSender{
		send: func(_ context.Context, _ email.Message) error {
			return wantErr
		},
	}

	svc := service.NewNotificationService(sender, repo)
	err := svc.ProcessInventoryReserved(context.Background(), dto.InventoryReservedEvent{
		OrderID:       uuid.New(),
		UserID:        uuid.New(),
		CustomerEmail: "test@example.com",
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected send error, got %v", err)
	}
}

func TestProcessInventoryReserved_UsesStubEmailWhenCustomerEmailEmpty(t *testing.T) {
	userID := uuid.New()
	var sentTo string

	sender := &mockEmailSender{
		send: func(_ context.Context, msg email.Message) error {
			sentTo = msg.To
			return nil
		},
	}
	repo := &mockProcessedNotificationRepo{
		tryClaim: func(_ context.Context, _ uuid.UUID, _ string) (bool, error) {
			return true, nil
		},
	}

	svc := service.NewNotificationService(sender, repo)
	err := svc.ProcessInventoryReserved(context.Background(), dto.InventoryReservedEvent{
		OrderID: uuid.New(),
		UserID:  userID,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	want := "user-" + userID.String() + "@stub.local"
	if sentTo != want {
		t.Fatalf("expected stub email %s, got %s", want, sentTo)
	}
}

func TestProcessInventoryReservationFailed_Success(t *testing.T) {
	orderID := uuid.New()
	var claimedEventType string
	var sentMsg email.Message

	repo := &mockProcessedNotificationRepo{
		tryClaim: func(_ context.Context, id uuid.UUID, eventType string) (bool, error) {
			if id != orderID {
				t.Fatalf("expected order %s, got %s", orderID, id)
			}
			claimedEventType = eventType
			return true, nil
		},
	}
	sender := &mockEmailSender{
		send: func(_ context.Context, msg email.Message) error {
			sentMsg = msg
			return nil
		},
	}

	svc := service.NewNotificationService(sender, repo)
	err := svc.ProcessInventoryReservationFailed(context.Background(), dto.InventoryReservationFailedEvent{
		OrderID:       orderID,
		UserID:        uuid.New(),
		CustomerEmail: "fail@example.com",
		Reason:        "insufficient_inventory",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if claimedEventType != "inventory_reservation_failed" {
		t.Fatalf("expected event type inventory_reservation_failed, got %s", claimedEventType)
	}
	if sentMsg.To != "fail@example.com" {
		t.Fatalf("expected email to fail@example.com, got %s", sentMsg.To)
	}
	if !strings.Contains(sentMsg.Body, "insufficient_inventory") {
		t.Fatalf("expected reason in body, got %q", sentMsg.Body)
	}
}

func TestProcessInventoryReservationFailed_Duplicate(t *testing.T) {
	sendCalled := false

	repo := &mockProcessedNotificationRepo{
		tryClaim: func(_ context.Context, _ uuid.UUID, _ string) (bool, error) {
			return false, nil
		},
	}
	sender := &mockEmailSender{
		send: func(_ context.Context, _ email.Message) error {
			sendCalled = true
			return nil
		},
	}

	svc := service.NewNotificationService(sender, repo)
	err := svc.ProcessInventoryReservationFailed(context.Background(), dto.InventoryReservationFailedEvent{
		OrderID: uuid.New(),
		UserID:  uuid.New(),
		Reason:  "insufficient_inventory",
	})
	if err != nil {
		t.Fatalf("expected no error on duplicate, got %v", err)
	}
	if sendCalled {
		t.Fatal("expected Send not to be called for duplicate event")
	}
}
