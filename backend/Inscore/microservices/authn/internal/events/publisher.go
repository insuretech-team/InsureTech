package events

import (
	"context"
	"time"

	"github.com/google/uuid"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	authneventsv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/events/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Kafka topic names for authn events
const (
	TopicUserRegistered                = "authn.user.registered"
	TopicUserLoggedIn                  = "authn.user.logged_in"
	TopicUserLoggedOut                 = "authn.user.logged_out"
	TopicOTPSent                       = "authn.otp.sent"
	TopicOTPVerified                   = "authn.otp.verified"
	TopicEmailVerificationSent         = "authn.email.verification_sent"
	TopicEmailVerified                 = "authn.email.verified"
	TopicEmailLoginSucceeded           = "authn.email.login_succeeded"
	TopicEmailLoginFailed              = "authn.email.login_failed"
	TopicPasswordResetByEmailRequested = "authn.email.password_reset_requested"
	TopicSessionRevoked                = "authn.session.revoked"
	TopicSessionExpired                = "authn.session.expired"
	TopicTokenRefreshed                = "authn.token.refreshed"
	TopicPasswordChanged               = "authn.password.changed"
	TopicPasswordResetRequested        = "authn.password.reset_requested"
	TopicLoginFailed                   = "authn.login.failed"
	TopicSMSDeliveryReport             = "authn.otp.sms_dlr"
	TopicCSRFValidationFailed          = "authn.csrf.validation_failed"
	TopicAccountLocked                 = "authn.account.locked"
)

// EventProducer interface for decoupling Kafka/MsgQueue
type EventProducer interface {
	Produce(ctx context.Context, topic string, key string, msg interface{}) error
	Close() error
}

type Publisher struct {
	producer EventProducer
}

func NewPublisher(producer EventProducer) *Publisher {
	return &Publisher{producer: producer}
}

func (p *Publisher) PublishUserRegistered(ctx context.Context, userID, mobile, email, ip, deviceType string) error {
	evt := &authneventsv1.UserRegisteredEvent{
		EventId:      uuid.New().String(),
		UserId:       userID,
		MobileNumber: mobile,
		Email:        email,
		Timestamp:    timestamppb.New(time.Now()),
		IpAddress:    ip,
		DeviceType:   deviceType,
	}
	if err := p.publish(ctx, TopicUserRegistered, userID, evt); err != nil {
		appLogger.Warnf("Failed to publish UserRegisteredEvent for user %s: %v", userID, err)
	}
	return nil
}

func (p *Publisher) PublishUserLoggedIn(ctx context.Context, userID, sessionID, sessionType, ip, deviceType, userAgent string) error {
	evt := &authneventsv1.UserLoggedInEvent{
		EventId:     uuid.New().String(),
		UserId:      userID,
		SessionId:   sessionID,
		SessionType: sessionType,
		Timestamp:   timestamppb.New(time.Now()),
		IpAddress:   ip,
		DeviceType:  deviceType,
		UserAgent:   userAgent,
	}
	if err := p.publish(ctx, TopicUserLoggedIn, userID, evt); err != nil {
		appLogger.Warnf("Failed to publish UserLoggedInEvent for user %s: %v", userID, err)
	}
	return nil
}

// PublishEmailVerificationSent publishes an event when an email verification OTP is sent
func (p *Publisher) PublishEmailVerificationSent(ctx context.Context, userID, email, otpID, otpType, ipAddress string) error {
	evt := &authneventsv1.EmailVerificationSentEvent{
		EventId:     uuid.New().String(),
		UserId:      userID,
		EmailMasked: maskEmailForEvent(email),
		OtpId:       otpID,
		Type:        otpType,
		Timestamp:   timestamppb.New(time.Now()),
		IpAddress:   ipAddress,
	}
	key := userID
	if key == "" {
		key = otpID
	}
	if err := p.publish(ctx, TopicEmailVerificationSent, key, evt); err != nil {
		appLogger.Warnf("Failed to publish EmailVerificationSentEvent: %v", err)
	}
	return nil
}

// PublishEmailVerified publishes an event when an email is verified
func (p *Publisher) PublishEmailVerified(ctx context.Context, userID, email string) error {
	evt := &authneventsv1.EmailVerifiedEvent{
		EventId:   uuid.New().String(),
		UserId:    userID,
		Email:     maskEmailForEvent(email),
		Timestamp: timestamppb.New(time.Now()),
	}
	if err := p.publish(ctx, TopicEmailVerified, userID, evt); err != nil {
		appLogger.Warnf("Failed to publish EmailVerifiedEvent for user %s: %v", userID, err)
	}
	return nil
}

// PublishEmailLoginSucceeded publishes an event when email OTP login succeeds
func (p *Publisher) PublishEmailLoginSucceeded(ctx context.Context, userID, sessionID, email, userType, ipAddress, userAgent, deviceName string) error {
	evt := &authneventsv1.EmailLoginSucceededEvent{
		EventId:     uuid.New().String(),
		UserId:      userID,
		SessionId:   sessionID,
		EmailMasked: maskEmailForEvent(email),
		UserType:    userType,
		Timestamp:   timestamppb.New(time.Now()),
		IpAddress:   ipAddress,
		UserAgent:   userAgent,
		DeviceName:  deviceName,
	}
	if err := p.publish(ctx, TopicEmailLoginSucceeded, userID, evt); err != nil {
		appLogger.Warnf("Failed to publish EmailLoginSucceededEvent for user %s: %v", userID, err)
	}
	return nil
}

// PublishEmailLoginFailed publishes an event when email OTP login fails
func (p *Publisher) PublishEmailLoginFailed(ctx context.Context, email, reason string, attempts int32, ipAddress, userAgent string) error {
	evt := &authneventsv1.EmailLoginFailedEvent{
		EventId:             uuid.New().String(),
		EmailMasked:         maskEmailForEvent(email),
		FailureReason:       reason,
		FailedAttemptsCount: attempts,
		Timestamp:           timestamppb.New(time.Now()),
		IpAddress:           ipAddress,
		UserAgent:           userAgent,
	}
	if err := p.publish(ctx, TopicEmailLoginFailed, maskEmailForEvent(email), evt); err != nil {
		appLogger.Warnf("Failed to publish EmailLoginFailedEvent: %v", err)
	}
	return nil
}

// PublishPasswordResetByEmailRequested publishes an event when a password reset is requested via email
func (p *Publisher) PublishPasswordResetByEmailRequested(ctx context.Context, userID, email, otpID, ipAddress string) error {
	evt := &authneventsv1.PasswordResetByEmailRequestedEvent{
		EventId:     uuid.New().String(),
		UserId:      userID,
		EmailMasked: maskEmailForEvent(email),
		OtpId:       otpID,
		Timestamp:   timestamppb.New(time.Now()),
		IpAddress:   ipAddress,
	}
	if err := p.publish(ctx, TopicPasswordResetByEmailRequested, userID, evt); err != nil {
		appLogger.Warnf("Failed to publish PasswordResetByEmailRequestedEvent for user %s: %v", userID, err)
	}
	return nil
}

// publish is the internal helper that sends to Kafka if a producer is wired,
// or logs a warning and no-ops if running without Kafka (e.g. dev/test).
func (p *Publisher) publish(ctx context.Context, topic, key string, msg interface{}) error {
	if p.producer == nil {
		appLogger.Infof("Kafka producer not configured - event dropped (topic=%s, key=%s)", topic, key)
		return nil
	}
	return p.producer.Produce(ctx, topic, key, msg)
}

// PublishOTPSent publishes a typed OTPSentEvent (proto) when an OTP is sent.
// Note: caller should pass provider/sender details when available.
func (p *Publisher) PublishOTPSent(ctx context.Context, otpID, recipientMasked, otpType, channel, provider, senderID, providerMessageID string, maskingUsed bool) error {
	evt := &authneventsv1.OTPSentEvent{
		EventId:           uuid.New().String(),
		OtpId:             otpID,
		Recipient:         recipientMasked,
		Type:              otpType,
		Timestamp:         timestamppb.New(time.Now()),
		Channel:           channel,
		Provider:          provider,
		SenderId:          senderID,
		ProviderMessageId: providerMessageID,
		MaskingUsed:       maskingUsed,
	}
	if err := p.publish(ctx, TopicOTPSent, otpID, evt); err != nil {
		appLogger.Warnf("Failed to publish OTPSentEvent (otp_id=%s): %v", otpID, err)
	}
	return nil
}

// PublishOTPVerified publishes a typed OTPVerifiedEvent (proto) when OTP verification succeeds.
func (p *Publisher) PublishOTPVerified(ctx context.Context, otpID, userID string, attempts int32) error {
	evt := &authneventsv1.OTPVerifiedEvent{
		EventId:   uuid.New().String(),
		OtpId:     otpID,
		UserId:    userID,
		Timestamp: timestamppb.New(time.Now()),
		Attempts:  attempts,
	}
	if err := p.publish(ctx, TopicOTPVerified, otpID, evt); err != nil {
		appLogger.Warnf("Failed to publish OTPVerifiedEvent (otp_id=%s): %v", otpID, err)
	}
	return nil
}

func (p *Publisher) PublishUserLoggedOut(ctx context.Context, userID, sessionID, sessionType, reason, ip, deviceType string) error {
	evt := &authneventsv1.UserLoggedOutEvent{
		EventId:      uuid.New().String(),
		UserId:       userID,
		SessionId:    sessionID,
		SessionType:  sessionType,
		LogoutReason: reason,
		Timestamp:    timestamppb.New(time.Now()),
		IpAddress:    ip,
		DeviceType:   deviceType,
	}
	if err := p.publish(ctx, TopicUserLoggedOut, userID, evt); err != nil {
		appLogger.Warnf("Failed to publish UserLoggedOutEvent for user %s: %v", userID, err)
	}
	return nil
}

func (p *Publisher) PublishLoginFailed(ctx context.Context, userID, mobile, reason, ip, deviceType, userAgent string, failedAttempts int32) error {
	evt := &authneventsv1.LoginFailedEvent{
		EventId:             uuid.New().String(),
		UserId:              userID,
		MobileNumber:        mobile,
		FailureReason:       reason,
		Timestamp:           timestamppb.New(time.Now()),
		IpAddress:           ip,
		DeviceType:          deviceType,
		UserAgent:           userAgent,
		FailedAttemptsCount: failedAttempts,
	}
	key := userID
	if key == "" {
		key = mobile
	}
	if err := p.publish(ctx, TopicLoginFailed, key, evt); err != nil {
		appLogger.Warnf("Failed to publish LoginFailedEvent (user_id=%s mobile=%s): %v", userID, mobile, err)
	}
	return nil
}

func (p *Publisher) PublishTokenRefreshed(ctx context.Context, userID, sessionID, oldJTI, newAccessJTI, newRefreshJTI, ip, deviceType, userAgent string) error {
	evt := &authneventsv1.TokenRefreshedEvent{
		EventId:            uuid.New().String(),
		UserId:             userID,
		SessionId:          sessionID,
		OldAccessTokenJti:  oldJTI,
		NewAccessTokenJti:  newAccessJTI,
		NewRefreshTokenJti: newRefreshJTI,
		Timestamp:          timestamppb.New(time.Now()),
		IpAddress:          ip,
		DeviceType:         deviceType,
		UserAgent:          userAgent,
	}
	if err := p.publish(ctx, TopicTokenRefreshed, userID, evt); err != nil {
		appLogger.Warnf("Failed to publish TokenRefreshedEvent (user_id=%s session_id=%s): %v", userID, sessionID, err)
	}
	return nil
}

func (p *Publisher) PublishPasswordChanged(ctx context.Context, userID, ip, changedBy string) error {
	evt := &authneventsv1.PasswordChangedEvent{
		EventId:   uuid.New().String(),
		UserId:    userID,
		Timestamp: timestamppb.New(time.Now()),
		IpAddress: ip,
		ChangedBy: changedBy,
	}
	if err := p.publish(ctx, TopicPasswordChanged, userID, evt); err != nil {
		appLogger.Warnf("Failed to publish PasswordChangedEvent (user_id=%s): %v", userID, err)
	}
	return nil
}

func (p *Publisher) PublishPasswordResetRequested(ctx context.Context, userID, mobile, ip, deviceType string) error {
	evt := &authneventsv1.PasswordResetRequestedEvent{
		EventId:      uuid.New().String(),
		UserId:       userID,
		MobileNumber: mobile,
		Timestamp:    timestamppb.New(time.Now()),
		IpAddress:    ip,
		DeviceType:   deviceType,
	}
	if err := p.publish(ctx, TopicPasswordResetRequested, userID, evt); err != nil {
		appLogger.Warnf("Failed to publish PasswordResetRequestedEvent (user_id=%s): %v", userID, err)
	}
	return nil
}

func (p *Publisher) PublishSessionRevoked(ctx context.Context, userID, sessionID, sessionType, revokedBy, reason string) error {
	evt := &authneventsv1.SessionRevokedEvent{
		EventId:     uuid.New().String(),
		UserId:      userID,
		SessionId:   sessionID,
		SessionType: sessionType,
		RevokedBy:   revokedBy,
		Reason:      reason,
		Timestamp:   timestamppb.New(time.Now()),
	}
	if err := p.publish(ctx, TopicSessionRevoked, userID, evt); err != nil {
		appLogger.Warnf("Failed to publish SessionRevokedEvent (user_id=%s session_id=%s): %v", userID, sessionID, err)
	}
	return nil
}

func (p *Publisher) PublishCSRFValidationFailed(ctx context.Context, userID, sessionID, expectedHash, receivedHash, ip, userAgent, path, method string) error {
	evt := &authneventsv1.CSRFValidationFailedEvent{
		EventId:               uuid.New().String(),
		UserId:                userID,
		SessionId:             sessionID,
		ExpectedCsrfTokenHash: expectedHash,
		ReceivedCsrfTokenHash: receivedHash,
		Timestamp:             timestamppb.New(time.Now()),
		IpAddress:             ip,
		UserAgent:             userAgent,
		RequestPath:           path,
		RequestMethod:         method,
	}
	if err := p.publish(ctx, TopicCSRFValidationFailed, userID, evt); err != nil {
		appLogger.Warnf("Failed to publish CSRFValidationFailedEvent (user_id=%s session_id=%s): %v", userID, sessionID, err)
	}
	return nil
}

func (p *Publisher) PublishSMSDeliveryReport(ctx context.Context, otpID, providerMsgID, msisdnMasked, status, errorCode, carrier string, deliveredAt time.Time) error {
	evt := &authneventsv1.SMSDeliveryReportEvent{
		EventId:           uuid.New().String(),
		OtpId:             otpID,
		ProviderMessageId: providerMsgID,
		Msisdn:            msisdnMasked,
		Status:            status,
		ErrorCode:         errorCode,
		DeliveredAt:       timestamppb.New(deliveredAt),
		Timestamp:         timestamppb.New(time.Now()),
		Carrier:           carrier,
	}
	if err := p.publish(ctx, TopicSMSDeliveryReport, otpID, evt); err != nil {
		appLogger.Warnf("Failed to publish SMSDeliveryReportEvent (otp_id=%s provider_msg_id=%s): %v", otpID, providerMsgID, err)
	}
	return nil
}

func (p *Publisher) PublishSessionExpired(ctx context.Context, userID, sessionID, sessionType string, expiredAt time.Time, inactivitySeconds int32) error {
	evt := &authneventsv1.SessionExpiredEvent{
		EventId:           uuid.New().String(),
		UserId:            userID,
		SessionId:         sessionID,
		SessionType:       sessionType,
		Timestamp:         timestamppb.New(time.Now()),
		ExpiredAt:         timestamppb.New(expiredAt),
		InactivitySeconds: inactivitySeconds,
	}
	if err := p.publish(ctx, TopicSessionExpired, userID, evt); err != nil {
		appLogger.Warnf("Failed to publish SessionExpiredEvent (user_id=%s session_id=%s): %v", userID, sessionID, err)
	}
	return nil
}

func (p *Publisher) PublishAccountLocked(ctx context.Context, userID, reason string, lockedUntil time.Time) error {
	evt := &authneventsv1.AccountLockedEvent{
		EventId:     uuid.New().String(),
		UserId:      userID,
		Reason:      reason,
		Timestamp:   timestamppb.New(time.Now()),
		LockedUntil: timestamppb.New(lockedUntil),
	}
	if err := p.publish(ctx, TopicAccountLocked, userID, evt); err != nil {
		appLogger.Warnf("Failed to publish AccountLockedEvent (user_id=%s): %v", userID, err)
	}
	return nil
}

// maskEmailForEvent masks an email address for safe event publishing
// e.g. user@domain.com → u***@domain.com
func maskEmailForEvent(email string) string {
	for i, c := range email {
		if c == '@' {
			if i == 0 {
				return "***" + email[i:]
			}
			return string(email[0]) + "***" + email[i:]
		}
	}
	return "***"
}
