package models

import "time"

// TaskStatus represents the current state of a task
type TaskStatus string

const (
	StatusPending    TaskStatus = "pending"
	StatusProcessing TaskStatus = "processing"
	StatusCompleted  TaskStatus = "completed"
	StatusFailed     TaskStatus = "failed"
)

// Task represents a unit of work for the agent
type Task struct {
	ID        string     `json:"id"`
	Request   string     `json:"request"`
	Status    TaskStatus `json:"status"`
	Result    string     `json:"result,omitempty"`
	Error     string     `json:"error,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}
