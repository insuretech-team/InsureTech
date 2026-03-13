package models


// NotificationStatusRetrievalResponse represents a notification_status_retrieval_response
type NotificationStatusRetrievalResponse struct {
	Notification *Notification `json:"notification,omitempty"`
	Error *Error `json:"error,omitempty"`
}
