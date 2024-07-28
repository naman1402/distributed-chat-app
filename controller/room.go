package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/naman1402/distributed-chat-app/database"
	"github.com/naman1402/distributed-chat-app/model"
	"github.com/segmentio/ksuid"
)

func CreateRoom(c *gin.Context) {
	newRoom := model.Room{}
	if err := c.ShouldBindBodyWithJSON(&newRoom); err != nil {
		fmt.Println(err)
	}
	Id := ksuid.New()
	query := `INSERT INTO room (id, room_name) VALUES (?, ?)`
	database.ExecuteQuery(query, Id, newRoom.Name)
	c.JSON(http.StatusOK, gin.H{"message": "done"})
}

func JoinRoom(c *gin.Context) {
	joiningRoom := model.Room{}
	if err := c.ShouldBindBodyWithJSON(&joiningRoom); err != nil {
		fmt.Println(err)
	}
	query := `INSERT INTO room_members(room_name,username)VALUES(?,?)`
	database.ExecuteQuery(query, joiningRoom.Name, joiningRoom.User)
	c.JSON(http.StatusOK, gin.H{"message": "room joined"})
}
