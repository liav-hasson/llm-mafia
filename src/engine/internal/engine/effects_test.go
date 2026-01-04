package engine

import (
	"testing"

	"mafia-engine/internal/events"
)

func TestNewPublishEffect(t *testing.T) {
	event := &events.AllChatMessage{
		BaseEvent: events.BaseEvent{GameID: "test", Type: events.TypeAllChatMessage},
		SenderID:  "p1",
		Message:   "hello",
	}

	effect := NewPublishEffect(event)

	if effect.Event != event {
		t.Error("event not stored")
	}
	if effect.Timestamp == 0 {
		t.Error("timestamp not set")
	}
}

func TestInjectTimestamp(t *testing.T) {
	timestamp := int64(1234567890)

	tests := []struct {
		name  string
		event any
	}{
		{"AllChat", &events.AllChatMessage{BaseEvent: events.BaseEvent{GameID: "test"}}},
		{"MafiaChat", &events.MafiaChatMessage{BaseEvent: events.BaseEvent{GameID: "test"}}},
		{"PhaseChanged", &events.PhaseChanged{BaseEvent: events.BaseEvent{GameID: "test"}}},
		{"GameStarted", &events.GameStarted{BaseEvent: events.BaseEvent{GameID: "test"}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := injectTimestamp(tt.event, timestamp); err != nil {
				t.Fatalf("inject failed: %v", err)
			}

			// Verify timestamp was set
			switch e := tt.event.(type) {
			case *events.AllChatMessage:
				if e.Timestamp != timestamp {
					t.Error("timestamp mismatch")
				}
			case *events.MafiaChatMessage:
				if e.Timestamp != timestamp {
					t.Error("timestamp mismatch")
				}
			case *events.PhaseChanged:
				if e.Timestamp != timestamp {
					t.Error("timestamp mismatch")
				}
			case *events.GameStarted:
				if e.Timestamp != timestamp {
					t.Error("timestamp mismatch")
				}
			}
		})
	}
}

func TestExtractGameID(t *testing.T) {
	event := &events.AllChatMessage{
		BaseEvent: events.BaseEvent{GameID: "game-123"},
	}

	gameID, err := extractGameID(event)
	if err != nil {
		t.Fatalf("extract failed: %v", err)
	}
	if gameID != "game-123" {
		t.Errorf("got %s, want game-123", gameID)
	}
}
