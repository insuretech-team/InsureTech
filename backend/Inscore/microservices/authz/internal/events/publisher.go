package events

// publisher.go — publishes AuthZ domain events to Kafka topic: authz.events
// Events published:
//   - RoleCreatedEvent, RoleUpdatedEvent, RoleDeletedEvent
//   - RoleAssignedEvent, RoleRemovedEvent
//   - PolicyRuleCreatedEvent, PolicyRuleUpdatedEvent, PolicyRuleDeletedEvent
//   - AccessDeniedEvent (every DENY → SIEM)
//   - AccessGrantedEvent (configurable — high volume)
//   - PortalConfigUpdatedEvent
//   - PolicyCacheInvalidatedEvent

import (
	"context"
	"strconv"

	"github.com/google/uuid"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	authzentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/entity/v1"
	authzeventsv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/events/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const TopicAuthZEvents = "authz.events"

// EventProducer interface for decoupling Kafka/MsgQueue — same pattern as authn.
type EventProducer interface {
	Produce(ctx context.Context, topic string, key string, msg interface{}) error
	Close() error
}

// Publisher publishes AuthZ domain events to Kafka.
type Publisher struct {
	producer EventProducer
}

// NewPublisher creates a new Publisher. producer may be nil (events will be dropped with a log).
func NewPublisher(producer EventProducer) *Publisher {
	return &Publisher{producer: producer}
}

func (p *Publisher) publish(ctx context.Context, topic, key string, msg interface{}) error {
	if p == nil || p.producer == nil {
		appLogger.Infof("authz: Kafka producer not configured — event dropped (topic=%s)", topic)
		return nil
	}
	return p.producer.Produce(ctx, topic, key, msg)
}

// ── Role events ───────────────────────────────────────────────────────────────

// PublishRoleCreated publishes a RoleCreatedEvent after a role is created.
func (p *Publisher) PublishRoleCreated(ctx context.Context, role *authzentityv1.Role) error {
	if role == nil {
		return nil
	}
	evt := &authzeventsv1.RoleCreatedEvent{
		EventId:    uuid.New().String(),
		RoleId:     role.RoleId,
		Name:       role.Name,
		Portal:     role.Portal,
		CreatedBy:  role.CreatedBy,
		OccurredAt: timestamppb.Now(),
	}
	if err := p.publish(ctx, TopicAuthZEvents, role.RoleId, evt); err != nil {
		appLogger.Warnf("authz: failed to publish RoleCreatedEvent (role_id=%s): %v", role.RoleId, err)
	}
	return nil
}

// PublishRoleUpdated publishes a RoleUpdatedEvent after a role is updated.
func (p *Publisher) PublishRoleUpdated(ctx context.Context, role *authzentityv1.Role) error {
	if role == nil {
		return nil
	}
	evt := &authzeventsv1.RoleUpdatedEvent{
		EventId:    uuid.New().String(),
		RoleId:     role.RoleId,
		OccurredAt: timestamppb.Now(),
	}
	if err := p.publish(ctx, TopicAuthZEvents, role.RoleId, evt); err != nil {
		appLogger.Warnf("authz: failed to publish RoleUpdatedEvent (role_id=%s): %v", role.RoleId, err)
	}
	return nil
}

// PublishRoleDeleted publishes a RoleDeletedEvent after a role is soft-deleted.
func (p *Publisher) PublishRoleDeleted(ctx context.Context, roleID, roleName string, portal authzentityv1.Portal, deletedBy string) error {
	evt := &authzeventsv1.RoleDeletedEvent{
		EventId:    uuid.New().String(),
		RoleId:     roleID,
		DeletedBy:  deletedBy,
		OccurredAt: timestamppb.Now(),
	}
	if err := p.publish(ctx, TopicAuthZEvents, roleID, evt); err != nil {
		appLogger.Warnf("authz: failed to publish RoleDeletedEvent (role_id=%s role=%s portal=%v): %v", roleID, roleName, portal, err)
	}
	return nil
}

// ── User-role assignment events ───────────────────────────────────────────────

// PublishRoleAssigned publishes a RoleAssignedEvent after a role is assigned to a user.
func (p *Publisher) PublishRoleAssigned(ctx context.Context, userID, roleID, roleName, domain, assignedBy string) error {
	evt := &authzeventsv1.RoleAssignedEvent{
		EventId:    uuid.New().String(),
		UserId:     userID,
		RoleId:     roleID,
		RoleName:   roleName,
		Domain:     domain,
		AssignedBy: assignedBy,
		OccurredAt: timestamppb.Now(),
	}
	if err := p.publish(ctx, TopicAuthZEvents, userID, evt); err != nil {
		appLogger.Warnf("authz: failed to publish RoleAssignedEvent (user_id=%s role_id=%s): %v", userID, roleID, err)
	}
	return nil
}

// PublishRoleRemoved publishes a RoleRemovedEvent after a role is removed from a user.
func (p *Publisher) PublishRoleRemoved(ctx context.Context, userID, roleID, roleName, domain, removedBy string) error {
	evt := &authzeventsv1.RoleRemovedEvent{
		EventId:    uuid.New().String(),
		UserId:     userID,
		RoleId:     roleID,
		RoleName:   roleName,
		Domain:     domain,
		RemovedBy:  removedBy,
		OccurredAt: timestamppb.Now(),
	}
	if err := p.publish(ctx, TopicAuthZEvents, userID, evt); err != nil {
		appLogger.Warnf("authz: failed to publish RoleRemovedEvent (user_id=%s role_id=%s): %v", userID, roleID, err)
	}
	return nil
}

// ── Policy rule events ────────────────────────────────────────────────────────

// PublishPolicyRuleCreated publishes a PolicyRuleCreatedEvent after a policy rule is created.
func (p *Publisher) PublishPolicyRuleCreated(ctx context.Context, pr *authzentityv1.PolicyRule) error {
	if pr == nil {
		return nil
	}
	evt := &authzeventsv1.PolicyRuleCreatedEvent{
		EventId:    uuid.New().String(),
		PolicyId:   pr.PolicyId,
		Subject:    pr.Subject,
		Domain:     pr.Domain,
		Object:     pr.Object,
		Action:     pr.Action,
		Effect:     pr.Effect,
		CreatedBy:  pr.CreatedBy,
		OccurredAt: timestamppb.Now(),
	}
	if err := p.publish(ctx, TopicAuthZEvents, pr.PolicyId, evt); err != nil {
		appLogger.Warnf("authz: failed to publish PolicyRuleCreatedEvent (policy_id=%s): %v", pr.PolicyId, err)
	}
	return nil
}

// PublishPolicyRuleUpdated publishes a PolicyRuleUpdatedEvent after a policy rule is updated.
func (p *Publisher) PublishPolicyRuleUpdated(ctx context.Context, pr *authzentityv1.PolicyRule) error {
	if pr == nil {
		return nil
	}
	evt := &authzeventsv1.PolicyRuleUpdatedEvent{
		EventId:    uuid.New().String(),
		PolicyId:   pr.PolicyId,
		OccurredAt: timestamppb.Now(),
	}
	if err := p.publish(ctx, TopicAuthZEvents, pr.PolicyId, evt); err != nil {
		appLogger.Warnf("authz: failed to publish PolicyRuleUpdatedEvent (policy_id=%s): %v", pr.PolicyId, err)
	}
	return nil
}

// PublishPolicyRuleDeleted publishes a PolicyRuleDeletedEvent after a policy rule is soft-deleted.
func (p *Publisher) PublishPolicyRuleDeleted(ctx context.Context, policyID, subject, domain string, deletedBy string) error {
	evt := &authzeventsv1.PolicyRuleDeletedEvent{
		EventId:    uuid.New().String(),
		PolicyId:   policyID,
		DeletedBy:  deletedBy,
		OccurredAt: timestamppb.Now(),
	}
	if err := p.publish(ctx, TopicAuthZEvents, policyID, evt); err != nil {
		appLogger.Warnf("authz: failed to publish PolicyRuleDeletedEvent (policy_id=%s subject=%s domain=%s): %v", policyID, subject, domain, err)
	}
	return nil
}

// ── Access decision events ────────────────────────────────────────────────────

// PublishAccessDenied publishes an AccessDeniedEvent for every DENY decision (→ SIEM).
func (p *Publisher) PublishAccessDenied(ctx context.Context, userID, domain, object, action, sessionID, ipAddr string) error {
	subject := "user:" + userID
	evt := &authzeventsv1.AccessDeniedEvent{
		EventId:    uuid.New().String(),
		UserId:     userID,
		SessionId:  sessionID,
		Domain:     domain,
		Subject:    subject,
		Object:     object,
		Action:     action,
		Reason:     "no matching policy (deny-by-default)",
		IpAddress:  ipAddr,
		OccurredAt: timestamppb.Now(),
	}
	if err := p.publish(ctx, TopicAuthZEvents, userID, evt); err != nil {
		appLogger.Warnf("authz: failed to publish AccessDeniedEvent (user_id=%s domain=%s object=%s action=%s): %v", userID, domain, object, action, err)
	}
	return nil
}

// PublishAccessGranted publishes an AccessGrantedEvent (high-volume — enable with care).
func (p *Publisher) PublishAccessGranted(ctx context.Context, userID, domain, object, action string) error {
	subject := "user:" + userID
	evt := &authzeventsv1.AccessGrantedEvent{
		EventId:    uuid.New().String(),
		UserId:     userID,
		Domain:     domain,
		Subject:    subject,
		Object:     object,
		Action:     action,
		OccurredAt: timestamppb.Now(),
	}
	if err := p.publish(ctx, TopicAuthZEvents, userID, evt); err != nil {
		appLogger.Warnf("authz: failed to publish AccessGrantedEvent (user_id=%s domain=%s object=%s action=%s): %v", userID, domain, object, action, err)
	}
	return nil
}

// ── Portal config events ──────────────────────────────────────────────────────

// PublishPortalConfigUpdated publishes a PortalConfigUpdatedEvent after portal config changes.
func (p *Publisher) PublishPortalConfigUpdated(ctx context.Context, portal authzentityv1.Portal, mfaRequired bool, accessTTL, refreshTTL, sessionTTL, idleTimeout int32, updatedBy string) error {
	changedFields := map[string]string{
		"mfa_required":              strconv.FormatBool(mfaRequired),
		"access_token_ttl_seconds":  strconv.Itoa(int(accessTTL)),
		"refresh_token_ttl_seconds": strconv.Itoa(int(refreshTTL)),
		"session_ttl_seconds":       strconv.Itoa(int(sessionTTL)),
		"idle_timeout_seconds":      strconv.Itoa(int(idleTimeout)),
	}
	evt := &authzeventsv1.PortalConfigUpdatedEvent{
		EventId:       uuid.New().String(),
		Portal:        portal,
		UpdatedBy:     updatedBy,
		ChangedFields: changedFields,
		OccurredAt:    timestamppb.Now(),
	}
	key := portal.String()
	if err := p.publish(ctx, TopicAuthZEvents, key, evt); err != nil {
		appLogger.Warnf("authz: failed to publish PortalConfigUpdatedEvent (portal=%v updated_by=%s): %v", portal, updatedBy, err)
	}
	return nil
}

// ── Cache invalidation events ─────────────────────────────────────────────────

// PublishPolicyCacheInvalidated publishes a PolicyCacheInvalidatedEvent after cache is cleared.
func (p *Publisher) PublishPolicyCacheInvalidated(ctx context.Context, domain, invalidatedBy string) error {
	evt := &authzeventsv1.PolicyCacheInvalidatedEvent{
		EventId:    uuid.New().String(),
		Domain:     domain,
		OccurredAt: timestamppb.Now(),
	}
	key := domain
	if key == "" {
		key = "global"
	}
	if err := p.publish(ctx, TopicAuthZEvents, key, evt); err != nil {
		appLogger.Warnf("authz: failed to publish PolicyCacheInvalidatedEvent (domain=%s invalidated_by=%s): %v", domain, invalidatedBy, err)
	}
	return nil
}
