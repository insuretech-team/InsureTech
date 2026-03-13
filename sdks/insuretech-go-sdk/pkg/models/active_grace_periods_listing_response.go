package models


// ActiveGracePeriodsListingResponse represents a active_grace_periods_listing_response
type ActiveGracePeriodsListingResponse struct {
	GracePeriods []*GracePeriod `json:"grace_periods,omitempty"`
}
