# Event Bus (Go)

events GO is a compact yet powerful distributed stream processing library for Go. It is designed to be used in conjunction with the [Redis Pub/Sub](http://redis.io/topics/pubsub) protocol.

## Installation

```bash
go get -u github.com/zipzoft/events-go
```

## Usage

```go
package main

import (
    "fmt"
    "github.com/zipzoft/events-go"
    "time"
)

func main() {
    opt, err := redis.ParseURL("redis://localhost:6379")
    if err != nil {
        panic(err)
    }

    // Create a new transport
	events.New("redis", events.NewRedisChannel(context.TODO(), opt))

    // Subscribe to the channel
    subscriber, err := events.Subscribe("some-event", func(topic string, message interface{}) {
        fmt.Println("Received message:", message)
    })

    // Publish a message to the channel
    err := events.Publish("some-event", "Hello World!")

    // Unsubscribe from the channel
    err := subscriber.Unsubscribe()
    
    // ...
}
```
