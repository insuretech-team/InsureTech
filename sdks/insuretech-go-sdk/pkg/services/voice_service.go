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
func (s *VoiceService) StartVoiceSession(ctx context.Context, req *models.VoiceSessionStartRequest) (*models.VoiceSessionStartResponse, error) {
	path := "/v1/voice-sessions"
	var result models.VoiceSessionStartResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ProcessVoiceCommand Process voice command
func (s *VoiceService) ProcessVoiceCommand(ctx context.Context, voiceSessionId string, req *models.VoiceCommandProcessingRequest) (*models.VoiceCommandProcessingResponse, error) {
	path := "/v1/voice-sessions/{voice_session_id}/commands"
	path = strings.ReplaceAll(path, "{voice_session_id}", voiceSessionId)
	var result models.VoiceCommandProcessingResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetVoiceSession Get voice session
func (s *VoiceService) GetVoiceSession(ctx context.Context, voiceSessionId string) (*models.VoiceVoiceSessionRetrievalResponse, error) {
	path := "/v1/voice-sessions/{voice_session_id}"
	path = strings.ReplaceAll(path, "{voice_session_id}", voiceSessionId)
	var result models.VoiceVoiceSessionRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// EndVoiceSession End voice session
func (s *VoiceService) EndVoiceSession(ctx context.Context, voiceSessionId string, req *models.VoiceEndVoiceSessionRequest) (*models.VoiceEndVoiceSessionResponse, error) {
	path := "/v1/voice-sessions/{voice_session_id}"
	path = strings.ReplaceAll(path, "{voice_session_id}", voiceSessionId)
	var result models.VoiceEndVoiceSessionResponse
	err := s.Client.DoRequest(ctx, "POST", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetTranscript Get transcript
func (s *VoiceService) GetTranscript(ctx context.Context, voiceSessionId string) (*models.TranscriptRetrievalResponse, error) {
	path := "/v1/voice-sessions/{voice_session_id}/transcript"
	path = strings.ReplaceAll(path, "{voice_session_id}", voiceSessionId)
	var result models.TranscriptRetrievalResponse
	err := s.Client.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

