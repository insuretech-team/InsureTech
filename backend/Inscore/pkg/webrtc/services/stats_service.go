package services

import (
	"context"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/webrtc/repository"
	"github.com/newage-saint/insuretech/gen/go/insuretech/webrtc/v1/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// StatsService implements the WebRTC StatsService gRPC API
type StatsService struct {
	service.UnimplementedStatsServiceServer
	repo *repository.Repository
}

// NewStatsService creates a new StatsService
func NewStatsService(repo *repository.Repository) *StatsService {
	return &StatsService{repo: repo}
}

// GetConnectionStats retrieves connection statistics for a peer
func (s *StatsService) GetConnectionStats(ctx context.Context, req *service.GetConnectionStatsRequest) (*service.GetConnectionStatsResponse, error) {
	if req.PeerId == "" {
		return nil, status.Error(codes.InvalidArgument, "peer_id is required")
	}

	// Verify peer exists
	peer, err := s.repo.GetPeer(ctx, req.PeerId)
	if err != nil || peer == nil {
		return nil, status.Error(codes.NotFound, "peer not found")
	}

	// Build connection stats
	stats := &service.ConnectionStats{
		PeerId:            req.PeerId,
		BytesSent:         0,
		BytesReceived:     0,
		BitrateKbps:       0,
		PacketLossPercent: 0,
		RttMs:             0,
	}

	return &service.GetConnectionStatsResponse{Stats: stats}, nil
}

// StreamStats streams real-time statistics
func (s *StatsService) StreamStats(req *service.StreamStatsRequest, stream service.StatsService_StreamStatsServer) error {
	if req.PeerId == "" {
		return status.Error(codes.InvalidArgument, "peer_id is required")
	}

	interval := time.Duration(1000) * time.Millisecond // Default to 1 second
	// req.IntervalSeconds is int32, but we'll stick to a default for now if it's 0
	if req.IntervalSeconds > 0 {
		interval = time.Duration(req.IntervalSeconds) * time.Second
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-stream.Context().Done():
			return stream.Context().Err()
		case <-ticker.C:
			// Get current stats
			statsResp, err := s.GetConnectionStats(stream.Context(), &service.GetConnectionStatsRequest{
				RoomId: req.RoomId,
				PeerId: req.PeerId,
			})
			if err != nil {
				return err
			}

			// Send stats update
			update := &service.StreamStatsResponse{
				Stats: statsResp.Stats,
			}

			if err := stream.Send(update); err != nil {
				return err
			}
		}
	}
}

// GetRoomAnalytics retrieves analytics for a room
func (s *StatsService) GetRoomAnalytics(ctx context.Context, req *service.GetRoomAnalyticsRequest) (*service.GetRoomAnalyticsResponse, error) {
	if req.RoomId == "" {
		return nil, status.Error(codes.InvalidArgument, "room_id is required")
	}

	// Get room
	room, err := s.repo.GetRoom(ctx, req.RoomId)
	if err != nil || room == nil {
		return nil, status.Error(codes.NotFound, "room not found")
	}

	// Get room session
	_, err = s.repo.GetActiveRoomSession(ctx, req.RoomId)
	if err != nil {
		// Return empty analytics if no session
		return &service.GetRoomAnalyticsResponse{
			Analytics: &service.RoomAnalytics{
				RoomId: req.RoomId,
			},
		}, nil
	}

	analytics := &service.RoomAnalytics{
		RoomId:           req.RoomId,
		PeakParticipants: 0, // Not tracked in new schema
		TotalSessions:    0, // Not available in session directly
	}

	return &service.GetRoomAnalyticsResponse{Analytics: analytics}, nil
}

// GetPeerAnalytics retrieves analytics for a specific peer
func (s *StatsService) GetPeerAnalytics(ctx context.Context, req *service.GetPeerAnalyticsRequest) (*service.GetPeerAnalyticsResponse, error) {
	if req.PeerId == "" {
		return nil, status.Error(codes.InvalidArgument, "peer_id is required")
	}

	// Get peer
	peer, err := s.repo.GetPeer(ctx, req.PeerId)
	if err != nil || peer == nil {
		return nil, status.Error(codes.NotFound, "peer not found")
	}

	// Calculate duration
	var duration int64
	if peer.LeftAt != nil {
		duration = peer.LeftAt.AsTime().Unix() - peer.JoinedAt.AsTime().Unix()
	} else {
		duration = time.Now().Unix() - peer.JoinedAt.AsTime().Unix()
	}

	analytics := &service.ParticipantAnalytics{
		PeerId:           req.PeerId,
		TotalTimeSeconds: duration,
	}

	return &service.GetPeerAnalyticsResponse{Analytics: analytics}, nil
}
