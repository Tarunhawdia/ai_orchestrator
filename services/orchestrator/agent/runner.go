package agent

import (
	"context"
	"fmt"

	localtools "github.com/Tarunhawdia/decentralized-ai-orchestrator/services/orchestrator/tools"
	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/llms/googleai"
	"github.com/tmc/langchaingo/tools"
)

// Runner manages the execution of tasks via langchaingo agents
type Runner struct {
	llm *googleai.GoogleAI
}

// NewRunner creates a new Runner instance
func NewRunner(llm *googleai.GoogleAI) *Runner {
	return &Runner{
		llm: llm,
	}
}

// Run executes a given task request and returns the result
func (r *Runner) Run(ctx context.Context, request string) (string, error) {
	mySearchTool := localtools.NewSearchTool()
	agentTools := []tools.Tool{mySearchTool}

	agentExecutor := agents.NewExecutor(
		agents.NewOneShotAgent(
			r.llm,
			agentTools,
			agents.WithMaxIterations(3),
		),
	)

	agentInput := map[string]any{
		"input": request,
	}

	agentResponse, err := agentExecutor.Call(ctx, agentInput)
	if err != nil {
		return "", fmt.Errorf("agent execution failed: %w", err)
	}

	result, ok := agentResponse["output"].(string)
	if !ok {
		return "", fmt.Errorf("unexpected agent response format")
	}

	return result, nil
}
