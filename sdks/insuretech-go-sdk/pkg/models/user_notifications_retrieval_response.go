package models


// UserNotificationsRetrievalResponse represents a user_notifications_retrieval_response
type UserNotificationsRetrievalResponse struct {
	TotalCount int `json:"total_count,omitempty"`
	UnreadCount int `json:"unread_count,omitempty"`
	Error *Error `json:"error,omitempty"`
	Notifications []*Notification `json:"notifications,omitempty"`
}
