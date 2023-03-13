package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	_ "embed"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	rt "github.com/papaburgs/almagest/pkg/redistools"
	redis "github.com/redis/go-redis/v9"
)

type discordHelper struct {
	Session  *discordgo.Session
	Channels map[string]string
	BotID    string
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

//go:embed gitc.txt
var gitCommit string

func main() {
	log.SetLevel(log.DebugLevel)
	done := make(chan bool, 1)
	Start(done)

	// hold application up indefinatly
	<-done
}

// Start starts the bot
func Start(done chan bool) {
	// define the redis client
	arc = rt.New()

	// setup the discord bot
	var goBot *discordgo.Session
	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		log.Fatal("Need DISCORD_BOT_TOKEN to be set")
	}
	goBot, err := discordgo.New("Bot " + token)

	if err != nil {
		log.Error("Error starting discord", "error", err)
		return
	}

	// find the details on the bot, so we have the ID for later
	u, err := goBot.User("@me")
	if err != nil {
		log.Error("Error starting finding details", "error", err)
		return
	}

	goBot.AddHandler(messageHandler)

	err = goBot.Open()
	if err != nil {
		log.Fatal(err)
	}
	dh = &discordHelper{
		Session:  goBot,
		Channels: make(map[string]string),
		BotID:    u.ID,
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
				log.Debug("adding channel %s [%s]", c.Name, c.ID)
				dh.Channels[c.Name] = c.ID
			}
		}
	}

	if err != nil {
		log.Error(err.Error())
		return
	}

	log.Info("Bot is running!")
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
			log.Error("Could not decode message ", "payload", msg.Payload)
			continue
		}

		log.Debug(fmt.Sprintf("Saw a message: %#v", psm))
		if psm.Service == "discord" {
			log.Debug("Sending message", "service", psm.Service, "channel", psm.Channel, "content", psm.Content)
			err = dh.Dispatch(psm.Channel, psm.Content)
			if err != nil {
				return
			}
		}
		if psm.Service == "healthcheck" {
			log.Debug("Saw a healthcheck")
			log.Debug(fmt.Sprintf("%#v", psm))
			if psm.ResponseTo == "" {
				log.Debug("looks like a new one, sending response")

				dsm := rt.PSMessage{
					Service:    "healthcheck",
					MessageID:  uuid.NewString(),
					ResponseTo: psm.MessageID,
				}
				dsm.Content = fmt.Sprintf("%s|%s|ok", "api", strings.TrimSpace(gitCommit))
				arc.Publish(dsm)
			}
		}
	}
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Print("running message handler")
	if m.Author.ID == dh.BotID {
		log.Print("message was from me, ignoring")
		return
	}

	if m.Content == "ping" {
		log.Debug("its a message for me, sending pong")
		log.Debug(m.ChannelID)
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
			log.Error("error posting to redis: %s", "error", err)
		}
		log.Print("published to uptime service")
		_, _ = s.ChannelMessageSend(m.ChannelID, "Uptime request recieved, please hold")
	}
}
