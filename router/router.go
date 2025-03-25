package router

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/naman1402/distributed-chat-app/config"
	"github.com/naman1402/distributed-chat-app/controller"
	"github.com/naman1402/distributed-chat-app/database"
)

func Start() {

	router := gin.Default()

	// Add logging middleware
	router.Use(gin.Logger())

	database.SetupConnection()

	/**
	sets redis client, establishing connection to a redis db
	starts new goroutine, PubSub() will run concurrently with main goroutine -> subscribe to redis client and continously handles publishing messages to a Pub/Sub system in a separate goroutine
	start another goroutine, Send() -> gets msg from broadcast channel and processes it, message -> process -> get receiver -> receiver conn -> send private OR group message
	config.PubSub() receiving data and config.Send() processing and distributing it

	*/
	config.NPool()     // create the redis.Client
	go config.PubSub() // receive message from pub sub and adds to broadcast channel
	go config.Send()   // gets message from broadcast channel, processes it and further sends it

	// Add health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	router.GET("/", home)
	router.POST("/create", controller.CreateRoom)
	router.POST("/join", controller.JoinRoom)
	router.POST("/signin", controller.CreateUser)
	router.POST("/login", controller.LoginUser)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is required")
	}

	log.Printf("Server starting on port: %s", port)
	if err := router.Run("0.0.0.0:" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "server started"})
}
