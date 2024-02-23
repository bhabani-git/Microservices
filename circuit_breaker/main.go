package main

import (
	"encoding/json"
	"net/http"

	"github.com/afex/hystrix-go/hystrix"
)

// Message represents a chat message
type Message struct {
	ID      string `json:"id"`
	UserID  string `json:"userId"`
	Content string `json:"content"`
}

// Mock message data (replace with database implementation)
var messages []Message

// SendMessageHandler handles message sending requests
func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	// Execute the service call with circuit breaker
	err := hystrix.Do("sendMessage", func() error {
		// Parse request body
		var msg Message
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			return err
		}

		// Store message (not shown here)
		// ...

		// Return success response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Message sent"))
		return nil
	}, nil)

	// Check for circuit breaker tripping
	if err != nil {
		http.Error(w, "Failed to send message1", http.StatusInternalServerError)
	}
}

func main() {
	// Configure Hystrix settings
	hystrix.ConfigureCommand("sendMessage", hystrix.CommandConfig{
		Timeout:                1,    // Timeout in milliseconds
		MaxConcurrentRequests:  2,    // Maximum concurrent requests
		RequestVolumeThreshold: 1,    // Minimum number of requests before tripping the circuit breaker
		ErrorPercentThreshold:  50,   // Error threshold percentage for tripping the circuit breaker
		SleepWindow:            5000, // Duration in milliseconds for which to sleep after circuit breaker trips
	})

	// Register send message handler
	http.HandleFunc("/send", SendMessageHandler)

	// Start HTTP server
	if err := http.ListenAndServe(":8081", nil); err != nil {
		panic(err)
	}
}
