package internal

import (
	"fmt"
	"github.com/go-redis/redis"
)

func NewRedis(local bool) *redis.Client {
	fmt.Println("Attempting to connect to redis")
	var address string
	if local {
		fmt.Printf("Running in LOCAL mode, connecting to localhost...\n")
		address = "localhost:6379"
	} else {
		fmt.Printf("Running in PRODUCTION mode, connecting to messaging...\n")
		address = "kv:6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "",
		DB:       0,
	})
	_, err := client.Ping().Result()
	FailOnError(err, "Could not connect to Redis")
	fmt.Println("Successfully connected!")
	return client
}
