package config

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// ctx is global context, Conn is redis Client used to connect and interact with the Redis server
var ctx = context.Background()
var Conn redis.Client

// var broadcast chan *redis.Message

// Redis setup function, initialise new Redis Client and connect it with redis server
// pinging the connection to ensure the connectivity,
// if connection is successful, set the global variable as the newly initiliased client
func NPool() {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	Conn = *rdb
}

// subscribing to the channel with ID using redis client, if successful subscriber will receive message published to that channel
// message received is stored in msg, and passed to the broadcast channel (redis.Message)
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
