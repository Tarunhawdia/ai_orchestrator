package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Tarunhawdia/decentralized-ai-orchestrator/services/orchestrator/agent"
	"github.com/Tarunhawdia/decentralized-ai-orchestrator/services/orchestrator/models"
	"github.com/Tarunhawdia/decentralized-ai-orchestrator/services/orchestrator/storage"
	"github.com/google/uuid"
)

type Handlers struct {
	store  storage.Store
	runner *agent.Runner
}

func NewHandlers(store storage.Store, runner *agent.Runner) *Handlers {
	return &Handlers{
		store:  store,
		runner: runner,
	}
}

type SubmitTaskRequest struct {
	Prompt string `json:"prompt"`
}

type SubmitTaskResponse struct {
	TaskID string `json:"task_id"`
}

func (h *Handlers) HandleSubmitTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SubmitTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	taskID := uuid.New().String()
	task := &models.Task{
		ID:        taskID,
		Request:   req.Prompt,
		Status:    models.StatusPending,
		CreatedAt: time.Now(),
	}

	if err := h.store.CreateTask(task); err != nil {
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	// Start processing asynchronously
	go h.processTask(taskID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(SubmitTaskResponse{TaskID: taskID})
}

func (h *Handlers) processTask(taskID string) {
	task, err := h.store.GetTask(taskID)
	if err != nil {
		return // Should not happen
	}

	task.Status = models.StatusProcessing
	h.store.UpdateTask(task)

	// Using a background context for the async job without timeout for now
	ctx := context.Background()
	result, err := h.runner.Run(ctx, task.Request)

	// Fetch fresh copy before update just in case, although we're the only ones editing it.
	task, _ = h.store.GetTask(taskID)

	if err != nil {
		task.Status = models.StatusFailed
		task.Error = err.Error()
	} else {
		task.Status = models.StatusCompleted
		task.Result = result
	}

	h.store.UpdateTask(task)
}

func (h *Handlers) HandleGetTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Simple path parsing, assuming /tasks/{id}
	path := r.URL.Path
	prefix := "/tasks/"
	if len(path) <= len(prefix) {
		http.Error(w, "Task ID required", http.StatusBadRequest)
		return
	}
	taskID := path[len(prefix):]

	task, err := h.store.GetTask(taskID)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}
