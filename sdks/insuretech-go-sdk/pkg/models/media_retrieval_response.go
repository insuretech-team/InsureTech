package models


// MediaRetrievalResponse represents a media_retrieval_response
type MediaRetrievalResponse struct {
	Media *MediaFile `json:"media,omitempty"`
	Error *Error `json:"error,omitempty"`
}
