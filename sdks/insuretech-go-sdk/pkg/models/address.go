package models


// Address represents a address
type Address struct {
	AddressLine2 string `json:"address_line2,omitempty"`
	City string `json:"city,omitempty"`
	Division string `json:"division,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
	Country string `json:"country,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
	AddressLine1 string `json:"address_line1,omitempty"`
	District string `json:"district,omitempty"`
	Latitude float64 `json:"latitude,omitempty"`
}
