package services

import (
	"context"
	"testing"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/webrtc/repository"

	"github.com/newage-saint/insuretech/gen/go/insuretech/webrtc/v1/entity"
	"github.com/newage-saint/insuretech/gen/go/insuretech/webrtc/v1/service"
)

func TestStatsService_GetConnectionStats(t *testing.T) {
	ctx := context.Background()
	defer cleanupTestData(ctx)

	repo := repository.NewRepository(getTestDB())
	roomService := NewRoomService(repo)
	peerService := NewPeerService(repo)
	trackService := NewTrackService(repo)
	statsService := NewStatsService(repo)

	// Create test room and peer
	roomResp, _ := roomService.CreateRoom(ctx, &service.CreateRoomRequest{Name: "Test Room"})
	joinResp, err := peerService.JoinRoom(ctx, &service.JoinRoomRequest{
		RoomId:      roomResp.Room.RoomId,
		DisplayName: "Test User",
	})
	if err != nil {
		t.Fatalf("Failed to join room: %v", err)
	}

	// Publish some tracks
	trackService.PublishTrack(ctx, &service.PublishTrackRequest{
		RoomId: roomResp.Room.RoomId,
		PeerId: joinResp.Peer.PeerId,
		Track: &entity.Track{
			TrackId: "audio-1",
			Type:    entity.TrackType_TRACK_TYPE_AUDIO,
		},
	})
	trackService.PublishTrack(ctx, &service.PublishTrackRequest{
		RoomId: roomResp.Room.RoomId,
		PeerId: joinResp.Peer.PeerId,
		Track: &entity.Track{
			TrackId: "video-1",
			Type:    entity.TrackType_TRACK_TYPE_VIDEO,
		},
	})

	t.Run("GetConnectionStats_Success", func(t *testing.T) {
		req := &service.GetConnectionStatsRequest{
			RoomId: roomResp.Room.RoomId,
			PeerId: joinResp.Peer.PeerId,
		}

		resp, err := statsService.GetConnectionStats(ctx, req)
		if err != nil {
			t.Fatalf("GetConnectionStats failed: %v", err)
		}

		if resp.Stats == nil {
			t.Fatal("Stats is nil")
		}

		if resp.Stats.PeerId != joinResp.Peer.PeerId {
			t.Errorf("Expected peer ID %s, got %s", joinResp.Peer.PeerId, resp.Stats.PeerId)
		}

	})

	t.Run("GetConnectionStats_MissingRoomId", func(t *testing.T) {
		req := &service.GetConnectionStatsRequest{
			PeerId: joinResp.Peer.PeerId,
		}

		_, err := statsService.GetConnectionStats(ctx, req)
		if err == nil {
			t.Fatal("Expected error for missing room ID")
		}
	})

	t.Run("GetConnectionStats_MissingPeerId", func(t *testing.T) {
		req := &service.GetConnectionStatsRequest{
			RoomId: roomResp.Room.RoomId,
		}

		_, err := statsService.GetConnectionStats(ctx, req)
		if err == nil {
			t.Fatal("Expected error for missing peer ID")
		}
	})

	t.Run("GetConnectionStats_PeerNotFound", func(t *testing.T) {
		req := &service.GetConnectionStatsRequest{
			RoomId: roomResp.Room.RoomId,
			PeerId: "00000000-0000-0000-0000-000000000000",
		}

		_, err := statsService.GetConnectionStats(ctx, req)
		if err == nil {
			t.Fatal("Expected error for non-existent peer")
		}
	})
}

func TestStatsService_GetRoomAnalytics(t *testing.T) {
	ctx := context.Background()
	defer cleanupTestData(ctx)

	repo := repository.NewRepository(getTestDB())
	roomService := NewRoomService(repo)
	peerService := NewPeerService(repo)
	statsService := NewStatsService(repo)

	// Create test room
	roomResp, _ := roomService.CreateRoom(ctx, &service.CreateRoomRequest{Name: "Test Room"})

	// Add multiple peers
	for i := 1; i <= 3; i++ {
		_, err := peerService.JoinRoom(ctx, &service.JoinRoomRequest{
			RoomId:      roomResp.Room.RoomId,
			DisplayName: "User " + string(rune('0'+i)),
		})
		if err != nil {
			t.Fatalf("Failed to join peer %d: %v", i, err)
		}
	}

	t.Run("GetRoomAnalytics_Success", func(t *testing.T) {
		req := &service.GetRoomAnalyticsRequest{
			RoomId: roomResp.Room.RoomId,
		}

		resp, err := statsService.GetRoomAnalytics(ctx, req)
		if err != nil {
			t.Fatalf("GetRoomAnalytics failed: %v", err)
		}

		if resp.Analytics == nil {
			t.Fatal("Analytics is nil")
		}

		if resp.Analytics.RoomId != roomResp.Room.RoomId {
			t.Errorf("Expected room ID %s, got %s", roomResp.Room.RoomId, resp.Analytics.RoomId)
		}

	})

	t.Run("GetRoomAnalytics_MissingRoomId", func(t *testing.T) {
		req := &service.GetRoomAnalyticsRequest{}

		_, err := statsService.GetRoomAnalytics(ctx, req)
		if err == nil {
			t.Fatal("Expected error for missing room ID")
		}
	})

	t.Run("GetRoomAnalytics_RoomNotFound", func(t *testing.T) {
		req := &service.GetRoomAnalyticsRequest{
			RoomId: "00000000-0000-0000-0000-000000000000",
		}

		_, err := statsService.GetRoomAnalytics(ctx, req)
		if err == nil {
			t.Fatal("Expected error for non-existent room")
		}
	})
}

func TestStatsService_GetPeerAnalytics(t *testing.T) {
	ctx := context.Background()
	defer cleanupTestData(ctx)

	repo := repository.NewRepository(getTestDB())
	roomService := NewRoomService(repo)
	peerService := NewPeerService(repo)
	statsService := NewStatsService(repo)

	// Create test room with peer
	roomResp, _ := roomService.CreateRoom(ctx, &service.CreateRoomRequest{Name: "Test Room"})
	joinResp, _ := peerService.JoinRoom(ctx, &service.JoinRoomRequest{
		RoomId:      roomResp.Room.RoomId,
		DisplayName: "Test User",
	})

	t.Run("GetPeerAnalytics_Success", func(t *testing.T) {
		req := &service.GetPeerAnalyticsRequest{
			PeerId: joinResp.Peer.PeerId,
		}

		resp, err := statsService.GetPeerAnalytics(ctx, req)
		if err != nil {
			t.Fatalf("GetPeerAnalytics failed: %v", err)
		}

		if resp.Analytics == nil {
			t.Fatal("Analytics is nil")
		}

		// Stats should have basic information even without time range
		if resp.Analytics.PeerId != joinResp.Peer.PeerId {
			t.Errorf("Expected peer ID %s, got %s", joinResp.Peer.PeerId, resp.Analytics.PeerId)
		}
	})

	t.Run("GetPeerAnalytics_MissingPeerId", func(t *testing.T) {
		req := &service.GetPeerAnalyticsRequest{}

		_, err := statsService.GetPeerAnalytics(ctx, req)
		if err == nil {
			t.Fatal("Expected error for missing peer ID")
		}
	})

	t.Run("GetPeerAnalytics_PeerNotFound", func(t *testing.T) {
		req := &service.GetPeerAnalyticsRequest{
			PeerId: "00000000-0000-0000-0000-000000000000",
		}

		_, err := statsService.GetPeerAnalytics(ctx, req)
		if err == nil {
			t.Fatal("Expected error for non-existent peer")
		}
	})
}
