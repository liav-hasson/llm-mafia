// constants and calculations to start a game

package domain

// set game constants
// available names in 'player.go'
const (
	MinPlayers = 6
	MaxPlayers = 12
)

func CanAddPlayer(currentPlayerCount int) bool {
	// true if current is smaller then max
	return currentPlayerCount < MaxPlayers
}

func CanStartGame(currentPlayerCount int) bool {
	// true if current is smaller or equal to max AND bigger or equal to min
	return currentPlayerCount <= MaxPlayers && currentPlayerCount >= MinPlayers
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
