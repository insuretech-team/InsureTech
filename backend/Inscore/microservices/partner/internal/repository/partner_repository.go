package repository

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/crypto"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	partnerv1 "github.com/newage-saint/insuretech/gen/go/insuretech/partner/entity/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type PartnerRepository struct {
	db            *gorm.DB
	encryptionKey string
}

func NewPartnerRepository(db *gorm.DB, encryptionKey string) *PartnerRepository {
	return &PartnerRepository{db: db, encryptionKey: encryptionKey}
}

// decrypt decrypts a partner's PII fields in-place.
func (r *PartnerRepository) decrypt(partner *partnerv1.Partner) {
	if partner == nil {
		return
	}
	if partner.BankAccount != "" {
		if plain, err := crypto.DecryptPII(partner.BankAccount, r.encryptionKey); err == nil {
			partner.BankAccount = plain
		} else {
			logger.Warnf("failed to decrypt partner bank account: %v", err)
		}
	}
	if partner.ContactPhone != "" {
		if plain, err := crypto.DecryptPII(partner.ContactPhone, r.encryptionKey); err == nil {
			partner.ContactPhone = plain
		} else {
			logger.Warnf("failed to decrypt partner contact phone: %v", err)
		}
	}
}

// encrypt returns a new PII-encrypted map of fields for saving.
func (r *PartnerRepository) encrypt(partner *partnerv1.Partner) (map[string]any, error) {
	encBank := partner.BankAccount
	encPhone := partner.ContactPhone
	var err error

	if encBank != "" {
		encBank, err = crypto.EncryptPII(encBank, r.encryptionKey)
		if err != nil {
			return nil, err
		}
	}
	if encPhone != "" {
		encPhone, err = crypto.EncryptPII(encPhone, r.encryptionKey)
		if err != nil {
			return nil, err
		}
	}

	return map[string]any{
		"bank_account":  encBank,
		"contact_phone": encPhone,
	}, nil
}

// Create stores a new partner in the database
func (r *PartnerRepository) Create(ctx context.Context, partner *partnerv1.Partner) error {
	if partner.PartnerId == "" {
		partner.PartnerId = uuid.New().String()
	}
	now := time.Now()
	partner.CreatedAt = timestamppb.New(now)
	partner.UpdatedAt = timestamppb.New(now)

	if partner.Status == partnerv1.PartnerStatus_PARTNER_STATUS_UNSPECIFIED {
		partner.Status = partnerv1.PartnerStatus_PARTNER_STATUS_PENDING_VERIFICATION
	}

	encFields, err := r.encrypt(partner)
	if err != nil {
		return err
	}

	values := map[string]any{
		"partner_id":                  partner.PartnerId,
		"organization_name":           partner.OrganizationName,
		"type":                        partner.Type.String(),
		"status":                      partner.Status.String(),
		"trade_license":               partner.TradeLicense,
		"tin_number":                  partner.TinNumber,
		"bank_account":                encFields["bank_account"],
		"bank_name":                   partner.BankName,
		"bank_branch":                 partner.BankBranch,
		"contact_email":               partner.ContactEmail,
		"contact_phone":               encFields["contact_phone"],
		"acquisition_commission_rate": partner.AcquisitionCommissionRate,
		"renewal_commission_rate":     partner.RenewalCommissionRate,
		"claims_assistance_rate":      partner.ClaimsAssistanceRate,
		"focal_person_id":             nil,
		"created_at":                  now,
		"updated_at":                  now,
	}

	if partner.FocalPersonId != "" {
		values["focal_person_id"] = partner.FocalPersonId
	}

	return r.db.WithContext(ctx).Table("partner_schema.partners").Create(values).Error
}

// GetByID retrieves a partner by UUID
func (r *PartnerRepository) GetByID(ctx context.Context, id string) (*partnerv1.Partner, error) {
	var partner partnerv1.Partner
	err := r.db.WithContext(ctx).Table("partner_schema.partners").
		Where("partner_id = ? AND deleted_at IS NULL", id).
		First(&partner).Error
	if err != nil {
		return nil, err
	}
	r.decrypt(&partner)
	return &partner, nil
}

// GetByTradeLicense retrieves a partner by their unique trade license
func (r *PartnerRepository) GetByTradeLicense(ctx context.Context, tradeLicense string) (*partnerv1.Partner, error) {
	var partner partnerv1.Partner
	err := r.db.WithContext(ctx).Table("partner_schema.partners").
		Where("trade_license = ? AND deleted_at IS NULL", tradeLicense).
		First(&partner).Error
	if err != nil {
		return nil, err
	}
	r.decrypt(&partner)
	return &partner, nil
}

// UpdateStatus updates the verification and operational status of a partner
func (r *PartnerRepository) UpdateStatus(ctx context.Context, partnerID string, status partnerv1.PartnerStatus) error {
	upd := map[string]interface{}{
		"status":     status.String(),
		"updated_at": time.Now(),
	}

	if status == partnerv1.PartnerStatus_PARTNER_STATUS_ACTIVE {
		upd["onboarded_at"] = time.Now()
	}

	res := r.db.WithContext(ctx).Table("partner_schema.partners").
		Where("partner_id = ?", partnerID).
		Updates(upd)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// List retrieves partners with basic pagination
func (r *PartnerRepository) List(ctx context.Context, limit, offset int) ([]*partnerv1.Partner, error) {
	var partners []*partnerv1.Partner
	err := r.db.WithContext(ctx).Table("partner_schema.partners").
		Where("deleted_at IS NULL").
		Limit(limit).Offset(offset).
		Find(&partners).Error
	if err == nil {
		for _, p := range partners {
			r.decrypt(p)
		}
	}
	return partners, err
}

// ListWithFilters retrieves partners with pagination, simple filter, and order.
// Supported filter grammar (comma separated): status=<ENUM>,type=<ENUM>,organization_name=<text>
func (r *PartnerRepository) ListWithFilters(ctx context.Context, limit, offset int, filter, orderBy string) ([]*partnerv1.Partner, int32, error) {
	var (
		partners []*partnerv1.Partner
		total    int64
	)

	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	if offset < 0 {
		offset = 0
	}

	q := r.db.WithContext(ctx).Table("partner_schema.partners").Where("deleted_at IS NULL")
	q = applyPartnerFilter(q, filter)

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	orderExpr := normalizePartnerOrderBy(orderBy)
	if orderExpr != "" {
		q = q.Order(orderExpr)
	} else {
		q = q.Order("created_at DESC")
	}

	if err := q.Limit(limit).Offset(offset).Find(&partners).Error; err != nil {
		return nil, 0, err
	}
	for _, p := range partners {
		r.decrypt(p)
	}
	return partners, int32(total), nil
}

// Update updates partner fields; when updateMask is empty, all mutable fields are considered.
func (r *PartnerRepository) Update(ctx context.Context, partnerID string, partner *partnerv1.Partner, updateMask []string) error {
	if partner == nil {
		return gorm.ErrInvalidData
	}
	updates := make(map[string]any)
	now := time.Now()
	updates["updated_at"] = now

	mask := normalizeUpdateMask(updateMask)
	useAll := len(mask) == 0

	set := func(keys []string, apply func()) {
		if useAll {
			apply()
			return
		}
		for _, k := range keys {
			if _, ok := mask[k]; ok {
				apply()
				return
			}
		}
	}

	set([]string{"organization_name", "organizationName"}, func() { updates["organization_name"] = partner.OrganizationName })
	set([]string{"type"}, func() { updates["type"] = partner.Type.String() })
	set([]string{"status"}, func() { updates["status"] = partner.Status.String() })
	set([]string{"trade_license", "tradeLicense"}, func() { updates["trade_license"] = partner.TradeLicense })
	set([]string{"tin_number", "tinNumber"}, func() { updates["tin_number"] = partner.TinNumber })
	set([]string{"bank_name", "bankName"}, func() { updates["bank_name"] = partner.BankName })
	set([]string{"bank_branch", "bankBranch"}, func() { updates["bank_branch"] = partner.BankBranch })
	set([]string{"contact_email", "contactEmail"}, func() { updates["contact_email"] = partner.ContactEmail })
	set([]string{"acquisition_commission_rate", "acquisitionCommissionRate"}, func() {
		updates["acquisition_commission_rate"] = partner.AcquisitionCommissionRate
	})
	set([]string{"renewal_commission_rate", "renewalCommissionRate"}, func() {
		updates["renewal_commission_rate"] = partner.RenewalCommissionRate
	})
	set([]string{"claims_assistance_rate", "claimsAssistanceRate"}, func() {
		updates["claims_assistance_rate"] = partner.ClaimsAssistanceRate
	})
	set([]string{"focal_person_id", "focalPersonId"}, func() {
		if strings.TrimSpace(partner.FocalPersonId) == "" {
			updates["focal_person_id"] = nil
			return
		}
		updates["focal_person_id"] = partner.FocalPersonId
	})
	set([]string{"commission"}, func() { updates["commission"] = partner.Commission })
	set([]string{"benefits"}, func() { updates["benefits"] = partner.Benefits })

	// PII fields need encryption.
	set([]string{"bank_account", "bankAccount", "contact_phone", "contactPhone"}, func() {
		encFields, err := r.encrypt(partner)
		if err != nil {
			updates["__encrypt_error__"] = err
			return
		}
		if useAll || hasUpdateMaskKey(mask, "bank_account", "bankAccount") {
			updates["bank_account"] = encFields["bank_account"]
		}
		if useAll || hasUpdateMaskKey(mask, "contact_phone", "contactPhone") {
			updates["contact_phone"] = encFields["contact_phone"]
		}
	})

	if encErr, ok := updates["__encrypt_error__"]; ok {
		if err, okCast := encErr.(error); okCast {
			return err
		}
	}
	delete(updates, "__encrypt_error__")

	res := r.db.WithContext(ctx).Table("partner_schema.partners").
		Where("partner_id = ? AND deleted_at IS NULL", partnerID).
		Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// SoftDelete marks a partner as deleted.
func (r *PartnerRepository) SoftDelete(ctx context.Context, partnerID string) error {
	now := time.Now()
	res := r.db.WithContext(ctx).Table("partner_schema.partners").
		Where("partner_id = ? AND deleted_at IS NULL", partnerID).
		Updates(map[string]any{
			"deleted_at": now,
			"updated_at": now,
			"status":     partnerv1.PartnerStatus_PARTNER_STATUS_TERMINATED.String(),
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// CreateWithAgent stores a new partner and an initial agent (focal person) within a single transaction
func (r *PartnerRepository) CreateWithAgent(ctx context.Context, partner *partnerv1.Partner, agent *partnerv1.Agent) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Prepare and insert Partner
		if partner.PartnerId == "" {
			partner.PartnerId = uuid.New().String()
		}
		now := time.Now()
		partner.CreatedAt = timestamppb.New(now)
		partner.UpdatedAt = timestamppb.New(now)

		if partner.Status == partnerv1.PartnerStatus_PARTNER_STATUS_UNSPECIFIED {
			partner.Status = partnerv1.PartnerStatus_PARTNER_STATUS_PENDING_VERIFICATION
		}

		encFields, encErr := r.encrypt(partner)
		if encErr != nil {
			return encErr
		}

		partnerValues := map[string]any{
			"partner_id":                  partner.PartnerId,
			"organization_name":           partner.OrganizationName,
			"type":                        partner.Type.String(),
			"status":                      partner.Status.String(),
			"trade_license":               partner.TradeLicense,
			"tin_number":                  partner.TinNumber,
			"bank_account":                encFields["bank_account"],
			"bank_name":                   partner.BankName,
			"bank_branch":                 partner.BankBranch,
			"contact_email":               partner.ContactEmail,
			"contact_phone":               encFields["contact_phone"],
			"acquisition_commission_rate": partner.AcquisitionCommissionRate,
			"renewal_commission_rate":     partner.RenewalCommissionRate,
			"claims_assistance_rate":      partner.ClaimsAssistanceRate,
			"created_at":                  now,
			"updated_at":                  now,
		}

		if partner.Commission != nil {
			partnerValues["commission"] = partner.Commission
		}
		if partner.Benefits != nil {
			partnerValues["benefits"] = partner.Benefits
		}
		if partner.FocalPersonId != "" {
			partnerValues["focal_person_id"] = partner.FocalPersonId
		}

		if err := tx.Table("partner_schema.partners").Create(partnerValues).Error; err != nil {
			return err
		}

		// 2. Prepare and insert Agent
		if agent != nil {
			if agent.AgentId == "" {
				agent.AgentId = uuid.New().String()
			}
			agent.PartnerId = partner.PartnerId
			agent.CreatedAt = timestamppb.New(now)
			agent.UpdatedAt = timestamppb.New(now)
			agent.JoinedAt = timestamppb.New(now)

			if agent.Status == partnerv1.AgentStatus_AGENT_STATUS_UNSPECIFIED {
				agent.Status = partnerv1.AgentStatus_AGENT_STATUS_ACTIVE
			}

			encAgentPhone := agent.PhoneNumber
			if encAgentPhone != "" {
				if enc, err := crypto.EncryptPII(encAgentPhone, r.encryptionKey); err == nil {
					encAgentPhone = enc
				}
			}
			encAgentEmail := agent.Email
			if encAgentEmail != "" {
				if enc, err := crypto.EncryptPII(encAgentEmail, r.encryptionKey); err == nil {
					encAgentEmail = enc
				}
			}
			encAgentNid := agent.NidNumber
			if encAgentNid != "" {
				if enc, err := crypto.EncryptPII(encAgentNid, r.encryptionKey); err == nil {
					encAgentNid = enc
				}
			}

			agentValues := map[string]any{
				"agent_id":        agent.AgentId,
				"partner_id":      agent.PartnerId,
				"user_id":         agent.UserId,
				"full_name":       agent.FullName,
				"phone_number":    encAgentPhone,
				"email":           encAgentEmail,
				"nid_number":      encAgentNid,
				"status":          agent.Status.String(),
				"commission_rate": agent.CommissionRate,
				"joined_at":       now,
				"created_at":      now,
				"updated_at":      now,
			}

			if err := tx.Table("partner_schema.agents").Create(agentValues).Error; err != nil {
				return err
			}

			// Optional: Update focal person ID on partner if this agent represents them
			if partner.FocalPersonId == "" && agent.UserId != "" {
				if err := tx.Table("partner_schema.partners").
					Where("partner_id = ?", partner.PartnerId).
					Update("focal_person_id", agent.UserId).Error; err != nil {
					return err
				}
				partner.FocalPersonId = agent.UserId
			}
		}

		return nil
	})
}

func normalizeUpdateMask(updateMask []string) map[string]struct{} {
	out := make(map[string]struct{}, len(updateMask))
	for _, m := range updateMask {
		trimmed := strings.TrimSpace(m)
		if trimmed == "" {
			continue
		}
		out[trimmed] = struct{}{}
	}
	return out
}

func hasUpdateMaskKey(mask map[string]struct{}, keys ...string) bool {
	for _, k := range keys {
		if _, ok := mask[k]; ok {
			return true
		}
	}
	return false
}

func applyPartnerFilter(q *gorm.DB, filter string) *gorm.DB {
	filter = strings.TrimSpace(filter)
	if filter == "" {
		return q
	}
	parts := strings.Split(filter, ",")
	for _, part := range parts {
		kv := strings.SplitN(strings.TrimSpace(part), "=", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.TrimSpace(kv[0])
		val := strings.TrimSpace(kv[1])
		if val == "" {
			continue
		}
		switch key {
		case "status":
			q = q.Where("status = ?", normalizePartnerStatusFilter(val))
		case "type":
			q = q.Where("type = ?", normalizePartnerTypeFilter(val))
		case "organization_name":
			q = q.Where("organization_name ILIKE ?", "%"+val+"%")
		}
	}
	return q
}

func normalizePartnerOrderBy(orderBy string) string {
	orderBy = strings.TrimSpace(orderBy)
	if orderBy == "" {
		return ""
	}

	direction := "ASC"
	field := orderBy
	if strings.HasPrefix(field, "-") {
		direction = "DESC"
		field = strings.TrimPrefix(field, "-")
	} else if strings.Contains(strings.ToUpper(orderBy), " DESC") {
		direction = "DESC"
		field = strings.TrimSpace(strings.TrimSuffix(strings.TrimSuffix(orderBy, " DESC"), " desc"))
	} else if strings.Contains(strings.ToUpper(orderBy), " ASC") {
		field = strings.TrimSpace(strings.TrimSuffix(strings.TrimSuffix(orderBy, " ASC"), " asc"))
	}

	switch field {
	case "created_at", "updated_at", "organization_name", "status", "type":
		return field + " " + direction
	default:
		return ""
	}
}

func normalizePartnerStatusFilter(v string) string {
	if iv, ok := partnerv1.PartnerStatus_value[v]; ok {
		return partnerv1.PartnerStatus(iv).String()
	}
	if iv, err := strconv.Atoi(v); err == nil {
		if _, ok := partnerv1.PartnerStatus_name[int32(iv)]; ok {
			return partnerv1.PartnerStatus(iv).String()
		}
	}
	if _, ok := partnerv1.PartnerStatus_value["PARTNER_STATUS_"+v]; ok {
		return "PARTNER_STATUS_" + v
	}
	return v
}

func normalizePartnerTypeFilter(v string) string {
	if iv, ok := partnerv1.PartnerType_value[v]; ok {
		return partnerv1.PartnerType(iv).String()
	}
	if iv, err := strconv.Atoi(v); err == nil {
		if _, ok := partnerv1.PartnerType_name[int32(iv)]; ok {
			return partnerv1.PartnerType(iv).String()
		}
	}
	if _, ok := partnerv1.PartnerType_value["PARTNER_TYPE_"+v]; ok {
		return "PARTNER_TYPE_" + v
	}
	return v
}
