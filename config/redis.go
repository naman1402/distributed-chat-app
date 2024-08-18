package config

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var Conn redis.Client

// manages websocket connections, connection lifecycle
// Ping is used for connection check
func NPool() {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	Conn = *rdb
}

// creates channel for individual chats and group chats
// subscribes to channel identified by serverid using client conn and ctx background connection
// subscriber is redis.PubSub !!!!!
// in infinite loop, continously listens for message on sub
// if message is received, add it to broadcast channel + error handling
func PubSub() {
	SERVERID := ""
	fmt.Println(SERVERID)
	subscriber := Conn.Subscribe(ctx, SERVERID)
	for {
		msg, err := subscriber.ReceiveMessage(ctx)
		if err != nil {
			panic(err)
		}
		broadcast <- msg
	}

}
