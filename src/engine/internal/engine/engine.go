package engine

import (
	"context"
	"errors"
	"sync"

	"mafia-engine/internal/domain"
	"mafia-engine/internal/events"
	"mafia-engine/internal/kafka"
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
) (*Engine, error) {

	if initialState == nil {
		return nil, errors.New("initial state must not be nil")
	}
	if producer == nil {
		return nil, errors.New("producer must not be nil")
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Engine{
		state:    initialState,
		producer: producer,
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

// HandleMessage is the single external entrypoint into the engine.
// It deserializes the event and delegates interpretation to handlers.
func (e *Engine) HandleMessage(ctx context.Context, msg kafka.Message) error {
	ev, err := events.Deserialize(msg.Value)
	if err != nil {
		return err
	}
	return HandleEvent(ctx, e.cmdCh, ev)
}
