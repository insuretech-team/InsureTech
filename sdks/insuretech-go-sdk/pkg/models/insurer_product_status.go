package models

// InsurerProductStatus represents a insurer_product_status
type InsurerProductStatus string

// InsurerProductStatus values
const (
	InsurerProductStatusPRODUCTSTATUSUNSPECIFIED InsurerProductStatus = "PRODUCT_STATUS_UNSPECIFIED"
	InsurerProductStatusPRODUCTSTATUSACTIVE  = "PRODUCT_STATUS_ACTIVE"
	InsurerProductStatusPRODUCTSTATUSINACTIVE  = "PRODUCT_STATUS_INACTIVE"
	InsurerProductStatusPRODUCTSTATUSDISCONTINUED  = "PRODUCT_STATUS_DISCONTINUED"
)
