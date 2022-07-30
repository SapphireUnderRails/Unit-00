package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {

	// Create a new Discord session using the provided bot token.
	session, err := discordgo.New("Bot " + "MTAwMTA3Njk0MjQ3NDMyMjAxMQ.GtnF2E.-7LU7kCBWdxLWxDrQ2Zk-aIIPD0hYaDcd39ECU")
	if err != nil {
		log.Fatal("Error creating Discord session,", err)
	}

	// Register a message create and channel create callback function.
	session.AddHandler(messageCreate)
	session.AddHandler(channelCreate)

	// Identify that we want all intents.
	session.Identify.Intents = discordgo.IntentsAll

	// Now we open a websocket connection to Discord and begin listening.
	err = session.Open()
	if err != nil {
		log.Fatal("Error opening Discord session,", err)
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	session.Close()
}
