package repository

import (
	"context"
	"time"

	apikeyv1 "github.com/newage-saint/insuretech/gen/go/insuretech/apikey/entity/v1"
	"gorm.io/gorm"
)

// ApiKeyUsageRepository provides access to authn_schema.api_key_usage.
// Uses proto-generated ApiKeyUsage struct directly with GORM.
type ApiKeyUsageRepository struct{ db *gorm.DB }

func NewApiKeyUsageRepository(db *gorm.DB) *ApiKeyUsageRepository {
	return &ApiKeyUsageRepository{db: db}
}

// Create inserts a new usage record using raw SQL (required for inet/jsonb casts).
func (r *ApiKeyUsageRepository) Create(ctx context.Context, u *apikeyv1.ApiKeyUsage) error {
	return r.db.WithContext(ctx).Exec(
		`insert into authn_schema.api_key_usage
			(usage_id, api_key_id, endpoint, http_method, status_code, response_time_ms, request_ip, user_agent, request_payload, response_payload, trace_id, timestamp)
		 values (?, ?, ?, ?, ?, ?, ?::inet, ?, ?::jsonb, ?::jsonb, ?, ?)`,
		u.Id,
		u.ApiKeyId,
		u.Endpoint,
		u.HttpMethod,
		u.StatusCode,
		nilIfZero(u.ResponseTimeMs),
		nullableString(u.RequestIp),
		nullableString(u.UserAgent),
		nullableJSON(u.RequestPayload),
		nullableJSON(u.ResponsePayload),
		nullableString(u.TraceId),
		nilOrTime(u.Timestamp),
	).Error
}

// GetByID returns a usage record by primary key.
func (r *ApiKeyUsageRepository) GetByID(ctx context.Context, id string) (*apikeyv1.ApiKeyUsage, error) {
	k := &apikeyv1.ApiKeyUsage{}
	if err := r.db.WithContext(ctx).
		Table("authn_schema.api_key_usage").
		Where("usage_id = ?", id).
		Limit(1).
		Scan(k).Error; err != nil {
		return nil, err
	}
	if k.Id == "" {
		return nil, gorm.ErrRecordNotFound
	}
	return k, nil
}

// ListByApiKey lists usage records for an API key with optional time range.
func (r *ApiKeyUsageRepository) ListByApiKey(ctx context.Context, apiKeyID string, from, to *time.Time, limit, offset int) ([]*apikeyv1.ApiKeyUsage, error) {
	q := r.db.WithContext(ctx).Table("authn_schema.api_key_usage").Where("api_key_id = ?", apiKeyID)
	if from != nil {
		q = q.Where("timestamp >= ?", *from)
	}
	if to != nil {
		q = q.Where("timestamp <= ?", *to)
	}
	if limit > 0 {
		q = q.Limit(limit)
	}
	if offset > 0 {
		q = q.Offset(offset)
	}
	var out []*apikeyv1.ApiKeyUsage
	if err := q.Order("timestamp desc").Scan(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteByApiKey hard-deletes all usage records for an API key.
func (r *ApiKeyUsageRepository) DeleteByApiKey(ctx context.Context, apiKeyID string) (int64, error) {
	res := r.db.WithContext(ctx).Table("authn_schema.api_key_usage").Where("api_key_id = ?", apiKeyID).Delete(map[string]any{})
	return res.RowsAffected, res.Error
}
