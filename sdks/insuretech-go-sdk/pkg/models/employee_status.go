package models

// EmployeeStatus represents a employee_status
type EmployeeStatus string

// EmployeeStatus values
const (
	EmployeeStatusEMPLOYEESTATUSUNSPECIFIED EmployeeStatus = "EMPLOYEE_STATUS_UNSPECIFIED"
	EmployeeStatusEMPLOYEESTATUSACTIVE  = "EMPLOYEE_STATUS_ACTIVE"
	EmployeeStatusEMPLOYEESTATUSINACTIVE  = "EMPLOYEE_STATUS_INACTIVE"
)
