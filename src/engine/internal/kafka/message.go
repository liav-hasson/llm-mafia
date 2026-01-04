package kafka

type Message struct {
	Topic string
	Key   []byte // gameID
	Value []byte // payload

	// maybe add this for metadata: 'Headers map[string][]byte'
	// headers will allow for getting data such as:
	// timestamp, version, event_type, etc...
	// headers dont require parsing logic (deserialization), which allows for:
	// faster filtering/routing, scheme versioning, better logging/tracing
}
