namespace PoliSync.Quoting.Models;

public class Quote
{
    public string Id { get; set; } = Guid.NewGuid().ToString();
    public string QuoteNumber { get; set; } = string.Empty;
    public string ProductId { get; set; } = string.Empty;
    public string CustomerId { get; set; } = string.Empty;
    public string? AgentId { get; set; }
    
    public QuoteStatus Status { get; set; } = QuoteStatus.Draft;
    public DateTime CreatedAt { get; set; } = DateTime.UtcNow;
    public DateTime ValidUntil { get; set; }
    public DateTime? ConvertedAt { get; set; }
    public string? ConvertedPolicyId { get; set; }
    
    public QuoteParameters Parameters { get; set; } = new();
    public PremiumCalculation Premium { get; set; } = new();
    public List<Coverage> Coverages { get; set; } = [];
    public List<Discount> Discounts { get; set; } = [];
    public List<QuoteRevision> Revisions { get; set; } = [];
}

public class QuoteParameters
{
    public CoverageType CoverageType { get; set; }
    public string CoveragePlan { get; set; } = string.Empty;
    public decimal AssetValue { get; set; }
    public List<OptionalCoverage> OptionalCoverages { get; set; } = [];
    public int CoverageDurationMonths { get; set; }
    public Dictionary<string, object> AdditionalData { get; set; } = [];
}

public class PremiumCalculation
{
    public decimal BasePremium { get; set; }
    public decimal RiskAdjustment { get; set; }
    public decimal OptionalCoveragesTotal { get; set; }
    public decimal DiscountsTotal { get; set; }
    public decimal Taxes { get; set; }
    public decimal Fees { get; set; }
    public decimal TotalPremium { get; set; }
    public string Currency { get; set; } = "USD";
    public List<PremiumBreakdown> Breakdown { get; set; } = [];
}

public class PremiumBreakdown
{
    public string Category { get; set; } = string.Empty;
    public string Description { get; set; } = string.Empty;
    public decimal Amount { get; set; }
    public bool IsDiscount { get; set; }
}

public class Coverage
{
    public string CoverageId { get; set; } = string.Empty;
    public string Name { get; set; } = string.Empty;
    public string Description { get; set; } = string.Empty;
    public decimal Limit { get; set; }
    public decimal Deductible { get; set; }
    public decimal Premium { get; set; }
    public bool IsIncluded { get; set; }
    public bool IsOptional { get; set; }
}

public class OptionalCoverage
{
    public string CoverageId { get; set; } = string.Empty;
    public string Name { get; set; } = string.Empty;
    public decimal SelectedLimit { get; set; }
    public decimal SelectedDeductible { get; set; }
    decimal Premium { get; set; }
}

public class Discount
{
    public string DiscountId { get; set; } = string.Empty;
    public string Name { get; set; } = string.Empty;
    public string Description { get; set; } = string.Empty;
    public decimal Amount { get; set; }
    public decimal Percentage { get; set; }
    public string DiscountType { get; set; } = string.Empty;
}

public class QuoteRevision
{
    public int RevisionNumber { get; set; }
    public DateTime RevisedAt { get; set; }
    public string RevisedBy { get; set; } = string.Empty;
    public string ChangeReason { get; set; } = string.Empty;
    public PremiumCalculation PreviousPremium { get; set; } = new();
    public PremiumCalculation NewPremium { get; set; } = new();
}

public enum QuoteStatus
{
    Draft,
    Generated,
    Sent,
    Viewed,
    Accepted,
    Declined,
    Expired,
    Converted
}

public enum CoverageType
{
    FullCoverage,
    LiabilityOnly,
    Comprehensive,
    Collision,
    PersonalInjury,
    UninsuredMotorist
}
