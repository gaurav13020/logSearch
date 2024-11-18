package main

import (
	"io"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

func MainFunction(uuid string) {
	// Initialize Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379", // Adjust the address if needed
	})

	// Read logs from logs.json
	file, err := os.Open("logs.json")
	if err != nil {
		log.Fatalf("Failed to open logs.json: %v", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read logs.json: %v", err)
	}

	// Write logs to Redis using the provided UUID
	err = rdb.Set(ctx, uuid, data, 0).Err()
	if err != nil {
		log.Printf("Failed to write logs to Redis: %v", err)
	}else{log.Printf("Logs written to Redis with UUID: %s", uuid)}

	
	
	
}
