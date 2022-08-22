package events_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis/v9"
	events "github.com/zipzoft/events-go"
)

func TestRedisChannel(t *testing.T) {
	opt, err := redis.ParseURL("redis://localhost:6379/0")
	if err != nil {
		t.Fatal(err)
	}

	expectMessage := func(expected string, got chan string) {
		select {
		case msg := <-got:
			if msg != expected {
				t.Errorf("expected %s, got %s", expected, msg)
			}
		case <-time.After(2 * time.Second):
			t.Error("timeout")
		}
	}

	channel := events.NewRedisChannel(context.Background(), opt)

	t.Cleanup(func() {
		_ = channel.Close()
	})

	t.Run("Should be publish to single subscriber success", func(t *testing.T) {
		message := "Hello World"
		got := make(chan string)

		_, err := channel.Subscribe("test", func(topic string, message interface{}) {
			if topic != "test" {
				t.Errorf("Expected topic to be 'test', got '%s'", topic)
			}

			got <- message.(string)
		})

		if err != nil {
			t.Errorf("Expected no error, got %s", err)
		}

		if err := channel.Publish("test", message); err != nil {
			t.Errorf("Expected no error, got %s", err)
		}

		expectMessage(message, got)
	})

	t.Run("Should unsubscribe success", func(t *testing.T) {
		message := "Hello World"
		got := make(chan string)

		subscriber, err := channel.Subscribe("test", func(topic string, message interface{}) {
			if topic != "test" {
				t.Errorf("Expected topic to be 'test', got '%s'", topic)
			}

			got <- message.(string)
		})
		if err != nil {
			t.Errorf("Expected no error, got %s", err)
		}

		if err := channel.Publish("test", message); err != nil {
			t.Errorf("Expected no error, got %s", err)
		}

		expectMessage(message, got)

		if err := subscriber.Unsubscribe(); err != nil {
			t.Errorf("Expected no error, got %s", err)
		}

		// Reset message
		got <- "empty"

		time.Sleep(1 * time.Second)

		if err := channel.Publish("test", message); err != nil {
			t.Errorf("Expected no error, got %s", err)
		}

		expectMessage("empty", got)
	})

	t.Run("Should be publish to multiple subscriber success", func(t *testing.T) {
		message := "Hello World"
		got1 := make(chan string)
		got2 := make(chan string)

		_, _ = channel.Subscribe("test", func(topic string, message interface{}) {
			if topic != "test" {
				t.Errorf("Expected topic to be 'test', got '%s'", topic)
			}

			got1 <- message.(string)
		})

		_, _ = channel.Subscribe("test", func(topic string, message interface{}) {
			if topic != "test" {
				t.Errorf("Expected topic to be 'test', got '%s'", topic)
			}

			got2 <- message.(string)
		})

		if err != nil {
			t.Errorf("Expected no error, got %s", err)
		}

		if err := channel.Publish("test", message); err != nil {
			t.Errorf("Expected no error, got %s", err)
		}

		expectMessage(message, got1)
		expectMessage(message, got2)
	})

	t.Run("Should be publish to multiple subscriber with different topic success", func(t *testing.T) {
		message1 := "Hello World"
		message2 := "Hello World 2"
		got1 := make(chan string)
		got2 := make(chan string)

		_, _ = channel.Subscribe("test", func(topic string, message interface{}) {
			if topic != "test" {
				t.Errorf("Expected topic to be 'test', got '%s'", topic)
			}

			got1 <- message.(string)
		})

		_, _ = channel.Subscribe("test2", func(topic string, message interface{}) {
			if topic != "test2" {
				t.Errorf("Expected topic to be 'test2', got '%s'", topic)
			}

			got2 <- message.(string)
		})

		if err := channel.Publish("test", message1); err != nil {
			t.Errorf("Expected no error, got %s", err)
		}

		expectMessage(message1, got1)
		expectMessage("", got2)

		if err := channel.Publish("test2", message2); err != nil {
			t.Errorf("Expected no error, got %s", err)
		}

		expectMessage(message1, got1)
		expectMessage(message2, got2)
	})
}
