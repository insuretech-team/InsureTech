package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/newage-saint/insuretech/backend/inscore/microservices/b2b/internal/domain"
	b2bv1 "github.com/newage-saint/insuretech/gen/go/insuretech/b2b/entity/v1"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type PortalRepository struct {
	db *gorm.DB
}

func NewPortalRepository(db *gorm.DB) *PortalRepository {
	return &PortalRepository{db: db}
}

func (r *PortalRepository) ListDepartments(
	ctx context.Context,
	pageSize, offset int,
	businessID string,
) ([]*b2bv1.Department, int64, error) {
	q := r.db.WithContext(ctx).Table("b2b_schema.departments")
	if businessID != "" {
		q = q.Where("business_id = ?", businessID)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []departmentRow
	if err := q.
		Select(departmentSelectColumns()).
		Order("name ASC").
		Limit(pageSize).
		Offset(offset).
		Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	departments, err := mapDepartmentRows(rows)
	if err != nil {
		return nil, 0, err
	}

	return departments, total, nil
}

func (r *PortalRepository) GetDepartment(ctx context.Context, departmentID string) (*b2bv1.Department, error) {
	var row departmentRow
	err := r.db.WithContext(ctx).
		Table("b2b_schema.departments").
		Select(departmentSelectColumns()).
		Where("department_id = ?", departmentID).
		First(&row).Error
	if err != nil {
		return nil, err
	}
	return mapDepartmentRow(row)
}

func (r *PortalRepository) ListEmployees(
	ctx context.Context,
	pageSize, offset int,
	departmentID, businessID string,
) ([]*b2bv1.Employee, int64, error) {
	q := r.db.WithContext(ctx).Table("b2b_schema.employees")
	if departmentID != "" {
		q = q.Where("department_id = ?", departmentID)
	}
	if businessID != "" {
		q = q.Where("business_id = ?", businessID)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []employeeRow
	if err := q.
		Select(employeeSelectColumns()).
		Order("employee_id ASC").
		Limit(pageSize).
		Offset(offset).
		Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	employees, err := mapEmployeeRows(rows)
	if err != nil {
		return nil, 0, err
	}

	return employees, total, nil
}

func (r *PortalRepository) GetEmployee(ctx context.Context, employeeUUID string) (*b2bv1.Employee, error) {
	var row employeeRow
	err := r.db.WithContext(ctx).
		Table("b2b_schema.employees").
		Select(employeeSelectColumns()).
		Where("employee_uuid = ?", employeeUUID).
		First(&row).Error
	if err != nil {
		return nil, err
	}
	return mapEmployeeRow(row)
}

type departmentNameRow struct {
	DepartmentID string `gorm:"column:department_id"`
	Name         string `gorm:"column:name"`
}

func (r *PortalRepository) GetDepartmentNames(ctx context.Context, departmentIDs []string) (map[string]string, error) {
	result := make(map[string]string)
	if len(departmentIDs) == 0 {
		return result, nil
	}

	var rows []departmentNameRow
	if err := r.db.WithContext(ctx).
		Table("b2b_schema.departments").
		Select("department_id, name").
		Where("department_id IN ?", departmentIDs).
		Find(&rows).Error; err != nil {
		return nil, err
	}

	for _, row := range rows {
		result[row.DepartmentID] = row.Name
	}
	return result, nil
}

type employeeRow struct {
	EmployeeUUID      string    `gorm:"column:employee_uuid"`
	Name              string    `gorm:"column:name"`
	EmployeeID        string    `gorm:"column:employee_id"`
	DepartmentID      string    `gorm:"column:department_id"`
	BusinessID        string    `gorm:"column:business_id"`
	InsuranceCategory string    `gorm:"column:insurance_category"`
	AssignedPlanID    string    `gorm:"column:assigned_plan_id"`
	CoverageAmount    []byte    `gorm:"column:coverage_amount"`
	PremiumAmount     []byte    `gorm:"column:premium_amount"`
	Status            string    `gorm:"column:status"`
	CreatedAt         time.Time `gorm:"column:created_at"`
	UpdatedAt         time.Time `gorm:"column:updated_at"`
	NumberOfDependent int32     `gorm:"column:number_of_dependent"`
}

type departmentRow struct {
	DepartmentID string    `gorm:"column:department_id"`
	Name         string    `gorm:"column:name"`
	BusinessID   string    `gorm:"column:business_id"`
	EmployeeNo   int32     `gorm:"column:employee_no"`
	TotalPremium []byte    `gorm:"column:total_premium"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

type catalogPlanRow struct {
	ProductID         string `gorm:"column:product_id"`
	ProductName       string `gorm:"column:product_name"`
	PlanID            string `gorm:"column:plan_id"`
	PlanName          string `gorm:"column:plan_name"`
	InsuranceCategory string `gorm:"column:insurance_category"`
	PremiumAmount     []byte `gorm:"column:premium_amount"`
}

func employeeSelectColumns() string {
	return strings.Join([]string{
		"employee_uuid",
		"name",
		"employee_id",
		"department_id",
		"business_id",
		"insurance_category",
		"assigned_plan_id",
		"coverage_amount",
		"premium_amount",
		"status",
		"created_at",
		"updated_at",
		"number_of_dependent",
	}, ", ")
}

func departmentSelectColumns() string {
	return strings.Join([]string{
		"department_id",
		"name",
		"business_id",
		"employee_no",
		"total_premium",
		"created_at",
		"updated_at",
	}, ", ")
}

func mapDepartmentRows(rows []departmentRow) ([]*b2bv1.Department, error) {
	departments := make([]*b2bv1.Department, 0, len(rows))
	for _, row := range rows {
		department, err := mapDepartmentRow(row)
		if err != nil {
			return nil, err
		}
		departments = append(departments, department)
	}
	return departments, nil
}

func mapDepartmentRow(row departmentRow) (*b2bv1.Department, error) {
	totalPremium, err := parseMoney(row.TotalPremium)
	if err != nil {
		return nil, fmt.Errorf("parse total_premium for department %s: %w", row.DepartmentID, err)
	}

	department := &b2bv1.Department{
		DepartmentId: row.DepartmentID,
		Name:         row.Name,
		BusinessId:   row.BusinessID,
		EmployeeNo:   row.EmployeeNo,
		TotalPremium: totalPremium,
	}

	if !row.CreatedAt.IsZero() {
		department.CreatedAt = timestamppb.New(row.CreatedAt)
	}
	if !row.UpdatedAt.IsZero() {
		department.UpdatedAt = timestamppb.New(row.UpdatedAt)
	}

	return department, nil
}

func mapEmployeeRows(rows []employeeRow) ([]*b2bv1.Employee, error) {
	employees := make([]*b2bv1.Employee, 0, len(rows))
	for _, row := range rows {
		employee, err := mapEmployeeRow(row)
		if err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}
	return employees, nil
}

func mapEmployeeRow(row employeeRow) (*b2bv1.Employee, error) {
	coverageAmount, err := parseMoney(row.CoverageAmount)
	if err != nil {
		return nil, fmt.Errorf("parse coverage_amount for employee %s: %w", row.EmployeeUUID, err)
	}
	premiumAmount, err := parseMoney(row.PremiumAmount)
	if err != nil {
		return nil, fmt.Errorf("parse premium_amount for employee %s: %w", row.EmployeeUUID, err)
	}

	employee := &b2bv1.Employee{
		EmployeeUuid:      row.EmployeeUUID,
		Name:              row.Name,
		EmployeeId:        row.EmployeeID,
		DepartmentId:      row.DepartmentID,
		BusinessId:        row.BusinessID,
		InsuranceCategory: parseInsuranceType(row.InsuranceCategory),
		AssignedPlanId:    row.AssignedPlanID,
		CoverageAmount:    coverageAmount,
		PremiumAmount:     premiumAmount,
		Status:            parseEmployeeStatus(row.Status),
		NumberOfDependent: row.NumberOfDependent,
	}

	if !row.CreatedAt.IsZero() {
		employee.CreatedAt = timestamppb.New(row.CreatedAt)
	}
	if !row.UpdatedAt.IsZero() {
		employee.UpdatedAt = timestamppb.New(row.UpdatedAt)
	}

	return employee, nil
}

func parseMoney(raw []byte) (*commonv1.Money, error) {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		return nil, nil
	}

	var money commonv1.Money
	if err := json.Unmarshal(trimmed, &money); err != nil {
		return nil, err
	}
	return &money, nil
}

func mustMarshalMoney(money *commonv1.Money) ([]byte, error) {
	if money == nil {
		return []byte("null"), nil
	}
	return json.Marshal(money)
}

func parseInsuranceType(value string) commonv1.InsuranceType {
	value = strings.TrimSpace(value)
	if value == "" {
		return commonv1.InsuranceType_INSURANCE_TYPE_UNSPECIFIED
	}
	if enumValue, ok := commonv1.InsuranceType_value[value]; ok {
		return commonv1.InsuranceType(enumValue)
	}
	normalized := "INSURANCE_TYPE_" + strings.ToUpper(value)
	if enumValue, ok := commonv1.InsuranceType_value[normalized]; ok {
		return commonv1.InsuranceType(enumValue)
	}
	return commonv1.InsuranceType_INSURANCE_TYPE_UNSPECIFIED
}

func parseEmployeeStatus(value string) b2bv1.EmployeeStatus {
	value = strings.TrimSpace(value)
	if value == "" {
		return b2bv1.EmployeeStatus_EMPLOYEE_STATUS_UNSPECIFIED
	}
	if enumValue, ok := b2bv1.EmployeeStatus_value[value]; ok {
		return b2bv1.EmployeeStatus(enumValue)
	}
	normalized := "EMPLOYEE_STATUS_" + strings.ToUpper(value)
	if enumValue, ok := b2bv1.EmployeeStatus_value[normalized]; ok {
		return b2bv1.EmployeeStatus(enumValue)
	}
	return b2bv1.EmployeeStatus_EMPLOYEE_STATUS_UNSPECIFIED
}

func parsePurchaseOrderStatus(value string) b2bv1.PurchaseOrderStatus {
	value = strings.TrimSpace(value)
	if value == "" {
		return b2bv1.PurchaseOrderStatus_PURCHASE_ORDER_STATUS_UNSPECIFIED
	}
	if enumValue, ok := b2bv1.PurchaseOrderStatus_value[value]; ok {
		return b2bv1.PurchaseOrderStatus(enumValue)
	}
	normalized := "PURCHASE_ORDER_STATUS_" + strings.ToUpper(value)
	if enumValue, ok := b2bv1.PurchaseOrderStatus_value[normalized]; ok {
		return b2bv1.PurchaseOrderStatus(enumValue)
	}
	return b2bv1.PurchaseOrderStatus_PURCHASE_ORDER_STATUS_UNSPECIFIED
}

func purchaseOrderStatusToDB(value b2bv1.PurchaseOrderStatus) string {
	switch value {
	case b2bv1.PurchaseOrderStatus_PURCHASE_ORDER_STATUS_DRAFT:
		return "DRAFT"
	case b2bv1.PurchaseOrderStatus_PURCHASE_ORDER_STATUS_SUBMITTED:
		return "SUBMITTED"
	case b2bv1.PurchaseOrderStatus_PURCHASE_ORDER_STATUS_APPROVED:
		return "APPROVED"
	case b2bv1.PurchaseOrderStatus_PURCHASE_ORDER_STATUS_FULFILLED:
		return "FULFILLED"
	case b2bv1.PurchaseOrderStatus_PURCHASE_ORDER_STATUS_REJECTED:
		return "REJECTED"
	default:
		return ""
	}
}

func moneyFromMinor(amount int64, currency string) *commonv1.Money {
	if currency == "" {
		currency = "BDT"
	}
	return &commonv1.Money{
		Amount:        amount,
		Currency:      currency,
		DecimalAmount: float64(amount) / 100,
	}
}

func parseMoneyFlexible(raw []byte, defaultCurrency string) (*commonv1.Money, error) {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		return nil, nil
	}

	if trimmed[0] == '{' {
		return parseMoney(trimmed)
	}

	if amount, err := strconv.ParseInt(string(trimmed), 10, 64); err == nil {
		return moneyFromMinor(amount, defaultCurrency), nil
	}

	if amount, err := strconv.ParseFloat(string(trimmed), 64); err == nil {
		return &commonv1.Money{
			Amount:        int64(amount),
			Currency:      defaultCurrency,
			DecimalAmount: amount / 100,
		}, nil
	}

	return nil, fmt.Errorf("unsupported money payload %q", string(trimmed))
}

func mapCatalogPlanRow(row catalogPlanRow) (*domain.CatalogPlan, error) {
	premiumAmount, err := parseMoneyFlexible(row.PremiumAmount, "BDT")
	if err != nil {
		return nil, fmt.Errorf("parse premium_amount for plan %s: %w", row.PlanID, err)
	}

	return &domain.CatalogPlan{
		ProductID:         row.ProductID,
		ProductName:       row.ProductName,
		PlanID:            row.PlanID,
		PlanName:          row.PlanName,
		InsuranceCategory: parseInsuranceType(row.InsuranceCategory),
		PremiumAmount:     premiumAmount,
	}, nil
}

func (r *PortalRepository) ListCatalogPlans(ctx context.Context) ([]*domain.CatalogPlan, error) {
	var rows []catalogPlanRow
	err := r.db.WithContext(ctx).
		Table("insurance_schema.product_plans AS pp").
		Select(strings.Join([]string{
			"pp.product_id AS product_id",
			"p.product_name AS product_name",
			"pp.plan_id AS plan_id",
			"pp.plan_name AS plan_name",
			"p.category AS insurance_category",
			"pp.premium_amount AS premium_amount",
		}, ", ")).
		Joins("JOIN insurance_schema.products AS p ON p.product_id = pp.product_id").
		Where("p.status IN ?", []string{"ACTIVE", "PRODUCT_STATUS_ACTIVE", "2"}).
		Order("p.product_name ASC, pp.plan_name ASC").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	items := make([]*domain.CatalogPlan, 0, len(rows))
	for _, row := range rows {
		item, err := mapCatalogPlanRow(row)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *PortalRepository) GetCatalogPlansByPlanIDs(ctx context.Context, planIDs []string) (map[string]*domain.CatalogPlan, error) {
	result := make(map[string]*domain.CatalogPlan)
	if len(planIDs) == 0 {
		return result, nil
	}

	var rows []catalogPlanRow
	err := r.db.WithContext(ctx).
		Table("insurance_schema.product_plans AS pp").
		Select(strings.Join([]string{
			"pp.product_id AS product_id",
			"p.product_name AS product_name",
			"pp.plan_id AS plan_id",
			"pp.plan_name AS plan_name",
			"p.category AS insurance_category",
			"pp.premium_amount AS premium_amount",
		}, ", ")).
		Joins("JOIN insurance_schema.products AS p ON p.product_id = pp.product_id").
		Where("pp.plan_id IN ?", planIDs).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		item, err := mapCatalogPlanRow(row)
		if err != nil {
			return nil, err
		}
		result[row.PlanID] = item
	}
	return result, nil
}
