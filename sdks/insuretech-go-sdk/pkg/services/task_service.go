package services

import (
	"context"
	"strings"
	"github.com/newage-saint/insuretech-go-sdk/pkg/models"
)

// TaskService handles task-related API calls
type TaskService struct {
	Client Client
}

// CreateTask Create task
func (s *TaskService) CreateTask(ctx context.Context, req *models.TaskCreationRequest) error {
	path := "/v1/tasks"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// ListMyTasks List my tasks
func (s *TaskService) ListMyTasks(ctx context.Context) error {
	path := "/v1/tasks/my-tasks"
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// GetTask Get task
func (s *TaskService) GetTask(ctx context.Context, taskId string) error {
	path := "/v1/tasks/{task_id}"
	path = strings.ReplaceAll(path, "{task_id}", taskId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// UpdateTask Update task
func (s *TaskService) UpdateTask(ctx context.Context, taskId string, req *models.TaskUpdateRequest) error {
	path := "/v1/tasks/{task_id}"
	path = strings.ReplaceAll(path, "{task_id}", taskId)
	return s.Client.DoRequest(ctx, "PATCH", path, req, nil)
}

// AssignTask Assign task
func (s *TaskService) AssignTask(ctx context.Context, taskId string, req *models.TaskAssignmentRequest) error {
	path := "/v1/tasks/{task_id}:assign"
	path = strings.ReplaceAll(path, "{task_id}", taskId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// CompleteTask Complete task
func (s *TaskService) CompleteTask(ctx context.Context, taskId string, req *models.TaskCompletionRequest) error {
	path := "/v1/tasks/{task_id}:complete"
	path = strings.ReplaceAll(path, "{task_id}", taskId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

