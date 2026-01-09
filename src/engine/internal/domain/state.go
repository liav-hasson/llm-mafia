// This file containes game state structs and supporting methods

package domain

import (
	"fmt"
	"math/rand"

	"github.com/xyproto/randomstring"
)

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

	// night action history for rule enforcement
	PreviousDoctorTarget string // Track last save (can't save same person twice in a row)
	SheriffUsedBullet    bool   // Sheriff only has one bullet
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

// --- reading game state --- //

// GetPlayer retrieves a player by ID from the game
// Returns nil if player doesn't exist
func (g *GameState) GetPlayer(id string) *Player {
	return g.Players[id]
}

// GetAlivePlayers returns a slice of all players who are still alive
// Used internally by ShufflePlayerOrder()
func (g *GameState) GetAlivePlayers() []*Player {
	// create empty slice to collect alive players
	// using var instead of make() — starts as nil, append works on nil slices
	var alive []*Player

	// range over map: gives key (id) and value (player)
	// underscore (_) ignores the key since we don't need it
	for _, player := range g.Players {
		if player.Alive {
			alive = append(alive, player)
		}
	}

	return alive
}

// GetPlayerCount returns the total number of players in the game
func (g *GameState) GetPlayerCount() int {
	return len(g.Players)
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

// --- mutating game state --- //

// NewGameState initializes a new game state with the given ID prefix.
func NewGameState(idPrefix string) *GameState {
	return &GameState{
		ID:     CreateGameID(idPrefix),
		Round:  1,
		Phase:  PhaseWaiting,
		Winner: WinnerNone,
		// must use make() to init before use, otherwise nil
		Players: make(map[string]*Player),
		Votes:   make(map[string]string),
		// no need to init MafiaTarget, DoctorTarget SheriffTarget
	}
}

// CreateGameID creates a random game ID with the given prefix.
// Format: {prefix}-{random-string}
// Example: "game-a3k9m" or "dev-x7p2q"
func CreateGameID(prefix string) string {
	const idlength = 5
	randomSuffix := randomstring.String(idlength)
	return fmt.Sprintf("%s-%s", prefix, randomSuffix)
}

func (g *GameState) ShufflePlayerOrder() []*Player {
	players := g.GetAlivePlayers()

	rand.Shuffle(len(players), func(i, j int) {
		players[i], players[j] = players[j], players[i]
	})

	return players
}

// assign roles to players, takes map[Role]int from 'GetRoleDistribution()'
func (g *GameState) AssignRolesToPlayers(roleDistribution map[Role]int) {
	shuffledPlayers := g.ShufflePlayerOrder()
	playerIndex := 0

	for role, count := range roleDistribution {
		for range count {
			shuffledPlayers[playerIndex].Role = role
			playerIndex++
		}
	}
}

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
// Preserves PreviousDoctorTarget for rule enforcement
func (g *GameState) ResetPhaseData() {
	// clear day votes — create new empty map
	g.Votes = make(map[string]string)

	// save doctor target before clearing (for consecutive save rule)
	g.PreviousDoctorTarget = g.DoctorTarget

	// clear night actions — reset to zero value (empty string)
	g.MafiaTarget = ""
	g.DoctorTarget = ""
	g.SheriffTarget = ""
	// Note: SheriffUsedBullet persists across rounds (one bullet per game)
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
//   - doctor tries to save same person as last round
//   - sheriff already used their bullet
//   - mafia/sheriff tries to target themselves (doctor CAN self-save)
func (g *GameState) SetNightAction(role Role, actorID, targetID string) bool {
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
		// Mafia cannot target themselves
		if actorID == targetID {
			return false
		}
		g.MafiaTarget = targetID

	case RoleDoctor:
		if g.DoctorTarget != "" {
			return false // already set
		}
		// Doctor cannot save the same person two rounds in a row
		if g.PreviousDoctorTarget == targetID {
			return false
		}
		// Doctor CAN save themselves
		g.DoctorTarget = targetID

	case RoleSheriff:
		if g.SheriffTarget != "" {
			return false // already set
		}
		// Sheriff only has one bullet
		if g.SheriffUsedBullet {
			return false
		}
		// Sheriff cannot investigate themselves
		if actorID == targetID {
			return false
		}
		g.SheriffTarget = targetID
		g.SheriffUsedBullet = true // Mark bullet as used

	default:
		return false // unknown role with night action
	}

	return true
}

// ResolveNightActions processes night actions and returns eliminated player ID
// Returns empty string if no one was eliminated
// Logic:
//   - Mafia kills their target
//   - Doctor saves their target
//   - If saved target == killed target, no elimination
func (g *GameState) ResolveNightActions() string {
	// If no mafia target, no one dies
	if g.MafiaTarget == "" {
		return ""
	}

	// If doctor saved the mafia target, no one dies
	if g.DoctorTarget == g.MafiaTarget {
		return ""
	}

	// Mafia target was not saved — they die
	return g.MafiaTarget
}

// ResolveVotingPhase tallies votes and returns eliminated player ID
// Returns empty string if no one was eliminated (tie or no votes)
func (g *GameState) ResolveVotingPhase() string {
	// If no votes, no one is eliminated
	if len(g.Votes) == 0 {
		return ""
	}

	// Use existing voting logic
	winner, hasWinner := GetVoteWinner(g.Votes)
	if !hasWinner {
		return "" // tie — no elimination
	}

	return winner
}
