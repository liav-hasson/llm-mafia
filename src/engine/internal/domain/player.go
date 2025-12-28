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

// NewPlayer creates a new player with the given parameters
// If id is empty, generates one using CreatePlayerID()
// If name is empty, generates one using CreatePlayerName()
func NewPlayer(id, name string, role Role) (*Player, error) {
	if id == "" {
		id = CreatePlayerID()
	}
	if name == "" {
		var err error
		name, err = CreatePlayerName()
		if err != nil {
			return nil, err
		}
	}
	return &Player{
		ID:    id,
		Name:  name,
		Role:  role,
		Alive: true, // new players start alive
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

// --- ID and Name Generation --- //

// package-level counters for generating player IDs and names
// using mutex to ensure thread-safety (concurrent Kafka events)
var (
	playerCounter int
	nameCounter   int
	counterMutex  sync.Mutex
)

// ErrNoMoreNames is returned when all available names have been used
var ErrNoMoreNames = errors.New("no more names available")

// availableNames is a list of names to assign to players
var availableNames = []string{
	"Gilbert McDonald", "Dorothy Bird",
	"Ernest Preston", "Vincent Schultz",
	"Joanne Sloan", "Lana Moran",
	"Adrienne Fuller", "Greg Bennett",
	"Curt Simon", "Rachel McMillan",
	"Dustin Eastman", "Willard Mendez",
}

// CreatePlayerID generates sequential player IDs: player-1, player-2, etc.
// Thread-safe: uses mutex to protect counter from race conditions
func CreatePlayerID() string {
	counterMutex.Lock()
	defer counterMutex.Unlock()

	playerCounter++
	return fmt.Sprintf("player-%d", playerCounter)
}

// CreatePlayerName returns sequential names from availableNames
// Thread-safe: uses mutex to protect counter
// Returns ErrNoMoreNames if all names have been used
func CreatePlayerName() (string, error) {
	counterMutex.Lock()
	defer counterMutex.Unlock()

	if nameCounter >= len(availableNames) {
		return "", ErrNoMoreNames
	}

	name := availableNames[nameCounter]
	nameCounter++
	return name, nil
}

// ResetPlayerCounters resets the ID and name counters to zero
// This is intended for use in tests to ensure clean state between test runs
func ResetPlayerCounters() {
	counterMutex.Lock()
	defer counterMutex.Unlock()

	playerCounter = 0
	nameCounter = 0
}

// --- Player helpers --- //

func (p Player) IsAlive() bool {
	return p.Alive
}

// --- Role helpers --- //

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
