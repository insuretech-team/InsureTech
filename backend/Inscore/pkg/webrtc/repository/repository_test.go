package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/newage-saint/insuretech/gen/go/insuretech/webrtc/v1/entity"
)

func TestRepository_RoomOperations(t *testing.T) {
	if testDB == nil {
		t.Skip("Test database not available")
	}

	ctx := context.Background()
	repo := NewRepository(testDB)

	// Cleanup before test
	if err := cleanupTestData(ctx); err != nil {
		t.Fatalf("Failed to cleanup test data: %v", err)
	}

	// Cleanup after test
	defer func() {
		if err := cleanupTestData(ctx); err != nil {
			t.Errorf("Failed to cleanup test data: %v", err)
		}
	}()

	t.Run("CreateRoom", func(t *testing.T) {
		room := &entity.Room{
			Name: "Test Room",
			Config: &entity.RoomConfig{
				MaxParticipants:       10,
				RequireToken:          false,
				EnableRecording:       false,
				EnableTranscription:   false,
				SessionTimeoutSeconds: 3600,
			},
			State:            entity.RoomState_ROOM_STATE_ACTIVE,
			MaxParticipants:  10,
			ParticipantCount: 0,
			Metadata: map[string]string{
				"test": "true",
			},
		}

		err := repo.CreateRoom(ctx, room)
		if err != nil {
			t.Fatalf("CreateRoom failed: %v", err)
		}

		if room.RoomId == "" {
			t.Fatal("Room ID not generated")
		}

		if room.CreatedAt == nil {
			t.Fatal("CreatedAt not set")
		}

		t.Logf("✅ Room created with ID: %s", room.RoomId)
	})

	t.Run("GetRoom", func(t *testing.T) {
		// Create a room first
		room := &entity.Room{
			Name:             "Test Room Get",
			Config:           &entity.RoomConfig{MaxParticipants: 10},
			State:            entity.RoomState_ROOM_STATE_ACTIVE,
			MaxParticipants:  10,
			ParticipantCount: 0,
		}

		if err := repo.CreateRoom(ctx, room); err != nil {
			t.Fatalf("Failed to create room: %v", err)
		}

		// Get the room
		retrieved, err := repo.GetRoom(ctx, room.RoomId)
		if err != nil {
			t.Fatalf("GetRoom failed: %v", err)
		}

		if retrieved == nil {
			t.Fatal("Room not found")
		}

		if retrieved.RoomId != room.RoomId {
			t.Errorf("Room ID mismatch: got %s, want %s", retrieved.RoomId, room.RoomId)
		}

		if retrieved.Name != room.Name {
			t.Errorf("Room name mismatch: got %s, want %s", retrieved.Name, room.Name)
		}

		if retrieved.State != room.State {
			t.Errorf("Room state mismatch: got %v, want %v", retrieved.State, room.State)
		}

		t.Logf("✅ Room retrieved successfully: %s", retrieved.Name)
	})

	t.Run("UpdateRoom", func(t *testing.T) {
		// Create a room first
		room := &entity.Room{
			Name:             "Test Room Update",
			Config:           &entity.RoomConfig{MaxParticipants: 10},
			State:            entity.RoomState_ROOM_STATE_ACTIVE,
			MaxParticipants:  10,
			ParticipantCount: 0,
		}

		if err := repo.CreateRoom(ctx, room); err != nil {
			t.Fatalf("Failed to create room: %v", err)
		}

		// Update the room
		room.Name = "Updated Room Name"
		room.ParticipantCount = 5
		room.Metadata = map[string]string{
			"updated": "true",
		}

		if err := repo.UpdateRoom(ctx, room); err != nil {
			t.Fatalf("UpdateRoom failed: %v", err)
		}

		// Verify update
		retrieved, err := repo.GetRoom(ctx, room.RoomId)
		if err != nil {
			t.Fatalf("Failed to get updated room: %v", err)
		}

		if retrieved.Name != "Updated Room Name" {
			t.Errorf("Name not updated: got %s", retrieved.Name)
		}

		if retrieved.ParticipantCount != 5 {
			t.Errorf("ParticipantCount not updated: got %d", retrieved.ParticipantCount)
		}

		t.Logf("✅ Room updated successfully")
	})

	t.Run("CloseRoom", func(t *testing.T) {
		// Create a room first
		room := &entity.Room{
			Name:             "Test Room Close",
			Config:           &entity.RoomConfig{MaxParticipants: 10},
			State:            entity.RoomState_ROOM_STATE_ACTIVE,
			MaxParticipants:  10,
			ParticipantCount: 0,
		}

		if err := repo.CreateRoom(ctx, room); err != nil {
			t.Fatalf("Failed to create room: %v", err)
		}

		// Close the room
		if err := repo.CloseRoom(ctx, room.RoomId); err != nil {
			t.Fatalf("CloseRoom failed: %v", err)
		}

		// Verify closed
		retrieved, err := repo.GetRoom(ctx, room.RoomId)
		if err != nil {
			t.Fatalf("Failed to get closed room: %v", err)
		}

		if retrieved.State != entity.RoomState_ROOM_STATE_CLOSED {
			t.Errorf("Room not closed: state is %v", retrieved.State)
		}

		if retrieved.ClosedAt == nil {
			t.Error("ClosedAt not set")
		}

		t.Logf("✅ Room closed successfully")
	})

	t.Run("ListRooms", func(t *testing.T) {
		// Create multiple rooms
		for i := 0; i < 5; i++ {
			room := &entity.Room{
				Name:             fmt.Sprintf("Test Room %d", i),
				Config:           &entity.RoomConfig{MaxParticipants: 10},
				State:            entity.RoomState_ROOM_STATE_ACTIVE,
				MaxParticipants:  10,
				ParticipantCount: 0,
			}
			if err := repo.CreateRoom(ctx, room); err != nil {
				t.Fatalf("Failed to create room %d: %v", i, err)
			}
		}

		// List all rooms
		rooms, total, err := repo.ListRooms(ctx, 10, 0, entity.RoomState_ROOM_STATE_UNSPECIFIED)
		if err != nil {
			t.Fatalf("ListRooms failed: %v", err)
		}

		if len(rooms) == 0 {
			t.Fatal("No rooms returned")
		}

		if total == 0 {
			t.Fatal("Total count is 0")
		}

		t.Logf("✅ Listed %d rooms (total: %d)", len(rooms), total)

		// Test filtering by state
		activeRooms, activeTotal, err := repo.ListRooms(ctx, 10, 0, entity.RoomState_ROOM_STATE_ACTIVE)
		if err != nil {
			t.Fatalf("ListRooms with filter failed: %v", err)
		}

		if len(activeRooms) == 0 {
			t.Error("No active rooms returned")
		}

		t.Logf("✅ Listed %d active rooms (total: %d)", len(activeRooms), activeTotal)
	})
}

func TestRepository_PeerOperations(t *testing.T) {
	if testDB == nil {
		t.Skip("Test database not available")
	}

	ctx := context.Background()
	repo := NewRepository(testDB)

	// Cleanup before test
	if err := cleanupTestData(ctx); err != nil {
		t.Fatalf("Failed to cleanup test data: %v", err)
	}

	// Cleanup after test
	defer func() {
		if err := cleanupTestData(ctx); err != nil {
			t.Errorf("Failed to cleanup test data: %v", err)
		}
	}()

	// Create a room for peer tests
	room := &entity.Room{
		Name:             "Test Room for Peers",
		Config:           &entity.RoomConfig{MaxParticipants: 10},
		State:            entity.RoomState_ROOM_STATE_ACTIVE,
		MaxParticipants:  10,
		ParticipantCount: 0,
	}

	if err := repo.CreateRoom(ctx, room); err != nil {
		t.Fatalf("Failed to create room: %v", err)
	}

	t.Run("AddPeer", func(t *testing.T) {
		peer := &entity.Peer{
			RoomId:      room.RoomId,
			DisplayName: "Test Peer",
			State:       entity.PeerConnectionState_PEER_CONNECTION_STATE_NEW,
			UserAgent:   "Test Agent",
			Metadata: map[string]string{
				"test": "true",
			},
		}

		err := repo.AddPeer(ctx, peer)
		if err != nil {
			t.Fatalf("AddPeer failed: %v", err)
		}

		if peer.PeerId == "" {
			t.Fatal("Peer ID not generated")
		}

		if peer.JoinedAt == nil {
			t.Fatal("JoinedAt not set")
		}

		if peer.LastSeenAt == nil {
			t.Fatal("LastSeenAt not set")
		}

		t.Logf("✅ Peer created with ID: %s", peer.PeerId)
	})

	t.Run("GetPeer", func(t *testing.T) {
		// Create a peer first
		peer := &entity.Peer{
			RoomId:      room.RoomId,
			DisplayName: "Test Peer Get",
			State:       entity.PeerConnectionState_PEER_CONNECTION_STATE_NEW,
		}

		if err := repo.AddPeer(ctx, peer); err != nil {
			t.Fatalf("Failed to create peer: %v", err)
		}

		// Get the peer
		retrieved, err := repo.GetPeer(ctx, peer.PeerId)
		if err != nil {
			t.Fatalf("GetPeer failed: %v", err)
		}

		if retrieved == nil {
			t.Fatal("Peer not found")
		}

		if retrieved.PeerId != peer.PeerId {
			t.Errorf("Peer ID mismatch: got %s, want %s", retrieved.PeerId, peer.PeerId)
		}

		if retrieved.DisplayName != peer.DisplayName {
			t.Errorf("DisplayName mismatch: got %s, want %s", retrieved.DisplayName, peer.DisplayName)
		}

		t.Logf("✅ Peer retrieved successfully: %s", retrieved.DisplayName)
	})

	t.Run("UpdatePeerState", func(t *testing.T) {
		// Create a peer first
		peer := &entity.Peer{
			RoomId:      room.RoomId,
			DisplayName: "Test Peer State",
			State:       entity.PeerConnectionState_PEER_CONNECTION_STATE_NEW,
		}

		if err := repo.AddPeer(ctx, peer); err != nil {
			t.Fatalf("Failed to create peer: %v", err)
		}

		// Update state
		newState := entity.PeerConnectionState_PEER_CONNECTION_STATE_CONNECTED
		if err := repo.UpdatePeerState(ctx, peer.PeerId, newState); err != nil {
			t.Fatalf("UpdatePeerState failed: %v", err)
		}

		// Verify update
		retrieved, err := repo.GetPeer(ctx, peer.PeerId)
		if err != nil {
			t.Fatalf("Failed to get updated peer: %v", err)
		}

		if retrieved.State != newState {
			t.Errorf("State not updated: got %v, want %v", retrieved.State, newState)
		}

		t.Logf("✅ Peer state updated successfully")
	})

	t.Run("RemovePeer", func(t *testing.T) {
		// Create a peer first
		peer := &entity.Peer{
			RoomId:      room.RoomId,
			DisplayName: "Test Peer Remove",
			State:       entity.PeerConnectionState_PEER_CONNECTION_STATE_NEW,
		}

		if err := repo.AddPeer(ctx, peer); err != nil {
			t.Fatalf("Failed to create peer: %v", err)
		}

		// Remove the peer
		if err := repo.RemovePeer(ctx, peer.PeerId); err != nil {
			t.Fatalf("RemovePeer failed: %v", err)
		}

		// Verify removed
		retrieved, err := repo.GetPeer(ctx, peer.PeerId)
		if err != nil {
			t.Fatalf("Failed to get removed peer: %v", err)
		}

		if retrieved.State != entity.PeerConnectionState_PEER_CONNECTION_STATE_CLOSED {
			t.Errorf("Peer not marked as closed: state is %v", retrieved.State)
		}

		if retrieved.LeftAt == nil {
			t.Error("LeftAt not set")
		}

		t.Logf("✅ Peer removed successfully")
	})

	t.Run("ListPeersInRoom", func(t *testing.T) {
		// Create multiple peers
		for i := 0; i < 3; i++ {
			peer := &entity.Peer{
				RoomId:      room.RoomId,
				DisplayName: fmt.Sprintf("Test Peer %d", i),
				State:       entity.PeerConnectionState_PEER_CONNECTION_STATE_CONNECTED,
			}
			if err := repo.AddPeer(ctx, peer); err != nil {
				t.Fatalf("Failed to create peer %d: %v", i, err)
			}
		}

		// List peers
		peers, err := repo.ListPeersInRoom(ctx, room.RoomId)
		if err != nil {
			t.Fatalf("ListPeersInRoom failed: %v", err)
		}

		if len(peers) == 0 {
			t.Fatal("No peers returned")
		}

		t.Logf("✅ Listed %d peers in room", len(peers))
	})
}

func TestRepository_TrackOperations(t *testing.T) {
	if testDB == nil {
		t.Skip("Test database not available")
	}

	ctx := context.Background()
	repo := NewRepository(testDB)

	// Cleanup before test
	if err := cleanupTestData(ctx); err != nil {
		t.Fatalf("Failed to cleanup test data: %v", err)
	}

	// Cleanup after test
	defer func() {
		if err := cleanupTestData(ctx); err != nil {
			t.Errorf("Failed to cleanup test data: %v", err)
		}
	}()

	// Create a room and peer for track tests
	room := &entity.Room{
		Name:             "Test Room for Tracks",
		Config:           &entity.RoomConfig{MaxParticipants: 10},
		State:            entity.RoomState_ROOM_STATE_ACTIVE,
		MaxParticipants:  10,
		ParticipantCount: 0,
	}

	if err := repo.CreateRoom(ctx, room); err != nil {
		t.Fatalf("Failed to create room: %v", err)
	}

	peer := &entity.Peer{
		RoomId:      room.RoomId,
		DisplayName: "Test Peer",
		State:       entity.PeerConnectionState_PEER_CONNECTION_STATE_CONNECTED,
	}

	if err := repo.AddPeer(ctx, peer); err != nil {
		t.Fatalf("Failed to create peer: %v", err)
	}

	t.Run("PublishTrack", func(t *testing.T) {
		track := &entity.Track{
			TrackId: "track-test-001",
			PeerId:  peer.PeerId,
			Type:    entity.TrackType_TRACK_TYPE_VIDEO,
			Label:   "Test Camera",
			Muted:   false,
			State:   entity.TrackState_TRACK_STATE_ACTIVE,
			Settings: &entity.TrackSettings{
				Width:     1920,
				Height:    1080,
				FrameRate: 30.0,
				Bitrate:   2500,
				Codec:     "VP8",
			},
			Metadata: map[string]string{
				"device": "camera",
			},
		}

		err := repo.PublishTrack(ctx, track)
		if err != nil {
			t.Fatalf("PublishTrack failed: %v", err)
		}

		t.Logf("✅ Track published: %s", track.TrackId)
	})

	t.Run("GetTrack", func(t *testing.T) {
		// Publish a track first
		track := &entity.Track{
			TrackId: "track-test-002",
			PeerId:  peer.PeerId,
			Type:    entity.TrackType_TRACK_TYPE_AUDIO,
			Label:   "Test Microphone",
			Muted:   false,
			State:   entity.TrackState_TRACK_STATE_ACTIVE,
		}

		if err := repo.PublishTrack(ctx, track); err != nil {
			t.Fatalf("Failed to publish track: %v", err)
		}

		// Get the track
		retrieved, err := repo.GetTrack(ctx, track.TrackId)
		if err != nil {
			t.Fatalf("GetTrack failed: %v", err)
		}

		if retrieved == nil {
			t.Fatal("Track not found")
		}

		if retrieved.TrackId != track.TrackId {
			t.Errorf("Track ID mismatch: got %s, want %s", retrieved.TrackId, track.TrackId)
		}

		if retrieved.Type != track.Type {
			t.Errorf("Type mismatch: got %v, want %v", retrieved.Type, track.Type)
		}

		t.Logf("✅ Track retrieved successfully: %s", retrieved.Label)
	})

	t.Run("MuteTrack", func(t *testing.T) {
		// Publish a track first
		track := &entity.Track{
			TrackId: "track-test-003",
			PeerId:  peer.PeerId,
			Type:    entity.TrackType_TRACK_TYPE_VIDEO,
			Label:   "Test Video",
			Muted:   false,
			State:   entity.TrackState_TRACK_STATE_ACTIVE,
		}

		if err := repo.PublishTrack(ctx, track); err != nil {
			t.Fatalf("Failed to publish track: %v", err)
		}

		// Mute the track
		if err := repo.MuteTrack(ctx, track.TrackId, true); err != nil {
			t.Fatalf("MuteTrack failed: %v", err)
		}

		// Verify muted
		retrieved, err := repo.GetTrack(ctx, track.TrackId)
		if err != nil {
			t.Fatalf("Failed to get track: %v", err)
		}

		if !retrieved.Muted {
			t.Error("Track not muted")
		}

		t.Logf("✅ Track muted successfully")
	})

	t.Run("UpdateTrack", func(t *testing.T) {
		// Publish a track first
		track := &entity.Track{
			TrackId: "track-test-004",
			PeerId:  peer.PeerId,
			Type:    entity.TrackType_TRACK_TYPE_VIDEO,
			Label:   "Test Video Update",
			Muted:   false,
			State:   entity.TrackState_TRACK_STATE_ACTIVE,
			Settings: &entity.TrackSettings{
				Width:  640,
				Height: 480,
			},
		}

		if err := repo.PublishTrack(ctx, track); err != nil {
			t.Fatalf("Failed to publish track: %v", err)
		}

		// Update settings
		track.Settings.Width = 1920
		track.Settings.Height = 1080
		track.Metadata = map[string]string{
			"updated": "true",
		}

		if err := repo.UpdateTrack(ctx, track); err != nil {
			t.Fatalf("UpdateTrack failed: %v", err)
		}

		// Verify update
		retrieved, err := repo.GetTrack(ctx, track.TrackId)
		if err != nil {
			t.Fatalf("Failed to get track: %v", err)
		}

		if retrieved.Settings.Width != 1920 {
			t.Errorf("Width not updated: got %d", retrieved.Settings.Width)
		}

		t.Logf("✅ Track updated successfully")
	})

	t.Run("UnpublishTrack", func(t *testing.T) {
		// Publish a track first
		track := &entity.Track{
			TrackId: "track-test-005",
			PeerId:  peer.PeerId,
			Type:    entity.TrackType_TRACK_TYPE_VIDEO,
			Label:   "Test Video Delete",
			Muted:   false,
			State:   entity.TrackState_TRACK_STATE_ACTIVE,
		}

		if err := repo.PublishTrack(ctx, track); err != nil {
			t.Fatalf("Failed to publish track: %v", err)
		}

		// Unpublish the track
		if err := repo.UnpublishTrack(ctx, track.TrackId); err != nil {
			t.Fatalf("UnpublishTrack failed: %v", err)
		}

		// Verify deleted
		retrieved, err := repo.GetTrack(ctx, track.TrackId)
		if err != nil {
			t.Fatalf("GetTrack failed: %v", err)
		}

		if retrieved != nil {
			t.Error("Track still exists after unpublish")
		}

		t.Logf("✅ Track unpublished successfully")
	})

	t.Run("ListTracks", func(t *testing.T) {
		// Publish multiple tracks
		for i := 0; i < 3; i++ {
			track := &entity.Track{
				TrackId: fmt.Sprintf("track-list-%d", i),
				PeerId:  peer.PeerId,
				Type:    entity.TrackType_TRACK_TYPE_VIDEO,
				Label:   fmt.Sprintf("Test Track %d", i),
				Muted:   false,
				State:   entity.TrackState_TRACK_STATE_ACTIVE,
			}
			if err := repo.PublishTrack(ctx, track); err != nil {
				t.Fatalf("Failed to publish track %d: %v", i, err)
			}
		}

		// List tracks
		tracks, err := repo.ListTracks(ctx, room.RoomId, "", entity.TrackType_TRACK_TYPE_UNSPECIFIED)
		if err != nil {
			t.Fatalf("ListTracks failed: %v", err)
		}

		if len(tracks) == 0 {
			t.Fatal("No tracks returned")
		}

		t.Logf("✅ Listed %d tracks", len(tracks))
	})
}
