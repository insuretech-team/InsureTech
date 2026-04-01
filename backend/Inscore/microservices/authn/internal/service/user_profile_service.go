package service

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/grpcmeta"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	authnentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/entity/v1"
	authnservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
	storageentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/storage/entity/v1"
	storageservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/storage/service/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ─── UserProfile ──────────────────────────────────────────────────────────────

// CreateUserProfile creates a new user profile for the given user.
func (s *AuthService) CreateUserProfile(ctx context.Context, req *authnservicev1.CreateUserProfileRequest) (*authnservicev1.CreateUserProfileResponse, error) {
	if s.userProfileRepo == nil {
		return nil, errors.New("user profile repository not configured")
	}

	if _, err := s.userRepo.GetByID(ctx, req.UserId); err != nil {
		logger.Errorf("user not found: %v", err)
		return nil, errors.New("user not found")
	}

	if existing, err := s.userProfileRepo.GetByUserID(ctx, req.UserId); err == nil && existing != nil {
		return nil, errors.New("profile already exists for this user")
	}

	profile := &authnentityv1.UserProfile{
		UserId:                 uuid.NewString(), // gorm primaryKey is user_id per proto tag
		FullName:               req.FullName,
		DateOfBirth:            req.DateOfBirth,
		Gender:                 parseGender(req.Gender),
		Occupation:             req.Occupation,
		Employer:               req.Employer,
		AddressLine1:           req.AddressLine1,
		AddressLine2:           req.AddressLine2,
		City:                   req.City,
		District:               req.District,
		Division:               req.Division,
		Country:                req.Country,
		PostalCode:             req.PostalCode,
		PermanentAddress:       req.PermanentAddress,
		NidNumber:              req.NidNumber,
		MaritalStatus:          req.MaritalStatus,
		EmergencyContactName:   req.EmergencyContactName,
		EmergencyContactNumber: req.EmergencyContactNumber,
		KycVerified:            false,
		CreatedAt:              timestamppb.Now(),
		UpdatedAt:              timestamppb.Now(),
	}
	// Override UserId with the requesting user's actual ID (profile PK = user_id).
	profile.UserId = req.UserId

	if err := s.userProfileRepo.Create(ctx, profile); err != nil {
		appLogger.Errorf("CreateUserProfile: failed to create profile for user %s: %v", req.UserId, err)
		logger.Errorf("failed to create profile: %v", err)
		return nil, errors.New("failed to create profile")
	}

	appLogger.Infof("CreateUserProfile: created profile for user %s", req.UserId)

	return &authnservicev1.CreateUserProfileResponse{
		Profile: profile,
		Message: "Profile created successfully",
	}, nil
}

// GetUserProfile retrieves the profile for a user.
// If no profile row exists yet (new user who has never visited My Profile),
// a minimal empty profile is created on-the-fly and returned so the
// Settings page always has something to display / save into.
func (s *AuthService) GetUserProfile(ctx context.Context, req *authnservicev1.GetUserProfileRequest) (*authnservicev1.GetUserProfileResponse, error) {
	if s.userProfileRepo == nil {
		return nil, errors.New("user profile repository not configured")
	}

	profile, err := s.userProfileRepo.GetByUserID(ctx, req.UserId)
	if err != nil {
		// Profile does not exist yet — auto-create a minimal empty one so
		// the My Profile tab renders correctly on first visit.
		appLogger.Infof("GetUserProfile: no profile for user %s, auto-creating empty profile", req.UserId)
		// date_of_birth is NOT NULL in the DB — supply a zero-epoch sentinel.
		// The user can update it to a real value from the My Profile form.
		zeroDOB := timestamppb.New(time.Time{})
		profile = &authnentityv1.UserProfile{
			UserId:      req.UserId,
			DateOfBirth: zeroDOB,
			CreatedAt:   timestamppb.Now(),
			UpdatedAt:   timestamppb.Now(),
		}
		if createErr := s.userProfileRepo.Create(ctx, profile); createErr != nil {
			// Creation failed — log and surface the original not-found error
			// rather than the internal create error so callers get a clean message.
			appLogger.Errorf("GetUserProfile: auto-create profile failed for user %s: %v", req.UserId, createErr)
			logger.Errorf("profile not found: %v", err)
			return nil, errors.New("profile not found")
		}
		appLogger.Infof("GetUserProfile: auto-created empty profile for user %s", req.UserId)
	}

	return &authnservicev1.GetUserProfileResponse{Profile: profile}, nil
}

// UpdateUserProfile updates existing profile fields (non-zero values only).
// If no profile row exists yet (new user saving for the first time), one is
// created with the supplied fields so the upsert always succeeds.
func (s *AuthService) UpdateUserProfile(ctx context.Context, req *authnservicev1.UpdateUserProfileRequest) (*authnservicev1.UpdateUserProfileResponse, error) {
	if s.userProfileRepo == nil {
		return nil, errors.New("user profile repository not configured")
	}

	existing, err := s.userProfileRepo.GetByUserID(ctx, req.UserId)
	if err != nil || existing == nil {
		// No profile yet — create a new one seeded with the incoming fields.
		appLogger.Infof("UpdateUserProfile: no profile for user %s, auto-creating on first save", req.UserId)
		// date_of_birth is NOT NULL in DB — use zero-epoch sentinel for auto-create.
		zeroDOB := timestamppb.New(time.Time{})
		existing = &authnentityv1.UserProfile{
			UserId:      req.UserId,
			DateOfBirth: zeroDOB,
			CreatedAt:   timestamppb.Now(),
			UpdatedAt:   timestamppb.Now(),
		}
		if createErr := s.userProfileRepo.Create(ctx, existing); createErr != nil {
			appLogger.Errorf("UpdateUserProfile: auto-create failed for user %s: %v", req.UserId, createErr)
			return nil, errors.New("failed to create profile")
		}
		appLogger.Infof("UpdateUserProfile: auto-created profile for user %s", req.UserId)
	}

	if req.FullName != "" {
		existing.FullName = req.FullName
	}
	if req.DateOfBirth != nil {
		existing.DateOfBirth = req.DateOfBirth
	}
	if req.Gender != "" {
		existing.Gender = parseGender(req.Gender)
	}
	if req.Occupation != "" {
		existing.Occupation = req.Occupation
	}
	if req.Employer != "" {
		existing.Employer = req.Employer
	}
	if req.AddressLine1 != "" {
		existing.AddressLine1 = req.AddressLine1
	}
	if req.AddressLine2 != "" {
		existing.AddressLine2 = req.AddressLine2
	}
	if req.City != "" {
		existing.City = req.City
	}
	if req.District != "" {
		existing.District = req.District
	}
	if req.Division != "" {
		existing.Division = req.Division
	}
	if req.Country != "" {
		existing.Country = req.Country
	}
	if req.PostalCode != "" {
		existing.PostalCode = req.PostalCode
	}
	if req.PermanentAddress != "" {
		existing.PermanentAddress = req.PermanentAddress
	}
	if req.NidNumber != "" {
		existing.NidNumber = req.NidNumber
	}
	if req.MaritalStatus != "" {
		existing.MaritalStatus = req.MaritalStatus
	}
	if req.EmergencyContactName != "" {
		existing.EmergencyContactName = req.EmergencyContactName
	}
	if req.EmergencyContactNumber != "" {
		existing.EmergencyContactNumber = req.EmergencyContactNumber
	}
	if req.ProfilePhotoUrl != "" {
		existing.ProfilePhotoUrl = req.ProfilePhotoUrl
	}
	existing.UpdatedAt = timestamppb.Now()

	if err := s.userProfileRepo.Update(ctx, existing); err != nil {
		appLogger.Errorf("UpdateUserProfile: failed for user %s: %v", req.UserId, err)
		logger.Errorf("failed to update profile: %v", err)
		return nil, errors.New("failed to update profile")
	}

	appLogger.Infof("UpdateUserProfile: updated profile for user %s", req.UserId)

	return &authnservicev1.UpdateUserProfileResponse{
		Profile: existing,
		Message: "Profile updated successfully",
	}, nil
}

// parseGender converts a string like "MALE" or "GENDER_MALE" to the proto enum.
func parseGender(s string) authnentityv1.Gender {
	// Try exact enum name first (e.g. "GENDER_MALE")
	if v, ok := authnentityv1.Gender_value[s]; ok {
		return authnentityv1.Gender(v)
	}
	// Try with prefix (e.g. "MALE" -> "GENDER_MALE")
	if v, ok := authnentityv1.Gender_value["GENDER_"+s]; ok {
		return authnentityv1.Gender(v)
	}
	return authnentityv1.Gender_GENDER_UNSPECIFIED
}

// ── GetProfilePhotoUploadURL ──────────────────────────────────────────────────

// GetProfilePhotoUploadURL generates a presigned S3/Storage upload URL for the user's profile photo.
// Tries the Storage microservice first; falls back to direct S3 presign if unavailable.
func (s *AuthService) GetProfilePhotoUploadURL(ctx context.Context, req *authnservicev1.GetProfilePhotoUploadURLRequest) (*authnservicev1.GetProfilePhotoUploadURLResponse, error) {
	contentType := req.ContentType
	if contentType == "" {
		contentType = "image/jpeg"
	}
	fileExt := ".jpg"
	if strings.Contains(contentType, "png") {
		fileExt = ".png"
	} else if strings.Contains(contentType, "webp") {
		fileExt = ".webp"
	}

	legacyFallback := func() (*authnservicev1.GetProfilePhotoUploadURLResponse, error) {
		bucket := os.Getenv("S3_BUCKET")
		if bucket == "" {
			bucket = "insuretech-user-media"
		}
		region := os.Getenv("AWS_REGION")
		if region == "" {
			region = "ap-southeast-1"
		}
		objectKey := "profile-photos/" + req.UserId + "/" + uuid.New().String() + fileExt
		fileURL := "https://" + bucket + ".s3." + region + ".amazonaws.com/" + objectKey

		loadOpts := []func(*awscfg.LoadOptions) error{awscfg.WithRegion(region)}
		accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
		secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
		if accessKey != "" && secretKey != "" {
			sessionToken := os.Getenv("AWS_SESSION_TOKEN")
			loadOpts = append(loadOpts, awscfg.WithCredentialsProvider(
				credentials.NewStaticCredentialsProvider(accessKey, secretKey, sessionToken),
			))
		}
		awsConfig, err := awscfg.LoadDefaultConfig(ctx, loadOpts...)
		if err != nil {
			return nil, errors.New("load AWS config")
		}
		s3Client := s3.NewFromConfig(awsConfig)
		presignClient := s3.NewPresignClient(s3Client)
		presignedReq, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
			Bucket:      aws.String(bucket),
			Key:         aws.String(objectKey),
			ContentType: aws.String(contentType),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = 15 * time.Minute
		})
		if err != nil {
			return nil, errors.New("generate presigned upload url")
		}
		return &authnservicev1.GetProfilePhotoUploadURLResponse{
			UploadUrl:        presignedReq.URL,
			FileUrl:          fileURL,
			ExpiresInSeconds: 900,
		}, nil
	}

	storageAddr := os.Getenv("STORAGE_SERVICE_ADDRESS")
	if storageAddr == "" {
		port := os.Getenv("STORAGE_GRPC_PORT")
		if port == "" {
			port = "50290"
		}
		storageAddr = "localhost:" + port
	}

	dialCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(dialCtx, storageAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock()) //nolint:staticcheck
	if err != nil {
		logger.Errorf("GetProfilePhotoUploadURL: dial storage service: %v", err)
		return legacyFallback()
	}
	defer func() { _ = conn.Close() }()

	tenantID := grpcmeta.TenantID(ctx, os.Getenv("DEFAULT_TENANT_ID"))
	if tenantID == "" {
		tenantID = "00000000-0000-0000-0000-000000000001"
	}

	filename := "profile-" + req.UserId + fileExt
	client := storageservicev1.NewStorageServiceClient(conn)
	uploadResp, err := client.GetUploadURL(ctx, &storageservicev1.GetUploadURLRequest{
		TenantId:         tenantID,
		Filename:         filename,
		ContentType:      contentType,
		FileType:         storageentityv1.FileType_FILE_TYPE_IMAGE,
		ExpiresInMinutes: 15,
		ReferenceId:      req.UserId,
		ReferenceType:    "USER_KYC_PROFILE",
		IsPublic:         false,
	})
	if err != nil {
		logger.Errorf("GetProfilePhotoUploadURL: storage get upload url: %v", err)
		return legacyFallback()
	}

	fileURL := ""
	if cdn := strings.TrimRight(os.Getenv("SPACES_CDN_ENDPOINT"), "/"); cdn != "" {
		fileURL = cdn + "/" + strings.TrimLeft(uploadResp.StorageKey, "/")
	} else if endpoint := strings.TrimRight(os.Getenv("SPACES_ENDPOINT"), "/"); endpoint != "" {
		fileURL = endpoint + "/" + strings.TrimLeft(uploadResp.StorageKey, "/")
	}

	return &authnservicev1.GetProfilePhotoUploadURLResponse{
		UploadUrl:        uploadResp.UploadUrl,
		FileUrl:          fileURL,
		ExpiresInSeconds: 900,
	}, nil
}

// ── UpdateNotificationPreferences ────────────────────────────────────────────

// GetNotificationPreferences returns the current notification preferences for a user.
// BUG-010 FIX: Added missing GET endpoint — previously only PATCH/update existed.
func (s *AuthService) GetNotificationPreferences(ctx context.Context, req *authnservicev1.GetNotificationPreferencesRequest) (*authnservicev1.GetNotificationPreferencesResponse, error) {
	if req.GetUserId() == "" {
		return nil, errors.New("user_id is required")
	}
	user, err := s.userRepo.GetByID(ctx, req.GetUserId())
	if err != nil {
		if err.Error() == "record not found" || err.Error() == "sql: no rows in result set" {
			return nil, errors.New("user not found")
		}
		appLogger.Errorf("GetNotificationPreferences: %v", err)
		return nil, errors.New("get notification preferences")
	}
	return &authnservicev1.GetNotificationPreferencesResponse{
		UserId:                 req.GetUserId(),
		NotificationPreference: user.NotificationPreference,
		PreferredLanguage:      user.PreferredLanguage,
	}, nil
}

// UpdateNotificationPreferences updates the user's notification channel and language preferences.
func (s *AuthService) UpdateNotificationPreferences(ctx context.Context, req *authnservicev1.UpdateNotificationPreferencesRequest) (*authnservicev1.UpdateNotificationPreferencesResponse, error) {
	user, err := s.userRepo.GetByID(ctx, req.UserId)
	if err != nil {
		appLogger.Errorf("UpdateNotificationPreferences: user not found: %v", err)
		return nil, errors.New("user not found")
	}
	_ = user
	if err := s.userRepo.UpdateNotificationPreferences(ctx, req.UserId, req.NotificationPreference, req.PreferredLanguage); err != nil {
		appLogger.Errorf("UpdateNotificationPreferences: %v", err)
		return nil, errors.New("update notification preferences")
	}
	return &authnservicev1.UpdateNotificationPreferencesResponse{
		Message: "Notification preferences updated",
	}, nil
}
