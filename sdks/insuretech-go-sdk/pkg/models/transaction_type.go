package models

// TransactionType represents a transaction_type
type TransactionType string

// TransactionType values
const (
	TransactionTypeTRANSACTIONTYPEUNSPECIFIED TransactionType = "TRANSACTION_TYPE_UNSPECIFIED"
	TransactionTypeTRANSACTIONTYPEPAYMENT  = "TRANSACTION_TYPE_PAYMENT"
	TransactionTypeTRANSACTIONTYPEREFUND  = "TRANSACTION_TYPE_REFUND"
)
