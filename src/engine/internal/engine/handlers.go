package engine

import (
	"context"
	"fmt"

	"mafia-engine/internal/events"
)

// HandleEvent converts external events into internal commands.
// It performs type assertion on the event and creates the appropriate Command.
// Commands are sent to cmdCh for execution in the run loop.
//
// This function should be fast and non-blocking - it just validates and routes.
// All actual state mutation happens in Command.Apply().
func HandleEvent(ctx context.Context, cmdCh chan Command, ev any) error {
	// Type switch on the event to create appropriate commands
	switch e := ev.(type) {

	case *events.VoteSubmitted:
		// Create command from event data
		cmd := &VoteCommand{
			VoterID:  e.VoterID,
			TargetID: e.TargetID,
		}

		// Send to command channel
		// Note: This is non-blocking because cmdCh is buffered (size 64)
		select {
		case cmdCh <- cmd:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}

	case *events.AllChatMessage:
		// Create command from event data
		cmd := &ChatCommand{
			SenderID: e.SenderID,
			Message:  e.Message,
		}

		// Send to command channel
		select {
		case cmdCh <- cmd:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}

	case *events.NightAction:
		// Create command from event data
		cmd := &NightActionCommand{
			Role:     e.Role,
			ActorID:  e.ActorID,
			TargetID: e.TargetID,
		}

		// Send to command channel
		select {
		case cmdCh <- cmd:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}

	case *events.MafiaChatMessage:
		// Create mafia-specific chat command
		cmd := &MafiaChatCommand{
			SenderID: e.SenderID,
			Message:  e.Message,
		}

		// Send to command channel
		select {
		case cmdCh <- cmd:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}

	case *events.PlayerThoughts:
		// Player thoughts don't mutate game state
		// They're for AI agent reasoning/debugging
		// We could log them or emit them back to Kafka for observability
		// For now, we just acknowledge and ignore (no state change needed)
		return nil

	default:
		// Unknown event type - this should never happen if Deserialize is correct
		return fmt.Errorf("unknown event type: %T", ev)
	}
}
