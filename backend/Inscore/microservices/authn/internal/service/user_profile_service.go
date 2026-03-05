package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	authnentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/entity/v1"
	authnservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
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
func (s *AuthService) GetUserProfile(ctx context.Context, req *authnservicev1.GetUserProfileRequest) (*authnservicev1.GetUserProfileResponse, error) {
	if s.userProfileRepo == nil {
		return nil, errors.New("user profile repository not configured")
	}

	profile, err := s.userProfileRepo.GetByUserID(ctx, req.UserId)
	if err != nil {
		logger.Errorf("profile not found: %v", err)
		return nil, errors.New("profile not found")
	}

	return &authnservicev1.GetUserProfileResponse{Profile: profile}, nil
}

// UpdateUserProfile updates existing profile fields (non-zero values only).
func (s *AuthService) UpdateUserProfile(ctx context.Context, req *authnservicev1.UpdateUserProfileRequest) (*authnservicev1.UpdateUserProfileResponse, error) {
	if s.userProfileRepo == nil {
		return nil, errors.New("user profile repository not configured")
	}

	existing, err := s.userProfileRepo.GetByUserID(ctx, req.UserId)
	if err != nil || existing == nil {
		return nil, errors.New("profile not found")
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
