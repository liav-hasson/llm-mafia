// This file containes player and roles structs and supporting methods

package domain

import (
	"errors"
	"fmt"
	"sync"
)

// --- Player struct --- //

// Player holds base player data
type Player struct {
	ID    string
	Name  string
	Role  Role
	Alive bool
	// TODO: add personallity trait (e.g. timid, agressive, nuetral...)
}

// NewPlayer creates a new player with the provided id, name, and role.
// Both id and name are required - auto-generation has been removed.
// Returns an error if id or name is empty.
func NewPlayer(id, name string, role Role) (*Player, error) {
	if id == "" {
		return nil, errors.New("player id is required")
	}
	if name == "" {
		return nil, errors.New("player name is required")
	}
	return &Player{
		ID:    id,
		Name:  name,
		Role:  role,
		Alive: true,
	}, nil
}

// --- Role enum --- //

// Role represents possible player roles
type Role int

const (
	RoleUnknown Role = iota
	RoleVillager
	RoleMafia
	RoleDoctor
	RoleSheriff
)

func (r Role) String() string {
	switch r {
	case RoleUnknown:
		return "unknown"
	case RoleVillager:
		return "villager"
	case RoleMafia:
		return "mafia"
	case RoleDoctor:
		return "doctor"
	case RoleSheriff:
		return "sheriff"
	default:
		return "invalid"
	}
}

// --- ID Generation --- //

// package-level counter for generating player IDs
// using mutex to ensure thread-safety (concurrent operations)
var (
	playerCounter int
	counterMutex  sync.Mutex
)

// CreatePlayerID generates sequential player IDs: player-1, player-2, etc.
// Thread-safe: uses mutex to protect counter from race conditions.
func CreatePlayerID() string {
	counterMutex.Lock()
	defer counterMutex.Unlock()

	playerCounter++
	return fmt.Sprintf("player-%d", playerCounter)
}

// ResetPlayerCounter resets the ID counter to zero.
// This is intended for use in tests to ensure clean state between test runs.
func ResetPlayerCounter() {
	counterMutex.Lock()
	defer counterMutex.Unlock()
	playerCounter = 0
}

// --- Player helpers --- //

// --- Role helpers --- //

// IsVillagerTeam returns true if the role is on the villager team
// Note: Not yet called in engine commands, kept for future use
func (r Role) IsVillagerTeam() bool {
	return r == RoleVillager ||
		r == RoleDoctor ||
		r == RoleSheriff
}

func (r Role) IsMafiaTeam() bool {
	return r == RoleMafia
}

func (r Role) HasNightAction() bool {
	return r == RoleMafia ||
		r == RoleDoctor ||
		r == RoleSheriff
}
