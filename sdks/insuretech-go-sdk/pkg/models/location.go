package models


// Location represents a location
type Location struct {
	Accuracy float64 `json:"accuracy,omitempty"`
	Altitude float64 `json:"altitude,omitempty"`
	Latitude float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
}
