package main

import (
	"strings"

	"github.com/charmbracelet/log"
	rt "github.com/papaburgs/almagest/pkg/redistools"
	redis "github.com/redis/go-redis/v9"
)

func redisListener() {
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

		if class == rt.SendToDiscordRequest {
			log.Debug("Sending message", "service", psm.Service, "channel", psm.Channel, "content", psm.Content)
			err = dh.Dispatch(psm.Channel, psm.Content)
			if err != nil {
				log.Error("dispatch to discord error", "error", err)
				continue
			}
		}
		if class == rt.HealthCheckRequest {
			log.Debug("replying to health check request")
			arc.PostStatus("discord", strings.TrimSpace(gitCommit), psm.MessageID)
		}
		if class == rt.ControlUpdateLogging {
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
		}
	}
}
