package domain

import (
	"context"
	"errors"

	workflowv1 "github.com/newage-saint/insuretech/gen/go/insuretech/workflow/entity/v1"
	workflowservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/workflow/services/v1"
)

// ─── Sentinel errors ──────────────────────────────────────────────────────────

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)

// ─── INPUT TYPES ──────────────────────────────────────────────────────────────

// DefinitionCreateInput carries all fields to create a workflow definition.
type DefinitionCreateInput struct {
	DefinitionID string
	Name         string
	Description  string
	Type         workflowv1.WorkflowType
	EntityType   string
	Steps        string // JSONB
	Conditions   string // JSONB
	Version      int32
	Status       workflowv1.WorkflowStatus
	CreatedBy    string
}

// InstanceCreateInput carries all fields to start a workflow instance.
type InstanceCreateInput struct {
	InstanceID           string
	WorkflowDefinitionID string
	EntityType           string
	EntityID             string
	Status               workflowv1.InstanceStatus
	CurrentStep          string
	Context              string // JSONB
	InitiatedBy          string
	CorrelationID        string
}

// TaskCreateInput carries fields to create a workflow task (step).
type TaskCreateInput struct {
	TaskID             string
	WorkflowInstanceID string
	StepName           string
	Type               workflowv1.WorkflowTaskType
	AssignedTo         string
	Status             workflowv1.WorkflowTaskStatus
	DueHours           int // hours from now for due_date
}

// TaskCompleteInput carries fields to complete a task.
type TaskCompleteInput struct {
	TaskID      string
	Decision    string // APPROVED, REJECTED, RETURNED
	Comments    string
	CompletedBy string
}

// ─── REPOSITORY INTERFACE ────────────────────────────────────────────────────

type WorkflowRepository interface {
	// Definitions
	CreateDefinition(ctx context.Context, input DefinitionCreateInput) (*workflowv1.WorkflowDefinition, error)
	GetDefinition(ctx context.Context, definitionID string) (*workflowv1.WorkflowDefinition, error)
	GetDefinitionByName(ctx context.Context, name string) (*workflowv1.WorkflowDefinition, error)
	ListDefinitions(ctx context.Context, entityType string, wfType workflowv1.WorkflowType, pageSize, offset int) ([]*workflowv1.WorkflowDefinition, int64, error)

	// Instances
	CreateInstance(ctx context.Context, input InstanceCreateInput) (*workflowv1.WorkflowInstance, error)
	GetInstance(ctx context.Context, instanceID string) (*workflowv1.WorkflowInstance, error)
	ListInstancesByEntity(ctx context.Context, entityType, entityID string) ([]*workflowv1.WorkflowInstance, error)
	UpdateInstanceStatus(ctx context.Context, instanceID string, status workflowv1.InstanceStatus, currentStep string) error
	CompleteInstance(ctx context.Context, instanceID string, status workflowv1.InstanceStatus) error

	// Tasks
	CreateTask(ctx context.Context, input TaskCreateInput) (*workflowv1.WorkflowTask, error)
	GetTask(ctx context.Context, taskID string) (*workflowv1.WorkflowTask, error)
	ListTasksByInstance(ctx context.Context, instanceID string) ([]*workflowv1.WorkflowTask, error)
	ListTasksByAssignee(ctx context.Context, assignedTo string, status workflowv1.WorkflowTaskStatus, pageSize, offset int) ([]*workflowv1.WorkflowTask, int64, error)
	CompleteTask(ctx context.Context, input TaskCompleteInput) (*workflowv1.WorkflowTask, error)
}

// ─── SERVICE INTERFACE ────────────────────────────────────────────────────────

type WorkflowService interface {
	CreateWorkflowDefinition(ctx context.Context, req *workflowservicev1.CreateWorkflowDefinitionRequest) (*workflowservicev1.CreateWorkflowDefinitionResponse, error)
	GetWorkflowDefinition(ctx context.Context, req *workflowservicev1.GetWorkflowDefinitionRequest) (*workflowservicev1.GetWorkflowDefinitionResponse, error)
	StartWorkflow(ctx context.Context, req *workflowservicev1.StartWorkflowRequest) (*workflowservicev1.StartWorkflowResponse, error)
	GetWorkflowInstance(ctx context.Context, req *workflowservicev1.GetWorkflowInstanceRequest) (*workflowservicev1.GetWorkflowInstanceResponse, error)
	GetMyTasks(ctx context.Context, req *workflowservicev1.GetMyTasksRequest, userID string) (*workflowservicev1.GetMyTasksResponse, error)
	CompleteTask(ctx context.Context, req *workflowservicev1.CompleteWorkflowTaskRequest, userID string) (*workflowservicev1.CompleteWorkflowTaskResponse, error)
	GetWorkflowHistory(ctx context.Context, req *workflowservicev1.GetWorkflowHistoryRequest) (*workflowservicev1.GetWorkflowHistoryResponse, error)
}
