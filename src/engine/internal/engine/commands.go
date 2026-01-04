package engine

import (
	"fmt"

	"mafia-engine/internal/domain"
	"mafia-engine/internal/events"
)

// AddPlayerCommand adds a new player to the game before it starts.
// This can only be called during the waiting phase.
type AddPlayerCommand struct {
	PlayerName string // Optional - will auto-generate if empty
	MaxPlayers int    // Maximum allowed players
}

func (c *AddPlayerCommand) Apply(state *domain.GameState) ([]Effect, error) {
	// Validation 1: Game must be in waiting phase
	if state.Phase != domain.PhaseWaiting {
		return nil, fmt.Errorf("cannot add players after game has started")
	}

	// Validation 2: Check if we can add another player
	currentCount := state.GetPlayerCount()
	if !domain.CanAddPlayer(currentCount, c.MaxPlayers) {
		return nil, fmt.Errorf("cannot add player: max players (%d) reached", c.MaxPlayers)
	}

	// Create new player using domain helper
	// NewPlayer will auto-generate ID and name if needed
	player, err := domain.NewPlayer("", c.PlayerName, domain.RoleUnknown)
	if err != nil {
		return nil, fmt.Errorf("failed to create player: %w", err)
	}

	// Add player to state using domain helper
	addedPlayer := state.AddPlayer(player)
	if addedPlayer == nil {
		return nil, fmt.Errorf("failed to add player: duplicate ID")
	}

	// No effects - player addition is silent
	// Clients will see player list when game starts
	return []Effect{}, nil
}

// StartGameCommand initializes the game by assigning roles and emitting GameStarted.
// This should be called after all players have been added.
type StartGameCommand struct {
	MinPlayers int // Minimum required players to start
	MaxPlayers int // Maximum allowed players
}

func (c *StartGameCommand) Apply(state *domain.GameState) ([]Effect, error) {
	// Validation 1: Game must be in waiting phase
	if state.Phase != domain.PhaseWaiting {
		return nil, fmt.Errorf("cannot start game in phase %s", state.Phase)
	}

	// Validation 2: Check player count using domain helper
	currentCount := state.GetPlayerCount()
	if !domain.CanStartGame(currentCount, c.MinPlayers, c.MaxPlayers) {
		return nil, fmt.Errorf("cannot start game: need %d-%d players, have %d",
			c.MinPlayers, c.MaxPlayers, currentCount)
	}

	// Calculate role distribution using domain helper
	roleDistribution := domain.GetRoleDistribution(currentCount)

	// Use domain helper to assign roles
	state.AssignRolesToPlayers(roleDistribution)

	// Transition to Night phase (game starts at night for mafia coordination)
	state.Phase = domain.PhaseNight
	state.Round = 1

	// Build effects
	effects := []Effect{}

	// Emit GameStarted event
	playerIDs := make([]string, 0, len(state.Players))
	for id := range state.Players {
		playerIDs = append(playerIDs, id)
	}

	gameStartedEvent := &events.GameStarted{
		BaseEvent: events.BaseEvent{
			GameID: state.ID,
			Type:   events.TypeGameStarted,
		},
		Players: playerIDs,
	}
	effects = append(effects, NewPublishEffect(gameStartedEvent))

	// Emit RoleAssigned events (one per player)
	for _, player := range state.Players {
		roleEvent := &events.RoleAssigned{
			BaseEvent: events.BaseEvent{
				GameID: state.ID,
				Type:   events.TypeRoleAssigned,
			},
			PlayerID: player.ID,
			Role:     player.Role.String(),
		}
		effects = append(effects, NewPublishEffect(roleEvent))
	}

	// Emit PhaseChanged to indicate game has started in Night phase
	phaseEvent := &events.PhaseChanged{
		BaseEvent: events.BaseEvent{
			GameID: state.ID,
			Type:   events.TypePhaseChanged,
		},
		Round:    state.Round,
		OldPhase: domain.PhaseWaiting.String(),
		NewPhase: domain.PhaseNight.String(),
	}
	effects = append(effects, NewPublishEffect(phaseEvent))

	return effects, nil
}

// VoteCommand records a player's vote during the voting phase.
// This is a pure state mutation - no effects are emitted.
// Votes are tallied silently and resolved at phase change.
type VoteCommand struct {
	VoterID  string
	TargetID string
}

// Apply implements the Command interface.
// It validates the vote and records it in game state.
// Returns empty effects slice - voting is silent.
func (c *VoteCommand) Apply(state *domain.GameState) ([]Effect, error) {
	// Validation: Check voting phase
	if state.Phase != domain.PhaseVoting {
		return nil, fmt.Errorf("cannot vote in phase %s", state.Phase)
	}

	// Use domain helper for validation and mutation
	// RegisterVote handles:
	// - Voter exists and is alive
	// - Target exists and is alive
	// - No duplicate votes
	success := state.RegisterVote(c.VoterID, c.TargetID)
	if !success {
		return nil, fmt.Errorf("vote rejected: invalid voter/target or duplicate vote")
	}

	// No effects - votes are silent until tallied
	return []Effect{}, nil
}

// ChatCommand handles public chat messages.
// It doesn't mutate state but returns a PublishEffect for the engine to execute.
type ChatCommand struct {
	SenderID string
	Message  string
}

// Apply implements the Command interface.
// It validates the sender and returns a PublishEffect.
// Note: Does NOT call time.Now() - engine provides timestamp via effect.
func (c *ChatCommand) Apply(state *domain.GameState) ([]Effect, error) {
	// Validation: Sender exists and is alive
	sender := state.GetPlayer(c.SenderID)
	if sender == nil {
		return nil, fmt.Errorf("sender %s not found", c.SenderID)
	}
	if !sender.Alive {
		return nil, fmt.Errorf("sender %s is dead and cannot speak", c.SenderID)
	}

	// No state mutation - chat is stateless

	// Create the event (without timestamp - engine will inject it)
	event := &events.AllChatMessage{
		BaseEvent: events.BaseEvent{
			GameID: state.ID,
			Type:   events.TypeAllChatMessage,
			// Timestamp will be injected by PublishEffect.Execute
		},
		Message:  c.Message,
		SenderID: c.SenderID,
	}

	// Return effect for engine to execute
	effect := NewPublishEffect(event)
	return []Effect{effect}, nil
}

// MafiaChatCommand handles private mafia chat messages.
// Only mafia members can send these during night phase.
type MafiaChatCommand struct {
	SenderID string
	Message  string
}

func (c *MafiaChatCommand) Apply(state *domain.GameState) ([]Effect, error) {
	// Validation 1: Sender exists and is alive
	sender := state.GetPlayer(c.SenderID)
	if sender == nil {
		return nil, fmt.Errorf("sender %s not found", c.SenderID)
	}
	if !sender.Alive {
		return nil, fmt.Errorf("sender %s is dead and cannot speak", c.SenderID)
	}

	// Validation 2: Sender is mafia
	if !sender.Role.IsMafiaTeam() {
		return nil, fmt.Errorf("sender %s is not mafia and cannot use mafia chat", c.SenderID)
	}

	// Validation 3: Must be night phase (mafia chat only at night)
	if state.Phase != domain.PhaseNight {
		return nil, fmt.Errorf("mafia chat only available during night phase")
	}

	// No state mutation - chat is stateless

	// Create the event
	event := &events.MafiaChatMessage{
		BaseEvent: events.BaseEvent{
			GameID: state.ID,
			Type:   events.TypeMafiaChatMessage,
		},
		Message:  c.Message,
		SenderID: c.SenderID,
	}

	// Return effect for engine to execute
	effect := NewPublishEffect(event)
	return []Effect{effect}, nil
}

// NightActionCommand handles mafia kills, doctor saves, sheriff investigations.
// This is a pure state mutation - no effects until phase resolves.
type NightActionCommand struct {
	Role     string // "mafia", "doctor", "sheriff"
	ActorID  string
	TargetID string
}

func (c *NightActionCommand) Apply(state *domain.GameState) ([]Effect, error) {
	// Validation 1: Check night phase
	if state.Phase != domain.PhaseNight {
		return nil, fmt.Errorf("cannot perform night action in phase %s", state.Phase)
	}

	// Validation 2: Actor exists and is alive
	actor := state.GetPlayer(c.ActorID)
	if actor == nil {
		return nil, fmt.Errorf("actor %s not found", c.ActorID)
	}
	if !actor.Alive {
		return nil, fmt.Errorf("actor %s is dead", c.ActorID)
	}

	// Validation 3: Actor's role matches the action role
	if actor.Role.String() != c.Role {
		return nil, fmt.Errorf("actor %s has role %s but tried to act as %s",
			c.ActorID, actor.Role, c.Role)
	}

	// Rules enforced by SetNightAction:
	// - Target exists and is alive
	// - actor has night action
	// - No duplicate actions this round
	// - Doctor can't save same person twice in a row
	// - Sheriff only has one bullet
	// - Mafia/Sheriff can't self-target (Doctor can)
	success := state.SetNightAction(actor.Role, c.ActorID, c.TargetID)
	if !success {
		return nil, fmt.Errorf("night action rejected: rules violated (check target validity, consecutive saves, or bullet usage)")
	}

	// No effects - night actions are secret until phase resolves
	return []Effect{}, nil
}

// PhaseChangeCommand transitions the game to a new phase.
// This is complex - it mutates state AND returns multiple effects.
type PhaseChangeCommand struct {
	NewPhase domain.Phase
}

func (c *PhaseChangeCommand) Apply(state *domain.GameState) ([]Effect, error) {
	// Track eliminated player for event emission
	var eliminatedPlayerID string
	var eliminationReason string

	// Step 1: Resolve actions from PREVIOUS phase
	switch state.Phase {
	case domain.PhaseNight:
		// Resolve night actions using domain helper
		eliminatedPlayerID = state.ResolveNightActions()
		if eliminatedPlayerID != "" {
			eliminationReason = "killed_by_mafia"
			// Use domain helper to mark player as dead
			state.EliminatePlayer(eliminatedPlayerID)
		}

	case domain.PhaseVoting:
		// Resolve voting using domain helper
		eliminatedPlayerID = state.ResolveVotingPhase()
		if eliminatedPlayerID != "" {
			eliminationReason = "voted_out"
			// Use domain helper to mark player as dead
			state.EliminatePlayer(eliminatedPlayerID)
		}
	}

	// Step 2: Clear phase data (votes and night actions)
	state.ResetPhaseData()

	// Step 3: Update phase
	oldPhase := state.Phase
	state.Phase = c.NewPhase

	// Increment round when entering night phase
	if c.NewPhase == domain.PhaseNight {
		state.Round++
	}

	// Step 4: Check win conditions
	gameEnded := state.IsGameOver()

	// Step 5: Build effects
	effects := []Effect{}

	// Always emit PhaseChanged event
	phaseEvent := &events.PhaseChanged{
		BaseEvent: events.BaseEvent{
			GameID: state.ID,
			Type:   events.TypePhaseChanged,
		},
		Round:    state.Round,
		OldPhase: oldPhase.String(),
		NewPhase: c.NewPhase.String(),
	}
	effects = append(effects, NewPublishEffect(phaseEvent))

	// If someone was eliminated, emit event
	if eliminatedPlayerID != "" {
		eliminatedEvent := &events.PlayerEliminated{
			BaseEvent: events.BaseEvent{
				GameID: state.ID,
				Type:   events.TypePlayerEliminated,
			},
			PlayerID: eliminatedPlayerID,
			Reason:   eliminationReason,
		}
		effects = append(effects, NewPublishEffect(eliminatedEvent))
	}

	// If game ended, emit GameEnded event
	if gameEnded {
		gameEndedEvent := &events.GameEnded{
			BaseEvent: events.BaseEvent{
				GameID: state.ID,
				Type:   events.TypeGameEnded,
			},
			Winner: state.Winner.String(),
		}
		effects = append(effects, NewPublishEffect(gameEndedEvent))
	}

	return effects, nil
}

// EliminatePlayerCommand removes a player from the game.
// This mutates state AND returns effects (elimination + maybe game end).
type EliminatePlayerCommand struct {
	PlayerID string
	Reason   string // "voted_out", "killed_by_mafia", etc.
}

func (c *EliminatePlayerCommand) Apply(state *domain.GameState) ([]Effect, error) {
	// Use domain helper to eliminate player
	// EliminatePlayer validates existence and alive status
	player := state.EliminatePlayer(c.PlayerID)
	if player == nil {
		return nil, fmt.Errorf("player %s not found or already dead", c.PlayerID)
	}

	// Check win conditions using domain helper
	gameEnded := state.IsGameOver()

	// Build effects
	effects := []Effect{}

	// Always emit PlayerEliminated event
	eliminatedEvent := &events.PlayerEliminated{
		BaseEvent: events.BaseEvent{
			GameID: state.ID,
			Type:   events.TypePlayerEliminated,
		},
		PlayerID: c.PlayerID,
		Reason:   c.Reason,
	}
	effects = append(effects, NewPublishEffect(eliminatedEvent))

	// If game ended, emit GameEnded event
	if gameEnded {
		gameEndedEvent := &events.GameEnded{
			BaseEvent: events.BaseEvent{
				GameID: state.ID,
				Type:   events.TypeGameEnded,
			},
			Winner: state.Winner.String(),
		}
		effects = append(effects, NewPublishEffect(gameEndedEvent))
	}

	return effects, nil
}
