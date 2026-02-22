package storage

import (
	"fmt"
	"sync"

	"github.com/Tarunhawdia/decentralized-ai-orchestrator/services/orchestrator/models"
)

// Store defines the interface for task storage
type Store interface {
	CreateTask(task *models.Task) error
	GetTask(id string) (*models.Task, error)
	UpdateTask(task *models.Task) error
}

// InMemoryStore is a thread-safe in-memory implementation of Store
type InMemoryStore struct {
	mu    sync.RWMutex
	tasks map[string]*models.Task
}

// NewInMemoryStore creates a new InMemoryStore
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		tasks: make(map[string]*models.Task),
	}
}

// CreateTask saves a new task
func (s *InMemoryStore) CreateTask(task *models.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.tasks[task.ID]; exists {
		return fmt.Errorf("task with ID %s already exists", task.ID)
	}
	// Store a copy
	taskCopy := *task
	s.tasks[task.ID] = &taskCopy
	return nil
}

// GetTask retrieves a task by ID
func (s *InMemoryStore) GetTask(id string) (*models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	task, exists := s.tasks[id]
	if !exists {
		return nil, fmt.Errorf("task %s not found", id)
	}
	// Return a copy
	taskCopy := *task
	return &taskCopy, nil
}

// UpdateTask updates an existing task
func (s *InMemoryStore) UpdateTask(task *models.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.tasks[task.ID]; !exists {
		return fmt.Errorf("task %s not found", task.ID)
	}
	// Store a copy
	taskCopy := *task
	s.tasks[task.ID] = &taskCopy
	return nil
}
