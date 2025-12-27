// This file containes game phase variables

package domain

// creates type "Phase" with underlying int value
type Phase int

// iota increments each var starting from 0
const (
	PhaseUnknown Phase = iota // debug
	PhaseWaiting
	PhaseNight
	PhaseDay
	PhaseVoting
	PhaseEnded
)

func (p Phase) String() string {
	switch p {
	case PhaseUnknown:
		return "unknown"
	case PhaseWaiting:
		return "waiting"
	case PhaseNight:
		return "night"
	case PhaseDay:
		return "day"
	case PhaseVoting:
		return "voting"
	case PhaseEnded:
		return "ended"
	default:
		return "invalid"
	}
}
