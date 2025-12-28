package domain

import "testing"

// TestPhaseString tests the String() method for all Phase values
func TestPhaseString(t *testing.T) {
	tests := []struct {
		phase    Phase
		expected string
	}{
		{PhaseUnknown, "unknown"},
		{PhaseWaiting, "waiting"},
		{PhaseNight, "night"},
		{PhaseDay, "day"},
		{PhaseVoting, "voting"},
		{PhaseEnded, "ended"},
		{Phase(99), "invalid"}, // unknown phase value
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.phase.String()
			if result != tt.expected {
				t.Errorf("got %s, expected %s", result, tt.expected)
			}
		})
	}
}

// TestPhaseIotaValues verifies the iota values are as expected
func TestPhaseIotaValues(t *testing.T) {
	// This test documents the expected int values
	// Important if these are stored in a database or sent over network
	tests := []struct {
		phase         Phase
		expectedValue int
	}{
		{PhaseUnknown, 0},
		{PhaseWaiting, 1},
		{PhaseNight, 2},
		{PhaseDay, 3},
		{PhaseVoting, 4},
		{PhaseEnded, 5},
	}

	for _, tt := range tests {
		t.Run(tt.phase.String(), func(t *testing.T) {
			if int(tt.phase) != tt.expectedValue {
				t.Errorf("Phase %s: got value %d, expected %d",
					tt.phase, int(tt.phase), tt.expectedValue)
			}
		})
	}
}
