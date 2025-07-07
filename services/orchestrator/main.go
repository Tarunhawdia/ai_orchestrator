package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
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

	//----llm integration----//
	geminiAPIKey := os.Getenv("GOOGLE_API_KEY")
	geminiModelName := os.Getenv("GOOGLE_GEMINI_MODEL")

	if geminiAPIKey == "" {
		log.Fatalf("GOOGLE_API_KEY environment variable not set. Please set it in your .env file or environment.")
	}
	if geminiModelName == "" {
		log.Println("GOOGLE_GEMINI_MODEL not set, defaulting to 'gemini-pro'.")
		geminiModelName = "gemini-pro"
	}

	//context for the LLM call
	// This context is used for the LLM call and can be extended with timeouts or
	// cancellation if needed in the future.
	ctx := context.Background()

	llm, err := googleai.New(
		ctx,
		googleai.WithAPIKey(geminiAPIKey),
		googleai.WithDefaultModel(geminiModelName),
	)

	if err != nil {
		log.Fatalf("Failed to create Google Gemini LLM client: %v", err)
	}

	log.Printf("Successfully initialized Google Gemini LLM with model: %s", geminiModelName)

	// prompt for the llm
	prompt := "what is the capital of india?"

	llmResponse, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("LLM Response: %s", llmResponse)
	// --- End LLM Integration ---

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
