package models

// ReserveType represents a reserve_type
type ReserveType string

// ReserveType values
const (
	ReserveTypeRESERVETYPEUNSPECIFIED ReserveType = "RESERVE_TYPE_UNSPECIFIED"
	ReserveTypeRESERVETYPECASE  = "RESERVE_TYPE_CASE"
	ReserveTypeRESERVETYPEIBNR  = "RESERVE_TYPE_IBNR"
	ReserveTypeRESERVETYPEIBNER  = "RESERVE_TYPE_IBNER"
	ReserveTypeRESERVETYPEEXPENSE  = "RESERVE_TYPE_EXPENSE"
	ReserveTypeRESERVETYPETOTAL  = "RESERVE_TYPE_TOTAL"
)
