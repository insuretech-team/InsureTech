package models


// TransactionRetrievalResponse represents a transaction_retrieval_response
type TransactionRetrievalResponse struct {
	Transaction *MFSTransaction `json:"transaction,omitempty"`
	Error *Error `json:"error,omitempty"`
}
