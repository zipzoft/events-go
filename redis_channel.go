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

func NewRedisChannel() *RedisChannel {
	return &RedisChannel{
		context: context.Background(),
		client: redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		}),
	}
}
