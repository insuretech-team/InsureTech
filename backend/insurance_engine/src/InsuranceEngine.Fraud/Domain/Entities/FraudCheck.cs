using System;
using System.Collections.Generic;
using InsuranceEngine.Fraud.Domain.Enums;

namespace InsuranceEngine.Fraud.Domain.Entities;

public class FraudCheck
{
    public Guid Id { get; set; }
    public Guid EntityId { get; set; } // ClaimId or PolicyId
    public string EntityType { get; set; } = string.Empty; // "Claim", "Policy"
    
    public FraudRiskLevel RiskLevel { get; set; }
    public FraudCheckStatus Status { get; set; }
    public double RiskScore { get; set; }
    
    public string? FindingsJson { get; set; } // Detailed reasons for flagging
    public string? CheckedRulesJson { get; set; }
    
    public DateTime CheckedAt { get; set; }
    public string? InspectorNotes { get; set; }
    
    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }
}
