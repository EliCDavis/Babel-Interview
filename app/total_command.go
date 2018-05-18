package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/dghubble/go-twitter/twitter"
)

var totalCommand = BotCommand{
	Name:        "Total Change",
	Command:     "total",
	Description: "Determine how much change you've recieved over the use of the bot",
	Example:     "/total",
	Execute: func(contents string, message twitter.DirectMessage, database *sql.DB) string {
		result, err := database.Query(
			"SELECT amount FROM receipt WHERE twitterUserId=?",
			message.SenderID,
		)

		if err != nil {
			return "Error looking up receipt, please try again later"
		}

		var amount float64
		var total float64
		totalChange := NewChange(0)
		for result.Next() {
			err = result.Scan(&amount)
			if err != nil {
				log.Printf("Error retrieving: %s\n", err.Error())
				return "Error retrieving receipt"
			}
			total += amount
			totalChange.Add(*NewChange(amount))
		}

		return fmt.Sprintf("Total Change Given: $%.2f\nCoins: %v", total, totalChange)
	},
}
