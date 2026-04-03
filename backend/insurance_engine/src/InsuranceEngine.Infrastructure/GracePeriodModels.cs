namespace InsuranceEngine.Infrastructure.Renewals;

public class GracePeriodSettings
{
    public int GracePeriodDays { get; set; } = 30;
    public int ReinstatementWindowDays { get; set; } = 90;
    public decimal ReinstatementPenaltyPercent { get; set; } = 10.0m;
    public bool EnableDailyReminders { get; set; } = true;
}

public class GracePeriodInfo
{
    public string PolicyId { get; set; } = string.Empty;
    public string PolicyNumber { get; set; } = string.Empty;
    public string Status { get; set; } = string.Empty;
    public DateTime GracePeriodStartDate { get; set; }
    public DateTime GracePeriodEndDate { get; set; }
    public int DaysRemaining { get; set; }
    public decimal ReinstatementPenaltyAmount { get; set; }
    public bool CanReinstate { get; set; }
    public bool CanRenew { get; set; }
    public string NextAction { get; set; } = string.Empty;
}

public class ReinstatementRequest
{
    public string PolicyId { get; set; } = string.Empty;
    public int TenureMonths { get; set; } = 12;
    public bool ApplyReinstatementPenalty { get; set; } = true;
    public string? Notes { get; set; }
}

public class ReinstatementResult
{
    public bool Success { get; set; }
    public string? NewPolicyId { get; set; }
    public string? NewPolicyNumber { get; set; }
    public decimal ReinstatementPenalty { get; set; }
    public decimal TotalAmountDue { get; set; }
    public string? ErrorMessage { get; set; }
}
