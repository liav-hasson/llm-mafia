package events

// base data for all events, embedded in all other structs
// struct tags allow to use Json snake_case format instead of Go PastalCase
type BaseEvent struct {
	GameID    string `json:"game_id"`
	Timestamp int64  `json:"timestamp"`
	Type      string `json:"type"`
}

// engine -> players events
type GameStarted struct {
	BaseEvent
	Players []string `json:"players"`
}

type PhaseChanged struct {
	BaseEvent
	Round    int    `json:"round"`
	NewPhase string `json:"new_phase"`
}

type PlayerEliminated struct {
	BaseEvent
	PlayerID string `json:"player_id"`
	Reason   string `json:"reason"`
}

type GameEnded struct {
	BaseEvent
	WinnerID string `json:"winner"`
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
