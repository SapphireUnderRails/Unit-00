package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Making a struct to hold the token data.
type Token struct {
	Token string
}

func main() {

	// Retrieve the token from token.json file.
	fileContent, err := os.ReadFile("token.json")
	if err != nil {
		log.Println("Could not open token file.")
		log.Fatal(err)
	}

	// Unmarshal the token from the file contnet to grab the token.
	var token Token
	json.Unmarshal(fileContent, &token)

	// Create a new Discord session using the provided bot token.
	session, err := discordgo.New("Bot " + token.Token)
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
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	session.Close()
}
