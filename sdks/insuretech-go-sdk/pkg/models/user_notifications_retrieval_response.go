package models


// UserNotificationsRetrievalResponse represents a user_notifications_retrieval_response
type UserNotificationsRetrievalResponse struct {
	Notifications []*Notification `json:"notifications,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
	UnreadCount int `json:"unread_count,omitempty"`
}
