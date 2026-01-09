package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

// HandlerFunc processes a single Kafka message.
// Returning an error means the message was NOT processed successfully.
type HandlerFunc func(ctx context.Context, msg Message) error

// Consumer subscribes to topics and delivers messages to a handler.
type Consumer interface {
	// Consume starts consuming messages and blocks until context is canceled
	Consume(ctx context.Context, handler HandlerFunc) error
	Close() error
}

// KafkaConsumer is a concrete implementation of the Consumer interface
// using segmentio/kafka-go Reader.
type KafkaConsumer struct {
	reader *kafka.Reader
}

// NewKafkaConsumer creates a new Kafka consumer subscribed to the given topic.
// It uses consumer groups for scalability and automatic partition assignment.
// Multiple consumers with the same groupID will share the work.
func NewKafkaConsumer(brokers []string, topic string, groupID string) (*KafkaConsumer, error) {
	if topic == "" {
		return nil, fmt.Errorf("topic is required")
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: groupID,
		Topic:   topic,

		// Start reading from the earliest unread message
		// If consumer group has existing offset, it will resume from there
		StartOffset: kafka.FirstOffset,

		// Commit offsets automatically after successful handler execution
		// This ensures at-least-once delivery semantics
		CommitInterval: 0, // We'll commit manually for better control

		// MinBytes and MaxBytes control fetch size
		MinBytes: 1,    // Fetch immediately if any data available
		MaxBytes: 10e6, // 10MB max per fetch

		// MaxWait limits how long to wait for MinBytes
		// If timeout expires, return whatever data is available
		// MaxWait: 500 * time.Millisecond, // Optional: reduce latency
	})

	return &KafkaConsumer{reader: reader}, nil
}

// Consume starts consuming messages and blocks until context is canceled.
// For each message:
//  1. Fetch message from Kafka
//  2. Call handler with the message
//  3. If handler succeeds, commit offset
//  4. If handler fails, log error and continue (message will be reprocessed on restart)
func (c *KafkaConsumer) Consume(ctx context.Context, handler HandlerFunc) error {
	for {
		// Check if context is cancelled before fetching
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Fetch next message (blocks until message available or context cancelled)
		kafkaMsg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			// Context cancelled or connection error
			if ctx.Err() != nil {
				return ctx.Err()
			}
			// Log error and continue trying
			// TODO: Add proper logging
			continue
		}

		// Convert kafka-go Message to our Message format
		msg := Message{
			Topic: kafkaMsg.Topic,
			Key:   kafkaMsg.Key,
			Value: kafkaMsg.Value,
		}

		// Call handler
		if err := handler(ctx, msg); err != nil {
			// Handler failed - message will be reprocessed on restart
			// TODO: Add proper logging and metrics
			// TODO: Consider dead letter queue for poison messages
			_ = err
			continue
		}

		// Handler succeeded - commit offset
		if err := c.reader.CommitMessages(ctx, kafkaMsg); err != nil {
			// Commit failed - message might be reprocessed (at-least-once semantics)
			// TODO: Add proper logging
			_ = err
		}
	}
}

// Close stops consuming and closes the Kafka connection.
// Should be called during graceful shutdown.
func (c *KafkaConsumer) Close() error {
	return c.reader.Close()
}
