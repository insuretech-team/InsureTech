package models

// VehicleType represents a vehicle_type
type VehicleType string

// VehicleType values
const (
	VehicleTypeVEHICLETYPEUNSPECIFIED VehicleType = "VEHICLE_TYPE_UNSPECIFIED"
	VehicleTypeVEHICLETYPEBIKE  = "VEHICLE_TYPE_BIKE"
	VehicleTypeVEHICLETYPECAR  = "VEHICLE_TYPE_CAR"
	VehicleTypeVEHICLETYPECOMMERCIAL  = "VEHICLE_TYPE_COMMERCIAL"
	VehicleTypeVEHICLETYPETRUCK  = "VEHICLE_TYPE_TRUCK"
	VehicleTypeVEHICLETYPEBUS  = "VEHICLE_TYPE_BUS"
)
