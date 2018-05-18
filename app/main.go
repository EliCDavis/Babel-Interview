package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"

	_ "github.com/go-sql-driver/mysql"
)

const rateLimitWindow = 15 * 60

const directMessageRateLimit = 180

var botCommands = []BotCommand{
	changeCommand,
	receiptCommand,
	totalCommand,
	clearCommand,
}

func initializeDatabase() (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(db)/%s?parseTime=true", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_DATABASE")))
	if err != nil {
		return nil, err
	}

	sql, err := ioutil.ReadFile("./init.sql")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(string(sql))
	if err != nil {
		return nil, err
	}

	return db, nil
}

func scanForCommands(babelBot *BabelBot) {
	log.Println("Bot watching for messages!")

	var err error

	updateRate := time.Second * (rateLimitWindow / directMessageRateLimit)

	for {
		// Sleep first, so if we panic and continue to the next loop, we still end up waiting
		time.Sleep(updateRate)

		log.Printf("tick..\n")

		err = babelBot.RespondToNewMessages()

		if err != nil {
			log.Printf("Error Responding To Messages: %s\n", err.Error())
		}
	}
}

func main() {

	log.Println("Starting bot..")

	// Connect to our database..
	db, err := initializeDatabase()
	for err != nil {
		log.Printf("Error starting database:\n\t%s\n\ttrying again in 1 second..", err.Error())
		time.Sleep(time.Second)
		db, err = initializeDatabase()
	}
	log.Println("Initialized Database..")

	// Create our twitter client..
	config := oauth1.NewConfig(os.Getenv("TWITTER_CONSUMER_KEY"), os.Getenv("TWITTER_CONSUMER_SECRET"))
	token := oauth1.NewToken(os.Getenv("TWITTER_ACCESS_TOKEN"), os.Getenv("TWITTER_ACCESS_SECRET"))
	twitterClient := twitter.NewClient(config.Client(oauth1.NoContext, token))
	log.Println("Created twitter client..")

	// Start our bot..
	scanForCommands(NewBabelBot(botCommands, twitterClient, db))
}
