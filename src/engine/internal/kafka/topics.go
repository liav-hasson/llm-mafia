package kafka

// Topic names.
// These represent durable Kafka logs, NOT event semantics.
const (
	// EngineEventsTopic is the stream of authoritative game events
	// emitted by the engine and consumed by players.
	EngineEventsTopic = "game.engine.events"

	// PlayerActionsTopic is the stream of player intents
	// (votes, night actions, thoughts) consumed by the engine.
	PlayerActionsTopic = "game.player.actions"
)

// Consumer group names.
// These identify who is consuming a topic, not what is being consumed.
const (
	EngineConsumerGroup = "mafia-engine"
)

// GameKey returns the Kafka partition key for a given game.
// All events for the same game MUST use the same key to preserve ordering.
func GameKey(gameID string) []byte {
	return []byte(gameID)
}
