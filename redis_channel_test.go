package events_test

import (
	"testing"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/zipzoft/events-go"
)

var _ events.Event = (*testEvent1)(nil)

type testEvent1 struct {
	message string
}

// BroadcastOnChannels implements events.Event
func (evt *testEvent1) BroadcastOnChannels() []string {
	return []string{"redis"}
}

// Topic implements events.Event
func (evt *testEvent1) Topic() string {
	return "test"
}

// Payload implements events.Event
func (evt *testEvent1) Payload() interface{} {
	return evt.message
}

var _ events.Event = (*testEvent2)(nil)

type testEvent2 struct {
	message string
}

// BroadcastOnChannels implements events.Event
func (evt *testEvent2) BroadcastOnChannels() []string {
	return []string{"redis"}
}

// Topic implements events.Event
func (evt *testEvent2) Topic() string {
	return "test2"
}

// Payload implements events.Event
func (evt *testEvent2) Payload() interface{} {
	return evt.message
}

func TestRedisChannel(t *testing.T) {
	opt := &redis.Options{
		Addr: "localhost:6379",
	}

	events.RegisterChannel("redis", events.NewRedisChannel(opt))

	t.Run("Single event", func(t *testing.T) {
		expectMessage := "test message 1"

		gotEvent := make(chan events.Event, 1)
		got := make(chan string, 1)

		events.Listen("test", func(evt events.Event) error {
			gotEvent <- evt
			return nil
		})

		if err := events.Dispatch(&testEvent1{message: expectMessage}); err != nil {
			t.Fatal(err)
		}

		select {
		case evt := <-gotEvent:
			if _, ok := evt.(*testEvent1); !ok {
				t.Errorf("expect event type %T, got %T", &testEvent1{}, evt)
			}
			break

		case <-got:
			t.Error("expect event, got message")
			break

		case <-time.After(2 * time.Second):
			t.Error("expect event, got timeout")
		}
	})

	t.Run("Subscribe", func(t *testing.T) {
		expectMessage := "test message"

		subscriber := events.SubscribeOn(events.NewRedisChannel(opt))
		subscriber.BindEvent(func(topic string, message interface{}) events.Event {
			if topic == "test" {
				return &testEvent1{message: message.(string)}
			}

			return nil
		})

		gotEvent := make(chan events.Event, 1)
		got := make(chan string, 1)

		events.Listen("test", func(evt events.Event) error {
			gotEvent <- evt
			return nil
		})

		if err := events.Dispatch(&testEvent1{message: expectMessage}); err != nil {
			t.Fatal(err)
		}

		select {
		case evt := <-gotEvent:
			if _, ok := evt.(*testEvent1); !ok {
				t.Errorf("expect event type %T, got %T", &testEvent1{}, evt)
			}

			if evt.Payload().(string) != expectMessage {
				t.Errorf("expect message %s, got %s", expectMessage, evt.Payload())
			}
			break

		case <-got:
			t.Error("expect event, got message")
			break

		case <-time.After(2 * time.Second):
			t.Error("expect event, got timeout")
		}
	})
}
