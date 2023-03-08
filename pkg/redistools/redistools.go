package redistools

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	redis "github.com/redis/go-redis/v9"
)

type ServiceType uint

const pubsub string = "almagest"

type PSMessage struct {
	// Service should be on all messages
	Service string `json:"service,omitempty"`
	// Channel can be used for discord channel
	Channel string `json:"channel,omitempty"`
	// Content to be posted to discord
	Content string `json:"content,omitempty"`
	// MessageID is a uuid used to track transactions
	MessageID string `json:"mid,omitempty"`
	// ResponseTo is populated with the originating message id to indicate its a response or follow up
	ResponseTo string `json:"rid,omitempty"`
	// Health is used when responding to a healthcheck
	Health string `json:"health,omitempty"`
}

type AlmagestRedisClient struct {
	Subs []*redis.PubSub
	c    *redis.Client
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
	if m.MessageID == "" {
		log.Error("MessageID is required")
		return fmt.Errorf("MessageID is required")
	}
	message, err := json.Marshal(m)
	if err != nil {
		return err
	}
	log.Debug("Publishing message", "content", string(message))
	return a.c.Publish(context.Background(), pubsub, string(message)).Err()
}

// PostStatus makes a status message
func (a AlmagestRedisClient) PostStatus(s, v, mid string) error {
	dsm := PSMessage{
		Service:    "healthcheck",
		ResponseTo: mid,
		Content:    fmt.Sprintf("%s|%s|ok", s, v),
	}
	return a.Publish(dsm)
}

func (a AlmagestRedisClient) PublishWatchdog(service string) error {
	m := PSMessage{
		Service: "watchdog",
		Content: service,
	}
	t := time.NewTicker(30 * time.Second)
	for {
		_ = <-t.C
		m.MessageID = uuid.NewString()
		a.Publish(m)
	}
}
