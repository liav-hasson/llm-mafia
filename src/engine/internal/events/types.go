package events

// Event type constants - stable contract strings used for serialization routing.
// These must match what the Python players and other services emit.
const (
	TypeGameStarted      = "game_started"
	TypePhaseChanged     = "phase_changed"
	TypePlayerEliminated = "player_eliminated"
	TypeGameEnded        = "game_ended"
	TypeAllChatMessage   = "all_chat"
	TypeMafiaChatMessage = "mafia_chat"
	TypePlayerThoughts   = "player_thoughts"
	TypeVoteSubmitted    = "vote_submitted"
	TypeNightAction      = "night_action"
	TypeRoleAssigned     = "role_assigned"
)

// base data for all events, embedded in all other structs
// struct tags allow to use Json snake_case format instead of Go PastalCase
// BaseEvent is the common header for all events.
// Timestamp is Unix time in milliseconds (int64).
// Type is a stable event type string (not runtime-configurable).
type BaseEvent struct {
	GameID    string `json:"game_id"`
	Timestamp int64  `json:"timestamp"` // Unix ms
	Type      string `json:"type"`      // stable contract string
}

// engine -> players events
type GameStarted struct {
	BaseEvent
	Players []string `json:"players"`
}

type PhaseChanged struct {
	BaseEvent
	Round    int    `json:"round"`
	OldPhase string `json:"old_phase"`
	NewPhase string `json:"new_phase"`
}

type PlayerEliminated struct {
	BaseEvent
	PlayerID string `json:"player_id"`
	Reason   string `json:"reason"`
}

type GameEnded struct {
	BaseEvent
	Winner string `json:"winner"`
}

// players -> players + engine events
type AllChatMessage struct {
	BaseEvent
	Message  string `json:"message"`
	SenderID string `json:"sender"`
}

type MafiaChatMessage struct {
	BaseEvent
	Message  string `json:"message"`
	SenderID string `json:"sender"`
}

// players -> engine events
type PlayerThoughts struct {
	BaseEvent
	Thought  string `json:"thought"`
	SenderID string `json:"sender"`
}

type VoteSubmitted struct {
	BaseEvent
	VoterID  string `json:"voter"`
	TargetID string `json:"target"`
}

type NightAction struct {
	BaseEvent
	// mafia, sheriff, doctor
	Role     string `json:"role"`
	ActorID  string `json:"actor"`
	TargetID string `json:"target"`
}

// Private - sent per-player
type RoleAssigned struct {
	BaseEvent
	PlayerID string `json:"player_id"`
	Role     string `json:"role"`
}
