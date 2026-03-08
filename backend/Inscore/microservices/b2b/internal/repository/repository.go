// Package repository provides data access implementations for the B2B microservice.
// All repositories use the proto-generated GORM-tagged structs directly.
//
// Pattern (matches authn service):
//   - Simple scalar fields: GORM Find/First with proto struct directly (GORM tags handle mapping)
//   - Timestamp fields: serializer:proto_timestamp registered in db package handles *timestamppb.Timestamp
//   - Enum fields: serializer:proto_enum registered in db package handles proto enum ↔ DB string
//   - Money fields (*v1.Money, stored as JSONB): raw SQL scanning via sql.NullString + json.Unmarshal
//   - Create/Update: map[string]interface{} for partial updates; proto struct for full inserts
//
// The B2BRepository interface is satisfied by composing all sub-repositories
// into a single PortalRepository that is wired in main.go.
package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	b2bv1 "github.com/newage-saint/insuretech/gen/go/insuretech/b2b/entity/v1"
	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	// Ensure proto_timestamp + proto_enum serializers are registered
	_ "github.com/newage-saint/insuretech/backend/inscore/db"
)

// PortalRepository satisfies domain.B2BRepository by composing all sub-repos.
type PortalRepository struct {
	db *gorm.DB
}

func NewPortalRepository(db *gorm.DB) *PortalRepository {
	return &PortalRepository{db: db}
}

// ─── Money helpers ────────────────────────────────────────────────────────────

// protoJSONMarshaler marshals proto messages to JSON using snake_case field names
// (consistent with protojson wire format used by the gateway).
var protoJSONMarshaler = protojson.MarshalOptions{UseProtoNames: true, EmitUnpopulated: false}
var protoJSONUnmarshaler = protojson.UnmarshalOptions{DiscardUnknown: true}

// scanMoney reads a JSONB money column from a *sql.NullString into *commonv1.Money.
// Uses protojson so field names (snake_case) match what was written by marshalMoney.
func scanMoney(raw sql.NullString) *commonv1.Money {
	if !raw.Valid || strings.TrimSpace(raw.String) == "" || raw.String == "null" {
		return nil
	}
	var m commonv1.Money
	if err := protoJSONUnmarshaler.Unmarshal([]byte(raw.String), &m); err != nil {
		// Fallback: try standard json (for rows written before this fix)
		_ = json.Unmarshal([]byte(raw.String), &m)
	}
	return &m
}

// marshalMoney encodes *commonv1.Money to JSON bytes for DB insert/update.
// Uses protojson with UseProtoNames=true so stored JSON uses snake_case keys
// (e.g. "decimal_amount" not "decimalAmount") consistent with protojson wire format.
func marshalMoney(m *commonv1.Money) ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	b, err := protoJSONMarshaler.Marshal(m)
	return b, err
}

// zeroMoneyJSON returns a zero BDT money JSON blob for new records.
func zeroMoneyJSON() []byte {
	b, _ := protoJSONMarshaler.Marshal(&commonv1.Money{Amount: 0, Currency: "BDT", DecimalAmount: 0})
	return b
}

// ─── Enum string helpers ──────────────────────────────────────────────────────

func employeeStatusStr(s b2bv1.EmployeeStatus) string {
	switch s {
	case b2bv1.EmployeeStatus_EMPLOYEE_STATUS_ACTIVE:
		return "EMPLOYEE_STATUS_ACTIVE"
	case b2bv1.EmployeeStatus_EMPLOYEE_STATUS_INACTIVE:
		return "EMPLOYEE_STATUS_INACTIVE"
	default:
		return "EMPLOYEE_STATUS_ACTIVE"
	}
}

func employeeGenderStr(g b2bv1.EmployeeGender) string {
	switch g {
	case b2bv1.EmployeeGender_EMPLOYEE_GENDER_MALE:
		return "MALE"
	case b2bv1.EmployeeGender_EMPLOYEE_GENDER_FEMALE:
		return "FEMALE"
	case b2bv1.EmployeeGender_EMPLOYEE_GENDER_OTHER:
		return "OTHER"
	default:
		return ""
	}
}

func organisationStatusStr(s b2bv1.OrganisationStatus) string {
	if s == b2bv1.OrganisationStatus_ORGANISATION_STATUS_UNSPECIFIED {
		return "ORGANISATION_STATUS_ACTIVE"
	}
	return s.String()
}

func orgMemberRoleStr(r b2bv1.OrgMemberRole) string {
	if r == b2bv1.OrgMemberRole_ORG_MEMBER_ROLE_UNSPECIFIED {
		return "ORG_MEMBER_ROLE_HR_MANAGER"
	}
	return r.String()
}

func purchaseOrderStatusStr(s b2bv1.PurchaseOrderStatus) string {
	if s == b2bv1.PurchaseOrderStatus_PURCHASE_ORDER_STATUS_UNSPECIFIED {
		return "PURCHASE_ORDER_STATUS_SUBMITTED"
	}
	return s.String()
}

// parseInsuranceType maps DB string → commonv1.InsuranceType.
func parseInsuranceType(value string) commonv1.InsuranceType {
	value = strings.TrimSpace(value)
	if value == "" {
		return commonv1.InsuranceType_INSURANCE_TYPE_UNSPECIFIED
	}
	if v, ok := commonv1.InsuranceType_value[value]; ok {
		return commonv1.InsuranceType(v)
	}
	normalized := "INSURANCE_TYPE_" + strings.ToUpper(value)
	if v, ok := commonv1.InsuranceType_value[normalized]; ok {
		return commonv1.InsuranceType(v)
	}
	return commonv1.InsuranceType_INSURANCE_TYPE_UNSPECIFIED
}

// ─── Shared time helpers ──────────────────────────────────────────────────────

func toProtoTS(t time.Time) *timestamppb.Timestamp {
	if t.IsZero() {
		return nil
	}
	return timestamppb.New(t)
}

// ─── Ensure uuid import used ──────────────────────────────────────────────────

func newUUID() string { return uuid.NewString() }

// ─── Ensure context import used ───────────────────────────────────────────────
var _ = context.Background
var _ = fmt.Sprintf
