package kafka

import (
	"context"

	"github.com/OnurCeliiik/ecommerce/services/order/dto"
)

// NoopPublisher drops events — useful for local runs without Kafka.
type NoopPublisher struct{}

func (NoopPublisher) PublishOrderCreated(context.Context, dto.OrderCreatedEvent) error {
	return nil
}
