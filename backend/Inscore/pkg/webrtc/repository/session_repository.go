package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/newage-saint/insuretech/gen/go/insuretech/webrtc/v1/entity"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ═══════════════════════════════════════════════════════════════════════════
// Session Operations - For Analytics and History
// ═══════════════════════════════════════════════════════════════════════════

// CreateRoomSession creates a new room session record for analytics
func (r *Repository) CreateRoomSession(ctx context.Context, session *entity.RoomSession) error {
	query := `
		INSERT INTO webrtc.room_sessions (
			session_id, room_id, started_at, ended_at, duration_seconds,
			peak_participants, total_participants, avg_call_quality, metrics
		) VALUES (
			gen_random_uuid(), $1, CURRENT_TIMESTAMP, $2, $3,
			$4, $5, $6, $7
		) RETURNING session_id, started_at`

	var sessionID string
	var startedAt sql.NullTime

	err := r.db.WithContext(ctx).Raw(query,
		session.RoomId,
		nil, // ended_at
		0,   // duration_seconds
		0,   // peak_participants
		0,   // total_participants
		0,   // avg_call_quality
		nil, // metrics
	).Row().Scan(&sessionID, &startedAt)

	if err != nil {
		return fmt.Errorf("failed to create room session: %w", err)
	}

	session.SessionId = sessionID
	if startedAt.Valid {
		session.StartedAt = timestamppb.New(startedAt.Time)
	}

	return nil
}

// UpdateRoomSession updates an existing room session
func (r *Repository) UpdateRoomSession(ctx context.Context, session *entity.RoomSession) error {
	query := `
		UPDATE webrtc.room_sessions
		SET ended_at = $1
		WHERE session_id = $2`

	var endedAt interface{}
	if session.EndedAt != nil {
		endedAt = session.EndedAt.AsTime()
	}

	result := r.db.WithContext(ctx).Exec(query,
		endedAt,
		session.SessionId,
	)

	if result.Error != nil {
		return fmt.Errorf("failed to update room session: %w", result.Error)
	}

	return nil
}

// GetRoomSession retrieves a room session by ID
func (r *Repository) GetRoomSession(ctx context.Context, sessionID string) (*entity.RoomSession, error) {
	query := `
		SELECT session_id, room_id, started_at, ended_at
		FROM webrtc.room_sessions
		WHERE session_id = $1`

	var (
		id, roomID         sql.NullString
		startedAt, endedAt sql.NullTime
	)

	err := r.db.WithContext(ctx).Raw(query, sessionID).Row().Scan(
		&id, &roomID, &startedAt, &endedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get room session: %w", err)
	}

	session := &entity.RoomSession{
		SessionId: id.String,
		RoomId:    roomID.String,
	}

	if startedAt.Valid {
		session.StartedAt = timestamppb.New(startedAt.Time)
	}

	if endedAt.Valid {
		session.EndedAt = timestamppb.New(endedAt.Time)
	}

	return session, nil
}

// GetActiveRoomSession retrieves the active session for a room
func (r *Repository) GetActiveRoomSession(ctx context.Context, roomID string) (*entity.RoomSession, error) {
	query := `
		SELECT session_id, room_id, started_at, ended_at
		FROM webrtc.room_sessions
		WHERE room_id = $1 AND ended_at IS NULL
		ORDER BY started_at DESC
		LIMIT 1`

	var (
		id, rID            sql.NullString
		startedAt, endedAt sql.NullTime
	)

	err := r.db.WithContext(ctx).Raw(query, roomID).Row().Scan(
		&id, &rID, &startedAt, &endedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get active room session: %w", err)
	}

	session := &entity.RoomSession{
		SessionId: id.String,
		RoomId:    rID.String,
	}

	if startedAt.Valid {
		session.StartedAt = timestamppb.New(startedAt.Time)
	}

	if endedAt.Valid {
		session.EndedAt = timestamppb.New(endedAt.Time)
	}

	return session, nil
}

// EndRoomSession marks a room session as ended
func (r *Repository) EndRoomSession(ctx context.Context, sessionID string) error {
	query := `
		UPDATE webrtc.room_sessions
		SET ended_at = CURRENT_TIMESTAMP,
		    duration_seconds = EXTRACT(EPOCH FROM (CURRENT_TIMESTAMP - started_at))
		WHERE session_id = $1`

	result := r.db.WithContext(ctx).Exec(query, sessionID)
	if result.Error != nil {
		return fmt.Errorf("failed to end room session: %w", result.Error)
	}

	return nil
}

// CreatePeerSession creates a new peer session record
func (r *Repository) CreatePeerSession(ctx context.Context, session *entity.PeerSession) error {
	query := `
		INSERT INTO webrtc.peer_sessions (
			peer_session_id, session_id, peer_id, display_name,
			joined_at, left_at, duration_seconds, avg_packet_loss,
			avg_bitrate_kbps, reconnection_count, user_agent, metrics
		) VALUES (
			gen_random_uuid(), $1, $2, $3,
			CURRENT_TIMESTAMP, $4, $5, $6,
			$7, $8, $9, $10
		) RETURNING peer_session_id, joined_at`

	var peerSessionID string
	var joinedAt sql.NullTime

	// Note: Inserting defaults/dummies for fields not in entity but present in DB
	err := r.db.WithContext(ctx).Raw(query,
		session.SessionId,
		session.PeerId,
		"",  // display_name (missing in entity)
		nil, // left_at
		0,   // duration_seconds
		0,   // avg_packet_loss
		0,   // avg_bitrate_kbps
		0,   // reconnection_count
		"",  // user_agent (missing in entity)
		nil, // metrics
	).Row().Scan(&peerSessionID, &joinedAt)

	if err != nil {
		return fmt.Errorf("failed to create peer session: %w", err)
	}

	session.PeerSessionId = peerSessionID
	if joinedAt.Valid {
		session.JoinedAt = timestamppb.New(joinedAt.Time)
	}

	return nil
}

// UpdatePeerSession updates an existing peer session
func (r *Repository) UpdatePeerSession(ctx context.Context, session *entity.PeerSession) error {
	query := `
		UPDATE webrtc.peer_sessions
		SET left_at = $1
		WHERE peer_session_id = $2`

	var leftAt interface{}
	if session.LeftAt != nil {
		leftAt = session.LeftAt.AsTime()
	}

	result := r.db.WithContext(ctx).Exec(query,
		leftAt,
		session.PeerSessionId,
	)

	if result.Error != nil {
		return fmt.Errorf("failed to update peer session: %w", result.Error)
	}

	return nil
}

// GetPeerSession retrieves a peer session by ID
func (r *Repository) GetPeerSession(ctx context.Context, peerSessionID string) (*entity.PeerSession, error) {
	query := `
		SELECT peer_session_id, session_id, peer_id,
		       joined_at, left_at
		FROM webrtc.peer_sessions
		WHERE peer_session_id = $1`

	var (
		id, sessionID, peerID sql.NullString
		joinedAt, leftAt      sql.NullTime
	)

	err := r.db.WithContext(ctx).Raw(query, peerSessionID).Row().Scan(
		&id, &sessionID, &peerID,
		&joinedAt, &leftAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get peer session: %w", err)
	}

	session := &entity.PeerSession{
		PeerSessionId: id.String,
		SessionId:     sessionID.String,
		PeerId:        peerID.String,
	}

	if joinedAt.Valid {
		session.JoinedAt = timestamppb.New(joinedAt.Time)
	}

	if leftAt.Valid {
		session.LeftAt = timestamppb.New(leftAt.Time)
	}

	return session, nil
}

// EndPeerSession marks a peer session as ended
func (r *Repository) EndPeerSession(ctx context.Context, peerSessionID string) error {
	query := `
		UPDATE webrtc.peer_sessions
		SET left_at = CURRENT_TIMESTAMP,
		    duration_seconds = EXTRACT(EPOCH FROM (CURRENT_TIMESTAMP - joined_at))
		WHERE peer_session_id = $1`

	result := r.db.WithContext(ctx).Exec(query, peerSessionID)
	if result.Error != nil {
		return fmt.Errorf("failed to end peer session: %w", result.Error)
	}

	return nil
}

// GetActivePeerSession retrieves the active peer session for a peer in a session
func (r *Repository) GetActivePeerSession(ctx context.Context, sessionID, peerID string) (*entity.PeerSession, error) {
	query := `
		SELECT peer_session_id, session_id, peer_id,
		       joined_at, left_at
		FROM webrtc.peer_sessions
		WHERE session_id = $1 AND peer_id = $2 AND left_at IS NULL
		ORDER BY joined_at DESC
		LIMIT 1`

	var (
		id, sID, pID     sql.NullString
		joinedAt, leftAt sql.NullTime
	)

	err := r.db.WithContext(ctx).Raw(query, sessionID, peerID).Row().Scan(
		&id, &sID, &pID,
		&joinedAt, &leftAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get active peer session: %w", err)
	}

	session := &entity.PeerSession{
		PeerSessionId: id.String,
		SessionId:     sID.String,
		PeerId:        pID.String,
	}

	if joinedAt.Valid {
		session.JoinedAt = timestamppb.New(joinedAt.Time)
	}

	if leftAt.Valid {
		session.LeftAt = timestamppb.New(leftAt.Time)
	}

	return session, nil
}

// UpdateSessionMetrics updates session metrics with calculated values
func (r *Repository) UpdateSessionMetrics(ctx context.Context, sessionID string) error {
	// Calculate and update peak participants
	query := `
		UPDATE webrtc.room_sessions
		SET peak_participants = (
			SELECT MAX(concurrent_count)
			FROM (
				SELECT COUNT(*) as concurrent_count
				FROM webrtc.peer_sessions
				WHERE session_id = $1
				GROUP BY date_trunc('minute', joined_at)
			) subq
		),
		total_participants = (
			SELECT COUNT(DISTINCT peer_id)
			FROM webrtc.peer_sessions
			WHERE session_id = $1
		)
		WHERE session_id = $1`

	result := r.db.WithContext(ctx).Exec(query, sessionID)
	if result.Error != nil {
		return fmt.Errorf("failed to update session metrics: %w", result.Error)
	}

	return nil
}

// IncrementReconnectionCount increments the reconnection count for a peer session
func (r *Repository) IncrementReconnectionCount(ctx context.Context, peerSessionID string) error {
	query := `
		UPDATE webrtc.peer_sessions
		SET reconnection_count = reconnection_count + 1
		WHERE peer_session_id = $1`

	result := r.db.WithContext(ctx).Exec(query, peerSessionID)
	if result.Error != nil {
		return fmt.Errorf("failed to increment reconnection count: %w", result.Error)
	}

	return nil
}

// UpdatePeerSessionStats updates statistics for a peer session
func (r *Repository) UpdatePeerSessionStats(ctx context.Context, peerSessionID string, packetLoss, bitrate float64) error {
	// Calculate running average
	query := `
		UPDATE webrtc.peer_sessions
		SET avg_packet_loss = COALESCE(
			(avg_packet_loss * (duration_seconds / GREATEST(EXTRACT(EPOCH FROM (CURRENT_TIMESTAMP - joined_at)), 1))) + 
			($2 * (1 - (duration_seconds / GREATEST(EXTRACT(EPOCH FROM (CURRENT_TIMESTAMP - joined_at)), 1)))),
			$2
		),
		avg_bitrate_kbps = COALESCE(
			(avg_bitrate_kbps * (duration_seconds / GREATEST(EXTRACT(EPOCH FROM (CURRENT_TIMESTAMP - joined_at)), 1))) + 
			($3 * (1 - (duration_seconds / GREATEST(EXTRACT(EPOCH FROM (CURRENT_TIMESTAMP - joined_at)), 1)))),
			$3
		)
		WHERE peer_session_id = $1`

	result := r.db.WithContext(ctx).Exec(query, peerSessionID, packetLoss, bitrate)
	if result.Error != nil {
		return fmt.Errorf("failed to update peer session stats: %w", result.Error)
	}

	return nil
}

// ListRoomSessions retrieves all sessions for a room
func (r *Repository) ListRoomSessions(ctx context.Context, roomID string, limit int, startTime, endTime time.Time) ([]*entity.RoomSession, error) {
	query := `
		SELECT session_id, room_id, started_at, ended_at
		FROM webrtc.room_sessions
		WHERE room_id = $1`

	args := []interface{}{roomID}

	if !startTime.IsZero() {
		query += " AND started_at >= $" + fmt.Sprintf("%d", len(args)+1)
		args = append(args, startTime)
	}

	if !endTime.IsZero() {
		query += " AND started_at <= $" + fmt.Sprintf("%d", len(args)+1)
		args = append(args, endTime)
	}

	query += " ORDER BY started_at DESC"

	if limit > 0 {
		query += " LIMIT $" + fmt.Sprintf("%d", len(args)+1)
		args = append(args, limit)
	}

	rows, err := r.db.WithContext(ctx).Raw(query, args...).Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to list room sessions: %w", err)
	}
	defer rows.Close()

	var sessions []*entity.RoomSession

	for rows.Next() {
		var (
			id, roomID         sql.NullString
			startedAt, endedAt sql.NullTime
		)

		if err := rows.Scan(
			&id, &roomID, &startedAt, &endedAt,
		); err != nil {
			continue
		}

		session := &entity.RoomSession{
			SessionId: id.String,
			RoomId:    roomID.String,
		}

		if startedAt.Valid {
			session.StartedAt = timestamppb.New(startedAt.Time)
		}

		if endedAt.Valid {
			session.EndedAt = timestamppb.New(endedAt.Time)
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}
