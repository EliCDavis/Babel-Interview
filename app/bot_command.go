package main

import (
	"database/sql"

	"github.com/dghubble/go-twitter/twitter"
)

type BotCommand struct {
	Command     string
	Execute     func(contents string, message twitter.DirectMessage, database *sql.DB) string
	Name        string
	Example     string
	Description string
}
