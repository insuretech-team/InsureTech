package models

// ProductsProductStatus represents a products_product_status
type ProductsProductStatus string

// ProductsProductStatus values
const (
	ProductsProductStatusPRODUCTSTATUSUNSPECIFIED ProductsProductStatus = "PRODUCT_STATUS_UNSPECIFIED"
	ProductsProductStatusPRODUCTSTATUSDRAFT  = "PRODUCT_STATUS_DRAFT"
	ProductsProductStatusPRODUCTSTATUSACTIVE  = "PRODUCT_STATUS_ACTIVE"
	ProductsProductStatusPRODUCTSTATUSINACTIVE  = "PRODUCT_STATUS_INACTIVE"
	ProductsProductStatusPRODUCTSTATUSDISCONTINUED  = "PRODUCT_STATUS_DISCONTINUED"
)
