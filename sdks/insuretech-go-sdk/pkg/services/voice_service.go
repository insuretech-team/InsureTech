package services

import (
	"context"
	"strings"
	"github.com/newage-saint/insuretech-go-sdk/pkg/models"
)

// VoiceService handles voice-related API calls
type VoiceService struct {
	Client Client
}

// StartVoiceSession Start voice session
func (s *VoiceService) StartVoiceSession(ctx context.Context, req *models.VoiceSessionStartRequest) error {
	path := "/v1/voice-sessions"
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// GetVoiceSession Get voice session
func (s *VoiceService) GetVoiceSession(ctx context.Context, voiceSessionId string) error {
	path := "/v1/voice-sessions/{voice_session_id}"
	path = strings.ReplaceAll(path, "{voice_session_id}", voiceSessionId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// ProcessVoiceCommand Process voice command
func (s *VoiceService) ProcessVoiceCommand(ctx context.Context, voiceSessionId string, req *models.VoiceCommandProcessingRequest) error {
	path := "/v1/voice-sessions/{voice_session_id}/commands"
	path = strings.ReplaceAll(path, "{voice_session_id}", voiceSessionId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

// GetTranscript Get transcript
func (s *VoiceService) GetTranscript(ctx context.Context, voiceSessionId string) error {
	path := "/v1/voice-sessions/{voice_session_id}/transcript"
	path = strings.ReplaceAll(path, "{voice_session_id}", voiceSessionId)
	return s.Client.DoRequest(ctx, "GET", path, nil, nil)
}

// EndVoiceSession End voice session
func (s *VoiceService) EndVoiceSession(ctx context.Context, voiceSessionId string, req *models.EndVoiceSessionRequest) error {
	path := "/v1/voice-sessions/{voice_session_id}:end"
	path = strings.ReplaceAll(path, "{voice_session_id}", voiceSessionId)
	return s.Client.DoRequest(ctx, "POST", path, req, nil)
}

