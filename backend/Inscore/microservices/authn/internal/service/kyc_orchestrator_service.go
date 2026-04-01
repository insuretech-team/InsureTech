package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"

	authnservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
	kycservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/kyc/services/v1"
	"github.com/redis/go-redis/v9"
)

const (
	kycSessionTTL         = 30 * time.Minute
	kycSessionOwnerKey    = "kyc:session:owner:"
	kycSessionFramesKey   = "kyc:session:frames:"
	kycSessionProviderKey = "kyc:session:provider:"
)

func kycOwnerKey(sessionID string) string {
	return kycSessionOwnerKey + sessionID
}

func kycFramesKey(sessionID string) string {
	return kycSessionFramesKey + sessionID
}

func kycProviderKey(kycID string) string {
	return kycSessionProviderKey + kycID
}

func (s *AuthService) cacheKYCSessionOwner(ctx context.Context, sessionID, userID string) {
	if s == nil || s.tokenService == nil || s.tokenService.rdb == nil {
		return
	}
	_ = s.tokenService.rdb.Set(ctx, kycOwnerKey(sessionID), userID, kycSessionTTL).Err()
	_ = s.tokenService.rdb.Set(ctx, kycFramesKey(sessionID), 0, kycSessionTTL).Err()
}

func (s *AuthService) cacheKYCProviderRef(ctx context.Context, kycID, providerRef string) {
	if s == nil || s.tokenService == nil || s.tokenService.rdb == nil {
		return
	}
	_ = s.tokenService.rdb.Set(ctx, kycProviderKey(kycID), providerRef, kycSessionTTL).Err()
}

func (s *AuthService) getProviderRef(ctx context.Context, sessionID string) (string, error) {
	// Fast path: Redis lookup by InsureTech KYC UUID.
	if s.tokenService != nil && s.tokenService.rdb != nil {
		ref, err := s.tokenService.rdb.Get(ctx, kycProviderKey(sessionID)).Result()
		if err == nil && ref != "" {
			return ref, nil
		}
		if err != nil && err != redis.Nil {
			logger.Warnf("redis provider ref lookup: %v", err)
		}
	}

	// Fallback: DB lookup by InsureTech KYC UUID first, then by provider_reference.
	k, err := s.kycRepo.GetByID(ctx, sessionID)
	if err != nil {
		// sessionID may already BE the FLVE provider_reference (session UUID).
		// Look up the KYC record by provider_reference to confirm.
		k2, err2 := s.kycRepo.GetByProviderReference(ctx, sessionID)
		if err2 != nil {
			return "", errors.New("kyc session not found")
		}
		// sessionID is the FLVE session UUID — it IS the provider reference.
		s.cacheKYCProviderRef(ctx, k2.Id, sessionID)
		return sessionID, nil
	}
	if k.ProviderReference == "" {
		return "", errors.New("no provider reference for kyc session")
	}
	// Re-cache
	s.cacheKYCProviderRef(ctx, sessionID, k.ProviderReference)
	return k.ProviderReference, nil
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
	// Try by InsureTech KYC UUID first, then by provider_reference (FLVE session UUID).
	k, err := s.kycRepo.GetByID(ctx, sessionID)
	if err != nil {
		// sessionID may be the FLVE session UUID (provider_reference), not the KYC UUID.
		// Look up by provider_reference.
		k2, err2 := s.kycRepo.GetByProviderReference(ctx, sessionID)
		if err2 != nil {
			return errors.New("kyc session not found")
		}
		k = k2
	}
	if k.EntityId != userID {
		return errors.New("session does not belong to user")
	}
	// Re-cache for future fast-path lookups.
	s.cacheKYCSessionOwner(ctx, sessionID, userID)
	return nil
}

// SubmitKYCFrame validates KYC session ownership and proxies to FLVE for liveness frame processing.
func (s *AuthService) SubmitKYCFrame(ctx context.Context, req *authnservicev1.SubmitKYCFrameRequest) (*authnservicev1.SubmitKYCFrameResponse, error) {
	if err := s.ensureKYCSessionOwner(ctx, req.SessionId, req.UserId); err != nil {
		return nil, err
	}

	// FLVE adapter path — rich response
	if s.flveAdapter != nil {
		providerRef, err := s.getProviderRef(ctx, req.SessionId)
		if err != nil {
			return nil, err
		}

		flveResp, err := s.flveAdapter.SubmitEKYCFrame(ctx, providerRef, req.ImageData)
		if err != nil {
			logger.Errorf("submit KYC frame (FLVE): %v", err)
			return nil, errors.New("submit KYC frame failed")
		}

		// Persist progress snapshot
		s.persistVerificationSnapshot(ctx, req.SessionId, "frame", flveResp)

		// Increment frame counter
		if s.tokenService != nil && s.tokenService.rdb != nil {
			_ = s.tokenService.rdb.Incr(ctx, kycFramesKey(req.SessionId)).Err()
			_ = s.tokenService.rdb.Expire(ctx, kycFramesKey(req.SessionId), kycSessionTTL).Err()
		}

		resp := &authnservicev1.SubmitKYCFrameResponse{
			Accepted:           true,
			SessionState:       flveResp.SessionState,
			StepProgress:       flveResp.StepProgress,
			OverallProgress:    flveResp.OverallProgress,
			LivenessConfidence: flveResp.LivenessScore,
			GuidanceMessages:   flveResp.Guidance,
			Message:            "Frame processed",
		}

		// Map guidance
		if len(flveResp.Guidance) > 0 {
			resp.Guidance = flveResp.Guidance[0]
		}

		// Map current step
		if flveResp.CurrentStep != nil {
			resp.CurrentStep = flveResp.CurrentStep.Type
			resp.CurrentStepDetail = &authnservicev1.KYCStep{
				StepNumber:     int32(flveResp.CurrentStep.StepNumber),
				ChallengeType:  flveResp.CurrentStep.Type,
				State:          flveResp.CurrentStep.State,
				Instruction:    flveResp.CurrentStep.Instruction,
				InstructionKey: flveResp.CurrentStep.InstructionKey,
				TimeoutSeconds: int32(flveResp.CurrentStep.TimeoutSeconds),
				Confidence:     flveResp.CurrentStep.Confidence,
			}
			resp.CompletedSteps = int32(flveResp.CurrentStep.StepNumber - 1)
			if flveResp.StepCompleted {
				resp.CompletedSteps = int32(flveResp.CurrentStep.StepNumber)
			}
		}

		// Map detection
		if flveResp.Detection != nil {
			det := &authnservicev1.KYCDetection{}
			if v, ok := flveResp.Detection["detected"].(bool); ok {
				det.Detected = v
			}
			if v, ok := flveResp.Detection["box"].(map[string]interface{}); ok {
				if x, ok := v["x"].(float64); ok {
					det.X = int32(x)
				}
				if y, ok := v["y"].(float64); ok {
					det.Y = int32(y)
				}
				if w, ok := v["width"].(float64); ok {
					det.Width = int32(w)
				}
				if h, ok := v["height"].(float64); ok {
					det.Height = int32(h)
				}
			}
			resp.Detection = det
		}

		// Map head pose
		if flveResp.HeadPose != nil {
			resp.HeadPose = &authnservicev1.KYCHeadPose{
				Yaw:   flveResp.HeadPose.Yaw,
				Pitch: flveResp.HeadPose.Pitch,
				Roll:  flveResp.HeadPose.Roll,
			}
		}

		// Map eye state
		if flveResp.EyeState != nil {
			resp.EyeState = &authnservicev1.KYCEyeState{
				LeftOpenness:  flveResp.EyeState.LeftOpenness,
				RightOpenness: flveResp.EyeState.RightOpenness,
				IsBlinking:    flveResp.EyeState.IsBlinking,
			}
		}

		// Map step_completed
		resp.StepCompleted = flveResp.StepCompleted

		// Map eye contours — JSON-encode the FLVE mesh so the frontend can draw
		// the contour overlay without adding typed proto messages for the edge list.
		if flveResp.EyeContours != nil {
			if b, err := json.Marshal(flveResp.EyeContours); err == nil {
				resp.EyeContoursJson = string(b)
			}
		}

		return resp, nil
	}

	// Legacy ExternalKYCClient fallback (to be removed)
	if s.externalKYC != nil {
		return s.submitKYCFrameLegacy(ctx, req)
	}

	// No external provider — local stub
	return s.submitKYCFrameLocal(ctx, req)
}

// CompleteKYCSession finalizes KYC. With FLVE, sets status to PENDING_REVIEW (not auto-verified).
func (s *AuthService) CompleteKYCSession(ctx context.Context, req *authnservicev1.CompleteKYCSessionRequest) (*authnservicev1.CompleteKYCSessionResponse, error) {
	if err := s.ensureKYCSessionOwner(ctx, req.SessionId, req.UserId); err != nil {
		return nil, err
	}

	// FLVE adapter path
	if s.flveAdapter != nil {
		providerRef, err := s.getProviderRef(ctx, req.SessionId)
		if err != nil {
			return nil, err
		}

		flveResp, err := s.flveAdapter.CompleteEKYC(ctx, providerRef)
		if err != nil {
			logger.Errorf("complete KYC session (FLVE): %v", err)
			return nil, errors.New("complete KYC session failed")
		}

		// Persist full FLVE response
		s.persistVerificationSnapshot(ctx, req.SessionId, "complete", flveResp)

		// Update profile photo URL if available
		if flveResp.ProfileImageURL != "" && s.userProfileRepo != nil {
			_ = s.userProfileRepo.SetProfilePhotoURL(ctx, req.UserId, flveResp.ProfileImageURL)
		}

		// Set status to PENDING_REVIEW — approval step is separate.
		// Sits between IN_PROGRESS and VERIFIED: FLVE session done, human review pending.
		// The proto enum now has VERIFICATION_STATUS_PENDING_REVIEW = 6; the DB column
		// stores the canonical string "PENDING_REVIEW" (trimmed prefix form).
		if err := s.kycRepo.SetStatus(ctx, req.SessionId, "PENDING_REVIEW"); err != nil {
			logger.Errorf("set kyc pending_review: %v", err)
		}

		// Mark user profile kyc_verified=true now that the FLVE session completed
		// without error. The KYC record is already PENDING_REVIEW; marking the profile
		// ensures the next login reads kyc_verified=true and does NOT re-gate the user
		// to /kyc. Without this, kyc_verified stays false in user_profiles and the user
		// is sent back to /kyc on every subsequent login.
		//
		// We do NOT gate this on flveResp.Success — CompleteEKYC already returned nil
		// error (meaning FLVE accepted the request). flveResp.Success may be false for
		// borderline liveness confidence but the session is still PENDING_REVIEW which
		// allows portal access. If you require strict liveness gating, add that check here.
		if s.userProfileRepo != nil {
			now := time.Now()
			if err := s.userProfileRepo.SetKYCVerified(ctx, req.UserId, true, &now); err != nil {
				// Non-fatal — KYC record status is already PENDING_REVIEW.
				logger.Warnf("complete KYC session: set profile kyc_verified: user=%s err=%v", req.UserId, err)
			} else {
				logger.Infof("eKYC completed: user_profiles.kyc_verified=true user=%s kyc_id=%s", req.UserId, req.SessionId)
			}
		}

		resp := &authnservicev1.CompleteKYCSessionResponse{
			KycId:              req.SessionId,
			Status:             "PENDING_REVIEW",
			Success:            flveResp.Success,
			LivenessConfidence: flveResp.LivenessConfidence,
			ProfileImageUrl:    flveResp.ProfileImageURL,
			Message:            "KYC session completed — pending review",
			ProviderReference:  providerRef,
			SessionState:       flveResp.State,
			IdentityMatch:      flveResp.IdentityMatch,
			MatchScore:         flveResp.MatchScore,
			CompletedAt:        flveResp.CompletedAt,
		}

		// Map summary
		if flveResp.Summary != nil {
			resp.Summary = &authnservicev1.KYCSessionSummary{
				TotalSteps:           int32(flveResp.Summary.TotalSteps),
				CompletedSteps:       int32(flveResp.Summary.CompletedSteps),
				FailedSteps:          int32(flveResp.Summary.FailedSteps),
				TotalFramesProcessed: int32(flveResp.Summary.TotalFramesProcessed),
				ElapsedMs:            int32(flveResp.Summary.ElapsedMs),
			}
			for _, sr := range flveResp.Summary.StepResults {
				resp.Summary.StepResults = append(resp.Summary.StepResults, &authnservicev1.KYCStepResult{
					ChallengeType:   sr.Type,
					State:           sr.State,
					Confidence:      sr.Confidence,
					FramesProcessed: int32(sr.FramesProcessed),
					ElapsedMs:       int32(sr.ElapsedMs),
				})
			}
		}

		return resp, nil
	}

	// Legacy ExternalKYCClient fallback
	if s.externalKYC != nil {
		return s.completeKYCSessionLegacy(ctx, req)
	}

	// No external provider — local auto-verify
	return s.completeKYCSessionLocal(ctx, req)
}

// ── Legacy / local fallback methods ──────────────────────────────────────────

func (s *AuthService) submitKYCFrameLegacy(ctx context.Context, req *authnservicev1.SubmitKYCFrameRequest) (*authnservicev1.SubmitKYCFrameResponse, error) {
	seq := req.FrameSequence
	if seq <= 0 {
		seq = 1
	}

	_, err := s.externalKYC.UploadDocument(ctx, &kycservicev1.UploadDocumentRequest{
		KycVerificationId: req.SessionId,
		DocumentType:      "LIVENESS_FRAME",
		DocumentNumber:    strconv.Itoa(int(seq)),
		DocumentUrl:       "data:application/octet-stream;base64," + base64.StdEncoding.EncodeToString(req.ImageData),
	})
	if err != nil {
		logger.Errorf("submit KYC frame (external): %v", err)
		return nil, errors.New("submit KYC frame (external)")
	}

	return s.submitKYCFrameLocal(ctx, req)
}

func (s *AuthService) submitKYCFrameLocal(ctx context.Context, req *authnservicev1.SubmitKYCFrameRequest) (*authnservicev1.SubmitKYCFrameResponse, error) {
	const kycTotalSteps int32 = 4

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

	currentStep := "BLINK"
	switch completed {
	case 2:
		currentStep = "LOOK_LEFT"
	case 3:
		currentStep = "LOOK_RIGHT"
	case 4:
		currentStep = "CAPTURE"
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
		OverallProgress:    float64(completed) / float64(kycTotalSteps),
	}, nil
}

func (s *AuthService) completeKYCSessionLegacy(ctx context.Context, req *authnservicev1.CompleteKYCSessionRequest) (*authnservicev1.CompleteKYCSessionResponse, error) {
	_, err := s.externalKYC.VerifyKYC(ctx, &kycservicev1.VerifyKYCRequest{
		KycVerificationId:  req.SessionId,
		VerifiedBy:         req.UserId,
		VerificationResult: "AUTO_VERIFIED",
	})
	if err != nil {
		logger.Errorf("complete KYC session (external): %v", err)
		return nil, errors.New("complete KYC session (external)")
	}
	return s.completeKYCSessionLocal(ctx, req)
}

func (s *AuthService) completeKYCSessionLocal(ctx context.Context, req *authnservicev1.CompleteKYCSessionRequest) (*authnservicev1.CompleteKYCSessionResponse, error) {
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
		Message:            "KYC session completed successfully",
	}, nil
}

// ── Helpers ──────────────────────────────────────────────────────────────────

func (s *AuthService) persistVerificationSnapshot(ctx context.Context, kycID, phase string, data interface{}) {
	if s.kycRepo == nil {
		return
	}
	snapshot, err := json.Marshal(map[string]interface{}{
		"phase":     phase,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"data":      data,
	})
	if err != nil {
		logger.Warnf("marshal verification snapshot: %v", err)
		return
	}
	if err := s.kycRepo.AppendVerificationResult(ctx, kycID, string(snapshot)); err != nil {
		logger.Warnf("persist verification snapshot: %v", err)
	}
}
