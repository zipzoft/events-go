package events

type Broker struct {
	// name is Current channel name
	name string

	channels map[string]MessageChannel
}

func (broker *Broker) Publish(topic string, message interface{}) error {
	return broker.Channel().Publish(topic, message)
}

func (broker *Broker) Subscribe(topic string, receiver MessageReceiver) (*Subscriber, error) {
	return broker.Channel().Subscribe(topic, receiver)
}

func (broker *Broker) Unsubscribe(topic string, receiver *MessageReceiver) error {
	return broker.Channel().Unsubscribe(topic, receiver)
}

func (broker *Broker) Receivers(topic string) []*MessageReceiver {
	return broker.Channel().Receivers(topic)
}

func (broker *Broker) Channel() MessageChannel {
	return broker.channels[broker.name]
}

func (broker *Broker) Extend(channelName string, channel MessageChannel) {
	broker.channels[channelName] = channel
}

func New(channelName string, channel MessageChannel) *Broker {
	broker := &Broker{
		name:     channelName,
		channels: make(map[string]MessageChannel),
	}

	if _, ok := broker.channels[channelName]; !ok {
		broker.channels[channelName] = channel
	}

	return broker
}
