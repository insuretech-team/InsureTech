package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	entityv1 "github.com/newage-saint/insuretech/gen/go/insuretech/authz/entity/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// RoleRepo implements domain.RoleRepository using proto structs directly as GORM models.
type RoleRepo struct{ db *gorm.DB }

func NewRoleRepo(db *gorm.DB) *RoleRepo { return &RoleRepo{db: db} }

func (r *RoleRepo) Create(ctx context.Context, role *entityv1.Role) (*entityv1.Role, error) {
	if role == nil {
		return nil, errors.New("role.Create: nil role")
	}
	now := time.Now()
	values := map[string]any{
		"role_id":     role.RoleId,
		"name":        role.Name,
		"portal":      role.Portal.String(),
		"description": role.Description,
		"is_system":   role.IsSystem,
		"is_active":   role.IsActive,
		"created_by":  nullableUUID(role.CreatedBy),
		"created_at":  now,
		"updated_at":  now,
	}
	if role.RoleId == "" {
		delete(values, "role_id")
	}
	if err := r.db.WithContext(ctx).Table("authz_schema.roles").Create(values).Error; err != nil {
		return nil, errors.New("role.Create: " + err.Error())
	}
	if role.RoleId == "" {
		created, err := r.GetByNameAndPortal(ctx, role.Name, role.Portal)
		if err != nil {
			return nil, err
		}
		return created, nil
	}
	return r.GetByID(ctx, role.RoleId)
}

func (r *RoleRepo) GetByID(ctx context.Context, id string) (*entityv1.Role, error) {
	role, err := r.queryOne(ctx,
		`SELECT role_id, name, portal, description, is_system, is_active, created_by, created_at, updated_at, deleted_at
		   FROM authz_schema.roles
		  WHERE role_id = ?
		  LIMIT 1`,
		id,
	)
	if err != nil {
		return nil, errors.New("role.GetByID: " + err.Error())
	}
	return role, nil
}

func (r *RoleRepo) GetByNameAndPortal(ctx context.Context, name string, portal entityv1.Portal) (*entityv1.Role, error) {
	role, err := r.queryOne(ctx,
		`SELECT role_id, name, portal, description, is_system, is_active, created_by, created_at, updated_at, deleted_at
		   FROM authz_schema.roles
		  WHERE name = ? AND (portal = ? OR portal = ?)
		  LIMIT 1`,
		name,
		portal.String(),
		strings.TrimPrefix(portal.String(), "PORTAL_"),
	)
	if err != nil {
		return nil, errors.New("role.GetByNameAndPortal: " + err.Error())
	}
	return role, nil
}

func (r *RoleRepo) GetByName(ctx context.Context, portal string, name string) (*entityv1.Role, error) {
	q := `SELECT role_id, name, portal, description, is_system, is_active, created_by, created_at, updated_at, deleted_at
	       FROM authz_schema.roles
	      WHERE name = ?`
	args := []any{name}
	if v, ok := entityv1.Portal_value[portal]; ok {
		_ = v
		q += " AND (portal = ? OR portal = ?)"
		args = append(args, portal, strings.TrimPrefix(portal, "PORTAL_"))
	} else {
		q += " AND portal = ?"
		args = append(args, portal)
	}
	q += " LIMIT 1"
	role, err := r.queryOne(ctx, q, args...)
	if err != nil {
		return nil, errors.New("role.GetByName: " + err.Error())
	}
	return role, nil
}

func (r *RoleRepo) SoftDelete(ctx context.Context, roleID string) error {
	return r.db.WithContext(ctx).Table("authz_schema.roles").
		Where("role_id = ?", roleID).
		Update("is_active", false).Error
}

func (r *RoleRepo) List(ctx context.Context, portal entityv1.Portal, activeOnly bool, limit, offset int) ([]*entityv1.Role, error) {
	q := `SELECT role_id, name, portal, description, is_system, is_active, created_by, created_at, updated_at, deleted_at
	       FROM authz_schema.roles
	      WHERE 1=1`
	args := make([]any, 0, 4)
	if portal != entityv1.Portal_PORTAL_UNSPECIFIED {
		q += " AND (portal = ? OR portal = ?)"
		args = append(args, portal.String(), strings.TrimPrefix(portal.String(), "PORTAL_"))
	}
	if activeOnly {
		q += " AND is_active = true"
	}
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	q += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.WithContext(ctx).Raw(q, args...).Rows()
	if err != nil {
		return nil, errors.New("role.List: " + err.Error())
	}
	defer rows.Close()

	roles := make([]*entityv1.Role, 0, limit)
	for rows.Next() {
		role, scanErr := scanRole(rows)
		if scanErr != nil {
			return nil, errors.New("role.List scan: " + scanErr.Error())
		}
		roles = append(roles, role)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.New("role.List: " + err.Error())
	}
	return roles, nil
}

func (r *RoleRepo) Update(ctx context.Context, role *entityv1.Role) (*entityv1.Role, error) {
	if err := r.db.WithContext(ctx).Table("authz_schema.roles").
		Where("role_id = ?", role.RoleId).
		Updates(map[string]any{
			"name":        role.Name,
			"portal":      role.Portal.String(),
			"description": role.Description,
			"is_system":   role.IsSystem,
			"is_active":   role.IsActive,
			"created_by":  role.CreatedBy,
			"updated_at":  gorm.Expr("NOW()"),
		}).Error; err != nil {
		return nil, errors.New("role.Update: " + err.Error())
	}
	return role, nil
}

func (r *RoleRepo) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Exec(`DELETE FROM authz_schema.roles WHERE role_id = ?`, id).Error; err != nil {
		return errors.New("role.Delete: " + err.Error())
	}
	return nil
}

func (r *RoleRepo) queryOne(ctx context.Context, query string, args ...any) (*entityv1.Role, error) {
	row := r.db.WithContext(ctx).Raw(query, args...).Row()
	role, err := scanRole(row)
	if err != nil {
		return nil, err
	}
	return role, nil
}

type roleRowScanner interface {
	Scan(dest ...any) error
}

func scanRole(scanner roleRowScanner) (*entityv1.Role, error) {
	var (
		role entityv1.Role

		portalStr   sql.NullString
		description sql.NullString
		createdBy   sql.NullString

		isSystem sql.NullBool
		isActive sql.NullBool

		createdAt sql.NullTime
		updatedAt sql.NullTime
		deletedAt sql.NullTime
	)

	if err := scanner.Scan(
		&role.RoleId,
		&role.Name,
		&portalStr,
		&description,
		&isSystem,
		&isActive,
		&createdBy,
		&createdAt,
		&updatedAt,
		&deletedAt,
	); err != nil {
		return nil, err
	}

	if portalStr.Valid {
		if v, ok := entityv1.Portal_value[portalStr.String]; ok {
			role.Portal = entityv1.Portal(v)
		} else if v, ok := entityv1.Portal_value["PORTAL_"+portalStr.String]; ok {
			role.Portal = entityv1.Portal(v)
		}
	}
	if description.Valid {
		role.Description = description.String
	}
	if isSystem.Valid {
		role.IsSystem = isSystem.Bool
	}
	if isActive.Valid {
		role.IsActive = isActive.Bool
	}
	if createdBy.Valid {
		role.CreatedBy = createdBy.String
	}
	if createdAt.Valid {
		role.CreatedAt = timestamppb.New(createdAt.Time)
	}
	if updatedAt.Valid {
		role.UpdatedAt = timestamppb.New(updatedAt.Time)
	}
	if deletedAt.Valid {
		role.DeletedAt = timestamppb.New(deletedAt.Time)
	}

	return &role, nil
}

func nullableUUID(v string) any {
	trimmed := strings.TrimSpace(v)
	if trimmed == "" {
		return nil
	}
	if _, err := uuid.Parse(trimmed); err != nil {
		return nil
	}
	return trimmed
}
