package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type LookupRepository struct {
	db *gorm.DB
}

type PolicyLookup struct {
	PolicyID     string
	PolicyNumber string
	CustomerID   string
}

type ClaimLookup struct {
	ClaimID         string
	ClaimNumber     string
	PolicyID        string
	CustomerID      string
	RejectionReason string
}

type RenewalScheduleLookup struct {
	ScheduleID     string
	PolicyID       string
	RenewalDueDate *time.Time
}

type GracePeriodLookup struct {
	GracePeriodID string
	PolicyID      string
	EndDate       *time.Time
}

type OrderLookup struct {
	OrderID     string
	OrderNumber string
	CustomerID  string
	PolicyID    string
	PaymentID   string
}

type PartnerLookup struct {
	PartnerID        string
	OrganizationName string
	FocalPersonID    string
	ContactEmail     string
	Status           string
}

type AgentLookup struct {
	AgentID   string
	PartnerID string
	UserID    string
	FullName  string
	Email     string
	Phone     string
	Status    string
}

type FileLookup struct {
	FileID        string
	UploadedBy    string
	ReferenceID   string
	ReferenceType string
	Filename      string
}

type MediaLookup struct {
	MediaID    string
	FileID     string
	UploadedBy string
	EntityType string
	EntityID   string
	MediaType  string
	MimeType   string
}

type DocumentGenerationLookup struct {
	GenerationID       string
	DocumentTemplateID string
	EntityType         string
	EntityID           string
	GeneratedBy        string
	Status             string
	FileURL            string
}

func NewLookupRepository(db *gorm.DB) *LookupRepository {
	return &LookupRepository{db: db}
}

func (r *LookupRepository) GetPaymentPayer(ctx context.Context, paymentID string) (string, error) {
	var payerID string
	if err := r.db.WithContext(ctx).Raw(`
		SELECT payer_id
		FROM payment_schema.payments
		WHERE payment_id = $1`, paymentID).Row().Scan(&payerID); err != nil {
		if err == sql.ErrNoRows {
			return "", gorm.ErrRecordNotFound
		}
		return "", fmt.Errorf("get payment payer: %w", err)
	}
	return payerID, nil
}

func (r *LookupRepository) GetPolicy(ctx context.Context, policyID string) (*PolicyLookup, error) {
	var policy PolicyLookup
	if err := r.db.WithContext(ctx).Raw(`
		SELECT policy_id, policy_number, customer_id
		FROM insurance_schema.policies
		WHERE policy_id = $1 AND deleted_at IS NULL`, policyID).Row().Scan(
		&policy.PolicyID,
		&policy.PolicyNumber,
		&policy.CustomerID,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("get policy lookup: %w", err)
	}
	return &policy, nil
}

func (r *LookupRepository) GetClaim(ctx context.Context, claimID string) (*ClaimLookup, error) {
	var claim ClaimLookup
	if err := r.db.WithContext(ctx).Raw(`
		SELECT claim_id, claim_number, policy_id, customer_id, COALESCE(rejection_reason, '')
		FROM insurance_schema.claims
		WHERE claim_id = $1 AND deleted_at IS NULL`, claimID).Row().Scan(
		&claim.ClaimID,
		&claim.ClaimNumber,
		&claim.PolicyID,
		&claim.CustomerID,
		&claim.RejectionReason,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("get claim lookup: %w", err)
	}
	return &claim, nil
}

func (r *LookupRepository) GetRenewalSchedule(ctx context.Context, scheduleID string) (*RenewalScheduleLookup, error) {
	var (
		schedule RenewalScheduleLookup
		dueDate  sql.NullTime
	)
	if err := r.db.WithContext(ctx).Raw(`
		SELECT schedule_id, policy_id, renewal_due_date
		FROM insurance_schema.renewal_schedules
		WHERE schedule_id = $1`, scheduleID).Row().Scan(
		&schedule.ScheduleID,
		&schedule.PolicyID,
		&dueDate,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("get renewal schedule lookup: %w", err)
	}
	if dueDate.Valid {
		schedule.RenewalDueDate = &dueDate.Time
	}
	return &schedule, nil
}

func (r *LookupRepository) GetGracePeriod(ctx context.Context, gracePeriodID string) (*GracePeriodLookup, error) {
	var (
		gracePeriod GracePeriodLookup
		endDate     sql.NullTime
	)
	if err := r.db.WithContext(ctx).Raw(`
		SELECT grace_period_id, policy_id, end_date
		FROM insurance_schema.grace_periods
		WHERE grace_period_id = $1`, gracePeriodID).Row().Scan(
		&gracePeriod.GracePeriodID,
		&gracePeriod.PolicyID,
		&endDate,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("get grace period lookup: %w", err)
	}
	if endDate.Valid {
		gracePeriod.EndDate = &endDate.Time
	}
	return &gracePeriod, nil
}

func (r *LookupRepository) GetOrder(ctx context.Context, orderID string) (*OrderLookup, error) {
	var order OrderLookup
	if err := r.db.WithContext(ctx).Raw(`
		SELECT order_id, order_number, customer_id, COALESCE(policy_id::text, ''), COALESCE(payment_id::text, '')
		FROM insurance_schema.orders
		WHERE order_id = $1`,
		orderID,
	).Row().Scan(
		&order.OrderID,
		&order.OrderNumber,
		&order.CustomerID,
		&order.PolicyID,
		&order.PaymentID,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("get order lookup: %w", err)
	}
	return &order, nil
}

func (r *LookupRepository) GetPartner(ctx context.Context, partnerID string) (*PartnerLookup, error) {
	var partner PartnerLookup
	if err := r.db.WithContext(ctx).Raw(`
		SELECT partner_id, organization_name, COALESCE(focal_person_id::text, ''), COALESCE(contact_email, ''), COALESCE(status, '')
		FROM partner_schema.partners
		WHERE partner_id = $1 AND deleted_at IS NULL`,
		partnerID,
	).Row().Scan(
		&partner.PartnerID,
		&partner.OrganizationName,
		&partner.FocalPersonID,
		&partner.ContactEmail,
		&partner.Status,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("get partner lookup: %w", err)
	}
	return &partner, nil
}

func (r *LookupRepository) GetAgent(ctx context.Context, agentID string) (*AgentLookup, error) {
	var agent AgentLookup
	if err := r.db.WithContext(ctx).Raw(`
		SELECT agent_id, partner_id, COALESCE(user_id::text, ''), full_name, COALESCE(email, ''), COALESCE(phone_number, ''), COALESCE(status, '')
		FROM partner_schema.agents
		WHERE agent_id = $1 AND deleted_at IS NULL`,
		agentID,
	).Row().Scan(
		&agent.AgentID,
		&agent.PartnerID,
		&agent.UserID,
		&agent.FullName,
		&agent.Email,
		&agent.Phone,
		&agent.Status,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("get agent lookup: %w", err)
	}
	return &agent, nil
}

func (r *LookupRepository) GetFile(ctx context.Context, fileID string) (*FileLookup, error) {
	var file FileLookup
	if err := r.db.WithContext(ctx).Raw(`
		SELECT file_id, COALESCE(uploaded_by::text, ''), COALESCE(reference_id::text, ''), COALESCE(reference_type, ''), filename
		FROM storage_schema.files
		WHERE file_id = $1`,
		fileID,
	).Row().Scan(
		&file.FileID,
		&file.UploadedBy,
		&file.ReferenceID,
		&file.ReferenceType,
		&file.Filename,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("get file lookup: %w", err)
	}
	return &file, nil
}

func (r *LookupRepository) GetMedia(ctx context.Context, mediaID string) (*MediaLookup, error) {
	var media MediaLookup
	if err := r.db.WithContext(ctx).Raw(`
		SELECT media_id, COALESCE(file_id::text, ''), COALESCE(uploaded_by::text, ''), COALESCE(entity_type, ''), COALESCE(entity_id::text, ''), COALESCE(media_type, ''), COALESCE(mime_type, '')
		FROM media_schema.media_files
		WHERE media_id = $1`,
		mediaID,
	).Row().Scan(
		&media.MediaID,
		&media.FileID,
		&media.UploadedBy,
		&media.EntityType,
		&media.EntityID,
		&media.MediaType,
		&media.MimeType,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("get media lookup: %w", err)
	}
	return &media, nil
}

func (r *LookupRepository) GetDocumentGeneration(ctx context.Context, generationID string) (*DocumentGenerationLookup, error) {
	var generation DocumentGenerationLookup
	if err := r.db.WithContext(ctx).Raw(`
		SELECT generation_id, document_template_id, COALESCE(entity_type, ''), COALESCE(entity_id::text, ''), COALESCE(generated_by::text, ''), COALESCE(status, ''), COALESCE(file_url, '')
		FROM storage_schema.document_generations
		WHERE generation_id = $1`,
		generationID,
	).Row().Scan(
		&generation.GenerationID,
		&generation.DocumentTemplateID,
		&generation.EntityType,
		&generation.EntityID,
		&generation.GeneratedBy,
		&generation.Status,
		&generation.FileURL,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("get document generation lookup: %w", err)
	}
	return &generation, nil
}

func (r *LookupRepository) ListOrganisationAdminUserIDs(ctx context.Context, organisationID string) ([]string, error) {
	rows, err := r.db.WithContext(ctx).Raw(`
		SELECT DISTINCT user_id
		FROM b2b_schema.org_members
		WHERE organisation_id = $1
		  AND deleted_at IS NULL
		  AND status = 'ORG_MEMBER_STATUS_ACTIVE'
		  AND role IN ('ORG_MEMBER_ROLE_BUSINESS_ADMIN', 'ORG_MEMBER_ROLE_ADMIN')`,
		organisationID,
	).Rows()
	if err != nil {
		return nil, fmt.Errorf("list organisation admin user ids: %w", err)
	}
	defer rows.Close()

	userIDs := make([]string, 0)
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return nil, fmt.Errorf("scan organisation admin user id: %w", err)
		}
		userIDs = append(userIDs, userID)
	}
	if len(userIDs) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return userIDs, nil
}
