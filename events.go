package events

var broker *Broker

func init() {
	broker = New("memory", NewMemoryChannel())
}

func Publish(topic string, message interface{}) error {
	return broker.Publish(topic, message)
}

func Subscribe(topic string, receiver MessageReceiver) (*Subscriber, error) {
	return broker.Subscribe(topic, receiver)
}

func Unsubscribe(topic string, receiver *MessageReceiver) error {
	return broker.Unsubscribe(topic, receiver)
}

func Receivers(topic string) []*MessageReceiver {
	return broker.Receivers(topic)
}

func Channel() MessageChannel {
	return broker.Channel()
}

func On(channelName string) MessageChannel {
	return broker.channels[channelName]
}
