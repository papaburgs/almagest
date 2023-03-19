package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
	redis "github.com/redis/go-redis/v9"

	_ "embed"

	"github.com/papaburgs/almagest/pkg/cfg"
	rt "github.com/papaburgs/almagest/pkg/redistools"
)

var arc *rt.AlmagestRedisClient

//go:embed gitc.txt
var gitCommit string

func main() {
	gitCommit = strings.TrimSpace(gitCommit)
	arc = rt.New()
	newLevel := cfg.GetParam(arc, cfg.LogLevel)
	switch newLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	}
	log.Info("LogLevel updated", "level", newLevel)

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

		switch class {
		case rt.Watchdog:
			log.Debug("Picked up watchdog message", "service", psm.Service)
		case rt.HealthCheckRequest:
			log.Debug("replying to health check request")
			arc.PostStatus("watchdog", gitCommit, psm.MessageID)
		case rt.ControlUpdateLogging:
			switch psm.Content {
			case "debug":
				log.SetLevel(log.DebugLevel)
			case "info":
				log.SetLevel(log.InfoLevel)
			case "warn":
				log.SetLevel(log.WarnLevel)
			case "error":
				log.SetLevel(log.ErrorLevel)
			}
			log.Info("LogLevel updated", "level", psm.Content)
		default:
			log.Info("Got message for someone else", "message", fmt.Sprintf("%v", psm))
		}
	}
}
