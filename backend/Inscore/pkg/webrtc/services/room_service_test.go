package services

import (
	"context"
	"testing"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/webrtc/repository"
	"github.com/newage-saint/insuretech/gen/go/insuretech/webrtc/v1/entity"
	"github.com/newage-saint/insuretech/gen/go/insuretech/webrtc/v1/service"
)

func TestRoomService_CreateRoom(t *testing.T) {
	ctx := context.Background()
	defer cleanupTestData(ctx)

	repo := repository.NewRepository(getTestDB())
	roomService := NewRoomService(repo)

	t.Run("CreateRoom_Success", func(t *testing.T) {
		req := &service.CreateRoomRequest{
			Name: "Test Room",
			Config: &entity.RoomConfig{
				MaxParticipants: 10,
			},
		}

		resp, err := roomService.CreateRoom(ctx, req)
		if err != nil {
			t.Fatalf("CreateRoom failed: %v", err)
		}

		if resp.Room == nil {
			t.Fatal("Room is nil")
		}

		if resp.Room.Name != "Test Room" {
			t.Errorf("Expected room name 'Test Room', got '%s'", resp.Room.Name)
		}

		if resp.Room.RoomId == "" {
			t.Error("Room ID should be set")
		}

		if resp.Room.State != entity.RoomState_ROOM_STATE_ACTIVE {
			t.Errorf("Expected room state ACTIVE, got %v", resp.Room.State)
		}
	})

	t.Run("CreateRoom_WithDefaultConfig", func(t *testing.T) {
		req := &service.CreateRoomRequest{
			Name: "Room with Defaults",
		}

		resp, err := roomService.CreateRoom(ctx, req)
		if err != nil {
			t.Fatalf("CreateRoom failed: %v", err)
		}

		if resp.Room.MaxParticipants != 10 {
			t.Errorf("Expected default max participants 10, got %d", resp.Room.MaxParticipants)
		}
	})

	t.Run("CreateRoom_MissingName", func(t *testing.T) {
		req := &service.CreateRoomRequest{
			Name: "",
		}

		_, err := roomService.CreateRoom(ctx, req)
		if err == nil {
			t.Fatal("Expected error for missing room name")
		}
	})
}

func TestRoomService_GetRoom(t *testing.T) {
	ctx := context.Background()
	defer cleanupTestData(ctx)

	repo := repository.NewRepository(getTestDB())
	roomService := NewRoomService(repo)

	// Create a test room
	createReq := &service.CreateRoomRequest{
		Name: "Test Room",
	}
	createResp, err := roomService.CreateRoom(ctx, createReq)
	if err != nil {
		t.Fatalf("Failed to create test room: %v", err)
	}

	t.Run("GetRoom_Success", func(t *testing.T) {
		req := &service.GetRoomRequest{
			RoomId: createResp.Room.RoomId,
		}

		resp, err := roomService.GetRoom(ctx, req)
		if err != nil {
			t.Fatalf("GetRoom failed: %v", err)
		}

		if resp.Room.RoomId != createResp.Room.RoomId {
			t.Errorf("Expected room ID %s, got %s", createResp.Room.RoomId, resp.Room.RoomId)
		}

		if resp.Room.Name != "Test Room" {
			t.Errorf("Expected room name 'Test Room', got '%s'", resp.Room.Name)
		}
	})

	t.Run("GetRoom_NotFound", func(t *testing.T) {
		req := &service.GetRoomRequest{
			RoomId: "00000000-0000-0000-0000-000000000000",
		}

		_, err := roomService.GetRoom(ctx, req)
		if err == nil {
			t.Fatal("Expected error for non-existent room")
		}
	})

	t.Run("GetRoom_MissingRoomId", func(t *testing.T) {
		req := &service.GetRoomRequest{}

		_, err := roomService.GetRoom(ctx, req)
		if err == nil {
			t.Fatal("Expected error for missing room ID")
		}
	})
}

func TestRoomService_UpdateRoom(t *testing.T) {
	ctx := context.Background()
	defer cleanupTestData(ctx)

	repo := repository.NewRepository(getTestDB())
	roomService := NewRoomService(repo)

	// Create a test room
	createReq := &service.CreateRoomRequest{
		Name: "Test Room",
	}
	createResp, err := roomService.CreateRoom(ctx, createReq)
	if err != nil {
		t.Fatalf("Failed to create test room: %v", err)
	}

	t.Run("UpdateRoom_Success", func(t *testing.T) {
		req := &service.UpdateRoomRequest{
			RoomId: createResp.Room.RoomId,
			Config: &entity.RoomConfig{
				MaxParticipants: 20,
			},
			Metadata: map[string]string{"updated": "true"},
		}

		resp, err := roomService.UpdateRoom(ctx, req)
		if err != nil {
			t.Fatalf("UpdateRoom failed: %v", err)
		}

		if resp.Room.MaxParticipants != 20 {
			t.Errorf("Expected max participants 20, got %d", resp.Room.MaxParticipants)
		}
	})

	t.Run("UpdateRoom_NotFound", func(t *testing.T) {
		req := &service.UpdateRoomRequest{
			RoomId: "00000000-0000-0000-0000-000000000000",
			Config: &entity.RoomConfig{MaxParticipants: 10},
		}

		_, err := roomService.UpdateRoom(ctx, req)
		if err == nil {
			t.Fatal("Expected error for non-existent room")
		}
	})
}

func TestRoomService_CloseRoom(t *testing.T) {
	ctx := context.Background()
	defer cleanupTestData(ctx)

	repo := repository.NewRepository(getTestDB())
	roomService := NewRoomService(repo)

	// Create a test room
	createReq := &service.CreateRoomRequest{
		Name: "Test Room",
	}
	createResp, err := roomService.CreateRoom(ctx, createReq)
	if err != nil {
		t.Fatalf("Failed to create test room: %v", err)
	}

	t.Run("CloseRoom_Success", func(t *testing.T) {
		req := &service.CloseRoomRequest{
			RoomId: createResp.Room.RoomId,
		}

		resp, err := roomService.CloseRoom(ctx, req)
		if err != nil {
			t.Fatalf("CloseRoom failed: %v", err)
		}

		if !resp.Success {
			t.Error("Expected success to be true")
		}

		// Verify room is closed
		getReq := &service.GetRoomRequest{
			RoomId: createResp.Room.RoomId,
		}
		getResp, err := roomService.GetRoom(ctx, getReq)
		if err != nil {
			t.Fatalf("GetRoom failed: %v", err)
		}

		if getResp.Room.State != entity.RoomState_ROOM_STATE_CLOSED {
			t.Errorf("Expected room state CLOSED, got %v", getResp.Room.State)
		}
	})

	t.Run("CloseRoom_NotFound", func(t *testing.T) {
		req := &service.CloseRoomRequest{
			RoomId: "00000000-0000-0000-0000-000000000000",
		}

		_, err := roomService.CloseRoom(ctx, req)
		if err == nil {
			t.Fatal("Expected error for non-existent room")
		}
	})
}

func TestRoomService_ListRooms(t *testing.T) {
	ctx := context.Background()
	defer cleanupTestData(ctx)

	repo := repository.NewRepository(getTestDB())
	roomService := NewRoomService(repo)

	// Create multiple test rooms
	for i := 1; i <= 5; i++ {
		createReq := &service.CreateRoomRequest{
			Name: "Test Room " + string(rune('0'+i)),
		}
		_, err := roomService.CreateRoom(ctx, createReq)
		if err != nil {
			t.Fatalf("Failed to create test room %d: %v", i, err)
		}
	}

	t.Run("ListRooms_Success", func(t *testing.T) {
		req := &service.ListRoomsRequest{
			PageSize: 10,
		}

		resp, err := roomService.ListRooms(ctx, req)
		if err != nil {
			t.Fatalf("ListRooms failed: %v", err)
		}

		if len(resp.Rooms) == 0 {
			t.Error("Expected at least one room")
		}

		if resp.TotalCount < 5 {
			t.Errorf("Expected at least 5 rooms, got %d", resp.TotalCount)
		}
	})

	t.Run("ListRooms_WithPagination", func(t *testing.T) {
		req := &service.ListRoomsRequest{
			PageSize: 2,
		}

		resp, err := roomService.ListRooms(ctx, req)
		if err != nil {
			t.Fatalf("ListRooms failed: %v", err)
		}

		if len(resp.Rooms) > 2 {
			t.Errorf("Expected at most 2 rooms, got %d", len(resp.Rooms))
		}
	})

	t.Run("ListRooms_FilterByState", func(t *testing.T) {
		req := &service.ListRoomsRequest{
			PageSize:    10,
			StateFilter: entity.RoomState_ROOM_STATE_ACTIVE,
		}

		resp, err := roomService.ListRooms(ctx, req)
		if err != nil {
			t.Fatalf("ListRooms failed: %v", err)
		}

		for _, room := range resp.Rooms {
			if room.State != entity.RoomState_ROOM_STATE_ACTIVE {
				t.Errorf("Expected all rooms to be ACTIVE, got %v", room.State)
			}
		}
	})
}
