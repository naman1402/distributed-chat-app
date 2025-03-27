// Package config implements Redis configuration and pub/sub mechanisms
// for distributed message handling across multiple server instances
package config

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// Global variables for Redis operations
var (
	// ctx provides the background context for Redis operations
	ctx = context.Background()

	// Conn represents the main Redis client instance
	// Used globally for all Redis operations including pub/sub
	Conn redis.Client
)

// NPool initializes and configures the Redis connection pool
// Implementation:
// 1. Creates new Redis client with server configuration
// 2. Verifies connection with ping
// 3. Sets global Redis client instance
// Panics if connection fails
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

// PubSub implements Redis publish/subscribe pattern
// Implementation:
// 1. Creates subscription to server-specific channel
// 2. Continuously listens for incoming messages
// 3. Forwards received messages to broadcast channel
// 4. Handles connection errors
// Note: SERVERID is used as the subscription channel
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
