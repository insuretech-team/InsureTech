package consumers

import (
	"context"
	"fmt"
	"testing"

	authzentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/entity/v1"
	authzservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/services/v1"
	authneventsv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/events/v1"
	authzeventsv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/events/v1"
	b2bv1 "github.com/newage-saint/insuretech/gen/go/insuretech/b2b/entity/v1"
	b2beventsv1 "github.com/newage-saint/insuretech/gen/go/insuretech/b2b/events/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// ── Test doubles ──────────────────────────────────────────────────────────────

type fakeAuthZClient struct {
	assignReq  *authzservicev1.AssignRoleRequest
	assignErr  error
	roleName   string // role name to return from ListRoles
}

func newFakeClient(roleName string) *fakeAuthZClient {
	return &fakeAuthZClient{roleName: roleName}
}

func (f *fakeAuthZClient) AssignRole(_ context.Context, req *authzservicev1.AssignRoleRequest) (*authzservicev1.AssignRoleResponse, error) {
	f.assignReq = req
	return &authzservicev1.AssignRoleResponse{}, f.assignErr
}

func (f *fakeAuthZClient) CreatePolicyRule(_ context.Context, _ *authzservicev1.CreatePolicyRuleRequest) (*authzservicev1.CreatePolicyRuleResponse, error) {
	return &authzservicev1.CreatePolicyRuleResponse{}, nil
}

func (f *fakeAuthZClient) ListRoles(_ context.Context, _ *authzservicev1.ListRolesRequest) (*authzservicev1.ListRolesResponse, error) {
	name := f.roleName
	if name == "" {
		name = "b2b_org_admin"
	}
	return &authzservicev1.ListRolesResponse{
		Roles: []*authzentityv1.Role{
			{RoleId: "role-b2b-admin", Name: "b2b_org_admin"},
			{RoleId: "role-partner-user", Name: "partner_user"},
		},
	}, nil
}

func mustMarshalJSON(t *testing.T, m proto.Message) []byte {
	t.Helper()
	b, err := protojson.Marshal(m)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return b
}

// ── B2BAdminAssigned ──────────────────────────────────────────────────────────

func TestHandleB2BAdminAssigned_AcceptsJSONEncodedProtoEvent(t *testing.T) {
	client := newFakeClient("b2b_org_admin")
	consumer := NewEventConsumer(client)

	payload := mustMarshalJSON(t, &b2beventsv1.B2BAdminAssignedEvent{
		OrganisationId: "org-123",
		UserId:         "user-456",
		AssignedBy:     "system-admin",
	})

	if err := consumer.HandleB2BAdminAssigned(context.Background(), payload); err != nil {
		t.Fatalf("HandleB2BAdminAssigned returned error: %v", err)
	}
	if client.assignReq == nil {
		t.Fatal("expected AssignRole to be called")
	}
	if client.assignReq.GetUserId() != "user-456" {
		t.Fatalf("unexpected user id: %s", client.assignReq.GetUserId())
	}
	if client.assignReq.GetRoleId() != "role-b2b-admin" {
		t.Fatalf("unexpected role id: %s", client.assignReq.GetRoleId())
	}
	if client.assignReq.GetDomain() != "b2b:org-123" {
		t.Fatalf("unexpected domain: %s", client.assignReq.GetDomain())
	}
}

func TestHandleB2BAdminAssigned_EmptyPayload_ReturnsError(t *testing.T) {
	consumer := NewEventConsumer(newFakeClient("b2b_org_admin"))
	err := consumer.HandleB2BAdminAssigned(context.Background(), []byte("   "))
	if err == nil {
		t.Fatal("expected error for empty payload")
	}
}

func TestHandleB2BAdminAssigned_InvalidPayload_ReturnsError(t *testing.T) {
	consumer := NewEventConsumer(newFakeClient("b2b_org_admin"))
	err := consumer.HandleB2BAdminAssigned(context.Background(), []byte("{not-valid-json!!!}"))
	if err == nil {
		t.Fatal("expected error for invalid payload")
	}
}

func TestHandleB2BAdminAssigned_DuplicateAssignment_IsIgnored(t *testing.T) {
	client := newFakeClient("b2b_org_admin")
	client.assignErr = fmt.Errorf("already exists")
	consumer := NewEventConsumer(client)

	payload := mustMarshalJSON(t, &b2beventsv1.B2BAdminAssignedEvent{
		OrganisationId: "org-dup",
		UserId:         "user-dup",
	})
	if err := consumer.HandleB2BAdminAssigned(context.Background(), payload); err != nil {
		t.Fatalf("duplicate assignment should be silently ignored, got: %v", err)
	}
}

// ── OrgMemberAdded ────────────────────────────────────────────────────────────

func TestHandleOrgMemberAdded_BusinessAdminMapsToB2BOrgAdminRole(t *testing.T) {
	client := newFakeClient("b2b_org_admin")
	consumer := NewEventConsumer(client)

	payload := mustMarshalJSON(t, &b2beventsv1.OrgMemberAddedEvent{
		OrganisationId: "org-789",
		UserId:         "user-999",
		AddedBy:        "system-admin",
		Role:           b2bv1.OrgMemberRole_ORG_MEMBER_ROLE_BUSINESS_ADMIN,
	})

	if err := consumer.HandleOrgMemberAdded(context.Background(), payload); err != nil {
		t.Fatalf("HandleOrgMemberAdded returned error: %v", err)
	}
	if client.assignReq == nil {
		t.Fatal("expected AssignRole to be called")
	}
	if client.assignReq.GetDomain() != "b2b:org-789" {
		t.Fatalf("unexpected domain: %s", client.assignReq.GetDomain())
	}
	if client.assignReq.GetRoleId() != "role-b2b-admin" {
		t.Fatalf("expected b2b_org_admin role, got role id: %s", client.assignReq.GetRoleId())
	}
}

func TestHandleOrgMemberAdded_HRManagerMapsToPartnerUserRole(t *testing.T) {
	client := newFakeClient("partner_user")
	consumer := NewEventConsumer(client)

	payload := mustMarshalJSON(t, &b2beventsv1.OrgMemberAddedEvent{
		OrganisationId: "org-hr",
		UserId:         "user-hr",
		AddedBy:        "admin",
		Role:           b2bv1.OrgMemberRole_ORG_MEMBER_ROLE_HR_MANAGER,
	})

	if err := consumer.HandleOrgMemberAdded(context.Background(), payload); err != nil {
		t.Fatalf("HandleOrgMemberAdded returned error: %v", err)
	}
	if client.assignReq == nil {
		t.Fatal("expected AssignRole to be called for HR_MANAGER role")
	}
	if client.assignReq.GetDomain() != "b2b:org-hr" {
		t.Fatalf("unexpected domain: %s", client.assignReq.GetDomain())
	}
	if client.assignReq.GetRoleId() != "role-partner-user" {
		t.Fatalf("expected partner_user role, got role id: %s", client.assignReq.GetRoleId())
	}
}

func TestHandleOrgMemberAdded_ViewerMapsToPartnerUserRole(t *testing.T) {
	client := newFakeClient("partner_user")
	consumer := NewEventConsumer(client)

	payload := mustMarshalJSON(t, &b2beventsv1.OrgMemberAddedEvent{
		OrganisationId: "org-viewer",
		UserId:         "user-viewer",
		AddedBy:        "admin",
		Role:           b2bv1.OrgMemberRole_ORG_MEMBER_ROLE_VIEWER,
	})

	if err := consumer.HandleOrgMemberAdded(context.Background(), payload); err != nil {
		t.Fatalf("HandleOrgMemberAdded returned error: %v", err)
	}
	if client.assignReq == nil {
		t.Fatal("expected AssignRole to be called for VIEWER role")
	}
	if client.assignReq.GetRoleId() != "role-partner-user" {
		t.Fatalf("expected partner_user role, got role id: %s", client.assignReq.GetRoleId())
	}
}

func TestHandleOrgMemberAdded_UnknownRole_NoAssignment_NoError(t *testing.T) {
	client := newFakeClient("")
	consumer := NewEventConsumer(client)

	payload := mustMarshalJSON(t, &b2beventsv1.OrgMemberAddedEvent{
		OrganisationId: "org-x",
		UserId:         "user-x",
		Role:           b2bv1.OrgMemberRole_ORG_MEMBER_ROLE_UNSPECIFIED,
	})

	if err := consumer.HandleOrgMemberAdded(context.Background(), payload); err != nil {
		t.Fatalf("unknown role should not return error: %v", err)
	}
	if client.assignReq != nil {
		t.Fatal("AssignRole should NOT be called for unknown role")
	}
}

// ── OrganisationCreated ───────────────────────────────────────────────────────

func TestHandleOrganisationCreated_ProcessesWithoutError(t *testing.T) {
	consumer := NewEventConsumer(newFakeClient(""))
	payload := mustMarshalJSON(t, &b2beventsv1.OrganisationCreatedEvent{
		OrganisationId: "org-new",
		Name:           "Acme Corp",
	})
	if err := consumer.HandleOrganisationCreated(context.Background(), payload); err != nil {
		t.Fatalf("HandleOrganisationCreated returned error: %v", err)
	}
}

// ── OrganisationApproved ──────────────────────────────────────────────────────

func TestHandleOrganisationApproved_ProcessesWithoutError(t *testing.T) {
	consumer := NewEventConsumer(newFakeClient(""))
	payload := mustMarshalJSON(t, &b2beventsv1.OrganisationApprovedEvent{
		OrganisationId: "org-approved",
		ApprovedBy:     "super-admin",
	})
	if err := consumer.HandleOrganisationApproved(context.Background(), payload); err != nil {
		t.Fatalf("HandleOrganisationApproved returned error: %v", err)
	}
}

// ── UserRegistered ────────────────────────────────────────────────────────────

func TestHandleUserRegistered_ProcessesWithoutError(t *testing.T) {
	consumer := NewEventConsumer(newFakeClient(""))
	payload := mustMarshalJSON(t, &authneventsv1.UserRegisteredEvent{
		UserId: "user-new",
		Email:  "test@example.com",
	})
	if err := consumer.HandleUserRegistered(context.Background(), payload); err != nil {
		t.Fatalf("HandleUserRegistered returned error: %v", err)
	}
}

// ── RoleAssigned ──────────────────────────────────────────────────────────────

func TestHandleRoleAssigned_ProcessesWithoutError(t *testing.T) {
	consumer := NewEventConsumer(newFakeClient(""))
	payload := mustMarshalJSON(t, &authzeventsv1.RoleAssignedEvent{
		UserId:   "user-role",
		RoleName: "b2b_org_admin",
		Domain:   "b2b:org-1",
	})
	if err := consumer.HandleRoleAssigned(context.Background(), payload); err != nil {
		t.Fatalf("HandleRoleAssigned returned error: %v", err)
	}
}

// ── unmarshalProtoEvent ───────────────────────────────────────────────────────

func TestUnmarshalProtoEvent_BinaryProto(t *testing.T) {
	original := &b2beventsv1.B2BAdminAssignedEvent{
		OrganisationId: "org-bin",
		UserId:         "user-bin",
	}
	bin, err := proto.Marshal(original)
	if err != nil {
		t.Fatalf("proto.Marshal: %v", err)
	}
	var target b2beventsv1.B2BAdminAssignedEvent
	if err := unmarshalProtoEvent(bin, &target); err != nil {
		t.Fatalf("unmarshalProtoEvent (binary): %v", err)
	}
	if target.GetOrganisationId() != "org-bin" {
		t.Fatalf("unexpected org_id: %s", target.GetOrganisationId())
	}
}

func TestUnmarshalProtoEvent_EmptyPayload_ReturnsError(t *testing.T) {
	var target b2beventsv1.B2BAdminAssignedEvent
	if err := unmarshalProtoEvent([]byte("  "), &target); err == nil {
		t.Fatal("expected error for empty payload")
	}
}

// ── isDuplicateAssignment ─────────────────────────────────────────────────────

func TestIsDuplicateAssignment(t *testing.T) {
	cases := []struct {
		err      error
		expected bool
	}{
		{nil, false},
		{fmt.Errorf("already exists"), true},
		{fmt.Errorf("duplicate key violates unique constraint"), true},
		{fmt.Errorf("uq_user_roles_unique"), true},
		{fmt.Errorf("some other error"), false},
	}
	for _, tc := range cases {
		got := isDuplicateAssignment(tc.err)
		if got != tc.expected {
			t.Errorf("isDuplicateAssignment(%v) = %v, want %v", tc.err, got, tc.expected)
		}
	}
}
