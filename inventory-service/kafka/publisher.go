package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/OnurCeliiik/ecommerce/services/inventory/dto"
	"github.com/segmentio/kafka-go"
)

const (
	inventoryReservedTopic          = "inventory.reserved"
	inventoryReservationFailedTopic = "inventory.reservation_failed"
)

// InventoryEventPublisher sends inventory outcome events to Kafka.
// It satisfies service.InventoryEventPublisher.
type InventoryEventPublisher struct {
	writer *kafka.Writer
}

func NewInventoryEventPublisher(brokersCSV string) (*InventoryEventPublisher, error) {
	brokers := splitBrokers(brokersCSV)
	if len(brokers) == 0 {
		return nil, fmt.Errorf("KAFKA_BROKERS is not set")
	}

	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Balancer: &kafka.LeastBytes{},
	}

	return &InventoryEventPublisher{writer: writer}, nil
}

func (p *InventoryEventPublisher) PublishInventoryReserved(ctx context.Context, event dto.InventoryReservedEvent) error {
	return p.publish(ctx, inventoryReservedTopic, event.OrderID.String(), event)
}

func (p *InventoryEventPublisher) PublishInventoryReservationFailed(ctx context.Context, event dto.InventoryReservationFailedEvent) error {
	return p.publish(ctx, inventoryReservationFailedTopic, event.OrderID.String(), event)
}

func (p *InventoryEventPublisher) publish(ctx context.Context, topic, key string, event any) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal %s: %w", topic, err)
	}

	msg := kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: payload,
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("write %s: %w", topic, err)
	}

	return nil
}

func (p *InventoryEventPublisher) Close() error {
	return p.writer.Close()
}
