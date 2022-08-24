package events

type Subscriber struct {
	creators  []CreateEvent
	listeners map[string][]*Listener
	channel   BroadcastChannel
}

type CreateEvent func(topic string, message interface{}) Event

func (s *Subscriber) BindEvent(creator CreateEvent) {
	s.creators = append(s.creators, creator)
}

func (s *Subscriber) Subscribe(topic ...string) error {
	return s.channel.Subscribe(topic...)
}
