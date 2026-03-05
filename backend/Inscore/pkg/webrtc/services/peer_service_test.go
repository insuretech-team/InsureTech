package services

import (
	"context"
	"testing"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/webrtc/repository"

	"github.com/newage-saint/insuretech/gen/go/insuretech/webrtc/v1/entity"
	"github.com/newage-saint/insuretech/gen/go/insuretech/webrtc/v1/service"
)

func TestPeerService_JoinRoom(t *testing.T) {
	ctx := context.Background()
	defer cleanupTestData(ctx)

	repo := repository.NewRepository(getTestDB())
	roomService := NewRoomService(repo)
	peerService := NewPeerService(repo)

	// Create a test room
	createRoomReq := &service.CreateRoomRequest{
		Name: "Test Room",
		Config: &entity.RoomConfig{
			MaxParticipants: 5,
		},
	}
	roomResp, err := roomService.CreateRoom(ctx, createRoomReq)
	if err != nil {
		t.Fatalf("Failed to create test room: %v", err)
	}

	t.Run("JoinRoom_Success", func(t *testing.T) {
		req := &service.JoinRoomRequest{
			RoomId:      roomResp.Room.RoomId,
			DisplayName: "Test User",
			Metadata:    map[string]string{"user_agent": "Test Browser"},
		}

		resp, err := peerService.JoinRoom(ctx, req)
		if err != nil {
			t.Fatalf("JoinRoom failed: %v", err)
		}

		if resp.Peer == nil {
			t.Fatal("Peer is nil")
		}

		if resp.Peer.DisplayName != "Test User" {
			t.Errorf("Expected display name 'Test User', got '%s'", resp.Peer.DisplayName)
		}

		if resp.Peer.PeerId == "" {
			t.Error("Peer ID should be set")
		}

		if resp.Peer.State != entity.PeerConnectionState_PEER_CONNECTION_STATE_NEW {
			t.Errorf("Expected peer state NEW, got %v", resp.Peer.State)
		}
	})

	t.Run("JoinRoom_MissingRoomId", func(t *testing.T) {
		req := &service.JoinRoomRequest{
			DisplayName: "Test User",
		}

		_, err := peerService.JoinRoom(ctx, req)
		if err == nil {
			t.Fatal("Expected error for missing room ID")
		}
	})

	t.Run("JoinRoom_MissingDisplayName", func(t *testing.T) {
		req := &service.JoinRoomRequest{
			RoomId: roomResp.Room.RoomId,
		}

		_, err := peerService.JoinRoom(ctx, req)
		if err == nil {
			t.Fatal("Expected error for missing display name")
		}
	})

	t.Run("JoinRoom_RoomNotFound", func(t *testing.T) {
		req := &service.JoinRoomRequest{
			RoomId:      "00000000-0000-0000-0000-000000000000",
			DisplayName: "Test User",
		}

		_, err := peerService.JoinRoom(ctx, req)
		if err == nil {
			t.Fatal("Expected error for non-existent room")
		}
	})

	t.Run("JoinRoom_RoomFull", func(t *testing.T) {
		// Create a room with max 1 participant
		smallRoomReq := &service.CreateRoomRequest{
			Name: "Small Room",
			Config: &entity.RoomConfig{
				MaxParticipants: 1,
			},
		}
		smallRoomResp, err := roomService.CreateRoom(ctx, smallRoomReq)
		if err != nil {
			t.Fatalf("Failed to create small room: %v", err)
		}

		// Join first peer
		joinReq1 := &service.JoinRoomRequest{
			RoomId:      smallRoomResp.Room.RoomId,
			DisplayName: "User 1",
		}
		_, err = peerService.JoinRoom(ctx, joinReq1)
		if err != nil {
			t.Fatalf("First join failed: %v", err)
		}

		// Try to join second peer (should fail)
		joinReq2 := &service.JoinRoomRequest{
			RoomId:      smallRoomResp.Room.RoomId,
			DisplayName: "User 2",
		}
		_, err = peerService.JoinRoom(ctx, joinReq2)
		if err == nil {
			t.Fatal("Expected error for room full")
		}
	})
}

func TestPeerService_GetPeer(t *testing.T) {
	ctx := context.Background()
	defer cleanupTestData(ctx)

	repo := repository.NewRepository(getTestDB())
	roomService := NewRoomService(repo)
	peerService := NewPeerService(repo)

	// Create test room and peer
	roomResp, _ := roomService.CreateRoom(ctx, &service.CreateRoomRequest{Name: "Test Room"})
	joinResp, err := peerService.JoinRoom(ctx, &service.JoinRoomRequest{
		RoomId:      roomResp.Room.RoomId,
		DisplayName: "Test User",
	})
	if err != nil {
		t.Fatalf("Failed to join room: %v", err)
	}

	t.Run("GetPeer_Success", func(t *testing.T) {
		req := &service.GetPeerRequest{
			PeerId: joinResp.Peer.PeerId,
		}

		resp, err := peerService.GetPeer(ctx, req)
		if err != nil {
			t.Fatalf("GetPeer failed: %v", err)
		}

		if resp.Peer.PeerId != joinResp.Peer.PeerId {
			t.Errorf("Expected peer ID %s, got %s", joinResp.Peer.PeerId, resp.Peer.PeerId)
		}

		if resp.Peer.DisplayName != "Test User" {
			t.Errorf("Expected display name 'Test User', got '%s'", resp.Peer.DisplayName)
		}
	})

	t.Run("GetPeer_NotFound", func(t *testing.T) {
		req := &service.GetPeerRequest{
			PeerId: "00000000-0000-0000-0000-000000000000",
		}

		_, err := peerService.GetPeer(ctx, req)
		if err == nil {
			t.Fatal("Expected error for non-existent peer")
		}
	})
}

func TestPeerService_UpdatePeer(t *testing.T) {
	ctx := context.Background()
	defer cleanupTestData(ctx)

	repo := repository.NewRepository(getTestDB())
	roomService := NewRoomService(repo)
	peerService := NewPeerService(repo)

	// Create test room and peer
	roomResp, _ := roomService.CreateRoom(ctx, &service.CreateRoomRequest{Name: "Test Room"})
	joinResp, _ := peerService.JoinRoom(ctx, &service.JoinRoomRequest{
		RoomId:      roomResp.Room.RoomId,
		DisplayName: "Test User",
	})

	t.Run("UpdatePeer_Success", func(t *testing.T) {
		req := &service.UpdatePeerRequest{
			PeerId:      joinResp.Peer.PeerId,
			DisplayName: "Updated User",
			Metadata:    map[string]string{"status": "active"},
		}

		resp, err := peerService.UpdatePeer(ctx, req)
		if err != nil {
			t.Fatalf("UpdatePeer failed: %v", err)
		}

		if resp.Peer.DisplayName != "Updated User" {
			t.Errorf("Expected display name 'Updated User', got '%s'", resp.Peer.DisplayName)
		}

		// Verify update persisted
		getResp, err := peerService.GetPeer(ctx, &service.GetPeerRequest{
			PeerId: joinResp.Peer.PeerId,
		})
		if err != nil {
			t.Fatalf("GetPeer failed: %v", err)
		}

		if getResp.Peer.DisplayName != "Updated User" {
			t.Errorf("Expected display name 'Updated User', got '%s'", getResp.Peer.DisplayName)
		}
	})
}

func TestPeerService_LeaveRoom(t *testing.T) {
	ctx := context.Background()
	defer cleanupTestData(ctx)

	repo := repository.NewRepository(getTestDB())
	roomService := NewRoomService(repo)
	peerService := NewPeerService(repo)

	// Create test room and peer
	roomResp, _ := roomService.CreateRoom(ctx, &service.CreateRoomRequest{Name: "Test Room"})
	joinResp, _ := peerService.JoinRoom(ctx, &service.JoinRoomRequest{
		RoomId:      roomResp.Room.RoomId,
		DisplayName: "Test User",
	})

	t.Run("LeaveRoom_Success", func(t *testing.T) {
		req := &service.LeaveRoomRequest{
			RoomId: roomResp.Room.RoomId,
			PeerId: joinResp.Peer.PeerId,
		}

		resp, err := peerService.LeaveRoom(ctx, req)
		if err != nil {
			t.Fatalf("LeaveRoom failed: %v", err)
		}

		if !resp.Success {
			t.Error("Expected success to be true")
		}

		// Verify peer state is closed
		getResp, err := peerService.GetPeer(ctx, &service.GetPeerRequest{
			PeerId: joinResp.Peer.PeerId,
		})
		if err != nil {
			t.Fatalf("GetPeer failed: %v", err)
		}

		if getResp.Peer.State != entity.PeerConnectionState_PEER_CONNECTION_STATE_CLOSED {
			t.Errorf("Expected peer state CLOSED, got %v", getResp.Peer.State)
		}
	})
}

func TestPeerService_ListPeers(t *testing.T) {
	ctx := context.Background()
	defer cleanupTestData(ctx)

	repo := repository.NewRepository(getTestDB())
	roomService := NewRoomService(repo)
	peerService := NewPeerService(repo)

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

	t.Run("ListPeers_Success", func(t *testing.T) {
		req := &service.ListPeersRequest{
			RoomId: roomResp.Room.RoomId,
		}

		resp, err := peerService.ListPeers(ctx, req)
		if err != nil {
			t.Fatalf("ListPeers failed: %v", err)
		}

		if len(resp.Peers) != 3 {
			t.Errorf("Expected 3 peers, got %d", len(resp.Peers))
		}
	})
}
