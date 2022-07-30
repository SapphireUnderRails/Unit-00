package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {

	// Create a new Discord session using the provided bot token.
	session, err := discordgo.New("Bot " + "MTAwMTA3Njk0MjQ3NDMyMjAxMQ.GtnF2E.-7LU7kCBWdxLWxDrQ2Zk-aIIPD0hYaDcd39ECU")
	if err != nil {
		fmt.Println("Error creating Discord session,", err)
	}

	// Register a message create callback function.
	session.AddHandler(messageCreate)

	// Identify that we want all intents.
	session.Identify.Intents = discordgo.IntentsAll

	// Now we open a websocket connection to Discord and begin listening.
	err = session.Open()
	if err != nil {
		fmt.Println("Error opening Discord session,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	session.Close()
}

// This function will be called every time a new message created.
func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {

	// Ignore all messages that were created by the bot itself.
	if message.Author.ID == session.State.User.ID {
		return
	}

	//Retrieves all the information of the channel the message was created in.
	channel, err := session.Channel(message.ChannelID)
	if err != nil {
		fmt.Println("Could not retrieve channel details: ", err)
	}

	//Formatted output.
	fmt.Printf("Channel: %s\nChannel ID: %s\nAuthor: %s#%s\nAuthor ID: %s\nContent: %s\nTime: %d",
		channel.Name,
		channel.ID,
		message.Author.Username,
		message.Author.Discriminator,
		message.Author.ID,
		message.Content,
		message.Timestamp.Unix())
}
