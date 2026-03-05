package service

import (
	"context"
	"encoding/base64"
	"errors"
	"strconv"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"

	authnservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
	kycservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/kyc/services/v1"
	"github.com/redis/go-redis/v9"
)

const (
	kycSessionTTL             = 30 * time.Minute
	kycSessionOwnerKey        = "kyc:session:owner:"
	kycSessionFramesKey       = "kyc:session:frames:"
	kycTotalSteps       int32 = 3
)

func kycOwnerKey(sessionID string) string {
	return kycSessionOwnerKey + sessionID
}

func kycFramesKey(sessionID string) string {
	return kycSessionFramesKey + sessionID
}

func frameDataURL(imageData []byte) string {
	return "data:application/octet-stream;base64," + base64.StdEncoding.EncodeToString(imageData)
}

func (s *AuthService) cacheKYCSessionOwner(ctx context.Context, sessionID, userID string) {
	if s == nil || s.tokenService == nil || s.tokenService.rdb == nil {
		return
	}
	_ = s.tokenService.rdb.Set(ctx, kycOwnerKey(sessionID), userID, kycSessionTTL).Err()
	_ = s.tokenService.rdb.Set(ctx, kycFramesKey(sessionID), 0, kycSessionTTL).Err()
}

func (s *AuthService) ensureKYCSessionOwner(ctx context.Context, sessionID, userID string) error {
	if s == nil || s.kycRepo == nil {
		return errors.New("kyc repository not configured")
	}

	// Fast path: Redis mapping set during InitiateKYC.
	if s.tokenService != nil && s.tokenService.rdb != nil {
		owner, err := s.tokenService.rdb.Get(ctx, kycOwnerKey(sessionID)).Result()
		if err == nil {
			if owner != userID {
				return errors.New("session does not belong to user")
			}
			return nil
		}
		if err != redis.Nil {
			logger.Errorf("kyc session lookup failed: %v", err)
			return errors.New("kyc session lookup failed")
		}
	}

	// Fallback: DB ownership check via KYC record.
	k, err := s.kycRepo.GetByID(ctx, sessionID)
	if err != nil {
		return errors.New("kyc session not found")
	}
	if k.EntityId != userID {
		return errors.New("session does not belong to user")
	}
	return nil
}

// SubmitKYCFrame validates KYC session ownership and returns orchestration progress.
// FLVE proxying is introduced in the next phase; this keeps API/state stable now.
func (s *AuthService) SubmitKYCFrame(ctx context.Context, req *authnservicev1.SubmitKYCFrameRequest) (*authnservicev1.SubmitKYCFrameResponse, error) {
	if err := s.ensureKYCSessionOwner(ctx, req.SessionId, req.UserId); err != nil {
		return nil, err
	}

	if s.externalKYC != nil {
		seq := req.FrameSequence
		if seq <= 0 {
			seq = 1
		}
		_, err := s.externalKYC.UploadDocument(ctx, &kycservicev1.UploadDocumentRequest{
			KycVerificationId: req.SessionId,
			DocumentType:      "LIVENESS_FRAME",
			DocumentNumber:    strconv.Itoa(int(seq)),
			DocumentUrl:       frameDataURL(req.ImageData),
		})
		if err != nil {
			logger.Errorf("submit KYC frame (external): %v", err)
			return nil, errors.New("submit KYC frame (external)")
		}
	}

	completed := req.FrameSequence
	if completed <= 0 && s.tokenService != nil && s.tokenService.rdb != nil {
		n, err := s.tokenService.rdb.Incr(ctx, kycFramesKey(req.SessionId)).Result()
		if err == nil {
			_ = s.tokenService.rdb.Expire(ctx, kycFramesKey(req.SessionId), kycSessionTTL).Err()
			completed = int32(n)
		}
	}
	if completed <= 0 {
		completed = 1
	}
	if completed > kycTotalSteps {
		completed = kycTotalSteps
	}

	currentStep := "LOOK_CENTER"
	switch completed {
	case 1:
		currentStep = "BLINK"
	case 2:
		currentStep = "LOOK_LEFT"
	default:
		currentStep = "LOOK_RIGHT"
	}

	confidence := float64(completed) / float64(kycTotalSteps)
	return &authnservicev1.SubmitKYCFrameResponse{
		Accepted:           true,
		Guidance:           "Hold steady and follow challenge prompt",
		CurrentStep:        currentStep,
		CompletedSteps:     completed,
		TotalSteps:         kycTotalSteps,
		LivenessConfidence: confidence,
		Message:            "Frame processed",
	}, nil
}

// CompleteKYCSession finalizes KYC for the session after ownership validation.
// In the FLVE phase this will proxy CompleteEKYC and persist FLVE outputs.
func (s *AuthService) CompleteKYCSession(ctx context.Context, req *authnservicev1.CompleteKYCSessionRequest) (*authnservicev1.CompleteKYCSessionResponse, error) {
	if err := s.ensureKYCSessionOwner(ctx, req.SessionId, req.UserId); err != nil {
		return nil, err
	}

	if s.externalKYC != nil {
		_, err := s.externalKYC.VerifyKYC(ctx, &kycservicev1.VerifyKYCRequest{
			KycVerificationId:  req.SessionId,
			VerifiedBy:         req.UserId,
			VerificationResult: "AUTO_VERIFIED",
		})
		if err != nil {
			logger.Errorf("complete KYC session (external): %v", err)
			return nil, errors.New("complete KYC session (external)")
		}
	}

	now := time.Now()
	if err := s.kycRepo.MarkVerified(ctx, req.SessionId, req.UserId, now, nil); err != nil {
		logger.Errorf("complete kyc session: %v", err)
		return nil, errors.New("complete kyc session")
	}
	if s.userProfileRepo != nil {
		_ = s.userProfileRepo.SetKYCVerified(ctx, req.UserId, true, &now)
	}

	return &authnservicev1.CompleteKYCSessionResponse{
		KycId:              req.SessionId,
		Status:             "VERIFIED",
		Success:            true,
		LivenessConfidence: 1.0,
		ProfileImageUrl:    "",
		Message:            "KYC session completed successfully",
	}, nil
}
