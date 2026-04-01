package models


// UserNotificationsRetrievalRequest represents a user_notifications_retrieval_request
type UserNotificationsRetrievalRequest struct {
	Limit int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
	UnreadOnly bool `json:"unread_only,omitempty"`
	UserId string `json:"user_id"`
}
