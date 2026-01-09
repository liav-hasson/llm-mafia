package names

import (
	"errors"
	"sync"
)

// ErrNoMoreNames is returned when all available names have been used.
var ErrNoMoreNames = errors.New("no more names available")

// Generator assigns names to players sequentially from a provided list.
// It is thread-safe and tracks which names have been used.
type Generator struct {
	names   []string
	counter int
	mu      sync.Mutex
}

// NewGenerator creates a new name generator with the provided list of names.
// Returns an error if the names list is empty.
func NewGenerator(names []string) (*Generator, error) {
	if len(names) == 0 {
		return nil, errors.New("names list must not be empty")
	}

	return &Generator{
		names: names,
	}, nil
}

// Next returns the next available name.
// Returns ErrNoMoreNames if all names have been used.
// Thread-safe: uses mutex to protect counter.
func (ng *Generator) Next() (string, error) {
	ng.mu.Lock()
	defer ng.mu.Unlock()

	if ng.counter >= len(ng.names) {
		return "", ErrNoMoreNames
	}

	name := ng.names[ng.counter]
	ng.counter++
	return name, nil
}

// Reset resets the counter to zero, allowing names to be reused.
// This is primarily intended for testing.
func (ng *Generator) Reset() {
	ng.mu.Lock()
	defer ng.mu.Unlock()
	ng.counter = 0
}

// Remaining returns the number of unused names.
func (ng *Generator) Remaining() int {
	ng.mu.Lock()
	defer ng.mu.Unlock()
	return len(ng.names) - ng.counter
}
