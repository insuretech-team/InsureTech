package models

// OrgMemberRole represents a org_member_role
type OrgMemberRole string

// OrgMemberRole values
const (
	OrgMemberRoleORGMEMBERROLEUNSPECIFIED OrgMemberRole = "ORG_MEMBER_ROLE_UNSPECIFIED"
	OrgMemberRoleORGMEMBERROLEBUSINESSADMIN  = "ORG_MEMBER_ROLE_BUSINESS_ADMIN"
	OrgMemberRoleORGMEMBERROLEHRMANAGER  = "ORG_MEMBER_ROLE_HR_MANAGER"
	OrgMemberRoleORGMEMBERROLEVIEWER  = "ORG_MEMBER_ROLE_VIEWER"
)
