package engine

import (
	"context"
	"sync"
	"time"

	"mafia-engine/internal/domain"
)

// TimerManager tracks and manages phase timeout timers.
// It ensures only one phase timer is active at a time and provides
// cancellation support for graceful shutdown and manual phase changes.
type TimerManager struct {
	mu           sync.Mutex
	phaseTimer   *time.Timer // current phase timer (nil if none active)
	phaseTimerID string      // identifier for debugging (e.g., "night-round-2")
}

// NewTimerManager creates a new TimerManager with no active timers.
func NewTimerManager() *TimerManager {
	return &TimerManager{}
}

// SchedulePhaseTimeout schedules a timer to automatically advance to the next phase.
// If a previous phase timer exists, it is cancelled first.
// When the timer fires, it sends a PhaseChangeCommand to the command channel.
func (tm *TimerManager) SchedulePhaseTimeout(
	currentPhase domain.Phase,
	round int,
	duration time.Duration,
	nextPhase domain.Phase,
	cmdCh chan Command,
	ctx context.Context,
) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Cancel previous phase timer if it exists
	if tm.phaseTimer != nil {
		tm.phaseTimer.Stop()
		tm.phaseTimer = nil
	}

	// Create timer ID for debugging
	tm.phaseTimerID = currentPhase.String() + "-round-" + string(rune(round))

	// Schedule new timer
	tm.phaseTimer = time.AfterFunc(duration, func() {
		// Send phase change command when timer fires
		cmd := &PhaseChangeCommand{NewPhase: nextPhase}

		// Non-blocking send with context check
		select {
		case cmdCh <- cmd:
			// Timer fired successfully, command sent
		case <-ctx.Done():
			// Engine stopped, ignore
		default:
			// Command channel full (should not happen with buffered channel)
			// TODO: Add logging/metrics for this edge case
		}
	})
}

// CancelPhaseTimer stops the current phase timer if one is active.
// This should be called when a phase changes manually (before the timer expires).
// It is safe to call even if no timer is active.
func (tm *TimerManager) CancelPhaseTimer() {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if tm.phaseTimer != nil {
		tm.phaseTimer.Stop()
		tm.phaseTimer = nil
		tm.phaseTimerID = ""
	}
}

// Shutdown stops all active timers.
// This should be called during engine shutdown.
func (tm *TimerManager) Shutdown() {
	tm.CancelPhaseTimer()
}

// GetPhaseTimeout returns the timeout duration for a given phase.
// These are the default durations - can be made configurable later.
func GetPhaseTimeout(phase domain.Phase) time.Duration {
	switch phase {
	case domain.PhaseNight:
		return 2 * time.Minute // Mafia coordination time
	case domain.PhaseDay:
		return 5 * time.Minute // Discussion time
	case domain.PhaseVoting:
		return 1 * time.Minute // Voting deadline
	default:
		return 0 // No timeout for Waiting/Ended phases
	}
}

// GetNextPhase returns the next phase in the game cycle.
// Night -> Day -> Voting -> Night (with round increment)
func GetNextPhase(current domain.Phase) domain.Phase {
	switch current {
	case domain.PhaseNight:
		return domain.PhaseDay
	case domain.PhaseDay:
		return domain.PhaseVoting
	case domain.PhaseVoting:
		return domain.PhaseNight
	default:
		return domain.PhaseWaiting
	}
}
