package main

import (
	"database/sql"

	"github.com/dghubble/go-twitter/twitter"
)

var clearCommand = BotCommand{
	Name:        "Clear History",
	Command:     "clear",
	Description: "Delete all recorded change transactions",
	Example:     "/clear",
	Execute: func(contents string, message twitter.DirectMessage, database *sql.DB) string {
		_, err := database.Query(
			"DELETE FROM receipt WHERE twitterUserId=?",
			message.SenderID,
		)

		if err != nil {
			return "Error deleting transactions"
		}

		return "Cleared all transactions"
	},
}
