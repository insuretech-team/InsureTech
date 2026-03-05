package service

// voice_service.go — Voice Biometric Session (Sprint 1.10)
// Implements a challenge-response voice authentication flow:
//   1. InitiateVoiceSession  — generates a random challenge phrase, creates a PENDING session
//   2. SubmitVoiceSample     — validates transcript + confidence score, marks COMPLETED or FAILED
//   3. VerifyVoiceSession    — confirms the session is COMPLETED and returns the user ID

import (
	"context"
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	appLogger "github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	authnservicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/services/v1"
	voicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/voice/entity/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// voiceChallengePhrases is the fixed bank of challenge phrases a user must speak.
var voiceChallengePhrases = [10]string{
	"My voice is my password",
	"Verify me with my unique voice",
	"I am logging in with my voice today",
	"InsureTech authenticates me by voice",
	"My voice unlocks my account securely",
	"Speak clearly to confirm your identity",
	"Voice biometrics keeps my account safe",
	"I authorize this login with my voice",
	"My spoken words verify who I am",
	"Trust my voice to open the door",
}

// voiceConfidenceThreshold is the minimum confidence score to accept a voice sample.
const voiceConfidenceThreshold = 0.85

// voiceSessionExpiry is the time window within which a session must be completed.
const voiceSessionExpiry = 15 * time.Minute

// ── InitiateVoiceSession ──────────────────────────────────────────────────────

// InitiateVoiceSession creates a new voice biometric challenge session.
// Returns a sessionID and a challenge phrase the user must speak.
func (s *AuthService) InitiateVoiceSession(ctx context.Context, req *authnservicev1.InitiateVoiceSessionRequest) (*authnservicev1.InitiateVoiceSessionResponse, error) {
	if s.voiceRepo == nil {
		return nil, errors.New("voice session repository not configured")
	}
	if req.UserId == "" {
		return nil, errors.New("user_id is required")
	}

	// Verify user exists.
	if _, err := s.userRepo.GetByID(ctx, req.UserId); err != nil {
		logger.Errorf("user not found: %v", err)
		return nil, errors.New("user not found")
	}

	// Pick a random challenge phrase.
	//nolint:gosec // math/rand is intentional — challenge phrases are not security-sensitive secrets
	challenge := voiceChallengePhrases[rand.Intn(len(voiceChallengePhrases))]

	// Generate a stable internal session UUID.
	sessionID := uuid.New().String()

	// Use ExternalSessionID (SessionId field) as the challenge token / nonce.
	// We repurpose the Intent field to store the challenge text.
	vs := &voicev1.VoiceSession{
		Id:        sessionID,
		SessionId: uuid.New().String(), // external_session_id: opaque nonce
		UserId:    req.UserId,
		Language:  "en",
		Status:    voicev1.SessionStatus_SESSION_STATUS_ACTIVE,
		Intent:    challenge,
		StartedAt: timestamppb.Now(),
	}

	if err := s.voiceRepo.Create(ctx, vs); err != nil {
		appLogger.Errorf("InitiateVoiceSession: failed to create session for user %s: %v", req.UserId, err)
		logger.Errorf("failed to create voice session: %v", err)
		return nil, errors.New("failed to create voice session")
	}

	appLogger.Infof("InitiateVoiceSession: session %s created for user %s", sessionID, req.UserId)

	return &authnservicev1.InitiateVoiceSessionResponse{
		SessionId: sessionID,
		Challenge: challenge,
		Message:   "Voice session initiated. Please speak the challenge phrase.",
	}, nil
}

// ── SubmitVoiceSample ─────────────────────────────────────────────────────────

// SubmitVoiceSample validates the client-submitted transcript and confidence score.
// Confidence must be >= 0.85 and the transcript must match the stored challenge
// (case-insensitive, trimmed whitespace). On success the session is marked COMPLETED;
// on failure it is marked FAILED.
func (s *AuthService) SubmitVoiceSample(ctx context.Context, req *authnservicev1.SubmitVoiceSampleRequest) (*authnservicev1.SubmitVoiceSampleResponse, error) {
	if s.voiceRepo == nil {
		return nil, errors.New("voice session repository not configured")
	}
	if req.SessionId == "" {
		return nil, errors.New("session_id is required")
	}

	vs, err := s.voiceRepo.GetByID(ctx, req.SessionId)
	if err != nil {
		logger.Errorf("voice session not found: %v", err)
		return nil, errors.New("voice session not found")
	}

	// Validate session is still active (PENDING / ACTIVE).
	if vs.Status != voicev1.SessionStatus_SESSION_STATUS_ACTIVE {
		return &authnservicev1.SubmitVoiceSampleResponse{
			Verified: false,
			Message:  "Voice session is no longer active",
		}, nil
	}

	// Check expiry: session must have been started within the last 15 minutes.
	if vs.StartedAt != nil {
		started := vs.StartedAt.AsTime()
		if time.Since(started) > voiceSessionExpiry {
			// Expire the session.
			endTime := time.Now()
			_ = s.voiceRepo.Complete(ctx, vs.Id, voicev1.SessionStatus_SESSION_STATUS_FAILED, endTime, nil)
			return &authnservicev1.SubmitVoiceSampleResponse{
				Verified: false,
				Message:  "Voice session has expired",
			}, nil
		}
	}

	// Validate confidence score and transcript match.
	transcriptMatch := strings.EqualFold(
		strings.TrimSpace(req.Transcript),
		strings.TrimSpace(vs.Intent),
	)
	confidenceOK := req.ConfidenceScore >= voiceConfidenceThreshold

	endTime := time.Now()
	if confidenceOK && transcriptMatch {
		// Mark session as COMPLETED.
		if err := s.voiceRepo.Complete(ctx, vs.Id, voicev1.SessionStatus_SESSION_STATUS_COMPLETED, endTime, nil); err != nil {
			appLogger.Errorf("SubmitVoiceSample: failed to complete session %s: %v", vs.Id, err)
			logger.Errorf("failed to update voice session: %v", err)
			return nil, errors.New("failed to update voice session")
		}
		appLogger.Infof("SubmitVoiceSample: session %s verified for user %s (confidence=%.2f)", vs.Id, vs.UserId, req.ConfidenceScore)
		return &authnservicev1.SubmitVoiceSampleResponse{
			Verified: true,
			Message:  "Voice sample verified successfully",
		}, nil
	}

	// Mark session as FAILED.
	if err := s.voiceRepo.Complete(ctx, vs.Id, voicev1.SessionStatus_SESSION_STATUS_FAILED, endTime, nil); err != nil {
		appLogger.Errorf("SubmitVoiceSample: failed to fail session %s: %v", vs.Id, err)
	}

	reason := "voice sample did not match"
	if !confidenceOK {
		reason = "confidence score " + strconv.FormatFloat(float64(req.ConfidenceScore), 'f', 2, 64) + " is below threshold " + strconv.FormatFloat(voiceConfidenceThreshold, 'f', 2, 64)
	} else if !transcriptMatch {
		reason = "transcript does not match challenge phrase"
	}

	appLogger.Warnf("SubmitVoiceSample: session %s failed for user %s: %s", vs.Id, vs.UserId, reason)
	return &authnservicev1.SubmitVoiceSampleResponse{
		Verified: false,
		Message:  "Voice verification failed: " + reason,
	}, nil
}

// ── VerifyVoiceSession ────────────────────────────────────────────────────────

// VerifyVoiceSession checks whether a voice session completed successfully.
// Returns authenticated=true and the associated user_id on success.
func (s *AuthService) VerifyVoiceSession(ctx context.Context, req *authnservicev1.VerifyVoiceSessionRequest) (*authnservicev1.VerifyVoiceSessionResponse, error) {
	if s.voiceRepo == nil {
		return nil, errors.New("voice session repository not configured")
	}
	if req.SessionId == "" {
		return nil, errors.New("session_id is required")
	}

	vs, err := s.voiceRepo.GetByID(ctx, req.SessionId)
	if err != nil {
		logger.Errorf("voice session not found: %v", err)
		return nil, errors.New("voice session not found")
	}

	// Check that the session completed within the allowed window.
	if vs.StartedAt != nil {
		started := vs.StartedAt.AsTime()
		if time.Since(started) > voiceSessionExpiry {
			return &authnservicev1.VerifyVoiceSessionResponse{
				Authenticated: false,
				Message:       "Voice session has expired",
			}, nil
		}
	}

	if vs.Status != voicev1.SessionStatus_SESSION_STATUS_COMPLETED {
		return &authnservicev1.VerifyVoiceSessionResponse{
			Authenticated: false,
			Message:       "Voice session has not been successfully verified",
		}, nil
	}

	appLogger.Infof("VerifyVoiceSession: session %s authenticated user %s", vs.Id, vs.UserId)
	return &authnservicev1.VerifyVoiceSessionResponse{
		Authenticated: true,
		UserId:        vs.UserId,
		Message:       "Voice authentication successful",
	}, nil
}
