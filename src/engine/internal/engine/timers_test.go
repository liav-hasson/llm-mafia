package engine

import (
	"testing"
	"time"

	"mafia-engine/internal/domain"
)

func TestGetPhaseTimeout(t *testing.T) {
	// Test timeout values
	nightTimeout := 2 * time.Minute
	dayTimeout := 5 * time.Minute
	votingTimeout := 1 * time.Minute

	tests := []struct {
		name     string
		phase    domain.Phase
		expected time.Duration
	}{
		{"Night phase", domain.PhaseNight, nightTimeout},
		{"Day phase", domain.PhaseDay, dayTimeout},
		{"Voting phase", domain.PhaseVoting, votingTimeout},
		{"Waiting phase", domain.PhaseWaiting, 0},
		{"Ended phase", domain.PhaseEnded, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetPhaseTimeout(tt.phase, nightTimeout, dayTimeout, votingTimeout)
			if result != tt.expected {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetNextPhase(t *testing.T) {
	tests := []struct {
		current  domain.Phase
		expected domain.Phase
	}{
		{domain.PhaseNight, domain.PhaseDay},
		{domain.PhaseDay, domain.PhaseVoting},
		{domain.PhaseVoting, domain.PhaseNight},
		{domain.PhaseWaiting, domain.PhaseWaiting},
	}

	for _, tt := range tests {
		result := GetNextPhase(tt.current)
		if result != tt.expected {
			t.Errorf("GetNextPhase(%s) = %s, want %s", tt.current, result, tt.expected)
		}
	}
}

func TestTimerManagerCancelBeforeFire(t *testing.T) {
	tm := NewTimerManager()
	cmdCh := make(chan Command, 1)

	tm.SchedulePhaseTimeout(domain.PhaseNight, 1, 50*time.Millisecond, domain.PhaseDay, cmdCh, nil)
	tm.CancelPhaseTimer()

	time.Sleep(100 * time.Millisecond)

	if len(cmdCh) > 0 {
		t.Error("Timer fired after cancellation")
	}
}
