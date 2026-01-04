// constants and calculations to start a game

package domain

// Rule helpers are pure functions. Minimum/maximum player limits are provided
// by the caller (engine) so they can be configured at runtime.

// CanAddPlayer returns true if we can add a player given the configured max.
func CanAddPlayer(currentPlayerCount, maxPlayers int) bool {
	return currentPlayerCount < maxPlayers
}

// CanStartGame returns true if currentPlayerCount is between min and max inclusive.
func CanStartGame(currentPlayerCount, minPlayers, maxPlayers int) bool {
	return currentPlayerCount <= maxPlayers && currentPlayerCount >= minPlayers
}

// mafia - villager ratio is 1/3, and always 1 doctor and 1 sheriff
func GetRoleDistribution(currentPlayerCount int) map[Role]int {
	mafiaCount := currentPlayerCount / 3
	doctorCount := 1
	sheriffCount := 1

	// rest players are villagers
	villagerCount := currentPlayerCount - mafiaCount - doctorCount - sheriffCount

	return map[Role]int{
		RoleVillager: villagerCount,
		RoleMafia:    mafiaCount,
		RoleDoctor:   doctorCount,
		RoleSheriff:  sheriffCount,
	}
}
