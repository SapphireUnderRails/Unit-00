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
		return
	}

	// Register a message create and channel create callback function.
	session.AddHandler(messageCreate)
	session.AddHandler(channelCreate)

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

	// Ignore all messages with the discriminator #0000.
	if message.Author.Discriminator == "0000" {
		return
	}

	// Retrieves all the information of the channel the message was created in.
	channel, err := session.Channel(message.ChannelID)
	if err != nil {
		fmt.Println("Could not retrieve channel details: ", err)
	}

	//Formatted output.
	fmt.Printf("----------\nCreated Message Details:\nChannel: %s\nChannel ID: %s\nAuthor: %s#%s\nAuthor ID: %s\nContent: %s\nTime: %d\n",
		channel.Name,
		channel.ID,
		message.Author.Username,
		message.Author.Discriminator,
		message.Author.ID,
		message.Content,
		message.Timestamp.Unix())

	// Retrieves all the webhooks attached to the channel.
	webhooks, err := session.ChannelWebhooks(channel.ID)
	if err != nil {
		fmt.Println("Could not retrieve channel webhooks: ", err)
		return
	}

	// Snagging the Webhook ID and Webhook Token to mimic the user.
	var webhook_id, webhook_token string
	for _, webhooks := range webhooks {
		if webhooks.Name == channel.Name {
			webhook_id = webhooks.ID
			webhook_token = webhooks.Token
		}
	}

	// Setting the parameters for the webhook that will mimic the user.
	params := discordgo.WebhookParams{}
	params.Content = message.Content
	params.Username = message.Author.Username
	params.AvatarURL = message.Author.AvatarURL(message.Author.Avatar)

	// Executing the webhook that will mimic the user.
	session.WebhookExecute(webhook_id, webhook_token, true, &params)
}

// This function will be called every time a new channel is made.
func channelCreate(session *discordgo.Session, channel *discordgo.ChannelCreate) {

	// Details of the newly created channel.
	fmt.Printf("----------\nCreated Channel Details:\nChannel ID: %s\nGuild ID: %s\nChannel Name: %s\n",
		channel.ID,
		channel.GuildID,
		channel.Name)
	switch channel.Type {
	case 0:
		fmt.Printf("Channel Topic: %s\nChannel Type: Text\n", channel.Topic)

		// Want to create a new web hook in the text channel that was just created.
		webhook, err := session.WebhookCreate(channel.ID, channel.Name, "")
		if err != nil {
			fmt.Println("Webhook Creation Error ", err)
		}
		fmt.Printf("Webhook successfully created with the name '%s'\n", webhook.Name)
	case 2:
		fmt.Println("Channel Type: Voice")
	case 4:
		fmt.Println("Channel Type: Category")
	default:
		fmt.Println("Channel Type: Other")
	}

}
