package service

import (
	"context"
	"errors"
	"strconv"
	"time"

	authnservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	defaultMaxResendsPerDay = 3
)

// ResendOTP invalidates a previous OTP and issues a new OTP to the same recipient/channel.
//
// Production constraints:
// - Enforces cooldown: cannot resend until OTPCooldown has elapsed since the original OTP creation.
// - Enforces coarse daily limit using CountRecentOTPs (DB-based). (No dedicated resend counter column exists.)
// - Expires the original OTP by setting expires_at=now.
func (s *AuthService) ResendOTP(ctx context.Context, req *authnservicev1.ResendOTPRequest) (*authnservicev1.ResendOTPResponse, error) {
	if req == nil || req.OriginalOtpId == "" {
		return nil, errors.New("original_otp_id is required")
	}

	if s.otpRepo == nil || s.otpService == nil {
		return nil, errors.New("OTP subsystem not configured")
	}

	orig, err := s.otpRepo.GetByID(ctx, req.OriginalOtpId)
	if err != nil || orig == nil {
		return nil, errors.New("original OTP not found")
	}

	// Do not resend verified OTPs.
	if orig.Verified {
		return nil, errors.New("OTP already verified")
	}

	cooldown := s.config.Security.OTPCooldown
	if cooldown <= 0 {
		cooldown = 60 * time.Second
	}

	createdAt := time.Time{}
	if orig.CreatedAt != nil {
		createdAt = orig.CreatedAt.AsTime()
	}
	elapsed := time.Since(createdAt)
	if !createdAt.IsZero() && elapsed < cooldown {
		wait := cooldown - elapsed
		return &authnservicev1.ResendOTPResponse{
			Message:           "Cooldown active. Try again in " + strconv.Itoa(int(wait.Seconds())) + " seconds",
			CooldownSeconds:   int32(wait.Seconds()),
			AttemptsRemaining: 0,
			CanRetryAt:        timestamppb.New(time.Now().Add(wait)),
		}, nil
	}

	// Coarse daily limit by recipient+pupose.
	maxResends := defaultMaxResendsPerDay
	if s.config.Security.RateLimitPerDay > 0 {
		// Keep resend limit <= overall OTP daily limit
		if s.config.Security.RateLimitPerDay < maxResends {
			maxResends = s.config.Security.RateLimitPerDay
		}
	}

	since := time.Now().Add(-24 * time.Hour)
	count, err := s.otpRepo.CountRecentOTPs(ctx, orig.Recipient, since)
	if err == nil && int(count) >= maxResends {
		return &authnservicev1.ResendOTPResponse{
			Message:           "Resend limit exceeded. Please try later.",
			AttemptsRemaining: 0,
			CooldownSeconds:   int32(cooldown.Seconds()),
			CanRetryAt:        timestamppb.New(time.Now().Add(24 * time.Hour)),
			ExpiresInSeconds:  0,
			OtpId:             "",
			SenderId:          orig.SenderId,
		}, nil
	}

	// Expire old OTP (best effort)
	_ = s.otpRepo.ExpireOTP(ctx, orig.OtpId)

	// Issue new OTP with same recipient/purpose/channel.
	sendReq := &authnservicev1.SendOTPRequest{
		Recipient:  orig.Recipient,
		Type:       orig.Purpose,
		Channel:    orig.Channel,
		UseMasking: orig.SenderId != "", // heuristic; original request flag isn't stored
	}
	newOTP, err := s.otpService.SendOTP(ctx, sendReq)
	if err != nil {
		return nil, err
	}

	// attempts remaining: approximate based on daily cap
	attemptsRemaining := int32(maxResends)
	if err == nil {
		attemptsRemaining = int32(maxResends - int(count) - 1)
		if attemptsRemaining < 0 {
			attemptsRemaining = 0
		}
	}

	canRetryAt := time.Now().Add(cooldown)
	return &authnservicev1.ResendOTPResponse{
		OtpId:             newOTP.OtpId,
		Message:           "OTP resent successfully",
		ExpiresInSeconds:  newOTP.ExpiresInSeconds,
		CooldownSeconds:   int32(cooldown.Seconds()),
		AttemptsRemaining: attemptsRemaining,
		CanRetryAt:        timestamppb.New(canRetryAt),
		SenderId:          newOTP.SenderId,
	}, nil
}
