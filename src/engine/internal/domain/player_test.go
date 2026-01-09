package domain

import (
	"testing"
)

// --- ID Generation Tests --- //

func TestCreatePlayerID(t *testing.T) {
	// reset counter before test to ensure clean state
	ResetPlayerCounter()

	tests := []struct {
		name     string
		expected string
	}{
		{name: "first ID", expected: "player-1"},
		{name: "second ID", expected: "player-2"},
		{name: "third ID", expected: "player-3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CreatePlayerID()
			if result != tt.expected {
				t.Errorf("got %s, expected %s", result, tt.expected)
			}
		})
	}
}

func TestResetPlayerCounter(t *testing.T) {
	// generate some IDs
	CreatePlayerID()
	CreatePlayerID()

	// reset
	ResetPlayerCounter()

	// should start fresh
	id := CreatePlayerID()
	if id != "player-1" {
		t.Errorf("after reset, got %s, expected player-1", id)
	}
}

// --- NewPlayer Tests --- //

func TestNewPlayer(t *testing.T) {
	tests := []struct {
		name         string
		id           string
		playerName   string
		role         Role
		expectedID   string
		expectedName string
		expectedRole Role
	}{
		{
			name:         "creates player with all fields",
			id:           "test-id",
			playerName:   "Test Name",
			role:         RoleVillager,
			expectedID:   "test-id",
			expectedName: "Test Name",
			expectedRole: RoleVillager,
		},
		{
			name:         "creates mafia player",
			id:           "mafia-1",
			playerName:   "Evil Bob",
			role:         RoleMafia,
			expectedID:   "mafia-1",
			expectedName: "Evil Bob",
			expectedRole: RoleMafia,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player, err := NewPlayer(tt.id, tt.playerName, tt.role)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if player.ID != tt.expectedID {
				t.Errorf("ID: got %s, expected %s", player.ID, tt.expectedID)
			}
			if player.Name != tt.expectedName {
				t.Errorf("Name: got %s, expected %s", player.Name, tt.expectedName)
			}
			if player.Role != tt.expectedRole {
				t.Errorf("Role: got %v, expected %v", player.Role, tt.expectedRole)
			}
			if !player.Alive {
				t.Error("new player should be alive")
			}
		})
	}
}

func TestNewPlayer_RequiresID(t *testing.T) {
	_, err := NewPlayer("", "Test Name", RoleVillager)
	if err == nil {
		t.Error("expected error when ID is empty")
	}
}

func TestNewPlayer_RequiresName(t *testing.T) {
	_, err := NewPlayer("test-id", "", RoleVillager)
	if err == nil {
		t.Error("expected error when name is empty")
	}
}

// --- Role Tests --- //

func TestRoleString(t *testing.T) {
	tests := []struct {
		name     string
		role     Role
		expected string
	}{
		{name: "unknown", role: RoleUnknown, expected: "unknown"},
		{name: "villager", role: RoleVillager, expected: "villager"},
		{name: "mafia", role: RoleMafia, expected: "mafia"},
		{name: "doctor", role: RoleDoctor, expected: "doctor"},
		{name: "sheriff", role: RoleSheriff, expected: "sheriff"},
		{name: "invalid", role: Role(999), expected: "invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.role.String()
			if result != tt.expected {
				t.Errorf("got %s, expected %s", result, tt.expected)
			}
		})
	}
}

func TestRoleIsVillagerTeam(t *testing.T) {
	tests := []struct {
		name     string
		role     Role
		expected bool
	}{
		{name: "villager", role: RoleVillager, expected: true},
		{name: "doctor", role: RoleDoctor, expected: true},
		{name: "sheriff", role: RoleSheriff, expected: true},
		{name: "mafia", role: RoleMafia, expected: false},
		{name: "unknown", role: RoleUnknown, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.role.IsVillagerTeam()
			if result != tt.expected {
				t.Errorf("got %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestRoleIsMafiaTeam(t *testing.T) {
	tests := []struct {
		name     string
		role     Role
		expected bool
	}{
		{name: "mafia", role: RoleMafia, expected: true},
		{name: "villager", role: RoleVillager, expected: false},
		{name: "doctor", role: RoleDoctor, expected: false},
		{name: "sheriff", role: RoleSheriff, expected: false},
		{name: "unknown", role: RoleUnknown, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.role.IsMafiaTeam()
			if result != tt.expected {
				t.Errorf("got %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestRoleHasNightAction(t *testing.T) {
	tests := []struct {
		name     string
		role     Role
		expected bool
	}{
		{name: "mafia", role: RoleMafia, expected: true},
		{name: "doctor", role: RoleDoctor, expected: true},
		{name: "sheriff", role: RoleSheriff, expected: true},
		{name: "villager", role: RoleVillager, expected: false},
		{name: "unknown", role: RoleUnknown, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.role.HasNightAction()
			if result != tt.expected {
				t.Errorf("got %v, expected %v", result, tt.expected)
			}
		})
	}
}
