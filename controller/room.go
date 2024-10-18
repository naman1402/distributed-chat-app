package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/naman1402/distributed-chat-app/database"
	"github.com/naman1402/distributed-chat-app/model"
	"github.com/segmentio/ksuid"
)

// context holds information about the upcoming HTTP request and provides method to handle the response
// storing the context into newRoom (JSON format) and inserting into cassandra session by calling execute query from database package
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

// get room from context and execute query to add member into room_name (room and user name from context -> joiningRoom)
func JoinRoom(c *gin.Context) {
	joiningRoom := model.Room{}
	if err := c.ShouldBindBodyWithJSON(&joiningRoom); err != nil {
		fmt.Println(err)
	}
	query := `INSERT INTO room_members(room_name,username)VALUES(?,?)`
	database.ExecuteQuery(query, joiningRoom.Name, joiningRoom.User)
	c.JSON(http.StatusOK, gin.H{"message": "room joined"})
}

// data is iterator containing room members (collected by calling query directly through cassandra session)
// in the for loop, we scan username from data and store it into members ([]string)
// return the array containing members
func GetMembersFromRoom(groupname string) []string {
	var username string
	members := []string{}
	query := `SELECT username FROM room_members WHERE room_name = ?`
	data := database.Connection.Session.Query(query, groupname).Iter()
	for data.Scan(&username) {
		members = append(members, username)
	}
	return members
}
