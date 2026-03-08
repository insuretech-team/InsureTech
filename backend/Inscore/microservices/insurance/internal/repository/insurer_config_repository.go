package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"gorm.io/gorm"

	insurerv1 "github.com/newage-saint/insuretech/gen/go/insuretech/insurer/entity/v1"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
)

type InsurerConfigRepository struct {
	db *gorm.DB
}

func NewInsurerConfigRepository(db *gorm.DB) *InsurerConfigRepository {
	return &InsurerConfigRepository{db: db}
}

func (r *InsurerConfigRepository) Create(ctx context.Context, config *insurerv1.InsurerConfig) (*insurerv1.InsurerConfig, error) {
	if config.Id == "" {
		return nil, fmt.Errorf("config_id is required")
	}

	// Handle JSONB fields
	var businessModel interface{}
	if config.BusinessModel != "" {
		businessModel = config.BusinessModel
	}
	var paymentTerms interface{}
	if config.PaymentTerms != "" {
		paymentTerms = config.PaymentTerms
	}
	var auditInfo interface{}
	if config.AuditInfo != nil {
		auditInfo = "{}"
	}

	err := r.db.WithContext(ctx).Exec(`
		INSERT INTO insurance_schema.insurer_configs
			(config_id, insurer_id, api_base_url, api_version, auth_type, auth_credentials,
			 webhook_url, webhook_secret, business_model, auto_underwriting_enabled,
			 underwriting_threshold, real_time_claim_notification, claim_settlement_days,
			 payment_terms, audit_info)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`,
		config.Id,
		config.InsurerId,
		config.ApiBaseUrl,
		config.ApiVersion,
		strings.ToUpper(config.AuthType.String()),
		config.AuthCredentials,
		config.WebhookUrl,
		config.WebhookSecret,
		businessModel,
		config.AutoUnderwritingEnabled,
		config.UnderwritingThreshold,
		config.RealTimeClaimNotification,
		config.ClaimSettlementDays,
		paymentTerms,
		auditInfo,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to insert insurer config: %w", err)
	}

	return r.GetByID(ctx, config.Id)
}

func (r *InsurerConfigRepository) GetByID(ctx context.Context, configID string) (*insurerv1.InsurerConfig, error) {
	var (
		cfg                        insurerv1.InsurerConfig
		apiBaseUrl                 sql.NullString
		apiVersion                 sql.NullString
		authTypeStr                sql.NullString
		authCredentials            sql.NullString
		webhookUrl                 sql.NullString
		webhookSecret              sql.NullString
		businessModel              sql.NullString
		paymentTerms               sql.NullString
		auditInfo                  sql.NullString
	)

	err := r.db.WithContext(ctx).Raw(`
		SELECT config_id, insurer_id, api_base_url, api_version, auth_type, auth_credentials,
		       webhook_url, webhook_secret, business_model, auto_underwriting_enabled,
		       underwriting_threshold, real_time_claim_notification, claim_settlement_days,
		       payment_terms, audit_info
		FROM insurance_schema.insurer_configs
		WHERE config_id = $1`,
		configID,
	).Row().Scan(
		&cfg.Id,
		&cfg.InsurerId,
		&apiBaseUrl,
		&apiVersion,
		&authTypeStr,
		&authCredentials,
		&webhookUrl,
		&webhookSecret,
		&businessModel,
		&cfg.AutoUnderwritingEnabled,
		&cfg.UnderwritingThreshold,
		&cfg.RealTimeClaimNotification,
		&cfg.ClaimSettlementDays,
		&paymentTerms,
		&auditInfo,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get insurer config: %w", err)
	}

	// Set optional fields
	if apiBaseUrl.Valid {
		cfg.ApiBaseUrl = apiBaseUrl.String
	}
	if apiVersion.Valid {
		cfg.ApiVersion = apiVersion.String
	}
	if authCredentials.Valid {
		cfg.AuthCredentials = authCredentials.String
	}
	if webhookUrl.Valid {
		cfg.WebhookUrl = webhookUrl.String
	}
	if webhookSecret.Valid {
		cfg.WebhookSecret = webhookSecret.String
	}
	if businessModel.Valid {
		cfg.BusinessModel = businessModel.String
	}
	if paymentTerms.Valid {
		cfg.PaymentTerms = paymentTerms.String
	}

	// Parse enum
	if authTypeStr.Valid {
		k := strings.ToUpper(authTypeStr.String)
		if v, ok := insurerv1.AuthenticationType_value[k]; ok {
			cfg.AuthType = insurerv1.AuthenticationType(v)
		}
	}

	// Set audit info
	if auditInfo.Valid {
		cfg.AuditInfo = &commonv1.AuditInfo{}
	}

	return &cfg, nil
}

func (r *InsurerConfigRepository) Update(ctx context.Context, config *insurerv1.InsurerConfig) (*insurerv1.InsurerConfig, error) {
	// Handle JSONB fields
	var businessModel interface{}
	if config.BusinessModel != "" {
		businessModel = config.BusinessModel
	}
	var paymentTerms interface{}
	if config.PaymentTerms != "" {
		paymentTerms = config.PaymentTerms
	}

	err := r.db.WithContext(ctx).Exec(`
		UPDATE insurance_schema.insurer_configs
		SET insurer_id = $2,
		    api_base_url = $3,
		    api_version = $4,
		    auth_type = $5,
		    auth_credentials = $6,
		    webhook_url = $7,
		    webhook_secret = $8,
		    business_model = $9,
		    auto_underwriting_enabled = $10,
		    underwriting_threshold = $11,
		    real_time_claim_notification = $12,
		    claim_settlement_days = $13,
		    payment_terms = $14
		WHERE config_id = $1`,
		config.Id,
		config.InsurerId,
		config.ApiBaseUrl,
		config.ApiVersion,
		strings.ToUpper(config.AuthType.String()),
		config.AuthCredentials,
		config.WebhookUrl,
		config.WebhookSecret,
		businessModel,
		config.AutoUnderwritingEnabled,
		config.UnderwritingThreshold,
		config.RealTimeClaimNotification,
		config.ClaimSettlementDays,
		paymentTerms,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to update insurer config: %w", err)
	}

	return r.GetByID(ctx, config.Id)
}

func (r *InsurerConfigRepository) Delete(ctx context.Context, configID string) error {
	err := r.db.WithContext(ctx).Exec(`
		DELETE FROM insurance_schema.insurer_configs
		WHERE config_id = $1`,
		configID,
	).Error

	if err != nil {
		return fmt.Errorf("failed to delete insurer config: %w", err)
	}

	return nil
}

func (r *InsurerConfigRepository) GetByInsurerID(ctx context.Context, insurerID string) (*insurerv1.InsurerConfig, error) {
	var (
		cfg                        insurerv1.InsurerConfig
		apiBaseUrl                 sql.NullString
		apiVersion                 sql.NullString
		authTypeStr                sql.NullString
		authCredentials            sql.NullString
		webhookUrl                 sql.NullString
		webhookSecret              sql.NullString
		businessModel              sql.NullString
		paymentTerms               sql.NullString
		auditInfo                  sql.NullString
	)

	err := r.db.WithContext(ctx).Raw(`
		SELECT config_id, insurer_id, api_base_url, api_version, auth_type, auth_credentials,
		       webhook_url, webhook_secret, business_model, auto_underwriting_enabled,
		       underwriting_threshold, real_time_claim_notification, claim_settlement_days,
		       payment_terms, audit_info
		FROM insurance_schema.insurer_configs
		WHERE insurer_id = $1`,
		insurerID,
	).Row().Scan(
		&cfg.Id,
		&cfg.InsurerId,
		&apiBaseUrl,
		&apiVersion,
		&authTypeStr,
		&authCredentials,
		&webhookUrl,
		&webhookSecret,
		&businessModel,
		&cfg.AutoUnderwritingEnabled,
		&cfg.UnderwritingThreshold,
		&cfg.RealTimeClaimNotification,
		&cfg.ClaimSettlementDays,
		&paymentTerms,
		&auditInfo,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get insurer config by insurer_id: %w", err)
	}

	// Set optional fields
	if apiBaseUrl.Valid {
		cfg.ApiBaseUrl = apiBaseUrl.String
	}
	if apiVersion.Valid {
		cfg.ApiVersion = apiVersion.String
	}
	if authCredentials.Valid {
		cfg.AuthCredentials = authCredentials.String
	}
	if webhookUrl.Valid {
		cfg.WebhookUrl = webhookUrl.String
	}
	if webhookSecret.Valid {
		cfg.WebhookSecret = webhookSecret.String
	}
	if businessModel.Valid {
		cfg.BusinessModel = businessModel.String
	}
	if paymentTerms.Valid {
		cfg.PaymentTerms = paymentTerms.String
	}

	// Parse enum
	if authTypeStr.Valid {
		k := strings.ToUpper(authTypeStr.String)
		if v, ok := insurerv1.AuthenticationType_value[k]; ok {
			cfg.AuthType = insurerv1.AuthenticationType(v)
		}
	}

	// Set audit info
	if auditInfo.Valid {
		cfg.AuditInfo = &commonv1.AuditInfo{}
	}

	return &cfg, nil
}
