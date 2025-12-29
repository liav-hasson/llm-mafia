package events

// file responsible for go <-> json convertions
import (
	"encoding/json"
)

// encodes all Go structs to json
func Marshal(event any) ([]byte, error) {
	return json.Marshal(event)
}

// Unmarshal funcs for all events the engine receives
// returns nil pointer on error for explicit failure signaling
func UnmarshalAllChatMessage(data []byte) (*AllChatMessage, error) {
	var event AllChatMessage
	err := json.Unmarshal(data, &event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func UnmarshalMafiaChatMessage(data []byte) (*MafiaChatMessage, error) {
	var event MafiaChatMessage
	err := json.Unmarshal(data, &event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func UnmarshalVoteSubmitted(data []byte) (*VoteSubmitted, error) {
	var event VoteSubmitted
	err := json.Unmarshal(data, &event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func UnmarshalNightAction(data []byte) (*NightAction, error) {
	var event NightAction
	err := json.Unmarshal(data, &event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func UnmarshalPlayerThoughts(data []byte) (*PlayerThoughts, error) {
	var event PlayerThoughts
	err := json.Unmarshal(data, &event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}
