package main

import (
	"fmt"
	"log"
	"strings"

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

	// Ignore all messages with no content in them.
	if message.Content == "" {
		return
	}

	// Create a database connection.
	message_database := message_database()

	// Create a table (if it doesn't already exist) in the database specific to the user to store the message in.
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS user_%s
	(id BIGINT unsigned AUTO_INCREMENT PRIMARY KEY,
	message_id BIGINT NOT NULL,
	channel_id BIGINT NOT NULL,
	content LONGTEXT NOT NULL);`, message.Author.ID)
	result, err := message_database.Exec(query)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(result)

	// Store the message data into the table.
	query = fmt.Sprintf(`INSERT INTO user_%s(message_id, channel_id, content) VALUES("%s", "%s", "%s");`,
		message.Author.ID, message.ID, message.ChannelID, message.Content)
	result, err = message_database.Exec(query)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(result)

	// Retrieves all the information of the channel the message was created in.
	channel, err := session.Channel(message.ChannelID)
	if err != nil {
		fmt.Println("Could not retrieve channel details: ", err)
		return
	}

	// Retrieves all the webhooks attached to the channel.
	webhooks, err := session.ChannelWebhooks(channel.ID)
	if err != nil {
		fmt.Println("Could not retrieve channel webhooks: ", err)
		return
	}

	// Snagging the Webhook ID and Webhook Token to mimic the user.
	var webhook_id, webhook_token string
	for _, webhooks := range webhooks {
		if webhooks.Name == fmt.Sprintf("%s-mimic", channel.Name) {
			webhook_id = webhooks.ID
			webhook_token = webhooks.Token
		}
	}

	// Grab all the user's messages from the database.
	query = fmt.Sprintf(`SELECT content FROM user_%s;`, message.Author.ID)
	rows, err := message_database.Query(query)
	if err != nil {
		log.Println(err)
		return
	}

	// Add all the snagged messages to one giant string.
	content := ""
	for rows.Next() {
		var message string
		rows.Scan(&message)
		content = content + message + " "
	}

	// Create the Markov chain with an order of 3 to mimic the user.
	chain := NewChain(1)

	// Feed in the giant string of messages for training.
	chain.Build(strings.NewReader(content))

	// Generate the chain.
	content_chain := chain.Generate(16)
	if err != nil {
		log.Println(err)
		return
	}

	// Setting the parameters for the webhook that will mimic the user.
	params := discordgo.WebhookParams{}
	params.Content = content_chain
	params.Username = message.Author.Username
	params.AvatarURL = message.Author.AvatarURL(message.Author.Avatar)

	// Executing the webhook that will mimic the user.
	webhook_message, err := session.WebhookExecute(webhook_id, webhook_token, true, &params)
	if err != nil {
		fmt.Println("Could not execute webhook: ", err)
		return
	}
	log.Println(webhook_message)
}
