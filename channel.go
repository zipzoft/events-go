package events

// MessageChannel is the interface for message channel
type MessageChannel interface {
	// Receivers returns the receivers of the topic
	Receivers(topic string) []*MessageReceiver

	// Publish a message to a topic
	Publish(topic string, message interface{}) error

	// Subscribe to a topic
	Subscribe(topic string, receiver MessageReceiver) (*Subscriber, error)

	// Unsubscribe from a topic
	Unsubscribe(topic string, receiver *MessageReceiver) error

	// Close the channel
	Close() error
}

// MessageReceiver is a function that receives a message
type MessageReceiver func(topic string, message interface{})

type Subscriber struct {
	channel  MessageChannel
	topic    string
	receiver *MessageReceiver
}

func (subscriber *Subscriber) Unsubscribe() error {
	return subscriber.channel.Unsubscribe(subscriber.topic, subscriber.receiver)
}

func NewSubscriber(channel MessageChannel, topic string, receiver *MessageReceiver) *Subscriber {
	return &Subscriber{
		channel:  channel,
		topic:    topic,
		receiver: receiver,
	}
}
