package events

import (
	"errors"
	"regexp"
	"strings"
)

var _ MessageChannel = (*memoryMessageChannel)(nil)

type memoryMessageChannel struct {
	topics map[string][]*MessageReceiver
}

func (channel *memoryMessageChannel) Receivers(topic string) []*MessageReceiver {
	return channel.topics[topic]
}

// Publish implements MessageChannel
func (channel *memoryMessageChannel) Publish(topic string, message interface{}) error {
	if topic == "*" {
		return errors.New("topic cannot be *")
	}

	if strings.HasPrefix(topic, "/") && strings.HasSuffix(topic, "/") {
		return errors.New("topic cannot be regex")
	}

	// Split the topic by comma
	topics := strings.Split(topic, ",")

	for _, t := range topics {
		t = strings.TrimSpace(t)

		if receivers, ok := channel.topics[t]; ok {
			for _, receiver := range receivers {
				(*receiver)(t, message)
			}
		}
	}

	if receivers, ok := channel.topics["*"]; ok {
		for _, receiver := range receivers {
			(*receiver)(topic, message)
		}
	}

	return nil
}

// Subscribe implements MessageChannel
func (channel *memoryMessageChannel) Subscribe(topic string, receiver MessageReceiver) (*Subscriber, error) {
	topic = strings.TrimSpace(topic)
	if len(topic) == 0 {
		return nil, errors.New("topic cannot be empty")
	}

	subscriber := NewSubscriber(channel, topic, &receiver)

	// If the topic is need match with regex
	if strings.HasPrefix(topic, "/") && strings.HasSuffix(topic, "/") {
		for t := range channel.topics {
			reg := strings.Trim(topic, "/")
			if match, _ := regexp.MatchString(reg, t); match {
				channel.topics[t] = append(channel.topics[t], &receiver)
			}
		}

		return subscriber, nil
	}

	// Split the topic by comma
	topics := strings.Split(topic, ",")
	for _, t := range topics {
		t = strings.TrimSpace(t)
		if len(t) == 0 {
			continue
		}

		channel.topics[t] = append(channel.topics[t], &receiver)
	}

	return subscriber, nil
}

// Unsubscribe implements MessageChannel
func (channel *memoryMessageChannel) Unsubscribe(topic string, receiver *MessageReceiver) error {
	topic = strings.TrimSpace(topic)
	if len(topic) == 0 {
		return errors.New("topic cannot be empty")
	}

	// If the topic is need match with regex
	if strings.HasPrefix(topic, "/") && strings.HasSuffix(topic, "/") {
		for t := range channel.topics {
			reg := strings.Trim(topic, "/")
			if match, _ := regexp.MatchString(reg, t); match {
				for i, _receiver := range channel.topics[t] {
					if _receiver == receiver {
						channel.topics[t] = append(channel.topics[t][:i], channel.topics[t][i+1:]...)
						break
					}
				}
			}
		}

		return nil
	}

	// Split the topic by comma
	topics := strings.Split(topic, ",")
	for _, t := range topics {
		t = strings.TrimSpace(t)
		if len(t) == 0 {
			continue
		}

		if receivers, ok := channel.topics[t]; ok {
			for i, _receiver := range receivers {
				if _receiver == receiver {
					channel.topics[t] = append(receivers[:i], receivers[i+1:]...)

					if len(channel.topics[t]) == 0 {
						delete(channel.topics, t)
					}
					break
				}
			}
		}
	}

	return nil
}

func (channel *memoryMessageChannel) Close() error {
	channel.topics = make(map[string][]*MessageReceiver)
	return nil
}

func NewMemoryChannel() *memoryMessageChannel {
	return &memoryMessageChannel{
		topics: make(map[string][]*MessageReceiver),
	}
}
