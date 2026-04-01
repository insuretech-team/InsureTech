package service

// kyc_service.go - KYC (Know Your Customer) verification methods.
// InitiateKYC, GetKYCStatus, ApproveKYC, rejectKYCInternal, VerifyDocument.
// TOTP -> totp_service.go | Voice -> voice_service.go | Profile -> user_profile_service.go

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	authnservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
	kycv1 "github.com/newage-saint/insuretech/gen/go/insuretech/kyc/entity/v1"
	kycservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/kyc/services/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// ── KYC ──────────────────────────────────────────────────────────────────────

// InitiateKYC creates a new KYC verification record and starts an FLVE eKYC session.
// Delegates to the KYC microservice (synchronous gRPC call) which calls FLVE HF Space
// and returns the session_id. Publishes KYCVerificationStarted Kafka event (best-effort).
func (s *AuthService) InitiateKYC(ctx context.Context, req *authnservicev1.InitiateKYCRequest) (*authnservicev1.InitiateKYCResponse, error) {
	if s.kycRepo == nil {
		return nil, errors.New("kyc repository not configured")
	}

	// Delegate to KYC microservice — it calls FLVE synchronously and returns session_id.
	if s.externalKYC != nil {
		extResp, err := s.externalKYC.StartKYCVerification(ctx, &kycservicev1.StartKYCVerificationRequest{
			Type:       "KYC",
			EntityType: "user",
			EntityId:   req.UserId,
			Method:     "FLVE_EKYC",
		})
		if err != nil {
			logger.Errorf("InitiateKYC: KYC service failed: %v", err)
			return nil, status.Errorf(codes.Internal, "KYC initiation failed: %v", err)
		}

		kycID := extResp.GetKycVerificationId()
		// KYC service encodes the FLVE session_id in the Message field.
		sessionID := extResp.GetMessage()

		// Cache the session ownership and provider reference so SubmitKYCFrame
		// and CompleteKYCSession can validate the session without hitting DB.
		// sessionID = FLVE session UUID, kycID = InsureTech KYC record UUID.
		if sessionID != "" {
			s.cacheKYCSessionOwner(ctx, sessionID, req.UserId)
			s.cacheKYCProviderRef(ctx, kycID, sessionID)
		}

		// Publish KYCVerificationStarted event (non-blocking, best-effort)
		if s.eventPublisher != nil {
			_ = s.eventPublisher.PublishKYCVerificationStarted(ctx, kycID, "user", req.UserId, "FLVE_EKYC")
		}

		return &authnservicev1.InitiateKYCResponse{
			KycId:             kycID,
			Status:            "IN_PROGRESS",
			Message:           "FLVE eKYC session started. Follow the challenge steps.",
			Provider:          "FLVE",
			ProviderReference: sessionID,
			SessionState:      "EKYC_SESSION_ACTIVE",
		}, nil
	}

	// Fallback: no KYC service configured — create record locally with no FLVE session.
	kycID := uuid.New().String()
	k := &kycv1.KYCVerification{
		Id:         kycID,
		Type:       kycv1.VerificationType_VERIFICATION_TYPE_KYC,
		EntityType: "user",
		EntityId:   req.UserId,
		Method:     kycv1.VerificationMethod_VERIFICATION_METHOD_FLVE_EKYC,
		Status:     kycv1.VerificationStatus_VERIFICATION_STATUS_IN_PROGRESS,
	}
	if err := s.kycRepo.Create(ctx, k); err != nil {
		logger.Errorf("initiate KYC (local): %v", err)
		return nil, errors.New("initiate KYC failed")
	}
	return &authnservicev1.InitiateKYCResponse{
		KycId:    kycID,
		Status:   "IN_PROGRESS",
		Message:  "KYC verification initiated.",
		Provider: "FLVE",
	}, nil
}

// mdFirst returns the first value for a gRPC metadata key, or empty string.
func mdFirst(md metadata.MD, key string) string {
	if vals := md.Get(key); len(vals) > 0 {
		return vals[0]
	}
	return ""
}

// GetKYCStatus returns the current KYC status for a user.
// When FLVE adapter is configured and status is IN_PROGRESS, fetches live data from FLVE.
func (s *AuthService) GetKYCStatus(ctx context.Context, req *authnservicev1.GetKYCStatusRequest) (*authnservicev1.GetKYCStatusResponse, error) {
	if s.kycRepo == nil {
		return nil, errors.New("kyc repository not configured")
	}
	k, err := s.kycRepo.GetByEntity(ctx, "user", req.UserId)
	if err != nil {
		// BUG-007 FIX: User has no KYC record yet — return NOT_STARTED instead of 500.
		// gorm.ErrRecordNotFound (or "sql: no rows") means KYC not initiated yet, which
		// is a normal state for new B2C users. Previously this crashed with 500.
		errStr := err.Error()
		if errStr == "record not found" || errStr == "sql: no rows in result set" || strings.Contains(errStr, "no rows") {
			return &authnservicev1.GetKYCStatusResponse{
				Status: "NOT_STARTED",
			}, nil
		}
		logger.Errorf("get KYC status: %v", err)
		return nil, errors.New("get KYC status")
	}
	// Map proto enum to client-facing status string.
	// PENDING_REVIEW surfaces as "PENDING_REVIEW" so clients can show "Under Review" UI.
	resp := &authnservicev1.GetKYCStatusResponse{
		KycId:  k.Id,
		Status: strings.TrimPrefix(k.Status.String(), "VERIFICATION_STATUS_"),
	}
	if k.RejectionReason != "" {
		resp.RejectionReason = k.RejectionReason
	}
	if k.VerifiedAt != nil {
		resp.ReviewedAt = k.VerifiedAt
	}

	// Enrich with live FLVE data when session is IN_PROGRESS.
	// PENDING_REVIEW means FLVE is done — no need to poll live status anymore.
	isActiveSession := k.Status == kycv1.VerificationStatus_VERIFICATION_STATUS_IN_PROGRESS
	if s.flveAdapter != nil && k.ProviderReference != "" && isActiveSession {
		flveStatus, err := s.flveAdapter.GetEKYCStatus(ctx, k.ProviderReference)
		if err == nil {
			resp.Provider = "FLVE"
			resp.ProviderReference = k.ProviderReference
			resp.SessionState = flveStatus.State
			resp.OverallProgress = flveStatus.OverallProgress
			resp.RemainingSeconds = int32(flveStatus.RemainingSeconds)

			for _, st := range flveStatus.Steps {
				// Match FLVE step states — also handle proto-prefixed forms
				state := strings.TrimPrefix(st.State, "EKYC_STEP_")
				if state == "PENDING" || state == "IN_PROGRESS" {
					resp.CurrentStep = &authnservicev1.KYCStep{
						StepNumber:     int32(st.StepNumber),
						ChallengeType:  strings.TrimPrefix(st.Type, "EKYC_CHALLENGE_"),
						State:          state,
						Instruction:    st.Instruction,
						InstructionKey: st.InstructionKey,
						TimeoutSeconds: int32(st.TimeoutSeconds),
						Confidence:     st.Confidence,
					}
					break
				}
			}

			completed := 0
			for _, st := range flveStatus.Steps {
				if strings.TrimPrefix(st.State, "EKYC_STEP_") == "COMPLETED" {
					completed++
				}
			}
			resp.CompletedSteps = int32(completed)
			resp.TotalSteps = int32(len(flveStatus.Steps))
		} else {
			logger.Warnf("get FLVE live status: %v", err)
		}
	} else if k.ProviderReference != "" {
		// Session completed/rejected — return stored provider reference for audit
		resp.Provider = "FLVE"
		resp.ProviderReference = k.ProviderReference
	}

	return resp, nil
}

// ApproveKYC sets KYC status to VERIFIED and marks user profile kyc_verified=true.
//
// Full approval flow:
//  1. MarkVerified — sets status=VERIFIED, verified_by, verified_at in kyc_verifications
//  2. SetKYCVerified — sets kyc_verified=true, kyc_verified_at on user_profiles
//  3. Publish KYCVerifiedEvent to Kafka so downstream services react (policies, notifications)
func (s *AuthService) ApproveKYC(ctx context.Context, req *authnservicev1.ApproveKYCRequest) (*authnservicev1.ApproveKYCResponse, error) {
	if s.kycRepo == nil {
		return nil, errors.New("kyc repository not configured")
	}
	if req.GetKycId() == "" {
		return nil, status.Error(codes.InvalidArgument, "kyc_id is required")
	}
	if req.GetReviewerId() == "" {
		return nil, status.Error(codes.InvalidArgument, "reviewer_id is required")
	}

	now := time.Now()

	// 1. Mark KYC record as VERIFIED
	if err := s.kycRepo.MarkVerified(ctx, req.KycId, req.ReviewerId, now, nil); err != nil {
		logger.Errorf("approve KYC mark verified: %v", err)
		return nil, status.Error(codes.Internal, "approve KYC failed")
	}

	// 2. Fetch KYC record to get entity_id (user_id) and method
	k, fetchErr := s.kycRepo.GetByID(ctx, req.KycId)

	// 3. Mark user profile as KYC verified
	if fetchErr == nil && k != nil && s.userProfileRepo != nil {
		if err := s.userProfileRepo.SetKYCVerified(ctx, k.EntityId, true, &now); err != nil {
			// Non-fatal — KYC record is already VERIFIED; profile update is best-effort
			logger.Warnf("approve KYC set profile kyc_verified: entity=%s err=%v", k.EntityId, err)
		} else {
			logger.Infof("KYC approved: kyc_id=%s entity_id=%s reviewer=%s", req.KycId, k.EntityId, req.ReviewerId)
		}

		// 4. Publish KYCVerifiedEvent so downstream services react (policies, notifications).
		if s.eventPublisher != nil {
			go func() {
				pubErr := s.eventPublisher.PublishKYCVerified(context.Background(), k.EntityId, req.KycId, req.ReviewerId, now)
				if pubErr != nil {
					logger.Warnf("publish KYC verified event: %v", pubErr)
				}
			}()
		}
	} else if fetchErr != nil {
		logger.Warnf("approve KYC fetch record: %v", fetchErr)
	}

	return &authnservicev1.ApproveKYCResponse{
		Message: "KYC approved successfully",
	}, nil
}

// RejectKYC sets KYC status to REJECTED with a rejection reason.
// NOTE: RejectKYC RPC was removed from auth_service.proto (API path conflict).
// This method is kept as internal logic callable via ApproveKYC flow or admin tooling.
func (s *AuthService) rejectKYCInternal(ctx context.Context, kycID string, rejectionReason string) error {
	if s.kycRepo == nil {
		return errors.New("kyc repository not configured")
	}
	reason := rejectionReason
	if err := s.kycRepo.UpdateStatus(ctx, kycID, kycv1.VerificationStatus_VERIFICATION_STATUS_REJECTED, &reason); err != nil {
		logger.Errorf("reject KYC: %v", err)
		return errors.New("reject KYC")
	}
	return nil
}

// VerifyDocument marks a user document as verified by the given reviewer.
func (s *AuthService) VerifyDocument(ctx context.Context, req *authnservicev1.VerifyDocumentRequest) (*authnservicev1.VerifyDocumentResponse, error) {
	if s.userDocumentRepo == nil {
		return nil, errors.New("user document repository not configured")
	}
	if err := s.userDocumentRepo.MarkVerified(ctx, req.UserDocumentId, req.VerifiedBy, req.VerificationStatus, req.RejectionReason); err != nil {
		logger.Errorf("verify document: %v", err)
		return nil, errors.New("verify document")
	}
	doc, err := s.userDocumentRepo.GetByID(ctx, req.UserDocumentId)
	if err != nil {
		return &authnservicev1.VerifyDocumentResponse{Message: "Document verified"}, nil
	}
	return &authnservicev1.VerifyDocumentResponse{
		Document: doc,
		Message:  "Document verified successfully",
	}, nil
}

