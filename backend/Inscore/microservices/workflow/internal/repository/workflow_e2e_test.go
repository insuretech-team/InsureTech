package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/workflow/internal/events"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/workflow/internal/service"
	workflowv1 "github.com/newage-saint/insuretech/gen/go/insuretech/workflow/entity/v1"
	workflowservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/workflow/services/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/structpb"
)

// TestWorkflowE2E_ClaimApprovalFlow tests the full lifecycle of a workflow:
//  1. Register a claim approval definition (template)
//  2. Start a workflow instance for a claim entity
//  3. Complete the first task (technical review) → APPROVED
//  4. Complete the second task (manager approval) → APPROVED
//  5. Verify instance is COMPLETED with correct status
//
// This mirrors the PoliSync "claim.standard-approval" template (2-stage).
func TestWorkflowE2E_ClaimApprovalFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E test")
	}

	ctx := context.Background()
	dbConn := testWorkflowDB(t)
	repo := New(dbConn)

	// Nil publisher — Kafka not needed for E2E test
	publisher := events.NewPublisher(nil)
	svc := service.New(repo, publisher)

	// ── Step 1: Register claim approval definition ─────────────────────────────
	defName := fmt.Sprintf("e2e.claim.standard-approval.%d", time.Now().UnixNano())
	cleanupWorkflowDefinitionByName(ctx, t, dbConn, defName)

	stepsJSON := `[
		{"name":"technical_review","type":"REVIEW","assign_role":"claims_officer","due_hours":72,"order":1},
		{"name":"final_approval","type":"APPROVAL","assign_role":"claims_manager","due_hours":48,"order":2}
	]`

	createResp, err := svc.CreateWorkflowDefinition(ctx, &workflowservicev1.CreateWorkflowDefinitionRequest{
		Name:        defName,
		Description: "E2E test: 2-stage claim approval",
		Type:        "APPROVAL",
		EntityType:  "CLAIM",
		Steps:       mustParseStruct(stepsJSON),
	})
	require.NoError(t, err)
	require.NotEmpty(t, createResp.WorkflowDefinitionId)

	definitionID := createResp.WorkflowDefinitionId
	t.Logf("✅ Definition created: %s", definitionID)

	// Cleanup definition + all instances after test
	t.Cleanup(func() { cleanupWorkflowDefinitionByName(ctx, t, dbConn, defName) })

	// ── Step 2: Create initiator user ─────────────────────────────────────────
	initiatorID := insertTestUser(t, dbConn)
	t.Cleanup(func() { cleanupTestUser(ctx, t, dbConn, initiatorID) })

	// ── Step 3: Start workflow instance for a claim ────────────────────────────
	claimID := uuid.NewString()
	ctxWithUser := metadata.NewIncomingContext(ctx, metadata.Pairs("x-user-id", initiatorID))
	startResp, err := svc.StartWorkflow(ctxWithUser, &workflowservicev1.StartWorkflowRequest{
		WorkflowDefinitionId: definitionID,
		EntityType:           "CLAIM",
		EntityId:             claimID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, startResp.WorkflowInstanceId)

	instanceID := startResp.WorkflowInstanceId
	t.Logf("✅ Workflow started: instance=%s claim=%s", instanceID, claimID)

	// ── Step 4: Verify instance is IN_PROGRESS with tasks created ─────────────
	getInstResp, err := svc.GetWorkflowInstance(ctx, &workflowservicev1.GetWorkflowInstanceRequest{
		WorkflowInstanceId: instanceID,
	})
	require.NoError(t, err)
	require.NotNil(t, getInstResp.WorkflowInstance)
	assert.Equal(t, workflowv1.InstanceStatus_INSTANCE_STATUS_IN_PROGRESS, getInstResp.WorkflowInstance.Status)
	assert.Equal(t, "CLAIM", getInstResp.WorkflowInstance.EntityType)
	assert.Equal(t, claimID, getInstResp.WorkflowInstance.EntityId)
	require.Len(t, getInstResp.Tasks, 2, "expected 2 tasks (technical_review + final_approval)")
	t.Logf("✅ Instance verified IN_PROGRESS with %d tasks", len(getInstResp.Tasks))

	// ── Step 5: Get my tasks for initiator (assigned_to is empty in template)───
	// In real flow tasks would be assigned to role-resolved users
	// For E2E: complete tasks directly via repo
	tasks, err := repo.ListTasksByInstance(ctx, instanceID)
	require.NoError(t, err)
	require.Len(t, tasks, 2)

	// Find tasks by step name (order is non-deterministic when due_date is NULL)
	var task1, task2 *workflowv1.WorkflowTask
	for _, task := range tasks {
		switch task.StepName {
		case "technical_review":
			task1 = task
		case "final_approval":
			task2 = task
		}
	}
	require.NotNil(t, task1, "technical_review task not found")
	require.NotNil(t, task2, "final_approval task not found")
	assert.Equal(t, workflowv1.WorkflowTaskStatus_WORKFLOW_TASK_STATUS_PENDING, task1.Status)
	assert.Equal(t, workflowv1.WorkflowTaskStatus_WORKFLOW_TASK_STATUS_PENDING, task2.Status)
	t.Logf("✅ Tasks verified: [%s=%s, %s=%s]", task1.StepName, task1.Status, task2.StepName, task2.Status)

	// ── Step 6: Complete task 1 (technical_review → APPROVED) ─────────────────
	completeResp1, err := svc.CompleteTask(ctx,
		&workflowservicev1.CompleteWorkflowTaskRequest{
			TaskId:   task1.Id,
			Decision: "APPROVED",
			Comments: "Medical records reviewed, claim looks valid.",
		},
		initiatorID)
	require.NoError(t, err)
	require.NotNil(t, completeResp1)
	t.Logf("✅ Task 1 (%s) completed: APPROVED", task1.StepName)

	// ── Step 7: Verify instance still IN_PROGRESS (task 2 still pending) ───────
	time.Sleep(100 * time.Millisecond) // allow advanceInstance to run
	inst2, err := repo.GetInstance(ctx, instanceID)
	require.NoError(t, err)
	assert.Equal(t, workflowv1.InstanceStatus_INSTANCE_STATUS_IN_PROGRESS, inst2.Status,
		"instance should still be IN_PROGRESS while task 2 is pending")
	t.Logf("✅ Instance still IN_PROGRESS after task 1 approval (task 2 pending)")

	// ── Step 8: Complete task 2 (final_approval → APPROVED) ───────────────────
	completeResp2, err := svc.CompleteTask(ctx,
		&workflowservicev1.CompleteWorkflowTaskRequest{
			TaskId:   task2.Id,
			Decision: "APPROVED",
			Comments: "Amount within policy limits, approved for settlement.",
		},
		initiatorID)
	require.NoError(t, err)
	require.NotNil(t, completeResp2)
	t.Logf("✅ Task 2 (%s) completed: APPROVED", task2.StepName)

	// ── Step 9: Verify instance is COMPLETED ──────────────────────────────────
	time.Sleep(200 * time.Millisecond) // allow advanceInstance to complete
	finalInst, err := repo.GetInstance(ctx, instanceID)
	require.NoError(t, err)
	assert.Equal(t, workflowv1.InstanceStatus_INSTANCE_STATUS_COMPLETED, finalInst.Status,
		"instance should be COMPLETED after all tasks approved")
	assert.NotNil(t, finalInst.CompletedAt, "completed_at should be set")
	t.Logf("✅ Workflow COMPLETED at %v", finalInst.CompletedAt.AsTime())

	// ── Step 10: Verify workflow history ─────────────────────────────────────
	histResp, err := svc.GetWorkflowHistory(ctx, &workflowservicev1.GetWorkflowHistoryRequest{
		EntityType: "CLAIM",
		EntityId:   claimID,
	})
	require.NoError(t, err)
	require.Len(t, histResp.WorkflowInstances, 1)
	assert.Equal(t, instanceID, histResp.WorkflowInstances[0].Id)
	t.Logf("✅ Workflow history: %d instance(s) for claim %s", len(histResp.WorkflowInstances), claimID)

	t.Logf("\n🎉 E2E CLAIM APPROVAL WORKFLOW COMPLETE — all steps passed!")
}

// TestWorkflowE2E_RejectionFlow tests that a rejected task fails the entire workflow.
func TestWorkflowE2E_RejectionFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E test")
	}

	ctx := context.Background()
	dbConn := testWorkflowDB(t)
	repo := New(dbConn)
	publisher := events.NewPublisher(nil)
	svc := service.New(repo, publisher)

	// Setup definition
	defName := fmt.Sprintf("e2e.refund.rejection.%d", time.Now().UnixNano())
	cleanupWorkflowDefinitionByName(ctx, t, dbConn, defName)
	t.Cleanup(func() { cleanupWorkflowDefinitionByName(ctx, t, dbConn, defName) })

	stepsJSON := `[{"name":"review","type":"APPROVAL","assign_role":"finance_officer","due_hours":48,"order":1}]`
	createResp, err := svc.CreateWorkflowDefinition(ctx, &workflowservicev1.CreateWorkflowDefinitionRequest{
		Name:       defName,
		Type:       "APPROVAL",
		EntityType: "REFUND",
		Steps:      mustParseStruct(stepsJSON),
	})
	require.NoError(t, err)

	initiatorID := insertTestUser(t, dbConn)
	t.Cleanup(func() { cleanupTestUser(ctx, t, dbConn, initiatorID) })

	// Start instance — pass user UUID via context so initiated_by is a valid UUID
	ctxWithUser := metadata.NewIncomingContext(ctx, metadata.Pairs("x-user-id", initiatorID))
	startResp, err := svc.StartWorkflow(ctxWithUser, &workflowservicev1.StartWorkflowRequest{
		WorkflowDefinitionId: createResp.WorkflowDefinitionId,
		EntityType:           "REFUND",
		EntityId:             uuid.NewString(),
	})
	require.NoError(t, err)

	// Get the single task
	tasks, err := repo.ListTasksByInstance(ctx, startResp.WorkflowInstanceId)
	require.NoError(t, err)
	require.Len(t, tasks, 1)

	// Complete with REJECTED
	_, err = svc.CompleteTask(ctx,
		&workflowservicev1.CompleteWorkflowTaskRequest{
			TaskId:   tasks[0].Id,
			Decision: "REJECTED",
			Comments: "Refund not within policy terms.",
		},
		initiatorID)
	require.NoError(t, err)

	// Verify instance is FAILED
	time.Sleep(200 * time.Millisecond)
	inst, err := repo.GetInstance(ctx, startResp.WorkflowInstanceId)
	require.NoError(t, err)
	assert.Equal(t, workflowv1.InstanceStatus_INSTANCE_STATUS_FAILED, inst.Status,
		"instance should be FAILED after task rejection")
	t.Logf("✅ Rejection flow: instance correctly FAILED after REJECTED task")
}

// mustParseStruct parses a JSON array of step objects into a proto Struct
// that is passed as Steps in CreateWorkflowDefinitionRequest.
// The service calls req.Steps.MarshalJSON() to get back the JSON,
// so we need a Struct whose MarshalJSON output is the steps JSON array.
// We wrap it as {"steps": [...]} and the service's parseSteps handles both formats.
func mustParseStruct(stepsJSON string) *structpb.Struct {
	// Parse the JSON array into Go interface
	var steps []any
	if err := json.Unmarshal([]byte(stepsJSON), &steps); err != nil {
		panic(fmt.Sprintf("mustParseStruct: invalid JSON: %v", err))
	}
	// Convert to proto Value list
	values := make([]*structpb.Value, len(steps))
	for i, s := range steps {
		v, err := structpb.NewValue(s)
		if err != nil {
			panic(fmt.Sprintf("mustParseStruct: step[%d]: %v", i, err))
		}
		values[i] = v
	}
	return &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"steps": structpb.NewListValue(&structpb.ListValue{Values: values}),
		},
	}
}
