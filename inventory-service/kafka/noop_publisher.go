package kafka

import (
	"context"

	"github.com/OnurCeliiik/ecommerce/services/inventory/dto"
)

// NoopPublisher drops events — useful for local runs without Kafka.
type NoopPublisher struct{}

func (NoopPublisher) PublishInventoryReserved(context.Context, dto.InventoryReservedEvent) error {
	return nil
}

func (NoopPublisher) PublishInventoryReservationFailed(context.Context, dto.InventoryReservationFailedEvent) error {
	return nil
}
