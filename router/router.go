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
