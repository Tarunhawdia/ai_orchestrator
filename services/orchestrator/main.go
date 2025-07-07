package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv" // Import the godotenv library
)

func main() {
	// Load environment variables from .env file if it exists.
	// godotenv.Load() will look for a .env file in the current directory.
	// It's fine if the file doesn't exist (e.g., in a production container),
	// as os.Getenv will then fall back to actual environment variables.
	err := godotenv.Load()
	if err != nil {
		// Log a non-fatal error if .env file is not found,
		// as it might be expected in containerized environments.
		log.Println("No .env file found, relying on environment variables or defaults.")
	}

	// Simple HTTP handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from the Orchestrator Service!")
	})

	// Determine the port to listen on from environment variable or default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	log.Printf("Orchestrator Service starting on port %s...", port)
	// Start the HTTP server
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Orchestrator Service failed to start: %v", err)
	}
}
