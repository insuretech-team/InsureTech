package grpc

import (
	"context"
	"errors"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/workflow/internal/domain"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/workflow/internal/service"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/grpcmeta"
	workflowservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/workflow/services/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// WorkflowHandler implements workflowservicev1.WorkflowServiceServer.
type WorkflowHandler struct {
	workflowservicev1.UnimplementedWorkflowServiceServer
	svc domain.WorkflowService
}

func NewWorkflowHandler(svc domain.WorkflowService) *WorkflowHandler {
	return &WorkflowHandler{svc: svc}
}

// ─── DEFINITIONS ─────────────────────────────────────────────────────────────

func (h *WorkflowHandler) CreateWorkflowDefinition(ctx context.Context, req *workflowservicev1.CreateWorkflowDefinitionRequest) (*workflowservicev1.CreateWorkflowDefinitionResponse, error) {
	resp, err := h.svc.CreateWorkflowDefinition(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *WorkflowHandler) GetWorkflowDefinition(ctx context.Context, req *workflowservicev1.GetWorkflowDefinitionRequest) (*workflowservicev1.GetWorkflowDefinitionResponse, error) {
	resp, err := h.svc.GetWorkflowDefinition(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

// ─── INSTANCES ───────────────────────────────────────────────────────────────

func (h *WorkflowHandler) StartWorkflow(ctx context.Context, req *workflowservicev1.StartWorkflowRequest) (*workflowservicev1.StartWorkflowResponse, error) {
	resp, err := h.svc.StartWorkflow(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *WorkflowHandler) GetWorkflowInstance(ctx context.Context, req *workflowservicev1.GetWorkflowInstanceRequest) (*workflowservicev1.GetWorkflowInstanceResponse, error) {
	resp, err := h.svc.GetWorkflowInstance(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

// ─── TASKS ───────────────────────────────────────────────────────────────────

func (h *WorkflowHandler) GetMyTasks(ctx context.Context, req *workflowservicev1.GetMyTasksRequest) (*workflowservicev1.GetMyTasksResponse, error) {
	userID := grpcmeta.ActorID(ctx, "")
	resp, err := h.svc.GetMyTasks(ctx, req, userID)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *WorkflowHandler) CompleteTask(ctx context.Context, req *workflowservicev1.CompleteWorkflowTaskRequest) (*workflowservicev1.CompleteWorkflowTaskResponse, error) {
	userID := grpcmeta.ActorID(ctx, "")
	resp, err := h.svc.CompleteTask(ctx, req, userID)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *WorkflowHandler) GetWorkflowHistory(ctx context.Context, req *workflowservicev1.GetWorkflowHistoryRequest) (*workflowservicev1.GetWorkflowHistoryResponse, error) {
	resp, err := h.svc.GetWorkflowHistory(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func mapError(err error) error {
	switch {
	case errors.Is(err, service.ErrInvalidArgument):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, service.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, service.ErrAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, service.ErrInvalidTransition):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, service.ErrFeatureUnavailable):
		return status.Error(codes.Unimplemented, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
