package events

// BroadcastChannel is the interface for message channel
type BroadcastChannel interface {
	// Publish a message to a topic
	Publish(topic string, message interface{}) error

	// Subscribe to a topic
	Subscribe(topics ...string) error
}
