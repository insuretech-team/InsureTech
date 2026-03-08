package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/b2b/internal/domain"
	b2bv1 "github.com/newage-saint/insuretech/gen/go/insuretech/b2b/entity/v1"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

// ─── SQL ──────────────────────────────────────────────────────────────────────

// employeeCols is the canonical SELECT column list for the employees table.
// Money columns (coverage_amount, premium_amount) are cast to TEXT so we can
// scan them into sql.NullString and unmarshal the JSON manually — GORM cannot
// scan JSONB into nested proto structs natively.
const employeeCols = `
	employee_uuid,
	name,
	employee_id,
	department_id,
	business_id,
	COALESCE(insurance_category, '') AS insurance_category,
	COALESCE(assigned_plan_id::TEXT, '') AS assigned_plan_id,
	COALESCE(coverage_amount::TEXT, 'null') AS coverage_amount,
	COALESCE(premium_amount::TEXT, 'null') AS premium_amount,
	status,
	created_at,
	updated_at,
	COALESCE(number_of_dependent, 0) AS number_of_dependent,
	COALESCE(email, '') AS email,
	COALESCE(mobile_number, '') AS mobile_number,
	COALESCE(CAST(date_of_birth AS TEXT), '') AS date_of_birth,
	COALESCE(CAST(date_of_joining AS TEXT), '') AS date_of_joining,
	COALESCE(gender, '') AS gender,
	COALESCE(user_id::TEXT, '') AS user_id
`

// ─── Scanner ──────────────────────────────────────────────────────────────────

// scanEmployee reads one row from a raw SQL query into *b2bv1.Employee.
// Uses the same sql.Null* pattern as authn's user_repository.go.
func scanEmployee(row interface{ Scan(...any) error }) (*b2bv1.Employee, error) {
	var (
		e                   b2bv1.Employee
		insuranceCategoryStr sql.NullString
		assignedPlanID       sql.NullString
		coverageAmountJSON   sql.NullString
		premiumAmountJSON    sql.NullString
		statusStr            sql.NullString
		createdAt            time.Time
		updatedAt            time.Time
		numberOfDependent    sql.NullInt32
		email                sql.NullString
		mobileNumber         sql.NullString
		dateOfBirth          sql.NullString
		dateOfJoining        sql.NullString
		genderStr            sql.NullString
		userID               sql.NullString
	)

	if err := row.Scan(
		&e.EmployeeUuid,
		&e.Name,
		&e.EmployeeId,
		&e.DepartmentId,
		&e.BusinessId,
		&insuranceCategoryStr,
		&assignedPlanID,
		&coverageAmountJSON,
		&premiumAmountJSON,
		&statusStr,
		&createdAt,
		&updatedAt,
		&numberOfDependent,
		&email,
		&mobileNumber,
		&dateOfBirth,
		&dateOfJoining,
		&genderStr,
		&userID,
	); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	// Enum: InsuranceCategory
	if insuranceCategoryStr.Valid {
		e.InsuranceCategory = parseInsuranceType(insuranceCategoryStr.String)
	}

	// String fields
	if assignedPlanID.Valid { e.AssignedPlanId = assignedPlanID.String }

	// Money: JSONB → *commonv1.Money
	e.CoverageAmount = scanMoney(coverageAmountJSON)
	e.PremiumAmount = scanMoney(premiumAmountJSON)

	// Enum: Status
	if statusStr.Valid {
		k := strings.ToUpper(statusStr.String)
		if v, ok := b2bv1.EmployeeStatus_value[k]; ok {
			e.Status = b2bv1.EmployeeStatus(v)
		} else if v, ok := b2bv1.EmployeeStatus_value["EMPLOYEE_STATUS_"+k]; ok {
			e.Status = b2bv1.EmployeeStatus(v)
		}
	}

	// Timestamps
	if !createdAt.IsZero() { e.CreatedAt = timestamppb.New(createdAt) }
	if !updatedAt.IsZero() { e.UpdatedAt = timestamppb.New(updatedAt) }

	// Scalar
	if numberOfDependent.Valid { e.NumberOfDependent = numberOfDependent.Int32 }
	if email.Valid         { e.Email = email.String }
	if mobileNumber.Valid  { e.MobileNumber = mobileNumber.String }
	if dateOfBirth.Valid   { e.DateOfBirth = dateOfBirth.String }
	if dateOfJoining.Valid { e.DateOfJoining = dateOfJoining.String }

	// Enum: Gender — stored as short string ("MALE","FEMALE","OTHER") in varchar(20)
	if genderStr.Valid && genderStr.String != "" {
		k := strings.ToUpper(strings.TrimSpace(genderStr.String))
		// Try short form first: "MALE" → "EMPLOYEE_GENDER_MALE"
		if v, ok := b2bv1.EmployeeGender_value["EMPLOYEE_GENDER_"+k]; ok {
			e.Gender = b2bv1.EmployeeGender(v)
		} else if v, ok := b2bv1.EmployeeGender_value[k]; ok {
			// Try full proto name as fallback
			e.Gender = b2bv1.EmployeeGender(v)
		}
	}

	if userID.Valid { e.UserId = userID.String }

	return &e, nil
}

// ─── Queries ──────────────────────────────────────────────────────────────────

func (r *PortalRepository) GetEmployee(ctx context.Context, employeeUUID string) (*b2bv1.Employee, error) {
	query := fmt.Sprintf(`SELECT %s FROM b2b_schema.employees WHERE employee_uuid = $1 AND deleted_at IS NULL LIMIT 1`, employeeCols)
	row := r.db.WithContext(ctx).Raw(query, employeeUUID).Row()
	return scanEmployee(row)
}

func (r *PortalRepository) ListEmployees(
	ctx context.Context,
	pageSize, offset int,
	departmentID, businessID string,
	status b2bv1.EmployeeStatus,
) ([]*b2bv1.Employee, int64, error) {
	// Build WHERE clause
	where := "deleted_at IS NULL"
	args := []interface{}{}
	argIdx := 1

	if businessID != "" {
		where += fmt.Sprintf(" AND business_id = $%d", argIdx)
		args = append(args, businessID)
		argIdx++
	}
	if departmentID != "" {
		where += fmt.Sprintf(" AND department_id = $%d", argIdx)
		args = append(args, departmentID)
		argIdx++
	}
	if status != b2bv1.EmployeeStatus_EMPLOYEE_STATUS_UNSPECIFIED {
		where += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, employeeStatusStr(status))
		argIdx++
	}

	// Count
	var total int64
	countArgs := make([]interface{}, len(args))
	copy(countArgs, args)
	if err := r.db.WithContext(ctx).Raw(
		fmt.Sprintf("SELECT COUNT(*) FROM b2b_schema.employees WHERE %s", where),
		countArgs...,
	).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	// Paginated fetch
	query := fmt.Sprintf(
		`SELECT %s FROM b2b_schema.employees WHERE %s ORDER BY employee_id ASC LIMIT $%d OFFSET $%d`,
		employeeCols, where, argIdx, argIdx+1,
	)
	args = append(args, pageSize, offset)

	rows, err := r.db.WithContext(ctx).Raw(query, args...).Rows()
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var employees []*b2bv1.Employee
	for rows.Next() {
		e, err := scanEmployee(rows)
		if err != nil {
			return nil, 0, err
		}
		employees = append(employees, e)
	}
	return employees, total, rows.Err()
}

func (r *PortalRepository) CreateEmployee(ctx context.Context, input domain.EmployeeCreateInput) (*b2bv1.Employee, error) {
	id := input.EmployeeUUID
	if id == "" {
		id = newUUID()
	}

	coverageJSONBytes, err := marshalMoney(input.CoverageAmount)
	if err != nil {
		return nil, fmt.Errorf("marshal coverage_amount: %w", err)
	}
	// Pass as string so PostgreSQL receives valid JSON text for the JSONB column.
	coverageJSON := string(coverageJSONBytes)

	premiumJSONBytes, err := marshalMoney(input.PremiumAmount)
	if err != nil {
		return nil, fmt.Errorf("marshal premium_amount: %w", err)
	}
	premiumJSON := string(premiumJSONBytes)

	if err := r.db.WithContext(ctx).Exec(`
		INSERT INTO b2b_schema.employees (
			employee_uuid, name, employee_id, department_id, business_id,
			insurance_category, assigned_plan_id,
			coverage_amount, premium_amount,
			status, number_of_dependent,
			email, mobile_number, date_of_birth, date_of_joining, gender, user_id
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7,
			$8, $9,
			$10, $11,
			$12, $13, $14::date, $15::date, $16, $17
		)`,
		id,
		input.Name,
		input.EmployeeID,
		input.DepartmentID,
		input.BusinessID,
		nullableStr(input.InsuranceCategory.String()),
		nullableStr(input.AssignedPlanID),
		coverageJSON,
		premiumJSON, // plan premium amount resolved from catalog before insert
		"EMPLOYEE_STATUS_ACTIVE",
		input.NumberOfDependent,
		nullableStr(input.Email),
		nullableStr(input.MobileNumber),
		nullableStr(input.DateOfBirth),
		nullableStr(input.DateOfJoining),
		nullableStr(employeeGenderStr(input.Gender)),
		nullableStr(input.UserID),
	).Error; err != nil {
		return nil, fmt.Errorf("insert employee: %w", err)
	}

	return r.GetEmployee(ctx, id)
}

func (r *PortalRepository) UpdateEmployee(ctx context.Context, input domain.EmployeeUpdateInput) (*b2bv1.Employee, error) {
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
	addRaw := func(col string, val interface{}) {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", col, idx))
		args = append(args, val)
		idx++
	}

	addStr("name", input.Name)
	addStr("department_id", input.DepartmentID)
	addStr("email", input.Email)
	addStr("mobile_number", input.MobileNumber)
	if input.DateOfBirth != "" {
		setClauses = append(setClauses, fmt.Sprintf("date_of_birth = $%d::date", idx))
		args = append(args, input.DateOfBirth)
		idx++
	}
	if input.DateOfJoining != "" {
		setClauses = append(setClauses, fmt.Sprintf("date_of_joining = $%d::date", idx))
		args = append(args, input.DateOfJoining)
		idx++
	}
	if input.Gender != b2bv1.EmployeeGender_EMPLOYEE_GENDER_UNSPECIFIED {
		addStr("gender", employeeGenderStr(input.Gender))
	}
	if input.InsuranceCategory != commonv1.InsuranceType_INSURANCE_TYPE_UNSPECIFIED {
		addStr("insurance_category", input.InsuranceCategory.String())
	}
	addStr("assigned_plan_id", input.AssignedPlanID)
	if input.CoverageAmount != nil {
		coverageJSON, err := marshalMoney(input.CoverageAmount)
		if err != nil {
			return nil, fmt.Errorf("marshal coverage_amount: %w", err)
		}
		// Pass as string (not []byte) so PostgreSQL receives valid JSON text for the JSONB column.
		// GORM sends []byte as bytea which PostgreSQL cannot implicitly cast to jsonb.
		addRaw("coverage_amount", string(coverageJSON))
	}
	if input.NumberOfDependent > 0 {
		addRaw("number_of_dependent", input.NumberOfDependent)
	}
	if input.Status != b2bv1.EmployeeStatus_EMPLOYEE_STATUS_UNSPECIFIED {
		addStr("status", employeeStatusStr(input.Status))
	}

	if len(setClauses) == 0 {
		return r.GetEmployee(ctx, input.EmployeeUUID)
	}

	// Add updated_at
	setClauses = append(setClauses, "updated_at = NOW()")

	query := fmt.Sprintf(
		"UPDATE b2b_schema.employees SET %s WHERE employee_uuid = $%d AND deleted_at IS NULL",
		strings.Join(setClauses, ", "), idx,
	)
	args = append(args, input.EmployeeUUID)

	if err := r.db.WithContext(ctx).Exec(query, args...).Error; err != nil {
		return nil, fmt.Errorf("update employee: %w", err)
	}
	return r.GetEmployee(ctx, input.EmployeeUUID)
}

func (r *PortalRepository) DeleteEmployee(ctx context.Context, employeeUUID string) error {
	return r.db.WithContext(ctx).Exec(
		"UPDATE b2b_schema.employees SET deleted_at = NOW() WHERE employee_uuid = $1 AND deleted_at IS NULL",
		employeeUUID,
	).Error
}

// GetDepartmentNames batch-fetches department names for enrichment views.
func (r *PortalRepository) GetDepartmentNames(ctx context.Context, departmentIDs []string) (map[string]string, error) {
	result := make(map[string]string)
	if len(departmentIDs) == 0 {
		return result, nil
	}

	type nameRow struct {
		DepartmentID string `gorm:"column:department_id"`
		Name         string `gorm:"column:name"`
	}
	var rows []nameRow
	if err := r.db.WithContext(ctx).Raw(
		"SELECT department_id, name FROM b2b_schema.departments WHERE department_id = ANY($1) AND deleted_at IS NULL",
		departmentIDs,
	).Scan(&rows).Error; err != nil {
		return nil, err
	}
	for _, row := range rows {
		result[row.DepartmentID] = row.Name
	}
	return result, nil
}

// ─── helpers ─────────────────────────────────────────────────────────────────

// nullableStr returns nil if s is empty, otherwise returns &s.
// Used to insert NULL into nullable VARCHAR columns instead of empty string.
func nullableStr(s string) interface{} {
	if s == "" || s == "INSURANCE_TYPE_UNSPECIFIED" || s == "EMPLOYEE_GENDER_UNSPECIFIED" {
		return nil
	}
	return s
}
