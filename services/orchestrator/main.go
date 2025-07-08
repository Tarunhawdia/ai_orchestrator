package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	localtools "github.com/Tarunhawdia/decentralized-ai-orchestrator/services/orchestrator/tools"
	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
	"github.com/tmc/langchaingo/tools"
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

	// --- Agent and Tool Integration ---
	// 1. Create an instance of your custom SearchTool
	mySearchTool := localtools.NewSearchTool()

	// 2. Define the tools available to the agent
	agentTools := []tools.Tool{mySearchTool}

	// 3. Create an agent executor
	// We'll use a ZeroShotAgent for simplicity, which decides what to do based on the prompt.
	// The prompt needs to clearly instruct the LLM on tool usage.
	agentExecutor := agents.NewExecutor(
		agents.NewOneShotAgent(
			llm,
			agentTools,
			agents.WithMaxIterations(3), 
		),
	)
	if err != nil {
		log.Fatalf("Failed to create agent executor: %v", err)
	}

	// 4. Define the agent's input (the question that might require a tool)
	agentInput := map[string]any{
		"input": "What is the current state of AI agent development according to recent news?",
	}

	// 5. Run the agent
	log.Println("\n--- Running Agent ---")
	agentResponse, err := agentExecutor.Call(ctx, agentInput)
	if err != nil {
		log.Fatalf("Agent execution failed: %v", err)
	}

	// 6. Print the agent's final answer
	fmt.Println("\n--- Agent's Final Answer ---")
	fmt.Printf("Agent Output: %s\n", agentResponse["output"])
	fmt.Println("---------------------------\n")
	// --- End Agent and Tool Integration ---

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
