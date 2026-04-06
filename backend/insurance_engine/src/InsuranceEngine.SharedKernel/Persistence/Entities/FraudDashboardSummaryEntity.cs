namespace InsuranceEngine.SharedKernel.Persistence.Entities;

public class FraudDashboardSummaryEntity
{
    public string SummaryId { get; set; } = string.Empty;
    public DateTime SummaryDate { get; set; }
    public int TotalFlagsToday { get; set; }
    public int HighRiskFlags { get; set; }
    public int MediumRiskFlags { get; set; }
    public int LowRiskFlags { get; set; }
    public int PendingReviewCount { get; set; }
    public int ResolvedCount { get; set; }
    public decimal AverageFraudScore { get; set; }
    public string TopFraudTypes { get; set; } = string.Empty;
    public DateTime GeneratedAt { get; set; }
}
