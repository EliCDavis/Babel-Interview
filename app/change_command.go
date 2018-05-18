package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/dghubble/go-twitter/twitter"
)

var changeCommand = BotCommand{
	Name:        "Change",
	Command:     "change",
	Description: "Computes the fewest number of American coins that represents the given monetary value. Returns a receipt number that can be used to retrieve this transaction",
	Example:     "/change 5.27",
	Execute: func(contents string, message twitter.DirectMessage, database *sql.DB) string {

		amount, err := strconv.ParseFloat(contents, 64)

		if err != nil {
			log.Printf("Error Parsing Monitary Value: %s\n", err.Error())
			return "Error interpreting monitary value, please enter valid input\nex: /change 5.27"
		}

		_, err = database.Exec(
			"INSERT INTO receipt (messageId, amount, twitterUserId) VALUES (?, ?, ?)",
			message.IDStr,
			amount,
			message.SenderID,
		)

		if err != nil {
			log.Printf("Error saving change: %s\n", err.Error())
			return "Error saving you're transaction, please try again at another time."
		}

		change := NewChange(amount)
		return fmt.Sprintf("Change:\n%v\n\nReceipt:\n%s", change, message.IDStr)
	},
}
