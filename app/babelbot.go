package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/dghubble/go-twitter/twitter"
)

type BabelBot struct {
	twitterClient *twitter.Client
	database      *sql.DB
}

func (bot BabelBot) CreateChangeWithRecipt(amount float64, message twitter.DirectMessage) (*Change, error) {
	change := NewChange(amount)

	_, err := bot.database.Exec(
		"INSERT INTO receipt (messageId, amount, twitterUserId) VALUES (?, ?, ?)",
		message.IDStr,
		amount,
		message.SenderID,
	)

	if err != nil {
		log.Printf("Error saving change: %s\n", err.Error())
		return nil, err
	}
	_, _, err = bot.twitterClient.DirectMessages.New(&twitter.DirectMessageNewParams{
		Text:   fmt.Sprintf("Change:\n%d Quater(s), %d Dimes, %d Nickle(s), %d Penny(s)\n\nRecipt:\n%s", change.Quarters, change.Dimes, change.Nickles, change.Pennies, message.IDStr),
		UserID: message.SenderID,
	})

	return change, nil
}

func (bot BabelBot) RetrieveRecipt(recipt string, message twitter.DirectMessage) {
	result, err := bot.database.Query(
		"SELECT amount, date FROM receipt WHERE messageId=? AND twitterUserId=?",
		recipt,
		message.SenderID,
	)

	if result.Next() {
		var amount float64
		var date time.Time
		err = result.Scan(&amount, &date)
		if err != nil {
			log.Printf("Error retrieving: %s\n", err.Error())
			return
		}

		log.Println(date)

		change := NewChange(amount)
		_, _, err = bot.twitterClient.DirectMessages.New(&twitter.DirectMessageNewParams{
			Text:   fmt.Sprintf("Change:\n%d Quater(s), %d Dimes, %d Nickle(s), %d Penny(s)\n\nRecipt:\n%s", change.Quarters, change.Dimes, change.Nickles, change.Pennies, recipt),
			UserID: message.SenderID,
		})
	}
}

func (bot BabelBot) GetMessages() ([]twitter.DirectMessage, error) {
	messages, _, err := bot.twitterClient.DirectMessages.Get(&twitter.DirectMessageGetParams{
		Count: 50, // 50 is twitter max
	})
	return messages, err
}

func NewBabelBot(twitterClient *twitter.Client, db *sql.DB) *BabelBot {
	bot := new(BabelBot)
	bot.twitterClient = twitterClient
	bot.database = db
	return bot
}
