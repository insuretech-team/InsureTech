package service

// email_auth_service.go
// Email-based authentication methods for AuthService.
// RESTRICTED to: BUSINESS_BENEFICIARY, SYSTEM_USER, AGENT, and B2B_ORG_ADMIN.
// Web portal: always produces SERVER_SIDE session (cookie-based).
//
// Flows:
//   RegisterEmailUser → sends verification OTP → VerifyEmail → user active
//   SendEmailOTP(type=email_login) → EmailLogin → SERVER_SIDE session
//   RequestPasswordResetByEmail → ResetPasswordByEmail

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	authnentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/entity/v1"
	authnservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
)

// isEmailAuthUser returns true if the user_type is allowed to use email auth.
// BUSINESS_BENEFICIARY, SYSTEM_USER, AGENT, and B2B_ORG_ADMIN are permitted.
func isEmailAuthUser(userType authnentityv1.UserType) bool {
	return userType == authnentityv1.UserType_USER_TYPE_BUSINESS_BENEFICIARY ||
		userType == authnentityv1.UserType_USER_TYPE_B2B_BENEFICIARY ||
		userType == authnentityv1.UserType_USER_TYPE_SYSTEM_USER ||
		userType == authnentityv1.UserType_USER_TYPE_AGENT ||
		userType == authnentityv1.UserType_USER_TYPE_B2B_ORG_ADMIN
}

func managedEmployeeActivationAllowed(user *authnentityv1.User) bool {
	return user != nil && user.UserType == authnentityv1.UserType_USER_TYPE_B2B_BENEFICIARY
}

// maskEmail returns a masked email for safe logging: user@domain.com → u***@domain.com
func maskEmail(email string) string {
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

// RegisterEmailUser registers a new portal user with email login support.
// - Validates user_type is allowed
// - Creates user with PENDING_VERIFICATION status
// - Sends email verification OTP automatically
func (s *AuthService) RegisterEmailUser(ctx context.Context, req *authnservicev1.RegisterEmailUserRequest) (*authnservicev1.RegisterEmailUserResponse, error) {
	reqMeta := s.metadata.ExtractAll(ctx)

	// Validate user_type
	userType := parseUserType(req.UserType)
	if !isEmailAuthUser(userType) {
		appLogger.Warnf("RegisterEmailUser: rejected user_type=%s from IP %s", req.UserType, reqMeta.IPAddress)
		return nil, errors.New("email registration is only available for BUSINESS_BENEFICIARY, SYSTEM_USER, AGENT, and B2B_ORG_ADMIN")
	}

	// Email is mandatory for email-based users
	if req.Email == "" {
		return nil, errors.New("email is required for business/system user registration")
	}

	// Check if email already registered
	existing, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existing != nil {
		appLogger.Warnf("RegisterEmailUser: email %s already exists", maskEmail(req.Email))
		return nil, errors.New("email already registered")
	}

	// Optionally check mobile if provided
	if req.MobileNumber != "" {
		existingMobile, _ := s.userRepo.GetByMobileNumber(ctx, req.MobileNumber)
		if existingMobile != nil {
			return nil, errors.New("mobile number already registered")
		}
	}

	// Enforce password policy and hash with Argon2id.
	if err := validatePasswordStrength(req.Password); err != nil {
		logger.Errorf("weak password: %v", err)
		return nil, errors.New("weak password")
	}
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		logger.Errorf("failed to hash password: %v", err)
		return nil, errors.New("failed to hash password")
	}

	// Create user with PENDING_VERIFICATION — email not verified yet
	mobile := req.MobileNumber
	if mobile == "" {
		// Business users may not have a mobile; use a placeholder that passes DB constraint
		// The check constraint only applies to mobile_number format; if empty we store empty
		mobile = ""
	}

	user := &authnentityv1.User{
		UserId:        uuid.New().String(),
		Email:         req.Email,
		MobileNumber:  mobile,
		PasswordHash:  hashedPassword,
		Status:        authnentityv1.UserStatus_USER_STATUS_PENDING_VERIFICATION,
		UserType:      userType,
		EmailVerified: false,
	}

	if err := s.userRepo.CreateFull(ctx, user); err != nil {
		appLogger.Errorf("RegisterEmailUser: failed to create user: %v", err)
		logger.Errorf("failed to create user: %v", err)
		return nil, errors.New("failed to create user")
	}

	// Send email verification OTP
	otpResp, err := s.otpService.SendEmailOTP(ctx, &authnservicev1.SendEmailOTPRequest{
		Email: req.Email,
		Type:  "email_verification",
	})
	if err != nil {
		appLogger.Errorf("RegisterEmailUser: failed to send verification email to %s: %v", maskEmail(req.Email), err)
		// Don't fail registration; user can resend via SendEmailOTP
		return &authnservicev1.RegisterEmailUserResponse{
			UserId:                user.UserId,
			Message:               "Registration successful. Failed to send verification email - please use resend.",
			VerificationEmailSent: false,
		}, nil
	}

	// Publish event
	_ = s.eventPublisher.PublishUserRegistered(ctx, user.UserId, mobile, req.Email, reqMeta.IPAddress, "WEB", portalForUserType(userType.String()), "root")

	appLogger.Infof("RegisterEmailUser: user %s (%s) registered, verification email sent from IP %s", user.UserId, maskEmail(req.Email), reqMeta.IPAddress)

	return &authnservicev1.RegisterEmailUserResponse{
		UserId:                user.UserId,
		Message:               "Registration successful. Please check your email for the verification code.",
		VerificationEmailSent: true,
		OtpId:                 otpResp.OtpId,
		OtpExpiresInSeconds:   otpResp.ExpiresInSeconds,
	}, nil
}

// SendEmailOTP sends an OTP to an email address.
// Validates that the email belongs to a BUSINESS_BENEFICIARY or SYSTEM_USER
// (except for email_verification type during registration).
func (s *AuthService) SendEmailOTP(ctx context.Context, req *authnservicev1.SendEmailOTPRequest) (*authnservicev1.SendEmailOTPResponse, error) {
	reqMeta := s.metadata.ExtractAll(ctx)

	// For email_login and password_reset_email: validate user exists and is correct type
	if req.Type == "email_login" || req.Type == "password_reset_email" {
		user, err := s.userRepo.GetByEmail(ctx, req.Email)
		if err != nil {
			// Return generic message to avoid email enumeration
			appLogger.Warnf("SendEmailOTP: email not found %s from IP %s", maskEmail(req.Email), reqMeta.IPAddress)
			return &authnservicev1.SendEmailOTPResponse{
				Message: "If this email is registered, an OTP has been sent.",
			}, nil
		}
		if !isEmailAuthUser(user.UserType) {
			appLogger.Warnf("SendEmailOTP: user_type %s not allowed for email OTP from IP %s", user.UserType, reqMeta.IPAddress)
			return &authnservicev1.SendEmailOTPResponse{
				Message: "If this email is registered, an OTP has been sent.",
			}, nil
		}
		// For email_login, email must be verified first
		if req.Type == "email_login" && !user.EmailVerified {
			return nil, errors.New("email address not verified. Please verify your email first.")
		}
		// Check email auth lock
		if user.EmailLockedUntil != nil && time.Now().Before(user.EmailLockedUntil.AsTime()) {
			remaining := time.Until(user.EmailLockedUntil.AsTime()).Round(time.Second)
			return nil, fmt.Errorf("account locked due to too many failed attempts. Try again in %s", remaining)
		}
	}

	appLogger.Infof("SendEmailOTP: sending %s OTP to %s from IP %s", req.Type, maskEmail(req.Email), reqMeta.IPAddress)

	resp, err := s.otpService.SendEmailOTP(ctx, req)
	if err != nil {
		appLogger.Errorf("SendEmailOTP: failed for %s: %v", maskEmail(req.Email), err)
		return nil, err
	}

	// Publish event
	_ = s.eventPublisher.PublishEmailVerificationSent(ctx, "", req.Email, resp.OtpId, req.Type, reqMeta.IPAddress)

	return resp, nil
}

// VerifyEmail verifies an email address using an OTP.
// On success: sets email_verified=true, activates account if PENDING_VERIFICATION.
func (s *AuthService) VerifyEmail(ctx context.Context, req *authnservicev1.VerifyEmailRequest) (*authnservicev1.VerifyEmailResponse, error) {
	reqMeta := s.metadata.ExtractAll(ctx)

	// Verify the OTP
	verifyResp, err := s.otpService.VerifyOTP(ctx, &authnservicev1.VerifyOTPRequest{
		OtpId: req.OtpId,
		Code:  req.Code,
	})
	if err != nil {
		return nil, err
	}
	if !verifyResp.Verified {
		return &authnservicev1.VerifyEmailResponse{
			Verified: false,
			Message:  verifyResp.Message,
		}, nil
	}

	// Fetch the OTP to get the recipient email
	otpEntity, err := s.otpRepo.GetByID(ctx, req.OtpId)
	if err != nil {
		logger.Errorf("failed to retrieve OTP details: %v", err)
		return nil, errors.New("failed to retrieve OTP details")
	}

	// Find user by email
	user, err := s.userRepo.GetByEmail(ctx, otpEntity.Recipient)
	if err != nil {
		logger.Errorf("user not found for email: %v", err)
		return nil, errors.New("user not found for email")
	}

	// Mark email as verified
	if err := s.userRepo.UpdateEmailVerified(ctx, user.UserId); err != nil {
		logger.Errorf("failed to mark email as verified: %v", err)
		return nil, errors.New("failed to mark email as verified")
	}

	// Activate account if still pending
	if user.Status == authnentityv1.UserStatus_USER_STATUS_PENDING_VERIFICATION {
		if err := s.userRepo.UpdateStatus(ctx, user.UserId, authnentityv1.UserStatus_USER_STATUS_ACTIVE); err != nil {
			appLogger.Warnf("VerifyEmail: failed to activate user %s: %v", user.UserId, err)
		}
	}

	appLogger.Infof("VerifyEmail: email verified for user %s (%s) from IP %s", user.UserId, maskEmail(otpEntity.Recipient), reqMeta.IPAddress)

	// Publish event
	_ = s.eventPublisher.PublishEmailVerified(ctx, user.UserId, otpEntity.Recipient)

	return &authnservicev1.VerifyEmailResponse{
		Verified: true,
		UserId:   user.UserId,
		Message:  "Email verified successfully. You can now log in.",
	}, nil
}

// EmailLogin authenticates a BUSINESS_BENEFICIARY or SYSTEM_USER via email OTP.
// Always produces a SERVER_SIDE session (web portal).
// Flow: SendEmailOTP(type=email_login) → EmailLogin
func (s *AuthService) EmailLogin(ctx context.Context, req *authnservicev1.EmailLoginRequest) (*authnservicev1.EmailLoginResponse, error) {
	reqMeta := s.metadata.ExtractAll(ctx)
	maxEmailLoginAttempts := s.config.Security.LoginMaxAttempts
	if maxEmailLoginAttempts <= 0 {
		maxEmailLoginAttempts = 8
	}
	emailLockDuration := s.config.Security.LoginLockoutDuration
	if emailLockDuration <= 0 {
		emailLockDuration = 10 * time.Minute
	}

	// 1. Find user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		appLogger.Warnf("EmailLogin: email not found %s from IP %s", maskEmail(req.Email), reqMeta.IPAddress)
		return nil, errors.New("invalid credentials")
	}

	// 2. Enforce user_type restriction
	if !isEmailAuthUser(user.UserType) {
		appLogger.Warnf("EmailLogin: rejected user_type=%s for %s from IP %s", user.UserType, maskEmail(req.Email), reqMeta.IPAddress)
		return nil, errors.New("email login is not available for this account type")
	}

	// 3. Check account status
	if user.Status == authnentityv1.UserStatus_USER_STATUS_SUSPENDED {
		return nil, errors.New("account is suspended")
	}
	if user.Status == authnentityv1.UserStatus_USER_STATUS_DELETED {
		return nil, errors.New("invalid credentials")
	}
	if user.Status == authnentityv1.UserStatus_USER_STATUS_PENDING_VERIFICATION {
		return nil, errors.New("email not verified. Please verify your email first.")
	}

	// 4. Check email auth lock
	if user.EmailLockedUntil != nil && time.Now().Before(user.EmailLockedUntil.AsTime()) {
		remaining := time.Until(user.EmailLockedUntil.AsTime()).Round(time.Second)
		appLogger.Warnf("EmailLogin: account locked for %s until %s from IP %s", maskEmail(req.Email), user.EmailLockedUntil.AsTime(), reqMeta.IPAddress)
		return nil, fmt.Errorf("account locked due to too many failed attempts. Try again in %s", remaining)
	}

	// 5. Verify email must be confirmed
	if !user.EmailVerified {
		return nil, errors.New("email address not verified. Please verify your email first.")
	}

	// 6. Verify OTP
	verifyResp, err := s.otpService.VerifyOTP(ctx, &authnservicev1.VerifyOTPRequest{
		OtpId: req.OtpId,
		Code:  req.Code,
	})
	if err != nil {
		return nil, err
	}
	if !verifyResp.Verified {
		// Increment failed attempts and potentially lock
		newAttempts, _ := s.userRepo.IncrementEmailLoginAttempts(ctx, user.UserId)
		appLogger.Warnf("EmailLogin: invalid OTP for %s, attempts=%d from IP %s", maskEmail(req.Email), newAttempts, reqMeta.IPAddress)

		if int(newAttempts) >= maxEmailLoginAttempts {
			_ = s.userRepo.LockEmailAuth(ctx, user.UserId, emailLockDuration)
			appLogger.Warnf("EmailLogin: account locked for %s after %d failed attempts from IP %s", maskEmail(req.Email), newAttempts, reqMeta.IPAddress)
			return nil, fmt.Errorf("too many failed attempts. Account locked for %s", emailLockDuration)
		}

		remaining := maxEmailLoginAttempts - int(newAttempts)
		return nil, fmt.Errorf("invalid OTP. %d attempts remaining", remaining)
	}

	// 7. Reset failed attempt counter on success
	_ = s.userRepo.ResetEmailLoginAttempts(ctx, user.UserId)

	// 8. Create SERVER_SIDE session (always for web portal email login)
	serverSession, err := s.tokenService.GenerateServerSideSession(
		ctx,
		user.UserId,
		req.DeviceId,
		authnentityv1.DeviceType_DEVICE_TYPE_WEB,
		reqMeta.IPAddress,
		reqMeta.UserAgent,
	)
	if err != nil {
		appLogger.Errorf("EmailLogin: failed to create server-side session for %s: %v", user.UserId, err)
		logger.Errorf("failed to create session: %v", err)
		return nil, errors.New("failed to create session")
	}

	// 9. Update last login
	_ = s.userRepo.UpdateLastLogin(ctx, user.UserId, "SERVER_SIDE")

	appLogger.Infof("EmailLogin: successful login for user %s (%s) session=%s from IP %s", user.UserId, maskEmail(req.Email), serverSession.SessionID, reqMeta.IPAddress)

	// 10. Publish event
	_ = s.eventPublisher.PublishEmailLoginSucceeded(ctx, user.UserId, serverSession.SessionID, req.Email, user.UserType.String(), reqMeta.IPAddress, reqMeta.UserAgent, req.DeviceName)

	return &authnservicev1.EmailLoginResponse{
		UserId:       user.UserId,
		SessionId:    serverSession.SessionID,
		SessionToken: serverSession.SessionToken,
		CsrfToken:    serverSession.CSRFToken,
		User:         user,
		SessionType:  "SERVER_SIDE",
	}, nil
}

// EmailPasswordLogin authenticates an email-based user with email + password.
// Used by B2B beneficiary self-service after OTP-based first-time activation.
func (s *AuthService) EmailPasswordLogin(ctx context.Context, req *authnservicev1.EmailPasswordLoginRequest) (*authnservicev1.EmailPasswordLoginResponse, error) {
	reqMeta := s.metadata.ExtractAll(ctx)
	maxLoginAttempts := s.config.Security.LoginMaxAttempts
	if maxLoginAttempts <= 0 {
		maxLoginAttempts = 8
	}
	lockDuration := s.config.Security.LoginLockoutDuration
	if lockDuration <= 0 {
		lockDuration = 10 * time.Minute
	}

	email := strings.TrimSpace(strings.ToLower(req.GetEmail()))
	if email == "" {
		return nil, errors.New("email is required")
	}

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		appLogger.Warnf("EmailPasswordLogin: email not found %s from IP %s", maskEmail(email), reqMeta.IPAddress)
		return nil, errors.New("invalid credentials")
	}
	if !isEmailAuthUser(user.UserType) {
		return nil, errors.New("email login is not available for this account type")
	}
	if user.Status == authnentityv1.UserStatus_USER_STATUS_SUSPENDED {
		return nil, errors.New("account is suspended")
	}
	if user.Status == authnentityv1.UserStatus_USER_STATUS_DELETED {
		return nil, errors.New("invalid credentials")
	}
	if user.Status == authnentityv1.UserStatus_USER_STATUS_PENDING_VERIFICATION && !managedEmployeeActivationAllowed(user) {
		return nil, errors.New("email not verified. Please verify your email first.")
	}
	if user.LockedUntil != nil && time.Now().Before(user.LockedUntil.AsTime()) {
		remaining := time.Until(user.LockedUntil.AsTime()).Round(time.Second)
		return nil, fmt.Errorf("account is locked. Try again in %s", remaining)
	}
	if !user.EmailVerified {
		return nil, errors.New("account activation is required before email login")
	}

	valid, needsRehash, verifyErr := verifyPassword(req.GetPassword(), user.GetPasswordHash())
	if verifyErr != nil || !valid {
		attempts, _ := s.userRepo.IncrementLoginAttempts(ctx, user.GetUserId())
		if int(attempts) >= maxLoginAttempts {
			_ = s.userRepo.LockAccount(ctx, user.GetUserId(), lockDuration)
			return nil, fmt.Errorf("account is locked. Try again in %s", lockDuration)
		}
		remaining := maxLoginAttempts - int(attempts)
		return nil, fmt.Errorf("invalid credentials. %d attempts remaining", remaining)
	}
	if needsRehash {
		if newHash, err := hashPassword(req.GetPassword()); err == nil {
			_ = s.userRepo.UpdatePassword(ctx, user.GetUserId(), newHash)
		}
	}
	_ = s.userRepo.ResetLoginAttempts(ctx, user.GetUserId())

	serverSession, err := s.tokenService.GenerateServerSideSession(
		ctx,
		user.GetUserId(),
		req.GetDeviceId(),
		authnentityv1.DeviceType_DEVICE_TYPE_WEB,
		reqMeta.IPAddress,
		reqMeta.UserAgent,
	)
	if err != nil {
		return nil, errors.New("failed to create session")
	}
	_ = s.userRepo.UpdateLastLogin(ctx, user.GetUserId(), "SERVER_SIDE")
	_ = s.eventPublisher.PublishUserLoggedIn(ctx, user.GetUserId(), serverSession.SessionID, "SERVER_SIDE", reqMeta.IPAddress, "WEB", reqMeta.UserAgent)

	return &authnservicev1.EmailPasswordLoginResponse{
		UserId:                 user.GetUserId(),
		SessionId:              serverSession.SessionID,
		SessionToken:           serverSession.SessionToken,
		CsrfToken:              serverSession.CSRFToken,
		User:                   user,
		SessionType:            "SERVER_SIDE",
		PasswordChangeRequired: user.GetPasswordChangeRequired(),
	}, nil
}

// RequestPasswordResetByEmail sends a password reset OTP to the user's email.
// Requires: user exists, email is verified, user_type allows email auth.
func (s *AuthService) RequestPasswordResetByEmail(ctx context.Context, req *authnservicev1.RequestPasswordResetByEmailRequest) (*authnservicev1.RequestPasswordResetByEmailResponse, error) {
	reqMeta := s.metadata.ExtractAll(ctx)

	// Always return success to prevent email enumeration
	genericResp := &authnservicev1.RequestPasswordResetByEmailResponse{
		Message: "If this email is registered, a password reset code has been sent.",
	}

	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		appLogger.Infof("RequestPasswordResetByEmail: email not found %s from IP %s (no-op)", maskEmail(req.Email), reqMeta.IPAddress)
		return genericResp, nil
	}

	if !isEmailAuthUser(user.UserType) {
		appLogger.Warnf("RequestPasswordResetByEmail: user_type %s not allowed from IP %s", user.UserType, reqMeta.IPAddress)
		return genericResp, nil
	}

	if !user.EmailVerified && !managedEmployeeActivationAllowed(user) {
		appLogger.Warnf("RequestPasswordResetByEmail: email not verified for %s from IP %s", maskEmail(req.Email), reqMeta.IPAddress)
		return genericResp, nil
	}

	otpResp, err := s.otpService.SendEmailOTP(ctx, &authnservicev1.SendEmailOTPRequest{
		Email: req.Email,
		Type:  "password_reset_email",
	})
	if err != nil {
		appLogger.Errorf("RequestPasswordResetByEmail: failed to send reset email to %s: %v", maskEmail(req.Email), err)
		return genericResp, nil
	}

	appLogger.Infof("RequestPasswordResetByEmail: reset OTP sent to %s from IP %s", maskEmail(req.Email), reqMeta.IPAddress)
	_ = s.eventPublisher.PublishPasswordResetByEmailRequested(ctx, user.UserId, req.Email, otpResp.OtpId, reqMeta.IPAddress)

	return &authnservicev1.RequestPasswordResetByEmailResponse{
		OtpId:            otpResp.OtpId,
		Message:          "Password reset code sent to your email address.",
		ExpiresInSeconds: otpResp.ExpiresInSeconds,
	}, nil
}

// ResetPasswordByEmail completes password reset using an email OTP.
func (s *AuthService) ResetPasswordByEmail(ctx context.Context, req *authnservicev1.ResetPasswordByEmailRequest) (*authnservicev1.ResetPasswordByEmailResponse, error) {
	reqMeta := s.metadata.ExtractAll(ctx)

	// Find user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		appLogger.Warnf("ResetPasswordByEmail: email not found %s from IP %s", maskEmail(req.Email), reqMeta.IPAddress)
		return nil, errors.New("invalid request")
	}

	if !isEmailAuthUser(user.UserType) {
		return nil, errors.New("email password reset is not available for this account type")
	}

	// Verify OTP
	verifyResp, err := s.otpService.VerifyOTP(ctx, &authnservicev1.VerifyOTPRequest{
		OtpId: req.OtpId,
		Code:  req.OtpCode,
	})
	if err != nil || !verifyResp.Verified {
		appLogger.Warnf("ResetPasswordByEmail: invalid OTP for %s from IP %s", maskEmail(req.Email), reqMeta.IPAddress)
		return nil, errors.New("invalid or expired OTP")
	}

	// Enforce password policy and hash with Argon2id.
	if err := validatePasswordStrength(req.NewPassword); err != nil {
		logger.Errorf("weak password: %v", err)
		return nil, errors.New("weak password")
	}
	hashedPassword, err := hashPassword(req.NewPassword)
	if err != nil {
		logger.Errorf("failed to hash password: %v", err)
		return nil, errors.New("failed to hash password")
	}

	// Update password
	if err := s.userRepo.UpdatePassword(ctx, user.UserId, hashedPassword); err != nil {
		logger.Errorf("failed to update password: %v", err)
		return nil, errors.New("failed to update password")
	}
	if !user.GetEmailVerified() {
		if err := s.userRepo.UpdateEmailVerified(ctx, user.GetUserId()); err != nil {
			logger.Errorf("failed to mark email as verified after password reset: %v", err)
			return nil, errors.New("failed to update password")
		}
	}
	if user.GetStatus() == authnentityv1.UserStatus_USER_STATUS_PENDING_VERIFICATION {
		if err := s.userRepo.UpdateStatus(ctx, user.GetUserId(), authnentityv1.UserStatus_USER_STATUS_ACTIVE); err != nil {
			logger.Errorf("failed to activate user after password reset: %v", err)
			return nil, errors.New("failed to update password")
		}
	}
	if err := s.userRepo.SetPasswordChangeRequired(ctx, user.UserId, false); err != nil {
		logger.Errorf("failed to clear password_change_required: %v", err)
		return nil, errors.New("failed to update password")
	}

	// Revoke all sessions (force re-login)
	_ = s.sessionRepo.RevokeAllByUserID(ctx, user.UserId, "")

	appLogger.Infof("ResetPasswordByEmail: password reset for user %s (%s) from IP %s", user.UserId, maskEmail(req.Email), reqMeta.IPAddress)

	return &authnservicev1.ResetPasswordByEmailResponse{
		Message: "Password reset successfully. Please log in with your new password.",
	}, nil
}

// ProvisionEmployeeUser provisions a managed B2B beneficiary user for employee self-service.
// The account is created once, starts in PENDING_VERIFICATION, and is activated later
// through RequestPasswordResetByEmail + ResetPasswordByEmail.
func (s *AuthService) ProvisionEmployeeUser(ctx context.Context, req *authnservicev1.ProvisionEmployeeUserRequest) (*authnservicev1.ProvisionEmployeeUserResponse, error) {
	email := strings.TrimSpace(strings.ToLower(req.GetEmail()))
	fullName := strings.TrimSpace(req.GetFullName())
	if email == "" {
		return nil, errors.New("email is required")
	}
	if fullName == "" {
		return nil, errors.New("full_name is required")
	}

	if existing, err := s.userRepo.GetByEmail(ctx, email); err == nil && existing != nil {
		if existing.GetUserType() != authnentityv1.UserType_USER_TYPE_B2B_BENEFICIARY {
			return nil, errors.New("email is already used by another account type")
		}
		return &authnservicev1.ProvisionEmployeeUserResponse{
			UserId:                 existing.GetUserId(),
			Created:                false,
			PasswordChangeRequired: existing.GetPasswordChangeRequired(),
		}, nil
	}

	randomSeed := uuid.NewString() + "Aa1!"
	hashedPassword, err := hashPassword(randomSeed)
	if err != nil {
		return nil, errors.New("failed to create employee account")
	}

	user := &authnentityv1.User{
		UserId:                 uuid.NewString(),
		Email:                  email,
		MobileNumber:           strings.TrimSpace(req.GetMobileNumber()),
		PasswordHash:           hashedPassword,
		Status:                 authnentityv1.UserStatus_USER_STATUS_PENDING_VERIFICATION,
		UserType:               authnentityv1.UserType_USER_TYPE_B2B_BENEFICIARY,
		EmailVerified:          false,
		PasswordChangeRequired: true,
	}
	if err := s.userRepo.CreateFull(ctx, user); err != nil {
		return nil, errors.New("failed to create employee account")
	}

	reqMeta := s.metadata.ExtractAll(ctx)
	_ = s.eventPublisher.PublishUserRegistered(ctx, user.GetUserId(), user.GetMobileNumber(), user.GetEmail(), reqMeta.IPAddress, "WEB", portalForUserType(user.GetUserType().String()), "root")

	return &authnservicev1.ProvisionEmployeeUserResponse{
		UserId:                 user.GetUserId(),
		Created:                true,
		PasswordChangeRequired: true,
	}, nil
}

// parseUserType converts string user_type to enum
func parseUserType(userTypeStr string) authnentityv1.UserType {
	switch userTypeStr {
	case "B2C_CUSTOMER":
		return authnentityv1.UserType_USER_TYPE_B2C_CUSTOMER
	case "AGENT":
		return authnentityv1.UserType_USER_TYPE_AGENT
	case "BUSINESS_BENEFICIARY":
		return authnentityv1.UserType_USER_TYPE_BUSINESS_BENEFICIARY
	case "SYSTEM_USER":
		return authnentityv1.UserType_USER_TYPE_SYSTEM_USER
	case "B2B_BENEFICIARY":
		return authnentityv1.UserType_USER_TYPE_B2B_BENEFICIARY
	case "B2B_ORG_ADMIN":
		return authnentityv1.UserType_USER_TYPE_B2B_ORG_ADMIN
	default:
		return authnentityv1.UserType_USER_TYPE_UNSPECIFIED
	}
}
