package services

import (
	"context"
	"testing"

	"github.com/newage-saint/insuretech/backend/inscore/pkg/webrtc/repository"

	"github.com/newage-saint/insuretech/gen/go/insuretech/webrtc/v1/entity"
	"github.com/newage-saint/insuretech/gen/go/insuretech/webrtc/v1/service"
)

func TestTrackService_PublishTrack(t *testing.T) {
	ctx := context.Background()
	defer cleanupTestData(ctx)

	repo := repository.NewRepository(getTestDB())
	roomService := NewRoomService(repo)
	peerService := NewPeerService(repo)
	trackService := NewTrackService(repo)

	// Create test room and peer
	roomResp, _ := roomService.CreateRoom(ctx, &service.CreateRoomRequest{Name: "Test Room"})
	joinResp, err := peerService.JoinRoom(ctx, &service.JoinRoomRequest{
		RoomId:      roomResp.Room.RoomId,
		DisplayName: "Test User",
	})
	if err != nil {
		t.Fatalf("Failed to join room: %v", err)
	}

	t.Run("PublishTrack_Audio_Success", func(t *testing.T) {
		req := &service.PublishTrackRequest{
			RoomId: roomResp.Room.RoomId,
			PeerId: joinResp.Peer.PeerId,
			Track: &entity.Track{
				TrackId: "audio-track-1",
				Type:    entity.TrackType_TRACK_TYPE_AUDIO,
				Label:   "Audio Track",
				Muted:   false,
				State:   entity.TrackState_TRACK_STATE_ACTIVE,
			},
		}

		resp, err := trackService.PublishTrack(ctx, req)
		if err != nil {
			t.Fatalf("PublishTrack failed: %v", err)
		}

		if resp.Track == nil {
			t.Fatal("Track is nil")
		}

		if resp.Track.Type != entity.TrackType_TRACK_TYPE_AUDIO {
			t.Errorf("Expected track type AUDIO, got %v", resp.Track.Type)
		}

		if resp.Track.TrackId != "audio-track-1" {
			t.Errorf("Expected track ID 'audio-track-1', got '%s'", resp.Track.TrackId)
		}
	})

	t.Run("PublishTrack_Video_Success", func(t *testing.T) {
		req := &service.PublishTrackRequest{
			RoomId: roomResp.Room.RoomId,
			PeerId: joinResp.Peer.PeerId,
			Track: &entity.Track{
				TrackId: "video-track-1",
				Type:    entity.TrackType_TRACK_TYPE_VIDEO,
				Label:   "Video Track",
				Muted:   false,
				State:   entity.TrackState_TRACK_STATE_ACTIVE,
			},
		}

		resp, err := trackService.PublishTrack(ctx, req)
		if err != nil {
			t.Fatalf("PublishTrack failed: %v", err)
		}

		if resp.Track.Type != entity.TrackType_TRACK_TYPE_VIDEO {
			t.Errorf("Expected track type VIDEO, got %v", resp.Track.Type)
		}
	})

	t.Run("PublishTrack_MissingRoomId", func(t *testing.T) {
		req := &service.PublishTrackRequest{
			PeerId: joinResp.Peer.PeerId,
			Track: &entity.Track{
				TrackId: "track-1",
				Type:    entity.TrackType_TRACK_TYPE_AUDIO,
			},
		}

		_, err := trackService.PublishTrack(ctx, req)
		if err == nil {
			t.Fatal("Expected error for missing room ID")
		}
	})

	t.Run("PublishTrack_MissingPeerId", func(t *testing.T) {
		req := &service.PublishTrackRequest{
			RoomId: roomResp.Room.RoomId,
			Track: &entity.Track{
				TrackId: "track-1",
				Type:    entity.TrackType_TRACK_TYPE_AUDIO,
			},
		}

		_, err := trackService.PublishTrack(ctx, req)
		if err == nil {
			t.Fatal("Expected error for missing peer ID")
		}
	})

	t.Run("PublishTrack_MissingTrack", func(t *testing.T) {
		req := &service.PublishTrackRequest{
			RoomId: roomResp.Room.RoomId,
			PeerId: joinResp.Peer.PeerId,
		}

		_, err := trackService.PublishTrack(ctx, req)
		if err == nil {
			t.Fatal("Expected error for missing track")
		}
	})

	t.Run("PublishTrack_PeerNotFound", func(t *testing.T) {
		req := &service.PublishTrackRequest{
			RoomId: roomResp.Room.RoomId,
			PeerId: "00000000-0000-0000-0000-000000000000",
			Track: &entity.Track{
				TrackId: "track-1",
				Type:    entity.TrackType_TRACK_TYPE_AUDIO,
			},
		}

		_, err := trackService.PublishTrack(ctx, req)
		if err == nil {
			t.Fatal("Expected error for non-existent peer")
		}
	})
}

func TestTrackService_GetTrack(t *testing.T) {
	ctx := context.Background()
	defer cleanupTestData(ctx)

	repo := repository.NewRepository(getTestDB())
	roomService := NewRoomService(repo)
	peerService := NewPeerService(repo)
	trackService := NewTrackService(repo)

	// Create test room, peer, and track
	roomResp, _ := roomService.CreateRoom(ctx, &service.CreateRoomRequest{Name: "Test Room"})
	joinResp, _ := peerService.JoinRoom(ctx, &service.JoinRoomRequest{
		RoomId:      roomResp.Room.RoomId,
		DisplayName: "Test User",
	})
	publishResp, err := trackService.PublishTrack(ctx, &service.PublishTrackRequest{
		RoomId: roomResp.Room.RoomId,
		PeerId: joinResp.Peer.PeerId,
		Track: &entity.Track{
			TrackId: "test-track-1",
			Type:    entity.TrackType_TRACK_TYPE_AUDIO,
			Label:   "Test Track",
		},
	})
	if err != nil {
		t.Fatalf("Failed to publish track: %v", err)
	}

	t.Run("GetTrack_Success", func(t *testing.T) {
		req := &service.GetTrackRequest{
			TrackId: publishResp.Track.TrackId,
		}

		resp, err := trackService.GetTrack(ctx, req)
		if err != nil {
			t.Fatalf("GetTrack failed: %v", err)
		}

		if resp.Track.TrackId != publishResp.Track.TrackId {
			t.Errorf("Expected track ID %s, got %s", publishResp.Track.TrackId, resp.Track.TrackId)
		}

		if resp.Track.Label != "Test Track" {
			t.Errorf("Expected label 'Test Track', got '%s'", resp.Track.Label)
		}
	})

	t.Run("GetTrack_NotFound", func(t *testing.T) {
		req := &service.GetTrackRequest{
			TrackId: "non-existent-track",
		}

		_, err := trackService.GetTrack(ctx, req)
		if err == nil {
			t.Fatal("Expected error for non-existent track")
		}
	})
}

func TestTrackService_UnpublishTrack(t *testing.T) {
	ctx := context.Background()
	defer cleanupTestData(ctx)

	repo := repository.NewRepository(getTestDB())
	roomService := NewRoomService(repo)
	peerService := NewPeerService(repo)
	trackService := NewTrackService(repo)

	// Create test room, peer, and track
	roomResp, _ := roomService.CreateRoom(ctx, &service.CreateRoomRequest{Name: "Test Room"})
	joinResp, _ := peerService.JoinRoom(ctx, &service.JoinRoomRequest{
		RoomId:      roomResp.Room.RoomId,
		DisplayName: "Test User",
	})
	publishResp, _ := trackService.PublishTrack(ctx, &service.PublishTrackRequest{
		RoomId: roomResp.Room.RoomId,
		PeerId: joinResp.Peer.PeerId,
		Track: &entity.Track{
			TrackId: "test-track-1",
			Type:    entity.TrackType_TRACK_TYPE_AUDIO,
		},
	})

	t.Run("UnpublishTrack_Success", func(t *testing.T) {
		req := &service.UnpublishTrackRequest{
			RoomId:  roomResp.Room.RoomId,
			PeerId:  joinResp.Peer.PeerId,
			TrackId: publishResp.Track.TrackId,
		}

		resp, err := trackService.UnpublishTrack(ctx, req)
		if err != nil {
			t.Fatalf("UnpublishTrack failed: %v", err)
		}

		if !resp.Success {
			t.Error("Expected success to be true")
		}

		// Verify track no longer exists
		getReq := &service.GetTrackRequest{
			TrackId: publishResp.Track.TrackId,
		}
		_, err = trackService.GetTrack(ctx, getReq)
		if err == nil {
			t.Fatal("Expected error when getting unpublished track")
		}
	})
}

func TestTrackService_MuteTrack(t *testing.T) {
	ctx := context.Background()
	defer cleanupTestData(ctx)

	repo := repository.NewRepository(getTestDB())
	roomService := NewRoomService(repo)
	peerService := NewPeerService(repo)
	trackService := NewTrackService(repo)

	// Create test room, peer, and track
	roomResp, _ := roomService.CreateRoom(ctx, &service.CreateRoomRequest{Name: "Test Room"})
	joinResp, _ := peerService.JoinRoom(ctx, &service.JoinRoomRequest{
		RoomId:      roomResp.Room.RoomId,
		DisplayName: "Test User",
	})
	publishResp, _ := trackService.PublishTrack(ctx, &service.PublishTrackRequest{
		RoomId: roomResp.Room.RoomId,
		PeerId: joinResp.Peer.PeerId,
		Track: &entity.Track{
			TrackId: "test-track-1",
			Type:    entity.TrackType_TRACK_TYPE_AUDIO,
			Muted:   false,
		},
	})

	t.Run("MuteTrack_Success", func(t *testing.T) {
		req := &service.MuteTrackRequest{
			RoomId:  roomResp.Room.RoomId,
			PeerId:  joinResp.Peer.PeerId,
			TrackId: publishResp.Track.TrackId,
			Muted:   true,
		}

		resp, err := trackService.MuteTrack(ctx, req)
		if err != nil {
			t.Fatalf("MuteTrack failed: %v", err)
		}

		if resp.Track == nil {
			t.Fatal("Expected track in response")
		}

		// Verify track is muted
		getResp, err := trackService.GetTrack(ctx, &service.GetTrackRequest{
			TrackId: publishResp.Track.TrackId,
		})
		if err != nil {
			t.Fatalf("GetTrack failed: %v", err)
		}

		if !getResp.Track.Muted {
			t.Error("Expected track to be muted")
		}
	})

	t.Run("UnmuteTrack_Success", func(t *testing.T) {
		req := &service.MuteTrackRequest{
			RoomId:  roomResp.Room.RoomId,
			PeerId:  joinResp.Peer.PeerId,
			TrackId: publishResp.Track.TrackId,
			Muted:   false,
		}

		resp, err := trackService.MuteTrack(ctx, req)
		if err != nil {
			t.Fatalf("MuteTrack failed: %v", err)
		}

		if resp.Track == nil {
			t.Fatal("Expected track in response")
		}

		// Verify track is unmuted
		getResp, err := trackService.GetTrack(ctx, &service.GetTrackRequest{
			TrackId: publishResp.Track.TrackId,
		})
		if err != nil {
			t.Fatalf("GetTrack failed: %v", err)
		}

		if getResp.Track.Muted {
			t.Error("Expected track to be unmuted")
		}
	})
}

func TestTrackService_ListTracks(t *testing.T) {
	ctx := context.Background()
	defer cleanupTestData(ctx)

	repo := repository.NewRepository(getTestDB())
	roomService := NewRoomService(repo)
	peerService := NewPeerService(repo)
	trackService := NewTrackService(repo)

	// Create test room and peer
	roomResp, _ := roomService.CreateRoom(ctx, &service.CreateRoomRequest{Name: "Test Room"})
	joinResp, _ := peerService.JoinRoom(ctx, &service.JoinRoomRequest{
		RoomId:      roomResp.Room.RoomId,
		DisplayName: "Test User",
	})

	// Publish multiple tracks
	for i := 1; i <= 3; i++ {
		trackType := entity.TrackType_TRACK_TYPE_AUDIO
		if i == 3 {
			trackType = entity.TrackType_TRACK_TYPE_VIDEO
		}

		_, err := trackService.PublishTrack(ctx, &service.PublishTrackRequest{
			RoomId: roomResp.Room.RoomId,
			PeerId: joinResp.Peer.PeerId,
			Track: &entity.Track{
				TrackId: "track-" + string(rune('0'+i)),
				Type:    trackType,
				Label:   "Track " + string(rune('0'+i)),
			},
		})
		if err != nil {
			t.Fatalf("Failed to publish track %d: %v", i, err)
		}
	}

	t.Run("ListTracks_AllTracks", func(t *testing.T) {
		req := &service.ListTracksRequest{
			RoomId: roomResp.Room.RoomId,
			PeerId: joinResp.Peer.PeerId,
		}

		resp, err := trackService.ListTracks(ctx, req)
		if err != nil {
			t.Fatalf("ListTracks failed: %v", err)
		}

		if len(resp.Tracks) != 3 {
			t.Errorf("Expected 3 tracks, got %d", len(resp.Tracks))
		}
	})

	t.Run("ListTracks_FilterByType", func(t *testing.T) {
		req := &service.ListTracksRequest{
			RoomId:     roomResp.Room.RoomId,
			PeerId:     joinResp.Peer.PeerId,
			TypeFilter: entity.TrackType_TRACK_TYPE_AUDIO,
		}

		resp, err := trackService.ListTracks(ctx, req)
		if err != nil {
			t.Fatalf("ListTracks failed: %v", err)
		}

		if len(resp.Tracks) != 2 {
			t.Errorf("Expected 2 audio tracks, got %d", len(resp.Tracks))
		}

		for _, track := range resp.Tracks {
			if track.Type != entity.TrackType_TRACK_TYPE_AUDIO {
				t.Errorf("Expected all tracks to be AUDIO, got %v", track.Type)
			}
		}
	})

	t.Run("ListTracks_FilterByVideo", func(t *testing.T) {
		req := &service.ListTracksRequest{
			RoomId:     roomResp.Room.RoomId,
			PeerId:     joinResp.Peer.PeerId,
			TypeFilter: entity.TrackType_TRACK_TYPE_VIDEO,
		}

		resp, err := trackService.ListTracks(ctx, req)
		if err != nil {
			t.Fatalf("ListTracks failed: %v", err)
		}

		if len(resp.Tracks) != 1 {
			t.Errorf("Expected 1 video track, got %d", len(resp.Tracks))
		}

		if resp.Tracks[0].Type != entity.TrackType_TRACK_TYPE_VIDEO {
			t.Errorf("Expected track type VIDEO, got %v", resp.Tracks[0].Type)
		}
	})
}
