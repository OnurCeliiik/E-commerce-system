package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/OnurCeliiik/ecommerce/services/order/dto"
	"github.com/segmentio/kafka-go"
)

// The topic names for events.
const orderCreatedTopic = "order.created"

// OrderEventPublisher sends domain events to Kafka.
// It satisfies service.OrderEventPublisher.
type OrderEventPublisher struct {
	writer *kafka.Writer
}

func NewOrderEventPublisher(brokersCSV string) (*OrderEventPublisher, error) {
	brokers := splitBrokers(brokersCSV)
	if len(brokers) == 0 {
		return nil, fmt.Errorf("KAFKA_BROKERS is not set")
	}

	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    orderCreatedTopic,
		Balancer: &kafka.LeastBytes{},
	}

	return &OrderEventPublisher{writer: writer}, nil
}

func (p *OrderEventPublisher) PublishOrderCreated(ctx context.Context, event dto.OrderCreatedEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal order.created: %w", err)
	}

	// Key = order_id so all events for one order land in the same partition.
	msg := kafka.Message{
		Key:   []byte(event.OrderID.String()),
		Value: payload,
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("write order.created: %w", err)
	}

	return nil
}

func (p *OrderEventPublisher) Close() error {
	return p.writer.Close()
}

func splitBrokers(brokersCSV string) []string {
	parts := strings.Split(brokersCSV, ",")
	brokers := make([]string, 0, len(parts))
	for _, part := range parts {
		broker := strings.TrimSpace(part)
		if broker != "" {
			brokers = append(brokers, broker)
		}
	}
	return brokers
}
