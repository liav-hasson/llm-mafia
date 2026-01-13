package engine

import (
	"context"
	"log"
	"strconv"
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
// ctx must not be nil.
func (tm *TimerManager) SchedulePhaseTimeout(
	currentPhase domain.Phase,
	round int,
	duration time.Duration,
	nextPhase domain.Phase,
	cmdCh chan Command,
	ctx context.Context,
) {
	if ctx == nil {
		log.Printf("[TIMER] ERROR: nil context passed to SchedulePhaseTimeout, skipping timer")
		return
	}

	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Cancel previous phase timer if it exists
	if tm.phaseTimer != nil {
		tm.phaseTimer.Stop()
		tm.phaseTimer = nil
	}

	// Create timer ID for debugging (fixed: was using string(rune(round)) which is wrong, converts numbers to unicode)
	tm.phaseTimerID = currentPhase.String() + "-round-" + strconv.Itoa(round)

	// Schedule new timer
	timerID := tm.phaseTimerID // capture for closure
	tm.phaseTimer = time.AfterFunc(duration, func() {
		// Send phase change command when timer fires
		cmd := &PhaseChangeCommand{NewPhase: nextPhase}

		// Blocking send with context check only
		// If channel is full, we block — dropping phase changes silently is dangerous
		select {
		case cmdCh <- cmd:
			log.Printf("[TIMER] Phase timeout fired: %s, sent PhaseChangeCommand to %s", timerID, nextPhase)
		case <-ctx.Done():
			log.Printf("[TIMER] Phase timeout fired but engine stopped: %s", timerID)
		}
	})

	log.Printf("[TIMER] Scheduled phase timeout: %s, duration=%v, nextPhase=%s", tm.phaseTimerID, duration, nextPhase)
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
// Timeout values are provided from the engine's configuration.
func GetPhaseTimeout(phase domain.Phase, nightTimeout, dayTimeout, votingTimeout time.Duration) time.Duration {
	switch phase {
	case domain.PhaseNight:
		return nightTimeout
	case domain.PhaseDay:
		return dayTimeout
	case domain.PhaseVoting:
		return votingTimeout
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
