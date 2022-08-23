package events

var _dispatcher = NewDispatcher()

func RegisterChannel(name string, channel BroadcastChannel) {
	_dispatcher.RegisterChannel(name, channel)
}

func Listen(topic string, listener ...Listener) {
	_dispatcher.Listen(topic, listener...)
}

func UnregisterChannel(name string) {
	_dispatcher.UnregisterChannel(name)
}

func Unlisten(topic string, listener *Listener) {
	_dispatcher.Unlisten(topic, listener)
}

func Dispatch(event Event) error {
	return _dispatcher.Dispatch(event)
}

type Event interface {
	// Topic returns the topic of the event
	Topic() string

	// BroadcastOnChannels returns the channels to broadcast the event
	BroadcastOnChannels() []string

	// Payload returns the payload of the event
	Payload() interface{}
}

type Listener func(Event) error
