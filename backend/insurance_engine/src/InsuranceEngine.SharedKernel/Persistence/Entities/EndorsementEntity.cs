namespace InsuranceEngine.SharedKernel.Persistence.Entities;

public class EndorsementEntity
{
    public string EndorsementId { get; set; } = string.Empty;
    public string EndorsementNumber { get; set; } = string.Empty;
    public string PolicyId { get; set; } = string.Empty;
    public string Type { get; set; } = string.Empty;
    public string Status { get; set; } = string.Empty;
    public string Reason { get; set; } = string.Empty;
    public string Changes { get; set; } = string.Empty;
    public decimal? OldSumAssured { get; set; }
    public decimal? NewSumAssured { get; set; }
    public decimal? RefundAmount { get; set; }
    public decimal? AdditionalPremium { get; set; }
    public DateTime EffectiveDate { get; set; }
    public string ApprovedBy { get; set; } = string.Empty;
    public string RejectedBy { get; set; } = string.Empty;
    public string RejectionReason { get; set; } = string.Empty;
    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }
    
    public PolicyEntity? Policy { get; set; }
}
