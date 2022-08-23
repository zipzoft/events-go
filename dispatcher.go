package events

type Dispatcher struct {
	channels  map[string]BroadcastChannel
	listeners map[string][]*Listener
}

func (dispatcher *Dispatcher) RegisterChannel(name string, channel BroadcastChannel) {
	dispatcher.channels[name] = channel
}

func (dispatcher *Dispatcher) Dispatch(event Event) error {
	data := event.Payload()

	for _, listener := range dispatcher.listeners[event.Topic()] {
		go func(listener *Listener) {
			(*listener)(event)
		}(listener)
	}

	// Broadcast
	for _, channel := range event.BroadcastOnChannels() {
		if ch, ok := _dispatcher.channels[channel]; ok {
			err := ch.Publish(event.Topic(), data)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (dispatcher *Dispatcher) UnregisterChannel(name string) {
	delete(dispatcher.channels, name)
}

func (dispatcher *Dispatcher) Listen(topic string, listener ...Listener) {
	for _, l := range listener {
		dispatcher.listeners[topic] = append(dispatcher.listeners[topic], &l)
	}
}

func (dispatcher *Dispatcher) Unlisten(topic string, listener *Listener) {
	for i, l := range dispatcher.listeners[topic] {
		if l == listener {
			dispatcher.listeners[topic] = append(dispatcher.listeners[topic][:i], dispatcher.listeners[topic][i+1:]...)
			break
		}
	}
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		channels:  make(map[string]BroadcastChannel),
		listeners: make(map[string][]*Listener),
	}
}
