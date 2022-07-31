package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

type Config struct {
	Username string
	Password string
	Database string
}

func message_database() sql.DB {

	// Retrieve the configuration information from sql_data.json file.
	fileContent, err := os.ReadFile("sql_data.json")
	if err != nil {
		log.Println("Could not open sql_data file.")
		log.Fatal(err)
	}

	// Unmarshal the token from the file contnet to grab the token.
	var config Config
	json.Unmarshal(fileContent, &config)

	// Set up the configuration for the database connection.
	configuration := mysql.Config{
		User:   config.Username,
		Passwd: config.Password,
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: config.Database,
	}

	// Open a connection to the discord_messages database.
	db, err := sql.Open("mysql", configuration.FormatDSN())
	if err != nil {
		log.Println(err)
	}

	return *db
}
