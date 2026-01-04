package engine

import (
	"strings"
	"testing"

	"mafia-engine/internal/domain"
)

func TestAddPlayerCommand_Success(t *testing.T) {
	state := &domain.GameState{
		Phase:   domain.PhaseWaiting,
		Players: make(map[string]*domain.Player),
	}

	cmd := &AddPlayerCommand{PlayerName: "Alice", MaxPlayers: 10}
	effects, err := cmd.Apply(state)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(effects) != 0 {
		t.Errorf("expected 0 effects, got %d", len(effects))
	}
	if state.GetPlayerCount() != 1 {
		t.Errorf("expected 1 player, got %d", state.GetPlayerCount())
	}
}

func TestAddPlayerCommand_WrongPhase(t *testing.T) {
	state := &domain.GameState{
		Phase:   domain.PhaseNight,
		Players: make(map[string]*domain.Player),
	}

	cmd := &AddPlayerCommand{PlayerName: "Bob", MaxPlayers: 10}
	_, err := cmd.Apply(state)

	if err == nil {
		t.Fatal("expected error for wrong phase")
	}
	if !strings.Contains(err.Error(), "cannot add players") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestStartGameCommand_Success(t *testing.T) {
	state := &domain.GameState{
		ID:      "test-game",
		Phase:   domain.PhaseWaiting,
		Players: make(map[string]*domain.Player),
	}

	// Add 6 players
	for i := 0; i < 6; i++ {
		player, _ := domain.NewPlayer("", "", domain.RoleUnknown)
		state.AddPlayer(player)
	}

	cmd := &StartGameCommand{MinPlayers: 6, MaxPlayers: 12}
	effects, err := cmd.Apply(state)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state.Phase != domain.PhaseNight {
		t.Errorf("expected phase Night, got %s", state.Phase)
	}
	if state.Round != 1 {
		t.Errorf("expected round 1, got %d", state.Round)
	}
	if len(effects) != 8 {
		t.Errorf("expected 8 effects, got %d", len(effects))
	}
}

func TestVoteCommand_Success(t *testing.T) {
	state := &domain.GameState{
		Phase:   domain.PhaseVoting,
		Players: make(map[string]*domain.Player),
		Votes:   make(map[string]string),
	}

	p1, _ := domain.NewPlayer("p1", "Alice", domain.RoleVillager)
	p2, _ := domain.NewPlayer("p2", "Bob", domain.RoleVillager)
	state.AddPlayer(p1)
	state.AddPlayer(p2)

	cmd := &VoteCommand{VoterID: "p1", TargetID: "p2"}
	effects, err := cmd.Apply(state)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(effects) != 0 {
		t.Errorf("expected 0 effects, got %d", len(effects))
	}
	if state.Votes["p1"] != "p2" {
		t.Error("vote not registered")
	}
}
