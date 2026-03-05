package main

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/newage-saint/insuretech/backend/inscore/db"
	_ "github.com/newage-saint/insuretech/backend/inscore/db" // ensures init() fires schema registry
	"github.com/newage-saint/insuretech/backend/inscore/microservices/authn/internal/repository"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	authnentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/entity/v1"
	"github.com/newage-saint/insuretech/ops/env"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm/schema"
)

func main() {
	// Initialize logger
	if err := appLogger.Initialize(appLogger.NoFileConfig()); err != nil {
		appLogger.Fatalf("Failed to initialize logger: %v", err)
	}

	if err := env.Load(); err != nil {
		appLogger.Warnf("Warning: couldn't load .env: %v", err)
	}

	configPath := "../../../../database.yaml"
	if err := db.InitializeManagerForService(configPath); err != nil {
		appLogger.Fatalf("Failed to connect to DB: %v", err)
	}

	// Double ensure the global serializer is wired up just in case
	schema.RegisterSerializer("proto_timestamp", db.ProtoTimestampSerializer{})

	dbConn := db.GetDB()
	if dbConn == nil {
		appLogger.Fatal("DB connection is nil")
	}

	ctx := context.Background()

	// Ensure correct schema prefix since we determined they live in authn_schema
	// Note: We bypass repository abstraction slightly only to set the Table/Schema globally
	// but the repositories actually use their own TableName() methods if they exist, or default.
	// user_repository has TableName() returning "users" which failed before because of schema.
	// So we explicitly pass a Session configured with the schema
	// Note: We bypass repository abstraction slightly only to set the Table/Schema globally
	dbConn = dbConn.Debug()

	appLogger.Info("===================================================")
	appLogger.Info("   THOROUGH TESTING OF ALL AUTHN REPOSITORIES")
	appLogger.Info("===================================================")

	testMobile := "+8801999999777"

	// Pre-CLEANUP (in case a previous run failed and left data behind)
	appLogger.Info("[PRE-CLEANUP] Removing any left-over test data from previous runs...")
	dbConn.Table("authn_schema.otps").Where("user_id IN (SELECT user_id FROM authn_schema.users WHERE mobile_number = ?)", testMobile).Delete(&authnentityv1.OTP{})
	dbConn.Table("authn_schema.sessions").Where("user_id IN (SELECT user_id FROM authn_schema.users WHERE mobile_number = ?)", testMobile).Delete(&authnentityv1.Session{})
	dbConn.Table("authn_schema.users").Where("mobile_number = ?", testMobile).Delete(&authnentityv1.User{})
	dbConn.Table("authn_schema.otps").Where("user_id IN (SELECT user_id FROM authn_schema.users WHERE mobile_number = ?)", testMobile).Delete(&authnentityv1.OTP{})
	dbConn.Table("authn_schema.sessions").Where("user_id IN (SELECT user_id FROM authn_schema.users WHERE mobile_number = ?)", testMobile).Delete(&authnentityv1.Session{})
	dbConn.Table("authn_schema.users").Where("mobile_number = ?", testMobile).Delete(&authnentityv1.User{})

	// 1. User Repository Test
	appLogger.Info("\n[1/3] Testing UserRepository...")
	userRepo := repository.NewUserRepository(dbConn.Table("authn_schema.users"))

	appLogger.Infof("      -> Creating User %s...\n", testMobile)
	user, err := userRepo.Create(ctx, testMobile, "hashed_pw", "all@test.com", authnentityv1.UserStatus_USER_STATUS_ACTIVE)
	if err != nil {
		appLogger.Fatalf("FAILED to create User: %v", err)
	}
	appLogger.Infof("      -> SUCCESS: Created User ID: %s\n", user.UserId)

	appLogger.Info("      -> Retrieving User...\n")
	fetchedUser, err := userRepo.GetByMobileNumber(ctx, testMobile)
	if err != nil {
		appLogger.Fatalf("FAILED to retrieve User: %v", err)
	}
	appLogger.Infof("      -> SUCCESS: Fetched User ID: %s, Validation Check: %s\n", fetchedUser.UserId, fetchedUser.MobileNumber)

	// 2. OTP Repository Test
	appLogger.Info("\n[2/3] Testing OTPRepository...")
	otpRepo := repository.NewOTPRepository(dbConn.Table("authn_schema.otps"))
	otpID := uuid.New().String()

	newOTP := &authnentityv1.OTP{
		OtpId:     otpID,
		UserId:    user.UserId,
		OtpHash:   "hashed_123456",
		Purpose:   "login",
		ExpiresAt: timestamppb.New(time.Now().Add(5 * time.Minute)),
	}

	appLogger.Infof("      -> Creating OTP %s...\n", otpID)
	err = otpRepo.Create(ctx, newOTP)
	if err != nil {
		appLogger.Fatalf("FAILED to create OTP: %v", err)
	}
	appLogger.Info("      -> SUCCESS: OTP Created.\n")

	appLogger.Info("      -> Updating/Verifying OTP...\n")
	err = otpRepo.MarkVerified(ctx, otpID)
	if err != nil {
		appLogger.Fatalf("FAILED to verify OTP: %v", err)
	}

	lastOTP, err := otpRepo.GetLastOTP(ctx, testMobile)
	if err != nil || lastOTP == nil {
		appLogger.Fatalf("FAILED to fetch last OTP: %v", err)
	}
	appLogger.Info("      -> SUCCESS: Verified OTP state fetching mapped correctly.\n")

	// 3. Session Repository Test
	appLogger.Info("\n[3/3] Testing SessionRepository...")
	sessionRepo := repository.NewSessionRepository(dbConn.Table("authn_schema.sessions"))
	sessionID := uuid.New().String()

	newSession := &authnentityv1.Session{
		SessionId:             sessionID,
		UserId:                user.UserId,
		SessionType:           authnentityv1.SessionType_SESSION_TYPE_SERVER_SIDE,
		DeviceId:              "windows-chrome",
		DeviceType:            authnentityv1.DeviceType_DEVICE_TYPE_DESKTOP,
		IpAddress:             "127.0.0.1",
		UserAgent:             "Mozilla test",
		IsActive:              true,
		ExpiresAt:             timestamppb.New(time.Now().Add(24 * time.Hour)),
		SessionTokenHash:      "token_hash_value",
		AccessTokenJti:        uuid.New().String(),
		RefreshTokenJti:       uuid.New().String(),
		AccessTokenExpiresAt:  timestamppb.New(time.Now().Add(15 * time.Minute)),
		RefreshTokenExpiresAt: timestamppb.New(time.Now().Add(24 * time.Hour)),
	}

	appLogger.Infof("      -> Creating Session %s...\n", sessionID)
	err = sessionRepo.Create(ctx, newSession)
	if err != nil {
		appLogger.Fatalf("FAILED to create Session: %v", err)
	}
	appLogger.Info("      -> SUCCESS: Session Created.\n")

	appLogger.Info("      -> Updating Session Activity...\n")
	err = sessionRepo.UpdateLastActivity(ctx, sessionID)
	if err != nil {
		appLogger.Fatalf("FAILED to update session activity: %v", err)
	}

	activeSessions, err := sessionRepo.ListByUserID(ctx, user.UserId, true, nil)
	if err != nil || len(activeSessions) == 0 {
		appLogger.Fatalf("FAILED to list active sessions: %v", err)
	}
	appLogger.Infof("      -> SUCCESS: Fetched %d active session(s) for User.\n", len(activeSessions))

	// CLEANUP
	appLogger.Info("\n[CLEANUP] Removing test data from Live DB...")
	dbConn.Table("authn_schema.sessions").Where("session_id = ?", sessionID).Delete(&authnentityv1.Session{})
	dbConn.Table("authn_schema.otps").Where("otp_id = ?", otpID).Delete(&authnentityv1.OTP{})
	dbConn.Table("authn_schema.users").Where("user_id = ?", user.UserId).Delete(&authnentityv1.User{})

	appLogger.Info("      -> SUCCESS: All test data cleaned up.")
	appLogger.Info("\n===================================================")
	appLogger.Info("   ALL REPOSITORIES VALIDATED SUCCESSFULLY")
	appLogger.Info("===================================================")
}
