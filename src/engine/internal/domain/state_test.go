package domain

import (
	"fmt"
	"testing"
)

// --- Helper function for tests ---

// createTestGame creates a game with n players for testing
// All players start with RoleUnknown
func createTestGame(n int) *GameState {
	ResetPlayerCounter() // ensure clean state
	game := NewGameState("test")

	for i := 0; i < n; i++ {
		id := CreatePlayerID()
		name := fmt.Sprintf("TestPlayer%d", i+1)
		player, _ := NewPlayer(id, name, RoleUnknown)
		game.AddPlayer(player)
	}

	return game
}

// --- NewGameState Tests ---

func TestNewGameState(t *testing.T) {
	game := NewGameState("test")

	if game.ID == "" {
		t.Error("game ID should not be empty")
	}
	if game.Round != 1 {
		t.Errorf("initial round: got %d, expected 1", game.Round)
	}
	if game.Phase != PhaseWaiting {
		t.Errorf("initial phase: got %v, expected PhaseWaiting", game.Phase)
	}
	if game.Winner != WinnerNone {
		t.Errorf("initial winner: got %v, expected WinnerNone", game.Winner)
	}
	if game.Players == nil {
		t.Error("Players map should be initialized")
	}
	if game.Votes == nil {
		t.Error("Votes map should be initialized")
	}
}

func TestCreateGameID(t *testing.T) {
	id1 := CreateGameID("test")
	id2 := CreateGameID("prod")

	// IDs should be prefix + "-" + 5 characters = 10 total for "test"
	if len(id1) != 10 { // "test-" (5) + random (5)
		t.Errorf("game ID length: got %d, expected 10", len(id1))
	}

	// IDs should start with prefix
	if id1[:5] != "test-" {
		t.Errorf("game ID should start with 'test-', got %s", id1)
	}

	if id2[:5] != "prod-" {
		t.Errorf("game ID should start with 'prod-', got %s", id2)
	}

	// IDs with same prefix should have different random suffixes
	id3 := CreateGameID("test")
	if id1 == id3 {
		t.Error("two game IDs with same prefix should have different suffixes")
	}
}

// --- Winner Tests ---

func TestWinnerString(t *testing.T) {
	tests := []struct {
		winner   Winner
		expected string
	}{
		{WinnerNone, "none"},
		{WinnerMafia, "mafia"},
		{WinnerVillage, "village"},
		{Winner(99), "invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.winner.String()
			if result != tt.expected {
				t.Errorf("got %s, expected %s", result, tt.expected)
			}
		})
	}
}

// --- AddPlayer Tests ---

func TestAddPlayer(t *testing.T) {
	game := NewGameState("test")
	player := &Player{ID: "test-1", Name: "Test Player", Role: RoleVillager, Alive: true}

	result := game.AddPlayer(player)

	if result == nil {
		t.Fatal("AddPlayer should return the added player")
	}
	if game.GetPlayerCount() != 1 {
		t.Errorf("player count: got %d, expected 1", game.GetPlayerCount())
	}
	if game.GetPlayer("test-1") != player {
		t.Error("should be able to retrieve added player")
	}
}

func TestAddPlayer_RejectsDuplicate(t *testing.T) {
	game := NewGameState("test")
	player1 := &Player{ID: "test-1", Name: "First", Role: RoleVillager, Alive: true}
	player2 := &Player{ID: "test-1", Name: "Duplicate", Role: RoleMafia, Alive: true}

	game.AddPlayer(player1)
	result := game.AddPlayer(player2)

	if result != nil {
		t.Error("AddPlayer should reject duplicate ID")
	}
	if game.GetPlayerCount() != 1 {
		t.Errorf("player count: got %d, expected 1", game.GetPlayerCount())
	}
}

// --- GetPlayer Tests ---

func TestGetPlayer(t *testing.T) {
	game := NewGameState("test")
	player := &Player{ID: "test-1", Name: "Test", Role: RoleVillager, Alive: true}
	game.AddPlayer(player)

	result := game.GetPlayer("test-1")
	if result != player {
		t.Error("GetPlayer should return the player")
	}

	notFound := game.GetPlayer("nonexistent")
	if notFound != nil {
		t.Error("GetPlayer should return nil for nonexistent player")
	}
}

// --- GetAlivePlayers Tests ---

func TestGetAlivePlayers(t *testing.T) {
	game := NewGameState("test")
	game.AddPlayer(&Player{ID: "1", Name: "Alive1", Alive: true})
	game.AddPlayer(&Player{ID: "2", Name: "Dead", Alive: false})
	game.AddPlayer(&Player{ID: "3", Name: "Alive2", Alive: true})

	alive := game.GetAlivePlayers()

	if len(alive) != 2 {
		t.Errorf("alive count: got %d, expected 2", len(alive))
	}
}

func TestGetAlivePlayers_EmptyGame(t *testing.T) {
	game := NewGameState("test")
	alive := game.GetAlivePlayers()

	if len(alive) != 0 {
		t.Errorf("empty game should return empty slice, got %v", alive)
	}
}

// --- GetPlayersByRole Tests ---

// --- EliminatePlayer Tests ---

func TestEliminatePlayer(t *testing.T) {
	game := NewGameState("test")
	player := &Player{ID: "test-1", Name: "Test", Alive: true}
	game.AddPlayer(player)

	result := game.EliminatePlayer("test-1")

	if result == nil {
		t.Fatal("EliminatePlayer should return the eliminated player")
	}
	if player.Alive {
		t.Error("player should be marked as dead")
	}
}

func TestEliminatePlayer_NonexistentPlayer(t *testing.T) {
	game := NewGameState("test")

	result := game.EliminatePlayer("nonexistent")

	if result != nil {
		t.Error("EliminatePlayer should return nil for nonexistent player")
	}
}

func TestEliminatePlayer_AlreadyDead(t *testing.T) {
	game := NewGameState("test")
	player := &Player{ID: "test-1", Name: "Test", Alive: false}
	game.AddPlayer(player)

	result := game.EliminatePlayer("test-1")

	if result != nil {
		t.Error("EliminatePlayer should return nil for already dead player")
	}
}

// --- RegisterVote Tests ---

func TestRegisterVote(t *testing.T) {
	game := NewGameState("test")
	game.AddPlayer(&Player{ID: "voter", Name: "Voter", Alive: true})
	game.AddPlayer(&Player{ID: "target", Name: "Target", Alive: true})

	result := game.RegisterVote("voter", "target")

	if !result {
		t.Error("RegisterVote should return true for valid vote")
	}
	if game.Votes["voter"] != "target" {
		t.Error("vote should be recorded")
	}
}

func TestRegisterVote_DeadVoter(t *testing.T) {
	game := NewGameState("test")
	game.AddPlayer(&Player{ID: "voter", Name: "Voter", Alive: false})
	game.AddPlayer(&Player{ID: "target", Name: "Target", Alive: true})

	result := game.RegisterVote("voter", "target")

	if result {
		t.Error("dead player should not be able to vote")
	}
}

func TestRegisterVote_DeadTarget(t *testing.T) {
	game := NewGameState("test")
	game.AddPlayer(&Player{ID: "voter", Name: "Voter", Alive: true})
	game.AddPlayer(&Player{ID: "target", Name: "Target", Alive: false})

	result := game.RegisterVote("voter", "target")

	if result {
		t.Error("should not be able to vote for dead player")
	}
}

func TestRegisterVote_DuplicateVote(t *testing.T) {
	game := NewGameState("test")
	game.AddPlayer(&Player{ID: "voter", Name: "Voter", Alive: true})
	game.AddPlayer(&Player{ID: "target1", Name: "Target1", Alive: true})
	game.AddPlayer(&Player{ID: "target2", Name: "Target2", Alive: true})

	game.RegisterVote("voter", "target1")
	result := game.RegisterVote("voter", "target2")

	if result {
		t.Error("should not be able to vote twice")
	}
	if game.Votes["voter"] != "target1" {
		t.Error("original vote should be preserved")
	}
}

// --- SetNightAction Tests ---

func TestSetNightAction_Mafia(t *testing.T) {
	game := NewGameState("test")
	game.AddPlayer(&Player{ID: "mafia-1", Name: "Mafia", Role: RoleMafia, Alive: true})
	game.AddPlayer(&Player{ID: "target", Name: "Target", Alive: true})

	result := game.SetNightAction(RoleMafia, "mafia-1", "target")

	if !result {
		t.Error("SetNightAction should return true for valid action")
	}
	if game.MafiaTarget != "target" {
		t.Error("mafia target should be recorded")
	}
}

func TestSetNightAction_Doctor(t *testing.T) {
	game := NewGameState("test")
	game.AddPlayer(&Player{ID: "target", Name: "Target", Alive: true})

	result := game.SetNightAction(RoleDoctor, "target", "target")

	if !result {
		t.Error("SetNightAction should return true")
	}
	if game.DoctorTarget != "target" {
		t.Error("doctor target should be recorded")
	}
}

func TestSetNightAction_Sheriff(t *testing.T) {
	game := NewGameState("test")
	game.AddPlayer(&Player{ID: "sheriff-1", Name: "Sheriff", Role: RoleSheriff, Alive: true})
	game.AddPlayer(&Player{ID: "target", Name: "Target", Alive: true})

	result := game.SetNightAction(RoleSheriff, "sheriff-1", "target")

	if !result {
		t.Error("SetNightAction should return true")
	}
	if game.SheriffTarget != "target" {
		t.Error("sheriff target should be recorded")
	}
}

func TestSetNightAction_VillagerCannot(t *testing.T) {
	game := NewGameState("test")
	game.AddPlayer(&Player{ID: "villager-1", Name: "Villager", Role: RoleVillager, Alive: true})
	game.AddPlayer(&Player{ID: "target", Name: "Target", Alive: true})

	result := game.SetNightAction(RoleVillager, "villager-1", "target")

	if result {
		t.Error("villager should not have night action")
	}
}

func TestSetNightAction_AlreadySet(t *testing.T) {
	game := NewGameState("test")
	game.AddPlayer(&Player{ID: "target1", Name: "Target1", Alive: true})
	game.AddPlayer(&Player{ID: "mafia-1", Name: "Mafia1", Role: RoleMafia, Alive: true})
	game.AddPlayer(&Player{ID: "mafia-2", Name: "Mafia2", Role: RoleMafia, Alive: true})
	game.AddPlayer(&Player{ID: "target2", Name: "Target2", Alive: true})

	game.SetNightAction(RoleMafia, "mafia-1", "target2")
	result := game.SetNightAction(RoleMafia, "mafia-2", "target1")

	if result {
		t.Error("should not be able to change night action")
	}
	if game.MafiaTarget != "target2" {
		t.Error("original target should be preserved")
	}
}

// --- ResetPhaseData Tests ---

func TestResetPhaseData(t *testing.T) {
	game := NewGameState("test")
	game.AddPlayer(&Player{ID: "voter", Name: "Voter", Alive: true})
	game.AddPlayer(&Player{ID: "target", Name: "Target", Alive: true})

	// Set some data
	game.RegisterVote("voter", "target")
	game.SetNightAction(RoleMafia, "voter", "target")
	game.SetNightAction(RoleDoctor, "voter", "target")
	game.SetNightAction(RoleSheriff, "voter", "target")

	// Reset
	game.ResetPhaseData()

	if len(game.Votes) != 0 {
		t.Error("votes should be cleared")
	}
	if game.MafiaTarget != "" {
		t.Error("mafia target should be cleared")
	}
	if game.DoctorTarget != "" {
		t.Error("doctor target should be cleared")
	}
	if game.SheriffTarget != "" {
		t.Error("sheriff target should be cleared")
	}
}

// --- ShufflePlayerOrder Tests ---

func TestShufflePlayerOrder_ReturnsSameCount(t *testing.T) {
	game := createTestGame(6)

	shuffled := game.ShufflePlayerOrder()

	if len(shuffled) != 6 {
		t.Errorf("shuffled count: got %d, expected 6", len(shuffled))
	}
}

func TestShufflePlayerOrder_ContainsSamePlayers(t *testing.T) {
	game := createTestGame(6)
	original := game.GetAlivePlayers()

	shuffled := game.ShufflePlayerOrder()

	// Check all original players are in shuffled
	for _, player := range original {
		found := false
		for _, s := range shuffled {
			if s.ID == player.ID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("player %s missing from shuffled result", player.ID)
		}
	}
}

func TestShufflePlayerOrder_ExcludesDeadPlayers(t *testing.T) {
	ResetPlayerCounter()
	game := NewGameState("test")

	alive, _ := NewPlayer(CreatePlayerID(), "Alive", RoleUnknown)
	dead, _ := NewPlayer(CreatePlayerID(), "Dead", RoleUnknown)
	dead.Alive = false

	game.AddPlayer(alive)
	game.AddPlayer(dead)

	shuffled := game.ShufflePlayerOrder()

	if len(shuffled) != 1 {
		t.Errorf("should only include alive players, got %d", len(shuffled))
	}
}

func TestAssignRolesToPlayers(t *testing.T) {
	ResetPlayerCounter()
	game := createTestGame(6)

	// should create 2 mafia, 2 villagers, 1 doctor, 1 sheriff
	game.AssignRolesToPlayers(GetRoleDistribution(6))
	expectedResult := map[Role]int{
		RoleVillager: 2,
		RoleMafia:    2,
		RoleDoctor:   1,
		RoleSheriff:  1,
	}

	for role, count := range expectedResult {
		// Count players with this role manually
		found := 0
		for _, player := range game.Players {
			if player.Role == role {
				found++
			}
		}
		if found != count {
			t.Errorf("expected %d %s , got %d", count, role, found)
		}
	}
}

// testing scenerios on 4 players
func TestIsGameOver(t *testing.T) {
	tests := []struct {
		name     string
		roles    []Role
		expected bool
	}{
		{
			name: "game over - villager win, 0 mafia",
			roles: []Role{
				RoleVillager, RoleVillager, RoleVillager, RoleVillager,
			},
			expected: true,
		},
		{
			name: "game over - mafia win, numeric advantage",
			roles: []Role{
				RoleVillager, RoleMafia, RoleMafia, RoleMafia,
			},
			expected: true,
		},
		{
			name: "game over - mafia win, equal players",
			roles: []Role{
				RoleVillager, RoleVillager, RoleMafia, RoleMafia,
			},
			expected: true,
		},
		{
			name: "game not over - mafia still remain",
			roles: []Role{
				RoleVillager, RoleVillager, RoleVillager, RoleMafia,
			},
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := NewGameState("test")
			for i, role := range tt.roles {
				id := fmt.Sprintf("p%d", i)
				game.AddPlayer(&Player{ID: id, Name: id, Role: role, Alive: true})
			}

			result := game.IsGameOver()
			if result != tt.expected {
				t.Errorf("got %v, expected %v", result, tt.expected)
			}
		})
	}
}
