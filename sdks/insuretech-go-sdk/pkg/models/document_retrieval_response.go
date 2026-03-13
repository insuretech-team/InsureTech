package models


// DocumentRetrievalResponse represents a document_retrieval_response
type DocumentRetrievalResponse struct {
	Document *DocumentGeneration `json:"document,omitempty"`
	Error *Error `json:"error,omitempty"`
}
