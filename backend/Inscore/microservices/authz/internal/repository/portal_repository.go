package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/lib/pq"
	entityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/entity/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

// PortalRepo implements domain.PortalConfigRepository using proto structs directly.
type PortalRepo struct{ db *gorm.DB }

func NewPortalRepo(db *gorm.DB) *PortalRepo { return &PortalRepo{db: db} }

func (r *PortalRepo) Upsert(ctx context.Context, cfg *entityv1.PortalConfig) (*entityv1.PortalConfig, error) {
	if cfg == nil {
		return nil, errors.New("portalConfig.Upsert: nil config")
	}
	portal := cfg.Portal.String()
	var updatedBy any = nil
	if cfg.UpdatedBy != "" {
		if _, err := uuid.Parse(cfg.UpdatedBy); err == nil {
			updatedBy = cfg.UpdatedBy
		}
	}
	err := r.db.WithContext(ctx).Exec(
		`INSERT INTO authz_schema.portal_configs
		 (portal, mfa_required, mfa_methods, access_token_ttl_seconds, refresh_token_ttl_seconds, session_ttl_seconds, idle_timeout_seconds, allow_concurrent_sessions, max_concurrent_sessions, updated_by)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		 ON CONFLICT (portal) DO UPDATE SET
		   mfa_required = EXCLUDED.mfa_required,
		   mfa_methods = EXCLUDED.mfa_methods,
		   access_token_ttl_seconds = EXCLUDED.access_token_ttl_seconds,
		   refresh_token_ttl_seconds = EXCLUDED.refresh_token_ttl_seconds,
		   session_ttl_seconds = EXCLUDED.session_ttl_seconds,
		   idle_timeout_seconds = EXCLUDED.idle_timeout_seconds,
		   allow_concurrent_sessions = EXCLUDED.allow_concurrent_sessions,
		   max_concurrent_sessions = EXCLUDED.max_concurrent_sessions,
		   updated_by = EXCLUDED.updated_by,
		   updated_at = NOW()`,
		portal,
		cfg.MfaRequired,
		pq.Array(cfg.MfaMethods),
		cfg.AccessTokenTtlSeconds,
		cfg.RefreshTokenTtlSeconds,
		cfg.SessionTtlSeconds,
		cfg.IdleTimeoutSeconds,
		cfg.AllowConcurrentSessions,
		cfg.MaxConcurrentSessions,
		updatedBy,
	).Error
	if err != nil {
		return nil, errors.New("portalConfig.Upsert: " + err.Error())
	}

	return cfg, nil
}

func (r *PortalRepo) GetByPortal(ctx context.Context, portal entityv1.Portal) (*entityv1.PortalConfig, error) {
	var (
		portalStr string
		methods   pq.StringArray
		cfg       entityv1.PortalConfig
	)

	err := r.db.WithContext(ctx).Raw(
		`SELECT portal, mfa_required, mfa_methods, access_token_ttl_seconds, refresh_token_ttl_seconds, session_ttl_seconds, idle_timeout_seconds, allow_concurrent_sessions, max_concurrent_sessions
		   FROM authz_schema.portal_configs
		  WHERE portal = ?
		  LIMIT 1`,
		portal.String(),
	).Row().Scan(
		&portalStr,
		&cfg.MfaRequired,
		&methods,
		&cfg.AccessTokenTtlSeconds,
		&cfg.RefreshTokenTtlSeconds,
		&cfg.SessionTtlSeconds,
		&cfg.IdleTimeoutSeconds,
		&cfg.AllowConcurrentSessions,
		&cfg.MaxConcurrentSessions,
	)
	if err != nil {
		return nil, errors.New("portalConfig.GetByPortal: " + err.Error())
	}
	cfg.MfaMethods = []string(methods)
	if v, ok := entityv1.Portal_value[portalStr]; ok {
		cfg.Portal = entityv1.Portal(v)
	}
	return &cfg, nil
}

func (r *PortalRepo) List(ctx context.Context) ([]*entityv1.PortalConfig, error) {
	rows, err := r.db.WithContext(ctx).Raw(
		`SELECT portal, mfa_required, mfa_methods, access_token_ttl_seconds, refresh_token_ttl_seconds, session_ttl_seconds, idle_timeout_seconds, allow_concurrent_sessions, max_concurrent_sessions
		   FROM authz_schema.portal_configs`,
	).Rows()
	if err != nil {
		return nil, errors.New("portalConfig.List: " + err.Error())
	}
	defer rows.Close()

	var cfgs []*entityv1.PortalConfig
	for rows.Next() {
		var (
			portalStr string
			methods   pq.StringArray
			cfg       entityv1.PortalConfig
		)
		if err := rows.Scan(
			&portalStr,
			&cfg.MfaRequired,
			&methods,
			&cfg.AccessTokenTtlSeconds,
			&cfg.RefreshTokenTtlSeconds,
			&cfg.SessionTtlSeconds,
			&cfg.IdleTimeoutSeconds,
			&cfg.AllowConcurrentSessions,
			&cfg.MaxConcurrentSessions,
		); err != nil {
			return nil, errors.New("portalConfig.List scan: " + err.Error())
		}
		cfg.MfaMethods = []string(methods)
		if v, ok := entityv1.Portal_value[portalStr]; ok {
			cfg.Portal = entityv1.Portal(v)
		}
		c := cfg
		cfgs = append(cfgs, &c)
	}

	return cfgs, nil
}

// TokenConfigRepo implements domain.TokenConfigRepository using proto structs directly.
type TokenConfigRepo struct{ db *gorm.DB }

func NewTokenConfigRepo(db *gorm.DB) *TokenConfigRepo { return &TokenConfigRepo{db: db} }

func (r *TokenConfigRepo) GetActive(ctx context.Context) (*entityv1.TokenConfig, error) {
	var cfg entityv1.TokenConfig
	if err := r.db.WithContext(ctx).Table("authz_schema.token_configs").Where("is_active = true").First(&cfg).Error; err != nil {
		return nil, errors.New("tokenConfig.GetActive: " + err.Error())
	}
	return &cfg, nil
}

func (r *TokenConfigRepo) List(ctx context.Context) ([]*entityv1.TokenConfig, error) {
	var cfgs []*entityv1.TokenConfig
	if err := r.db.WithContext(ctx).Table("authz_schema.token_configs").Find(&cfgs).Error; err != nil {
		return nil, errors.New("tokenConfig.List: " + err.Error())
	}
	return cfgs, nil
}

func (r *TokenConfigRepo) Create(ctx context.Context, cfg *entityv1.TokenConfig) (*entityv1.TokenConfig, error) {
	if cfg == nil {
		return nil, errors.New("tokenConfig.Create: nil config")
	}
	// DB column is NOT NULL; ensure proto model always carries a timestamp.
	if cfg.CreatedAt == nil {
		cfg.CreatedAt = timestamppb.Now()
	}
	if cfg.Algorithm == "" {
		cfg.Algorithm = "RS256"
	}
	if err := r.db.WithContext(ctx).Table("authz_schema.token_configs").Create(cfg).Error; err != nil {
		return nil, errors.New("tokenConfig.Create: " + err.Error())
	}
	return cfg, nil
}
