package consumers

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	authneventsv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/events/v1"
	authzentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/entity/v1"
	authzeventsv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/events/v1"
	authzservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/services/v1"
	b2beventsv1 "github.com/newage-saint/insuretech/gen/go/insuretech/b2b/events/v1"
	b2bv1 "github.com/newage-saint/insuretech/gen/go/insuretech/b2b/entity/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// AuthZClient interface for calling AuthZ service
type AuthZClient interface {
	AssignRole(ctx context.Context, req *authzservicev1.AssignRoleRequest) (*authzservicev1.AssignRoleResponse, error)
	CreatePolicyRule(ctx context.Context, req *authzservicev1.CreatePolicyRuleRequest) (*authzservicev1.CreatePolicyRuleResponse, error)
	ListRoles(ctx context.Context, req *authzservicev1.ListRolesRequest) (*authzservicev1.ListRolesResponse, error)
}

// EventConsumer handles B2B events and triggers AuthZ operations
type EventConsumer struct {
	authzClient AuthZClient
}

func NewEventConsumer(authzClient AuthZClient) *EventConsumer {
	return &EventConsumer{
		authzClient: authzClient,
	}
}

// HandleOrganisationCreated processes organisation creation events
// Creates initial Casbin policies for the organisation
func (c *EventConsumer) HandleOrganisationCreated(ctx context.Context, msg []byte) error {
	var evt b2beventsv1.OrganisationCreatedEvent
	if err := unmarshalProtoEvent(msg, &evt); err != nil {
		return fmt.Errorf("unmarshal OrganisationCreatedEvent: %w", err)
	}

	appLogger.Infof("Processing OrganisationCreatedEvent: org_id=%s, name=%s", evt.OrganisationId, evt.Name)
	appLogger.Infof("Organisation authz bootstrap now relies on b2b:root seeded policies and scoped role assignment (org_id=%s)", evt.OrganisationId)
	return nil
}

// HandleB2BAdminAssigned processes B2B admin assignment events
// Assigns the b2b_admin role to the user in AuthZ
func (c *EventConsumer) HandleB2BAdminAssigned(ctx context.Context, msg []byte) error {
	var evt b2beventsv1.B2BAdminAssignedEvent
	if err := unmarshalProtoEvent(msg, &evt); err != nil {
		return fmt.Errorf("unmarshal B2BAdminAssignedEvent: %w", err)
	}

	appLogger.Infof("Processing B2BAdminAssignedEvent: org_id=%s, user_id=%s", evt.OrganisationId, evt.UserId)
	if err := c.assignB2BOrgAdminRole(ctx, evt.UserId, evt.OrganisationId, evt.AssignedBy); err != nil {
		return err
	}

	appLogger.Infof("Assigned b2b_org_admin role to user %s for org %s", evt.UserId, evt.OrganisationId)
	return nil
}

// HandleOrgMemberAdded processes org member addition events
// Assigns appropriate role based on member role
func (c *EventConsumer) HandleOrgMemberAdded(ctx context.Context, msg []byte) error {
	var evt b2beventsv1.OrgMemberAddedEvent
	if err := unmarshalProtoEvent(msg, &evt); err != nil {
		return fmt.Errorf("unmarshal OrgMemberAddedEvent: %w", err)
	}

	appLogger.Infof("Processing OrgMemberAddedEvent: org_id=%s, user_id=%s, role=%s",
		evt.OrganisationId, evt.UserId, evt.Role.String())

	switch evt.Role {
	case b2bv1.OrgMemberRole_ORG_MEMBER_ROLE_BUSINESS_ADMIN:
		// Full org admin — gets the b2b_org_admin Casbin role (all CRUD).
		if err := c.assignB2BOrgAdminRole(ctx, evt.UserId, evt.OrganisationId, evt.AddedBy); err != nil {
			return err
		}
		appLogger.Infof("Assigned b2b_org_admin role to user %s for org %s", evt.UserId, evt.OrganisationId)

	case b2bv1.OrgMemberRole_ORG_MEMBER_ROLE_HR_MANAGER,
		b2bv1.OrgMemberRole_ORG_MEMBER_ROLE_VIEWER:
		// HR managers and viewers get the partner_user Casbin role
		// (read + write b2b data within their org, no admin operations).
		if err := c.assignPartnerUserRole(ctx, evt.UserId, evt.OrganisationId, evt.AddedBy); err != nil {
			return err
		}
		appLogger.Infof("Assigned partner_user role to user %s for org %s (role=%s)",
			evt.UserId, evt.OrganisationId, evt.Role.String())

	default:
		appLogger.Infof("No authz role mapping configured for org member role=%s (user_id=%s org_id=%s)",
			evt.Role.String(), evt.UserId, evt.OrganisationId)
	}
	return nil
}

// assignPartnerUserRole assigns the partner_user Casbin role in domain b2b:{org_id}.
func (c *EventConsumer) assignPartnerUserRole(ctx context.Context, userID, organisationID, assignedBy string) error {
	roleID, err := c.lookupRoleID(ctx, "partner_user")
	if err != nil {
		return fmt.Errorf("lookup partner_user role: %w", err)
	}
	_, err = c.authzClient.AssignRole(ctx, &authzservicev1.AssignRoleRequest{
		UserId:     userID,
		RoleId:     roleID,
		Domain:     fmt.Sprintf("b2b:%s", organisationID),
		AssignedBy: firstNonEmpty(assignedBy, "system"),
	})
	if err != nil {
		if isDuplicateAssignment(err) {
			return nil
		}
		return fmt.Errorf("assign partner_user role: %w", err)
	}
	return nil
}

// HandleOrganisationApproved processes organisation approval events
// Updates organisation status and enables access
func (c *EventConsumer) HandleOrganisationApproved(ctx context.Context, msg []byte) error {
	var evt b2beventsv1.OrganisationApprovedEvent
	if err := unmarshalProtoEvent(msg, &evt); err != nil {
		return fmt.Errorf("unmarshal OrganisationApprovedEvent: %w", err)
	}

	appLogger.Infof("Processing OrganisationApprovedEvent: org_id=%s, approved_by=%s", 
		evt.OrganisationId, evt.ApprovedBy)

	// Additional logic can be added here (e.g., send notifications, update external systems)
	
	return nil
}

// HandleUserRegistered processes user registration events from AuthN
// Can be used to auto-create organisation memberships or send welcome emails
func (c *EventConsumer) HandleUserRegistered(ctx context.Context, msg []byte) error {
	var evt authneventsv1.UserRegisteredEvent
	if err := protojson.Unmarshal(msg, &evt); err != nil {
		return fmt.Errorf("unmarshal UserRegisteredEvent: %w", err)
	}

	appLogger.Infof("Processing UserRegisteredEvent: user_id=%s, email=%s", evt.UserId, evt.Email)

	// B2B service can react to new user registrations
	// For example: check if user email matches pending org invitations
	
	return nil
}

// HandleRoleAssigned processes role assignment events from AuthZ
// Can be used for audit logging or triggering notifications
func (c *EventConsumer) HandleRoleAssigned(ctx context.Context, msg []byte) error {
	var evt authzeventsv1.RoleAssignedEvent
	if err := protojson.Unmarshal(msg, &evt); err != nil {
		return fmt.Errorf("unmarshal RoleAssignedEvent: %w", err)
	}

	appLogger.Infof("Processing RoleAssignedEvent: user_id=%s, role=%s, domain=%s", 
		evt.UserId, evt.RoleName, evt.Domain)

	// B2B service can react to role assignments
	// For example: send notification to user about their new role
	
	return nil
}

func (c *EventConsumer) assignB2BOrgAdminRole(ctx context.Context, userID, organisationID, assignedBy string) error {
	roleID, err := c.lookupRoleID(ctx, "b2b_org_admin")
	if err != nil {
		return fmt.Errorf("lookup b2b_org_admin role: %w", err)
	}

	_, err = c.authzClient.AssignRole(ctx, &authzservicev1.AssignRoleRequest{
		UserId:     userID,
		RoleId:     roleID,
		Domain:     fmt.Sprintf("b2b:%s", organisationID),
		AssignedBy: firstNonEmpty(assignedBy, "system"),
	})
	if err != nil {
		if isDuplicateAssignment(err) {
			return nil
		}
		return fmt.Errorf("assign b2b_org_admin role: %w", err)
	}
	return nil
}

func (c *EventConsumer) lookupRoleID(ctx context.Context, roleName string) (string, error) {
	resp, err := c.authzClient.ListRoles(ctx, &authzservicev1.ListRolesRequest{
		Portal:     authzentityv1.Portal_PORTAL_B2B,
		ActiveOnly: true,
		PageSize:   200,
	})
	if err != nil {
		return "", err
	}

	for _, role := range resp.GetRoles() {
		if role.GetName() == roleName {
			return role.GetRoleId(), nil
		}
	}
	return "", fmt.Errorf("role %s not found", roleName)
}

func unmarshalProtoEvent(msg []byte, target proto.Message) error {
	trimmed := bytes.TrimSpace(msg)
	if len(trimmed) == 0 {
		return fmt.Errorf("empty message payload")
	}
	if trimmed[0] == '{' || trimmed[0] == '[' {
		if err := protojson.Unmarshal(trimmed, target); err == nil {
			return nil
		}
	}
	if err := proto.Unmarshal(trimmed, target); err == nil {
		return nil
	}
	return protojson.Unmarshal(trimmed, target)
}

func isDuplicateAssignment(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "already exists") ||
		strings.Contains(msg, "duplicate") ||
		strings.Contains(msg, "uq_")
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}
	return ""
}
