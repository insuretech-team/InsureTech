package models

// OrgMemberStatus represents a org_member_status
type OrgMemberStatus string

// OrgMemberStatus values
const (
	OrgMemberStatusORGMEMBERSTATUSUNSPECIFIED OrgMemberStatus = "ORG_MEMBER_STATUS_UNSPECIFIED"
	OrgMemberStatusORGMEMBERSTATUSACTIVE  = "ORG_MEMBER_STATUS_ACTIVE"
	OrgMemberStatusORGMEMBERSTATUSINACTIVE  = "ORG_MEMBER_STATUS_INACTIVE"
)
