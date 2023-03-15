package redistools

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	redis "github.com/redis/go-redis/v9"
)

type RedisMessageClass uint

const (
	UnknownClass RedisMessageClass = iota
	HealthCheckRequest
	HealthCheckResponse
	SendToDiscordRequest
	Watchdog
)

// ClassifyMessage takes in a redis message, returns a PSMessage,
// with a Classification, or an error
func ClassifyMessage(r *redis.Message) (PSMessage, RedisMessageClass, error) {
	var (
		rmc RedisMessageClass
		m   PSMessage
		err error
	)
	err = json.Unmarshal([]byte(r.Payload), &m)
	if err != nil {
		log.Error("Could not decode message ", "payload", r.Payload)
		return PSMessage{}, UnknownClass, err
	}

	switch m.Service {
	case "healthcheck":
		log.Debug("see healthcheck")
		if m.ResponseTo == "" {
			log.Debug("looks like a new one, setting request")
			rmc = HealthCheckRequest
		} else {
			log.Debug("looks like a response")
			rmc = HealthCheckResponse
		}
	case "discord":
		rmc = SendToDiscordRequest
	case "watchdog":
		rmc = Watchdog
	default:
		rmc = UnknownClass
	}

	return m, rmc, nil
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

// PostStatus makes a status message
func (a AlmagestRedisClient) PostStatus(s, v, mid string) error {
	dsm := PSMessage{
		Service:    "healthcheck",
		MessageID:  uuid.NewString(),
		ResponseTo: mid,
		Content:    fmt.Sprintf("%s|%s|ok", s, v),
	}
	return a.Publish(dsm)
}
