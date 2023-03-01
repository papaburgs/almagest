package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/papaburgs/almagest/pkg/config"
	rt "github.com/papaburgs/almagest/pkg/redistools"
	redis "github.com/redis/go-redis/v9"
)

type discordHelper struct {
	Session  *discordgo.Session
	Channels map[string]string
}

func (d *discordHelper) Dispatch(channel, message string) error {
	var (
		cID string
		ok  bool
		err error
	)

	cID, ok = d.Channels[channel]
	if !ok {
		return fmt.Errorf("channel not found: %s", channel)
	}

	_, err = d.Session.ChannelMessageSend(cID, message)
	return err
}

// Don't like these globals, but its just easier
var arc *rt.AlmagestRedisClient
var dh *discordHelper

// Start starts the bot and populates some config
func Start(done chan bool) {
	// define the redis client
	arc = rt.New()

	// setup the discord bot
	var goBot *discordgo.Session
	c := config.Config
	goBot, err := discordgo.New("Bot " + c.Token)

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
		log.Fatal(err)
	}
	dh = &discordHelper{
		Session:  goBot,
		Channels: make(map[string]string),
	}
	// build channel list
	for _, g := range goBot.State.Ready.Guilds {
		fmt.Println(g.ID)
		if err != nil {
			log.Fatal(err)
		}
		channels, err := goBot.GuildChannels(g.ID)
		if err != nil {
			log.Fatal(err)
		}
		for _, c := range channels {
			if c.Type == discordgo.ChannelTypeGuildText {
				log.Printf("adding channel %s [%s]", c.Name, c.ID)
				dh.Channels[c.Name] = c.ID
			}
		}
	}

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
		err = json.Unmarshal([]byte(msg.Payload), &psm)
		if err != nil {
			log.Println("Could not decode message ", msg.Payload)
			continue
		}
		if psm.Service == "discord" {
			log.Printf("Service: %s, send to %s, with message %s\n",
				psm.Service,
				psm.Channel,
				psm.Content,
			)
			err = dh.Dispatch(psm.Channel, psm.Content)
			if err != nil {
				return
			}
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
		dsm := rt.PSMessage{
			Service:   "uptime",
			MessageID: uuid.New().String(),
		}
		log.Print("its a message for me, someone is looking for uptime")
		err := arc.Publish(dsm)
		if err != nil {
			log.Printf("error posting to redis: %s", err)
		}
		log.Print("published to uptime service")
		_, _ = s.ChannelMessageSend(m.ChannelID, "Uptime request recieved, please hold")
	}
}
