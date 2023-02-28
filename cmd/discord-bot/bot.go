package main

import (
	"encoding/json"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/papaburgs/almagest/pkg/config"
	rt "github.com/papaburgs/almagest/pkg/redistools"
	redis "github.com/redis/go-redis/v9"
)

var arc *rt.AlmagestRedisClient

// Start starts the bot and populates some config
func Start(done chan bool) {
	// define the redis client
	arc = rt.New()

	// setup the discord bot
	var goBot *discordgo.Session
	c := config.Config
	goBot, err := discordgo.New("Bot " + c.Token)
	goBot.StateEnabled = true

	if err != nil {
		log.Println(err.Error())
		return
	}

	u, err := goBot.User("@me")
	if err != nil {
		log.Println(err.Error())
		return
	}
	c.BotID = u.ID

	goBot.AddHandler(messageHandler)

	err = goBot.Open()

	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println("Bot is running!")
	// redisListener is blocking
	redisListener()
	done <- true
}

func redisListener() {
	var (
		msg *redis.Message
		psm rt.PSMessage
		err error
	)

	c := arc.Subscribe()
	for {
		msg = <-c
		log.Println("got a command")
		err = json.Unmarshal([]byte(msg.Payload), &psm)
		if err != nil {
			log.Println("Could not decode message")
		}
		if psm.Service == "discord" {
			arc.Publish(psm)
			log.Printf("Service: %s, send to %s, with message %s\n",
				psm.Service,
				psm.Channel,
				psm.Content,
			)
		}
	}
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	c := config.Config
	log.Print("running message handler")
	if m.Author.ID == c.BotID {
		log.Print("message was from me, ignoring")
		return
	}

	if m.Content == "ping" {
		log.Print("its a message for me, sending pong")
		log.Println(m.ChannelID)
		_, _ = s.ChannelMessageSend(m.ChannelID, "pong")
	}

	if m.Content == "!uptime" {
		dsm := rt.PSMessage{Service: "uptime"}
		log.Print("its a message for me, someone is looking for uptime")
		err := arc.Publish(dsm)
		if err != nil {
			log.Printf("error posting to redis: %s", err)
		}
		log.Print("published to uptime service")
		_, _ = s.ChannelMessageSend(m.ChannelID, "Uptime request recieved, please hold")
	}
}
