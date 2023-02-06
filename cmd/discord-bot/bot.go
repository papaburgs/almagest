package main

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/papaburgs/almagest/pkg/config"
)

// Start starts the bot and populates some config
func Start() {
	var goBot *discordgo.Session
	c := config.Config
	goBot, err := discordgo.New("Bot " + c.Token)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	u, err := goBot.User("@me")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	c.BotID = u.ID

	goBot.AddHandler(pingHandler)

	err = goBot.Open()

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	log.Println("Bot is running !")
}

func pingHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	c := config.Config
	log.Print("running ping handler")
	if m.Author.ID == c.BotID {
		log.Print("message was from me, ignoring")
		return
	}

	if m.Content == "ping" {
		log.Print("its a message for me, sending pong")
		_, _ = s.ChannelMessageSend(m.ChannelID, "pong")
	}
}
