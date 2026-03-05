package services

import (
	"context"
	"encoding/json"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/webrtc/repository"
	"github.com/newage-saint/insuretech/gen/go/insuretech/webrtc/v1/entity"
	"github.com/newage-saint/insuretech/gen/go/insuretech/webrtc/v1/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// PeerService implements the WebRTC PeerService gRPC API
type PeerService struct {
	service.UnimplementedPeerServiceServer
	repo *repository.Repository
}

// NewPeerService creates a new PeerService
func NewPeerService(repo *repository.Repository) *PeerService {
	return &PeerService{repo: repo}
}

// JoinRoom allows a peer to join a conference room
func (s *PeerService) JoinRoom(ctx context.Context, req *service.JoinRoomRequest) (*service.JoinRoomResponse, error) {
	if req.RoomId == "" {
		return nil, status.Error(codes.InvalidArgument, "room_id is required")
	}

	if req.DisplayName == "" {
		return nil, status.Error(codes.InvalidArgument, "display_name is required")
	}

	// Verify room exists and is active
	room, err := s.repo.GetRoom(ctx, req.RoomId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get room: %v", err)
	}

	if room == nil {
		return nil, status.Error(codes.NotFound, "room not found")
	}

	if room.State == entity.RoomState_ROOM_STATE_CLOSED {
		return nil, status.Error(codes.FailedPrecondition, "room is closed")
	}

	// Check if room is full
	if room.ParticipantCount >= room.MaxParticipants {
		return nil, status.Error(codes.ResourceExhausted, "room is full")
	}

	// Validate join_token if room requires token
	if room.Config.GetRequireToken() {
		if req.JoinToken == "" {
			return nil, status.Error(codes.Unauthenticated, "join_token is required for this room")
		}
		// Simple token validation - in production, verify JWT signature and expiration
		if !s.validateJoinToken(req.JoinToken, req.RoomId) {
			return nil, status.Error(codes.Unauthenticated, "invalid join_token")
		}
	}

	// Create peer
	peer := &entity.Peer{
		RoomId:      req.RoomId,
		DisplayName: req.DisplayName,
		Metadata:    req.Metadata,
		State:       entity.PeerConnectionState_PEER_CONNECTION_STATE_NEW,
	}

	// Extract user agent from metadata if available
	if userAgent, ok := req.Metadata["user_agent"]; ok {
		peer.UserAgent = userAgent
	}

	if err := s.repo.AddPeer(ctx, peer); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add peer: %v", err)
	}

	// Increment participant count
	room.ParticipantCount++
	if err := s.repo.UpdateRoom(ctx, room); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update room: %v", err)
	}

	// Create peer session for analytics
	session, err := s.repo.GetActiveRoomSession(ctx, req.RoomId)
	if err == nil && session != nil {
		peerSession := &entity.PeerSession{
			SessionId: session.SessionId, // string
			PeerId:    peer.PeerId,       // string
		}
		_ = s.repo.CreatePeerSession(ctx, peerSession)

		// Update session metrics
		_ = s.repo.UpdateSessionMetrics(ctx, session.SessionId)
	}

	// Get existing peers in the room
	existingPeers, err := s.repo.ListPeersInRoom(ctx, req.RoomId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list peers: %v", err)
	}

	// Filter out the newly joined peer from existing peers
	var otherPeers []*entity.Peer
	for _, p := range existingPeers {
		if p.PeerId != peer.PeerId {
			otherPeers = append(otherPeers, p)
		}
	}

	return &service.JoinRoomResponse{
		Peer:          peer,
		Room:          room,
		ExistingPeers: otherPeers,
	}, nil
}

// LeaveRoom allows a peer to leave a conference room
func (s *PeerService) LeaveRoom(ctx context.Context, req *service.LeaveRoomRequest) (*service.LeaveRoomResponse, error) {
	if req.RoomId == "" {
		return nil, status.Error(codes.InvalidArgument, "room_id is required")
	}

	if req.PeerId == "" {
		return nil, status.Error(codes.InvalidArgument, "peer_id is required")
	}

	// Remove peer
	if err := s.repo.RemovePeer(ctx, req.PeerId); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to remove peer: %v", err)
	}

	// End peer session for analytics
	session, err := s.repo.GetActiveRoomSession(ctx, req.RoomId)
	if err == nil && session != nil {
		peerSession, err := s.repo.GetActivePeerSession(ctx, session.SessionId, req.PeerId)
		if err == nil && peerSession != nil {
			_ = s.repo.EndPeerSession(ctx, peerSession.PeerSessionId)
		}
	}

	// Decrement participant count
	room, err := s.repo.GetRoom(ctx, req.RoomId)
	if err == nil && room != nil {
		if room.ParticipantCount > 0 {
			room.ParticipantCount--
		}
		_ = s.repo.UpdateRoom(ctx, room)
	}

	return &service.LeaveRoomResponse{
		Success: true,
	}, nil
}

// validateJoinToken validates a join token for a room
func (s *PeerService) validateJoinToken(token, roomID string) bool {
	// Simple token validation - in production, verify JWT signature and expiration
	// Check if token starts with "tok_" and contains the room ID
	expectedPrefix := "tok_" + roomID
	return len(token) > len(expectedPrefix) && token[:len(expectedPrefix)] == expectedPrefix
}

// GetPeer retrieves peer details
func (s *PeerService) GetPeer(ctx context.Context, req *service.GetPeerRequest) (*service.GetPeerResponse, error) {
	if req.PeerId == "" {
		return nil, status.Error(codes.InvalidArgument, "peer_id is required")
	}

	peer, err := s.repo.GetPeer(ctx, req.PeerId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get peer: %v", err)
	}

	if peer == nil {
		return nil, status.Error(codes.NotFound, "peer not found")
	}

	return &service.GetPeerResponse{
		Peer: peer,
	}, nil
}

// UpdatePeer updates peer information
func (s *PeerService) UpdatePeer(ctx context.Context, req *service.UpdatePeerRequest) (*service.UpdatePeerResponse, error) {
	if req.PeerId == "" {
		return nil, status.Error(codes.InvalidArgument, "peer_id is required")
	}

	peer, err := s.repo.GetPeer(ctx, req.PeerId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get peer: %v", err)
	}

	if peer == nil {
		return nil, status.Error(codes.NotFound, "peer not found")
	}

	// Update fields
	if req.DisplayName != "" {
		peer.DisplayName = req.DisplayName
	}

	if req.Metadata != nil {
		peer.Metadata = req.Metadata
	}

	// Update in database
	metadataJSON, err := json.Marshal(peer.Metadata)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to marshal metadata: %v", err)
	}

	query := `
		UPDATE webrtc.peers
		SET display_name = $1, metadata = $2, last_seen_at = CURRENT_TIMESTAMP
		WHERE peer_id = $3`

	if err := s.repo.GetDB().WithContext(ctx).Exec(query, peer.DisplayName, metadataJSON, peer.PeerId).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update peer: %v", err)
	}

	return &service.UpdatePeerResponse{
		Peer: peer,
	}, nil
}

// ListPeers lists all peers in a room
func (s *PeerService) ListPeers(ctx context.Context, req *service.ListPeersRequest) (*service.ListPeersResponse, error) {
	if req.RoomId == "" {
		return nil, status.Error(codes.InvalidArgument, "room_id is required")
	}

	peers, err := s.repo.ListPeersInRoom(ctx, req.RoomId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list peers: %v", err)
	}

	// Filter by state if specified
	if req.StateFilter != entity.PeerConnectionState_PEER_CONNECTION_STATE_UNSPECIFIED {
		var filtered []*entity.Peer
		for _, p := range peers {
			if p.State == req.StateFilter {
				filtered = append(filtered, p)
			}
		}
		peers = filtered
	}

	return &service.ListPeersResponse{
		Peers:      peers,
		TotalCount: int32(len(peers)),
	}, nil
}
