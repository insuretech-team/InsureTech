package models


// DepartmentView represents a department_view
type DepartmentView struct {
	Department *Department `json:"department,omitempty"`
	OrganisationName string `json:"organisation_name,omitempty"`
}
