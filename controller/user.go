package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/naman1402/distributed-chat-app/database"
	"github.com/naman1402/distributed-chat-app/model"
	"github.com/segmentio/ksuid"
)

func CreateUser(c *gin.Context) {
	user := &model.User{}
	if err := c.ShouldBindBodyWithJSON(&user); err != nil {
		fmt.Println(err)
	}
	Id := ksuid.New()
	query := `INSERT INTO users(id, username) VALUES(?, ?)`
	database.ExecuteQuery(query, Id.String(), user.Username)
	c.JSON(http.StatusAccepted, gin.H{"message": "User created"})
}

func LoginUser(c *gin.Context) {

	// get user from context to know about the account details
	user := &model.LoginReq{}
	if err := c.ShouldBindJSON(&user); err != nil {
		fmt.Println(err)
	}
	// search for the account detail in db, and use database package function to check if user exist (using query)
	// if ID is "", then it does not exist and return bad request http status
	// else setcookies for successful login
	query := `SELECT id, username FROM users WHERE username = ?`
	ID, username := database.CheckIfExist(query, user.Id)
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user"})
		return
	}
	c.SetCookie("uid", ID, 36000, "/", "localhost", false, true)
	c.JSON(http.StatusAccepted, gin.H{"id": ID, "name": username})
}

func SetUser(userid, serverId string) {
	query := `INSERT INTO user_mapping (username, server_id) VALUES (?, ?)`
	err := database.ExecuteQuery(query, userid, serverId)
	if err != nil {
		fmt.Print(err)
	}
}

func GetServerId(userid string) string {
	var serverid string
	query := `SELECT server_id FROM user_mapping WHERE username = ?`
	iter := database.Connection.Session.Query(query, userid).Iter()
	iter.Scan(&serverid)
	return serverid
}

func CheckIfUserExist(userId string) (string, string) {
	query := `SELECT id,username FROM users WHERE username = ?`
	ID, username := database.CheckIfExist(query, userId)
	return ID, username
}
