package repository

import (
	"database/sql"
	"errors"
	"reflect"
	"testing"
	"time"

	b2bv1 "github.com/newage-saint/insuretech/gen/go/insuretech/b2b/entity/v1"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type fakeScanRow struct {
	values []any
	err    error
}

func (r fakeScanRow) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}
	for i := range dest {
		reflect.ValueOf(dest[i]).Elem().Set(reflect.ValueOf(r.values[i]))
	}
	return nil
}

func TestRepositoryMoneyAndEnumHelpers(t *testing.T) {
	require.Nil(t, scanMoney(sql.NullString{}))

	jsonMoney, err := marshalMoney(&commonv1.Money{Amount: 50000, Currency: "BDT", DecimalAmount: 500})
	require.NoError(t, err)

	money := scanMoney(sql.NullString{Valid: true, String: string(jsonMoney)})
	require.NotNil(t, money)
	require.Equal(t, int64(50000), money.Amount)
	require.Equal(t, "BDT", money.Currency)

	zeroMoney := scanMoney(sql.NullString{Valid: true, String: string(zeroMoneyJSON())})
	require.NotNil(t, zeroMoney)
	require.Equal(t, "BDT", zeroMoney.Currency)

	require.Equal(t, "EMPLOYEE_STATUS_ACTIVE", employeeStatusStr(b2bv1.EmployeeStatus_EMPLOYEE_STATUS_ACTIVE))
	require.Equal(t, "EMPLOYEE_STATUS_INACTIVE", employeeStatusStr(b2bv1.EmployeeStatus_EMPLOYEE_STATUS_INACTIVE))
	require.Equal(t, "EMPLOYEE_STATUS_ACTIVE", employeeStatusStr(b2bv1.EmployeeStatus_EMPLOYEE_STATUS_UNSPECIFIED))

	require.Equal(t, "MALE", employeeGenderStr(b2bv1.EmployeeGender_EMPLOYEE_GENDER_MALE))
	require.Equal(t, "FEMALE", employeeGenderStr(b2bv1.EmployeeGender_EMPLOYEE_GENDER_FEMALE))
	require.Equal(t, "", employeeGenderStr(b2bv1.EmployeeGender_EMPLOYEE_GENDER_UNSPECIFIED))

	require.Equal(t, "ORGANISATION_STATUS_ACTIVE", organisationStatusStr(b2bv1.OrganisationStatus_ORGANISATION_STATUS_UNSPECIFIED))
	require.Equal(t, b2bv1.OrganisationStatus_ORGANISATION_STATUS_ACTIVE.String(), organisationStatusStr(b2bv1.OrganisationStatus_ORGANISATION_STATUS_ACTIVE))

	require.Equal(t, "ORG_MEMBER_ROLE_HR_MANAGER", orgMemberRoleStr(b2bv1.OrgMemberRole_ORG_MEMBER_ROLE_UNSPECIFIED))
	require.Equal(t, b2bv1.OrgMemberRole_ORG_MEMBER_ROLE_BUSINESS_ADMIN.String(), orgMemberRoleStr(b2bv1.OrgMemberRole_ORG_MEMBER_ROLE_BUSINESS_ADMIN))

	require.Equal(t, "PURCHASE_ORDER_STATUS_SUBMITTED", purchaseOrderStatusStr(b2bv1.PurchaseOrderStatus_PURCHASE_ORDER_STATUS_UNSPECIFIED))
	require.Equal(t, b2bv1.PurchaseOrderStatus_PURCHASE_ORDER_STATUS_APPROVED.String(), purchaseOrderStatusStr(b2bv1.PurchaseOrderStatus_PURCHASE_ORDER_STATUS_APPROVED))

	require.Equal(t, commonv1.InsuranceType_INSURANCE_TYPE_HEALTH, parseInsuranceType("INSURANCE_TYPE_HEALTH"))
	require.Equal(t, commonv1.InsuranceType_INSURANCE_TYPE_LIFE, parseInsuranceType("life"))
	require.Equal(t, commonv1.InsuranceType_INSURANCE_TYPE_UNSPECIFIED, parseInsuranceType(""))
	require.Equal(t, commonv1.InsuranceType_INSURANCE_TYPE_UNSPECIFIED, parseInsuranceType("unknown"))

	require.Nil(t, toProtoTS(time.Time{}))
	require.NotNil(t, toProtoTS(time.Now()))

	require.Nil(t, nullableStr(""))
	require.Nil(t, nullableStr("INSURANCE_TYPE_UNSPECIFIED"))
	require.Nil(t, nullableStr("EMPLOYEE_GENDER_UNSPECIFIED"))
	require.Equal(t, "x", nullableStr("x"))
}

func TestRepositoryScanners(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	moneyJSON, err := marshalMoney(&commonv1.Money{Amount: 1000, Currency: "BDT", DecimalAmount: 10})
	require.NoError(t, err)

	plan, err := scanCatalogPlan(fakeScanRow{values: []any{
		"prod-1",
		"Health Insurance",
		"plan-1",
		"Seba",
		sql.NullString{Valid: true, String: "health"},
		sql.NullString{Valid: true, String: string(moneyJSON)},
	}})
	require.NoError(t, err)
	require.Equal(t, "plan-1", plan.PlanID)
	require.Equal(t, commonv1.InsuranceType_INSURANCE_TYPE_HEALTH, plan.InsuranceCategory)
	require.NotNil(t, plan.PremiumAmount)

	dept, err := scanDepartment(fakeScanRow{values: []any{
		"dept-1",
		"Engineering",
		"biz-1",
		sql.NullInt32{Valid: true, Int32: 7},
		sql.NullString{Valid: true, String: string(moneyJSON)},
		now,
		now,
	}})
	require.NoError(t, err)
	require.Equal(t, int32(7), dept.EmployeeNo)
	require.Equal(t, "Engineering", dept.Name)
	require.NotNil(t, dept.TotalPremium)

	emp, err := scanEmployee(fakeScanRow{values: []any{
		"emp-1",
		"Rahim",
		"E-1",
		"dept-1",
		"biz-1",
		sql.NullString{Valid: true, String: "life"},
		sql.NullString{Valid: true, String: "plan-1"},
		sql.NullString{Valid: true, String: string(moneyJSON)},
		sql.NullString{Valid: true, String: string(moneyJSON)},
		sql.NullString{Valid: true, String: "active"},
		now,
		now,
		sql.NullInt32{Valid: true, Int32: 2},
		sql.NullString{Valid: true, String: "rahim@example.com"},
		sql.NullString{Valid: true, String: "+8801712345678"},
		sql.NullString{Valid: true, String: "1990-01-01"},
		sql.NullString{Valid: true, String: "2024-01-01"},
		sql.NullString{Valid: true, String: "male"},
		sql.NullString{Valid: true, String: "user-1"},
	}})
	require.NoError(t, err)
	require.Equal(t, b2bv1.EmployeeStatus_EMPLOYEE_STATUS_ACTIVE, emp.Status)
	require.Equal(t, b2bv1.EmployeeGender_EMPLOYEE_GENDER_MALE, emp.Gender)
	require.Equal(t, commonv1.InsuranceType_INSURANCE_TYPE_LIFE, emp.InsuranceCategory)
	require.NotNil(t, emp.CoverageAmount)
	require.Equal(t, "user-1", emp.UserId)

	org, err := scanOrganisation(fakeScanRow{values: []any{
		"org-1",
		"tenant-1",
		"Acme",
		"ACME",
		"Tech",
		"admin@acme.test",
		"+8801",
		"Dhaka",
		sql.NullString{Valid: true, String: "active"},
		sql.NullInt32{Valid: true, Int32: 9},
		now,
		now,
	}})
	require.NoError(t, err)
	require.Equal(t, b2bv1.OrganisationStatus_ORGANISATION_STATUS_ACTIVE, org.Status)
	require.Equal(t, int32(9), org.TotalEmployees)

	member, err := scanOrgMember(fakeScanRow{values: []any{
		"member-1",
		"org-1",
		"user-1",
		sql.NullString{Valid: true, String: "business_admin"},
		sql.NullString{Valid: true, String: "active"},
		now,
		now,
		now,
	}})
	require.NoError(t, err)
	require.Equal(t, b2bv1.OrgMemberRole_ORG_MEMBER_ROLE_BUSINESS_ADMIN, member.Role)
	require.Equal(t, b2bv1.OrgMemberStatus_ORG_MEMBER_STATUS_ACTIVE, member.Status)
}

func TestRepositoryScanners_NotFound(t *testing.T) {
	noRows := errors.New("sql: no rows in result set")

	_, err := scanDepartment(fakeScanRow{err: noRows})
	require.ErrorIs(t, err, gorm.ErrRecordNotFound)

	_, err = scanEmployee(fakeScanRow{err: noRows})
	require.ErrorIs(t, err, gorm.ErrRecordNotFound)

	_, err = scanOrganisation(fakeScanRow{err: noRows})
	require.ErrorIs(t, err, gorm.ErrRecordNotFound)

	_, err = scanOrgMember(fakeScanRow{err: noRows})
	require.ErrorIs(t, err, gorm.ErrRecordNotFound)
}
