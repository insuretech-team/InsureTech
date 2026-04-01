package services

import (
	"context"
	"strings"
	"github.com/newage-saint/insuretech-go-sdk/pkg/models"
)

// WorkflowService handles workflow-related API calls
type WorkflowService struct {
	Client Client
}

// GetWorkflowHistory Get workflow history for entity
func (s *WorkflowService) GetWorkflowHistory(ctx context.Context, entityType string, entityId string) error {
	path := "/v1/entities/{entity_type}/{entity_id}/workflow-history"
	path = strings.ReplaceAll(path, "{entity_type}", entityType)
	path = strings.ReplaceAll(path, "{entity_id}", entityId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// CreateWorkflowDefinition Create workflow definition
func (s *WorkflowService) CreateWorkflowDefinition(ctx context.Context, req *models.WorkflowDefinitionCreationRequest) error {
	path := "/v1/workflow-definitions"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// GetWorkflowDefinition Get workflow definition
func (s *WorkflowService) GetWorkflowDefinition(ctx context.Context, workflowDefinitionId string) error {
	path := "/v1/workflow-definitions/{workflow_definition_id}"
	path = strings.ReplaceAll(path, "{workflow_definition_id}", workflowDefinitionId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// StartWorkflow Start workflow instance
func (s *WorkflowService) StartWorkflow(ctx context.Context, req *models.WorkflowStartRequest) error {
	path := "/v1/workflow-instances"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// GetWorkflowInstance Get workflow instance
func (s *WorkflowService) GetWorkflowInstance(ctx context.Context, workflowInstanceId string) error {
	path := "/v1/workflow-instances/{workflow_instance_id}"
	path = strings.ReplaceAll(path, "{workflow_instance_id}", workflowInstanceId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// GetMyTasks Get my tasks
func (s *WorkflowService) GetMyTasks(ctx context.Context) error {
	path := "/v1/workflow-tasks/my-tasks"
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// CompleteTask Complete task
func (s *WorkflowService) CompleteTask(ctx context.Context, taskId string, req *models.WorkflowTaskCompletionRequest) error {
	path := "/v1/workflow-tasks/{task_id}:complete"
	path = strings.ReplaceAll(path, "{task_id}", taskId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

