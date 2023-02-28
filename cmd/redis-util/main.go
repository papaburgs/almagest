package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	rt "github.com/papaburgs/almagest/pkg/redistools"
	redis "github.com/redis/go-redis/v9"
)

var arc *rt.AlmagestRedisClient

func main() {
	log.Println("starting")
	arc = rt.New()

	go redisListener()

	time.Sleep(3 * time.Second)

	m1 := rt.PSMessage{
		Service: "discord",
		Channel: "default",
		Content: "first message",
	}
	arc.Publish(m1)
	time.Sleep(3 * time.Second)
	m2 := rt.PSMessage{
		Service: "discord",
		Channel: "default",
		Content: "second message",
	}
	arc.Publish(m2)
	time.Sleep(2 * time.Second)
}

func redisListener() {
	var (
		msg *redis.Message
		psm rt.PSMessage
	)
	c := arc.Subscribe()
	for {
		msg = <-c
		log.Println("got a command")
		err := json.Unmarshal([]byte(msg.Payload), &psm)
		if err != nil {
			log.Println("Could not decode message")
		}
		fmt.Printf("send to %s, with message %s\n",
			psm.Channel,
			psm.Content,
		)
	}

}

// type discordMessage struct {
// 	Content string
// 	Channel string
// }

// func messageFromRedis(msg *redis.Message) discordMessage {
// 	dm := discordMessage{}
// 	channelSlice := strings.Split(msg.Channel, "|")
// 	if len(channelSlice) < 3 {
// 		log.Println("Not enough chunks")
// 		dm.Channel = ""
// 		dm.Content = "not enough chunks"
// 		return dm
// 	}
// 	dm.Channel = channelSlice[2]
// 	dm.Content = msg.Payload
// 	return dm
// }
