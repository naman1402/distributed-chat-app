// Package config implements WebSocket handling for a distributed chat system
// It manages real-time message distribution across multiple server instances using Redis pub/sub
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

// Message defines the structure for chat messages with support for both private and group communications
// Id: Unique message identifier generated using KSUID
// Message: The actual message content
// Sender: UserID of message sender
// Receiver: UserID of recipient (for private messages)
// Group: Flag indicating if message is for group chat
// GroupName: Name of the group for group messages
// GroupMembers: List of users in the group
// ServerId: ID of server handling the message
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

// ErrMessage defines the structure for error messages
type ErrMessage struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Global connection management variables
var (
	// broadcast acts as a message queue for Redis pub/sub messages
	broadcast = make(chan *redis.Message)

	// clients maintains mapping of active user connections
	// key: userId, value: websocket connection
	clients = make(map[string]*websocket.Conn)

	// upgrader configures WebSocket connection parameters
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	// SERVERID uniquely identifies this server instance in the distributed system
	SERVERID string = ""
)

// WSHandler establishes and manages WebSocket connections
// 1. Extracts user ID from request
// 2. Upgrades HTTP connection to WebSocket
// 3. Validates user authentication
// 4. Initializes client connection
// Parameters:
//   - w: HTTP response writer
//   - r: HTTP request
//   - c: Gin context containing user information
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
	// Initializes client connection
	NewClient(userId, conn)
	ReceiveMessage(conn, userId)
}

// ReceiveMessage processes incoming WebSocket messages
// Implementation:
// 1. Continuously reads messages from WebSocket
// 2. Deserializes JSON messages
// 3. Validates message content
// 4. Routes messages to appropriate handlers (group/private)
// 5. Persists messages to database
// 6. Publishes to Redis for cross-server communication
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
// NewClient registers a new WebSocket client connection
// 1. Records user-server mapping in database
// 2. Stores WebSocket connection in memory
// 3. Sends connection acknowledgment
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
// Send implements the message distribution system
// 1. Listens to Redis broadcast channel
// 2. Deserializes incoming messages
// 3. Routes to group/private message handlers
// 4. Handles offline user scenarios
func Send() {
	for {

		msg := <-broadcast
		message := Message{}
		err := json.Unmarshal([]byte(msg.Payload), &message)
		if err != nil {
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
// groupMessage implements group chat message distribution
// 1. Creates a new message instance for each recipient
// 2. Checks recipient connection status
// 3. Delivers message to all online group members
// 4. Handles connection failures and cleanup
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
// privateMessage handles one-to-one message delivery
// 1. Serializes message to JSON
// 2. Delivers to recipient's WebSocket connection
// 3. Handles connection errors and cleanup
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

// CloseWS performs graceful WebSocket connection termination
// Sends close frame with custom message before closing
func CloseWS(msg string, conn *websocket.Conn) {
	cm := websocket.FormatCloseMessage(websocket.CloseNormalClosure, msg)
	if err := conn.WriteMessage(websocket.CloseMessage, cm); err != nil {
		fmt.Println(err)
		return
	}
	conn.Close()
}

// MsgFailed notifies client of message delivery failure
// Sends standardized error JSON response
func MsgFailed(conn *websocket.Conn) {

	msg := `{"message": "Failed to send message"}`
	if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		fmt.Println(err)
		return
	}
}

// Validate implements message validation rules
// Validates:
// - Message content: Required, non-empty, length 1-1000 chars
// - Group flag: Must be non-nil
// - Group name: Length 1-25 chars when present
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
