package main

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

// This function will be called every time a new channel is made.
func channelCreate(session *discordgo.Session, channel *discordgo.ChannelCreate) {

	// Create a mimicry webhook in the newly created channel.
	if channel.Type == 0 {
		webhook, err := session.WebhookCreate(channel.ID, fmt.Sprintf("%s-mimic", channel.Name), "")
		if err != nil {
			fmt.Println("Webhook Creation Error ", err)
		}
		log.Printf("Webhook successfully created with the name '%s'\n", webhook.Name)
	}

}
