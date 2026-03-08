package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	authnv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
	authzentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/entity/v1"
	authzv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/services/v1"
	b2bv1 "github.com/newage-saint/insuretech/gen/go/insuretech/b2b/entity/v1"
	b2bservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/b2b/services/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type B2BServiceHandler struct {
	client      b2bservicev1.B2BServiceClient
	authnClient authnv1.AuthServiceClient
	authzClient authzv1.AuthZServiceClient
}

type assignOrgAdminPayload struct {
	MemberID       string `json:"memberId"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	FullName       string `json:"fullName"`
	MobileNumber   string `json:"mobileNumber"`
	DepartmentName string `json:"departmentName"`
	EmployeeID     string `json:"employeeId"`
}

func NewB2BServiceHandler(conn *grpc.ClientConn, authnClient authnv1.AuthServiceClient, authzClient authzv1.AuthZServiceClient) *B2BServiceHandler {
	return &B2BServiceHandler{
		client:      b2bservicev1.NewB2BServiceClient(conn),
		authnClient: authnClient,
		authzClient: authzClient,
	}
}

func parsePurchaseOrderStatusQuery(value string) b2bv1.PurchaseOrderStatus {
	value = strings.TrimSpace(value)
	if value == "" {
		return b2bv1.PurchaseOrderStatus_PURCHASE_ORDER_STATUS_UNSPECIFIED
	}
	if n, err := strconv.Atoi(value); err == nil {
		return b2bv1.PurchaseOrderStatus(n)
	}
	if enumValue, ok := b2bv1.PurchaseOrderStatus_value[strings.ToUpper(value)]; ok {
		return b2bv1.PurchaseOrderStatus(enumValue)
	}
	normalized := "PURCHASE_ORDER_STATUS_" + strings.ToUpper(value)
	if enumValue, ok := b2bv1.PurchaseOrderStatus_value[normalized]; ok {
		return b2bv1.PurchaseOrderStatus(enumValue)
	}
	return b2bv1.PurchaseOrderStatus_PURCHASE_ORDER_STATUS_UNSPECIFIED
}

func (h *B2BServiceHandler) ListPurchaseOrderCatalog(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.ListPurchaseOrderCatalog(ctx, &b2bservicev1.ListPurchaseOrderCatalogRequest{})
	})
}

func (h *B2BServiceHandler) ListPurchaseOrders(w http.ResponseWriter, r *http.Request) {
	if businessID := strings.TrimSpace(r.URL.Query().Get("business_id")); businessID != "" {
		r.Header.Set("X-Business-ID", businessID)
	}
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		req := &b2bservicev1.ListPurchaseOrdersRequest{
			PageToken:  r.URL.Query().Get("page_token"),
			BusinessId: r.URL.Query().Get("business_id"),
			Status:     parsePurchaseOrderStatusQuery(r.URL.Query().Get("status")),
		}
		if q := r.URL.Query().Get("page_size"); q != "" {
			if n, err := strconv.Atoi(q); err == nil && n > 0 {
				req.PageSize = int32(n)
			}
		}
		return h.client.ListPurchaseOrders(ctx, req)
	})
}

func (h *B2BServiceHandler) GetPurchaseOrder(w http.ResponseWriter, r *http.Request) {
	purchaseOrderID := r.PathValue("purchase_order_id")
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.GetPurchaseOrder(ctx, &b2bservicev1.GetPurchaseOrderRequest{
			PurchaseOrderId: purchaseOrderID,
		})
	})
}

func (h *B2BServiceHandler) CreatePurchaseOrder(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req b2bservicev1.CreatePurchaseOrderRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.CreatePurchaseOrder(ctx, &req)
	})
}

func (h *B2BServiceHandler) ListDepartments(w http.ResponseWriter, r *http.Request) {
	if businessID := strings.TrimSpace(r.URL.Query().Get("business_id")); businessID != "" {
		r.Header.Set("X-Business-ID", businessID)
	}
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		req := &b2bservicev1.ListDepartmentsRequest{
			PageToken:  r.URL.Query().Get("page_token"),
			BusinessId: r.URL.Query().Get("business_id"),
		}
		if q := r.URL.Query().Get("page_size"); q != "" {
			if n, err := strconv.Atoi(q); err == nil && n > 0 {
				req.PageSize = int32(n)
			}
		}
		return h.client.ListDepartments(ctx, req)
	})
}

func (h *B2BServiceHandler) GetDepartment(w http.ResponseWriter, r *http.Request) {
	departmentID := r.PathValue("department_id")
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.GetDepartment(ctx, &b2bservicev1.GetDepartmentRequest{
			DepartmentId: departmentID,
		})
	})
}

func (h *B2BServiceHandler) CreateDepartment(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req b2bservicev1.CreateDepartmentRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.CreateDepartment(ctx, &req)
	})
}

func (h *B2BServiceHandler) UpdateDepartment(w http.ResponseWriter, r *http.Request) {
	departmentID := r.PathValue("department_id")
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req b2bservicev1.UpdateDepartmentRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		req.DepartmentId = departmentID
		return h.client.UpdateDepartment(ctx, &req)
	})
}

func (h *B2BServiceHandler) DeleteDepartment(w http.ResponseWriter, r *http.Request) {
	departmentID := r.PathValue("department_id")
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.DeleteDepartment(ctx, &b2bservicev1.DeleteDepartmentRequest{
			DepartmentId: departmentID,
		})
	})
}

func (h *B2BServiceHandler) ListEmployees(w http.ResponseWriter, r *http.Request) {
	if businessID := strings.TrimSpace(r.URL.Query().Get("business_id")); businessID != "" {
		r.Header.Set("X-Business-ID", businessID)
	}
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		req := &b2bservicev1.ListEmployeesRequest{
			PageToken:    r.URL.Query().Get("page_token"),
			DepartmentId: r.URL.Query().Get("department_id"),
			BusinessId:   r.URL.Query().Get("business_id"),
		}
		if q := r.URL.Query().Get("page_size"); q != "" {
			if n, err := strconv.Atoi(q); err == nil && n > 0 {
				req.PageSize = int32(n)
			}
		}
		return h.client.ListEmployees(ctx, req)
	})
}

func (h *B2BServiceHandler) GetEmployee(w http.ResponseWriter, r *http.Request) {
	employeeUUID := r.PathValue("employee_uuid")
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.GetEmployee(ctx, &b2bservicev1.GetEmployeeRequest{
			EmployeeUuid: employeeUUID,
		})
	})
}

func (h *B2BServiceHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req b2bservicev1.CreateEmployeeRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.CreateEmployee(ctx, &req)
	})
}

func (h *B2BServiceHandler) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	employeeUUID := r.PathValue("employee_uuid")
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req b2bservicev1.UpdateEmployeeRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		req.EmployeeUuid = employeeUUID
		return h.client.UpdateEmployee(ctx, &req)
	})
}

func (h *B2BServiceHandler) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	employeeUUID := r.PathValue("employee_uuid")
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.DeleteEmployee(ctx, &b2bservicev1.DeleteEmployeeRequest{
			EmployeeUuid: employeeUUID,
		})
	})
}

// ── Organisations ─────────────────────────────────────────────────────────────

func (h *B2BServiceHandler) ListOrganisations(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		req := &b2bservicev1.ListOrganisationsRequest{
			PageToken: r.URL.Query().Get("page_token"),
			TenantId:  r.URL.Query().Get("tenant_id"),
		}
		if q := r.URL.Query().Get("page_size"); q != "" {
			if n, err := strconv.Atoi(q); err == nil && n > 0 {
				req.PageSize = int32(n)
			}
		}
		return h.client.ListOrganisations(ctx, req)
	})
}

func (h *B2BServiceHandler) GetOrganisation(w http.ResponseWriter, r *http.Request) {
	organisationID := r.PathValue("organisation_id")
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.GetOrganisation(ctx, &b2bservicev1.GetOrganisationRequest{
			OrganisationId: organisationID,
		})
	})
}

func (h *B2BServiceHandler) CreateOrganisation(w http.ResponseWriter, r *http.Request) {
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req b2bservicev1.CreateOrganisationRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		return h.client.CreateOrganisation(ctx, &req)
	})
}

func (h *B2BServiceHandler) UpdateOrganisation(w http.ResponseWriter, r *http.Request) {
	organisationID := r.PathValue("organisation_id")
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req b2bservicev1.UpdateOrganisationRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		req.OrganisationId = organisationID
		return h.client.UpdateOrganisation(ctx, &req)
	})
}

func (h *B2BServiceHandler) DeleteOrganisation(w http.ResponseWriter, r *http.Request) {
	organisationID := r.PathValue("organisation_id")
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.DeleteOrganisation(ctx, &b2bservicev1.DeleteOrganisationRequest{
			OrganisationId: organisationID,
		})
	})
}

func (h *B2BServiceHandler) ListOrgMembers(w http.ResponseWriter, r *http.Request) {
	organisationID := r.PathValue("organisation_id")
	r.Header.Set("X-Business-ID", organisationID)
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.ListOrgMembers(ctx, &b2bservicev1.ListOrgMembersRequest{
			OrganisationId: organisationID,
		})
	})
}

func (h *B2BServiceHandler) ResolveMyOrganisation(w http.ResponseWriter, r *http.Request) {
	userID := strings.TrimSpace(r.Header.Get("X-User-ID"))
	if userID == "" {
		writeJSONError(w, http.StatusUnauthorized, "missing user context")
		return
	}
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.ResolveMyOrganisation(ctx, &b2bservicev1.ResolveMyOrganisationRequest{
			UserId: userID,
		})
	})
}

func (h *B2BServiceHandler) AddOrgMember(w http.ResponseWriter, r *http.Request) {
	organisationID := r.PathValue("organisation_id")
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var req b2bservicev1.AddOrgMemberRequest
		if err := protojson.Unmarshal(body, &req); err != nil {
			return nil, err
		}
		req.OrganisationId = organisationID
		return h.client.AddOrgMember(ctx, &req)
	})
}

func (h *B2BServiceHandler) AssignOrgAdmin(w http.ResponseWriter, r *http.Request) {
	organisationID := r.PathValue("organisation_id")
	callUnary(w, r, func(ctx context.Context, body []byte) (proto.Message, error) {
		var payload assignOrgAdminPayload
		if len(body) > 0 {
			if err := json.Unmarshal(body, &payload); err != nil {
				return nil, status.Error(codes.InvalidArgument, "invalid admin payload")
			}
		}
		if strings.TrimSpace(organisationID) == "" {
			return nil, status.Error(codes.InvalidArgument, "organisation_id is required")
		}

		if memberID := strings.TrimSpace(payload.MemberID); memberID != "" {
			assignResp, err := h.client.AssignOrgAdmin(ctx, &b2bservicev1.AssignOrgAdminRequest{
				OrganisationId: organisationID,
				MemberId:       memberID,
			})
			if err != nil {
				return nil, wrapStepError("assign organisation admin", err)
			}
			if err := h.assignB2BOrgAdminRole(ctx, assignResp.GetMember().GetUserId(), organisationID, strings.TrimSpace(r.Header.Get("X-User-ID"))); err != nil {
				return nil, wrapStepError("assign authz role", err)
			}
			return assignResp, nil
		}

		return h.bootstrapOrgAdmin(ctx, organisationID, payload, strings.TrimSpace(r.Header.Get("X-User-ID")))
	})
}

func (h *B2BServiceHandler) RemoveOrgMember(w http.ResponseWriter, r *http.Request) {
	organisationID := r.PathValue("organisation_id")
	memberID := r.PathValue("member_id")
	callUnary(w, r, func(ctx context.Context, _ []byte) (proto.Message, error) {
		return h.client.RemoveOrgMember(ctx, &b2bservicev1.RemoveOrgMemberRequest{
			OrganisationId: organisationID,
			MemberId:       memberID,
		})
	})
}

func (h *B2BServiceHandler) bootstrapOrgAdmin(
	ctx context.Context,
	organisationID string,
	payload assignOrgAdminPayload,
	assignedBy string,
) (proto.Message, error) {
	email := strings.TrimSpace(payload.Email)
	password := payload.Password
	mobileNumber := strings.TrimSpace(payload.MobileNumber)
	fullName := strings.TrimSpace(payload.FullName)
	if email == "" || password == "" || mobileNumber == "" {
		return nil, status.Error(codes.InvalidArgument, "email, password, and mobileNumber are required")
	}

	registrationResp, err := h.authnClient.RegisterEmailUser(ctx, &authnv1.RegisterEmailUserRequest{
		Email:        email,
		Password:     password,
		FullName:     fullName,
		MobileNumber: mobileNumber,
		UserType:     "B2B_ORG_ADMIN",
		DeviceId:     "b2b-portal-" + organisationID,
	})
	if err != nil {
		return nil, wrapStepError("register admin user", err)
	}
	userID := strings.TrimSpace(registrationResp.GetUserId())
	if userID == "" {
		return nil, status.Error(codes.Internal, "registered user_id missing from auth service response")
	}

	orgCtx := withBusinessContext(ctx, organisationID)

	memberResp, err := h.client.AddOrgMember(orgCtx, &b2bservicev1.AddOrgMemberRequest{
		OrganisationId: organisationID,
		UserId:         userID,
		Role:           b2bv1.OrgMemberRole_ORG_MEMBER_ROLE_BUSINESS_ADMIN,
	})
	if err != nil {
		return nil, wrapStepError("add organisation member", err)
	}

	departmentID, err := h.ensureAdminDepartment(orgCtx, organisationID, payload.DepartmentName)
	if err != nil {
		return nil, wrapStepError("ensure admin department", err)
	}

	_, err = h.client.CreateEmployee(orgCtx, &b2bservicev1.CreateEmployeeRequest{
		Name:          chooseAdminDisplayName(fullName, email),
		EmployeeId:    chooseAdminEmployeeID(payload.EmployeeID, email, organisationID),
		DepartmentId:  departmentID,
		BusinessId:    organisationID,
		Email:         email,
		MobileNumber:  mobileNumber,
		DateOfJoining: time.Now().UTC().Format("2006-01-02"),
	})
	if err != nil {
		return nil, wrapStepError("create admin employee", err)
	}

	if err := h.assignB2BOrgAdminRole(ctx, userID, organisationID, assignedBy); err != nil {
		return nil, wrapStepError("assign authz role", err)
	}

	return &b2bservicev1.AssignOrgAdminResponse{
		Member:  memberResp.GetMember(),
		Message: "B2B admin created successfully",
	}, nil
}

func (h *B2BServiceHandler) ensureAdminDepartment(ctx context.Context, organisationID, requestedName string) (string, error) {
	departmentName := strings.TrimSpace(requestedName)
	if departmentName == "" {
		departmentName = "Administration"
	}

	listResp, err := h.client.ListDepartments(ctx, &b2bservicev1.ListDepartmentsRequest{
		BusinessId: organisationID,
		PageSize:   200,
	})
	if err != nil {
		return "", err
	}
	for _, department := range listResp.GetDepartments() {
		if strings.EqualFold(strings.TrimSpace(department.GetName()), departmentName) {
			return department.GetDepartmentId(), nil
		}
	}

	createResp, err := h.client.CreateDepartment(ctx, &b2bservicev1.CreateDepartmentRequest{
		Name:       departmentName,
		BusinessId: organisationID,
	})
	if err != nil {
		return "", err
	}
	if createResp.GetDepartment() == nil || strings.TrimSpace(createResp.GetDepartment().GetDepartmentId()) == "" {
		return "", status.Error(codes.Internal, "default department creation returned no department_id")
	}
	return createResp.GetDepartment().GetDepartmentId(), nil
}

func (h *B2BServiceHandler) assignB2BOrgAdminRole(ctx context.Context, userID, organisationID, assignedBy string) error {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return status.Error(codes.InvalidArgument, "user_id is required for org admin role assignment")
	}
	if strings.TrimSpace(organisationID) == "" {
		return status.Error(codes.InvalidArgument, "organisation_id is required for org admin role assignment")
	}

	authzCtx := withInternalServiceContext(ctx, "gateway")

	roleID, err := h.lookupB2BOrgAdminRoleID(authzCtx)
	if err != nil {
		return err
	}

	actor := strings.TrimSpace(assignedBy)
	if actor == "" {
		actor = userID
	}

	_, err = h.authzClient.AssignRole(authzCtx, &authzv1.AssignRoleRequest{
		UserId:     userID,
		RoleId:     roleID,
		Domain:     "b2b:" + organisationID,
		AssignedBy: actor,
	})
	if err != nil {
		errText := strings.ToLower(err.Error())
		if strings.Contains(errText, "duplicate") || strings.Contains(errText, "already exists") {
			return nil
		}
		return err
	}
	return nil
}

func (h *B2BServiceHandler) lookupB2BOrgAdminRoleID(ctx context.Context) (string, error) {
	resp, err := h.authzClient.ListRoles(ctx, &authzv1.ListRolesRequest{
		Portal:     authzentityv1.Portal_PORTAL_B2B,
		ActiveOnly: true,
		PageSize:   200,
	})
	if err != nil {
		return "", err
	}
	for _, role := range resp.GetRoles() {
		if role.GetName() == "b2b_org_admin" {
			return role.GetRoleId(), nil
		}
	}
	return "", status.Error(codes.NotFound, "authz role b2b_org_admin not found")
}

func chooseAdminDisplayName(fullName, email string) string {
	if strings.TrimSpace(fullName) != "" {
		return strings.TrimSpace(fullName)
	}
	localPart := email
	if at := strings.Index(localPart, "@"); at > 0 {
		localPart = localPart[:at]
	}
	localPart = strings.ReplaceAll(localPart, ".", " ")
	localPart = strings.ReplaceAll(localPart, "_", " ")
	localPart = strings.ReplaceAll(localPart, "-", " ")
	localPart = strings.TrimSpace(localPart)
	if localPart == "" {
		return "B2B Admin"
	}
	return localPart
}

func wrapStepError(step string, err error) error {
	if err == nil {
		return nil
	}
	if st, ok := status.FromError(err); ok {
		return status.Errorf(st.Code(), "%s: %s", step, st.Message())
	}
	return status.Errorf(codes.Internal, "%s: %v", step, err)
}

func withBusinessContext(ctx context.Context, organisationID string) context.Context {
	organisationID = strings.TrimSpace(organisationID)
	if organisationID == "" {
		return ctx
	}

	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		cloned := md.Copy()
		cloned.Set("x-business-id", organisationID)
		return metadata.NewOutgoingContext(ctx, cloned)
	}

	return metadata.NewOutgoingContext(ctx, metadata.Pairs("x-business-id", organisationID))
}

func withInternalServiceContext(ctx context.Context, serviceName string) context.Context {
	serviceName = strings.TrimSpace(serviceName)
	if serviceName == "" {
		return ctx
	}

	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		cloned := md.Copy()
		cloned.Set("x-internal-service", serviceName)
		return metadata.NewOutgoingContext(ctx, cloned)
	}

	return metadata.NewOutgoingContext(ctx, metadata.Pairs("x-internal-service", serviceName))
}

func chooseAdminEmployeeID(requestedID, email, organisationID string) string {
	if trimmed := strings.TrimSpace(requestedID); trimmed != "" {
		return trimmed
	}
	localPart := email
	if at := strings.Index(localPart, "@"); at > 0 {
		localPart = localPart[:at]
	}
	localPart = strings.ToUpper(localPart)
	localPart = strings.Map(func(r rune) rune {
		if (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			return r
		}
		return -1
	}, localPart)
	if len(localPart) > 6 {
		localPart = localPart[:6]
	}
	if localPart == "" {
		localPart = "ADMIN"
	}
	orgSuffix := strings.ToUpper(strings.ReplaceAll(organisationID, "-", ""))
	if len(orgSuffix) > 6 {
		orgSuffix = orgSuffix[len(orgSuffix)-6:]
	}
	return localPart + "-" + orgSuffix
}
