namespace InsuranceEngine.SharedKernel.Persistence.Entities;

public class CancellationEntity
{
    public string CancellationId { get; set; } = string.Empty;
    public string CancellationNumber { get; set; } = string.Empty;
    public string PolicyId { get; set; } = string.Empty;
    public string CancellationType { get; set; } = string.Empty;
    public string Reason { get; set; } = string.Empty;
    public string Status { get; set; } = string.Empty;
    public decimal? RefundAmount { get; set; }
    public string RefundStatus { get; set; } = string.Empty;
    public DateTime EffectiveDate { get; set; }
    public string ApprovedBy { get; set; } = string.Empty;
    public string RejectedBy { get; set; } = string.Empty;
    public string RejectionReason { get; set; } = string.Empty;
    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }
    
    public PolicyEntity? Policy { get; set; }
}
