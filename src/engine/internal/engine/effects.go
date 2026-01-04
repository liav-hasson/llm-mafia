package engine

import (
	"context"
	"fmt"
	"time"

	"mafia-engine/internal/events"
	"mafia-engine/internal/kafka"
)

// PublishEffect represents a Kafka event that should be published.
// This is the primary side effect - publishing authoritative game events.
type PublishEffect struct {
	// Event is the event struct to publish (must have BaseEvent embedded)
	Event any

	// Timestamp is set by the engine when creating the effect
	// Commands must NOT set this - engine provides deterministic timestamps
	Timestamp int64
}

// Execute implements the Effect interface.
// It marshals the event to JSON and publishes to Kafka.
func (e *PublishEffect) Execute(ctx context.Context, producer kafka.Producer) error {
	// Inject timestamp into the event's BaseEvent
	// This is a bit tricky because Event is any - we need type assertion
	if err := injectTimestamp(e.Event, e.Timestamp); err != nil {
		return fmt.Errorf("failed to inject timestamp: %w", err)
	}

	// Marshal event to JSON
	eventBytes, err := events.Marshal(e.Event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Extract GameID from event for partitioning
	gameID, err := extractGameID(e.Event)
	if err != nil {
		return fmt.Errorf("failed to extract game ID: %w", err)
	}

	// Create Kafka message
	msg := kafka.Message{
		Topic: kafka.EngineEventsTopic,
		Key:   kafka.GameKey(gameID),
		Value: eventBytes,
	}

	// Publish to Kafka
	if err := producer.Publish(ctx, msg); err != nil {
		return fmt.Errorf("failed to publish to kafka: %w", err)
	}

	return nil
}

// Helper: Inject timestamp into event's BaseEvent field
// This uses reflection-like patterns but we can optimize with type switches
func injectTimestamp(event any, timestamp int64) error {
	// Type switch on known event types
	switch e := event.(type) {
	case *events.AllChatMessage:
		e.Timestamp = timestamp
	case *events.MafiaChatMessage:
		e.Timestamp = timestamp
	case *events.PhaseChanged:
		e.Timestamp = timestamp
	case *events.PlayerEliminated:
		e.Timestamp = timestamp
	case *events.GameEnded:
		e.Timestamp = timestamp
	case *events.GameStarted:
		e.Timestamp = timestamp
	case *events.NightAction:
		e.Timestamp = timestamp
	case *events.RoleAssigned:
		e.Timestamp = timestamp
	case *events.VoteSubmitted:
		e.Timestamp = timestamp
	case *events.PlayerThoughts:
		e.Timestamp = timestamp
	default:
		return fmt.Errorf("unknown event type: %T", event)
	}
	return nil
}

// Helper: Extract GameID from event's BaseEvent field
func extractGameID(event any) (string, error) {
	// Type switch to extract GameID
	switch e := event.(type) {
	case *events.AllChatMessage:
		return e.GameID, nil
	case *events.MafiaChatMessage:
		return e.GameID, nil
	case *events.PhaseChanged:
		return e.GameID, nil
	case *events.PlayerEliminated:
		return e.GameID, nil
	case *events.GameEnded:
		return e.GameID, nil
	case *events.GameStarted:
		return e.GameID, nil
	case *events.NightAction:
		return e.GameID, nil
	case *events.RoleAssigned:
		return e.GameID, nil
	case *events.VoteSubmitted:
		return e.GameID, nil
	case *events.PlayerThoughts:
		return e.GameID, nil
	default:
		return "", fmt.Errorf("unknown event type: %T", event)
	}
}

// NewPublishEffect creates a PublishEffect with the current timestamp.
// The timestamp is injected by the engine, not by commands.
func NewPublishEffect(event any) *PublishEffect {
	return &PublishEffect{
		Event:     event,
		Timestamp: time.Now().UnixMilli(),
	}
}

// TimerEffect schedules a command to execute after a delay.
// Example: "After 5 minutes, advance to next phase"
// NOTE: This requires access to cmdCh, which is passed during effect creation.
type TimerEffect struct {
	Delay   time.Duration
	Command Command
	CmdCh   chan Command // Target channel to send command to after delay
}

func (e *TimerEffect) Execute(ctx context.Context, producer kafka.Producer) error {
	// Schedule the command to be sent after delay
	// time.AfterFunc returns a Timer that can be stopped if needed
	timer := time.AfterFunc(e.Delay, func() {
		// Non-blocking send with context check
		select {
		case e.CmdCh <- e.Command:
			// Command successfully scheduled
		case <-ctx.Done():
			// Engine stopped before timer fired - ignore
		}
	})

	// Store timer reference if we need to cancel it later
	// For now, we just let it run
	_ = timer

	return nil
}

// LogEffect logs a message (useful for debugging without Kafka)
type LogEffect struct {
	Message string
	Level   string // "info", "warn", "error"
}

func (e *LogEffect) Execute(ctx context.Context, producer kafka.Producer) error {
	// TODO: Integrate with proper logging library
	// For now, just print
	fmt.Printf("[%s] %s\n", e.Level, e.Message)
	return nil
}
