package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Tarunhawdia/decentralized-ai-orchestrator/services/orchestrator/agent"
	"github.com/Tarunhawdia/decentralized-ai-orchestrator/services/orchestrator/api"
	"github.com/Tarunhawdia/decentralized-ai-orchestrator/services/orchestrator/storage"
	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/llms/googleai"
)

func main() {
	// Load environment variables from .env file if it exists.
	err := godotenv.Load()
	if err != nil {
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

	// Initialize components
	store := storage.NewInMemoryStore()
	runner := agent.NewRunner(llm)
	handlers := api.NewHandlers(store, runner)

	// Setup routes
	http.HandleFunc("/tasks", handlers.HandleSubmitTask) // POST
	http.HandleFunc("/tasks/", handlers.HandleGetTask)   // GET

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		fmt.Fprintf(w, "Hello from the Orchestrator Service! Use /tasks to submit work.")
	})

	// Determine the port to listen on from environment variable or default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	log.Printf("Orchestrator Service starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Orchestrator Service failed to start: %v", err)
	}
}
