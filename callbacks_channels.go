package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

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
