package models

// GracePeriodStatus represents a grace_period_status
type GracePeriodStatus string

// GracePeriodStatus values
const (
	GracePeriodStatusGRACEPERIODSTATUSUNSPECIFIED GracePeriodStatus = "GRACE_PERIOD_STATUS_UNSPECIFIED"
	GracePeriodStatusGRACEPERIODSTATUSACTIVE  = "GRACE_PERIOD_STATUS_ACTIVE"
	GracePeriodStatusGRACEPERIODSTATUSREVIVED  = "GRACE_PERIOD_STATUS_REVIVED"
	GracePeriodStatusGRACEPERIODSTATUSLAPSED  = "GRACE_PERIOD_STATUS_LAPSED"
)
