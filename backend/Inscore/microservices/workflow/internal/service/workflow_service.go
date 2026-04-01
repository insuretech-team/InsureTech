// Package service implements workflow engine business logic.
// It handles: definition CRUD, instance lifecycle (start/advance/complete),
// task assignment and completion, and Kafka event publishing.
package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/workflow/internal/domain"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/workflow/internal/events"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/grpcmeta"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	workflowv1 "github.com/newage-saint/insuretech/gen/go/insuretech/workflow/entity/v1"
	workfloweventsv1 "github.com/newage-saint/insuretech/gen/go/insuretech/workflow/events/v1"
	workflowservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/workflow/services/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ─── STEP DEFINITION (parsed from definition.Steps JSONB) ────────────────────

// StepDef describes one step in a workflow template.
type StepDef struct {
	Name       string `json:"name"`
	Type       string `json:"type"`        // APPROVAL, REVIEW, NOTIFICATION, ACTION
	AssignRole string `json:"assign_role"` // role hint, resolved externally
	AssignTo   string `json:"assign_to"`   // explicit user UUID (optional)
	DueHours   int    `json:"due_hours"`   // 0 = no deadline
	Order      int    `json:"order"`
}

// ─── SERVICE IMPLEMENTATION ───────────────────────────────────────────────────

type workflowService struct {
	repo      domain.WorkflowRepository
	publisher *events.Publisher
}

func New(repo domain.WorkflowRepository, publisher *events.Publisher) domain.WorkflowService {
	return &workflowService{repo: repo, publisher: publisher}
}

// ─── DEFINITIONS ─────────────────────────────────────────────────────────────

func (s *workflowService) CreateWorkflowDefinition(
	ctx context.Context,
	req *workflowservicev1.CreateWorkflowDefinitionRequest,
) (*workflowservicev1.CreateWorkflowDefinitionResponse, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("%w: name is required", ErrInvalidArgument)
	}
	if req.EntityType == "" {
		return nil, fmt.Errorf("%w: entity_type is required", ErrInvalidArgument)
	}

	// Marshal steps struct to JSON string for JSONB storage
	stepsJSON := "{}"
	if req.Steps != nil {
		b, err := req.Steps.MarshalJSON()
		if err == nil {
			stepsJSON = string(b)
		}
	}

	// Validate steps are parseable
	if err := validateSteps(stepsJSON); err != nil {
		return nil, fmt.Errorf("%w: steps: %v", ErrInvalidArgument, err)
	}

	wfType := parseWorkflowTypeStr(req.Type)

	def, err := s.repo.CreateDefinition(ctx, domain.DefinitionCreateInput{
		DefinitionID: uuid.NewString(),
		Name:         req.Name,
		Description:  req.Description,
		Type:         wfType,
		EntityType:   req.EntityType,
		Steps:        stepsJSON,
		Version:      1,
		Status:       workflowv1.WorkflowStatus_WORKFLOW_STATUS_ACTIVE,
	})
	if err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			return nil, fmt.Errorf("%w: %v", ErrAlreadyExists, err)
		}
		return nil, fmt.Errorf("create workflow definition: %w", err)
	}

	appLogger.Infof("Created workflow definition %s (name=%s, entity=%s)", def.WorkflowDefinitionId, def.Name, def.EntityType)

	return &workflowservicev1.CreateWorkflowDefinitionResponse{
		WorkflowDefinitionId: def.WorkflowDefinitionId,
		Message:              "Workflow definition created successfully",
	}, nil
}

func (s *workflowService) GetWorkflowDefinition(
	ctx context.Context,
	req *workflowservicev1.GetWorkflowDefinitionRequest,
) (*workflowservicev1.GetWorkflowDefinitionResponse, error) {
	if req.WorkflowDefinitionId == "" {
		return nil, fmt.Errorf("%w: workflow_definition_id is required", ErrInvalidArgument)
	}
	def, err := s.repo.GetDefinition(ctx, req.WorkflowDefinitionId)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, fmt.Errorf("%w: %v", ErrNotFound, err)
		}
		return nil, fmt.Errorf("get workflow definition: %w", err)
	}
	return &workflowservicev1.GetWorkflowDefinitionResponse{WorkflowDefinition: def}, nil
}

// ─── INSTANCES ───────────────────────────────────────────────────────────────

func (s *workflowService) StartWorkflow(
	ctx context.Context,
	req *workflowservicev1.StartWorkflowRequest,
) (*workflowservicev1.StartWorkflowResponse, error) {
	if req.WorkflowDefinitionId == "" || req.EntityId == "" || req.EntityType == "" {
		return nil, fmt.Errorf("%w: workflow_definition_id, entity_type and entity_id are required", ErrInvalidArgument)
	}

	// Load the definition to get its steps template
	def, err := s.repo.GetDefinition(ctx, req.WorkflowDefinitionId)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, fmt.Errorf("%w: %v", ErrNotFound, err)
		}
		return nil, fmt.Errorf("load workflow definition: %w", err)
	}

	// Parse steps from definition JSONB
	steps, err := parseSteps(def.Steps)
	if err != nil {
		return nil, fmt.Errorf("parse workflow steps: %w", err)
	}

	// Serialise the caller-supplied context struct
	contextJSON := "{}"
	if req.Context != nil {
		b, _ := req.Context.MarshalJSON()
		contextJSON = string(b)
	}

	// Determine initiated_by from gRPC metadata.
	initiatedBy := metaUserID(ctx)
	if initiatedBy == "" {
		return nil, fmt.Errorf("%w: initiated_by is required (missing x-user-id metadata)", ErrInvalidArgument)
	}

	instanceID := uuid.NewString()
	firstStep := ""
	if len(steps) > 0 {
		firstStep = steps[0].Name
	}

	inst, err := s.repo.CreateInstance(ctx, domain.InstanceCreateInput{
		InstanceID:           instanceID,
		WorkflowDefinitionID: req.WorkflowDefinitionId,
		EntityType:           req.EntityType,
		EntityID:             req.EntityId,
		Status:               workflowv1.InstanceStatus_INSTANCE_STATUS_IN_PROGRESS,
		CurrentStep:          firstStep,
		Context:              contextJSON,
		InitiatedBy:          initiatedBy,
	})
	if err != nil {
		return nil, fmt.Errorf("create workflow instance: %w", err)
	}

	// Create tasks for all steps from the template
	for _, step := range steps {
		taskType := parseTaskTypeStr(step.Type)
		_, taskErr := s.repo.CreateTask(ctx, domain.TaskCreateInput{
			TaskID:             uuid.NewString(),
			WorkflowInstanceID: inst.Id,
			StepName:           step.Name,
			Type:               taskType,
			AssignedTo:         step.AssignTo,
			Status:             workflowv1.WorkflowTaskStatus_WORKFLOW_TASK_STATUS_PENDING,
			DueHours:           step.DueHours,
		})
		if taskErr != nil {
			appLogger.Errorf("Failed to create task for step %s in instance %s: %v", step.Name, inst.Id, taskErr)
		}
	}

	// Publish WorkflowStartedEvent
	s.publisher.Publish(ctx, events.TopicWorkflowStarted, inst.Id, &workfloweventsv1.WorkflowStartedEvent{
		EventId:              uuid.NewString(),
		WorkflowInstanceId:   inst.Id,
		WorkflowDefinitionId: req.WorkflowDefinitionId,
		EntityType:           req.EntityType,
		EntityId:             req.EntityId,
		InitiatedBy:          initiatedBy,
		Timestamp:            timestamppb.Now(),
	})

	appLogger.Infof("Started workflow instance %s for %s/%s using definition %s",
		inst.Id, req.EntityType, req.EntityId, req.WorkflowDefinitionId)

	return &workflowservicev1.StartWorkflowResponse{
		WorkflowInstanceId: inst.Id,
		Message:            "Workflow started successfully",
	}, nil
}

func (s *workflowService) GetWorkflowInstance(
	ctx context.Context,
	req *workflowservicev1.GetWorkflowInstanceRequest,
) (*workflowservicev1.GetWorkflowInstanceResponse, error) {
	if req.WorkflowInstanceId == "" {
		return nil, fmt.Errorf("%w: workflow_instance_id is required", ErrInvalidArgument)
	}
	inst, err := s.repo.GetInstance(ctx, req.WorkflowInstanceId)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, fmt.Errorf("%w: %v", ErrNotFound, err)
		}
		return nil, err
	}
	tasks, err := s.repo.ListTasksByInstance(ctx, inst.Id)
	if err != nil {
		appLogger.Warnf("Failed to list tasks for instance %s: %v", inst.Id, err)
		tasks = nil
	}
	return &workflowservicev1.GetWorkflowInstanceResponse{
		WorkflowInstance: inst,
		Tasks:            tasks,
	}, nil
}

// ─── TASKS ───────────────────────────────────────────────────────────────────

func (s *workflowService) GetMyTasks(
	ctx context.Context,
	req *workflowservicev1.GetMyTasksRequest,
	userID string,
) (*workflowservicev1.GetMyTasksResponse, error) {
	if userID == "" {
		return nil, fmt.Errorf("%w: user_id is required (missing x-user-id metadata)", ErrInvalidArgument)
	}
	pageSize := int(req.PageSize)
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	page := int(req.Page)
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * pageSize

	statusFilter := workflowv1.WorkflowTaskStatus_WORKFLOW_TASK_STATUS_UNSPECIFIED
	if req.Status != "" {
		statusFilter = parseTaskStatusStr(req.Status)
	}

	tasks, total, err := s.repo.ListTasksByAssignee(ctx, userID, statusFilter, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("list my tasks: %w", err)
	}
	return &workflowservicev1.GetMyTasksResponse{
		Tasks:      tasks,
		TotalCount: int32(total),
	}, nil
}

func (s *workflowService) CompleteTask(
	ctx context.Context,
	req *workflowservicev1.CompleteWorkflowTaskRequest,
	userID string,
) (*workflowservicev1.CompleteWorkflowTaskResponse, error) {
	if req.TaskId == "" {
		return nil, fmt.Errorf("%w: task_id is required", ErrInvalidArgument)
	}
	if req.Decision == "" {
		return nil, fmt.Errorf("%w: decision is required (APPROVED, REJECTED, RETURNED)", ErrInvalidArgument)
	}

	task, err := s.repo.GetTask(ctx, req.TaskId)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, fmt.Errorf("%w: task %s", ErrNotFound, req.TaskId)
		}
		return nil, err
	}

	if task.Status == workflowv1.WorkflowTaskStatus_WORKFLOW_TASK_STATUS_COMPLETED {
		return nil, fmt.Errorf("%w: task %s is already completed", ErrInvalidTransition, req.TaskId)
	}

	completed, err := s.repo.CompleteTask(ctx, domain.TaskCompleteInput{
		TaskID:      req.TaskId,
		Decision:    req.Decision,
		Comments:    req.Comments,
		CompletedBy: userID,
	})
	if err != nil {
		return nil, fmt.Errorf("complete task: %w", err)
	}

	// Publish task completed event
	s.publisher.Publish(ctx, events.TopicWorkflowTaskCompleted, completed.WorkflowInstanceId, &workfloweventsv1.WorkflowTaskCompletedEvent{
		EventId:            uuid.NewString(),
		TaskId:             completed.Id,
		WorkflowInstanceId: completed.WorkflowInstanceId,
		Decision:           req.Decision,
		CompletedBy:        userID,
		Timestamp:          timestamppb.Now(),
	})

	// Advance or complete the instance
	s.advanceInstance(ctx, completed.WorkflowInstanceId, req.Decision)

	appLogger.Infof("Task %s completed by %s with decision %s", req.TaskId, userID, req.Decision)

	return &workflowservicev1.CompleteWorkflowTaskResponse{
		Message: fmt.Sprintf("Task completed with decision: %s", req.Decision),
	}, nil
}

func (s *workflowService) GetWorkflowHistory(
	ctx context.Context,
	req *workflowservicev1.GetWorkflowHistoryRequest,
) (*workflowservicev1.GetWorkflowHistoryResponse, error) {
	if req.EntityType == "" || req.EntityId == "" {
		return nil, fmt.Errorf("%w: entity_type and entity_id are required", ErrInvalidArgument)
	}
	instances, err := s.repo.ListInstancesByEntity(ctx, req.EntityType, req.EntityId)
	if err != nil {
		return nil, fmt.Errorf("get workflow history: %w", err)
	}
	return &workflowservicev1.GetWorkflowHistoryResponse{WorkflowInstances: instances}, nil
}

// ─── INSTANCE ADVANCEMENT ────────────────────────────────────────────────────

// advanceInstance checks if all tasks in the instance are done and
// moves the instance to the next step or completes it.
func (s *workflowService) advanceInstance(ctx context.Context, instanceID, decision string) {
	tasks, err := s.repo.ListTasksByInstance(ctx, instanceID)
	if err != nil {
		appLogger.Errorf("advanceInstance: failed to list tasks for %s: %v", instanceID, err)
		return
	}

	pendingCount := 0
	allApproved := true
	anyRejected := false

	for _, t := range tasks {
		switch t.Status {
		case workflowv1.WorkflowTaskStatus_WORKFLOW_TASK_STATUS_PENDING,
			workflowv1.WorkflowTaskStatus_WORKFLOW_TASK_STATUS_IN_PROGRESS:
			pendingCount++
		case workflowv1.WorkflowTaskStatus_WORKFLOW_TASK_STATUS_COMPLETED:
			if t.Decision != "APPROVED" {
				allApproved = false
			}
			if t.Decision == "REJECTED" {
				anyRejected = true
			}
		}
	}

	if pendingCount > 0 {
		// Still have pending tasks — find next pending and mark IN_PROGRESS
		return
	}

	// All tasks resolved — determine outcome
	finalStatus := workflowv1.InstanceStatus_INSTANCE_STATUS_COMPLETED
	outcomeStr := "COMPLETED"
	if anyRejected || !allApproved {
		finalStatus = workflowv1.InstanceStatus_INSTANCE_STATUS_FAILED
		outcomeStr = "FAILED"
	}

	if err := s.repo.CompleteInstance(ctx, instanceID, finalStatus); err != nil {
		appLogger.Errorf("advanceInstance: failed to complete instance %s: %v", instanceID, err)
		return
	}

	// Publish WorkflowCompletedEvent
	inst, _ := s.repo.GetInstance(ctx, instanceID)
	entityType, entityID := "", ""
	if inst != nil {
		entityType = inst.EntityType
		entityID = inst.EntityId
	}

	s.publisher.Publish(ctx, events.TopicWorkflowCompleted, instanceID, &workfloweventsv1.WorkflowCompletedEvent{
		EventId:            uuid.NewString(),
		WorkflowInstanceId: instanceID,
		EntityType:         entityType,
		EntityId:           entityID,
		Status:             outcomeStr,
		Timestamp:          timestamppb.Now(),
	})

	appLogger.Infof("Workflow instance %s completed with status %s", instanceID, outcomeStr)
}

// ─── HELPERS ─────────────────────────────────────────────────────────────────

func parseSteps(stepsJSON string) ([]StepDef, error) {
	if stepsJSON == "" || stepsJSON == "{}" || stepsJSON == "[]" {
		return nil, nil
	}
	// Support both array format and object with "steps" key
	var steps []StepDef
	if err := json.Unmarshal([]byte(stepsJSON), &steps); err == nil {
		return steps, nil
	}
	// Try object format: {"steps": [...]}
	var wrapper struct {
		Steps []StepDef `json:"steps"`
	}
	if err := json.Unmarshal([]byte(stepsJSON), &wrapper); err != nil {
		return nil, fmt.Errorf("invalid steps JSON: %w", err)
	}
	return wrapper.Steps, nil
}

func validateSteps(stepsJSON string) error {
	if stepsJSON == "" || stepsJSON == "{}" {
		return nil
	}
	_, err := parseSteps(stepsJSON)
	return err
}

func parseWorkflowTypeStr(s string) workflowv1.WorkflowType {
	switch s {
	case "REVIEW":
		return workflowv1.WorkflowType_WORKFLOW_TYPE_REVIEW
	case "ESCALATION":
		return workflowv1.WorkflowType_WORKFLOW_TYPE_ESCALATION
	case "NOTIFICATION":
		return workflowv1.WorkflowType_WORKFLOW_TYPE_NOTIFICATION
	default:
		return workflowv1.WorkflowType_WORKFLOW_TYPE_APPROVAL
	}
}

func parseTaskTypeStr(s string) workflowv1.WorkflowTaskType {
	switch s {
	case "REVIEW":
		return workflowv1.WorkflowTaskType_WORKFLOW_TASK_TYPE_REVIEW
	case "NOTIFICATION":
		return workflowv1.WorkflowTaskType_WORKFLOW_TASK_TYPE_NOTIFICATION
	case "ACTION":
		return workflowv1.WorkflowTaskType_WORKFLOW_TASK_TYPE_ACTION
	default:
		return workflowv1.WorkflowTaskType_WORKFLOW_TASK_TYPE_APPROVAL
	}
}

func parseTaskStatusStr(s string) workflowv1.WorkflowTaskStatus {
	switch s {
	case "IN_PROGRESS":
		return workflowv1.WorkflowTaskStatus_WORKFLOW_TASK_STATUS_IN_PROGRESS
	case "COMPLETED":
		return workflowv1.WorkflowTaskStatus_WORKFLOW_TASK_STATUS_COMPLETED
	case "SKIPPED":
		return workflowv1.WorkflowTaskStatus_WORKFLOW_TASK_STATUS_SKIPPED
	default:
		return workflowv1.WorkflowTaskStatus_WORKFLOW_TASK_STATUS_PENDING
	}
}

func metaUserID(ctx context.Context) string {
	return grpcmeta.ActorID(ctx, "")
}
