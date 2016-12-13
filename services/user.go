package services

import (
	"log"
	"gopkg.in/redis.v5"
)

func RedisClient() *redis.Client{
    Client := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })

    pong, err := Client.Ping().Result()
    log.Println(pong, err)
    // Output: PONG <nil>
	return Client
	
}
