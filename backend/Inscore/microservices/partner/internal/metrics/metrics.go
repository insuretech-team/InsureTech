package metrics

import "sync/atomic"

// RuntimeMetrics stores lightweight counters for partner service operations.
type RuntimeMetrics struct {
	partnerCreated      atomic.Int64
	partnerFetched      atomic.Int64
	partnerUpdated      atomic.Int64
	partnerListed       atomic.Int64
	partnerDeleted      atomic.Int64
	partnerVerified     atomic.Int64
	partnerStatusUpdate atomic.Int64
	commissionCalculated atomic.Int64
	apiKeyRotated       atomic.Int64
}

func NewRuntimeMetrics() *RuntimeMetrics { return &RuntimeMetrics{} }

func (m *RuntimeMetrics) IncPartnerCreated()       { m.partnerCreated.Add(1) }
func (m *RuntimeMetrics) IncPartnerFetched()       { m.partnerFetched.Add(1) }
func (m *RuntimeMetrics) IncPartnerUpdated()       { m.partnerUpdated.Add(1) }
func (m *RuntimeMetrics) IncPartnerListed()        { m.partnerListed.Add(1) }
func (m *RuntimeMetrics) IncPartnerDeleted()       { m.partnerDeleted.Add(1) }
func (m *RuntimeMetrics) IncPartnerVerified()      { m.partnerVerified.Add(1) }
func (m *RuntimeMetrics) IncPartnerStatusUpdate()  { m.partnerStatusUpdate.Add(1) }
func (m *RuntimeMetrics) IncCommissionCalculated() { m.commissionCalculated.Add(1) }
func (m *RuntimeMetrics) IncAPIKeyRotated()        { m.apiKeyRotated.Add(1) }

func (m *RuntimeMetrics) Snapshot() map[string]int64 {
	return map[string]int64{
		"partner_created":        m.partnerCreated.Load(),
		"partner_fetched":        m.partnerFetched.Load(),
		"partner_updated":        m.partnerUpdated.Load(),
		"partner_listed":         m.partnerListed.Load(),
		"partner_deleted":        m.partnerDeleted.Load(),
		"partner_verified":       m.partnerVerified.Load(),
		"partner_status_updates": m.partnerStatusUpdate.Load(),
		"commission_calculated":  m.commissionCalculated.Load(),
		"api_key_rotated":        m.apiKeyRotated.Load(),
	}
}
