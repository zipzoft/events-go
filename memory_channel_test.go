package events_test

import (
	"testing"

	events "github.com/zipzoft/events-go"
)

func TestMemoryChannel(t *testing.T) {

	t.Run("Should be publish to single subscriber success", func(t *testing.T) {
		channel := events.NewMemoryChannel()
		expectMessage := "Hello World"
		gotMessage := ""

		channel.Subscribe("test", func(topic string, message interface{}) {
			if topic != "test" {
				t.Errorf("Expected topic to be 'test', got '%s'", topic)
			}

			gotMessage = message.(string)
		})

		channel.Publish("test", expectMessage)

		if gotMessage != expectMessage {
			t.Errorf("Expected message %s, but got %s", expectMessage, gotMessage)
		}
	})

	t.Run("Should be publish to multiple subscriber success", func(t *testing.T) {
		channel := events.NewMemoryChannel()
		expectMessage := "Hello World"
		gotMessage1 := ""
		gotMessage2 := ""

		channel.Subscribe("test", func(topic string, message interface{}) {
			if topic != "test" {
				t.Errorf("Expected topic to be 'test', got '%s'", topic)
			}

			gotMessage1 = message.(string)
		})

		channel.Subscribe("test", func(topic string, message interface{}) {
			if topic != "test" {
				t.Errorf("Expected topic to be 'test', got '%s'", topic)
			}

			gotMessage2 = message.(string)
		})

		channel.Publish("test", expectMessage)

		if gotMessage1 != expectMessage {
			t.Errorf("Expected message %s, but got %s", expectMessage, gotMessage1)
		}

		if gotMessage2 != expectMessage {
			t.Errorf("Expected message %s, but got %s", expectMessage, gotMessage2)
		}
	})

	t.Run("Should be publish to multiple subscriber with wildcard success", func(t *testing.T) {
		channel := events.NewMemoryChannel()
		expectMessage := "Hello World"
		gotMessage1 := ""
		gotMessage2 := ""

		channel.Subscribe("test", func(topic string, message interface{}) {
			if topic != "test" {
				t.Errorf("Expected topic to be 'test', got '%s'", topic)
			}

			gotMessage1 = message.(string)
		})

		channel.Subscribe("*", func(topic string, message interface{}) {
			if topic != "test" {
				t.Errorf("Expected topic to be 'test', got '%s'", topic)
			}

			gotMessage2 = message.(string)
		})

		channel.Publish("test", expectMessage)

		if gotMessage1 != expectMessage {
			t.Errorf("Expected message %s, but got %s", expectMessage, gotMessage1)
		}

		if gotMessage2 != expectMessage {
			t.Errorf("Expected message %s, but got %s", expectMessage, gotMessage2)
		}
	})

	t.Run("Should be publish to multiple subscriber with regex success", func(t *testing.T) {
		channel := events.NewMemoryChannel()
		expectMessage := "Hello World"
		gotMessage1 := ""
		gotMessage2 := ""

		channel.Subscribe("test", func(topic string, message interface{}) {
			if topic != "test" {
				t.Errorf("Expected topic to be 'test', got '%s'", topic)
			}

			gotMessage1 = message.(string)
		})

		channel.Subscribe("/^t.+/", func(topic string, message interface{}) {
			if topic != "test" {
				t.Errorf("Expected topic to be 'test', got '%s'", topic)
			}

			gotMessage2 = message.(string)
		})

		channel.Publish("test", expectMessage)

		if gotMessage1 != expectMessage {
			t.Errorf("Expected message %s, but got %s", expectMessage, gotMessage1)
		}

		if gotMessage2 != expectMessage {
			t.Errorf("Expected message %s, but got %s", expectMessage, gotMessage2)
		}
	})

	t.Run("Should be publish to multiple subscriber with comma success", func(t *testing.T) {
		channel := events.NewMemoryChannel()
		expectMessage := "Hello World"
		gotMessage1 := ""
		gotMessage2 := ""
		gotMessage3 := ""

		channel.Subscribe("test", func(topic string, message interface{}) {
			if topic != "test" {
				t.Errorf("Expected topic to be 'test', got '%s'", topic)
			}

			gotMessage1 = message.(string)
		})

		channel.Subscribe("a, test2", func(topic string, message interface{}) {
			if topic != "test2" {
				t.Errorf("Expected topic to be 'test', got '%s'", topic)
			}

			gotMessage2 = message.(string)
		})

		channel.Subscribe("a", func(topic string, message interface{}) {
			t.Errorf("Should not be called")

			gotMessage3 = message.(string)
		})

		channel.Publish("test,test2", expectMessage)

		if gotMessage1 != expectMessage {
			t.Errorf("Expected message %s, but got %s", expectMessage, gotMessage1)
		}

		if gotMessage2 != expectMessage {
			t.Errorf("Expected message %s, but got %s", expectMessage, gotMessage2)
		}

		if gotMessage3 != "" {
			t.Errorf("Expected message %s, but got %s", expectMessage, gotMessage3)
		}
	})

	t.Run("Should be unsubscribe success", func(t *testing.T) {
		topic := "test"
		channel := events.NewMemoryChannel()

		if len(channel.Receivers(topic)) != 0 {
			t.Errorf("Expected receivers to be 0, got %d", len(channel.Receivers(topic)))
		}

		subscriber, _ := channel.Subscribe(topic, func(topic string, message interface{}) {
			// gotMessage = message.(string)
		})

		if len(channel.Receivers(topic)) != 1 {
			t.Errorf("Expected receivers to be 1, got %d", len(channel.Receivers(topic)))
		}

		subscriber.Unsubscribe()

		if len(channel.Receivers(topic)) != 0 {
			t.Errorf("Expected receivers to be 0, got %d", len(channel.Receivers(topic)))
		}
	})

	t.Run("Should be subscribe with wildcard success", func(t *testing.T) {
		topic := "test"
		channel := events.NewMemoryChannel()
		expectMessage := "Hello World"
		gotMessage := ""

		if len(channel.Receivers(topic)) != 0 {
			t.Errorf("Expected receivers to be 0, got %d", len(channel.Receivers(topic)))
		}

		channel.Subscribe("*", func(topic string, message interface{}) {
			gotMessage = message.(string)
		})

		if len(channel.Receivers(topic)) != 0 {
			t.Errorf("Expected receivers to be 0, got %d", len(channel.Receivers(topic)))
		}

		if len(channel.Receivers("*")) != 1 {
			t.Errorf("Expected receivers to be 1, got %d", len(channel.Receivers("*")))
		}

		channel.Publish(topic, expectMessage)

		if gotMessage != expectMessage {
			t.Errorf("Expected message %s, but got %s", expectMessage, gotMessage)
		}
	})
}
