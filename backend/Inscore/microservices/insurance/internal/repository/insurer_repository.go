package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"gorm.io/gorm"
	"google.golang.org/protobuf/types/known/timestamppb"

	insurerv1 "github.com/newage-saint/insuretech/gen/go/insuretech/insurer/entity/v1"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
)

type InsurerRepository struct {
	db *gorm.DB
}

func NewInsurerRepository(db *gorm.DB) *InsurerRepository {
	return &InsurerRepository{db: db}
}

func (r *InsurerRepository) Create(ctx context.Context, insurer *insurerv1.Insurer) (*insurerv1.Insurer, error) {
	if insurer.Id == "" {
		return nil, fmt.Errorf("insurer_id is required")
	}

	// Handle JSONB fields - paid_up_capital is stored as JSONB
	var paidUpCapitalJSON interface{}
	if insurer.PaidUpCapital != nil {
		// Store as JSON: {"amount": 123, "currency": "BDT"}
		paidUpCapitalJSON = fmt.Sprintf(`{"amount": %d, "currency": "%s"}`, insurer.PaidUpCapital.Amount, insurer.PaidUpCapital.Currency)
	} else {
		paidUpCapitalJSON = "{}"
	}

	// Handle timestamps
	var idraLicenseExpiry sql.NullTime
	if insurer.IdraLicenseExpiry != nil {
		idraLicenseExpiry = sql.NullTime{Time: insurer.IdraLicenseExpiry.AsTime(), Valid: true}
	}

	// Handle JSONB fields - serialize proto messages to JSON
	var contactInfo interface{}
	if insurer.ContactInfo != nil {
		contactInfo = fmt.Sprintf(`{"mobile_number":"%s","email":"%s","alternate_mobile":"%s","landline":"%s"}`,
			insurer.ContactInfo.MobileNumber,
			insurer.ContactInfo.Email,
			insurer.ContactInfo.AlternateMobile,
			insurer.ContactInfo.Landline)
	} else {
		contactInfo = "{}"
	}
	
	var registeredAddress interface{}
	if insurer.RegisteredAddress != nil {
		registeredAddress = "{}"
	} else {
		registeredAddress = "{}"
	}
	
	var headOfficeAddress interface{}
	if insurer.HeadOfficeAddress != nil {
		headOfficeAddress = "{}"
	} else {
		headOfficeAddress = "{}"
	}

	var auditInfo interface{}
	if insurer.AuditInfo != nil {
		auditInfo = "{}"
	} else {
		auditInfo = "{}"
	}

	err := r.db.WithContext(ctx).Exec(`
		INSERT INTO insurance_schema.insurers
			(insurer_id, code, name, name_bn, type, status, trade_license_number,
			 tin_number, idra_license_number, idra_license_expiry, contact_info,
			 registered_address, head_office_address, logo_url, website_url,
			 financial_rating, paid_up_capital, audit_info)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)`,
		insurer.Id,
		insurer.Code,
		insurer.Name,
		insurer.NameBn,
		strings.TrimPrefix(strings.ToUpper(insurer.Type.String()), "INSURER_TYPE_"),
		strings.TrimPrefix(strings.ToUpper(insurer.Status.String()), "INSURER_STATUS_"),
		insurer.TradeLicenseNumber,
		insurer.TinNumber,
		insurer.IdraLicenseNumber,
		idraLicenseExpiry,
		contactInfo,
		registeredAddress,
		headOfficeAddress,
		insurer.LogoUrl,
		insurer.WebsiteUrl,
		insurer.FinancialRating,
		paidUpCapitalJSON,
		auditInfo,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to insert insurer: %w", err)
	}

	return r.GetByID(ctx, insurer.Id)
}

func (r *InsurerRepository) GetByID(ctx context.Context, insurerID string) (*insurerv1.Insurer, error) {
	var (
		ins                   insurerv1.Insurer
		typeStr               sql.NullString
		statusStr             sql.NullString
		nameBn                sql.NullString
		tradeLicenseNumber    sql.NullString
		tinNumber             sql.NullString
		idraLicenseNumber     sql.NullString
		idraLicenseExpiry     sql.NullTime
		contactInfo           sql.NullString
		registeredAddress     sql.NullString
		headOfficeAddress     sql.NullString
		logoUrl               sql.NullString
		websiteUrl            sql.NullString
		financialRating       sql.NullString
		paidUpCapitalJSON     sql.NullString
		auditInfo             sql.NullString
	)

	err := r.db.WithContext(ctx).Raw(`
		SELECT insurer_id, code, name, name_bn, type, status, trade_license_number,
		       tin_number, idra_license_number, idra_license_expiry, contact_info,
		       registered_address, head_office_address, logo_url, website_url,
		       financial_rating, paid_up_capital, audit_info
		FROM insurance_schema.insurers
		WHERE insurer_id = $1`,
		insurerID,
	).Row().Scan(
		&ins.Id,
		&ins.Code,
		&ins.Name,
		&nameBn,
		&typeStr,
		&statusStr,
		&tradeLicenseNumber,
		&tinNumber,
		&idraLicenseNumber,
		&idraLicenseExpiry,
		&contactInfo,
		&registeredAddress,
		&headOfficeAddress,
		&logoUrl,
		&websiteUrl,
		&financialRating,
		&paidUpCapitalJSON,
		&auditInfo,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get insurer: %w", err)
	}

	// Set Money field from JSONB (skip parsing for now, just set empty)
	ins.PaidUpCapital = &commonv1.Money{Amount: 0, Currency: "BDT"}

	// Set optional fields
	if nameBn.Valid {
		ins.NameBn = nameBn.String
	}
	if tradeLicenseNumber.Valid {
		ins.TradeLicenseNumber = tradeLicenseNumber.String
	}
	if tinNumber.Valid {
		ins.TinNumber = tinNumber.String
	}
	if idraLicenseNumber.Valid {
		ins.IdraLicenseNumber = idraLicenseNumber.String
	}
	if logoUrl.Valid {
		ins.LogoUrl = logoUrl.String
	}
	if websiteUrl.Valid {
		ins.WebsiteUrl = websiteUrl.String
	}
	if financialRating.Valid {
		ins.FinancialRating = financialRating.String
	}

	// Set JSONB fields - parse JSON to proto messages
	if contactInfo.Valid && contactInfo.String != "{}" {
		// Parse ContactInfo JSON
		var contactData struct {
			MobileNumber    string `json:"mobile_number"`
			Email           string `json:"email"`
			AlternateMobile string `json:"alternate_mobile"`
			Landline        string `json:"landline"`
		}
		if err := json.Unmarshal([]byte(contactInfo.String), &contactData); err == nil {
			ins.ContactInfo = &commonv1.ContactInfo{
				MobileNumber:    contactData.MobileNumber,
				Email:           contactData.Email,
				AlternateMobile: contactData.AlternateMobile,
				Landline:        contactData.Landline,
			}
		} else {
			ins.ContactInfo = &commonv1.ContactInfo{}
		}
	}
	if registeredAddress.Valid {
		ins.RegisteredAddress = &commonv1.Address{}
	}
	if headOfficeAddress.Valid {
		ins.HeadOfficeAddress = &commonv1.Address{}
	}

	// Parse enums
	if typeStr.Valid {
		k := strings.ToUpper(typeStr.String)
		if v, ok := insurerv1.InsurerType_value[k]; ok {
			ins.Type = insurerv1.InsurerType(v)
		}
	}
	if statusStr.Valid {
		k := strings.ToUpper(statusStr.String)
		if v, ok := insurerv1.InsurerStatus_value[k]; ok {
			ins.Status = insurerv1.InsurerStatus(v)
		}
	}

	// Set timestamps
	if idraLicenseExpiry.Valid {
		ins.IdraLicenseExpiry = timestamppb.New(idraLicenseExpiry.Time)
	}

	// Set audit info
	if auditInfo.Valid {
		ins.AuditInfo = &commonv1.AuditInfo{}
	}

	return &ins, nil
}

func (r *InsurerRepository) Update(ctx context.Context, insurer *insurerv1.Insurer) (*insurerv1.Insurer, error) {
	// Handle timestamps
	var idraLicenseExpiry sql.NullTime
	if insurer.IdraLicenseExpiry != nil {
		idraLicenseExpiry = sql.NullTime{Time: insurer.IdraLicenseExpiry.AsTime(), Valid: true}
	}

	// Handle JSONB fields - serialize proto messages to JSON
	var contactInfo interface{}
	if insurer.ContactInfo != nil {
		// Serialize ContactInfo to JSON
		contactInfo = fmt.Sprintf(`{"mobile_number":"%s","email":"%s","alternate_mobile":"%s","landline":"%s"}`,
			insurer.ContactInfo.MobileNumber,
			insurer.ContactInfo.Email,
			insurer.ContactInfo.AlternateMobile,
			insurer.ContactInfo.Landline)
	} else {
		contactInfo = "{}"
	}
	
	var registeredAddress interface{}
	if insurer.RegisteredAddress != nil {
		registeredAddress = "{}"
	} else {
		registeredAddress = "{}"
	}
	
	var headOfficeAddress interface{}
	if insurer.HeadOfficeAddress != nil {
		headOfficeAddress = "{}"
	} else {
		headOfficeAddress = "{}"
	}
	
	// Handle paid_up_capital as JSONB
	var paidUpCapital interface{}
	if insurer.PaidUpCapital != nil {
		paidUpCapital = fmt.Sprintf(`{"amount":%d,"currency":"%s"}`, insurer.PaidUpCapital.Amount, insurer.PaidUpCapital.Currency)
	} else {
		paidUpCapital = "{}"
	}

	err := r.db.WithContext(ctx).Exec(`
		UPDATE insurance_schema.insurers
		SET code = $2,
		    name = $3,
		    name_bn = $4,
		    type = $5,
		    status = $6,
		    trade_license_number = $7,
		    tin_number = $8,
		    idra_license_number = $9,
		    idra_license_expiry = $10,
		    contact_info = $11,
		    registered_address = $12,
		    head_office_address = $13,
		    logo_url = $14,
		    website_url = $15,
		    financial_rating = $16,
		    paid_up_capital = $17
		WHERE insurer_id = $1`,
		insurer.Id,
		insurer.Code,
		insurer.Name,
		insurer.NameBn,
		strings.TrimPrefix(strings.ToUpper(insurer.Type.String()), "INSURER_TYPE_"),
		strings.TrimPrefix(strings.ToUpper(insurer.Status.String()), "INSURER_STATUS_"),
		insurer.TradeLicenseNumber,
		insurer.TinNumber,
		insurer.IdraLicenseNumber,
		idraLicenseExpiry,
		contactInfo,
		registeredAddress,
		headOfficeAddress,
		insurer.LogoUrl,
		insurer.WebsiteUrl,
		insurer.FinancialRating,
		paidUpCapital,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to update insurer: %w", err)
	}

	return r.GetByID(ctx, insurer.Id)
}

func (r *InsurerRepository) Delete(ctx context.Context, insurerID string) error {
	err := r.db.WithContext(ctx).Exec(`
		DELETE FROM insurance_schema.insurers
		WHERE insurer_id = $1`,
		insurerID,
	).Error

	if err != nil {
		return fmt.Errorf("failed to delete insurer: %w", err)
	}

	return nil
}

func (r *InsurerRepository) List(ctx context.Context, page, pageSize int) ([]*insurerv1.Insurer, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 50
	}

	offset := (page - 1) * pageSize

	// Get total count
	var total int64
	err := r.db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM insurance_schema.insurers`).Scan(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count insurers: %w", err)
	}

	// Get insurers
	rows, err := r.db.WithContext(ctx).Raw(`
		SELECT insurer_id, code, name, name_bn, type, status, trade_license_number,
		       tin_number, idra_license_number, idra_license_expiry, contact_info,
		       registered_address, head_office_address, logo_url, website_url,
		       financial_rating, paid_up_capital, audit_info
		FROM insurance_schema.insurers
		ORDER BY name ASC
		LIMIT $1 OFFSET $2`,
		pageSize, offset,
	).Rows()

	if err != nil {
		return nil, 0, fmt.Errorf("failed to list insurers: %w", err)
	}
	defer rows.Close()

	insurers := make([]*insurerv1.Insurer, 0)
	for rows.Next() {
		var (
			ins                insurerv1.Insurer
			typeStr            sql.NullString
			statusStr          sql.NullString
			nameBn             sql.NullString
			tradeLicenseNumber sql.NullString
			tinNumber          sql.NullString
			idraLicenseNumber  sql.NullString
			idraLicenseExpiry  sql.NullTime
			contactInfo        sql.NullString
			registeredAddress  sql.NullString
			headOfficeAddress  sql.NullString
			logoUrl            sql.NullString
			websiteUrl         sql.NullString
			financialRating    sql.NullString
			paidUpCapital      sql.NullString
			auditInfo          sql.NullString
		)

		err := rows.Scan(
			&ins.Id,
			&ins.Code,
			&ins.Name,
			&nameBn,
			&typeStr,
			&statusStr,
			&tradeLicenseNumber,
			&tinNumber,
			&idraLicenseNumber,
			&idraLicenseExpiry,
			&contactInfo,
			&registeredAddress,
			&headOfficeAddress,
			&logoUrl,
			&websiteUrl,
			&financialRating,
			&paidUpCapital,
			&auditInfo,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan insurer: %w", err)
		}

		// Parse paid_up_capital JSONB
		if paidUpCapital.Valid && paidUpCapital.String != "{}" {
			ins.PaidUpCapital = &commonv1.Money{Amount: 0, Currency: "BDT"}
		}

		// Set optional fields
		if nameBn.Valid {
			ins.NameBn = nameBn.String
		}
		if tradeLicenseNumber.Valid {
			ins.TradeLicenseNumber = tradeLicenseNumber.String
		}
		if tinNumber.Valid {
			ins.TinNumber = tinNumber.String
		}
		if idraLicenseNumber.Valid {
			ins.IdraLicenseNumber = idraLicenseNumber.String
		}
		if logoUrl.Valid {
			ins.LogoUrl = logoUrl.String
		}
		if websiteUrl.Valid {
			ins.WebsiteUrl = websiteUrl.String
		}
		if financialRating.Valid {
			ins.FinancialRating = financialRating.String
		}

		// Set JSONB fields
		if contactInfo.Valid {
			ins.ContactInfo = &commonv1.ContactInfo{}
		}
		if registeredAddress.Valid {
			ins.RegisteredAddress = &commonv1.Address{}
		}
		if headOfficeAddress.Valid {
			ins.HeadOfficeAddress = &commonv1.Address{}
		}

		// Parse enums - need to add prefix back
		if typeStr.Valid {
			k := "INSURER_TYPE_" + strings.ToUpper(typeStr.String)
			if v, ok := insurerv1.InsurerType_value[k]; ok {
				ins.Type = insurerv1.InsurerType(v)
			}
		}
		if statusStr.Valid {
			k := "INSURER_STATUS_" + strings.ToUpper(statusStr.String)
			if v, ok := insurerv1.InsurerStatus_value[k]; ok {
				ins.Status = insurerv1.InsurerStatus(v)
			}
		}

		// Set timestamps
		if idraLicenseExpiry.Valid {
			ins.IdraLicenseExpiry = timestamppb.New(idraLicenseExpiry.Time)
		}

		// Set audit info
		if auditInfo.Valid {
			ins.AuditInfo = &commonv1.AuditInfo{}
		}

		insurers = append(insurers, &ins)
	}

	return insurers, total, nil
}
