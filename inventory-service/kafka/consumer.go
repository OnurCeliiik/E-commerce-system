package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/OnurCeliiik/ecommerce/services/inventory/dto"
	"github.com/segmentio/kafka-go"
)

const orderCreatedTopic = "order.created"

// OrderProcessor handles business logic for order events.
// Defined here (kafka package) — implemented by service.inventoryService.
type OrderProcessor interface {
	ProcessOrderCreated(ctx context.Context, event dto.OrderCreatedEvent) error
}

// OrderEventConsumer reads order.created from Kafka and delegates to the service.
type OrderEventConsumer struct {
	reader    *kafka.Reader
	processor OrderProcessor
}

func NewOrderEventConsumer(brokersCSV string, processor OrderProcessor) (*OrderEventConsumer, error) {
	brokers := splitBrokers(brokersCSV)
	if len(brokers) == 0 {
		return nil, fmt.Errorf("KAFKA_BROKERS is not set")
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   orderCreatedTopic,
		GroupID: "inventory-service",
	})

	return &OrderEventConsumer{
		reader:    reader,
		processor: processor,
	}, nil
}

func (c *OrderEventConsumer) Run(ctx context.Context) error {
	defer c.reader.Close()

	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			return fmt.Errorf("read message: %w", err)
		}

		var event dto.OrderCreatedEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("skip invalid order.created payload: %v", err)
			continue
		}

		if err := c.processor.ProcessOrderCreated(ctx, event); err != nil {
			log.Printf("process order.created order_id=%s: %v", event.OrderID, err)
			continue
		}

		log.Printf("processed order.created order_id=%s", event.OrderID)
	}
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
