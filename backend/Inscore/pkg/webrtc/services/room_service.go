package services

import (
	"context"
	"fmt"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/webrtc/repository"
	"github.com/newage-saint/insuretech/gen/go/insuretech/webrtc/v1/entity"
	"github.com/newage-saint/insuretech/gen/go/insuretech/webrtc/v1/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RoomService implements the WebRTC RoomService gRPC API
type RoomService struct {
	service.UnimplementedRoomServiceServer
	repo *repository.Repository
}

// NewRoomService creates a new RoomService
func NewRoomService(repo *repository.Repository) *RoomService {
	return &RoomService{repo: repo}
}

// CreateRoom creates a new conference room
func (s *RoomService) CreateRoom(ctx context.Context, req *service.CreateRoomRequest) (*service.CreateRoomResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "room name is required")
	}

	// Set default config if not provided
	config := req.Config
	if config == nil {
		config = &entity.RoomConfig{
			MaxParticipants:       10,
			RequireToken:          false,
			EnableRecording:       false,
			EnableTranscription:   false,
			SessionTimeoutSeconds: 3600,
		}
	}

	// Set default max participants
	maxParticipants := config.GetMaxParticipants()
	if maxParticipants <= 0 {
		maxParticipants = 10
	}

	room := &entity.Room{
		Name:             req.Name,
		Config:           config,
		Metadata:         req.Metadata,
		CreatorId:        req.CreatorId,
		State:            entity.RoomState_ROOM_STATE_ACTIVE,
		MaxParticipants:  maxParticipants,
		ParticipantCount: 0,
	}

	// Create room
	if err := s.repo.CreateRoom(ctx, room); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create room: %v", err)
	}

	// Create room session for analytics
	session := &entity.RoomSession{
		RoomId: room.RoomId, // string ID from entity wrapper
	}
	_ = s.repo.CreateRoomSession(ctx, session)

	// Generate secure join token (JWT-style token would be better in production)
	joinToken := s.generateJoinToken(room.RoomId, config.GetRequireToken())

	return &service.CreateRoomResponse{
		Room:      room,
		JoinToken: joinToken,
	}, nil
}

// GetRoom retrieves room details
func (s *RoomService) GetRoom(ctx context.Context, req *service.GetRoomRequest) (*service.GetRoomResponse, error) {
	if req.RoomId == "" {
		return nil, status.Error(codes.InvalidArgument, "room_id is required")
	}

	room, err := s.repo.GetRoom(ctx, req.RoomId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get room: %v", err)
	}

	if room == nil {
		return nil, status.Error(codes.NotFound, "room not found")
	}

	// Get peers in the room
	peers, err := s.repo.ListPeersInRoom(ctx, req.RoomId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list peers: %v", err)
	}

	return &service.GetRoomResponse{
		Room:  room,
		Peers: peers,
	}, nil
}

// UpdateRoom updates room configuration
func (s *RoomService) UpdateRoom(ctx context.Context, req *service.UpdateRoomRequest) (*service.UpdateRoomResponse, error) {
	if req.RoomId == "" {
		return nil, status.Error(codes.InvalidArgument, "room_id is required")
	}

	room, err := s.repo.GetRoom(ctx, req.RoomId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get room: %v", err)
	}

	if room == nil {
		return nil, status.Error(codes.NotFound, "room not found")
	}

	// Update fields
	if req.Config != nil {
		room.Config = req.Config
	}
	if req.Metadata != nil {
		room.Metadata = req.Metadata
	}

	if err := s.repo.UpdateRoom(ctx, room); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update room: %v", err)
	}

	return &service.UpdateRoomResponse{
		Room: room,
	}, nil
}

// CloseRoom closes a conference room
func (s *RoomService) CloseRoom(ctx context.Context, req *service.CloseRoomRequest) (*service.CloseRoomResponse, error) {
	if req.RoomId == "" {
		return nil, status.Error(codes.InvalidArgument, "room_id is required")
	}

	// Get room to verify it exists
	room, err := s.repo.GetRoom(ctx, req.RoomId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get room: %v", err)
	}

	if room == nil {
		return nil, status.Error(codes.NotFound, "room not found")
	}

	// Close the room
	if err := s.repo.CloseRoom(ctx, req.RoomId); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to close room: %v", err)
	}

	// End active room session
	session, err := s.repo.GetActiveRoomSession(ctx, req.RoomId)
	if err == nil && session != nil {
		_ = s.repo.EndRoomSession(ctx, session.SessionId)
		_ = s.repo.UpdateSessionMetrics(ctx, session.SessionId)
	}

	// Mark all peers as left
	peers, err := s.repo.ListPeersInRoom(ctx, req.RoomId)
	if err == nil {
		for _, peer := range peers {
			_ = s.repo.RemovePeer(ctx, peer.PeerId)
		}
	}

	return &service.CloseRoomResponse{
		Success: true,
	}, nil
}

// ListRooms lists all rooms with optional filtering
func (s *RoomService) ListRooms(ctx context.Context, req *service.ListRoomsRequest) (*service.ListRoomsResponse, error) {
	pageSize := req.PageSize
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	// Parse page token to get offset (simple implementation, JWT-based would be better)
	offset := 0
	if req.PageToken != "" {
		if parsedOffset, err := s.parsePageToken(req.PageToken); err == nil {
			offset = parsedOffset
		}
	}

	rooms, total, err := s.repo.ListRooms(ctx, int(pageSize), offset, req.StateFilter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list rooms: %v", err)
	}

	// Generate next page token if there are more results
	var nextPageToken string
	if int64(offset+int(pageSize)) < total {
		nextPageToken = s.generatePageToken(offset + int(pageSize))
	}

	return &service.ListRoomsResponse{
		Rooms:         rooms,
		NextPageToken: nextPageToken,
		TotalCount:    int32(total),
	}, nil
}

// generateJoinToken generates a join token for a room
func (s *RoomService) generateJoinToken(roomID string, required bool) string {
	if !required {
		return "" // No token required
	}
	// Simple token generation - in production, use JWT with expiration
	return fmt.Sprintf("tok_%s_%d", roomID, time.Now().Unix())
}

// parsePageToken parses a page token to extract the offset
func (s *RoomService) parsePageToken(token string) (int, error) {
	var offset int
	_, err := fmt.Sscanf(token, "offset_%d", &offset)
	return offset, err
}

// generatePageToken generates a page token for pagination
func (s *RoomService) generatePageToken(offset int) string {
	return fmt.Sprintf("offset_%d", offset)
}
