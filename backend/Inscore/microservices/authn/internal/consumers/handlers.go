// Package consumers wires Kafka consumer handlers for authn domain events.
// Each handler processes a specific topic with idiomatic error handling and
// structured logging.
package consumers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/email"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/repository"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/sms"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/kafka/consumer"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	"go.uber.org/zap"
)

// ─── authn.otp.sms_dlr ────────────────────────────────────────────────────────

// DLRPayload mirrors the SSLWireless DLR webhook JSON body.
type DLRPayload struct {
	MessageID string `json:"message_id"`
	Status    string `json:"status"`
	ErrorCode string `json:"error_code,omitempty"`
}

// NewSMSDLRHandler returns a HandlerFunc that persists SMS delivery reports.
func NewSMSDLRHandler(otpRepo *repository.OTPRepository) consumer.HandlerFunc {
	return func(ctx context.Context, msg *consumer.Message) error {
		var dlr DLRPayload
		if err := msg.Unmarshal(&dlr); err != nil {
			logger.Errorf("sms_dlr: invalid payload (offset=%d): %v", msg.Offset, err)
			return errors.New("sms_dlr: invalid payload (offset=%d)")
		}
		if dlr.MessageID == "" {
			appLogger.Warn("sms_dlr: missing message_id — skipping", zap.Int64("offset", msg.Offset))
			return nil // non-retryable; skip rather than DLQ
		}
		if err := otpRepo.UpdateDLRStatus(ctx, dlr.MessageID, dlr.Status, dlr.ErrorCode); err != nil {
			logger.Errorf("sms_dlr: UpdateDLRStatus failed (msg_id=%s): %v", dlr.MessageID, err)
			return errors.New("sms_dlr: UpdateDLRStatus failed (msg_id=%s)")
		}
		appLogger.Info("sms_dlr: delivery report persisted",
			zap.String("message_id", dlr.MessageID),
			zap.String("status", dlr.Status),
		)
		return nil
	}
}

// ─── authn.account.locked ─────────────────────────────────────────────────────

// AccountLockedPayload is the event payload for account lock events.
type AccountLockedPayload struct {
	UserID   string `json:"user_id"`
	Reason   string `json:"reason"`
	LockedAt string `json:"locked_at"`
}

// NewAccountLockedHandler returns a HandlerFunc that sends an SMS alert on account lock.
func NewAccountLockedHandler(userRepo *repository.UserRepository, smsClient *sms.SSLWirelessClient) consumer.HandlerFunc {
	return func(ctx context.Context, msg *consumer.Message) error {
		var evt AccountLockedPayload
		if err := msg.Unmarshal(&evt); err != nil {
			logger.Errorf("account_locked: invalid payload (offset=%d): %v", msg.Offset, err)
			return errors.New("account_locked: invalid payload (offset=%d)")
		}
		if evt.UserID == "" {
			appLogger.Warn("account_locked: missing user_id — skipping", zap.Int64("offset", msg.Offset))
			return nil
		}
		appLogger.Info("account_locked: received",
			zap.String("user_id", evt.UserID),
			zap.String("reason", evt.Reason),
			zap.String("locked_at", evt.LockedAt),
		)
		user, err := userRepo.GetByID(ctx, evt.UserID)
		if err != nil {
			appLogger.Warn("account_locked: user lookup failed", zap.String("user_id", evt.UserID), zap.Error(err))
			return nil // non-retryable: user may have been deleted
		}
		smsText := fmt.Sprintf(
			"Your InsureTech account has been temporarily locked due to multiple failed login attempts (reason: %s). It will auto-unlock at %s. If this was not you, contact support immediately.",
			evt.Reason, evt.LockedAt,
		)
		if _, err := smsClient.SendSMS(ctx, &sms.SendSMSRequest{MSISDN: user.MobileNumber, Message: smsText}); err != nil {
			appLogger.Warn("account_locked: SMS send failed", zap.String("user_id", evt.UserID), zap.Error(err))
		}
		return nil
	}
}

// ─── authn.user.registered ────────────────────────────────────────────────────

// UserRegisteredPayload is the event payload for new user registration.
type UserRegisteredPayload struct {
	UserID       string `json:"user_id"`
	MobileNumber string `json:"mobile_number"`
	Email        string `json:"email"`
	DeviceType   string `json:"device_type"`
}

// NewUserRegisteredHandler returns a HandlerFunc that sends welcome email + SMS on registration.
func NewUserRegisteredHandler(emailClient *email.Client, smsClient *sms.SSLWirelessClient) consumer.HandlerFunc {
	return func(ctx context.Context, msg *consumer.Message) error {
		var evt UserRegisteredPayload
		if err := msg.Unmarshal(&evt); err != nil {
			logger.Errorf("user_registered: invalid payload (offset=%d): %v", msg.Offset, err)
			return errors.New("user_registered: invalid payload (offset=%d)")
		}
		if evt.UserID == "" {
			appLogger.Warn("user_registered: missing user_id — skipping", zap.Int64("offset", msg.Offset))
			return nil
		}
		appLogger.Info("user_registered: received",
			zap.String("user_id", evt.UserID),
			zap.String("mobile_number", evt.MobileNumber),
		)
		// Send welcome SMS
		welcomeSMS := "Welcome to InsureTech! Your account is ready. Download our app to browse insurance products and protect your family."
		if _, err := smsClient.SendSMS(ctx, &sms.SendSMSRequest{MSISDN: evt.MobileNumber, Message: welcomeSMS}); err != nil {
			appLogger.Warn("user_registered: welcome SMS failed", zap.String("user_id", evt.UserID), zap.Error(err))
		}
		// Send welcome email if available
		if evt.Email != "" {
			if _, err := emailClient.SendOTP(&email.SendOTPRequest{
				To:        evt.Email,
				OTPCode:   "",
				Purpose:   "welcome",
				ExpiryMin: 0,
			}); err != nil {
				appLogger.Warn("user_registered: welcome email failed", zap.String("user_id", evt.UserID), zap.Error(err))
			}
		}
		return nil
	}
}

// ─── authn.password.changed ───────────────────────────────────────────────────

// PasswordChangedPayload is the event payload for password change events.
type PasswordChangedPayload struct {
	UserID    string `json:"user_id"`
	IP        string `json:"ip_address"`
	ChangedBy string `json:"changed_by"`
}

// NewPasswordChangedHandler returns a HandlerFunc that notifies users on password change.
func NewPasswordChangedHandler(userRepo *repository.UserRepository, smsClient *sms.SSLWirelessClient) consumer.HandlerFunc {
	return func(ctx context.Context, msg *consumer.Message) error {
		var evt PasswordChangedPayload
		if err := msg.Unmarshal(&evt); err != nil {
			logger.Errorf("password_changed: invalid payload (offset=%d): %v", msg.Offset, err)
			return errors.New("password_changed: invalid payload (offset=%d)")
		}
		if evt.UserID == "" {
			appLogger.Warn("password_changed: missing user_id — skipping", zap.Int64("offset", msg.Offset))
			return nil
		}
		user, err := userRepo.GetByID(ctx, evt.UserID)
		if err != nil {
			appLogger.Warn("password_changed: user lookup failed", zap.String("user_id", evt.UserID), zap.Error(err))
			return nil
		}
		smsMsg := fmt.Sprintf(
			"Your InsureTech password was changed from IP %s. If this was not you, contact support immediately and secure your account.",
			evt.IP,
		)
		if _, err := smsClient.SendSMS(ctx, &sms.SendSMSRequest{MSISDN: user.MobileNumber, Message: smsMsg}); err != nil {
			appLogger.Warn("password_changed: SMS send failed", zap.String("user_id", evt.UserID), zap.Error(err))
		}
		return nil
	}
}

// ─── authn.password.reset_requested ──────────────────────────────────────────

// PasswordResetRequestedPayload is the event payload for password reset requests.
type PasswordResetRequestedPayload struct {
	UserID       string `json:"user_id"`
	MobileNumber string `json:"mobile_number"`
	IP           string `json:"ip_address"`
}

// NewPasswordResetRequestedHandler notifies user when a reset is requested.
func NewPasswordResetRequestedHandler(userRepo *repository.UserRepository, smsClient *sms.SSLWirelessClient) consumer.HandlerFunc {
	return func(ctx context.Context, msg *consumer.Message) error {
		var evt PasswordResetRequestedPayload
		if err := msg.Unmarshal(&evt); err != nil {
			logger.Errorf("password_reset_requested: invalid payload (offset=%d): %v", msg.Offset, err)
			return errors.New("password_reset_requested: invalid payload (offset=%d)")
		}
		if evt.MobileNumber == "" && evt.UserID == "" {
			appLogger.Warn("password_reset_requested: missing identifiers — skipping", zap.Int64("offset", msg.Offset))
			return nil
		}
		mobile := evt.MobileNumber
		if mobile == "" {
			user, err := userRepo.GetByID(ctx, evt.UserID)
			if err != nil {
				appLogger.Warn("password_reset_requested: user lookup failed", zap.String("user_id", evt.UserID), zap.Error(err))
				return nil
			}
			mobile = user.MobileNumber
		}
		smsMsg := "A password reset was requested for your InsureTech account. If this was not you, please contact support immediately."
		if _, err := smsClient.SendSMS(ctx, &sms.SendSMSRequest{MSISDN: mobile, Message: smsMsg}); err != nil {
			appLogger.Warn("password_reset_requested: SMS send failed", zap.String("mobile", mobile), zap.Error(err))
		}
		return nil
	}
}

// ─── authn.session.revoked ────────────────────────────────────────────────────

// SessionRevokedPayload is the event payload for session revocation events.
type SessionRevokedPayload struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
	Reason    string `json:"reason"`
}

// NewSessionRevokedAllHandler notifies user when all sessions are revoked.
func NewSessionRevokedAllHandler(userRepo *repository.UserRepository, smsClient *sms.SSLWirelessClient) consumer.HandlerFunc {
	return func(ctx context.Context, msg *consumer.Message) error {
		var evt SessionRevokedPayload
		if err := msg.Unmarshal(&evt); err != nil {
			logger.Errorf("session_revoked: invalid payload (offset=%d): %v", msg.Offset, err)
			return errors.New("session_revoked: invalid payload (offset=%d)")
		}
		// Only send SMS when all sessions are revoked (no specific session_id = bulk revoke)
		if evt.UserID == "" || evt.SessionID != "" {
			return nil
		}
		user, err := userRepo.GetByID(ctx, evt.UserID)
		if err != nil {
			appLogger.Warn("session_revoked: user lookup failed", zap.String("user_id", evt.UserID), zap.Error(err))
			return nil
		}
		smsMsg := "All your InsureTech sessions have been logged out for security. If this was not you, contact support immediately."
		if _, err := smsClient.SendSMS(ctx, &sms.SendSMSRequest{MSISDN: user.MobileNumber, Message: smsMsg}); err != nil {
			appLogger.Warn("session_revoked: SMS send failed", zap.String("user_id", evt.UserID), zap.Error(err))
		}
		return nil
	}
}

// ─── Multi-topic fan-out helper ───────────────────────────────────────────────

// TopicHandlers maps topic names to their HandlerFuncs.
type TopicHandlers map[string]consumer.HandlerFunc

// FanOutHandler returns a HandlerFunc that dispatches to the correct handler
// based on msg.Topic. Useful when a single ConsumerGroup subscribes to multiple topics.
func FanOutHandler(handlers TopicHandlers) consumer.HandlerFunc {
	return func(ctx context.Context, msg *consumer.Message) error {
		h, ok := handlers[msg.Topic]
		if !ok {
			raw, _ := json.Marshal(msg)
			appLogger.Warn("fan_out: no handler for topic",
				zap.String("topic", msg.Topic),
				zap.ByteString("message", raw),
			)
			return nil // skip unknown topics
		}
		return h(ctx, msg)
	}
}

// ─── authz.events (PortalConfigUpdated) ──────────────────────────────────────

// PortalConfigUpdatedPayload mirrors the authz PortalConfigUpdatedEvent Kafka payload.
type PortalConfigUpdatedPayload struct {
	Portal             string `json:"portal"`
	MfaRequired        bool   `json:"mfa_required"`
	AccessTokenTTL     int32  `json:"access_token_ttl_seconds"`
	RefreshTokenTTL    int32  `json:"refresh_token_ttl_seconds"`
	SessionTTL         int32  `json:"session_ttl_seconds"`
	IdleTimeoutSeconds int32  `json:"idle_timeout_seconds"`
	UpdatedBy          string `json:"updated_by"`
	UpdatedAt          string `json:"updated_at"`
}

// PortalConfigCache is a simple thread-safe in-process cache for portal configs.
// Authn service reads this cache on every login to check MFA + session requirements
// without making a synchronous gRPC call to AuthZ on each request.
type PortalConfigCache struct {
	mu      sync.RWMutex
	configs map[string]*PortalConfigUpdatedPayload // keyed by portal name
}

var GlobalPortalConfigCache = &PortalConfigCache{
	configs: make(map[string]*PortalConfigUpdatedPayload),
}

// Set stores a portal config in the cache.
func (c *PortalConfigCache) Set(portal string, cfg *PortalConfigUpdatedPayload) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.configs[portal] = cfg
}

// Get retrieves a portal config from the cache. Returns nil if not present.
func (c *PortalConfigCache) Get(portal string) *PortalConfigUpdatedPayload {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.configs[portal]
}

// Invalidate removes a portal config from the cache (forces re-read from AuthZ on next login).
func (c *PortalConfigCache) Invalidate(portal string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.configs, portal)
}

// InvalidateAll clears the entire portal config cache.
func (c *PortalConfigCache) InvalidateAll() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.configs = make(map[string]*PortalConfigUpdatedPayload)
}

// NewPortalConfigUpdatedHandler returns a HandlerFunc that receives PortalConfigUpdatedEvent
// from the authz.events Kafka topic and updates the local in-memory portal config cache.
// This allows authn to enforce MFA + session policies from AuthZ without synchronous gRPC calls.
func NewPortalConfigUpdatedHandler() consumer.HandlerFunc {
	return func(ctx context.Context, msg *consumer.Message) error {
		var evt PortalConfigUpdatedPayload
		if err := msg.Unmarshal(&evt); err != nil {
			logger.Errorf("portal_config_updated: invalid payload (offset=%d): %v", msg.Offset, err)
			return errors.New("portal_config_updated: invalid payload (offset=%d)")
		}
		if evt.Portal == "" {
			appLogger.Warn("portal_config_updated: missing portal — skipping", zap.Int64("offset", msg.Offset))
			return nil
		}
		GlobalPortalConfigCache.Set(evt.Portal, &evt)
		appLogger.Info("portal_config_updated: cache refreshed",
			zap.String("portal", evt.Portal),
			zap.Bool("mfa_required", evt.MfaRequired),
			zap.String("updated_by", evt.UpdatedBy),
		)
		return nil
	}
}
