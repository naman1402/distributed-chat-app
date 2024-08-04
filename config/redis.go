package config

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var conn redis.Client

// manages websocket connections, connection lifecycle
func NPool() {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	conn = *rdb
}

// creates channel for individual chats and group chats
func PubSub() {
	SERVERID := ""
	fmt.Println(SERVERID)
	subscriber := conn.Subscribe(ctx, SERVERID)
	for {
		msg, err := subscriber.ReceiveMessage(ctx)
		if err != nil {
			panic(err)
		}
		broadcast <- msg
	}

}
