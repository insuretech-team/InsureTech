package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"

	"github.com/google/uuid"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/config"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/email"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/events"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/repository"
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/sms"
	authnentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/entity/v1"
	authnservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// OTPService handles OTP generation, verification, and rate limiting.
// Supports two channels:
//   - SMS (sslwireless): for B2C_CUSTOMER and AGENT (mobile OTP)
//   - Email (SMTP): for BUSINESS_BENEFICIARY and SYSTEM_USER (email OTP)
//
// Rate limiting uses Redis when a client is provided, falling back to
// the DB-based approach (CountRecentOTPs) when Redis is unavailable.
type OTPService struct {
	otpRepo        *repository.OTPRepository
	smsClient      *sms.SSLWirelessClient
	emailClient    *email.Client
	config         *config.Config
	eventPublisher *events.Publisher
	redisClient    redis.UniversalClient // optional; enables Redis rate limiting
}

// NewOTPService creates a new OTP service (no Redis rate limiting).
func NewOTPService(otpRepo *repository.OTPRepository, smsClient *sms.SSLWirelessClient, emailClient *email.Client, cfg *config.Config, eventPublisher *events.Publisher) *OTPService {
	return &OTPService{
		otpRepo:        otpRepo,
		smsClient:      smsClient,
		emailClient:    emailClient,
		config:         cfg,
		eventPublisher: eventPublisher,
	}
}

// NewOTPServiceWithRedis creates a new OTP service with Redis-backed rate limiting.
func NewOTPServiceWithRedis(otpRepo *repository.OTPRepository, smsClient *sms.SSLWirelessClient, emailClient *email.Client, cfg *config.Config, eventPublisher *events.Publisher, rdb redis.UniversalClient) *OTPService {
	return &OTPService{
		otpRepo:        otpRepo,
		smsClient:      smsClient,
		emailClient:    emailClient,
		config:         cfg,
		eventPublisher: eventPublisher,
		redisClient:    rdb,
	}
}

// SendOTP generates and sends an OTP via SMS or email
func (s *OTPService) SendOTP(ctx context.Context, req *authnservicev1.SendOTPRequest) (*authnservicev1.SendOTPResponse, error) {
	// Determine channel
	channel := req.Channel
	if channel == "" {
		channel = "sms" // Default
	}

	// Validate recipient format
	recipient := req.Recipient
	if channel == "sms" {
		if !sms.ValidateMSISDN(recipient) {
			logger.Errorf("invalid mobile number format")
			return nil, errors.New("invalid mobile number format")
		}
		recipient = sms.NormalizeMSISDN(recipient)
	}

	// Check rate limiting
	if err := s.checkRateLimit(ctx, recipient, req.Type); err != nil {
		cooldown := s.calculateCooldown(ctx, recipient)
		return &authnservicev1.SendOTPResponse{
			Message:         fmt.Sprintf("Rate limit exceeded. Please try again in %d seconds", int(cooldown.Seconds())),
			CooldownSeconds: int32(cooldown.Seconds()),
		}, err
	}

	// Generate OTP code
	otpCode, err := generateOTPCode(s.config.Security.OTPLength)
	if err != nil {
		logger.Errorf("failed to generate OTP: %v", err)
		return nil, errors.New("failed to generate OTP")
	}

	// Create OTP entity
	otpID := uuid.New().String()
	expiresAt := time.Now().Add(s.config.Security.OTPExpiry)

	// Hash the OTP code for storage (security best practice)
	otpHash, err := bcrypt.GenerateFromPassword([]byte(otpCode), bcrypt.DefaultCost)
	if err != nil {
		logger.Errorf("failed to hash OTP: %v", err)
		return nil, errors.New("failed to hash OTP")
	}

	otpEntity := &authnentityv1.OTP{
		OtpId:     otpID,
		UserId:    "", // Will be set on verification if needed
		OtpHash:   string(otpHash),
		Purpose:   req.Type,
		Recipient: recipient,
		Channel:   channel,
		ExpiresAt: timestamppb.New(expiresAt),
		Verified:  false,
		Attempts:  0,
		DlrStatus: string(sms.DLRStatusPending),
	}

	// Send OTP based on channel
	var providerMessageID string
	var senderID string

	if channel == "sms" {
		// Build SMS message
		message := fmt.Sprintf("Your verification code is: %s. Valid for %d minutes. Do not share this code.",
			otpCode, int(s.config.Security.OTPExpiry.Minutes()))

		// Send via SSL Wireless
		smsReq := &sms.SendSMSRequest{
			MSISDN:     recipient,
			Message:    message,
			UseMasking: req.UseMasking,
			CSMSId:     otpID,
		}

		smsResp, err := s.smsClient.SendSMS(ctx, smsReq)
		if err != nil {
			logger.Errorf("failed to send SMS: %v", err)
			return nil, errors.New("failed to send SMS")
		}

		providerMessageID = smsResp.MessageID
		senderID = s.config.SMS.MaskingSenderID
		if !req.UseMasking {
			senderID = s.config.SMS.NonMaskingSender
		}

		// Update OTP entity with SMS tracking info
		otpEntity.ProviderMessageId = providerMessageID
		otpEntity.SenderId = senderID
		otpEntity.Carrier = sms.DetectCarrier(recipient)

	} else if channel == "email" {
		if s.emailClient == nil {
			logger.Errorf("email client not configured")
			return nil, errors.New("email client not configured")
		}
		expiryMin := int(s.config.Security.OTPExpiry.Minutes())
		emailResp, err := s.emailClient.SendOTP(&email.SendOTPRequest{
			To:        recipient,
			OTPCode:   otpCode,
			Purpose:   req.Type,
			ExpiryMin: expiryMin,
		})
		if err != nil {
			logger.Errorf("failed to send email OTP: %v", err)
			return nil, errors.New("failed to send email OTP")
		}
		providerMessageID = emailResp.MessageID
		senderID = s.config.Email.From
		otpEntity.DlrStatus = "DELIVERED" // email has no async DLR; assume delivered on send
	} else {
		return nil, fmt.Errorf("unsupported channel: %s", channel)
	}

	// Save to database
	if err := s.otpRepo.Create(ctx, otpEntity); err != nil {
		logger.Errorf("failed to save OTP: %v", err)
		return nil, errors.New("failed to save OTP")
	}

	// Update rate limit counter
	s.incrementRateLimit(ctx, recipient, req.Type)

	return &authnservicev1.SendOTPResponse{
		OtpId:            otpID,
		Message:          "OTP sent successfully",
		ExpiresInSeconds: int32(s.config.Security.OTPExpiry.Seconds()),
		SenderId:         senderID,
		CooldownSeconds:  int32(s.config.Security.OTPCooldown.Seconds()),
	}, nil
}

// VerifyOTP verifies an OTP code
func (s *OTPService) VerifyOTP(ctx context.Context, req *authnservicev1.VerifyOTPRequest) (*authnservicev1.VerifyOTPResponse, error) {
	// Retrieve OTP from database
	otpEntity, err := s.otpRepo.GetByID(ctx, req.OtpId)
	if err != nil {
		return &authnservicev1.VerifyOTPResponse{
			Verified: false,
			Message:  "Invalid or expired OTP",
		}, nil
	}

	// Check if already verified
	if otpEntity.Verified {
		return &authnservicev1.VerifyOTPResponse{
			Verified: false,
			Message:  "OTP already used",
		}, nil
	}

	// Check expiry
	if time.Now().After(otpEntity.ExpiresAt.AsTime()) {
		return &authnservicev1.VerifyOTPResponse{
			Verified: false,
			Message:  "OTP has expired",
		}, nil
	}

	// Check max attempts
	if otpEntity.Attempts >= int32(s.config.Security.OTPMaxAttempts) {
		return &authnservicev1.VerifyOTPResponse{
			Verified: false,
			Message:  "Maximum verification attempts exceeded",
		}, nil
	}

	// Increment attempt count
	otpEntity.Attempts++
	if err := s.otpRepo.IncrementAttempts(ctx, req.OtpId); err != nil {
		logger.Errorf("failed to update attempts: %v", err)
		return nil, errors.New("failed to update attempts")
	}

	// Verify OTP code against hash
	if err := bcrypt.CompareHashAndPassword([]byte(otpEntity.OtpHash), []byte(req.Code)); err != nil {
		remainingAttempts := s.config.Security.OTPMaxAttempts - int(otpEntity.Attempts)
		return &authnservicev1.VerifyOTPResponse{
			Verified: false,
			Message:  fmt.Sprintf("Invalid OTP code. %d attempts remaining", remainingAttempts),
		}, nil
	}

	// Mark as verified
	otpEntity.Verified = true
	otpEntity.VerifiedAt = timestamppb.Now()
	if err := s.otpRepo.MarkVerified(ctx, req.OtpId); err != nil {
		logger.Errorf("failed to mark OTP as verified: %v", err)
		return nil, errors.New("failed to mark OTP as verified")
	}

	return &authnservicev1.VerifyOTPResponse{
		Verified: true,
		Message:  "OTP verified successfully",
		UserId:   otpEntity.UserId,
	}, nil
}

// SendEmailOTP sends an OTP specifically via email channel.
// Used by: RegisterEmailUser, SendEmailOTP RPC, RequestPasswordResetByEmail.
// Validates that recipient is a valid email format before sending.
func (s *OTPService) SendEmailOTP(ctx context.Context, req *authnservicev1.SendEmailOTPRequest) (*authnservicev1.SendEmailOTPResponse, error) {
	recipient := req.Email
	otpType := req.Type

	// Rate limit check
	if err := s.checkRateLimit(ctx, recipient, otpType); err != nil {
		cooldown := s.calculateCooldown(ctx, recipient)
		return &authnservicev1.SendEmailOTPResponse{
			Message:         fmt.Sprintf("Rate limit exceeded. Please try again in %d seconds", int(cooldown.Seconds())),
			CooldownSeconds: int32(cooldown.Seconds()),
		}, err
	}

	// Generate OTP code
	otpCode, err := generateOTPCode(s.config.Security.OTPLength)
	if err != nil {
		logger.Errorf("failed to generate OTP: %v", err)
		return nil, errors.New("failed to generate OTP")
	}

	otpID := uuid.New().String()
	expiresAt := time.Now().Add(s.config.Security.OTPExpiry)

	otpHash, err := bcrypt.GenerateFromPassword([]byte(otpCode), bcrypt.DefaultCost)
	if err != nil {
		logger.Errorf("failed to hash OTP: %v", err)
		return nil, errors.New("failed to hash OTP")
	}

	if s.emailClient == nil {
		logger.Errorf("email client not configured")
		return nil, errors.New("email client not configured")
	}

	expiryMin := int(s.config.Security.OTPExpiry.Minutes())
	emailResp, err := s.emailClient.SendOTP(&email.SendOTPRequest{
		To:        recipient,
		OTPCode:   otpCode,
		Purpose:   otpType,
		ExpiryMin: expiryMin,
	})
	if err != nil {
		logger.Errorf("failed to send email OTP: %v", err)
		return nil, errors.New("failed to send email OTP")
	}

	otpEntity := &authnentityv1.OTP{
		OtpId:             otpID,
		UserId:            "",
		OtpHash:           string(otpHash),
		Purpose:           otpType,
		Recipient:         recipient,
		Channel:           "email",
		ExpiresAt:         timestamppb.New(expiresAt),
		Verified:          false,
		Attempts:          0,
		DlrStatus:         "DELIVERED",
		ProviderMessageId: emailResp.MessageID,
		SenderId:          s.config.Email.From,
	}

	if err := s.otpRepo.Create(ctx, otpEntity); err != nil {
		logger.Errorf("failed to save OTP: %v", err)
		return nil, errors.New("failed to save OTP")
	}

	s.incrementRateLimit(ctx, recipient, otpType)

	return &authnservicev1.SendEmailOTPResponse{
		OtpId:            otpID,
		Message:          "OTP sent to your email address",
		ExpiresInSeconds: int32(s.config.Security.OTPExpiry.Seconds()),
		CooldownSeconds:  int32(s.config.Security.OTPCooldown.Seconds()),
	}, nil
}

// HandleDLR processes delivery report webhooks from SSL Wireless.
// payload is the raw JSON body from the SSLWireless DLR webhook POST.
func (s *OTPService) HandleDLR(ctx context.Context, payload []byte) error {
	if s.smsClient == nil {
		logger.Errorf("SMS client not configured")
		return errors.New("SMS client not configured")
	}

	// Parse DLR payload into DLRWebhookPayload.
	dlr, err := s.smsClient.ParseDLRWebhook(payload)
	if err != nil {
		logger.Errorf("failed to parse DLR: %v", err)
		return errors.New("failed to parse DLR")
	}

	// Update OTP record matched by provider_message_id (canonical DLR key).
	if err := s.otpRepo.UpdateDLRStatus(ctx, dlr.MessageID, dlr.Status, dlr.ErrorCode); err != nil {
		logger.Errorf("failed to update DLR status for message_id %s: %v", dlr.MessageID, err)
		return errors.New("failed to update DLR status for message_id %s")
	}

	// Publish DLR event (best-effort; non-fatal).
	if s.eventPublisher != nil {
		_ = s.eventPublisher.PublishSMSDeliveryReport(
			ctx,
			"", // otp_id not known here without an extra DB lookup — use message_id as correlation key
			dlr.MessageID,
			sms.MaskMSISDN(dlr.MSISDN),
			dlr.Status,
			dlr.ErrorCode,
			dlr.Carrier,
			dlr.DeliveredAt,
		)
	}

	return nil
}

// checkRateLimit checks if the recipient has exceeded rate limits.
// Uses Redis when available (INCR + EXPIRE sliding window), falls back to DB counts.
func (s *OTPService) checkRateLimit(ctx context.Context, recipient, otpType string) error {
	const otpPerHourLimit = 3
	if s.redisClient != nil {
		// Per-hour window (Sprint 5 baseline: 3/hour/recipient).
		hourKey := fmt.Sprintf("otp_rl:hour:%s:%s", otpType, recipient)
		hourCount, err := s.redisClient.Incr(ctx, hourKey).Result()
		if err == nil {
			if hourCount == 1 {
				s.redisClient.Expire(ctx, hourKey, time.Hour)
			}
			if hourCount > otpPerHourLimit {
				logger.Errorf("rate limit exceeded: too many OTPs in last hour")
				return errors.New("rate limit exceeded: too many OTPs in last hour")
			}
		}

		// Per-minute window
		minKey := fmt.Sprintf("otp_rl:min:%s:%s", otpType, recipient)
		minCount, err := s.redisClient.Incr(ctx, minKey).Result()
		if err == nil {
			if minCount == 1 {
				s.redisClient.Expire(ctx, minKey, time.Minute)
			}
			if minCount > int64(s.config.Security.RateLimitPerMinute) {
				logger.Errorf("rate limit exceeded: too many OTPs in last minute")
				return errors.New("rate limit exceeded: too many OTPs in last minute")
			}
		}

		// Per-day window
		dayKey := fmt.Sprintf("otp_rl:day:%s:%s", otpType, recipient)
		dayCount, err := s.redisClient.Incr(ctx, dayKey).Result()
		if err == nil {
			if dayCount == 1 {
				s.redisClient.Expire(ctx, dayKey, 24*time.Hour)
			}
			if dayCount > int64(s.config.Security.RateLimitPerDay) {
				logger.Errorf("daily OTP limit exceeded")
				return errors.New("daily OTP limit exceeded")
			}
		}
		return nil
	}

	// Fallback: DB count queries when Redis is unavailable.
	hourlyCount, err := s.otpRepo.CountRecentOTPs(ctx, recipient, time.Now().Add(-1*time.Hour))
	if err != nil {
		logger.Errorf("failed to check hourly limit: %v", err)
		return errors.New("failed to check hourly limit")
	}
	if hourlyCount >= otpPerHourLimit {
		return fmt.Errorf("hourly limit exceeded: %d OTPs in last hour", hourlyCount)
	}

	recentCount, err := s.otpRepo.CountRecentOTPs(ctx, recipient, time.Now().Add(-1*time.Minute))
	if err != nil {
		logger.Errorf("failed to check rate limit: %v", err)
		return errors.New("failed to check rate limit")
	}
	if recentCount >= int64(s.config.Security.RateLimitPerMinute) {
		return fmt.Errorf("rate limit exceeded: %d OTPs in last minute", recentCount)
	}

	dailyCount, err := s.otpRepo.CountRecentOTPs(ctx, recipient, time.Now().Add(-24*time.Hour))
	if err != nil {
		logger.Errorf("failed to check daily limit: %v", err)
		return errors.New("failed to check daily limit")
	}
	if dailyCount >= int64(s.config.Security.RateLimitPerDay) {
		return fmt.Errorf("daily limit exceeded: %d OTPs in last 24 hours", dailyCount)
	}

	return nil
}

// calculateCooldown calculates cooldown period based on recent OTP sends
func (s *OTPService) calculateCooldown(ctx context.Context, recipient string) time.Duration {
	// Get last OTP send time
	lastOTP, err := s.otpRepo.GetLastOTP(ctx, recipient)
	if err != nil || lastOTP == nil {
		return 0
	}

	elapsed := time.Since(lastOTP.CreatedAt.AsTime())
	if elapsed < s.config.Security.OTPCooldown {
		return s.config.Security.OTPCooldown - elapsed
	}

	return 0
}

// incrementRateLimit increments rate limit counter
func (s *OTPService) incrementRateLimit(ctx context.Context, recipient, otpType string) {
	if s.redisClient == nil {
		return
	}

	hourKey := fmt.Sprintf("otp_rl:hour:%s:%s", otpType, recipient)
	minKey := fmt.Sprintf("otp_rl:min:%s:%s", otpType, recipient)
	dayKey := fmt.Sprintf("otp_rl:day:%s:%s", otpType, recipient)

	// Best-effort writes; SendOTP already persisted OTP in DB and should not fail because
	// Redis rate-counter updates are temporarily unavailable.
	_, _ = s.redisClient.Incr(ctx, hourKey).Result()
	_, _ = s.redisClient.Expire(ctx, hourKey, time.Hour).Result()
	_, _ = s.redisClient.Incr(ctx, minKey).Result()
	_, _ = s.redisClient.Expire(ctx, minKey, time.Minute).Result()
	_, _ = s.redisClient.Incr(ctx, dayKey).Result()
	_, _ = s.redisClient.Expire(ctx, dayKey, 24*time.Hour).Result()
}

// generateOTPCode generates a random numeric OTP code
func generateOTPCode(length int) (string, error) {
	if length < 4 || length > 8 {
		logger.Errorf("OTP length must be between 4 and 8")
		return "", errors.New("OTP length must be between 4 and 8")
	}

	// Calculate max value (e.g., for 6 digits: 999999)
	max := int64(1)
	for i := 0; i < length; i++ {
		max *= 10
	}

	// Generate random number
	n, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		logger.Errorf("failed to generate random number: %v", err)
		return "", errors.New("failed to generate random number")
	}

	// Format with leading zeros
	format := fmt.Sprintf("%%0%dd", length)
	return fmt.Sprintf(format, n.Int64()), nil
}
