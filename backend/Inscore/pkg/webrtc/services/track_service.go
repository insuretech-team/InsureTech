package services

import (
	"context"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/webrtc/repository"
	"github.com/newage-saint/insuretech/gen/go/insuretech/webrtc/v1/entity"
	"github.com/newage-saint/insuretech/gen/go/insuretech/webrtc/v1/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TrackService implements the WebRTC TrackService gRPC API
type TrackService struct {
	service.UnimplementedTrackServiceServer
	repo *repository.Repository
}

// NewTrackService creates a new TrackService
func NewTrackService(repo *repository.Repository) *TrackService {
	return &TrackService{repo: repo}
}

// PublishTrack publishes a media track
func (s *TrackService) PublishTrack(ctx context.Context, req *service.PublishTrackRequest) (*service.PublishTrackResponse, error) {
	if req.RoomId == "" {
		return nil, status.Error(codes.InvalidArgument, "room_id is required")
	}

	if req.PeerId == "" {
		return nil, status.Error(codes.InvalidArgument, "peer_id is required")
	}

	if req.Track == nil {
		return nil, status.Error(codes.InvalidArgument, "track is required")
	}

	// Verify peer exists and is in the room
	peer, err := s.repo.GetPeer(ctx, req.PeerId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get peer: %v", err)
	}

	if peer == nil {
		return nil, status.Error(codes.NotFound, "peer not found")
	}

	if peer.RoomId != req.RoomId {
		return nil, status.Error(codes.PermissionDenied, "peer is not in the specified room")
	}

	// Set default state if not specified
	if req.Track.State == entity.TrackState_TRACK_STATE_UNSPECIFIED {
		req.Track.State = entity.TrackState_TRACK_STATE_ACTIVE
	}

	// Ensure peer_id is set
	req.Track.PeerId = req.PeerId

	if err := s.repo.PublishTrack(ctx, req.Track); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to publish track: %v", err)
	}

	return &service.PublishTrackResponse{
		Track: req.Track,
	}, nil
}

// UnpublishTrack unpublishes a media track
func (s *TrackService) UnpublishTrack(ctx context.Context, req *service.UnpublishTrackRequest) (*service.UnpublishTrackResponse, error) {
	if req.TrackId == "" {
		return nil, status.Error(codes.InvalidArgument, "track_id is required")
	}

	// Verify track exists and belongs to the peer
	track, err := s.repo.GetTrack(ctx, req.TrackId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get track: %v", err)
	}

	if track == nil {
		return nil, status.Error(codes.NotFound, "track not found")
	}

	if req.PeerId != "" && track.PeerId != req.PeerId {
		return nil, status.Error(codes.PermissionDenied, "track does not belong to the specified peer")
	}

	if err := s.repo.UnpublishTrack(ctx, req.TrackId); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to unpublish track: %v", err)
	}

	return &service.UnpublishTrackResponse{
		Success: true,
	}, nil
}

// MuteTrack mutes or unmutes a track
func (s *TrackService) MuteTrack(ctx context.Context, req *service.MuteTrackRequest) (*service.MuteTrackResponse, error) {
	if req.TrackId == "" {
		return nil, status.Error(codes.InvalidArgument, "track_id is required")
	}

	// Verify track exists
	track, err := s.repo.GetTrack(ctx, req.TrackId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get track: %v", err)
	}

	if track == nil {
		return nil, status.Error(codes.NotFound, "track not found")
	}

	if req.PeerId != "" && track.PeerId != req.PeerId {
		return nil, status.Error(codes.PermissionDenied, "track does not belong to the specified peer")
	}

	if err := s.repo.MuteTrack(ctx, req.TrackId, req.Muted); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to mute track: %v", err)
	}

	track.Muted = req.Muted

	return &service.MuteTrackResponse{
		Track: track,
	}, nil
}

// GetTrack retrieves track details
func (s *TrackService) GetTrack(ctx context.Context, req *service.GetTrackRequest) (*service.GetTrackResponse, error) {
	if req.TrackId == "" {
		return nil, status.Error(codes.InvalidArgument, "track_id is required")
	}

	track, err := s.repo.GetTrack(ctx, req.TrackId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get track: %v", err)
	}

	if track == nil {
		return nil, status.Error(codes.NotFound, "track not found")
	}

	return &service.GetTrackResponse{
		Track: track,
	}, nil
}

// ListTracks lists tracks in a room
func (s *TrackService) ListTracks(ctx context.Context, req *service.ListTracksRequest) (*service.ListTracksResponse, error) {
	if req.RoomId == "" {
		return nil, status.Error(codes.InvalidArgument, "room_id is required")
	}

	peerID := ""
	if req.PeerId != "" {
		peerID = req.PeerId
	}

	tracks, err := s.repo.ListTracks(ctx, req.RoomId, peerID, req.TypeFilter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list tracks: %v", err)
	}

	return &service.ListTracksResponse{
		Tracks:     tracks,
		TotalCount: int32(len(tracks)),
	}, nil
}

// UpdateTrack updates track settings
func (s *TrackService) UpdateTrack(ctx context.Context, req *service.UpdateTrackRequest) (*service.UpdateTrackResponse, error) {
	if req.TrackId == "" {
		return nil, status.Error(codes.InvalidArgument, "track_id is required")
	}

	track, err := s.repo.GetTrack(ctx, req.TrackId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get track: %v", err)
	}

	if track == nil {
		return nil, status.Error(codes.NotFound, "track not found")
	}

	// Update fields
	if req.Settings != nil {
		track.Settings = req.Settings
	}

	if req.Metadata != nil {
		track.Metadata = req.Metadata
	}

	if err := s.repo.UpdateTrack(ctx, track); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update track: %v", err)
	}

	return &service.UpdateTrackResponse{
		Track: track,
	}, nil
}
