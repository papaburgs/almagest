package redistools

import (
	"context"
	"encoding/json"

	redis "github.com/redis/go-redis/v9"
)

const pubsub string = "almagest"

type PSMessage struct {
	// Service should be on all messages
	Service string `json:"service,omitempty"`
	// Channel can be used for discord channel
	Channel string `json:"channel,omitempty"`
	// Content to be posted to discord
	Content string `json:"message,omitempty"`
}

type AlmagestRedisClient struct {
	Subs     []*redis.PubSub
	c        *redis.Client
	channels map[string]string
}

// New makes a new AlmagestRedisClient and returns pointer to it
func New() *AlmagestRedisClient {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	arc := AlmagestRedisClient{
		c: rdb,
		channels: map[string]string{
			"discord": "almagest|discord",
			"uptime":  "almagest|uptime",
		},
	}
	return &arc
}

func (a AlmagestRedisClient) Subscribe() <-chan *redis.Message {
	var (
		err error
		ctx context.Context = context.Background()
	)
	pSub := a.c.Subscribe(ctx, pubsub)
	a.Subs = append(a.Subs, pSub)

	// Wait for confirmation that subscription is created before publishing anything.
	_, err = pSub.Receive(ctx)
	if err != nil {
		panic(err)
	}
	return pSub.Channel()
}

// Publish takes in a PSMessage, converts to json and publishes
// it to the almagest channel
func (a AlmagestRedisClient) Publish(m PSMessage) error {
	message, err := json.Marshal(m)
	if err != nil {
		return err
	}

	return a.c.Publish(context.Background(), pubsub, string(message)).Err()
}
