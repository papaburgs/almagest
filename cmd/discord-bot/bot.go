package main

import (
	"fmt"
	"os"

	_ "embed"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	rt "github.com/papaburgs/almagest/pkg/redistools"
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
	go arc.PublishWatchdog("discord")

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
