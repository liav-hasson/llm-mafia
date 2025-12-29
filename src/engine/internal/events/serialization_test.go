package events

import (
	"testing"
)

func TestMarshal(t *testing.T) {
	event := VoteSubmitted{
		BaseEvent: BaseEvent{
			GameID:    "test-game",
			Timestamp: 1234567890,
			Type:      "vote_submitted",
		},
		VoterID:  "player-1",
		TargetID: "player-2",
	}

	data, err := Marshal(event)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	json := string(data)
	expectedFields := []string{
		`"game_id":"test-game"`,
		`"voter":"player-1"`,
		`"target":"player-2"`,
	}

	for _, field := range expectedFields {
		if !contains(json, field) {
			t.Errorf("JSON missing field: %s\nGot: %s", field, json)
		}
	}
}

func TestUnmarshalVoteSubmitted(t *testing.T) {
	tests := []struct {
		name       string
		input      []byte
		wantErr    bool
		wantVoter  string
		wantTarget string
	}{
		{
			name:       "valid json",
			input:      []byte(`{"game_id":"game-1","voter":"player-1","target":"player-2"}`),
			wantErr:    false,
			wantVoter:  "player-1",
			wantTarget: "player-2",
		},
		{
			name:    "invalid json",
			input:   []byte(`not valid json`),
			wantErr: true,
		},
		{
			name:    "empty input",
			input:   []byte{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := UnmarshalVoteSubmitted(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				if result != nil {
					t.Error("expected nil result on error")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result.VoterID != tt.wantVoter {
				t.Errorf("VoterID: got %s, want %s", result.VoterID, tt.wantVoter)
			}
			if result.TargetID != tt.wantTarget {
				t.Errorf("TargetID: got %s, want %s", result.TargetID, tt.wantTarget)
			}
		})
	}
}

func TestMarshalUnmarshalRoundTrip(t *testing.T) {
	original := VoteSubmitted{
		BaseEvent: BaseEvent{
			GameID:    "round-trip-game",
			Timestamp: 9999,
			Type:      "vote_submitted",
		},
		VoterID:  "voter-abc",
		TargetID: "target-xyz",
	}

	data, err := Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	result, err := UnmarshalVoteSubmitted(data)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if result.GameID != original.GameID {
		t.Errorf("GameID: got %s, want %s", result.GameID, original.GameID)
	}
	if result.VoterID != original.VoterID {
		t.Errorf("VoterID: got %s, want %s", result.VoterID, original.VoterID)
	}
	if result.TargetID != original.TargetID {
		t.Errorf("TargetID: got %s, want %s", result.TargetID, original.TargetID)
	}
}

// helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
