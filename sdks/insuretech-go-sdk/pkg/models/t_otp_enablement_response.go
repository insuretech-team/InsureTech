package models


// TOTPEnablementResponse represents a t_otp_enablement_response
type TOTPEnablementResponse struct {
	Error *Error `json:"error,omitempty"`
	TotpSecret string `json:"totp_secret,omitempty"`
	ProvisioningUri string `json:"provisioning_uri,omitempty"`
}
