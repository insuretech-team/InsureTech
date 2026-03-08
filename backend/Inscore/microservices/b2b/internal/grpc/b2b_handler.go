package grpc

import (
	"context"
	"errors"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/b2b/internal/domain"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/b2b/internal/service"
	b2bservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/b2b/services/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type B2BHandler struct {
	b2bservicev1.UnimplementedB2BServiceServer
	svc domain.B2BService
}

func NewB2BHandler(svc domain.B2BService) *B2BHandler {
	return &B2BHandler{svc: svc}
}

// ─── ORGANISATION ─────────────────────────────────────────────────────────────

func (h *B2BHandler) CreateOrganisation(ctx context.Context, req *b2bservicev1.CreateOrganisationRequest) (*b2bservicev1.CreateOrganisationResponse, error) {
	resp, err := h.svc.CreateOrganisation(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *B2BHandler) GetOrganisation(ctx context.Context, req *b2bservicev1.GetOrganisationRequest) (*b2bservicev1.GetOrganisationResponse, error) {
	resp, err := h.svc.GetOrganisation(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *B2BHandler) ListOrganisations(ctx context.Context, req *b2bservicev1.ListOrganisationsRequest) (*b2bservicev1.ListOrganisationsResponse, error) {
	resp, err := h.svc.ListOrganisations(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *B2BHandler) UpdateOrganisation(ctx context.Context, req *b2bservicev1.UpdateOrganisationRequest) (*b2bservicev1.UpdateOrganisationResponse, error) {
	resp, err := h.svc.UpdateOrganisation(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *B2BHandler) DeleteOrganisation(ctx context.Context, req *b2bservicev1.DeleteOrganisationRequest) (*b2bservicev1.DeleteOrganisationResponse, error) {
	resp, err := h.svc.DeleteOrganisation(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *B2BHandler) ListOrgMembers(ctx context.Context, req *b2bservicev1.ListOrgMembersRequest) (*b2bservicev1.ListOrgMembersResponse, error) {
	resp, err := h.svc.ListOrgMembers(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *B2BHandler) AddOrgMember(ctx context.Context, req *b2bservicev1.AddOrgMemberRequest) (*b2bservicev1.AddOrgMemberResponse, error) {
	resp, err := h.svc.AddOrgMember(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *B2BHandler) AssignOrgAdmin(ctx context.Context, req *b2bservicev1.AssignOrgAdminRequest) (*b2bservicev1.AssignOrgAdminResponse, error) {
	resp, err := h.svc.AssignOrgAdmin(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *B2BHandler) RemoveOrgMember(ctx context.Context, req *b2bservicev1.RemoveOrgMemberRequest) (*b2bservicev1.RemoveOrgMemberResponse, error) {
	resp, err := h.svc.RemoveOrgMember(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *B2BHandler) ResolveMyOrganisation(ctx context.Context, req *b2bservicev1.ResolveMyOrganisationRequest) (*b2bservicev1.ResolveMyOrganisationResponse, error) {
	resp, err := h.svc.ResolveMyOrganisation(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

// ─── DEPARTMENTS ──────────────────────────────────────────────────────────────

func (h *B2BHandler) ListDepartments(ctx context.Context, req *b2bservicev1.ListDepartmentsRequest) (*b2bservicev1.ListDepartmentsResponse, error) {
	resp, err := h.svc.ListDepartments(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *B2BHandler) GetDepartment(ctx context.Context, req *b2bservicev1.GetDepartmentRequest) (*b2bservicev1.GetDepartmentResponse, error) {
	resp, err := h.svc.GetDepartment(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *B2BHandler) CreateDepartment(ctx context.Context, req *b2bservicev1.CreateDepartmentRequest) (*b2bservicev1.CreateDepartmentResponse, error) {
	resp, err := h.svc.CreateDepartment(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *B2BHandler) UpdateDepartment(ctx context.Context, req *b2bservicev1.UpdateDepartmentRequest) (*b2bservicev1.UpdateDepartmentResponse, error) {
	resp, err := h.svc.UpdateDepartment(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *B2BHandler) DeleteDepartment(ctx context.Context, req *b2bservicev1.DeleteDepartmentRequest) (*b2bservicev1.DeleteDepartmentResponse, error) {
	resp, err := h.svc.DeleteDepartment(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

// ─── EMPLOYEES ────────────────────────────────────────────────────────────────

func (h *B2BHandler) ListEmployees(ctx context.Context, req *b2bservicev1.ListEmployeesRequest) (*b2bservicev1.ListEmployeesResponse, error) {
	resp, err := h.svc.ListEmployees(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *B2BHandler) GetEmployee(ctx context.Context, req *b2bservicev1.GetEmployeeRequest) (*b2bservicev1.GetEmployeeResponse, error) {
	resp, err := h.svc.GetEmployee(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *B2BHandler) CreateEmployee(ctx context.Context, req *b2bservicev1.CreateEmployeeRequest) (*b2bservicev1.CreateEmployeeResponse, error) {
	resp, err := h.svc.CreateEmployee(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *B2BHandler) UpdateEmployee(ctx context.Context, req *b2bservicev1.UpdateEmployeeRequest) (*b2bservicev1.UpdateEmployeeResponse, error) {
	resp, err := h.svc.UpdateEmployee(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *B2BHandler) DeleteEmployee(ctx context.Context, req *b2bservicev1.DeleteEmployeeRequest) (*b2bservicev1.DeleteEmployeeResponse, error) {
	resp, err := h.svc.DeleteEmployee(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

// ─── PURCHASE ORDERS ──────────────────────────────────────────────────────────

func (h *B2BHandler) ListPurchaseOrderCatalog(ctx context.Context, req *b2bservicev1.ListPurchaseOrderCatalogRequest) (*b2bservicev1.ListPurchaseOrderCatalogResponse, error) {
	resp, err := h.svc.ListPurchaseOrderCatalog(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *B2BHandler) ListPurchaseOrders(ctx context.Context, req *b2bservicev1.ListPurchaseOrdersRequest) (*b2bservicev1.ListPurchaseOrdersResponse, error) {
	resp, err := h.svc.ListPurchaseOrders(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *B2BHandler) GetPurchaseOrder(ctx context.Context, req *b2bservicev1.GetPurchaseOrderRequest) (*b2bservicev1.GetPurchaseOrderResponse, error) {
	resp, err := h.svc.GetPurchaseOrder(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (h *B2BHandler) CreatePurchaseOrder(ctx context.Context, req *b2bservicev1.CreatePurchaseOrderRequest) (*b2bservicev1.CreatePurchaseOrderResponse, error) {
	resp, err := h.svc.CreatePurchaseOrder(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

// ─── ERROR MAPPER ─────────────────────────────────────────────────────────────

func mapError(err error) error {
	switch {
	case errors.Is(err, service.ErrInvalidArgument):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, service.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
