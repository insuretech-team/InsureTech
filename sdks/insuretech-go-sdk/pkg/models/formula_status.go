package models

// FormulaStatus represents a formula_status
type FormulaStatus string

// FormulaStatus values
const (
	FormulaStatusFORMULASTATUSUNSPECIFIED FormulaStatus = "FORMULA_STATUS_UNSPECIFIED"
	FormulaStatusFORMULASTATUSDRAFT  = "FORMULA_STATUS_DRAFT"
	FormulaStatusFORMULASTATUSACTIVE  = "FORMULA_STATUS_ACTIVE"
	FormulaStatusFORMULASTATUSDEPRECATED  = "FORMULA_STATUS_DEPRECATED"
	FormulaStatusFORMULASTATUSRETIRED  = "FORMULA_STATUS_RETIRED"
)
