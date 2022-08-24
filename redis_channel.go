package events

import (
	"context"

	"github.com/go-redis/redis/v9"
)

var _ BroadcastChannel = (*RedisChannel)(nil)

type RedisChannel struct {
	context context.Context
	client  *redis.Client
}

// Publish implements BroadcastChannel
func (channel *RedisChannel) Publish(topic string, message interface{}) error {
	return channel.client.Publish(channel.context, topic, message).Err()
}

// Subscribe implements BroadcastChannel
func (channel *RedisChannel) Subscribe(topics ...string) error {
	pubsub := channel.client.Subscribe(channel.context, topics...)
	if _, err := pubsub.Receive(channel.context); err != nil {
		return err
	}

	go func() {
		for msg := range pubsub.Channel() {
			topic := msg.Channel
			message := msg.Payload

			for _, creator := range _subscriber.creators {
				if evt := creator(topic, message); evt != nil {
					if evt.Topic() == topic {
						_dispatcher.Dispatch(evt)
					}
				}
			}
		}
	}()

	return nil
}

func NewRedisChannel(opt *redis.Options) *RedisChannel {
	return &RedisChannel{
		context: context.Background(),
		client:  redis.NewClient(opt),
	}
}
