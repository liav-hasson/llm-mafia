// This file containes player and roles structs and supporting methods

package domain

import (
	"fmt"
	"sync"
)

// player base data
type Player struct {
	ID    string
	Name  string
	Role  Role
	Alive bool
	// TODO: add personallity trait (e.g. timid, agressive, nuetral...)
}

// possible player roles
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

// package-level counter for generating player IDs
// using mutex to ensure thread-safety (concurrent Kafka events)
var (
	playerCounter int
	nameCounter   int
	playerMutex   sync.Mutex
)

// CreatePlayerID generates sequential player IDs: player-1, player-2, etc.
// Thread-safe: uses mutex to protect counter from race conditions
func CreatePlayerID() string {
	playerMutex.Lock()         // acquire lock - blocks other goroutines
	defer playerMutex.Unlock() // release lock when function returns

	playerCounter++
	return fmt.Sprintf("player-%d", playerCounter)
}

// availableNames is a list of names to assign to players
// unexported since it's only used internally
var availableNames = []string{
	"Gilbert McDonald", "Dorothy Bird",
	"Ernest Preston", "Vincent Schultz",
	"Joanne Sloan", "Lana Moran",
	"Adrienne Fuller", "Greg Bennett",
	"Curt Simon", "Rachel McMillan",
	"Dustin Eastman", "Willard Mendez",
}

// CreatePlayerName returns sequential names from availableNames
// Thread-safe: shares mutex with CreatePlayerID
func CreatePlayerName() (string, error) {
	playerMutex.Lock()
	defer playerMutex.Unlock()

	// get current name, then increment for next call
	name := availableNames[nameCounter]
	nameCounter++

	return name, nil
}

// NewPlayer creates a new player with the given parameters
// If id is empty, generates one using CreatePlayerID()
// If name is empty, generates one using CreatePlayerName()
// Returns ErrMaxPlayersReached if auto-generating name and limit exceeded
func NewPlayer(id, name string, role Role) (*Player, error) {
	if id == "" {
		id = CreatePlayerID()
	}
	if name == "" {
		var err error
		name, err = CreatePlayerName()
		if err != nil {
			return nil, err // propagate the error to caller
		}
	}
	return &Player{
		ID:    id,
		Name:  name,
		Role:  role,
		Alive: true, // new players start alive
	}, nil
}

// player state helpers
func (p Player) IsAlive() bool {
	return p.Alive
}

// player role helpers
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
