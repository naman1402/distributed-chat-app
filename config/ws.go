package config

import (
	"fmt"
	"net/http"
	"syscall/js"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/naman1402/distributed-chat-app/controller"
	"github.com/naman1402/distributed-chat-app/model"
	"github.com/redis/go-redis/v9"
)

type Message struct {
	Id           string
	Message      string   `json:"msg"`
	Sender       string   `json:"sender"`
	Receiver     string   `json:"receiver,omitempty"`
	Group        bool     `json:"is_group"`
	GroupName    string   `json:"group_name,omitempty"`
	GroupMembers []string `json:"group_members,omitempty"`
	ServerId     string   `json:"server_id,omitempty"`
}

type ErrMessage struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

var broadcast = make(chan *redis.Message)
var clients = make(map[string]*websocket.Conn)
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var SERVERID string = ""

func WSHandler(w http.ResponseWriter, r *http.Request, c *gin.Context) {
	userId := c.Query("id")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Failed to upgrade: %+v", err)
		return
	}
	_, username := controller.CheckIfUserExist(userId)
	if username == ""{
		close ws func 
		return
	}


}

func ReceiveMessage() {}

func NewClient(userId string, conn *websocket.Conn) {

	controller.SetUser(userId, SERVERID)
	clients[userId] = conn
	clients[userId].WriteMessage(websocket.TextMessage, []byte("ok"))
}

func Send() {
	for {

		msg := <- broadcast
		message := Message{}
		if err := json.Unmarshal([]byte(msg.Payload), &message); err != nil {
			panic(err)
		}
		if message.Group {
			groupMessage(message)
			continue
		}
		client := clients[message.Receiver]
		if client == nil {
			fmt.Println("Reciever offline")
			continue
		}
		privateMessage ----- 
	}
}

func groupMessage() {}

func privateMessage() {}
