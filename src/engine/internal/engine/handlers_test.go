package engine

import (
	"context"
	"testing"
	"time"

	"mafia-engine/internal/domain"
	"mafia-engine/internal/events"
)

func TestHandleEvent_VoteSubmitted(t *testing.T) {
	cmdCh := make(chan Command, 1)
	ctx := context.Background()

	event := &events.VoteSubmitted{
		BaseEvent: events.BaseEvent{GameID: "test", Type: events.TypeVoteSubmitted},
		VoterID:   "p1",
		TargetID:  "p2",
	}

	if err := HandleEvent(ctx, cmdCh, event); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	select {
	case cmd := <-cmdCh:
		voteCmd, ok := cmd.(*VoteCommand)
		if !ok {
			t.Fatalf("expected VoteCommand, got %T", cmd)
		}
		if voteCmd.VoterID != "p1" || voteCmd.TargetID != "p2" {
			t.Error("wrong command data")
		}
	case <-time.After(50 * time.Millisecond):
		t.Fatal("no command sent")
	}
}

func TestHandleEvent_AllChat(t *testing.T) {
	cmdCh := make(chan Command, 1)
	ctx := context.Background()

	event := &events.AllChatMessage{
		BaseEvent: events.BaseEvent{GameID: "test", Type: events.TypeAllChatMessage},
		SenderID:  "p1",
		Message:   "hello",
	}

	if err := HandleEvent(ctx, cmdCh, event); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	select {
	case cmd := <-cmdCh:
		if _, ok := cmd.(*ChatCommand); !ok {
			t.Fatalf("expected ChatCommand, got %T", cmd)
		}
	case <-time.After(50 * time.Millisecond):
		t.Fatal("no command sent")
	}
}

func TestHandleEvent_NightAction(t *testing.T) {
	cmdCh := make(chan Command, 1)
	ctx := context.Background()

	event := &events.NightAction{
		BaseEvent: events.BaseEvent{GameID: "test", Type: events.TypeNightAction},
		Role:      domain.RoleMafia.String(),
		ActorID:   "m1",
		TargetID:  "v1",
	}

	if err := HandleEvent(ctx, cmdCh, event); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	select {
	case cmd := <-cmdCh:
		if _, ok := cmd.(*NightActionCommand); !ok {
			t.Fatalf("expected NightActionCommand, got %T", cmd)
		}
	case <-time.After(50 * time.Millisecond):
		t.Fatal("no command sent")
	}
}

func TestHandleEvent_PlayerThoughts(t *testing.T) {
	cmdCh := make(chan Command, 1)
	ctx := context.Background()

	event := &events.PlayerThoughts{
		BaseEvent: events.BaseEvent{GameID: "test", Type: events.TypePlayerThoughts},
		SenderID:  "p1",
		Thought:   "thinking",
	}

	if err := HandleEvent(ctx, cmdCh, event); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	select {
	case <-cmdCh:
		t.Fatal("PlayerThoughts should not send command")
	case <-time.After(20 * time.Millisecond):
		// Expected: no command
	}
}

func TestHandleEvent_UnknownType(t *testing.T) {
	cmdCh := make(chan Command, 1)
	ctx := context.Background()

	unknownEvent := struct{ Val string }{"test"}
	err := HandleEvent(ctx, cmdCh, unknownEvent)

	if err == nil {
		t.Fatal("expected error for unknown event")
	}
}
