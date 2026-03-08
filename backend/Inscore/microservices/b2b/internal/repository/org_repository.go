package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/b2b/internal/domain"
	b2bv1 "github.com/newage-saint/insuretech/gen/go/insuretech/b2b/entity/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

// ─── SQL ──────────────────────────────────────────────────────────────────────

const organisationCols = `
	organisation_id,
	tenant_id,
	name,
	code,
	COALESCE(industry, '') AS industry,
	COALESCE(contact_email, '') AS contact_email,
	COALESCE(contact_phone, '') AS contact_phone,
	COALESCE(address, '') AS address,
	status,
	COALESCE(total_employees, 0) AS total_employees,
	created_at,
	updated_at
`

const orgMemberCols = `
	member_id,
	organisation_id,
	user_id,
	role,
	status,
	joined_at,
	created_at,
	updated_at
`

// ─── Scanners ─────────────────────────────────────────────────────────────────

func scanOrganisation(row interface{ Scan(...any) error }) (*b2bv1.Organisation, error) {
	var (
		o              b2bv1.Organisation
		statusStr      sql.NullString
		totalEmployees sql.NullInt32
		createdAt      time.Time
		updatedAt      time.Time
	)

	if err := row.Scan(
		&o.OrganisationId,
		&o.TenantId,
		&o.Name,
		&o.Code,
		&o.Industry,
		&o.ContactEmail,
		&o.ContactPhone,
		&o.Address,
		&statusStr,
		&totalEmployees,
		&createdAt,
		&updatedAt,
	); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	if statusStr.Valid {
		k := strings.ToUpper(statusStr.String)
		if v, ok := b2bv1.OrganisationStatus_value[k]; ok {
			o.Status = b2bv1.OrganisationStatus(v)
		} else if v, ok := b2bv1.OrganisationStatus_value["ORGANISATION_STATUS_"+k]; ok {
			o.Status = b2bv1.OrganisationStatus(v)
		}
	}
	o.TotalEmployees = totalEmployees.Int32
	if !createdAt.IsZero() { o.CreatedAt = timestamppb.New(createdAt) }
	if !updatedAt.IsZero() { o.UpdatedAt = timestamppb.New(updatedAt) }
	return &o, nil
}

func scanOrgMember(row interface{ Scan(...any) error }) (*b2bv1.OrgMember, error) {
	var (
		m          b2bv1.OrgMember
		roleStr    sql.NullString
		statusStr  sql.NullString
		joinedAt   time.Time
		createdAt  time.Time
		updatedAt  time.Time
	)

	if err := row.Scan(
		&m.MemberId,
		&m.OrganisationId,
		&m.UserId,
		&roleStr,
		&statusStr,
		&joinedAt,
		&createdAt,
		&updatedAt,
	); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	if roleStr.Valid {
		k := strings.ToUpper(roleStr.String)
		if v, ok := b2bv1.OrgMemberRole_value[k]; ok {
			m.Role = b2bv1.OrgMemberRole(v)
		} else if v, ok := b2bv1.OrgMemberRole_value["ORG_MEMBER_ROLE_"+k]; ok {
			m.Role = b2bv1.OrgMemberRole(v)
		}
	}
	if statusStr.Valid {
		k := strings.ToUpper(statusStr.String)
		if v, ok := b2bv1.OrgMemberStatus_value[k]; ok {
			m.Status = b2bv1.OrgMemberStatus(v)
		} else if v, ok := b2bv1.OrgMemberStatus_value["ORG_MEMBER_STATUS_"+k]; ok {
			m.Status = b2bv1.OrgMemberStatus(v)
		}
	}
	if !joinedAt.IsZero()  { m.JoinedAt = timestamppb.New(joinedAt) }
	if !createdAt.IsZero() { m.CreatedAt = timestamppb.New(createdAt) }
	if !updatedAt.IsZero() { m.UpdatedAt = timestamppb.New(updatedAt) }
	return &m, nil
}

// ─── Organisation CRUD ────────────────────────────────────────────────────────

func (r *PortalRepository) GetOrganisation(ctx context.Context, organisationID string) (*b2bv1.Organisation, error) {
	query := fmt.Sprintf(
		`SELECT %s FROM b2b_schema.organisations WHERE organisation_id = $1 AND deleted_at IS NULL LIMIT 1`,
		organisationCols,
	)
	row := r.db.WithContext(ctx).Raw(query, organisationID).Row()
	return scanOrganisation(row)
}

func (r *PortalRepository) ListOrganisations(
	ctx context.Context,
	pageSize, offset int,
	tenantID string,
	status b2bv1.OrganisationStatus,
) ([]*b2bv1.Organisation, int64, error) {
	where := "deleted_at IS NULL"
	args := []interface{}{}
	idx := 1

	if tenantID != "" {
		where += fmt.Sprintf(" AND tenant_id = $%d", idx)
		args = append(args, tenantID)
		idx++
	}
	if status != b2bv1.OrganisationStatus_ORGANISATION_STATUS_UNSPECIFIED {
		where += fmt.Sprintf(" AND status = $%d", idx)
		args = append(args, organisationStatusStr(status))
		idx++
	}

	var total int64
	countArgs := make([]interface{}, len(args))
	copy(countArgs, args)
	if err := r.db.WithContext(ctx).Raw(
		fmt.Sprintf("SELECT COUNT(*) FROM b2b_schema.organisations WHERE %s", where),
		countArgs...,
	).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(
		`SELECT %s FROM b2b_schema.organisations WHERE %s ORDER BY name ASC LIMIT $%d OFFSET $%d`,
		organisationCols, where, idx, idx+1,
	)
	args = append(args, pageSize, offset)

	rows, err := r.db.WithContext(ctx).Raw(query, args...).Rows()
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var orgs []*b2bv1.Organisation
	for rows.Next() {
		o, err := scanOrganisation(rows)
		if err != nil {
			return nil, 0, err
		}
		orgs = append(orgs, o)
	}
	return orgs, total, rows.Err()
}

func (r *PortalRepository) CreateOrganisation(ctx context.Context, input domain.OrganisationCreateInput) (*b2bv1.Organisation, error) {
	id := input.OrganisationID
	if id == "" {
		id = newUUID()
	}

	if err := r.db.WithContext(ctx).Exec(`
		INSERT INTO b2b_schema.organisations
			(organisation_id, tenant_id, name, code, industry, contact_email, contact_phone, address, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		id,
		input.TenantID,
		input.Name,
		input.Code,
		nullableStr(input.Industry),
		nullableStr(input.ContactEmail),
		nullableStr(input.ContactPhone),
		nullableStr(input.Address),
		"ORGANISATION_STATUS_ACTIVE",
	).Error; err != nil {
		return nil, fmt.Errorf("insert organisation: %w", err)
	}

	return r.GetOrganisation(ctx, id)
}

func (r *PortalRepository) UpdateOrganisation(ctx context.Context, input domain.OrganisationUpdateInput) (*b2bv1.Organisation, error) {
	setClauses := []string{}
	args := []interface{}{}
	idx := 1

	addStr := func(col, val string) {
		if val != "" {
			setClauses = append(setClauses, fmt.Sprintf("%s = $%d", col, idx))
			args = append(args, val)
			idx++
		}
	}

	addStr("name", input.Name)
	addStr("industry", input.Industry)
	addStr("contact_email", input.ContactEmail)
	addStr("contact_phone", input.ContactPhone)
	addStr("address", input.Address)
	if input.Status != b2bv1.OrganisationStatus_ORGANISATION_STATUS_UNSPECIFIED {
		setClauses = append(setClauses, fmt.Sprintf("status = $%d", idx))
		args = append(args, organisationStatusStr(input.Status))
		idx++
	}
	if len(setClauses) == 0 {
		return r.GetOrganisation(ctx, input.OrganisationID)
	}

	setClauses = append(setClauses, "updated_at = NOW()")
	query := fmt.Sprintf(
		"UPDATE b2b_schema.organisations SET %s WHERE organisation_id = $%d",
		strings.Join(setClauses, ", "), idx,
	)
	args = append(args, input.OrganisationID)

	if err := r.db.WithContext(ctx).Exec(query, args...).Error; err != nil {
		return nil, fmt.Errorf("update organisation: %w", err)
	}
	return r.GetOrganisation(ctx, input.OrganisationID)
}

func (r *PortalRepository) DeleteOrganisation(ctx context.Context, organisationID string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(`
			UPDATE b2b_schema.org_members
			   SET status = $2, deleted_at = NOW(), updated_at = NOW()
			 WHERE organisation_id = $1
			   AND deleted_at IS NULL`,
			organisationID,
			"ORG_MEMBER_STATUS_INACTIVE",
		).Error; err != nil {
			return fmt.Errorf("deactivate org members: %w", err)
		}

		result := tx.Exec(`
			UPDATE b2b_schema.organisations
			   SET status = $2, deleted_at = NOW(), updated_at = NOW()
			 WHERE organisation_id = $1
			   AND deleted_at IS NULL`,
			organisationID,
			"ORGANISATION_STATUS_INACTIVE",
		)
		if result.Error != nil {
			return fmt.Errorf("delete organisation: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}

// ─── OrgMember CRUD ───────────────────────────────────────────────────────────

func (r *PortalRepository) ListOrgMembers(ctx context.Context, organisationID string) ([]*b2bv1.OrgMember, error) {
	rows, err := r.db.WithContext(ctx).Raw(
		fmt.Sprintf(
			`SELECT %s
			   FROM b2b_schema.org_members
			  WHERE organisation_id = $1
			    AND deleted_at IS NULL
			  ORDER BY joined_at ASC, created_at ASC`,
			orgMemberCols,
		),
		organisationID,
	).Rows()
	if err != nil {
		return nil, fmt.Errorf("list org members: %w", err)
	}
	defer rows.Close()

	var members []*b2bv1.OrgMember
	for rows.Next() {
		member, scanErr := scanOrgMember(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		members = append(members, member)
	}
	return members, rows.Err()
}

func (r *PortalRepository) AddOrgMember(ctx context.Context, input domain.OrgMemberCreateInput) (*b2bv1.OrgMember, error) {
	id := input.MemberID
	if id == "" {
		id = newUUID()
	}

	if err := r.db.WithContext(ctx).Exec(`
		INSERT INTO b2b_schema.org_members (member_id, organisation_id, user_id, role, status, joined_at)
		VALUES ($1, $2, $3, $4, $5, NOW())`,
		id,
		input.OrganisationID,
		input.UserID,
		orgMemberRoleStr(input.Role),
		"ORG_MEMBER_STATUS_ACTIVE",
	).Error; err != nil {
		return nil, fmt.Errorf("insert org_member: %w", err)
	}

	query := fmt.Sprintf(
		`SELECT %s FROM b2b_schema.org_members WHERE member_id = $1 LIMIT 1`,
		orgMemberCols,
	)
	row := r.db.WithContext(ctx).Raw(query, id).Row()
	return scanOrgMember(row)
}

func (r *PortalRepository) AssignOrgAdmin(ctx context.Context, organisationID, memberID string) (*b2bv1.OrgMember, error) {
	result := r.db.WithContext(ctx).Exec(`
		UPDATE b2b_schema.org_members
		   SET role = $3, updated_at = NOW()
		 WHERE organisation_id = $1
		   AND member_id = $2
		   AND deleted_at IS NULL`,
		organisationID,
		memberID,
		orgMemberRoleStr(b2bv1.OrgMemberRole_ORG_MEMBER_ROLE_BUSINESS_ADMIN),
	)
	if result.Error != nil {
		return nil, fmt.Errorf("assign org admin: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	row := r.db.WithContext(ctx).Raw(
		fmt.Sprintf(
			`SELECT %s FROM b2b_schema.org_members WHERE member_id = $1 AND deleted_at IS NULL LIMIT 1`,
			orgMemberCols,
		),
		memberID,
	).Row()
	return scanOrgMember(row)
}

func (r *PortalRepository) RemoveOrgMember(ctx context.Context, organisationID, memberID string) error {
	return r.db.WithContext(ctx).Exec(
		`DELETE FROM b2b_schema.org_members
		 WHERE organisation_id = $1 AND member_id = $2`,
		organisationID, memberID,
	).Error
}

// ResolveOrganisationByUserID is the key query that fixes the hardcoded business_id problem.
// It joins org_members → organisations to return the active organisation for a given user.
func (r *PortalRepository) ResolveOrganisationByUserID(
	ctx context.Context,
	userID string,
) (organisationID string, role b2bv1.OrgMemberRole, organisationName string, err error) {
	type resolveRow struct {
		OrganisationID   string `gorm:"column:organisation_id"`
		OrganisationName string `gorm:"column:organisation_name"`
		Role             string `gorm:"column:role"`
	}

	var row resolveRow
	err = r.db.WithContext(ctx).Raw(`
		SELECT om.organisation_id, o.name AS organisation_name, om.role
		FROM b2b_schema.org_members om
		JOIN b2b_schema.organisations o
		  ON o.organisation_id = om.organisation_id
		WHERE om.user_id = $1
		  AND om.status = 'ORG_MEMBER_STATUS_ACTIVE'
		ORDER BY om.joined_at ASC
		LIMIT 1`,
		userID,
	).Scan(&row).Error

	if err != nil {
		return "", b2bv1.OrgMemberRole_ORG_MEMBER_ROLE_UNSPECIFIED, "", fmt.Errorf("resolve organisation for user %s: %w", userID, err)
	}
	if row.OrganisationID == "" {
		return "", b2bv1.OrgMemberRole_ORG_MEMBER_ROLE_UNSPECIFIED, "", gorm.ErrRecordNotFound
	}

	k := strings.ToUpper(row.Role)
	var parsedRole b2bv1.OrgMemberRole
	if v, ok := b2bv1.OrgMemberRole_value[k]; ok {
		parsedRole = b2bv1.OrgMemberRole(v)
	} else if v, ok := b2bv1.OrgMemberRole_value["ORG_MEMBER_ROLE_"+k]; ok {
		parsedRole = b2bv1.OrgMemberRole(v)
	}

	return row.OrganisationID, parsedRole, row.OrganisationName, nil
}
