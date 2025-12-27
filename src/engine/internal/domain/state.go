package domain

import "github.com/xyproto/randomstring"

// live game status data
type GameState struct {
	// unique game ID
	ID string

	// current status
	Round  int
	Phase  Phase
	Winner Winner

	// game players maps player id -> player object
	// '*Player' used to not copy the entire struct
	Players map[string]*Player

	// day votes, maps voterID -> targetID
	Votes map[string]string

	// night actions (player ID target)
	// will initialize as empty string ""
	MafiaTarget   string
	DoctorTarget  string
	SheriffTarget string
}

// set winner type
type Winner int

const (
	WinnerNone Winner = iota
	WinnerMafia
	WinnerVillage
)

func (w Winner) String() string {
	switch w {
	case WinnerNone:
		return "none"
	case WinnerMafia:
		return "mafia"
	case WinnerVillage:
		return "village"
	default:
		return "invalid"
	}
}

// initialize new game state
func NewGameState() *GameState {
	return &GameState{
		ID:     CreateGameID(),
		Round:  1,
		Phase:  PhaseWaiting,
		Winner: WinnerNone,
		// must use make() to init before use, otherwise nil
		Players: make(map[string]*Player),
		Votes:   make(map[string]string),
		// no need to init MafiaTarget, DoctorTarget SheriffTarget
	}
}

// create random game ID
func CreateGameID() string {
	const idlength = 5
	return randomstring.String(idlength)
}

// GetPlayer retrieves a player by ID from the game
// Returns nil if player doesn't exist
func (g *GameState) GetPlayer(id string) *Player {
	return g.Players[id]
}

// GetAlivePlayers returns a slice of all players who are still alive
func (g *GameState) GetAlivePlayers() []*Player {
	// create empty slice to collect alive players
	// using var instead of make() — starts as nil, append works on nil slices
	var alive []*Player

	// range over map: gives key (id) and value (player)
	// underscore (_) ignores the key since we don't need it
	// potentially sort the output for prnting, Go map lookup is random
	for _, player := range g.Players {
		if player.Alive {
			alive = append(alive, player)
		}
	}

	return alive
}

// IsGameOver checks if win conditions are met and updates the Winner field
// Returns true if game has ended
// Win conditions:
//   - Villagers win: All Mafia players are eliminated
//   - Mafia wins: Mafia count >= Village team count (among alive players)
func (g *GameState) IsGameOver() bool {
	// count alive players by team
	var mafiaAlive, villageAlive int

	for _, player := range g.Players {
		if !player.Alive {
			continue // skip dead players
		}

		if player.Role.IsMafiaTeam() {
			mafiaAlive++
		} else {
			villageAlive++ // includes Villager, Doctor, Sheriff
		}
	}

	// check win conditions
	if mafiaAlive == 0 {
		g.Winner = WinnerVillage
		g.Phase = PhaseEnded
		return true
	}

	if mafiaAlive >= villageAlive {
		g.Winner = WinnerMafia
		g.Phase = PhaseEnded
		return true
	}

	// game continues
	return false
}

// mutate game state

// AddPlayer adds a player to the game
// Returns the added player, or nil if player with same ID already exists
func (g *GameState) AddPlayer(player *Player) *Player {
	// check if player already exists (prevent duplicates)
	if _, exists := g.Players[player.ID]; exists {
		return nil
	}

	// add player to map: key = ID, value = pointer to player
	g.Players[player.ID] = player
	return player
}

// EliminatePlayer marks a player as dead
// Returns the eliminated player, or nil if player not found or already dead
func (g *GameState) EliminatePlayer(id string) *Player {
	player := g.Players[id]

	// check if player exists
	if player == nil {
		return nil
	}

	// check if already dead (can't eliminate twice)
	if !player.Alive {
		return nil
	}

	// mark as dead
	player.Alive = false
	return player
}

// ResetPhaseData clears all votes and night actions
// Called between phases to start fresh
func (g *GameState) ResetPhaseData() {
	// clear day votes — create new empty map
	g.Votes = make(map[string]string)

	// clear night actions — reset to zero value (empty string)
	g.MafiaTarget = ""
	g.DoctorTarget = ""
	g.SheriffTarget = ""
}

// GetPlayerCount returns the total number of players in the game
func (g *GameState) GetPlayerCount() int {
	return len(g.Players)
}

// GetPlayersByRole returns all players with the specified role
// Includes both alive and dead players
func (g *GameState) GetPlayersByRole(role Role) []*Player {
	var players []*Player

	for _, player := range g.Players {
		if player.Role == role {
			players = append(players, player)
		}
	}

	return players
}

// RegisterVote records a day vote from voter to target
// Returns false if:
//   - voter doesn't exist or is dead
//   - target doesn't exist or is dead
//   - voter has already voted (no changing votes)
func (g *GameState) RegisterVote(voterID, targetID string) bool {
	// validate voter exists and is alive
	voter := g.Players[voterID]
	if voter == nil || !voter.Alive {
		return false
	}

	// validate target exists and is alive
	target := g.Players[targetID]
	if target == nil || !target.Alive {
		return false
	}

	// check if voter already voted (reject duplicate votes)
	if _, alreadyVoted := g.Votes[voterID]; alreadyVoted {
		return false
	}

	// record the vote
	g.Votes[voterID] = targetID
	return true
}

// SetNightAction records a night action for a role
// Returns false if:
//   - role doesn't have a night action
//   - action already set for this role (no changing actions)
//   - target doesn't exist or is dead
func (g *GameState) SetNightAction(role Role, targetID string) bool {
	// validate role has night action
	if !role.HasNightAction() {
		return false
	}

	// validate target exists and is alive
	target := g.Players[targetID]
	if target == nil || !target.Alive {
		return false
	}

	// check if action already set and record it
	// each role has its own target field
	switch role {
	case RoleMafia:
		if g.MafiaTarget != "" {
			return false // already set
		}
		g.MafiaTarget = targetID

	case RoleDoctor:
		if g.DoctorTarget != "" {
			return false // already set
		}
		g.DoctorTarget = targetID

	case RoleSheriff:
		if g.SheriffTarget != "" {
			return false // already set
		}
		g.SheriffTarget = targetID

	default:
		return false // unknown role with night action
	}

	return true
}
