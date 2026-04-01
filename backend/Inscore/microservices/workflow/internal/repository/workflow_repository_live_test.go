package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/workflow/internal/domain"
	workflowv1 "github.com/newage-saint/insuretech/gen/go/insuretech/workflow/entity/v1"
)

// TestWorkflowRepository_LiveDB_DefinitionCRUD tests the full definition lifecycle:
// Create → GetByID → GetByName → ListDefinitions
func TestWorkflowRepository_LiveDB_DefinitionCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	dbConn := testWorkflowDB(t)
	repo := New(dbConn)

	defName := fmt.Sprintf("test.claim-approval.%d", time.Now().UnixNano())

	// Pre-cleanup in case of leftover from a previous failed run
	cleanupWorkflowDefinitionByName(ctx, t, dbConn, defName)
	t.Cleanup(func() { cleanupWorkflowDefinitionByName(ctx, t, dbConn, defName) })

	// ── Create ────────────────────────────────────────────────────────────────
	def, err := repo.CreateDefinition(ctx, domain.DefinitionCreateInput{
		DefinitionID: uuid.NewString(),
		Name:         defName,
		Description:  "Test claim approval workflow",
		Type:         workflowv1.WorkflowType_WORKFLOW_TYPE_APPROVAL,
		EntityType:   "CLAIM",
		Steps:        `[{"name":"review","type":"APPROVAL","assign_role":"claims_officer","due_hours":48,"order":1}]`,
		Version:      1,
		Status:       workflowv1.WorkflowStatus_WORKFLOW_STATUS_ACTIVE,
		CreatedBy:    "test-system",
	})
	require.NoError(t, err)
	require.NotEmpty(t, def.WorkflowDefinitionId)
	require.Equal(t, defName, def.Name)
	require.Equal(t, workflowv1.WorkflowType_WORKFLOW_TYPE_APPROVAL, def.Type)
	require.Equal(t, workflowv1.WorkflowStatus_WORKFLOW_STATUS_ACTIVE, def.Status)
	require.Equal(t, "CLAIM", def.EntityType)

	// ── GetByID ───────────────────────────────────────────────────────────────
	got, err := repo.GetDefinition(ctx, def.WorkflowDefinitionId)
	require.NoError(t, err)
	require.Equal(t, def.WorkflowDefinitionId, got.WorkflowDefinitionId)
	require.Equal(t, defName, got.Name)
	require.Equal(t, workflowv1.WorkflowType_WORKFLOW_TYPE_APPROVAL, got.Type)

	// ── GetByName ─────────────────────────────────────────────────────────────
	byName, err := repo.GetDefinitionByName(ctx, defName)
	require.NoError(t, err)
	require.Equal(t, def.WorkflowDefinitionId, byName.WorkflowDefinitionId)

	// ── ListDefinitions ───────────────────────────────────────────────────────
	defs, total, err := repo.ListDefinitions(ctx, "CLAIM", workflowv1.WorkflowType_WORKFLOW_TYPE_APPROVAL, 10, 0)
	require.NoError(t, err)
	require.Positive(t, total)
	found := false
	for _, d := range defs {
		if d.WorkflowDefinitionId == def.WorkflowDefinitionId {
			found = true
			break
		}
	}
	require.True(t, found, "created definition should appear in list")

	// ── Not found ─────────────────────────────────────────────────────────────
	_, err = repo.GetDefinition(ctx, uuid.NewString())
	require.Error(t, err)
	require.ErrorIs(t, err, domain.ErrNotFound)

	// ── Duplicate name ────────────────────────────────────────────────────────
	_, err = repo.CreateDefinition(ctx, domain.DefinitionCreateInput{
		DefinitionID: uuid.NewString(),
		Name:         defName, // same name → must fail
		EntityType:   "CLAIM",
		Type:         workflowv1.WorkflowType_WORKFLOW_TYPE_APPROVAL,
		Status:       workflowv1.WorkflowStatus_WORKFLOW_STATUS_ACTIVE,
		CreatedBy:    "test-system",
	})
	require.Error(t, err)
	require.ErrorIs(t, err, domain.ErrAlreadyExists)
}

// TestWorkflowRepository_LiveDB_InstanceLifecycle tests the full instance lifecycle:
// Create → GetByID → UpdateStatus → Complete → ListByEntity
func TestWorkflowRepository_LiveDB_InstanceLifecycle(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	dbConn := testWorkflowDB(t)
	repo := New(dbConn)

	// ── Prerequisites: definition + user ──────────────────────────────────────
	defName := fmt.Sprintf("test.instance-lifecycle.%d", time.Now().UnixNano())
	cleanupWorkflowDefinitionByName(ctx, t, dbConn, defName)

	def, err := repo.CreateDefinition(ctx, domain.DefinitionCreateInput{
		DefinitionID: uuid.NewString(),
		Name:         defName,
		EntityType:   "ENDORSEMENT",
		Type:         workflowv1.WorkflowType_WORKFLOW_TYPE_APPROVAL,
		Status:       workflowv1.WorkflowStatus_WORKFLOW_STATUS_ACTIVE,
		CreatedBy:    "test-system",
	})
	require.NoError(t, err)

	userID := insertTestUser(t, dbConn)
	entityID := uuid.NewString()

	t.Cleanup(func() {
		cleanupWorkflowDefinitionByName(ctx, t, dbConn, defName)
		cleanupTestUser(ctx, t, dbConn, userID)
	})

	// ── CreateInstance ────────────────────────────────────────────────────────
	inst, err := repo.CreateInstance(ctx, domain.InstanceCreateInput{
		InstanceID:           uuid.NewString(),
		WorkflowDefinitionID: def.WorkflowDefinitionId,
		EntityType:           "ENDORSEMENT",
		EntityID:             entityID,
		Status:               workflowv1.InstanceStatus_INSTANCE_STATUS_IN_PROGRESS,
		CurrentStep:          "initial_review",
		Context:              `{"amount":"50000"}`,
		InitiatedBy:          userID,
		CorrelationID:        uuid.NewString(),
	})
	require.NoError(t, err)
	require.NotEmpty(t, inst.Id)
	require.Equal(t, "ENDORSEMENT", inst.EntityType)
	require.Equal(t, entityID, inst.EntityId)
	require.Equal(t, workflowv1.InstanceStatus_INSTANCE_STATUS_IN_PROGRESS, inst.Status)
	require.NotNil(t, inst.StartedAt)

	t.Cleanup(func() { cleanupWorkflowInstance(ctx, t, dbConn, inst.Id) })

	// ── GetInstance ───────────────────────────────────────────────────────────
	got, err := repo.GetInstance(ctx, inst.Id)
	require.NoError(t, err)
	require.Equal(t, inst.Id, got.Id)
	require.Equal(t, "initial_review", got.CurrentStep)
	require.Equal(t, userID, got.InitiatedBy)

	// ── UpdateInstanceStatus ──────────────────────────────────────────────────
	err = repo.UpdateInstanceStatus(ctx, inst.Id, workflowv1.InstanceStatus_INSTANCE_STATUS_IN_PROGRESS, "final_approval")
	require.NoError(t, err)

	updated, err := repo.GetInstance(ctx, inst.Id)
	require.NoError(t, err)
	require.Equal(t, "final_approval", updated.CurrentStep)

	// ── ListInstancesByEntity ─────────────────────────────────────────────────
	instances, err := repo.ListInstancesByEntity(ctx, "ENDORSEMENT", entityID)
	require.NoError(t, err)
	require.Len(t, instances, 1)
	require.Equal(t, inst.Id, instances[0].Id)

	// ── CompleteInstance ──────────────────────────────────────────────────────
	err = repo.CompleteInstance(ctx, inst.Id, workflowv1.InstanceStatus_INSTANCE_STATUS_COMPLETED)
	require.NoError(t, err)

	completed, err := repo.GetInstance(ctx, inst.Id)
	require.NoError(t, err)
	require.Equal(t, workflowv1.InstanceStatus_INSTANCE_STATUS_COMPLETED, completed.Status)
	require.NotNil(t, completed.CompletedAt)
}

// TestWorkflowRepository_LiveDB_TaskFlow tests the full task lifecycle:
// Create → GetByID → ListByInstance → ListByAssignee → CompleteTask
func TestWorkflowRepository_LiveDB_TaskFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	dbConn := testWorkflowDB(t)
	repo := New(dbConn)

	// ── Prerequisites ─────────────────────────────────────────────────────────
	defName := fmt.Sprintf("test.task-flow.%d", time.Now().UnixNano())
	cleanupWorkflowDefinitionByName(ctx, t, dbConn, defName)

	def, err := repo.CreateDefinition(ctx, domain.DefinitionCreateInput{
		DefinitionID: uuid.NewString(),
		Name:         defName,
		EntityType:   "CLAIM",
		Type:         workflowv1.WorkflowType_WORKFLOW_TYPE_APPROVAL,
		Status:       workflowv1.WorkflowStatus_WORKFLOW_STATUS_ACTIVE,
		CreatedBy:    "test-system",
	})
	require.NoError(t, err)

	initiatorID := insertTestUser(t, dbConn)
	assigneeID := insertTestUser(t, dbConn)

	t.Cleanup(func() {
		cleanupWorkflowDefinitionByName(ctx, t, dbConn, defName)
		cleanupTestUser(ctx, t, dbConn, initiatorID)
		cleanupTestUser(ctx, t, dbConn, assigneeID)
	})

	inst, err := repo.CreateInstance(ctx, domain.InstanceCreateInput{
		InstanceID:           uuid.NewString(),
		WorkflowDefinitionID: def.WorkflowDefinitionId,
		EntityType:           "CLAIM",
		EntityID:             uuid.NewString(),
		Status:               workflowv1.InstanceStatus_INSTANCE_STATUS_IN_PROGRESS,
		CurrentStep:          "review",
		InitiatedBy:          initiatorID,
	})
	require.NoError(t, err)
	t.Cleanup(func() { cleanupWorkflowInstance(ctx, t, dbConn, inst.Id) })

	// ── CreateTask ────────────────────────────────────────────────────────────
	task, err := repo.CreateTask(ctx, domain.TaskCreateInput{
		TaskID:             uuid.NewString(),
		WorkflowInstanceID: inst.Id,
		StepName:           "review",
		Type:               workflowv1.WorkflowTaskType_WORKFLOW_TASK_TYPE_APPROVAL,
		AssignedTo:         assigneeID,
		Status:             workflowv1.WorkflowTaskStatus_WORKFLOW_TASK_STATUS_PENDING,
		DueHours:           48,
	})
	require.NoError(t, err)
	require.NotEmpty(t, task.Id)
	require.Equal(t, assigneeID, task.AssignedTo)
	require.Equal(t, workflowv1.WorkflowTaskStatus_WORKFLOW_TASK_STATUS_PENDING, task.Status)
	require.NotNil(t, task.DueDate)

	// ── GetTask ───────────────────────────────────────────────────────────────
	got, err := repo.GetTask(ctx, task.Id)
	require.NoError(t, err)
	require.Equal(t, task.Id, got.Id)
	require.Equal(t, "review", got.StepName)

	// ── ListTasksByInstance ───────────────────────────────────────────────────
	tasks, err := repo.ListTasksByInstance(ctx, inst.Id)
	require.NoError(t, err)
	require.Len(t, tasks, 1)
	require.Equal(t, task.Id, tasks[0].Id)

	// ── ListTasksByAssignee ───────────────────────────────────────────────────
	myTasks, total, err := repo.ListTasksByAssignee(
		ctx, assigneeID,
		workflowv1.WorkflowTaskStatus_WORKFLOW_TASK_STATUS_PENDING,
		10, 0,
	)
	require.NoError(t, err)
	require.Positive(t, total)
	found := false
	for _, mt := range myTasks {
		if mt.Id == task.Id {
			found = true
			break
		}
	}
	require.True(t, found, "pending task should appear in assignee's task list")

	// ── CompleteTask ──────────────────────────────────────────────────────────
	completedTask, err := repo.CompleteTask(ctx, domain.TaskCompleteInput{
		TaskID:      task.Id,
		Decision:    "APPROVED",
		Comments:    "Looks good, approved.",
		CompletedBy: assigneeID,
	})
	require.NoError(t, err)
	require.Equal(t, workflowv1.WorkflowTaskStatus_WORKFLOW_TASK_STATUS_COMPLETED, completedTask.Status)
	require.Equal(t, "APPROVED", completedTask.Decision)
	require.Equal(t, "Looks good, approved.", completedTask.Comments)
	require.NotNil(t, completedTask.CompletedAt)

	// ── Verify assignee inbox no longer shows completed task ──────────────────
	afterComplete, _, err := repo.ListTasksByAssignee(
		ctx, assigneeID,
		workflowv1.WorkflowTaskStatus_WORKFLOW_TASK_STATUS_PENDING,
		10, 0,
	)
	require.NoError(t, err)
	for _, mt := range afterComplete {
		require.NotEqual(t, task.Id, mt.Id, "completed task should not appear in pending list")
	}

	// ── Double-complete guard (service layer prevents, but DB state check) ────
	completedAgain, err := repo.GetTask(ctx, task.Id)
	require.NoError(t, err)
	require.Equal(t, workflowv1.WorkflowTaskStatus_WORKFLOW_TASK_STATUS_COMPLETED, completedAgain.Status)
}

// TestWorkflowRepository_LiveDB_NotFound verifies domain.ErrNotFound for all Get methods.
func TestWorkflowRepository_LiveDB_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := context.Background()
	repo := New(testWorkflowDB(t))
	nonExistent := uuid.NewString()

	_, err := repo.GetDefinition(ctx, nonExistent)
	require.ErrorIs(t, err, domain.ErrNotFound)

	_, err = repo.GetInstance(ctx, nonExistent)
	require.ErrorIs(t, err, domain.ErrNotFound)

	_, err = repo.GetTask(ctx, nonExistent)
	require.ErrorIs(t, err, domain.ErrNotFound)
}
