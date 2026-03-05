package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	authnentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/entity/v1"
	authnservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func createEmailAuthUser(t *testing.T, svc *AuthService, ctx context.Context, email, mobile string, userType authnentityv1.UserType, status authnentityv1.UserStatus, emailVerified bool, password string) string {
	t.Helper()
	pwd, err := hashPassword(password)
	require.NoError(t, err)
	u := &authnentityv1.User{
		UserId:        uuid.NewString(),
		Email:         email,
		MobileNumber:  mobile,
		PasswordHash:  pwd,
		UserType:      userType,
		Status:        status,
		EmailVerified: emailVerified,
	}
	require.NoError(t, svc.userRepo.CreateFull(ctx, u))
	return u.UserId
}

func createEmailOTP(t *testing.T, svc *AuthService, ctx context.Context, otpID, userID, recipient, purpose, code string, expiresAt time.Time) {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
	require.NoError(t, err)
	err = svc.otpRepo.Create(ctx, &authnentityv1.OTP{
		OtpId:      otpID,
		UserId:     userID,
		OtpHash:    string(hash),
		Purpose:    purpose,
		Recipient:  recipient,
		Channel:    "email",
		ExpiresAt:  timestamppb.New(expiresAt),
		Verified:   false,
		Attempts:   0,
		DlrStatus:  "DELIVERED",
		DeviceType: "WEB",
		IpAddress:  "127.0.0.1",
	})
	require.NoError(t, err)
}

func TestAuthService_LiveDB_EmailAuthFlows(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live DB test")
	}
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		"x-forwarded-for", "127.0.0.1",
		"user-agent", "email-auth-live-test",
	))
	dbConn := testServiceLiveDB(t)
	svc := buildLiveAuthService(t, dbConn)
	// email client intentionally nil in buildLiveAuthService; coverage includes error-tolerant paths.

	emailReg := "email_reg_" + uuid.NewString()[:8] + "@example.com"
	mobileReg := fmt.Sprintf("+8801%09d", time.Now().UnixNano()%1_000_000_000)
	regResp, err := svc.RegisterEmailUser(ctx, &authnservicev1.RegisterEmailUserRequest{
		Email:        emailReg,
		MobileNumber: mobileReg,
		Password:     "Str0ng!Email1",
		UserType:     "BUSINESS_BENEFICIARY",
	})
	require.NoError(t, err)
	require.NotEmpty(t, regResp.UserId)
	require.False(t, regResp.VerificationEmailSent)
	t.Cleanup(func() { cleanupLiveAuthnUser(t, dbConn, regResp.UserId) })

	// SendEmailOTP path: unverified email_login should fail before SMTP send.
	_, err = svc.SendEmailOTP(ctx, &authnservicev1.SendEmailOTPRequest{
		Email: emailReg,
		Type:  "email_login",
	})
	require.Error(t, err)

	// VerifyEmail success path using a pre-inserted OTP
	verifyCode := "123456"
	verifyOTPID := uuid.NewString()
	createEmailOTP(t, svc, ctx, verifyOTPID, regResp.UserId, emailReg, "email_verification", verifyCode, time.Now().Add(5*time.Minute))
	vResp, err := svc.VerifyEmail(ctx, &authnservicev1.VerifyEmailRequest{
		OtpId: verifyOTPID,
		Code:  verifyCode,
	})
	require.NoError(t, err)
	require.True(t, vResp.Verified)

	verifiedUser, err := svc.userRepo.GetByID(ctx, regResp.UserId)
	require.NoError(t, err)
	require.True(t, verifiedUser.EmailVerified)

	// EmailLogin success path
	loginCode := "654321"
	loginOTPID := uuid.NewString()
	createEmailOTP(t, svc, ctx, loginOTPID, regResp.UserId, emailReg, "email_login", loginCode, time.Now().Add(5*time.Minute))
	loginResp, err := svc.EmailLogin(ctx, &authnservicev1.EmailLoginRequest{
		Email:    emailReg,
		OtpId:    loginOTPID,
		Code:     loginCode,
		DeviceId: "web-email-1",
	})
	require.NoError(t, err)
	require.Equal(t, "SERVER_SIDE", loginResp.SessionType)
	require.NotEmpty(t, loginResp.SessionId)

	// RequestPasswordResetByEmail generic-not-found path (anti-enumeration).
	resetReqResp, err := svc.RequestPasswordResetByEmail(ctx, &authnservicev1.RequestPasswordResetByEmailRequest{
		Email: "missing_" + uuid.NewString()[:8] + "@example.com",
	})
	require.NoError(t, err)
	require.Contains(t, resetReqResp.Message, "If this email is registered")

	// ResetPasswordByEmail success path
	resetCode := "112233"
	resetOTPID := uuid.NewString()
	createEmailOTP(t, svc, ctx, resetOTPID, regResp.UserId, emailReg, "password_reset_email", resetCode, time.Now().Add(5*time.Minute))
	_, err = svc.ResetPasswordByEmail(ctx, &authnservicev1.ResetPasswordByEmailRequest{
		Email:       emailReg,
		OtpId:       resetOTPID,
		OtpCode:     resetCode,
		NewPassword: "Str0ng!Email2",
	})
	require.NoError(t, err)

	uAfter, err := svc.userRepo.GetByID(ctx, regResp.UserId)
	require.NoError(t, err)
	ok, _, err := verifyPassword("Str0ng!Email2", uAfter.PasswordHash)
	require.NoError(t, err)
	require.True(t, ok)
}
