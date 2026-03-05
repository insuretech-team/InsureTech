package repository

import (
	"context"
	"strings"
	"time"

	"github.com/lib/pq"
	apikeyv1 "github.com/newage-saint/insuretech/gen/go/insuretech/apikey/entity/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

// ApiKeyRepository provides access to authn_schema.api_keys.
// Uses proto-generated ApiKey struct directly with GORM.
type ApiKeyRepository struct {
	db *gorm.DB
}

func NewApiKeyRepository(db *gorm.DB) *ApiKeyRepository {
	return &ApiKeyRepository{db: db}
}

// Create inserts a new API key using raw SQL (required for pq.Array and jsonb).
func (r *ApiKeyRepository) Create(ctx context.Context, k *apikeyv1.ApiKey) error {
	return r.db.WithContext(ctx).Exec(
		`insert into authn_schema.api_keys
			(api_key_id, key_hash, name, owner_type, owner_id, scopes, status, rate_limit_per_minute, expires_at, last_used_at, ip_whitelist, audit_info)
		 values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, coalesce(?, '{}'::jsonb))`,
		k.Id,
		k.KeyHash,
		k.Name,
		apiKeyOwnerTypeToString(k.OwnerType),
		k.OwnerId,
		pq.Array(k.Scopes),
		apiKeyStatusToString(k.Status),
		k.RateLimitPerMinute,
		nilOrTime(k.ExpiresAt),
		nilOrTime(k.LastUsedAt),
		pq.Array(k.IpWhitelist),
		nil, // audit_info: keep minimal
	).Error
}

// GetByID returns an API key by primary key.
func (r *ApiKeyRepository) GetByID(ctx context.Context, id string) (*apikeyv1.ApiKey, error) {
	return r.getOne(ctx, "api_key_id = ?", id)
}

// GetByKeyHash returns an API key by SHA-256 hash.
func (r *ApiKeyRepository) GetByKeyHash(ctx context.Context, keyHash string) (*apikeyv1.ApiKey, error) {
	return r.getOne(ctx, "key_hash = ?", keyHash)
}

// ListByOwner lists keys for an owner.
func (r *ApiKeyRepository) ListByOwner(ctx context.Context, ownerType apikeyv1.ApiKeyOwnerType, ownerID string, status *apikeyv1.ApiKeyStatus, limit, offset int) ([]*apikeyv1.ApiKey, error) {
	q := `select api_key_id, key_hash, name, owner_type, owner_id, scopes, status, rate_limit_per_minute, expires_at, last_used_at, ip_whitelist
	      from authn_schema.api_keys where owner_id = ?`
	args := []any{ownerID}

	if ownerType != apikeyv1.ApiKeyOwnerType_API_KEY_OWNER_TYPE_UNSPECIFIED {
		q += " and owner_type = ?"
		args = append(args, apiKeyOwnerTypeToString(ownerType))
	}
	if status != nil && *status != apikeyv1.ApiKeyStatus_API_KEY_STATUS_UNSPECIFIED {
		q += " and status = ?"
		args = append(args, apiKeyStatusToString(*status))
	}
	q += " order by api_key_id desc"
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

	var out []*apikeyv1.ApiKey
	for rows.Next() {
		k := &apikeyv1.ApiKey{}
		var ownerTypeStr, statusStr string
		var expiresAt, lastUsedAt *time.Time
		if err := rows.Scan(
			&k.Id, &k.KeyHash, &k.Name,
			&ownerTypeStr, &k.OwnerId,
			pq.Array(&k.Scopes),
			&statusStr,
			&k.RateLimitPerMinute,
			&expiresAt, &lastUsedAt,
			pq.Array(&k.IpWhitelist),
		); err != nil {
			return nil, err
		}
		k.OwnerType = apiKeyOwnerTypeFromString(ownerTypeStr)
		k.Status = apiKeyStatusFromString(statusStr)
		if expiresAt != nil {
			k.ExpiresAt = timestamppb.New(*expiresAt)
		}
		if lastUsedAt != nil {
			k.LastUsedAt = timestamppb.New(*lastUsedAt)
		}
		out = append(out, k)
	}
	return out, rows.Err()
}

// TouchLastUsed sets last_used_at.
func (r *ApiKeyRepository) TouchLastUsed(ctx context.Context, id string, at time.Time) error {
	return r.db.WithContext(ctx).Table("authn_schema.api_keys").Where("api_key_id = ?", id).Update("last_used_at", at).Error
}

// UpdateStatus updates status.
func (r *ApiKeyRepository) UpdateStatus(ctx context.Context, id string, status apikeyv1.ApiKeyStatus) error {
	return r.db.WithContext(ctx).Table("authn_schema.api_keys").Where("api_key_id = ?", id).Update("status", apiKeyStatusToString(status)).Error
}

// Revoke sets status=REVOKED.
func (r *ApiKeyRepository) Revoke(ctx context.Context, id string) error {
	return r.UpdateStatus(ctx, id, apikeyv1.ApiKeyStatus_API_KEY_STATUS_REVOKED)
}

// Delete hard-deletes the row.
func (r *ApiKeyRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Table("authn_schema.api_keys").Where("api_key_id = ?", id).Delete(map[string]any{}).Error
}

// SetExpiration sets the expires_at timestamp for an API key.
func (r *ApiKeyRepository) SetExpiration(ctx context.Context, id string, expiresAt time.Time) error {
	return r.db.WithContext(ctx).Table("authn_schema.api_keys").
		Where("api_key_id = ?", id).
		Update("expires_at", expiresAt).Error
}

// MarkAsRotating updates the status to ROTATING and sets expiration (single atomic update).
func (r *ApiKeyRepository) MarkAsRotating(ctx context.Context, id string, expiresAt time.Time) error {
	return r.db.WithContext(ctx).Table("authn_schema.api_keys").
		Where("api_key_id = ?", id).
		Updates(map[string]interface{}{
			"status":     apiKeyStatusToString(apikeyv1.ApiKeyStatus_API_KEY_STATUS_ROTATING),
			"expires_at": expiresAt,
		}).Error
}

func (r *ApiKeyRepository) getOne(ctx context.Context, where string, args ...any) (*apikeyv1.ApiKey, error) {
	q := `select api_key_id, key_hash, name, owner_type, owner_id, scopes, status, rate_limit_per_minute, expires_at, last_used_at, ip_whitelist
	      from authn_schema.api_keys where ` + where + ` limit 1`

	row := r.db.WithContext(ctx).Raw(q, args...).Row()
	if err := row.Err(); err != nil {
		return nil, err
	}

	k := &apikeyv1.ApiKey{}
	var ownerTypeStr, statusStr string
	var expiresAt, lastUsedAt *time.Time
	err := row.Scan(
		&k.Id, &k.KeyHash, &k.Name,
		&ownerTypeStr, &k.OwnerId,
		pq.Array(&k.Scopes),
		&statusStr,
		&k.RateLimitPerMinute,
		&expiresAt, &lastUsedAt,
		pq.Array(&k.IpWhitelist),
	)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	if k.Id == "" {
		return nil, gorm.ErrRecordNotFound
	}
	k.OwnerType = apiKeyOwnerTypeFromString(ownerTypeStr)
	k.Status = apiKeyStatusFromString(statusStr)
	if expiresAt != nil {
		k.ExpiresAt = timestamppb.New(*expiresAt)
	}
	if lastUsedAt != nil {
		k.LastUsedAt = timestamppb.New(*lastUsedAt)
	}
	return k, nil
}

func apiKeyOwnerTypeToString(t apikeyv1.ApiKeyOwnerType) string {
	s := t.String()
	s = strings.TrimPrefix(s, "API_KEY_OWNER_TYPE_")
	return s
}

func apiKeyOwnerTypeFromString(s string) apikeyv1.ApiKeyOwnerType {
	s = strings.ToUpper(strings.TrimSpace(s))
	s = strings.TrimPrefix(s, "API_KEY_OWNER_TYPE_")
	if v, ok := apikeyv1.ApiKeyOwnerType_value["API_KEY_OWNER_TYPE_"+s]; ok {
		return apikeyv1.ApiKeyOwnerType(v)
	}
	return apikeyv1.ApiKeyOwnerType_API_KEY_OWNER_TYPE_UNSPECIFIED
}

func apiKeyStatusToString(st apikeyv1.ApiKeyStatus) string {
	s := st.String()
	s = strings.TrimPrefix(s, "API_KEY_STATUS_")
	return s
}

func apiKeyStatusFromString(s string) apikeyv1.ApiKeyStatus {
	s = strings.ToUpper(strings.TrimSpace(s))
	s = strings.TrimPrefix(s, "API_KEY_STATUS_")
	if v, ok := apikeyv1.ApiKeyStatus_value["API_KEY_STATUS_"+s]; ok {
		return apikeyv1.ApiKeyStatus(v)
	}
	return apikeyv1.ApiKeyStatus_API_KEY_STATUS_UNSPECIFIED
}
