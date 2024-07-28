package controller

import (
	"fmt"

	"github.com/naman1402/distributed-chat-app/database"
)

func SaveMessagePrivateChat(id, msg, sender, reciever string) {
	query := `INSERT INTO private_chat(id, msg, sender, reciever, timestamp) VALUES (?, ?, ?, ?, toTimeStamp(now()))`
	err := database.ExecuteQuery(query, id, msg, sender, reciever)
	if err != nil {
		fmt.Println(err)
	}
}

func SaveMessageGroupChat(id, msg, sender, groupName string) {
	query := `INSERT INTO group_chat(id, msg, sender, group, timestamp) VALUES (?, ?, ?, ?, toTimeStamp(now()))`
	err := database.ExecuteQuery(query, id, msg, sender, groupName)
	if err != nil {
		fmt.Println(err)
	}
}
