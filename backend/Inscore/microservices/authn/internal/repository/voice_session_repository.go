package repository

import (
	"context"
	"strings"
	"time"

	voicev1 "github.com/newage-saint/insuretech/gen/go/insuretech/voice/entity/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

// VoiceSessionRepository provides access to authn_schema.voice_sessions.
type VoiceSessionRepository struct{ db *gorm.DB }

func NewVoiceSessionRepository(db *gorm.DB) *VoiceSessionRepository {
	return &VoiceSessionRepository{db: db}
}

func (r *VoiceSessionRepository) Create(ctx context.Context, s *voicev1.VoiceSession) error {
	return r.db.WithContext(ctx).Exec(
		`insert into authn_schema.voice_sessions (session_id, external_session_id, user_id, phone_number, language, status, intent, context, started_at, ended_at, duration_seconds, audit_info)
		 values (?, ?, ?, ?, ?, ?, ?, ?::jsonb, now(), ?, ?, '{}'::jsonb)`,
		s.Id,
		s.SessionId,
		nullableString(s.UserId),
		nullableString(s.PhoneNumber),
		s.Language,
		voiceStatusToString(s.Status),
		nullableString(s.Intent),
		nullableJSON(s.Context),
		nilOrTime(s.EndedAt),
		nilIfZero(s.DurationSeconds),
	).Error
}

func (r *VoiceSessionRepository) GetByID(ctx context.Context, id string) (*voicev1.VoiceSession, error) {
	return r.getOne(ctx, "session_id = ?", id)
}

func (r *VoiceSessionRepository) GetByExternalSessionID(ctx context.Context, externalID string) (*voicev1.VoiceSession, error) {
	return r.getOne(ctx, "external_session_id = ?", externalID)
}

const voiceCols = `session_id, external_session_id, user_id, phone_number, language, status, intent, started_at, ended_at, duration_seconds`

func scanVoiceSession(row interface{ Scan(...any) error }) (*voicev1.VoiceSession, error) {
	var s voicev1.VoiceSession
	var statusStr string
	var userID, phoneNumber, intent *string
	var startedAt time.Time
	var endedAt *time.Time
	var durationSeconds *int32
	if err := row.Scan(&s.Id, &s.SessionId, &userID, &phoneNumber, &s.Language, &statusStr, &intent, &startedAt, &endedAt, &durationSeconds); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	if userID != nil {
		s.UserId = *userID
	}
	if phoneNumber != nil {
		s.PhoneNumber = *phoneNumber
	}
	if intent != nil {
		s.Intent = *intent
	}
	if durationSeconds != nil {
		s.DurationSeconds = *durationSeconds
	}
	s.Status = voiceStatusFromString(statusStr)
	s.StartedAt = timestamppb.New(startedAt)
	if endedAt != nil {
		s.EndedAt = timestamppb.New(*endedAt)
	}
	return &s, nil
}

func (r *VoiceSessionRepository) ListByUser(ctx context.Context, userID string, limit, offset int) ([]*voicev1.VoiceSession, error) {
	q := `select ` + voiceCols + ` from authn_schema.voice_sessions where user_id = ? order by started_at desc`
	args := []any{userID}
	if limit > 0 {
		q += " limit ?"
		args = append(args, limit)
	}
	if offset > 0 {
		q += " offset ?"
		args = append(args, offset)
	}
	rows, err := r.db.WithContext(ctx).Raw(q, args...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*voicev1.VoiceSession
	for rows.Next() {
		s, err := scanVoiceSession(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func (r *VoiceSessionRepository) Complete(ctx context.Context, id string, status voicev1.SessionStatus, endedAt time.Time, durationSeconds *int32) error {
	upd := map[string]any{
		"status":   voiceStatusToString(status),
		"ended_at": endedAt,
	}
	if durationSeconds != nil {
		upd["duration_seconds"] = *durationSeconds
	}
	return r.db.WithContext(ctx).Table("authn_schema.voice_sessions").Where("session_id = ?", id).Updates(upd).Error
}

func (r *VoiceSessionRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Table("authn_schema.voice_sessions").Where("session_id = ?", id).Delete(map[string]any{}).Error
}

func (r *VoiceSessionRepository) getOne(ctx context.Context, where string, args ...any) (*voicev1.VoiceSession, error) {
	q := `select ` + voiceCols + ` from authn_schema.voice_sessions where ` + where + ` limit 1`
	row := r.db.WithContext(ctx).Raw(q, args...).Row()
	if err := row.Err(); err != nil {
		return nil, err
	}
	s, err := scanVoiceSession(row)
	if err != nil {
		return nil, err
	}
	if s.Id == "" {
		return nil, gorm.ErrRecordNotFound
	}
	return s, nil
}

func voiceStatusToString(st voicev1.SessionStatus) string {
	return strings.TrimPrefix(st.String(), "SESSION_STATUS_")
}

func voiceStatusFromString(s string) voicev1.SessionStatus {
	s = strings.ToUpper(strings.TrimSpace(s))
	s = strings.TrimPrefix(s, "SESSION_STATUS_")
	if v, ok := voicev1.SessionStatus_value["SESSION_STATUS_"+s]; ok {
		return voicev1.SessionStatus(v)
	}
	return voicev1.SessionStatus_SESSION_STATUS_UNSPECIFIED
}
