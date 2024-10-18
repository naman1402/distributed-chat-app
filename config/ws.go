package config

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/naman1402/distributed-chat-app/controller"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/ksuid"
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

var broadcast = make(chan *redis.Message)      // channel to receive message
var clients = make(map[string]*websocket.Conn) // userid to websocket connection
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}                        // converts HTTP connections into Websocket connections
var SERVERID string = "" // Serverid used for redis pubsub connection

/*
* gets userid from context, allows connections from any origin and upgrading http connection to a websocket connection, return websocket.Conn
check error in upgrading, getting username of id from controller. If username is empty then close websocket connection
Register new client (id) and their websocket connection (conn) and Listens to upcoming websocket messagess
*/
func WSHandler(w http.ResponseWriter, r *http.Request, c *gin.Context) {
	userId := c.Query("id")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	// upgrades to a websocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Failed to upgrade: %+v", err)
		return
	}
	_, username := controller.CheckIfUserExist(userId)
	if username == "" {
		CloseWS("Authentication failed - invalid username", conn)
		return
	}

	NewClient(userId, conn)
	ReceiveMessage(conn, userId)
}

/*
*
Listens to incoming websocket messages
reading messages from the websocket connection and check for errors
convert msg json into res Message using Unmarshal and check for errors
create new id and keep res.Id as the string form of this id, set sender as userID
if group of res exists then
send through private chat
*/
func ReceiveMessage(conn *websocket.Conn, userID string) {

	for {
		_, msg, errCon := conn.ReadMessage()
		if errCon != nil {
			log.Println("Read error: ", errCon)
			break
		}
		var res Message
		if err := json.Unmarshal(msg, &res); err != nil {
			log.Println("error: " + err.Error())
			MsgFailed(conn)
			continue
		}
		id := ksuid.New()
		res.Id = id.String()
		res.Sender = userID
		err := res.Validate()
		if err != nil {
			b, _ := json.Marshal(err)
			conn.WriteMessage(websocket.TextMessage, b)
			continue
		}

		// saves the message in db and get members of the groupname
		// for all members stored their serverid to member mapping
		// using loop, iterate through all servers and publish the message on redis client
		if res.Group {
			controller.SaveMessageGroupChat(res.Id, res.Message, res.Sender, res.GroupName)
			members := controller.GetMembersFromRoom(res.GroupName)
			servers := make(map[string][]string)
			for _, member := range members {
				serverId := controller.GetServerId(member)
				servers[serverId] = append(servers[serverId], member)
			}

			for key, element := range servers {
				res.ServerId = key
				res.GroupMembers = element
				jsonData, err := json.Marshal(res)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println("redis key ", key)
				//////////////////////////////////////////////////////
				// used to send messages to a specified Redis channel//
				/////////////////////////////////////////////////////
				// ctx is the context for the operaion, key is the name of Redis channel(serverid) to which the message will be published
				// jsonData is message we want to send
				Conn.Publish(ctx, key, jsonData)
			}
			continue
		}
		// logic to execute private chat, publishing message on redis Client
		controller.SaveMessagePrivateChat(res.Id, res.Message, res.Sender, res.Receiver)
		serverId := controller.GetServerId(res.Receiver)
		jsonData, err := json.Marshal(res)
		if err != nil {
			fmt.Println(err)
			return
		}
		Conn.Publish(ctx, serverId, jsonData)
	}

	cm := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "connection closing")
	if err := conn.WriteMessage(websocket.CloseMessage, cm); err != nil {
		fmt.Println(err)
		return
	}
	// closing websocket connection
	conn.Close()
}

// stores the user in db, id -> conn mapping, sends ack message ok to newly connected client
func NewClient(userId string, conn *websocket.Conn) {

	controller.SetUser(userId, SERVERID)
	clients[userId] = conn
	clients[userId].WriteMessage(websocket.TextMessage, []byte("ok"))
}

/*
*
infinite loop, get msg from broadcast channel (redis message), using json.unmarshal get msg (redis message) to message (json form)
if group exist, call groupMessage function, if not group then private message
get conn of receiver of message, it conn is nil then receiver is offline
send private message(json) to recevier conn
*/
func Send() {
	for {

		msg := <-broadcast
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
		privateMessage(message, client)
	}
}

/*
*
declares empty message, loop for all member. get their connection (client) and check if connection is online or not
sets attributes of response message, use marshal to convert res into json data and sends using client(conn).WriteMessage to member
after sending, remove receiver from clients and close respective connection
*/
func groupMessage(message Message) {
	res := Message{}
	for _, member := range message.GroupMembers {
		client := clients[member]
		if client == nil {
			fmt.Println("Reciever offline")
			continue
		}

		res.Id = message.Id
		res.Sender = message.Sender
		res.Message = message.Message
		res.Group = message.Group
		res.GroupName = message.GroupName
		data, err := json.Marshal(res)
		if err != nil {
			fmt.Println(err)
			return
		}
		// send message using websocket connection
		err = client.WriteMessage(websocket.TextMessage, []byte(data))
		if err != nil {
			delete(clients, message.Receiver)
			client.Close()
		}
	}
}

/*
*
directly to specific user, marshal data into jsonData and send using client.WriteMessage()
*/
func privateMessage(message Message, client *websocket.Conn) {
	jsonData, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err)
		return
	}
	// TO WRITE MESSAGE we send message using websocket connection
	err = client.WriteMessage(websocket.TextMessage, []byte(jsonData))
	if err != nil {
		delete(clients, message.Receiver)
		client.Close()
	}
}

func CloseWS(msg string, conn *websocket.Conn) {
	cm := websocket.FormatCloseMessage(websocket.CloseNormalClosure, msg)
	if err := conn.WriteMessage(websocket.CloseMessage, cm); err != nil {
		fmt.Println(err)
		return
	}
	conn.Close()
}

func MsgFailed(conn *websocket.Conn) {

	msg := `{"message": "Failed to send message"}`
	if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		fmt.Println(err)
		return
	}
}

func (m Message) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Message,
			validation.Required.Error("msg field is required"),
			validation.NotNil.Error("msg field cannot be empty"),
			validation.Length(1, 1000).Error("character length should be between 1 and 1000"),
		),
		validation.Field(&m.Group,
			validation.NotNil.Error("is_group field cannot be empty"),
		),
		validation.Field(&m.GroupName,
			validation.Length(1, 25).Error("character length should be between 1 and 25"),
		),
	)
}
