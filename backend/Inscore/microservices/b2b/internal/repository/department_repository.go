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

// departmentCols is the canonical SELECT column list.
// total_premium is JSONB — cast to TEXT for manual scanning.
const departmentCols = `
	department_id,
	name,
	business_id,
	COALESCE(employee_no, 0) AS employee_no,
	COALESCE(total_premium::TEXT, 'null') AS total_premium,
	created_at,
	updated_at
`

// ─── Scanner ──────────────────────────────────────────────────────────────────

func scanDepartment(row interface{ Scan(...any) error }) (*b2bv1.Department, error) {
	var (
		d                b2bv1.Department
		totalPremiumJSON sql.NullString
		employeeNo       sql.NullInt32
		createdAt        time.Time
		updatedAt        time.Time
	)

	if err := row.Scan(
		&d.DepartmentId,
		&d.Name,
		&d.BusinessId,
		&employeeNo,
		&totalPremiumJSON,
		&createdAt,
		&updatedAt,
	); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	d.EmployeeNo = employeeNo.Int32
	d.TotalPremium = scanMoney(totalPremiumJSON)
	if !createdAt.IsZero() {
		d.CreatedAt = timestamppb.New(createdAt)
	}
	if !updatedAt.IsZero() {
		d.UpdatedAt = timestamppb.New(updatedAt)
	}
	return &d, nil
}

// ─── Queries ──────────────────────────────────────────────────────────────────

func (r *PortalRepository) GetDepartment(ctx context.Context, departmentID string) (*b2bv1.Department, error) {
	query := fmt.Sprintf(
		`SELECT %s FROM b2b_schema.departments WHERE department_id = $1 AND deleted_at IS NULL LIMIT 1`,
		departmentCols,
	)
	row := r.db.WithContext(ctx).Raw(query, departmentID).Row()
	return scanDepartment(row)
}

func (r *PortalRepository) ListDepartments(
	ctx context.Context,
	pageSize, offset int,
	businessID string,
) ([]*b2bv1.Department, int64, error) {
	where := "deleted_at IS NULL"
	args := []interface{}{}
	argIdx := 1

	if businessID != "" {
		where += fmt.Sprintf(" AND business_id = $%d", argIdx)
		args = append(args, businessID)
		argIdx++
	}

	var total int64
	countArgs := make([]interface{}, len(args))
	copy(countArgs, args)
	if err := r.db.WithContext(ctx).Raw(
		fmt.Sprintf("SELECT COUNT(*) FROM b2b_schema.departments WHERE %s", where),
		countArgs...,
	).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(
		`SELECT %s FROM b2b_schema.departments WHERE %s ORDER BY name ASC LIMIT $%d OFFSET $%d`,
		departmentCols, where, argIdx, argIdx+1,
	)
	args = append(args, pageSize, offset)

	rows, err := r.db.WithContext(ctx).Raw(query, args...).Rows()
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var departments []*b2bv1.Department
	for rows.Next() {
		d, err := scanDepartment(rows)
		if err != nil {
			return nil, 0, err
		}
		departments = append(departments, d)
	}
	return departments, total, rows.Err()
}

func (r *PortalRepository) CreateDepartment(ctx context.Context, input domain.DepartmentCreateInput) (*b2bv1.Department, error) {
	id := input.DepartmentID
	if id == "" {
		id = newUUID()
	}

	if err := r.db.WithContext(ctx).Exec(`
		INSERT INTO b2b_schema.departments (department_id, name, business_id, employee_no, total_premium)
		VALUES ($1, $2, $3, 0, $4)`,
		id,
		input.Name,
		input.BusinessID,
		string(zeroMoneyJSON()), // pass as string so PostgreSQL receives JSON text for JSONB column
	).Error; err != nil {
		return nil, fmt.Errorf("insert department: %w", err)
	}

	return r.GetDepartment(ctx, id)
}

func (r *PortalRepository) UpdateDepartment(ctx context.Context, input domain.DepartmentUpdateInput) (*b2bv1.Department, error) {
	setClauses := []string{}
	args := []interface{}{}
	idx := 1

	if strings.TrimSpace(input.Name) != "" {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", idx))
		args = append(args, input.Name)
		idx++
	}
	if len(setClauses) == 0 {
		return r.GetDepartment(ctx, input.DepartmentID)
	}

	setClauses = append(setClauses, "updated_at = NOW()")
	query := fmt.Sprintf(
		"UPDATE b2b_schema.departments SET %s WHERE department_id = $%d AND deleted_at IS NULL",
		strings.Join(setClauses, ", "), idx,
	)
	args = append(args, input.DepartmentID)

	if err := r.db.WithContext(ctx).Exec(query, args...).Error; err != nil {
		return nil, fmt.Errorf("update department: %w", err)
	}
	return r.GetDepartment(ctx, input.DepartmentID)
}

func (r *PortalRepository) UpdateDepartmentTotalPremium(ctx context.Context, departmentID string) error {
	query := `
		WITH sum_data AS (
			SELECT
				COALESCE(SUM((premium_amount->>'amount')::bigint), 0) AS total_amount,
				COALESCE(SUM((premium_amount->>'decimal_amount')::numeric), 0) AS total_decimal,
				MAX(premium_amount->>'currency') AS currency
			FROM b2b_schema.employees
			WHERE department_id = $1
			  AND status = 'EMPLOYEE_STATUS_ACTIVE'
			  AND deleted_at IS NULL
			  AND premium_amount IS NOT NULL
		)
		UPDATE b2b_schema.departments 
		SET total_premium = jsonb_build_object(
			'amount', sum_data.total_amount,
			'decimal_amount', sum_data.total_decimal,
			'currency', COALESCE(sum_data.currency, 'BDT')
		),
		updated_at = NOW()
		FROM sum_data
		WHERE department_id = $1
	`
	return r.db.WithContext(ctx).Exec(query, departmentID).Error
}

func (r *PortalRepository) DeleteDepartment(ctx context.Context, departmentID string) error {
	// Safety: refuse if active employees exist
	var count int64
	if err := r.db.WithContext(ctx).Raw(
		`SELECT COUNT(*) FROM b2b_schema.employees
		 WHERE department_id = $1 AND status = 'EMPLOYEE_STATUS_ACTIVE' AND deleted_at IS NULL`,
		departmentID,
	).Scan(&count).Error; err != nil {
		return fmt.Errorf("check active employees: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("department has %d active employee(s): reassign before deleting", count)
	}

	return r.db.WithContext(ctx).Exec(
		"UPDATE b2b_schema.departments SET deleted_at = NOW() WHERE department_id = $1 AND deleted_at IS NULL",
		departmentID,
	).Error
}
