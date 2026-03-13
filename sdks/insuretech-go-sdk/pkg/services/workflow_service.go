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

// StartWorkflow Start workflow instance
func (s *WorkflowService) StartWorkflow(ctx context.Context, req *models.WorkflowStartRequest) (*models.WorkflowStartResponse, error) {
	path := "/v1/workflow-instances"
	var result models.WorkflowStartResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetWorkflowDefinition Get workflow definition
func (s *WorkflowService) GetWorkflowDefinition(ctx context.Context, workflowDefinitionId string) (*models.WorkflowDefinitionRetrievalResponse, error) {
	path := "/v1/workflow-definitions/{workflow_definition_id}"
	path = strings.ReplaceAll(path, "{workflow_definition_id}", workflowDefinitionId)
	var result models.WorkflowDefinitionRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetWorkflowInstance Get workflow instance
func (s *WorkflowService) GetWorkflowInstance(ctx context.Context, workflowInstanceId string) (*models.WorkflowInstanceRetrievalResponse, error) {
	path := "/v1/workflow-instances/{workflow_instance_id}"
	path = strings.ReplaceAll(path, "{workflow_instance_id}", workflowInstanceId)
	var result models.WorkflowInstanceRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateWorkflowDefinition Create workflow definition
func (s *WorkflowService) CreateWorkflowDefinition(ctx context.Context, req *models.WorkflowDefinitionCreationRequest) (*models.WorkflowDefinitionCreationResponse, error) {
	path := "/v1/workflow-definitions"
	var result models.WorkflowDefinitionCreationResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CompleteTask Complete task
func (s *WorkflowService) CompleteTask(ctx context.Context, taskId string, req *models.WorkflowTaskCompletionRequest) (*models.WorkflowTaskCompletionResponse, error) {
	path := "/v1/workflow-tasks/{task_id}"
	path = strings.ReplaceAll(path, "{task_id}", taskId)
	var result models.WorkflowTaskCompletionResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetMyTasks Get my tasks
func (s *WorkflowService) GetMyTasks(ctx context.Context) (*models.MyTasksRetrievalResponse, error) {
	path := "/v1/workflow-tasks/my-tasks"
	var result models.MyTasksRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetWorkflowHistory Get workflow history for entity
func (s *WorkflowService) GetWorkflowHistory(ctx context.Context, entityType string, entityId string) (*models.WorkflowHistoryRetrievalResponse, error) {
	path := "/v1/entities/{entity_type}/{entity_id}/workflow-history"
	path = strings.ReplaceAll(path, "{entity_type}", entityType)
	path = strings.ReplaceAll(path, "{entity_id}", entityId)
	var result models.WorkflowHistoryRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

