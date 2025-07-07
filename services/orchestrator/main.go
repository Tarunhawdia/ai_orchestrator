package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// Simple HTTP handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from the Orchestrator Service!")
	})

	// Determine the port to listen on
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