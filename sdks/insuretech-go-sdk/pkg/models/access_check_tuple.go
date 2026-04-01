package models


// AccessCheckTuple represents a access_check_tuple
type AccessCheckTuple struct {
	Action string `json:"action,omitempty"`
	Object string `json:"object,omitempty"`
}
