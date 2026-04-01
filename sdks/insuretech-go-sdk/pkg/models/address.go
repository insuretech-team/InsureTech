package models


// Address represents a address
type Address struct {
	AddressLine1 string `json:"address_line1,omitempty"`
	AddressLine2 string `json:"address_line2,omitempty"`
	City string `json:"city,omitempty"`
	Country string `json:"country,omitempty"`
	District string `json:"district,omitempty"`
	Division string `json:"division,omitempty"`
	Latitude float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
}
