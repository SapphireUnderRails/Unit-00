package main

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

func message_database() sql.DB {

	// Set up the configuration for the database connection.
	configuration := mysql.Config{
		User:   "root",
		Passwd: "KingArthur09052012?",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "discord_messages",
	}

	// Open a connection to the discord_messages database.
	db, err := sql.Open("mysql", configuration.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	return *db
}
