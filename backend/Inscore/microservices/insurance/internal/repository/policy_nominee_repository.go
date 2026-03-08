package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/gorm"
	"google.golang.org/protobuf/types/known/timestamppb"

	policyv1 "github.com/newage-saint/insuretech/gen/go/insuretech/policy/entity/v1"
)

type PolicyNomineeRepository struct {
	db *gorm.DB
}

func NewPolicyNomineeRepository(db *gorm.DB) *PolicyNomineeRepository {
	return &PolicyNomineeRepository{db: db}
}

func (r *PolicyNomineeRepository) Create(ctx context.Context, nominee *policyv1.Nominee) (*policyv1.Nominee, error) {
	if nominee.NomineeId == "" {
		return nil, fmt.Errorf("nominee_id is required")
	}
	
	if nominee.PolicyId == "" {
		return nil, fmt.Errorf("policy_id is required")
	}
	
	// Handle date_of_birth
	var dateOfBirth interface{}
	if nominee.DateOfBirth != nil {
		dateOfBirth = nominee.DateOfBirth.AsTime()
	}
	
	err := r.db.WithContext(ctx).Exec(`
		INSERT INTO insurance_schema.policy_nominees
			(nominee_id, policy_id, full_name, relationship, share_percentage,
			 date_of_birth, nid_number, phone_number, nominee_dob_text, nominee_share_percent,
			 created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())`,
		nominee.NomineeId,
		nominee.PolicyId,
		nominee.FullName,
		nominee.Relationship,
		nominee.SharePercentage,
		dateOfBirth,
		nominee.NidNumber,
		nominee.PhoneNumber,
		nominee.NomineeDobText,
		nominee.NomineeSharePercent,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to insert policy nominee: %w", err)
	}

	return r.GetByID(ctx, nominee.NomineeId)
}

func (r *PolicyNomineeRepository) GetByID(ctx context.Context, nomineeID string) (*policyv1.Nominee, error) {
	var (
		n             policyv1.Nominee
		dateOfBirth   time.Time
		nidNumber     sql.NullString
		phoneNumber   sql.NullString
		nomineeDobText sql.NullString
		nomineeSharePercent sql.NullFloat64
		createdAt     time.Time
		updatedAt     time.Time
	)

	err := r.db.WithContext(ctx).Raw(`
		SELECT nominee_id, policy_id, full_name, relationship, share_percentage,
		       date_of_birth, nid_number, phone_number, nominee_dob_text, nominee_share_percent,
		       created_at, updated_at
		FROM insurance_schema.policy_nominees
		WHERE nominee_id = $1`,
		nomineeID,
	).Row().Scan(
		&n.NomineeId,
		&n.PolicyId,
		&n.FullName,
		&n.Relationship,
		&n.SharePercentage,
		&dateOfBirth,
		&nidNumber,
		&phoneNumber,
		&nomineeDobText,
		&nomineeSharePercent,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get policy nominee: %w", err)
	}

	if !dateOfBirth.IsZero() {
		n.DateOfBirth = timestamppb.New(dateOfBirth)
	}
	if nidNumber.Valid {
		n.NidNumber = nidNumber.String
	}
	if phoneNumber.Valid {
		n.PhoneNumber = phoneNumber.String
	}
	if nomineeDobText.Valid {
		n.NomineeDobText = nomineeDobText.String
	}
	if nomineeSharePercent.Valid {
		n.NomineeSharePercent = nomineeSharePercent.Float64
	}
	if !createdAt.IsZero() {
		n.CreatedAt = timestamppb.New(createdAt)
	}
	if !updatedAt.IsZero() {
		n.UpdatedAt = timestamppb.New(updatedAt)
	}

	return &n, nil
}

func (r *PolicyNomineeRepository) Update(ctx context.Context, nominee *policyv1.Nominee) (*policyv1.Nominee, error) {
	// Handle date_of_birth
	var dateOfBirth interface{}
	if nominee.DateOfBirth != nil {
		dateOfBirth = nominee.DateOfBirth.AsTime()
	}
	
	err := r.db.WithContext(ctx).Exec(`
		UPDATE insurance_schema.policy_nominees
		SET policy_id = $2,
		    full_name = $3,
		    relationship = $4,
		    share_percentage = $5,
		    date_of_birth = $6,
		    nid_number = $7,
		    phone_number = $8,
		    nominee_dob_text = $9,
		    nominee_share_percent = $10,
		    updated_at = NOW()
		WHERE nominee_id = $1`,
		nominee.NomineeId,
		nominee.PolicyId,
		nominee.FullName,
		nominee.Relationship,
		nominee.SharePercentage,
		dateOfBirth,
		nominee.NidNumber,
		nominee.PhoneNumber,
		nominee.NomineeDobText,
		nominee.NomineeSharePercent,
	).Error

	if err != nil {
		return nil, fmt.Errorf("failed to update policy nominee: %w", err)
	}

	return r.GetByID(ctx, nominee.NomineeId)
}

func (r *PolicyNomineeRepository) Delete(ctx context.Context, nomineeID string) error {
	err := r.db.WithContext(ctx).Exec(`
		DELETE FROM insurance_schema.policy_nominees
		WHERE nominee_id = $1`,
		nomineeID,
	).Error

	if err != nil {
		return fmt.Errorf("failed to delete policy nominee: %w", err)
	}

	return nil
}

func (r *PolicyNomineeRepository) ListByPolicyID(ctx context.Context, policyID string) ([]*policyv1.Nominee, error) {
	rows, err := r.db.WithContext(ctx).Raw(`
		SELECT nominee_id, policy_id, full_name, relationship, share_percentage,
		       date_of_birth, nid_number, phone_number, nominee_dob_text, nominee_share_percent,
		       created_at, updated_at
		FROM insurance_schema.policy_nominees
		WHERE policy_id = $1
		ORDER BY created_at DESC`,
		policyID,
	).Rows()

	if err != nil {
		return nil, fmt.Errorf("failed to list policy nominees: %w", err)
	}
	defer rows.Close()

	nominees := make([]*policyv1.Nominee, 0)
	for rows.Next() {
		var (
			n             policyv1.Nominee
			dateOfBirth   time.Time
			nidNumber     sql.NullString
			phoneNumber   sql.NullString
			nomineeDobText sql.NullString
			nomineeSharePercent sql.NullFloat64
			createdAt     time.Time
			updatedAt     time.Time
		)

		err := rows.Scan(
			&n.NomineeId,
			&n.PolicyId,
			&n.FullName,
			&n.Relationship,
			&n.SharePercentage,
			&dateOfBirth,
			&nidNumber,
			&phoneNumber,
			&nomineeDobText,
			&nomineeSharePercent,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan policy nominee: %w", err)
		}

		if !dateOfBirth.IsZero() {
			n.DateOfBirth = timestamppb.New(dateOfBirth)
		}
		if nidNumber.Valid {
			n.NidNumber = nidNumber.String
		}
		if phoneNumber.Valid {
			n.PhoneNumber = phoneNumber.String
		}
		if nomineeDobText.Valid {
			n.NomineeDobText = nomineeDobText.String
		}
		if nomineeSharePercent.Valid {
			n.NomineeSharePercent = nomineeSharePercent.Float64
		}
		if !createdAt.IsZero() {
			n.CreatedAt = timestamppb.New(createdAt)
		}
		if !updatedAt.IsZero() {
			n.UpdatedAt = timestamppb.New(updatedAt)
		}

		nominees = append(nominees, &n)
	}

	return nominees, nil
}
