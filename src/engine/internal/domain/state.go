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

// read state functions
func (g *GameState) GetPlayer(id string) *Player
func (g *GameState) GetAlivePlayers() []*Player
func (g *GameState) IsGameOver(id string) bool

// mutate game state
func (g *GameState) AddPlayer(player *Player) *Player
func (g *GameState) EliminatePlayer(id string) *Player
func (g *GameState) ResetPhaseData()
