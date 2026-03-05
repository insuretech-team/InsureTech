namespace PoliSync.Underwriting.Domain;

public enum QuoteStatus
{
    Unspecified = 0,
    Draft = 1,
    PendingUnderwriting = 2,
    Approved = 3,
    Rejected = 4,
    Expired = 5,
    ConvertedToPolicy = 6,
    Withdrawn = 7
}

public enum DecisionType
{
    Unspecified = 0,
    Approved = 1,
    Rejected = 2,
    Referred = 3,
    ApprovedWithConditions = 4
}

public enum DecisionMethod
{
    Unspecified = 0,
    Automatic = 1,
    Manual = 2,
    Hybrid = 3
}

public enum RiskLevel
{
    Unspecified = 0,
    Low = 1,
    Medium = 2,
    High = 3,
    VeryHigh = 4,
    Declined = 5
}
