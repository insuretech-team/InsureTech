package repository

import (
	"testing"

	partnerv1 "github.com/newage-saint/insuretech/gen/go/insuretech/partner/entity/v1"
	"github.com/stretchr/testify/assert"
)

func TestRepositoryHelpers(t *testing.T) {
	assert.Nil(t, nullableOptionalString(""))
	assert.Equal(t, "value", nullableOptionalString("value"))

	assert.Equal(t, "CORPORATE", partnerTypeToDBValue(partnerv1.PartnerType_PARTNER_TYPE_CORPORATE))
	assert.Equal(t, "ACTIVE", partnerStatusToDBValue(partnerv1.PartnerStatus_PARTNER_STATUS_ACTIVE))
	assert.Equal(t, "ACTIVE", agentStatusToDBValue(partnerv1.InsuranceAgentStatus_AGENT_STATUS_ACTIVE))
	assert.Equal(t, "ACQUISITION", commissionTypeToDBValue(partnerv1.CommissionType_COMMISSION_TYPE_ACQUISITION))
	assert.Equal(t, "PENDING", commissionStatusToDBValue(partnerv1.CommissionStatus_COMMISSION_STATUS_PENDING))

	mask := normalizeUpdateMask([]string{" organization_name ", "", "contact_phone"})
	assert.True(t, hasUpdateMaskKey(mask, "organization_name"))
	assert.True(t, hasUpdateMaskKey(mask, "contact_phone"))
	assert.False(t, hasUpdateMaskKey(mask, "bank_account"))

	assert.Equal(t, "created_at DESC", normalizePartnerOrderBy("-created_at"))
	assert.Equal(t, "updated_at ASC", normalizePartnerOrderBy("updated_at asc"))
	assert.Equal(t, "", normalizePartnerOrderBy("DROP TABLE"))

	assert.Equal(t, "ACTIVE", normalizePartnerStatusFilter("PARTNER_STATUS_ACTIVE"))
	assert.Equal(t, "ACTIVE", normalizePartnerStatusFilter("ACTIVE"))
	assert.Equal(t, "AGENT_NETWORK", normalizePartnerTypeFilter("PARTNER_TYPE_AGENT_NETWORK"))
	assert.Equal(t, "AGENT_NETWORK", normalizePartnerTypeFilter("AGENT_NETWORK"))
}
