package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// This function will be called every time a new message created.
func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {

	// Ignore all messages that were created by the bot itself.
	if message.Author.ID == session.State.User.ID {
		return
	}

	// Ignore all messages with the discriminator #0000 (Webhooks).
	if message.Author.Discriminator == "0000" {
		return
	}

	// Retrieves all the information of the channel the message was created in.
	channel, err := session.Channel(message.ChannelID)
	if err != nil {
		fmt.Println("Could not retrieve channel details: ", err)
		return
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

	// // Retrieves all the webhooks attached to the channel.
	// webhooks, err := session.ChannelWebhooks(channel.ID)
	// if err != nil {
	// 	fmt.Println("Could not retrieve channel webhooks: ", err)
	// 	return
	// }

	// // Snagging the Webhook ID and Webhook Token to mimic the user.
	// var webhook_id, webhook_token string
	// for _, webhooks := range webhooks {
	// 	if webhooks.Name == channel.Name {
	// 		webhook_id = webhooks.ID
	// 		webhook_token = webhooks.Token
	// 	}
	// }

	// // Setting the parameters for the webhook that will mimic the user.
	// params := discordgo.WebhookParams{}
	// params.Content = message.Content
	// params.Username = message.Author.Username
	// params.AvatarURL = message.Author.AvatarURL(message.Author.Avatar)

	// // Executing the webhook that will mimic the user.
	// webhook_message, err := session.WebhookExecute(webhook_id, webhook_token, true, &params)
	// if err != nil {
	// 	fmt.Println("Could not execute webhook: ", err)
	// 	return
	// }
	// fmt.Println(webhook_message)
}
