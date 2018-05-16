package main

import (
	"database/sql"

	"github.com/dghubble/go-twitter/twitter"
)

type BabelBot struct {
	twitterClient *twitter.Client
	database      *sql.DB
}

func (bot BabelBot) CreateChangeWithRecipt(amount float64) {
	change := NewChange(amount)
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
