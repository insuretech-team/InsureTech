namespace InsuranceEngine.SharedKernel.Persistence.Entities;

public class RefundEntity
{
    public string RefundId { get; set; } = string.Empty;
    public string RefundNumber { get; set; } = string.Empty;
    public string CancellationId { get; set; } = string.Empty;
    public string PolicyId { get; set; } = string.Empty;
    public string RefundType { get; set; } = string.Empty;
    public decimal RefundAmount { get; set; }
    public string RefundCurrency { get; set; } = "BDT";
    public string Status { get; set; } = string.Empty;
    public string PaymentMethod { get; set; } = string.Empty;
    public string PaymentReference { get; set; } = string.Empty;
    public DateTime? ProcessedAt { get; set; }
    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }
    
    public CancellationEntity? Cancellation { get; set; }
}
