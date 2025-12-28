package domain

import "testing"

// TestCanAddPlayer tests the player limit check
func TestCanAddPlayer(t *testing.T) {
	tests := []struct {
		name        string
		playerCount int
		expected    bool
	}{
		{"zero players can add", 0, true},
		{"one player can add", 1, true},
		{"at min-1 can add", MinPlayers - 1, true},
		{"at min can add", MinPlayers, true},
		{"at max-1 can add", MaxPlayers - 1, true},
		{"at max cannot add", MaxPlayers, false},
		{"over max cannot add", MaxPlayers + 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CanAddPlayer(tt.playerCount)
			if result != tt.expected {
				t.Errorf("CanAddPlayer(%d): got %v, expected %v",
					tt.playerCount, result, tt.expected)
			}
		})
	}
}

// TestCanStartGame tests the game start requirements
func TestCanStartGame(t *testing.T) {
	tests := []struct {
		name        string
		playerCount int
		expected    bool
	}{
		{"zero players cannot start", 0, false},
		{"one player cannot start", 1, false},
		{"min-1 cannot start", MinPlayers - 1, false},
		{"at min can start", MinPlayers, true},
		{"between min and max can start", (MinPlayers + MaxPlayers) / 2, true},
		{"at max can start", MaxPlayers, true},
		{"over max cannot start", MaxPlayers + 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CanStartGame(tt.playerCount)
			if result != tt.expected {
				t.Errorf("CanStartGame(%d): got %v, expected %v",
					tt.playerCount, result, tt.expected)
			}
		})
	}
}

// TestMinMaxPlayersConstants verifies the game constants are sensible
// using constants
func TestMinMaxPlayersConstants(t *testing.T) {
	if MinPlayers < 6 {
		t.Errorf("MinPlayers should be at least 4 for a fun game, got %d", MinPlayers)
	}
	if MaxPlayers < MinPlayers {
		t.Errorf("MaxPlayers (%d) should be >= MinPlayers (%d)", MaxPlayers, MinPlayers)
	}
	if MaxPlayers > 12 {
		t.Errorf("MaxPlayers (%d) seems too high, might cause issues", MaxPlayers)
	}
}

// TODO: Write this test
//
// Things to verify:
// 1. Mafia count is playerCount / 3
// 2. Doctor count is 1
// 3. Sheriff count is 1
// 4. Villager count fills the rest
// 5. All counts add up to total playerCount
//
// Test cases to consider: 6 players, 9 players, 12 players

// TestGetRoleDistribution verifies role distribution follows the rules:
// - Mafia: playerCount / 3
// - Doctor: 1
// - Sheriff: 1
// - Villager: fills the rest
func TestGetRoleDistribution(t *testing.T) {
	tests := []struct {
		name        string
		playerCount int
		expected    map[Role]int
	}{
		{
			name:        "6 players",
			playerCount: 6,
			expected: map[Role]int{
				RoleMafia:    2, // 6/3 = 2
				RoleDoctor:   1,
				RoleSheriff:  1,
				RoleVillager: 2, // 6 - 2 - 1 - 1 = 2
			},
		},
		{
			name:        "9 players",
			playerCount: 9,
			expected: map[Role]int{
				RoleMafia:    3, // 9/3 = 3
				RoleDoctor:   1,
				RoleSheriff:  1,
				RoleVillager: 4, // 9 - 3 - 1 - 1 = 4
			},
		},
		{
			name:        "12 players (max)",
			playerCount: 12,
			expected: map[Role]int{
				RoleMafia:    4, // 12/3 = 4
				RoleDoctor:   1,
				RoleSheriff:  1,
				RoleVillager: 6, // 12 - 4 - 1 - 1 = 6
			},
		},
		{
			name:        "7 players (non-divisible by 3)",
			playerCount: 7,
			expected: map[Role]int{
				RoleMafia:    2, // 7/3 = 2 (integer division)
				RoleDoctor:   1,
				RoleSheriff:  1,
				RoleVillager: 3, // 7 - 2 - 1 - 1 = 3
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetRoleDistribution(tt.playerCount)

			// Check each role count matches expected
			for role, expectedCount := range tt.expected {
				if result[role] != expectedCount {
					t.Errorf("%s count: got %d, expected %d",
						role, result[role], expectedCount)
				}
			}

			// Verify total adds up to playerCount
			total := 0
			for _, count := range result {
				total += count
			}
			if total != tt.playerCount {
				t.Errorf("total roles: got %d, expected %d", total, tt.playerCount)
			}
		})
	}
}
