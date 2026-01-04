package domain

import (
	"testing"
)

// TestTallyVotes tests the vote counting function
func TestTallyVotes(t *testing.T) {
	tests := []struct {
		name     string
		votes    map[string]string // voterID -> targetID
		expected map[string]int    // targetID -> vote count
	}{
		{
			name:     "empty votes returns empty tally",
			votes:    map[string]string{},
			expected: map[string]int{},
		},
		{
			name:     "single vote",
			votes:    map[string]string{"voter1": "target1"},
			expected: map[string]int{"target1": 1},
		},
		{
			name:     "two voters same target",
			votes:    map[string]string{"voter1": "target1", "voter2": "target1"},
			expected: map[string]int{"target1": 2},
		},
		{
			name:     "two voters different targets",
			votes:    map[string]string{"voter1": "target1", "voter2": "target2"},
			expected: map[string]int{"target1": 1, "target2": 1},
		},
		{
			name: "multiple voters mixed targets",
			votes: map[string]string{
				"voter1": "target1",
				"voter2": "target1",
				"voter3": "target2",
				"voter4": "target1",
			},
			expected: map[string]int{"target1": 3, "target2": 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TallyVotes(tt.votes)

			// check same number of targets
			if len(result) != len(tt.expected) {
				t.Errorf("got %d targets, expected %d", len(result), len(tt.expected))
				return
			}

			// check each target has correct vote count
			for target, expectedCount := range tt.expected {
				if result[target] != expectedCount {
					t.Errorf("target %s: got %d votes, expected %d", target, result[target], expectedCount)
				}
			}
		})
	}
}

// TestGetTopVoted tests the internal helper that finds player(s) with most votes
// Note: we can test unexported functions because test is in same package
func TestGetTopVoted(t *testing.T) {
	tests := []struct {
		name     string
		votes    map[string]string
		expected []string // can have multiple winners if tie
	}{
		{
			name:     "empty votes returns nil",
			votes:    map[string]string{},
			expected: nil,
		},
		{
			name:     "single vote returns that target",
			votes:    map[string]string{"voter1": "target1"},
			expected: []string{"target1"},
		},
		{
			name:     "clear winner",
			votes:    map[string]string{"voter1": "target1", "voter2": "target1", "voter3": "target2"},
			expected: []string{"target1"},
		},
		{
			name:     "tie returns multiple players",
			votes:    map[string]string{"voter1": "target1", "voter2": "target2"},
			expected: []string{"target1", "target2"}, // order may vary
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getTopVoted(tt.votes)

			// check same number of top voted
			if len(result) != len(tt.expected) {
				t.Errorf("got %d top voted, expected %d", len(result), len(tt.expected))
				return
			}

			// for tie cases, check all expected players are in result
			// (order doesn't matter since map iteration is random)
			for _, exp := range tt.expected {
				found := false
				for _, res := range result {
					if res == exp {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected %s to be in top voted, but wasn't found", exp)
				}
			}
		})
	}
}

// TestGetVoteWinner tests finding a single winner (fails on ties)
func TestGetVoteWinner(t *testing.T) {
	tests := []struct {
		name           string
		votes          map[string]string
		expectedWinner string
		expectedOk     bool
	}{
		{
			name:           "empty votes returns no winner",
			votes:          map[string]string{},
			expectedWinner: "",
			expectedOk:     false,
		},
		{
			name:           "single vote has winner",
			votes:          map[string]string{"voter1": "target1"},
			expectedWinner: "target1",
			expectedOk:     true,
		},
		{
			name:           "clear winner",
			votes:          map[string]string{"voter1": "target1", "voter2": "target1", "voter3": "target2"},
			expectedWinner: "target1",
			expectedOk:     true,
		},
		{
			name:           "tie returns no winner",
			votes:          map[string]string{"voter1": "target1", "voter2": "target2"},
			expectedWinner: "",
			expectedOk:     false,
		},
		{
			name: "three-way tie returns no winner",
			votes: map[string]string{
				"voter1": "target1",
				"voter2": "target2",
				"voter3": "target3",
			},
			expectedWinner: "",
			expectedOk:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			winner, ok := GetVoteWinner(tt.votes)

			if ok != tt.expectedOk {
				t.Errorf("got ok=%v, expected ok=%v", ok, tt.expectedOk)
			}

			if winner != tt.expectedWinner {
				t.Errorf("got winner=%q, expected winner=%q", winner, tt.expectedWinner)
			}
		})
	}
}
