using System;

namespace InsuranceEngine.Fraud.Domain.Enums;

public enum FraudRiskLevel
{
    Low = 0,
    Medium = 1,
    High = 2,
    Critical = 3
}

public enum FraudCheckStatus
{
    Pending = 0,
    Approved = 1,
    Flagged = 2,
    Rejected = 3
}
