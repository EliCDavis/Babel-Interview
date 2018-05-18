package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
)

func messageContainsCommand(message string, command string) bool {
	return len(strings.SplitAfter(message, fmt.Sprintf("/%s", command))) > 1
}

func getContentFromCommand(message string, command string) string {
	commands := strings.SplitAfter(message, fmt.Sprintf("/%s", command))
	if len(commands) > 1 {
		return strings.TrimSpace(commands[1])
	}
	return ""
}

const twitterDateFormat = "Mon Jan 2 15:04:05 -0700 2006"

type BabelBot struct {
	twitterClient  *twitter.Client
	database       *sql.DB
	commands       []BotCommand
	lastUpdateTime time.Time
}

func (bot *BabelBot) RespondToNewMessages() error {

	messages, err := bot.getMessages()

	if err != nil {
		return err
	}

	for _, msg := range messages {
		messageTime, _ := time.Parse(twitterDateFormat, msg.CreatedAt)

		// This is a new message
		if messageTime.After(bot.lastUpdateTime) {

			bot.lastUpdateTime = messageTime
			recognizedCommand := false

			// Attempt to find a command that is associated with the user text
			for _, cmd := range bot.commands {
				if messageContainsCommand(msg.Text, cmd.Command) {
					recognizedCommand = true

					_, _, err = bot.twitterClient.DirectMessages.New(&twitter.DirectMessageNewParams{
						Text:   cmd.Execute(getContentFromCommand(msg.Text, cmd.Command), msg, bot.database),
						UserID: msg.SenderID,
					})

					if err != nil {
						return err
					}
				}
			}

			// We didn't recognize the command, let the twitter user know
			if recognizedCommand == false {
				_, _, err = bot.twitterClient.DirectMessages.New(&twitter.DirectMessageNewParams{
					Text:   "unrecognized input, type /help for help",
					UserID: msg.SenderID,
				})
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (bot BabelBot) getMessages() ([]twitter.DirectMessage, error) {
	messages, _, err := bot.twitterClient.DirectMessages.Get(&twitter.DirectMessageGetParams{
		Count: 50, // 50 is twitter max
	})
	return messages, err
}

func NewBabelBot(commands []BotCommand, twitterClient *twitter.Client, db *sql.DB) *BabelBot {
	bot := new(BabelBot)
	bot.commands = append(commands, BotCommand{
		Command:     "help",
		Description: "description of all available commands",
		Example:     "/help",
		Name:        "Help",
		Execute: func(contents string, message twitter.DirectMessage, database *sql.DB) string {
			var buffer bytes.Buffer
			for _, cmd := range commands {
				buffer.WriteString(fmt.Sprintf("%s: (/%s)\n%s\nExample: %s\n\n", cmd.Name, cmd.Command, cmd.Description, cmd.Example))
			}
			return buffer.String()
		},
	})
	bot.twitterClient = twitterClient
	bot.database = db
	bot.lastUpdateTime = time.Now()
	return bot
}
