namespace InsuranceEngine.Policy.Domain.Enums;

public enum ClaimStatus
{
    Unspecified = 0,
    Submitted = 1,         // Initial submission
    UnderReview = 2,      // Being reviewed
    PendingDocuments = 3, // Awaiting documents
    Approved = 4,          // Approved for payment
    Rejected = 5,          // Rejected
    Settled = 6,           // Payment completed
    Disputed = 7,          // Under dispute
    FraudCheck = 8         // Flagged for fraud check
}
