package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/naman1402/distributed-chat-app/config"
	"github.com/naman1402/distributed-chat-app/controller"
	"github.com/naman1402/distributed-chat-app/database"
)

func Start() {

	router := gin.Default()
	database.SetupConnection()

	/**
	sets redis client, establishing connection to a redis db
	starts new goroutine, PubSub() will run concurrently with main goroutine -> subscribe to redis client and continously handles publishing messages to a Pub/Sub system in a separate goroutine
	start another goroutine, Send() -> gets msg from broadcast channel and processes it, message -> process -> get receiver -> receiver conn -> send private OR group message
	config.PubSub() receiving data and config.Send() processing and distributing it

	*/
	config.NPool()
	go config.PubSub()
	go config.Send()

	router.GET("/", home)
	router.POST("/create", controller.CreateRoom)
	router.POST("/join", controller.JoinRoom)
	router.POST("/signin", controller.CreateUser)
	router.POST("/login", controller.LoginUser)
	router.Run(":" + "3000")
}

func home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "server started"})
}
