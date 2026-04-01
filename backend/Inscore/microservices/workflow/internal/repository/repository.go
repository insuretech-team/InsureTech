// Package repository provides PostgreSQL data access for the workflow-engine.
//
// Pattern:
//   - WRITE: map[string]any with enum name strings (proto_enum serializer only applies to struct fields)
//   - READ:  explicit SELECT column list to avoid GORM trying to parse AuditInfo as a relation
//
// Proto structs are the source of truth for column names and types via @inject_tag annotations.
// We use explicit column selects to read back into proto structs safely.
package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	// Import db package to trigger init() registration of proto_enum and proto_timestamp serializers
	_ "github.com/newage-saint/insuretech/backend/inscore/db"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/workflow/internal/domain"
	workflowv1 "github.com/newage-saint/insuretech/gen/go/insuretech/workflow/entity/v1"
	"gorm.io/gorm"
)

// ─── TABLE NAMES ─────────────────────────────────────────────────────────────

const (
	tableDefinitions = "workflow_schema.workflow_definitions"
	tableInstances   = "workflow_schema.workflow_instances"
	tableTasks       = "workflow_schema.workflow_tasks"
)

// ─── COLUMN SELECTS ───────────────────────────────────────────────────────────
// Explicit column lists that map to proto struct fields with GORM tags.
// Excludes audit_info (nested proto struct — causes GORM relation parse error).

const (
	definitionCols = "workflow_definition_id, name, description, type, entity_type, steps, conditions, version, status"
	instanceCols   = "instance_id, workflow_definition_id, entity_type, entity_id, status, current_step, context, initiated_by, correlation_id, started_at, completed_at"
	taskCols       = "task_id, workflow_instance_id, step_name, type, assigned_to, status, decision, comments, due_date, completed_at"
)

// ─── REPOSITORY ───────────────────────────────────────────────────────────────

type workflowRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) domain.WorkflowRepository {
	return &workflowRepository{db: db}
}

// ─── DEFINITIONS ─────────────────────────────────────────────────────────────

func (r *workflowRepository) CreateDefinition(ctx context.Context, in domain.DefinitionCreateInput) (*workflowv1.WorkflowDefinition, error) {
	row := map[string]any{
		"workflow_definition_id": in.DefinitionID,
		"name":                   in.Name,
		"description":            in.Description,
		"type":                   workflowv1.WorkflowType_name[int32(in.Type)],
		"entity_type":            in.EntityType,
		"steps":                  orDefault(in.Steps, "{}"),
		"conditions":             orDefault(in.Conditions, "{}"),
		"version":                orDefaultInt32(in.Version, 1),
		"status":                 workflowv1.WorkflowStatus_name[int32(in.Status)],
		"audit_info":             buildAuditJSON(in.CreatedBy),
	}
	if err := r.db.WithContext(ctx).Table(tableDefinitions).Create(row).Error; err != nil {
		if isUnique(err) {
			return nil, fmt.Errorf("%w: definition '%s'", domain.ErrAlreadyExists, in.Name)
		}
		return nil, fmt.Errorf("create definition: %w", err)
	}
	return r.GetDefinition(ctx, in.DefinitionID)
}

func (r *workflowRepository) GetDefinition(ctx context.Context, id string) (*workflowv1.WorkflowDefinition, error) {
	var row workflowv1.WorkflowDefinition
	err := r.db.WithContext(ctx).Table(tableDefinitions).
		Select(definitionCols).
		Where("workflow_definition_id = ?", id).
		First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: definition %s", domain.ErrNotFound, id)
		}
		return nil, fmt.Errorf("get definition: %w", err)
	}
	return &row, nil
}

func (r *workflowRepository) GetDefinitionByName(ctx context.Context, name string) (*workflowv1.WorkflowDefinition, error) {
	var row workflowv1.WorkflowDefinition
	err := r.db.WithContext(ctx).Table(tableDefinitions).
		Select(definitionCols).
		Where("name = ?", name).
		First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: definition '%s'", domain.ErrNotFound, name)
		}
		return nil, fmt.Errorf("get definition by name: %w", err)
	}
	return &row, nil
}

func (r *workflowRepository) ListDefinitions(ctx context.Context, entityType string, wfType workflowv1.WorkflowType, pageSize, offset int) ([]*workflowv1.WorkflowDefinition, int64, error) {
	q := r.db.WithContext(ctx).Table(tableDefinitions).
		Select(definitionCols).
		Where("status = ?", workflowv1.WorkflowStatus_name[int32(workflowv1.WorkflowStatus_WORKFLOW_STATUS_ACTIVE)])
	if entityType != "" {
		q = q.Where("entity_type = ?", entityType)
	}
	if wfType != workflowv1.WorkflowType_WORKFLOW_TYPE_UNSPECIFIED {
		q = q.Where("type = ?", workflowv1.WorkflowType_name[int32(wfType)])
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count definitions: %w", err)
	}
	var rows []*workflowv1.WorkflowDefinition
	if err := q.Order("name ASC").Limit(pageSize).Offset(offset).Find(&rows).Error; err != nil {
		return nil, 0, fmt.Errorf("list definitions: %w", err)
	}
	return rows, total, nil
}

// ─── INSTANCES ───────────────────────────────────────────────────────────────

func (r *workflowRepository) CreateInstance(ctx context.Context, in domain.InstanceCreateInput) (*workflowv1.WorkflowInstance, error) {
	row := map[string]any{
		"instance_id":            in.InstanceID,
		"workflow_definition_id": in.WorkflowDefinitionID,
		"entity_type":            in.EntityType,
		"entity_id":              in.EntityID,
		"status":                 workflowv1.InstanceStatus_name[int32(in.Status)],
		"current_step":           in.CurrentStep,
		"context":                orDefault(in.Context, "{}"),
		"initiated_by":           in.InitiatedBy,
		"correlation_id":         nullableUUID(in.CorrelationID),
		"started_at":             time.Now().UTC(),
		"audit_info":             buildAuditJSON(in.InitiatedBy),
	}
	if err := r.db.WithContext(ctx).Table(tableInstances).Create(row).Error; err != nil {
		return nil, fmt.Errorf("create instance: %w", err)
	}
	return r.GetInstance(ctx, in.InstanceID)
}

func (r *workflowRepository) GetInstance(ctx context.Context, instanceID string) (*workflowv1.WorkflowInstance, error) {
	var row workflowv1.WorkflowInstance
	err := r.db.WithContext(ctx).Table(tableInstances).
		Select(instanceCols).
		Where("instance_id = ?", instanceID).
		First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: instance %s", domain.ErrNotFound, instanceID)
		}
		return nil, fmt.Errorf("get instance: %w", err)
	}
	return &row, nil
}

func (r *workflowRepository) ListInstancesByEntity(ctx context.Context, entityType, entityID string) ([]*workflowv1.WorkflowInstance, error) {
	var rows []*workflowv1.WorkflowInstance
	err := r.db.WithContext(ctx).Table(tableInstances).
		Select(instanceCols).
		Where("entity_type = ? AND entity_id = ?", entityType, entityID).
		Order("started_at DESC").
		Find(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("list instances by entity: %w", err)
	}
	return rows, nil
}

func (r *workflowRepository) UpdateInstanceStatus(ctx context.Context, instanceID string, status workflowv1.InstanceStatus, currentStep string) error {
	result := r.db.WithContext(ctx).Table(tableInstances).
		Where("instance_id = ?", instanceID).
		Updates(map[string]any{
			"status":       workflowv1.InstanceStatus_name[int32(status)],
			"current_step": currentStep,
		})
	if result.Error != nil {
		return fmt.Errorf("update instance status: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("%w: instance %s", domain.ErrNotFound, instanceID)
	}
	return nil
}

func (r *workflowRepository) CompleteInstance(ctx context.Context, instanceID string, status workflowv1.InstanceStatus) error {
	result := r.db.WithContext(ctx).Table(tableInstances).
		Where("instance_id = ?", instanceID).
		Updates(map[string]any{
			"status":       workflowv1.InstanceStatus_name[int32(status)],
			"completed_at": time.Now().UTC(),
		})
	if result.Error != nil {
		return fmt.Errorf("complete instance: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("%w: instance %s", domain.ErrNotFound, instanceID)
	}
	return nil
}

// ─── TASKS ───────────────────────────────────────────────────────────────────

func (r *workflowRepository) CreateTask(ctx context.Context, in domain.TaskCreateInput) (*workflowv1.WorkflowTask, error) {
	row := map[string]any{
		"task_id":              in.TaskID,
		"workflow_instance_id": in.WorkflowInstanceID,
		"step_name":            in.StepName,
		"type":                 workflowv1.WorkflowTaskType_name[int32(in.Type)],
		"assigned_to":          nullableUUID(in.AssignedTo),
		"status":               workflowv1.WorkflowTaskStatus_name[int32(in.Status)],
		"audit_info":           buildAuditJSON(in.AssignedTo),
	}
	if in.DueHours > 0 {
		row["due_date"] = time.Now().UTC().Add(time.Duration(in.DueHours) * time.Hour)
	}
	if err := r.db.WithContext(ctx).Table(tableTasks).Create(row).Error; err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}
	return r.GetTask(ctx, in.TaskID)
}

func (r *workflowRepository) GetTask(ctx context.Context, taskID string) (*workflowv1.WorkflowTask, error) {
	var row workflowv1.WorkflowTask
	err := r.db.WithContext(ctx).Table(tableTasks).
		Select(taskCols).
		Where("task_id = ?", taskID).
		First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: task %s", domain.ErrNotFound, taskID)
		}
		return nil, fmt.Errorf("get task: %w", err)
	}
	return &row, nil
}

func (r *workflowRepository) ListTasksByInstance(ctx context.Context, instanceID string) ([]*workflowv1.WorkflowTask, error) {
	var rows []*workflowv1.WorkflowTask
	err := r.db.WithContext(ctx).Table(tableTasks).
		Select(taskCols).
		Where("workflow_instance_id = ?", instanceID).
		Order("due_date ASC NULLS LAST").
		Find(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("list tasks by instance: %w", err)
	}
	return rows, nil
}

func (r *workflowRepository) ListTasksByAssignee(ctx context.Context, assignedTo string, status workflowv1.WorkflowTaskStatus, pageSize, offset int) ([]*workflowv1.WorkflowTask, int64, error) {
	q := r.db.WithContext(ctx).Table(tableTasks).
		Select(taskCols).
		Where("assigned_to = ?", assignedTo)
	if status != workflowv1.WorkflowTaskStatus_WORKFLOW_TASK_STATUS_UNSPECIFIED {
		q = q.Where("status = ?", workflowv1.WorkflowTaskStatus_name[int32(status)])
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count tasks: %w", err)
	}
	var rows []*workflowv1.WorkflowTask
	if err := q.Order("due_date ASC NULLS LAST").Limit(pageSize).Offset(offset).Find(&rows).Error; err != nil {
		return nil, 0, fmt.Errorf("list tasks by assignee: %w", err)
	}
	return rows, total, nil
}

func (r *workflowRepository) CompleteTask(ctx context.Context, in domain.TaskCompleteInput) (*workflowv1.WorkflowTask, error) {
	result := r.db.WithContext(ctx).Table(tableTasks).
		Where("task_id = ?", in.TaskID).
		Updates(map[string]any{
			"status":       workflowv1.WorkflowTaskStatus_name[int32(workflowv1.WorkflowTaskStatus_WORKFLOW_TASK_STATUS_COMPLETED)],
			"decision":     in.Decision,
			"comments":     in.Comments,
			"completed_at": time.Now().UTC(),
		})
	if result.Error != nil {
		return nil, fmt.Errorf("complete task: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("%w: task %s", domain.ErrNotFound, in.TaskID)
	}
	return r.GetTask(ctx, in.TaskID)
}

// ─── HELPERS ─────────────────────────────────────────────────────────────────

// nullableUUID returns nil for empty strings so UUID columns get NULL instead of invalid input.
func nullableUUID(s string) any {
	if s == "" {
		return nil
	}
	return s
}

func buildAuditJSON(userID string) string {
	return fmt.Sprintf(`{"created_by":%q,"created_at":%q}`, userID, time.Now().UTC().Format(time.RFC3339))
}

func orDefault(s, def string) string {
	if s == "" {
		return def
	}
	return s
}

func orDefaultInt32(v, def int32) int32 {
	if v == 0 {
		return def
	}
	return v
}

func isUnique(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	for _, s := range []string{"unique", "duplicate", "23505"} {
		for i := 0; i <= len(msg)-len(s); i++ {
			if msg[i:i+len(s)] == s {
				return true
			}
		}
	}
	return false
}
