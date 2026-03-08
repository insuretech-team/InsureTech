package repository

import (
	"context"
	"strings"
	"time"

	authnv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authn/entity/v1"
	"gorm.io/gorm"
)

// UserProfileRepository provides access to authn_schema.user_profiles.
// NOTE: proto UserProfile uses user_id as primary key.
type UserProfileRepository struct{ db *gorm.DB }

func NewUserProfileRepository(db *gorm.DB) *UserProfileRepository {
	return &UserProfileRepository{db: db}
}

// genderDBValue returns the DB string for a proto Gender enum.
// The proto_enum serializer stores the name without the type prefix, e.g. "MALE".
// We replicate that here for raw-map inserts so the Scan round-trip works.
func genderDBValue(g authnv1.Gender) string {
	s := g.String() // e.g. "GENDER_MALE" or "GENDER_UNSPECIFIED"
	// Strip leading "GENDER_" prefix that the proto_enum serializer normally removes.
	if after, ok := strings.CutPrefix(s, "GENDER_"); ok {
		return after
	}
	return s
}

func (r *UserProfileRepository) Create(ctx context.Context, p *authnv1.UserProfile) error {
	now := time.Now().UTC()

	// date_of_birth is a PostgreSQL DATE column with:
	//   NOT NULL
	//   CHECK (date_of_birth <= CURRENT_DATE - '18 years')
	// Use a safe historical sentinel (1900-01-01) when caller supplies nil/zero.
	var dob time.Time
	if p.DateOfBirth != nil && p.DateOfBirth.IsValid() && p.DateOfBirth.AsTime().Year() > 1 {
		dob = p.DateOfBirth.AsTime().UTC()
	} else {
		dob = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	// nid_number is NOT NULL with CHECK (~ '^[0-9]{10}$|^[0-9]{13}$|^[0-9]{17}$').
	// Use supplied value if valid, otherwise fall back to a 10-zero placeholder so
	// the auto-create row passes the constraint. The user can update it later.
	nid := p.NidNumber
	if nid == "" {
		nid = "0000000000" // 10-digit placeholder — satisfies '^[0-9]{10}$'
	}

	// country defaults to 'Bangladesh' in DB; use supplied value or the default.
	country := p.Country
	if country == "" {
		country = "Bangladesh"
	}

	// Use a raw map so every NOT NULL column gets an explicit value, bypassing
	// GORM zero-value skipping and proto serializer edge cases.
	row := map[string]any{
		"user_id":                    p.UserId,
		"full_name":                  p.FullName,
		"date_of_birth":              dob,
		"gender":                     genderDBValue(p.Gender),
		"occupation":                 p.Occupation,
		"address_line1":              p.AddressLine1,
		"address_line2":              p.AddressLine2,
		"city":                       p.City,
		"district":                   p.District,
		"division":                   p.Division,
		"postal_code":                p.PostalCode,
		"country":                    country,
		"nid_number":                 nid,
		"kyc_verified":               p.KycVerified,
		"consent_privacy_acceptance": p.ConsentPrivacyAcceptance,
		"marital_status":             p.MaritalStatus,
		"employer":                   p.Employer,
		"created_at":                 now,
		"updated_at":                 now,
	}
	return r.db.WithContext(ctx).Table("authn_schema.user_profiles").Create(row).Error
}

func (r *UserProfileRepository) GetByUserID(ctx context.Context, userID string) (*authnv1.UserProfile, error) {
	var p authnv1.UserProfile
	err := r.db.WithContext(ctx).
		Table("authn_schema.user_profiles").
		Where("user_id = ?", userID).
		First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *UserProfileRepository) SetKYCVerified(ctx context.Context, userID string, verified bool, at *time.Time) error {
	upd := map[string]any{"kyc_verified": verified, "updated_at": time.Now()}
	if at != nil {
		upd["kyc_verified_at"] = *at
	}
	return r.db.WithContext(ctx).
		Table("authn_schema.user_profiles").
		Where("user_id = ?", userID).
		Updates(upd).Error
}

func (r *UserProfileRepository) Update(ctx context.Context, p *authnv1.UserProfile) error {
	// nid_number has CHECK (~ '^[0-9]{10}$|^[0-9]{13}$|^[0-9]{17}$') NOT NULL.
	// If caller sends empty string (profile not yet filled), use placeholder so
	// the constraint is not violated. User can provide real NID later.
	nid := p.NidNumber
	if nid == "" {
		nid = "0000000000"
	}

	// country NOT NULL — fall back to default if empty.
	country := p.Country
	if country == "" {
		country = "Bangladesh"
	}

	// Build a raw map of only the editable profile fields so GORM doesn't skip
	// zero values and proto serializers don't cause type mismatches.
	upd := map[string]any{
		"full_name":      p.FullName,
		"gender":         genderDBValue(p.Gender), // must be string, not int32
		"occupation":     p.Occupation,
		"address_line1":  p.AddressLine1,
		"address_line2":  p.AddressLine2,
		"city":           p.City,
		"district":       p.District,
		"division":       p.Division,
		"postal_code":    p.PostalCode,
		"country":        country,
		"nid_number":     nid,
		"marital_status": p.MaritalStatus,
		"employer":       p.Employer,
		"updated_at":     time.Now().UTC(),
	}
	// Only update date_of_birth if a real value was provided (not nil / sentinel).
	// The DB CHECK requires date_of_birth <= CURRENT_DATE - '18 years'.
	if p.DateOfBirth != nil && p.DateOfBirth.IsValid() && p.DateOfBirth.AsTime().Year() > 1 {
		upd["date_of_birth"] = p.DateOfBirth.AsTime().UTC()
	}
	return r.db.WithContext(ctx).
		Table("authn_schema.user_profiles").
		Where("user_id = ?", p.UserId).
		Updates(upd).Error
}

func (r *UserProfileRepository) DeleteByUserID(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).Table("authn_schema.user_profiles").Where("user_id = ?", userID).Delete(map[string]any{}).Error
}

// UserProfileRepoIface is the interface satisfied by UserProfileRepository.
type UserProfileRepoIface interface {
	Create(ctx context.Context, p *authnv1.UserProfile) error
	GetByUserID(ctx context.Context, userID string) (*authnv1.UserProfile, error)
	Update(ctx context.Context, p *authnv1.UserProfile) error
	SetKYCVerified(ctx context.Context, userID string, verified bool, at *time.Time) error
	DeleteByUserID(ctx context.Context, userID string) error
}
