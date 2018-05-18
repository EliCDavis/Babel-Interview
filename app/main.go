package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/gin-gonic/gin"

	_ "github.com/go-sql-driver/mysql"
)

const twitterDateFormat = "Mon Jan 2 15:04:05 -0700 2006"

const rateLimitWindow = 15 * 60

const directMessageRateLimit = 180

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

// Create our routes
func initRoutes(router *gin.Engine) {

	router.LoadHTMLGlob("templates/*.tmpl")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{})
	})

	router.StaticFS("/media", http.Dir("media"))
}

func scanForChangeRequests(babelBot *BabelBot) {

	mostRecentMessage := time.Now()

	for {
		// Sleep first, so if we panic and continue to the next loop, we still end up waiting
		time.Sleep(time.Second * (rateLimitWindow / directMessageRateLimit))

		log.Printf("tick %s\n", time.Now().Format("15:04:05"))

		messages, err := babelBot.GetMessages()

		if err != nil {
			log.Printf("Error getting messages: %s\n", err.Error())
		} else {
			for _, msg := range messages {
				messageTime, _ := time.Parse(twitterDateFormat, msg.CreatedAt)
				if messageTime.After(mostRecentMessage) {
					mostRecentMessage = messageTime
					log.Println(msg.Text)
					if messageContainsCommand(msg.Text, "change") {
						f, err := strconv.ParseFloat(getContentFromCommand(msg.Text, "change"), 64)
						if err != nil {
							log.Println("Could not parse number")
						} else {
							babelBot.CreateChangeWithRecipt(f, msg)
						}
					}

					if messageContainsCommand(msg.Text, "recipt") {
						babelBot.RetrieveRecipt(getContentFromCommand(msg.Text, "recipt"), msg)
					}
				}
			}
		}
	}

}

func main() {

	// Can't run a server without a port
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable was not set")
		return
	}
	log.Printf("Starting bot using port %s\n", port)

	db, err := initializeDatabase()
	for err != nil {
		log.Printf("Error starting database:\n\t%s\n\ttrying again in 1 second..", err.Error())
		time.Sleep(time.Second * 10)
		db, err = initializeDatabase()
	}
	log.Println("Initialized Database")

	config := oauth1.NewConfig(os.Getenv("TWITTER_CONSUMER_KEY"), os.Getenv("TWITTER_CONSUMER_SECRET"))
	token := oauth1.NewToken(os.Getenv("TWITTER_ACCESS_TOKEN"), os.Getenv("TWITTER_ACCESS_SECRET"))
	twitterClient := twitter.NewClient(config.Client(oauth1.NoContext, token))

	bot := NewBabelBot(twitterClient, db)

	go scanForChangeRequests(bot)

	// Create our engine
	r := gin.New()

	// Logging middleware
	r.Use(gin.Logger())

	// Recover from errors and return 500
	r.Use(gin.Recovery())

	initRoutes(r)
	r.Run(":" + port)

}
