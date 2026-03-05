package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/newage-saint/insuretech/gen/go/insuretech/webrtc/v1/entity"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

// Repository provides data access for WebRTC entities
// Tables are created and managed by proto-driven migrations
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new WebRTC repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// GetDB returns the underlying database connection
func (r *Repository) GetDB() *gorm.DB {
	return r.db
}

// ═══════════════════════════════════════════════════════════════════════════
// Room Operations
// ═══════════════════════════════════════════════════════════════════════════

// CreateRoom creates a new room in webrtc.rooms table
func (r *Repository) CreateRoom(ctx context.Context, room *entity.Room) error {
	configJSON, err := json.Marshal(room.Config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	query := `
		INSERT INTO webrtc.rooms (
			room_id, name, config, participant_count, max_participants, 
			created_at, state, closed_at, creator_id
		) VALUES (
			gen_random_uuid(), $1, $2, $3, $4, 
			CURRENT_TIMESTAMP, $5, $6, $7
		) RETURNING room_id, created_at`

	var roomID string
	var createdAt sql.NullTime

	// Convert state enum to integer for database
	stateValue := int32(room.State)

	// Handle optional creator_id - pass nil if not provided or empty
	var creatorID interface{}
	if room.CreatorId != "" {
		creatorID = room.CreatorId
	} else {
		creatorID = nil
	}

	err = r.db.WithContext(ctx).Raw(query,
		room.Name,
		configJSON,
		room.ParticipantCount,
		room.MaxParticipants,
		stateValue,
		nil, // closed_at
		creatorID,
	).Row().Scan(&roomID, &createdAt)

	if err != nil {
		return fmt.Errorf("failed to create room: %w", err)
	}

	room.RoomId = roomID
	if createdAt.Valid {
		room.CreatedAt = timestamppb.New(createdAt.Time)
	}

	return nil
}

// GetRoom retrieves a room by ID from webrtc.rooms table
func (r *Repository) GetRoom(ctx context.Context, roomID string) (*entity.Room, error) {
	query := `
		SELECT room_id, name, config, participant_count, max_participants,
		       created_at, metadata, state, closed_at, creator_id
		FROM webrtc.rooms
		WHERE room_id = $1`

	var (
		id, name, creatorID               sql.NullString
		state                             sql.NullInt32
		configJSON, metadataJSON          []byte
		participantCount, maxParticipants sql.NullInt32
		createdAt, closedAt               sql.NullTime
	)

	err := r.db.WithContext(ctx).Raw(query, roomID).Row().Scan(
		&id, &name, &configJSON, &participantCount, &maxParticipants,
		&createdAt, &metadataJSON, &state, &closedAt, &creatorID,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get room: %w", err)
	}

	room := &entity.Room{
		RoomId:           id.String,
		Name:             name.String,
		ParticipantCount: participantCount.Int32,
		MaxParticipants:  maxParticipants.Int32,
		State:            entity.RoomState(state.Int32),
	}

	if len(configJSON) > 0 {
		var config entity.RoomConfig
		if err := json.Unmarshal(configJSON, &config); err == nil {
			room.Config = &config
		}
	}

	if len(metadataJSON) > 0 {
		var metadata map[string]string
		if err := json.Unmarshal(metadataJSON, &metadata); err == nil {
			room.Metadata = metadata
		}
	}

	if createdAt.Valid {
		room.CreatedAt = timestamppb.New(createdAt.Time)
	}

	if closedAt.Valid {
		room.ClosedAt = timestamppb.New(closedAt.Time)
	}

	if creatorID.Valid {
		room.CreatorId = creatorID.String
	}

	return room, nil
}

// UpdateRoom updates an existing room in webrtc.rooms table
func (r *Repository) UpdateRoom(ctx context.Context, room *entity.Room) error {
	configJSON, err := json.Marshal(room.Config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	metadataJSON, err := json.Marshal(room.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	query := `
		UPDATE webrtc.rooms
		SET name = $1, config = $2, participant_count = $3, 
		    max_participants = $4, metadata = $5, state = $6
		WHERE room_id = $7`

	result := r.db.WithContext(ctx).Exec(query,
		room.Name,
		configJSON,
		room.ParticipantCount,
		room.MaxParticipants,
		metadataJSON,
		int32(room.State),
		room.RoomId,
	)

	if result.Error != nil {
		return fmt.Errorf("failed to update room: %w", result.Error)
	}

	return nil
}

// CloseRoom marks a room as closed in webrtc.rooms table
func (r *Repository) CloseRoom(ctx context.Context, roomID string) error {
	query := `
		UPDATE webrtc.rooms
		SET state = $1, closed_at = CURRENT_TIMESTAMP
		WHERE room_id = $2`

	result := r.db.WithContext(ctx).Exec(query, int32(entity.RoomState_ROOM_STATE_CLOSED), roomID)
	if result.Error != nil {
		return fmt.Errorf("failed to close room: %w", result.Error)
	}

	return nil
}

// ListRooms retrieves rooms with pagination from webrtc.rooms table
func (r *Repository) ListRooms(ctx context.Context, limit, offset int, stateFilter entity.RoomState) ([]*entity.Room, int64, error) {
	// Count total
	var total int64
	countQuery := "SELECT COUNT(*) FROM webrtc.rooms"
	if stateFilter != entity.RoomState_ROOM_STATE_UNSPECIFIED {
		countQuery += fmt.Sprintf(" WHERE state = %d", int32(stateFilter))
	}

	if err := r.db.WithContext(ctx).Raw(countQuery).Scan(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count rooms: %w", err)
	}

	// Fetch rooms
	query := `
		SELECT room_id, name, config, participant_count, max_participants,
		       created_at, metadata, state, closed_at, creator_id
		FROM webrtc.rooms`

	if stateFilter != entity.RoomState_ROOM_STATE_UNSPECIFIED {
		query += fmt.Sprintf(" WHERE state = %d", int32(stateFilter))
	}

	query += " ORDER BY created_at DESC LIMIT $1 OFFSET $2"

	rows, err := r.db.WithContext(ctx).Raw(query, limit, offset).Rows()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list rooms: %w", err)
	}
	defer rows.Close()

	var rooms []*entity.Room
	for rows.Next() {
		var (
			id, name, creatorID               sql.NullString
			state                             sql.NullInt32
			configJSON, metadataJSON          []byte
			participantCount, maxParticipants sql.NullInt32
			createdAt, closedAt               sql.NullTime
		)

		if err := rows.Scan(
			&id, &name, &configJSON, &participantCount, &maxParticipants,
			&createdAt, &metadataJSON, &state, &closedAt, &creatorID,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan room: %w", err)
		}

		room := &entity.Room{
			RoomId:           id.String,
			Name:             name.String,
			ParticipantCount: participantCount.Int32,
			MaxParticipants:  maxParticipants.Int32,
			State:            entity.RoomState(state.Int32),
		}

		if len(configJSON) > 0 {
			var config entity.RoomConfig
			if err := json.Unmarshal(configJSON, &config); err == nil {
				room.Config = &config
			}
		}

		if len(metadataJSON) > 0 {
			var metadata map[string]string
			if err := json.Unmarshal(metadataJSON, &metadata); err == nil {
				room.Metadata = metadata
			}
		}

		if createdAt.Valid {
			room.CreatedAt = timestamppb.New(createdAt.Time)
		}

		if closedAt.Valid {
			room.ClosedAt = timestamppb.New(closedAt.Time)
		}

		if creatorID.Valid {
			room.CreatorId = creatorID.String
		}

		rooms = append(rooms, room)
	}

	return rooms, total, nil
}

// ═══════════════════════════════════════════════════════════════════════════
// Peer Operations
// ═══════════════════════════════════════════════════════════════════════════

// AddPeer adds a new peer to webrtc.peers table
func (r *Repository) AddPeer(ctx context.Context, peer *entity.Peer) error {
	metadataJSON, err := json.Marshal(peer.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	query := `
		INSERT INTO webrtc.peers (
			peer_id, room_id, display_name, state, metadata,
			joined_at, last_seen_at, user_agent, left_at
		) VALUES (
			gen_random_uuid(), $1, $2, $3, $4,
			CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, $5, $6
		) RETURNING peer_id, joined_at, last_seen_at`

	var peerID string
	var joinedAt, lastSeenAt sql.NullTime

	// Validate that room_id is provided and not empty
	roomID := peer.RoomId
	if roomID == "" {
		return fmt.Errorf("room_id is required")
	}

	err = r.db.WithContext(ctx).Raw(query,
		roomID,
		peer.DisplayName,
		int32(peer.State),
		metadataJSON,
		peer.UserAgent,
		nil, // left_at
	).Row().Scan(&peerID, &joinedAt, &lastSeenAt)

	if err != nil {
		return fmt.Errorf("failed to add peer: %w", err)
	}

	peer.PeerId = peerID
	if joinedAt.Valid {
		peer.JoinedAt = timestamppb.New(joinedAt.Time)
	}
	if lastSeenAt.Valid {
		peer.LastSeenAt = timestamppb.New(lastSeenAt.Time)
	}

	return nil
}

// GetPeer retrieves a peer by ID from webrtc.peers table
func (r *Repository) GetPeer(ctx context.Context, peerID string) (*entity.Peer, error) {
	query := `
		SELECT peer_id, room_id, display_name, state, metadata,
		       joined_at, last_seen_at, user_agent, left_at
		FROM webrtc.peers
		WHERE peer_id = $1`

	var (
		id, roomID, displayName, userAgent sql.NullString
		state                              sql.NullInt32
		metadataJSON                       []byte
		joinedAt, lastSeenAt, leftAt       sql.NullTime
	)

	err := r.db.WithContext(ctx).Raw(query, peerID).Row().Scan(
		&id, &roomID, &displayName, &state, &metadataJSON,
		&joinedAt, &lastSeenAt, &userAgent, &leftAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get peer: %w", err)
	}

	peer := &entity.Peer{
		PeerId:      id.String,
		RoomId:      roomID.String,
		DisplayName: displayName.String,
		State:       entity.PeerConnectionState(state.Int32),
		UserAgent:   userAgent.String,
	}

	if len(metadataJSON) > 0 {
		var metadata map[string]string
		if err := json.Unmarshal(metadataJSON, &metadata); err == nil {
			peer.Metadata = metadata
		}
	}

	if joinedAt.Valid {
		peer.JoinedAt = timestamppb.New(joinedAt.Time)
	}

	if lastSeenAt.Valid {
		peer.LastSeenAt = timestamppb.New(lastSeenAt.Time)
	}

	if leftAt.Valid {
		peer.LeftAt = timestamppb.New(leftAt.Time)
	}

	return peer, nil
}

// UpdatePeerState updates a peer's state in webrtc.peers table
func (r *Repository) UpdatePeerState(ctx context.Context, peerID string, state entity.PeerConnectionState) error {
	query := `
		UPDATE webrtc.peers
		SET state = $1, last_seen_at = CURRENT_TIMESTAMP
		WHERE peer_id = $2`

	result := r.db.WithContext(ctx).Exec(query, int32(state), peerID)
	if result.Error != nil {
		return fmt.Errorf("failed to update peer state: %w", result.Error)
	}

	return nil
}

// RemovePeer marks a peer as left in webrtc.peers table
func (r *Repository) RemovePeer(ctx context.Context, peerID string) error {
	query := `
		UPDATE webrtc.peers
		SET state = $1, left_at = CURRENT_TIMESTAMP
		WHERE peer_id = $2`

	result := r.db.WithContext(ctx).Exec(query, int32(entity.PeerConnectionState_PEER_CONNECTION_STATE_CLOSED), peerID)
	if result.Error != nil {
		return fmt.Errorf("failed to remove peer: %w", result.Error)
	}

	return nil
}

// ListPeersInRoom retrieves all peers in a room from webrtc.peers table
func (r *Repository) ListPeersInRoom(ctx context.Context, roomID string) ([]*entity.Peer, error) {
	query := `
		SELECT peer_id, room_id, display_name, state, metadata,
		       joined_at, last_seen_at, user_agent, left_at
		FROM webrtc.peers
		WHERE room_id = $1 AND left_at IS NULL
		ORDER BY joined_at ASC`

	rows, err := r.db.WithContext(ctx).Raw(query, roomID).Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to list peers: %w", err)
	}
	defer rows.Close()

	var peers []*entity.Peer
	for rows.Next() {
		var (
			id, rID, displayName, userAgent sql.NullString
			state                           sql.NullInt32
			metadataJSON                    []byte
			joinedAt, lastSeenAt, leftAt    sql.NullTime
		)

		if err := rows.Scan(
			&id, &rID, &displayName, &state, &metadataJSON,
			&joinedAt, &lastSeenAt, &userAgent, &leftAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan peer: %w", err)
		}

		peer := &entity.Peer{
			PeerId:      id.String,
			RoomId:      rID.String,
			DisplayName: displayName.String,
			State:       entity.PeerConnectionState(state.Int32),
			UserAgent:   userAgent.String,
		}

		if len(metadataJSON) > 0 {
			var metadata map[string]string
			if err := json.Unmarshal(metadataJSON, &metadata); err == nil {
				peer.Metadata = metadata
			}
		}

		if joinedAt.Valid {
			peer.JoinedAt = timestamppb.New(joinedAt.Time)
		}

		if lastSeenAt.Valid {
			peer.LastSeenAt = timestamppb.New(lastSeenAt.Time)
		}

		if leftAt.Valid {
			peer.LeftAt = timestamppb.New(leftAt.Time)
		}

		peers = append(peers, peer)
	}

	return peers, nil
}

// ═══════════════════════════════════════════════════════════════════════════
// Track Operations
// ═══════════════════════════════════════════════════════════════════════════

// PublishTrack adds a new track to webrtc_tracks table
func (r *Repository) PublishTrack(ctx context.Context, track *entity.Track) error {
	metadataJSON, err := json.Marshal(track.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	settingsJSON, err := json.Marshal(track.Settings)
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	// Validate that peer_id is provided and not empty
	peerID := track.PeerId
	if peerID == "" {
		return fmt.Errorf("peer_id is required")
	}

	query := `
		INSERT INTO webrtc.tracks (
			track_id, peer_id, type, label, muted, state, metadata, settings
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	result := r.db.WithContext(ctx).Exec(query,
		track.TrackId,
		peerID,
		int32(track.Type),
		track.Label,
		track.Muted,
		int32(track.State),
		metadataJSON,
		settingsJSON,
	)

	if result.Error != nil {
		return fmt.Errorf("failed to publish track: %w", result.Error)
	}

	return nil
}

// UnpublishTrack removes a track from webrtc_tracks table
func (r *Repository) UnpublishTrack(ctx context.Context, trackID string) error {
	query := `DELETE FROM webrtc.tracks WHERE track_id = $1`

	result := r.db.WithContext(ctx).Exec(query, trackID)
	if result.Error != nil {
		return fmt.Errorf("failed to unpublish track: %w", result.Error)
	}

	return nil
}

// MuteTrack updates track mute status in webrtc_tracks table
func (r *Repository) MuteTrack(ctx context.Context, trackID string, muted bool) error {
	query := `UPDATE webrtc.tracks SET muted = $1 WHERE track_id = $2`

	result := r.db.WithContext(ctx).Exec(query, muted, trackID)
	if result.Error != nil {
		return fmt.Errorf("failed to mute track: %w", result.Error)
	}

	return nil
}

// GetTrack retrieves a track by ID from webrtc_tracks table
func (r *Repository) GetTrack(ctx context.Context, trackID string) (*entity.Track, error) {
	query := `
		SELECT track_id, peer_id, type, label, muted, state, metadata, settings
		FROM webrtc.tracks
		WHERE track_id = $1`

	var (
		id, peerID, label          sql.NullString
		trackType, state           sql.NullInt32
		muted                      sql.NullBool
		metadataJSON, settingsJSON []byte
	)

	err := r.db.WithContext(ctx).Raw(query, trackID).Row().Scan(
		&id, &peerID, &trackType, &label, &muted, &state, &metadataJSON, &settingsJSON,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get track: %w", err)
	}

	track := &entity.Track{
		TrackId: id.String,
		PeerId:  peerID.String,
		Type:    entity.TrackType(trackType.Int32),
		Label:   label.String,
		Muted:   muted.Bool,
		State:   entity.TrackState(state.Int32),
	}

	if len(metadataJSON) > 0 {
		var metadata map[string]string
		if err := json.Unmarshal(metadataJSON, &metadata); err == nil {
			track.Metadata = metadata
		}
	}

	if len(settingsJSON) > 0 {
		var settings entity.TrackSettings
		if err := json.Unmarshal(settingsJSON, &settings); err == nil {
			track.Settings = &settings
		}
	}

	return track, nil
}

// ListTracks retrieves tracks with filters from webrtc_tracks table
func (r *Repository) ListTracks(ctx context.Context, roomID, peerID string, typeFilter entity.TrackType) ([]*entity.Track, error) {
	query := `
		SELECT t.track_id, t.peer_id, t.type, t.label, t.muted, t.state, t.metadata, t.settings
		FROM webrtc.tracks t
		INNER JOIN webrtc.peers p ON t.peer_id = p.peer_id
		WHERE p.room_id = $1`

	args := []interface{}{roomID}

	if peerID != "" {
		query += " AND t.peer_id = $2"
		args = append(args, peerID)
	}

	if typeFilter != entity.TrackType_TRACK_TYPE_UNSPECIFIED {
		if peerID != "" {
			query += " AND t.type = $3"
		} else {
			query += " AND t.type = $2"
		}
		args = append(args, int32(typeFilter))
	}

	query += " ORDER BY t.track_id"

	rows, err := r.db.WithContext(ctx).Raw(query, args...).Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to list tracks: %w", err)
	}
	defer rows.Close()

	var tracks []*entity.Track
	for rows.Next() {
		var (
			id, pID, label             sql.NullString
			trackType, state           sql.NullInt32
			muted                      sql.NullBool
			metadataJSON, settingsJSON []byte
		)

		if err := rows.Scan(
			&id, &pID, &trackType, &label, &muted, &state, &metadataJSON, &settingsJSON,
		); err != nil {
			return nil, fmt.Errorf("failed to scan track: %w", err)
		}

		track := &entity.Track{
			TrackId: id.String,
			PeerId:  pID.String,
			Type:    entity.TrackType(trackType.Int32),
			Label:   label.String,
			Muted:   muted.Bool,
			State:   entity.TrackState(state.Int32),
		}

		if len(metadataJSON) > 0 {
			var metadata map[string]string
			if err := json.Unmarshal(metadataJSON, &metadata); err == nil {
				track.Metadata = metadata
			}
		}

		if len(settingsJSON) > 0 {
			var settings entity.TrackSettings
			if err := json.Unmarshal(settingsJSON, &settings); err == nil {
				track.Settings = &settings
			}
		}

		tracks = append(tracks, track)
	}

	return tracks, nil
}

// UpdateTrack updates track settings and metadata in webrtc_tracks table
func (r *Repository) UpdateTrack(ctx context.Context, track *entity.Track) error {
	metadataJSON, err := json.Marshal(track.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	settingsJSON, err := json.Marshal(track.Settings)
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	query := `
		UPDATE webrtc.tracks
		SET settings = $1, metadata = $2
		WHERE track_id = $3`

	result := r.db.WithContext(ctx).Exec(query, settingsJSON, metadataJSON, track.TrackId)
	if result.Error != nil {
		return fmt.Errorf("failed to update track: %w", result.Error)
	}

	return nil
}
