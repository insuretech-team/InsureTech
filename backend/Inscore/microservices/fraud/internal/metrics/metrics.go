package metrics

import "sync/atomic"

// RuntimeMetrics stores lightweight fraud service counters.
type RuntimeMetrics struct {
	fraudChecks      atomic.Int64
	fraudDetections  atomic.Int64
	alertsCreated    atomic.Int64
	casesCreated     atomic.Int64
	rulesActivated   atomic.Int64
	rulesDeactivated atomic.Int64
}

func NewRuntimeMetrics() *RuntimeMetrics {
	return &RuntimeMetrics{}
}

func (m *RuntimeMetrics) IncFraudChecks()      { m.fraudChecks.Add(1) }
func (m *RuntimeMetrics) IncFraudDetections()  { m.fraudDetections.Add(1) }
func (m *RuntimeMetrics) IncAlertsCreated()    { m.alertsCreated.Add(1) }
func (m *RuntimeMetrics) IncCasesCreated()     { m.casesCreated.Add(1) }
func (m *RuntimeMetrics) IncRulesActivated()   { m.rulesActivated.Add(1) }
func (m *RuntimeMetrics) IncRulesDeactivated() { m.rulesDeactivated.Add(1) }

func (m *RuntimeMetrics) Snapshot() map[string]int64 {
	return map[string]int64{
		"fraud_checks":      m.fraudChecks.Load(),
		"fraud_detections":  m.fraudDetections.Load(),
		"alerts_created":    m.alertsCreated.Load(),
		"cases_created":     m.casesCreated.Load(),
		"rules_activated":   m.rulesActivated.Load(),
		"rules_deactivated": m.rulesDeactivated.Load(),
	}
}
