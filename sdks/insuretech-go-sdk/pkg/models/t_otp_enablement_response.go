package models


// TOTPEnablementResponse represents a t_otp_enablement_response
type TOTPEnablementResponse struct {
	ProvisioningUri string `json:"provisioning_uri,omitempty"`
	TotpSecret string `json:"totp_secret,omitempty"`
}
