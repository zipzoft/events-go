package events

type Subscriber struct {
	topics    map[string][]Event
	listeners map[string][]*Listener
	channel   BroadcastChannel
}

type CreateEvent func(topic string, message interface{}) Event

func (s *Subscriber) RegisterEvent(evt Event) {
	s.topics[evt.Topic()] = append(s.topics[evt.Topic()], evt)
}

func (s *Subscriber) Subscribe(topic ...string) error {
	return s.channel.Subscribe(topic...)
}
