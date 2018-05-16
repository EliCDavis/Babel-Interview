package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/gin-gonic/gin"
)

import _ "github.com/go-sql-driver/mysql"

const twitterDateFormat = "Mon Jan 2 15:04:05 -0700 2006"

const rateLimitWindow = 15 * 60

const directMessageRateLimit = 450

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
		// Sleep first, so if we error out and continue to the next loop, we still end up waiting
		time.Sleep(time.Second * (rateLimitWindow / directMessageRateLimit))

		log.Printf("tick %s\n", time.Now().String())

		messages, err := babelBot.GetMessages()

		if err != nil {
			log.Printf("Error getting messages: %s\n", err.Error())
		} else {
			for _, msg := range messages {
				messageTime, _ := time.Parse(twitterDateFormat, msg.CreatedAt)
				if messageTime.After(mostRecentMessage) {
					mostRecentMessage = messageTime
					log.Println(msg.Text)
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

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/db", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD")))
	if err != nil {
		log.Printf("Error connecting to database: %s\n", err.Error())
		return
	}

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
