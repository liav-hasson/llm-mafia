package engine

import (
	"context"
	"errors"
	"sync"

	"mafia-engine/internal/config"
	"mafia-engine/internal/domain"
	"mafia-engine/internal/events"
	"mafia-engine/internal/kafka"
	"mafia-engine/internal/names"
)

// Effect represents a side effect that should be executed after state mutation.
// Effects are pure descriptions of what should happen - they don't execute themselves.
// The engine loop is responsible for executing effects.
type Effect interface {
	// Execute performs the side effect (Kafka publish, timer, etc.)
	// Context allows cancellation and timeout control
	Execute(ctx context.Context, producer kafka.Producer) error
}

// Command represents a pure state transformation.
// Commands validate inputs, mutate state, and return effects to be executed.
// Commands must NOT perform I/O or call time.Now() - they should be deterministic.
type Command interface {
	// Apply executes the command logic:
	// 1. Validate inputs against current state
	// 2. Mutate state if valid
	// 3. Return effects that should be executed
	// Returns effects to execute and error if validation fails
	Apply(state *domain.GameState) ([]Effect, error)
}

// Engine is the authoritative orchestrator of a single Mafia game.
// It owns game state, receives events, and emits new events.
// All state mutation is serialized through a single internal loop.
type Engine struct {
	// state is the authoritative in-memory game state.
	state *domain.GameState

	// producer emits authoritative events.
	producer kafka.Producer

	// cfg holds runtime configuration (timeouts, limits, etc.)
	cfg *config.Config

	// nameGen generates player names from configured list.
	nameGen *names.Generator

	// cmdCh carries internal commands that mutate state.
	cmdCh chan Command

	// timers manages phase timeout timers.
	timers *TimerManager

	// ctx controls engine lifecycle.
	ctx    context.Context
	cancel context.CancelFunc

	// wg ensures clean shutdown.
	wg sync.WaitGroup
}

// NewEngine constructs an Engine with its dependencies wired.
// It does not start any goroutines.
func NewEngine(
	initialState *domain.GameState,
	producer kafka.Producer,
	cfg *config.Config,
) (*Engine, error) {

	if initialState == nil {
		return nil, errors.New("initial state must not be nil")
	}
	if producer == nil {
		return nil, errors.New("producer must not be nil")
	}
	if cfg == nil {
		return nil, errors.New("config must not be nil")
	}

	// Create name generator from config
	nameGen, err := names.NewGenerator(cfg.PlayerNames)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Engine{
		state:    initialState,
		producer: producer,
		cfg:      cfg,
		nameGen:  nameGen,
		cmdCh:    make(chan Command, 64),
		timers:   NewTimerManager(),
		ctx:      ctx,
		cancel:   cancel,
	}, nil
}

// Start launches the engine loop.
// State mutation is possible only after Start is called.
func (e *Engine) Start() {
	e.wg.Add(1)
	go func() {
		defer e.wg.Done()
		e.run()
	}()
}

// Stop cancels the engine context and waits for shutdown.
func (e *Engine) Stop() {
	e.timers.Shutdown()
	e.cancel()
	e.wg.Wait()
}

// CreatePlayer creates a new player with auto-generated ID and name.
// Returns the created player or an error if name generation fails.
func (e *Engine) CreatePlayer() (*domain.Player, error) {
	id := domain.CreatePlayerID()
	name, err := e.nameGen.Next()
	if err != nil {
		return nil, err
	}
	return domain.NewPlayer(id, name, domain.RoleUnknown)
}

// AddPlayer adds a player to the game during the waiting phase.
// Creates a player with auto-generated ID and name, then sends AddPlayerCommand.
// Returns error if name generation fails or command validation fails.
func (e *Engine) AddPlayer() error {
	player, err := e.CreatePlayer()
	if err != nil {
		return err
	}

	cmd := &AddPlayerCommand{
		Player:     player,
		MaxPlayers: e.cfg.GameMaxPlayers,
	}

	// Send command to the engine loop
	select {
	case e.cmdCh <- cmd:
		return nil
	case <-e.ctx.Done():
		return e.ctx.Err()
	}
}

// StartGame sends a StartGameCommand to the engine.
// It uses min/max players from the configuration.
func (e *Engine) StartGame() error {
	cmd := &StartGameCommand{
		MinPlayers: e.cfg.GameMinPlayers,
		MaxPlayers: e.cfg.GameMaxPlayers,
	}

	select {
	case e.cmdCh <- cmd:
		return nil
	case <-e.ctx.Done():
		return e.ctx.Err()
	}
}

// HandleMessage is the single external entrypoint into the engine.
// It deserializes the event and delegates interpretation to handlers.
func (e *Engine) HandleMessage(ctx context.Context, msg kafka.Message) error {
	ev, err := events.Deserialize(msg.Value)
	if err != nil {
		return err
	}
	return HandleEvent(ctx, e.cmdCh, ev)
}
