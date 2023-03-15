package main

import (
	"strings"

	"github.com/charmbracelet/log"
	redis "github.com/redis/go-redis/v9"

	_ "embed"

	rt "github.com/papaburgs/almagest/pkg/redistools"
)

var arc *rt.AlmagestRedisClient

//go:embed gitc.txt
var gitCommit string

func main() {

	gitCommit = strings.TrimSpace(gitCommit)
	arc = rt.New()
	log.SetLevel(log.DebugLevel)

	var (
		msg   *redis.Message
		psm   rt.PSMessage
		class rt.RedisMessageClass
		err   error
	)

	c := arc.Subscribe()
	for {
		msg = <-c
		psm, class, err = rt.ClassifyMessage(msg)
		if err != nil {
			continue
		}

		if class == rt.Watchdog {
			log.Debug("Picked up watchdog message", "service", psm.Service)
		}
		if class == rt.HealthCheckRequest {
			log.Debug("replying to health check request")
			arc.PostStatus("watchdog", gitCommit, psm.MessageID)
		}
	}
}
