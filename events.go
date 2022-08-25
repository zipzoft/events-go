package events

var _dispatcher = NewDispatcher()
var _subscriber *Subscriber = nil

func SubscribeOn(channel BroadcastChannel) *Subscriber {
	_subscriber = &Subscriber{
		topics:    make(map[string][]Event),
		listeners: make(map[string][]*Listener),
		channel:   channel,
	}

	return _subscriber
}

func RegisterChannel(name string, channel BroadcastChannel) {
	_dispatcher.RegisterChannel(name, channel)
}

func Listen(topic string, listener ...Listener) {
	_dispatcher.Listen(topic, listener...)

	if _subscriber != nil {
		topics := make([]string, 0)

		for _topic, _ := range _dispatcher.listeners {
			topics = append(topics, _topic)
		}

		_subscriber.channel.Subscribe(topics...)
	}
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

	// OnBroadcastReceive is called when the event is received on a channel
	OnBroadcastReceive(message interface{}) error
}

type Listener func(Event) error
