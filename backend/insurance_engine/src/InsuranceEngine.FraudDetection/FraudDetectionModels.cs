namespace InsuranceEngine.FraudDetection;

public class FraudCheckSettings
{
    public int RapidClaimHoursThreshold { get; set; } = 48;
    public int ClaimFrequencyThreshold { get; set; } = 2;
    public int ClaimFrequencyWindowMonths { get; set; } = 12;
    public decimal FullCoverageClaimThreshold { get; set; } = 1.0m;
    public int DeviceAccountThreshold { get; set; } = 3;
    public bool EnablePatternAnalysis { get; set; } = true;
    public bool EnableProviderValidation { get; set; } = true;
}

public class FraudCheckRequest
{
    public string EntityId { get; set; } = string.Empty;
    public string EntityType { get; set; } = string.Empty;
    public string? PolicyId { get; set; }
    public string? ClaimId { get; set; }
    public string? CustomerId { get; set; }
    public string? ProviderId { get; set; }
    public decimal? ClaimAmount { get; set; }
    public decimal? PolicyCoverageAmount { get; set; }
    public DateTime? PolicyPurchaseDate { get; set; }
    public DateTime? ClaimSubmissionDate { get; set; }
    public string? ClaimType { get; set; }
    public string? DeviceFingerprint { get; set; }
    public string? IpAddress { get; set; }
    public Dictionary<string, string>? Metadata { get; set; }
}

public class FraudCheckResult
{
    public bool IsFraudDetected { get; set; }
    public int FraudScore { get; set; }
    public string RiskLevel { get; set; } = "LOW";
    public List<FraudIndicator> Indicators { get; set; } = new();
    public string? Recommendation { get; set; }
    public DateTime CheckedAt { get; set; } = DateTime.UtcNow;
}

public class FraudIndicator
{
    public string Code { get; set; } = string.Empty;
    public string Description { get; set; } = string.Empty;
    public int ScoreContribution { get; set; }
    public string Severity { get; set; } = "LOW";
    public string? FrId { get; set; }
}

public class ClaimPatternAnalysis
{
    public string CustomerId { get; set; } = string.Empty;
    public int TotalClaimsInPeriod { get; set; }
    public int SameTypeClaimsInPeriod { get; set; }
    public decimal TotalClaimAmount { get; set; }
    public List<ClaimPattern> Patterns { get; set; } = new();
    public bool HasSuspiciousPattern { get; set; }
}

public class ClaimPattern
{
    public string ClaimType { get; set; } = string.Empty;
    public int Count { get; set; }
    public decimal TotalAmount { get; set; }
    public List<string> ClaimIds { get; set; } = new();
}

public class ProviderValidationResult
{
    public string ProviderId { get; set; } = string.Empty;
    public bool IsApproved { get; set; }
    public bool RequiresVerification { get; set; }
    public string? NetworkName { get; set; }
    public string? VerificationStatus { get; set; }
}

public class FraudDashboardSummary
{
    public int TotalFlagsToday { get; set; }
    public int HighRiskFlags { get; set; }
    public int MediumRiskFlags { get; set; }
    public int LowRiskFlags { get; set; }
    public int PendingReviewCount { get; set; }
    public List<FraudAlert> RecentAlerts { get; set; } = new();
    public Dictionary<string, int> FraudByType { get; set; } = new();
}

public class FraudAlert
{
    public string AlertId { get; set; } = string.Empty;
    public string EntityId { get; set; } = string.Empty;
    public string EntityType { get; set; } = string.Empty;
    public string AlertType { get; set; } = string.Empty;
    public string Severity { get; set; } = string.Empty;
    public string Status { get; set; } = "PENDING";
    public string Description { get; set; } = string.Empty;
    public DateTime CreatedAt { get; set; }
    public string? AssignedTo { get; set; }
}
