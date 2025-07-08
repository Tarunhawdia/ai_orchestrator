package tools

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// SearchTool is a concrete implementation of the langchaingo/tools.Tool interface.
// It simulates a web search.
type SearchTool struct{}

// NewSearchTool creates and returns a pointer to a new instance of SearchTool.
// This is the constructor function you'll call from main.go.
func NewSearchTool() *SearchTool {
	return &SearchTool{}
}

// Name returns the name of the tool. This is how the LLM will refer to it.
func (t *SearchTool) Name() string {
	return "search_internet"
}

// Description returns a detailed description of what the tool does and its input.
// This description is crucial for the LLM to understand when and how to use the tool.
func (t *SearchTool) Description() string {
	return "Use this tool to search the internet for information. Input should be a concise search query."
}

// Call executes the tool's logic. For now, it returns a hardcoded simulated result.
func (t *SearchTool) Call(ctx context.Context, input string) (string, error) {
	log.Printf("SearchTool called with input: '%s'", input)

	serperAPIKey := os.Getenv("SERPER_API_KEY")
	if serperAPIKey == "" {
		return "", fmt.Errorf("SERPER_API_KEY environment variable is not set")
	}

	url := "https://google.serper.dev/search"
	method := "POST"
	payload := strings.NewReader(fmt.Sprintf(`{"q":"%s","gl":"in"}`, input))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("X-API-KEY", serperAPIKey)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	return string(body), nil
}
