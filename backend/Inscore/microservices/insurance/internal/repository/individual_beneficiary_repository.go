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

type IndividualBeneficiaryRepository struct {
	db *gorm.DB
}

func NewIndividualBeneficiaryRepository(db *gorm.DB) *IndividualBeneficiaryRepository {
	return &IndividualBeneficiaryRepository{db: db}
}

func (r *IndividualBeneficiaryRepository) Create(ctx context.Context, individual *beneficiaryv1.IndividualBeneficiary) (*beneficiaryv1.IndividualBeneficiary, error) {
	if individual.BeneficiaryId == "" {
		return nil, fmt.Errorf("beneficiary_id is required")
	}

	var createdBy string
	err := r.db.WithContext(ctx).Raw(`SELECT user_id FROM authn_schema.users LIMIT 1`).Scan(&createdBy).Error
	if err != nil || createdBy == "" {
		return nil, fmt.Errorf("failed to get valid user for created_by: %w", err)
	}

	auditInfoJSON := fmt.Sprintf(`{"created_by":"%s","created_at":"%s"}`, createdBy, time.Now().Format(time.RFC3339))

	var dateOfBirth sql.NullTime
	if individual.DateOfBirth != nil {
		dateOfBirth = sql.NullTime{Time: individual.DateOfBirth.AsTime(), Valid: true}
	}

	err = r.db.WithContext(ctx).Exec(`
		INSERT INTO insurance_schema.individual_beneficiaries
			(beneficiary_id, full_name, full_name_bn, date_of_birth, gender, nid_number,
			 passport_number, birth_certificate_number, tin_number, marital_status, occupation,
			 contact_info, permanent_address, present_address, nominee_name, nominee_relationship, audit_info)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)`,
		individual.BeneficiaryId,
		individual.FullName,
		individual.FullNameBn,
		dateOfBirth,
		strings.ToUpper(individual.Gender.String()),
		individual.NidNumber,
		individual.PassportNumber,
		individual.BirthCertificateNumber,
		individual.TinNumber,
		strings.ToUpper(individual.MaritalStatus.String()),
		individual.Occupation,
		"{}",
		"{}",
		"{}",
		individual.NomineeName,
		individual.NomineeRelationship,
		auditInfoJSON,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to insert individual beneficiary: %w", err)
	}

	return r.GetByID(ctx, individual.BeneficiaryId)
}

func (r *IndividualBeneficiaryRepository) GetByID(ctx context.Context, beneficiaryID string) (*beneficiaryv1.IndividualBeneficiary, error) {
	var (
		ind            beneficiaryv1.IndividualBeneficiary
		genderStr      sql.NullString
		maritalStr     sql.NullString
		dateOfBirth    sql.NullTime
		nidNumber      sql.NullString
		passportNumber sql.NullString
		birthCertNum   sql.NullString
		tinNumber      sql.NullString
		occupation     sql.NullString
		nomineeName    sql.NullString
		nomineeRel     sql.NullString
		contactInfo    sql.NullString
		permAddress    sql.NullString
		presAddress    sql.NullString
		auditInfo      sql.NullString
		createdAt      time.Time
		updatedAt      time.Time
	)

	err := r.db.WithContext(ctx).Raw(`
		SELECT beneficiary_id, full_name, COALESCE(full_name_bn, '') as full_name_bn,
		       date_of_birth, gender, nid_number, passport_number, birth_certificate_number,
		       tin_number, marital_status, occupation, contact_info, permanent_address,
		       present_address, nominee_name, nominee_relationship, audit_info,
		       created_at, updated_at
		FROM insurance_schema.individual_beneficiaries
		WHERE beneficiary_id = $1`,
		beneficiaryID,
	).Row().Scan(
		&ind.BeneficiaryId,
		&ind.FullName,
		&ind.FullNameBn,
		&dateOfBirth,
		&genderStr,
		&nidNumber,
		&passportNumber,
		&birthCertNum,
		&tinNumber,
		&maritalStr,
		&occupation,
		&contactInfo,
		&permAddress,
		&presAddress,
		&nomineeName,
		&nomineeRel,
		&auditInfo,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get individual beneficiary: %w", err)
	}

	if dateOfBirth.Valid {
		ind.DateOfBirth = timestamppb.New(dateOfBirth.Time)
	}

	if genderStr.Valid {
		k := strings.ToUpper(genderStr.String)
		if v, ok := beneficiaryv1.Gender_value[k]; ok {
			ind.Gender = beneficiaryv1.Gender(v)
		}
	}

	if maritalStr.Valid {
		k := strings.ToUpper(maritalStr.String)
		if v, ok := beneficiaryv1.MaritalStatus_value[k]; ok {
			ind.MaritalStatus = beneficiaryv1.MaritalStatus(v)
		}
	}

	if nidNumber.Valid {
		ind.NidNumber = nidNumber.String
	}
	if passportNumber.Valid {
		ind.PassportNumber = passportNumber.String
	}
	if birthCertNum.Valid {
		ind.BirthCertificateNumber = birthCertNum.String
	}
	if tinNumber.Valid {
		ind.TinNumber = tinNumber.String
	}
	if occupation.Valid {
		ind.Occupation = occupation.String
	}
	if nomineeName.Valid {
		ind.NomineeName = nomineeName.String
	}
	if nomineeRel.Valid {
		ind.NomineeRelationship = nomineeRel.String
	}

	if contactInfo.Valid {
		ind.ContactInfo = &commonv1.ContactInfo{}
	}
	if permAddress.Valid {
		ind.PermanentAddress = &commonv1.Address{}
	}
	if presAddress.Valid {
		ind.PresentAddress = &commonv1.Address{}
	}
	if auditInfo.Valid {
		ind.AuditInfo = &commonv1.AuditInfo{}
	}

	return &ind, nil
}

func (r *IndividualBeneficiaryRepository) Update(ctx context.Context, individual *beneficiaryv1.IndividualBeneficiary) (*beneficiaryv1.IndividualBeneficiary, error) {
	var dateOfBirth sql.NullTime
	if individual.DateOfBirth != nil {
		dateOfBirth = sql.NullTime{Time: individual.DateOfBirth.AsTime(), Valid: true}
	}

	err := r.db.WithContext(ctx).Exec(`
		UPDATE insurance_schema.individual_beneficiaries
		SET full_name = $2,
		    full_name_bn = $3,
		    date_of_birth = $4,
		    gender = $5,
		    nid_number = $6,
		    passport_number = $7,
		    birth_certificate_number = $8,
		    tin_number = $9,
		    marital_status = $10,
		    occupation = $11,
		    nominee_name = $12,
		    nominee_relationship = $13,
		    updated_at = NOW()
		WHERE beneficiary_id = $1`,
		individual.BeneficiaryId,
		individual.FullName,
		individual.FullNameBn,
		dateOfBirth,
		strings.ToUpper(individual.Gender.String()),
		individual.NidNumber,
		individual.PassportNumber,
		individual.BirthCertificateNumber,
		individual.TinNumber,
		strings.ToUpper(individual.MaritalStatus.String()),
		individual.Occupation,
		individual.NomineeName,
		individual.NomineeRelationship,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to update individual beneficiary: %w", err)
	}

	return r.GetByID(ctx, individual.BeneficiaryId)
}

func (r *IndividualBeneficiaryRepository) Delete(ctx context.Context, beneficiaryID string) error {
	err := r.db.WithContext(ctx).Exec(`
		DELETE FROM insurance_schema.individual_beneficiaries
		WHERE beneficiary_id = $1`,
		beneficiaryID,
	).Error

	if err != nil {
		return fmt.Errorf("failed to delete individual beneficiary: %w", err)
	}

	return nil
}
