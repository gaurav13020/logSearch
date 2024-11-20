package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"regexp"

	"github.com/go-redis/redis/v8"
)


type LogEntry struct {
	Data      map[string]interface{} `json:"data"`
	RawLog    string                 `json:"raw_log"`
	TimeStamp string                 `json:"timestamp"`
}

func MainFunction(uuid string, regex string, text string) {
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

	var logs []LogEntry
	if err := json.Unmarshal(data, &logs); err != nil {
		log.Fatalf("Failed to unmarshal logs: %v", err)
	}

	var filteredLogs []LogEntry
	for _, logEntry := range logs {
		if regex != "" {
			matched, err := regexp.MatchString(regex, logEntry.RawLog)
			if err != nil {
				log.Printf("Failed to match regex: %v", err)
				continue
			}
			if matched {
				log.Printf("Regex matched: %s", logEntry.RawLog)
				filteredLogs = append(filteredLogs, logEntry)
				continue
			}
		}
		if text != "" && containsText(logEntry.RawLog, text) {
			log.Printf("Text matched: %s", logEntry.RawLog)
			filteredLogs = append(filteredLogs, logEntry)
		}
	}

	filteredData, err := json.Marshal(filteredLogs)
	if err != nil {
		log.Fatalf("Failed to marshal filtered logs: %v", err)
	}

	// Write filtered logs to Redis using the provided UUID
	err = rdb.Set(ctx, uuid, filteredData, 0).Err()
	if err != nil {
		log.Printf("Failed to write logs to Redis: %v", err)
	} else {
		log.Printf("Logs written to Redis with UUID: %s", uuid)
	}
}

func containsText(rawLog, text string) bool {
	return regexp.MustCompile(`(?i)` + regexp.QuoteMeta(text)).MatchString(rawLog)
}
