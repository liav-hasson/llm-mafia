package domain

import (
	"testing"
)

// --- ID Generation Tests --- //

func TestCreatePlayerID(t *testing.T) {
	// reset counters before test to ensure clean state
	ResetPlayerCounters()

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

// --- Name Generation Tests --- //

func TestCreatePlayerName(t *testing.T) {
	// reset counters before test to ensure clean state
	ResetPlayerCounters()

	tests := []struct {
		name         string
		expectedName string
		expectedErr  error
	}{
		{name: "first name", expectedName: "Gilbert McDonald", expectedErr: nil},
		{name: "second name", expectedName: "Dorothy Bird", expectedErr: nil},
		{name: "third name", expectedName: "Ernest Preston", expectedErr: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, err := CreatePlayerName()

			if gotName != tt.expectedName {
				t.Errorf("got name %q, expected %q", gotName, tt.expectedName)
			}

			if err != tt.expectedErr {
				t.Errorf("got error %v, expected %v", err, tt.expectedErr)
			}
		})
	}
}

func TestCreatePlayerName_ExhaustsNames(t *testing.T) {
	// reset and use all names
	ResetPlayerCounters()

	// use all 12 available names
	for i := 0; i < 12; i++ {
		_, err := CreatePlayerName()
		if err != nil {
			t.Fatalf("unexpected error on name %d: %v", i+1, err)
		}
	}

	// 13th call should return error
	name, err := CreatePlayerName()
	if err != ErrNoMoreNames {
		t.Errorf("expected ErrNoMoreNames, got %v", err)
	}
	if name != "" {
		t.Errorf("expected empty name, got %q", name)
	}
}

func TestResetPlayerCounters(t *testing.T) {
	// generate some IDs and names
	CreatePlayerID()
	CreatePlayerID()
	_, _ = CreatePlayerName() // ignore error for test setup

	// reset
	ResetPlayerCounters()

	// should start fresh
	id := CreatePlayerID()
	if id != "player-1" {
		t.Errorf("after reset, got %s, expected player-1", id)
	}

	name, _ := CreatePlayerName()
	if name != "Gilbert McDonald" {
		t.Errorf("after reset, got %s, expected Gilbert McDonald", name)
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

func TestNewPlayer_AutoGeneratesID(t *testing.T) {
	ResetPlayerCounters()

	// empty ID should auto-generate
	player, err := NewPlayer("", "Test Name", RoleVillager)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if player.ID != "player-1" {
		t.Errorf("got ID %s, expected player-1", player.ID)
	}
}

func TestNewPlayer_AutoGeneratesName(t *testing.T) {
	ResetPlayerCounters()

	// empty name should auto-generate
	player, err := NewPlayer("test-id", "", RoleVillager)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if player.Name != "Gilbert McDonald" {
		t.Errorf("got name %s, expected Gilbert McDonald", player.Name)
	}
}

// --- Role Tests --- //

func TestRoleString(t *testing.T) {
	tests := []struct {
		role     Role
		expected string
	}{
		{RoleUnknown, "unknown"},
		{RoleVillager, "villager"},
		{RoleMafia, "mafia"},
		{RoleDoctor, "doctor"},
		{RoleSheriff, "sheriff"},
		{Role(99), "invalid"}, // unknown role value
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.role.String()
			if result != tt.expected {
				t.Errorf("got %s, expected %s", result, tt.expected)
			}
		})
	}
}

func TestRoleIsVillagerTeam(t *testing.T) {
	tests := []struct {
		role     Role
		expected bool
	}{
		{RoleVillager, true},
		{RoleDoctor, true},
		{RoleSheriff, true},
		{RoleMafia, false},
		{RoleUnknown, false},
	}

	for _, tt := range tests {
		t.Run(tt.role.String(), func(t *testing.T) {
			result := tt.role.IsVillagerTeam()
			if result != tt.expected {
				t.Errorf("%s.IsVillagerTeam(): got %v, expected %v",
					tt.role, result, tt.expected)
			}
		})
	}
}

func TestRoleIsMafiaTeam(t *testing.T) {
	tests := []struct {
		role     Role
		expected bool
	}{
		{RoleMafia, true},
		{RoleVillager, false},
		{RoleDoctor, false},
		{RoleSheriff, false},
		{RoleUnknown, false},
	}

	for _, tt := range tests {
		t.Run(tt.role.String(), func(t *testing.T) {
			result := tt.role.IsMafiaTeam()
			if result != tt.expected {
				t.Errorf("%s.IsMafiaTeam(): got %v, expected %v",
					tt.role, result, tt.expected)
			}
		})
	}
}

func TestRoleHasNightAction(t *testing.T) {
	tests := []struct {
		role     Role
		expected bool
	}{
		{RoleMafia, true},
		{RoleDoctor, true},
		{RoleSheriff, true},
		{RoleVillager, false},
		{RoleUnknown, false},
	}

	for _, tt := range tests {
		t.Run(tt.role.String(), func(t *testing.T) {
			result := tt.role.HasNightAction()
			if result != tt.expected {
				t.Errorf("%s.HasNightAction(): got %v, expected %v",
					tt.role, result, tt.expected)
			}
		})
	}
}

// --- Player Helper Tests --- //
