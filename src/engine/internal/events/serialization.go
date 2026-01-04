// Package events defines the event contract for the engine.
//
// Event type strings are stable and must not be runtime-configurable.
// Timestamp fields are always Unix time in milliseconds.
package events

import (
	"encoding/json"
	"fmt"
)

// encodes all Go structs to json
func Marshal(event any) ([]byte, error) {
	return json.Marshal(event)
}

// Unmarshal funcs for all events the engine receives
// returns nil pointer on error for explicit failure signaling
func UnmarshalAllChatMessage(data []byte) (*AllChatMessage, error) {
	var event AllChatMessage
	err := json.Unmarshal(data, &event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func UnmarshalMafiaChatMessage(data []byte) (*MafiaChatMessage, error) {
	var event MafiaChatMessage
	err := json.Unmarshal(data, &event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func UnmarshalVoteSubmitted(data []byte) (*VoteSubmitted, error) {
	var event VoteSubmitted
	err := json.Unmarshal(data, &event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func UnmarshalNightAction(data []byte) (*NightAction, error) {
	var event NightAction
	err := json.Unmarshal(data, &event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func UnmarshalPlayerThoughts(data []byte) (*PlayerThoughts, error) {
	var event PlayerThoughts
	err := json.Unmarshal(data, &event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

// Deserialize takes raw JSON bytes and routes to the appropriate unmarshaler
// based on the "type" field in the JSON. Returns the concrete event struct.
//
// This is the single entry point for converting Kafka message bytes into
// strongly-typed event structs that the engine can work with.
func Deserialize(data []byte) (any, error) {
	// First, parse just the BaseEvent to get the type field
	var base BaseEvent
	if err := json.Unmarshal(data, &base); err != nil {
		return nil, fmt.Errorf("failed to parse event type: %w", err)
	}

	// Route to the appropriate unmarshaler based on type
	switch base.Type {
	case TypeAllChatMessage:
		return UnmarshalAllChatMessage(data)
	case TypeMafiaChatMessage:
		return UnmarshalMafiaChatMessage(data)
	case TypeVoteSubmitted:
		return UnmarshalVoteSubmitted(data)
	case TypeNightAction:
		return UnmarshalNightAction(data)
	case TypePlayerThoughts:
		return UnmarshalPlayerThoughts(data)
	// Engine emits these but doesn't consume them - players do
	case TypeGameStarted, TypePhaseChanged, TypePlayerEliminated, TypeGameEnded, TypeRoleAssigned:
		return nil, fmt.Errorf("engine does not consume event type: %s", base.Type)
	default:
		return nil, fmt.Errorf("unknown event type: %s", base.Type)
	}
}
