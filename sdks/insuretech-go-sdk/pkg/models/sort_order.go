package models

// SortOrder represents a sort_order
type SortOrder string

// SortOrder values
const (
	SortOrderSORTORDERUNSPECIFIED SortOrder = "SORT_ORDER_UNSPECIFIED"
	SortOrderSORTORDERASC  = "SORT_ORDER_ASC"
	SortOrderSORTORDERDESC  = "SORT_ORDER_DESC"
)
