package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/OnurCeliiik/ecommerce/services/notification/dto"
	"github.com/segmentio/kafka-go"
)

const (
	inventoryReservedTopic          = "inventory.reserved"
	inventoryReservationFailedTopic = "inventory.reservation_failed"
)

type InventoryReservedProcessor interface {
	ProcessInventoryReserved(ctx context.Context, event dto.InventoryReservedEvent) error
}

type InventoryReservationFailedProcessor interface {
	ProcessInventoryReservationFailed(ctx context.Context, event dto.InventoryReservationFailedEvent) error
}

type InventoryReservedConsumer struct {
	reader    *kafka.Reader
	processor InventoryReservedProcessor
}

type InventoryReservationFailedConsumer struct {
	reader    *kafka.Reader
	processor InventoryReservationFailedProcessor
}

func NewInventoryReservedConsumer(brokersCSV string, processor InventoryReservedProcessor) (*InventoryReservedConsumer, error) {
	reader, err := newReader(brokersCSV, inventoryReservedTopic, "notification-service")
	if err != nil {
		return nil, err
	}
	return &InventoryReservedConsumer{reader: reader, processor: processor}, nil
}

func NewInventoryReservationFailedConsumer(brokersCSV string, processor InventoryReservationFailedProcessor) (*InventoryReservationFailedConsumer, error) {
	reader, err := newReader(brokersCSV, inventoryReservationFailedTopic, "notification-service")
	if err != nil {
		return nil, err
	}
	return &InventoryReservationFailedConsumer{reader: reader, processor: processor}, nil
}

func (c *InventoryReservedConsumer) Run(ctx context.Context) error {
	defer c.reader.Close()

	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			return fmt.Errorf("read inventory.reserved: %w", err)
		}

		var event dto.InventoryReservedEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("skip invalid inventory.reserved payload: %v", err)
			continue
		}

		if err := c.processor.ProcessInventoryReserved(ctx, event); err != nil {
			log.Printf("process inventory.reserved order_id=%s: %v", event.OrderID, err)
			continue
		}

		log.Printf("processed inventory.reserved order_id=%s", event.OrderID)
	}
}

func (c *InventoryReservationFailedConsumer) Run(ctx context.Context) error {
	defer c.reader.Close()

	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			return fmt.Errorf("read inventory.reservation_failed: %w", err)
		}

		var event dto.InventoryReservationFailedEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("skip invalid inventory.reservation_failed payload: %v", err)
			continue
		}

		if err := c.processor.ProcessInventoryReservationFailed(ctx, event); err != nil {
			log.Printf("process inventory.reservation_failed order_id=%s: %v", event.OrderID, err)
			continue
		}

		log.Printf("processed inventory.reservation_failed order_id=%s", event.OrderID)
	}
}

func newReader(brokersCSV, topic, groupID string) (*kafka.Reader, error) {
	brokers := splitBrokers(brokersCSV)
	if len(brokers) == 0 {
		return nil, fmt.Errorf("KAFKA_BROKERS is not set")
	}

	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	}), nil
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
