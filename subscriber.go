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

func (s *Subscriber) Subscribe(topics ...string) error {
	if len(topics) == 0 {
		for topic := range s.topics {
			topics = append(topics, topic)
		}
	}

	return s.channel.Subscribe(topics...)
}
