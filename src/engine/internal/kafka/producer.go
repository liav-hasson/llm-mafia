package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

// Producer publishes messages to Kafka.
type Producer interface {
	// context.Context: standard Go practice for handling timeouts and cancellations.
	// If the context is cancelled before the message is sent,
	// the function should return an error and stop the process.
	Publish(ctx context.Context, msg Message) error

	// graceful shutdown of kafka network connection, flushes buffer
	Close() error
}

// KafkaProducer is a concrete implementation of the Producer interface
// using segmentio/kafka-go Writer.
type KafkaProducer struct {
	writer *kafka.Writer
}

// NewKafkaProducer creates a new Kafka producer connected to the given brokers.
// It uses a hash-based partitioner to ensure all messages with the same key
// (e.g., same game ID) go to the same partition, preserving event order.
func NewKafkaProducer(brokers []string, clientID string) (*KafkaProducer, error) {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Balancer: &kafka.Hash{}, // Key-based partitioning for event ordering

		// RequireOne waits for leader acknowledgment (durability vs performance)
		RequiredAcks: kafka.RequireOne,

		// Writer will handle transient failures with retries
		// Synchronous writes ensure events are persisted before returning
		MaxAttempts: 3,

		// No specific Topic - set per message for flexibility
	}

	return &KafkaProducer{writer: writer}, nil
}

// Publish sends a message to Kafka.
// This is a synchronous operation - it waits for the leader to acknowledge.
// The context can be used to set timeouts or cancel the operation.
func (p *KafkaProducer) Publish(ctx context.Context, msg Message) error {
	// Convert our Message to kafka-go's Message format
	kafkaMsg := kafka.Message{
		Topic: msg.Topic,
		Key:   msg.Key,
		Value: msg.Value,
	}

	// WriteMessages is synchronous - waits for ack from Kafka
	return p.writer.WriteMessages(ctx, kafkaMsg)
}

// Close flushes any buffered messages and closes the Kafka connection.
// Should be called during graceful shutdown.
func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}
