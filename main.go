package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/go-sql-driver/mysql"
)

// Making a struct to hold the token data.
type Token struct {
	Token string
}

//Making a struct to hold the MySQL server logon parameters.
type Parameters struct {
	Username string
	Password string
	Database string
}

//Global variable to hold database connection, because why not?
var db *sql.DB

func main() {

	// Retrieve the parameters from sql_data.json file.
	sql_parameters, err := os.ReadFile("sql_data.json")
	if err != nil {
		log.Println("Could not open sql_data file.")
		log.Fatal(err)
	}

	// Unmarshal the parameters from the file contnet to grab the logon information.
	var parameters Parameters
	json.Unmarshal(sql_parameters, &parameters)

	// Set up the parameters for the database connection.
	configuration := mysql.Config{
		User:   parameters.Username,
		Passwd: parameters.Password,
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: parameters.Database,
	}

	// Open a connection to the discord_messages database.
	db, err = sql.Open("mysql", configuration.FormatDSN())
	if err != nil {
		log.Println(err)
	}

	// Retrieve the token from token.json file.
	discord_token, err := os.ReadFile("token.json")
	if err != nil {
		log.Println("Could not open token file.")
		log.Fatal(err)
	}

	// Unmarshal the token from the file contnet to grab the token.
	var token Token
	json.Unmarshal(discord_token, &token)

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

// This function will be called every time a new message created.
func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {

	// Ignore all messages that were created by the bot itself.
	if message.Author.ID == session.State.User.ID {
		return
	}

	// Ignore all messages with the discriminator #0000 (Webhooks).
	if message.WebhookID != "" {
		return
	}

	// Filter out all URLs in the message.
	re, err := regexp.Compile(`([\w+]+\:\/\/)?([\w\d-]+\.)*[\w-]+[\.\:]\w+([\/\?\=\&\#\.]?[\w-]+)*\/?`)
	if err != nil {
		log.Println(err)
	}
	message_content := re.ReplaceAllString(message.Content, "")

	// Ultimately ignore all messages with no content in them.
	if message_content == "" {
		return
	}

	// Create a table (if it doesn't already exist) in the database specific to the user to store the message in.
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS user_%s
	(id BIGINT unsigned AUTO_INCREMENT PRIMARY KEY,
	message_id BIGINT NOT NULL,
	content LONGTEXT NOT NULL);`, message.Author.ID)
	result, err := db.Exec(query)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(result)

	// Store the message data into the table.
	query = fmt.Sprintf(`INSERT INTO user_%s(message_id, content) VALUES("%s", "%s");`,
		message.Author.ID, message.ID, message_content)
	result, err = db.Exec(query)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(result)

	// Grab all the user's messages from the database.
	query = fmt.Sprintf(`SELECT content FROM user_%s;`, message.Author.ID)
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
	}

	// Add all the snagged messages to one giant string.
	content := ""
	for rows.Next() {
		var message string
		rows.Scan(&message)
		content = content + message + " "
	}

	// Create the Markov chain with an order of 3 to mimic the user.
	chain := NewChain(2)

	// Feed in the giant string of messages for training.
	chain.Build(strings.NewReader(content))

	// Generate the chain.
	content_chain := chain.Generate(64)
	if err != nil {
		log.Println(err)
	}
	log.Println(content_chain)

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

	// Setting the parameters for the webhook that will mimic the user.
	params := discordgo.WebhookParams{}
	params.Content = content_chain
	params.Username = message.Author.Username
	params.AvatarURL = message.Author.AvatarURL(message.Author.Avatar)

	// Executing the webhook that will mimic the user.
	_, err = session.WebhookExecute(webhook_id, webhook_token, true, &params)
	if err != nil {
		log.Println(err)
		return
	}
}
