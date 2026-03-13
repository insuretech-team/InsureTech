package models


// DocumentTemplateRetrievalResponse represents a document_template_retrieval_response
type DocumentTemplateRetrievalResponse struct {
	Template *DocumentTemplate `json:"template,omitempty"`
	Error *Error `json:"error,omitempty"`
}
