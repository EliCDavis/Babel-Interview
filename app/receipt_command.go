package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/dghubble/go-twitter/twitter"
)

var receiptCommand = BotCommand{
	Name:        "Receipt",
	Command:     "receipt",
	Description: "Retrieve details on a change transaction",
	Example:     "/receipt 997274705244622852",
	Execute: func(contents string, message twitter.DirectMessage, database *sql.DB) string {
		result, err := database.Query(
			"SELECT amount, date FROM receipt WHERE messageId=? AND twitterUserId=?",
			contents,
			message.SenderID,
		)

		if err != nil {
			return "Error looking up receipt, please try again later"
		}

		if result.Next() {
			var amount float64
			var date time.Time
			err = result.Scan(&amount, &date)
			if err != nil {
				log.Printf("Error retrieving: %s\n", err.Error())
				return "Error retrieving receipt"
			}

			change := NewChange(amount)
			return fmt.Sprintf("%s\n\nInput: $%.2f\nOutput: %v", date.Format("15:04:05 Monday January 2, 2006"), amount, change)
		}

		return fmt.Sprintf("Receipt '%s' not found", contents)
	},
}
