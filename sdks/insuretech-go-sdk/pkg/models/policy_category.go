package models

// PolicyCategory represents a policy_category
type PolicyCategory string

// PolicyCategory values
const (
	PolicyCategoryPOLICYCATEGORYUNSPECIFIED PolicyCategory = "POLICY_CATEGORY_UNSPECIFIED"
	PolicyCategoryPOLICYCATEGORYLIFE  = "POLICY_CATEGORY_LIFE"
	PolicyCategoryPOLICYCATEGORYNONLIFE  = "POLICY_CATEGORY_NONLIFE"
)
