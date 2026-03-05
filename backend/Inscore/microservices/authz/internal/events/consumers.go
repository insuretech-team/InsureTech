package events

// consumers.go — Kafka consumers for AuthZ service.
// Consumes:
//  1. authn.user.registered  → assign default role for new users
//  2. authz.events (PolicyCacheInvalidatedEvent) → reload Casbin enforcer

import (
	"context"
	"errors"
	"strconv"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/authz/internal/domain"
	kafkaconsumer "github.com/newage-saint/insuretech/backend/inscore/pkg/kafka/consumer"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	"go.uber.org/zap"
)

// portalDefaultRoles maps portal name → default role name assigned on registration.
var portalDefaultRoles = map[string]string{
	"b2c":       "customer",
	"agent":     "agent",
	"business":  "admin",
	"b2b":       "partner_user",
	"system":    "support",
	"regulator": "inspector",
}

// UserRegisteredPayload mirrors the authn UserRegisteredEvent Kafka payload for role assignment.
type UserRegisteredPayload struct {
	UserID   string `json:"user_id"`
	Portal   string `json:"portal"`    // portal string e.g. "b2c", "agent"
	TenantID string `json:"tenant_id"` // may be empty → use "root"
}

// PolicyCacheInvalidatedPayload mirrors the authz PolicyCacheInvalidatedEvent payload.
type PolicyCacheInvalidatedPayload struct {
	Domain        string `json:"domain"`
	InvalidatedBy string `json:"invalidated_by"`
}

// NewUserRegisteredHandler returns a HandlerFunc that assigns a default role to new users.
// On receiving authn.user.registered, it assigns the portal-default role (e.g. "customer" for b2c).
func NewUserRegisteredHandler(enforcer domain.EnforcerIface) kafkaconsumer.HandlerFunc {
	return func(ctx context.Context, msg *kafkaconsumer.Message) error {
		var payload UserRegisteredPayload
		if err := msg.Unmarshal(&payload); err != nil {
			return errors.New("user_registered: invalid payload (offset=" + strconv.FormatInt(msg.Offset, 10) + "): " + err.Error())
		}
		if payload.UserID == "" {
			appLogger.Warn("user_registered: missing user_id — skipping", zap.Int64("offset", msg.Offset))
			return nil
		}
		if payload.Portal == "" {
			appLogger.Warn("user_registered: missing portal — skipping", zap.String("user_id", payload.UserID))
			return nil
		}

		defaultRole, ok := portalDefaultRoles[payload.Portal]
		if !ok {
			appLogger.Warn("user_registered: no default role for portal — skipping",
				zap.String("user_id", payload.UserID),
				zap.String("portal", payload.Portal),
			)
			return nil
		}

		tenantID := payload.TenantID
		if tenantID == "" {
			tenantID = "root"
		}

		sub := "user:" + payload.UserID
		role := "role:" + defaultRole
		dom := payload.Portal + ":" + tenantID

		if err := enforcer.AddRoleForUserInDomain(sub, role, dom); err != nil {
			appLogger.Warn("user_registered: failed to assign default role",
				zap.String("user_id", payload.UserID),
				zap.String("role", role),
				zap.String("domain", dom),
				zap.Error(err),
			)
			// Non-fatal: log and continue — a retry or manual fix can be done later.
			return nil
		}

		appLogger.Info("user_registered: default role assigned",
			zap.String("user_id", payload.UserID),
			zap.String("role", role),
			zap.String("domain", dom),
		)
		return nil
	}
}

// NewPolicyCacheInvalidatedHandler returns a HandlerFunc that reloads the Casbin enforcer.
// On receiving authz.events PolicyCacheInvalidatedEvent, it calls enforcer.InvalidateCache().
func NewPolicyCacheInvalidatedHandler(enforcer domain.EnforcerIface) kafkaconsumer.HandlerFunc {
	return func(ctx context.Context, msg *kafkaconsumer.Message) error {
		var payload PolicyCacheInvalidatedPayload
		if err := msg.Unmarshal(&payload); err != nil {
			// Non-fatal: the message may be a different event type on the same topic.
			appLogger.Warn("policy_cache_invalidated: failed to unmarshal payload — skipping",
				zap.Int64("offset", msg.Offset),
				zap.Error(err),
			)
			return nil
		}

		if err := enforcer.InvalidateCache(); err != nil {
			appLogger.Warn("policy_cache_invalidated: enforcer reload failed",
				zap.String("domain", payload.Domain),
				zap.String("invalidated_by", payload.InvalidatedBy),
				zap.Error(err),
			)
			return errors.New("policy_cache_invalidated: reload enforcer: " + err.Error())
		}

		appLogger.Info("policy_cache_invalidated: enforcer cache reloaded",
			zap.String("domain", payload.Domain),
			zap.String("invalidated_by", payload.InvalidatedBy),
		)
		return nil
	}
}

// ── Multi-topic fan-out helper ────────────────────────────────────────────────

// TopicHandlers maps topic names to their HandlerFuncs.
type TopicHandlers map[string]kafkaconsumer.HandlerFunc

// FanOutHandler returns a HandlerFunc that dispatches to the correct handler
// based on msg.Topic. Useful when a single ConsumerGroup subscribes to multiple topics.
func FanOutHandler(handlers TopicHandlers) kafkaconsumer.HandlerFunc {
	return func(ctx context.Context, msg *kafkaconsumer.Message) error {
		h, ok := handlers[msg.Topic]
		if !ok {
			return nil
		}
		return h(ctx, msg)
	}
}
