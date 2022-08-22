package events

import (
	"context"
	"errors"
	"strings"

	"github.com/go-redis/redis/v9"
)

var _ MessageChannel = (*redisChannel)(nil)

type redisChannel struct {
	ctx       context.Context
	client    *redis.Client
	listening bool

	pubsub *redis.PubSub

	topics map[string][]*MessageReceiver
}

// Subscribe implements MessageChannel
func (channel *redisChannel) Subscribe(topic string, receiver MessageReceiver) (*Subscriber, error) {
	topic = strings.TrimSpace(topic)
	if len(topic) == 0 {
		return nil, errors.New("topic cannot be empty")
	}

	subscriber := NewSubscriber(channel, topic, &receiver)

	if _, ok := channel.topics[topic]; !ok {
		channel.topics[topic] = make([]*MessageReceiver, 0)
	}

	channel.topics[topic] = append(channel.topics[topic], &receiver)

	channel.listen()

	return subscriber, nil
}

// Publish implements MessageChannel
func (channel *redisChannel) Publish(topic string, message interface{}) error {
	topic = strings.TrimSpace(topic)
	if len(topic) == 0 {
		return errors.New("topic cannot be empty")
	}

	if err := channel.client.Publish(channel.ctx, topic, message).Err(); err != nil {
		return err
	}

	return nil
}

// Receivers implements MessageChannel
func (channel *redisChannel) Receivers(topic string) []*MessageReceiver {
	topic = strings.TrimSpace(topic)
	if len(topic) == 0 {
		return nil
	}

	if receivers, ok := channel.topics[topic]; ok {
		return receivers
	}

	return make([]*MessageReceiver, 0)
}

// Unsubscribe implements MessageChannel
func (channel *redisChannel) Unsubscribe(topic string, receiver *MessageReceiver) error {
	topic = strings.TrimSpace(topic)
	if len(topic) == 0 {
		return errors.New("topic cannot be empty")
	}

	if _, ok := channel.topics[topic]; ok {
		delete(channel.topics, topic)
		channel.listen()
	}

	return nil
}

func (channel *redisChannel) Close() error {
	if channel.pubsub != nil {
		return channel.pubsub.Close()
	}

	channel.topics = make(map[string][]*MessageReceiver)
	channel.listening = false

	return nil
}

func (channel *redisChannel) listen() {
	topics := make([]string, 0)
	for topic := range channel.topics {
		topics = append(topics, topic)
	}

	if len(topics) == 0 {
		channel.Close()
	} else {
		if channel.pubsub == nil {
			channel.pubsub = channel.client.Subscribe(channel.ctx, topics...)
		} else {
			_ = channel.pubsub.Close()

			_ = channel.pubsub.Subscribe(channel.ctx, topics...)
		}
	}

	if channel.pubsub != nil {
		go func() {
			for {
				if channel.pubsub == nil {
					return
				}

				msg, err := channel.pubsub.ReceiveMessage(channel.ctx)
				if err != nil {
					continue
				}

				if receivers, ok := channel.topics[msg.Channel]; ok {
					for _, receiver := range receivers {
						(*receiver)(msg.Channel, msg.Payload)
					}
				}
			}
		}()
	}
}

func NewRedisChannel(ctx context.Context, opt *redis.Options) *redisChannel {
	channel := &redisChannel{
		ctx:       ctx,
		client:    redis.NewClient(opt),
		topics:    make(map[string][]*MessageReceiver),
		listening: false,
		pubsub:    nil,
	}

	return channel
}

func ParseFromURI(uri string) *redis.Options {
	opt, err := redis.ParseURL(uri)
	if err != nil {
		panic(err)
	}

	return opt
}
