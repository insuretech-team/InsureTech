package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"google.golang.org/protobuf/types/known/timestamppb"

	beneficiaryv1 "github.com/newage-saint/insuretech/gen/go/insuretech/beneficiary/entity/v1"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
)

type BusinessBeneficiaryRepository struct {
	db *gorm.DB
}

func NewBusinessBeneficiaryRepository(db *gorm.DB) *BusinessBeneficiaryRepository {
	return &BusinessBeneficiaryRepository{db: db}
}

func (r *BusinessBeneficiaryRepository) Create(ctx context.Context, business *beneficiaryv1.BusinessBeneficiary) (*beneficiaryv1.BusinessBeneficiary, error) {
	if business.Id == "" {
		return nil, fmt.Errorf("beneficiary_id is required")
	}

	var createdBy string
	err := r.db.WithContext(ctx).Raw(`SELECT user_id FROM authn_schema.users LIMIT 1`).Scan(&createdBy).Error
	if err != nil || createdBy == "" {
		return nil, fmt.Errorf("failed to get valid user for created_by: %w", err)
	}

	auditInfoJSON := fmt.Sprintf(`{"created_by":"%s","created_at":"%s"}`, createdBy, time.Now().Format(time.RFC3339))

	// Handle Money type
	totalPremiumAmount := int64(0)
	totalPremiumCurrency := "BDT"
	if business.TotalPremiumAmount != nil {
		totalPremiumAmount = business.TotalPremiumAmount.Amount
		totalPremiumCurrency = business.TotalPremiumAmount.Currency
	}

	var tradeLicenseIssue, tradeLicenseExpiry, incorporationDate sql.NullTime
	if business.TradeLicenseIssueDate != nil {
		tradeLicenseIssue = sql.NullTime{Time: business.TradeLicenseIssueDate.AsTime(), Valid: true}
	}
	if business.TradeLicenseExpiryDate != nil {
		tradeLicenseExpiry = sql.NullTime{Time: business.TradeLicenseExpiryDate.AsTime(), Valid: true}
	}
	if business.IncorporationDate != nil {
		incorporationDate = sql.NullTime{Time: business.IncorporationDate.AsTime(), Valid: true}
	}

	err = r.db.WithContext(ctx).Exec(`
		INSERT INTO insurance_schema.business_beneficiaries
			(beneficiary_id, parent_beneficiary_id, business_name, business_name_bn, trade_license_number,
			 trade_license_issue_date, trade_license_expiry_date, tin_number, bin_number, business_type,
			 industry_sector, employee_count, incorporation_date, contact_info, registered_address,
			 business_address, focal_person_name, focal_person_designation, focal_person_nid,
			 focal_person_contact, audit_info, registration_number, tax_id, primary_contact,
			 total_employees_covered, active_policies_count, total_premium_amount, total_premium_currency, pending_actions_count)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29)`,
		business.Id,
		business.BeneficiaryId,
		business.BusinessName,
		business.BusinessNameBn,
		business.TradeLicenseNumber,
		tradeLicenseIssue,
		tradeLicenseExpiry,
		business.TinNumber,
		business.BinNumber,
		strings.ToUpper(business.BusinessType.String()),
		business.IndustrySector,
		business.EmployeeCount,
		incorporationDate,
		"{}",
		"{}",
		"{}",
		business.FocalPersonName,
		business.FocalPersonDesignation,
		business.FocalPersonNid,
		"{}",
		auditInfoJSON,
		business.RegistrationNumber,
		business.TaxId,
		"{}",
		business.TotalEmployeesCovered,
		business.ActivePoliciesCount,
		totalPremiumAmount,
		totalPremiumCurrency,
		business.PendingActionsCount,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to insert business beneficiary: %w", err)
	}

	return r.GetByID(ctx, business.Id)
}

func (r *BusinessBeneficiaryRepository) GetByID(ctx context.Context, beneficiaryID string) (*beneficiaryv1.BusinessBeneficiary, error) {
	var (
		bus                    beneficiaryv1.BusinessBeneficiary
		businessTypeStr        sql.NullString
		tradeLicenseIssue      sql.NullTime
		tradeLicenseExpiry     sql.NullTime
		incorporationDate      sql.NullTime
		binNumber              sql.NullString
		industrySector         sql.NullString
		employeeCount          sql.NullInt64
		focalPersonDesignation sql.NullString
		focalPersonNid         sql.NullString
		registrationNumber     sql.NullString
		taxId                  sql.NullString
		totalPremiumAmount     int64
		contactInfo            sql.NullString
		registeredAddress      sql.NullString
		businessAddress        sql.NullString
		focalPersonContact     sql.NullString
		primaryContact         sql.NullString
		auditInfo              sql.NullString
		createdAt              time.Time
		updatedAt              time.Time
	)

	err := r.db.WithContext(ctx).Raw(`
		SELECT beneficiary_id, parent_beneficiary_id, business_name, COALESCE(business_name_bn, '') as business_name_bn,
		       trade_license_number, trade_license_issue_date, trade_license_expiry_date,
		       tin_number, bin_number, business_type, industry_sector, employee_count,
		       incorporation_date, contact_info, registered_address, business_address,
		       focal_person_name, focal_person_designation, focal_person_nid, focal_person_contact,
		       audit_info, registration_number, tax_id, primary_contact,
		       total_employees_covered, active_policies_count, total_premium_amount, pending_actions_count,
		       created_at, updated_at
		FROM insurance_schema.business_beneficiaries
		WHERE beneficiary_id = $1`,
		beneficiaryID,
	).Row().Scan(
		&bus.Id,
		&bus.BeneficiaryId,
		&bus.BusinessName,
		&bus.BusinessNameBn,
		&bus.TradeLicenseNumber,
		&tradeLicenseIssue,
		&tradeLicenseExpiry,
		&bus.TinNumber,
		&binNumber,
		&businessTypeStr,
		&industrySector,
		&employeeCount,
		&incorporationDate,
		&contactInfo,
		&registeredAddress,
		&businessAddress,
		&bus.FocalPersonName,
		&focalPersonDesignation,
		&focalPersonNid,
		&focalPersonContact,
		&auditInfo,
		&registrationNumber,
		&taxId,
		&primaryContact,
		&bus.TotalEmployeesCovered,
		&bus.ActivePoliciesCount,
		&totalPremiumAmount,
		&bus.PendingActionsCount,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get business beneficiary: %w", err)
	}

	if businessTypeStr.Valid {
		k := strings.ToUpper(businessTypeStr.String)
		if v, ok := beneficiaryv1.BusinessType_value[k]; ok {
			bus.BusinessType = beneficiaryv1.BusinessType(v)
		}
	}

	if tradeLicenseIssue.Valid {
		bus.TradeLicenseIssueDate = timestamppb.New(tradeLicenseIssue.Time)
	}
	if tradeLicenseExpiry.Valid {
		bus.TradeLicenseExpiryDate = timestamppb.New(tradeLicenseExpiry.Time)
	}
	if incorporationDate.Valid {
		bus.IncorporationDate = timestamppb.New(incorporationDate.Time)
	}

	if binNumber.Valid {
		bus.BinNumber = binNumber.String
	}
	if industrySector.Valid {
		bus.IndustrySector = industrySector.String
	}
	if employeeCount.Valid {
		bus.EmployeeCount = int32(employeeCount.Int64)
	}
	if focalPersonDesignation.Valid {
		bus.FocalPersonDesignation = focalPersonDesignation.String
	}
	if focalPersonNid.Valid {
		bus.FocalPersonNid = focalPersonNid.String
	}
	if registrationNumber.Valid {
		bus.RegistrationNumber = registrationNumber.String
	}
	if taxId.Valid {
		bus.TaxId = taxId.String
	}

	bus.TotalPremiumAmount = &commonv1.Money{
		Amount:   totalPremiumAmount,
		Currency: "BDT",
	}

	if contactInfo.Valid {
		bus.ContactInfo = &commonv1.ContactInfo{}
	}
	if registeredAddress.Valid {
		bus.RegisteredAddress = &commonv1.Address{}
	}
	if businessAddress.Valid {
		bus.BusinessAddress = &commonv1.Address{}
	}
	if focalPersonContact.Valid {
		bus.FocalPersonContact = &commonv1.ContactInfo{}
	}
	if primaryContact.Valid {
		bus.PrimaryContact = &beneficiaryv1.PrimaryContact{}
	}
	if auditInfo.Valid {
		bus.AuditInfo = &commonv1.AuditInfo{}
	}

	return &bus, nil
}

func (r *BusinessBeneficiaryRepository) Update(ctx context.Context, business *beneficiaryv1.BusinessBeneficiary) (*beneficiaryv1.BusinessBeneficiary, error) {
	totalPremiumAmount := int64(0)
	if business.TotalPremiumAmount != nil {
		totalPremiumAmount = business.TotalPremiumAmount.Amount
	}

	var tradeLicenseIssue, tradeLicenseExpiry, incorporationDate sql.NullTime
	if business.TradeLicenseIssueDate != nil {
		tradeLicenseIssue = sql.NullTime{Time: business.TradeLicenseIssueDate.AsTime(), Valid: true}
	}
	if business.TradeLicenseExpiryDate != nil {
		tradeLicenseExpiry = sql.NullTime{Time: business.TradeLicenseExpiryDate.AsTime(), Valid: true}
	}
	if business.IncorporationDate != nil {
		incorporationDate = sql.NullTime{Time: business.IncorporationDate.AsTime(), Valid: true}
	}

	err := r.db.WithContext(ctx).Exec(`
		UPDATE insurance_schema.business_beneficiaries
		SET parent_beneficiary_id = $2,
		    business_name = $3,
		    business_name_bn = $4,
		    trade_license_number = $5,
		    trade_license_issue_date = $6,
		    trade_license_expiry_date = $7,
		    tin_number = $8,
		    bin_number = $9,
		    business_type = $10,
		    industry_sector = $11,
		    employee_count = $12,
		    incorporation_date = $13,
		    focal_person_name = $14,
		    focal_person_designation = $15,
		    focal_person_nid = $16,
		    registration_number = $17,
		    tax_id = $18,
		    total_employees_covered = $19,
		    active_policies_count = $20,
		    total_premium_amount = $21,
		    pending_actions_count = $22,
		    updated_at = NOW()
		WHERE beneficiary_id = $1`,
		business.Id,
		business.BeneficiaryId,
		business.BusinessName,
		business.BusinessNameBn,
		business.TradeLicenseNumber,
		tradeLicenseIssue,
		tradeLicenseExpiry,
		business.TinNumber,
		business.BinNumber,
		strings.ToUpper(business.BusinessType.String()),
		business.IndustrySector,
		business.EmployeeCount,
		incorporationDate,
		business.FocalPersonName,
		business.FocalPersonDesignation,
		business.FocalPersonNid,
		business.RegistrationNumber,
		business.TaxId,
		business.TotalEmployeesCovered,
		business.ActivePoliciesCount,
		totalPremiumAmount,
		business.PendingActionsCount,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to update business beneficiary: %w", err)
	}

	return r.GetByID(ctx, business.Id)
}

func (r *BusinessBeneficiaryRepository) Delete(ctx context.Context, beneficiaryID string) error {
	err := r.db.WithContext(ctx).Exec(`
		DELETE FROM insurance_schema.business_beneficiaries
		WHERE beneficiary_id = $1`,
		beneficiaryID,
	).Error

	if err != nil {
		return fmt.Errorf("failed to delete business beneficiary: %w", err)
	}

	return nil
}
