package repository

import (
	"context"
	"strings"
	"time"

	authnentityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/entity/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type SessionRepository struct {
	db *gorm.DB
}

const sessionCols = `session_id, user_id, session_type, session_token_hash, session_token_lookup, access_token_jti, refresh_token_jti, access_token_expires_at, refresh_token_expires_at, expires_at, ip_address, user_agent, device_id, device_name, device_type, created_at, last_activity_at, is_active, csrf_token`

func scanSession(row interface{ Scan(...any) error }) (*authnentityv1.Session, error) {
	var s authnentityv1.Session
	var sessionTypeStr, deviceTypeStr string
	var accessExp, refreshExp, expiresAt *time.Time
	var createdAt, lastActivityAt time.Time
	var ipAddress, userAgent, deviceID, deviceName, csrfToken *string
	var sessionTokenHash, sessionTokenLookup, accessJTI, refreshJTI *string
	if err := row.Scan(
		&s.SessionId, &s.UserId, &sessionTypeStr,
		&sessionTokenHash, &sessionTokenLookup,
		&accessJTI, &refreshJTI,
		&accessExp, &refreshExp, &expiresAt,
		&ipAddress, &userAgent, &deviceID, &deviceName, &deviceTypeStr,
		&createdAt, &lastActivityAt, &s.IsActive, &csrfToken,
	); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	s.SessionType = sessionTypeFromString(sessionTypeStr)
	s.DeviceType = deviceTypeFromString(deviceTypeStr)
	if accessExp != nil {
		s.AccessTokenExpiresAt = timestamppb.New(*accessExp)
	}
	if refreshExp != nil {
		s.RefreshTokenExpiresAt = timestamppb.New(*refreshExp)
	}
	if expiresAt != nil {
		s.ExpiresAt = timestamppb.New(*expiresAt)
	}
	s.CreatedAt = timestamppb.New(createdAt)
	s.LastActivityAt = timestamppb.New(lastActivityAt)
	if sessionTokenHash != nil {
		s.SessionTokenHash = *sessionTokenHash
	}
	if sessionTokenLookup != nil {
		s.SessionTokenLookup = *sessionTokenLookup
	}
	if accessJTI != nil {
		s.AccessTokenJti = *accessJTI
	}
	if refreshJTI != nil {
		s.RefreshTokenJti = *refreshJTI
	}
	if ipAddress != nil {
		s.IpAddress = *ipAddress
	}
	if userAgent != nil {
		s.UserAgent = *userAgent
	}
	if deviceID != nil {
		s.DeviceId = *deviceID
	}
	if deviceName != nil {
		s.DeviceName = *deviceName
	}
	if csrfToken != nil {
		s.CsrfToken = *csrfToken
	}
	return &s, nil
}

func sessionTypeFromString(s string) authnentityv1.SessionType {
	s = strings.ToUpper(strings.TrimSpace(s))
	s = strings.TrimPrefix(s, "SESSION_TYPE_")
	if v, ok := authnentityv1.SessionType_value["SESSION_TYPE_"+s]; ok {
		return authnentityv1.SessionType(v)
	}
	return authnentityv1.SessionType_SESSION_TYPE_JWT
}

func deviceTypeFromString(s string) authnentityv1.DeviceType {
	s = strings.ToUpper(strings.TrimSpace(s))
	s = strings.TrimPrefix(s, "DEVICE_TYPE_")
	if v, ok := authnentityv1.DeviceType_value["DEVICE_TYPE_"+s]; ok {
		return authnentityv1.DeviceType(v)
	}
	return authnentityv1.DeviceType_DEVICE_TYPE_UNSPECIFIED
}

func (r *SessionRepository) getOne(ctx context.Context, where string, args ...any) (*authnentityv1.Session, error) {
	q := `select ` + sessionCols + ` from authn_schema.sessions where ` + where + ` limit 1`
	row := r.db.WithContext(ctx).Raw(q, args...).Row()
	if err := row.Err(); err != nil {
		return nil, err
	}
	return scanSession(row)
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

// Create creates a new session in the database using raw SQL to ensure
// enum types are stored as strings (not integer proto values).
func (r *SessionRepository) Create(ctx context.Context, session *authnentityv1.Session) error {
	now := time.Now()
	session.CreatedAt = timestamppb.New(now)
	session.LastActivityAt = timestamppb.New(now)

	return r.db.WithContext(ctx).Exec(`
		INSERT INTO authn_schema.sessions
			(session_id, user_id, session_type, session_token_hash, session_token_lookup,
			 access_token_jti, refresh_token_jti,
			 access_token_expires_at, refresh_token_expires_at, expires_at,
			 ip_address, user_agent, device_id, device_name, device_type,
			 created_at, last_activity_at, is_active, csrf_token)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
		session.SessionId,
		session.UserId,
		sessionTypeToString(session.SessionType),
		nullableString(session.SessionTokenHash),
		nullableString(session.SessionTokenLookup),
		nullableString(session.AccessTokenJti),
		nullableString(session.RefreshTokenJti),
		nilOrTime(session.AccessTokenExpiresAt),
		nilOrTime(session.RefreshTokenExpiresAt),
		nilOrTime(session.ExpiresAt),
		nullableString(session.IpAddress),
		nullableString(session.UserAgent),
		nullableString(session.DeviceId),
		nullableString(session.DeviceName),
		deviceTypeToString(session.DeviceType),
		now,
		now,
		true,
		nullableString(session.CsrfToken),
	).Error
}

func deviceTypeToString(dt authnentityv1.DeviceType) string {
	s := dt.String()
	return strings.TrimPrefix(s, "DEVICE_TYPE_")
}

// GetByID retrieves a session by session_id
func (r *SessionRepository) GetByID(ctx context.Context, id string) (*authnentityv1.Session, error) {
	return r.getOne(ctx, "session_id = ? AND is_active = ?", id, true)
}

// GetByTokenLookup retrieves a server-side session by deterministic lookup hash (sha256 hex).
// session_token_lookup is already unique so no session_type filter needed.
func (r *SessionRepository) GetByTokenLookup(ctx context.Context, lookup string) (*authnentityv1.Session, error) {
	return r.getOne(ctx, "session_token_lookup = ? AND is_active = ?", lookup, true)
}

// GetByRefreshToken retrieves a session by refresh_token_jti (for JWT sessions)
func (r *SessionRepository) GetByRefreshToken(ctx context.Context, refreshToken string) (*authnentityv1.Session, error) {
	return r.getOne(ctx, "refresh_token_jti = ? AND session_type = ? AND is_active = ?", refreshToken, "JWT", true)
}

// UpdateLastActivity updates the last_activity_at timestamp (for sliding expiration)
func (r *SessionRepository) UpdateLastActivity(ctx context.Context, sessionID string) error {
	return r.db.WithContext(ctx).
		Table("authn_schema.sessions").
		Where("session_id = ?", sessionID).
		Update("last_activity_at", time.Now()).Error
}

// UpdateTokens updates access and refresh token JTIs (for token rotation)
func (r *SessionRepository) UpdateTokens(ctx context.Context, sessionID, newAccessJTI, newRefreshJTI string) error {
	now := time.Now()
	updates := map[string]interface{}{
		"access_token_jti":         newAccessJTI,
		"refresh_token_jti":        newRefreshJTI,
		"access_token_expires_at":  now.Add(15 * time.Minute),
		"refresh_token_expires_at": now.Add(7 * 24 * time.Hour),
		"last_activity_at":         now,
	}

	return r.db.WithContext(ctx).
		Table("authn_schema.sessions").
		Where("session_id = ?", sessionID).
		Updates(updates).Error
}

// Revoke marks a session as inactive (soft delete)
func (r *SessionRepository) Revoke(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Table("authn_schema.sessions").
		Where("session_id = ?", id).
		Update("is_active", false).Error
}

// RevokeAllByUserID revokes all sessions for a user (logout from all devices)
func (r *SessionRepository) RevokeAllByUserID(ctx context.Context, userID string, excludeSessionID string) error {
	_, err := r.RevokeAllByUserIDWithCount(ctx, userID, excludeSessionID)
	return err
}

// RevokeAllByUserIDWithCount revokes all sessions for a user and returns how many rows were affected.
func (r *SessionRepository) RevokeAllByUserIDWithCount(ctx context.Context, userID string, excludeSessionID string) (int64, error) {
	query := r.db.WithContext(ctx).
		Table("authn_schema.sessions").
		Where("user_id = ?", userID)

	if excludeSessionID != "" {
		query = query.Where("session_id != ?", excludeSessionID)
	}

	result := query.Update("is_active", false)
	return result.RowsAffected, result.Error
}

// ListByUserID lists all sessions for a user with optional filters
func (r *SessionRepository) ListByUserID(ctx context.Context, userID string, activeOnly bool, sessionType *authnentityv1.SessionType) ([]*authnentityv1.Session, error) {
	q := `select ` + sessionCols + ` from authn_schema.sessions where user_id = ?`
	args := []any{userID}

	if activeOnly {
		q += " and is_active = ?"
		args = append(args, true)
	}
	if sessionType != nil {
		q += " and session_type = ?"
		args = append(args, sessionTypeToString(*sessionType))
	}
	q += " order by created_at desc"

	rows, err := r.db.WithContext(ctx).Raw(q, args...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*authnentityv1.Session
	for rows.Next() {
		s, err := scanSession(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

// sessionTypeToString converts a SessionType proto enum to the string stored in the DB.
func sessionTypeToString(st authnentityv1.SessionType) string {
	switch st {
	case authnentityv1.SessionType_SESSION_TYPE_JWT:
		return "JWT"
	case authnentityv1.SessionType_SESSION_TYPE_SERVER_SIDE:
		return "SERVER_SIDE"
	default:
		return "JWT"
	}
}

// CleanupExpiredSessions deletes expired sessions (background job)
func (r *SessionRepository) CleanupExpiredSessions(ctx context.Context) (int64, error) {
	result := r.db.WithContext(ctx).
		Table("authn_schema.sessions").
		Where("expires_at < ? AND is_active = ?", time.Now(), true).
		Delete(map[string]any{})

	return result.RowsAffected, result.Error
}
