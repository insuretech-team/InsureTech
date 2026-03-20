using System;
using System.Collections.Generic;

namespace InsuranceEngine.Claims.Domain.Entities;

/// <summary>
/// AI-powered fraud detection result for claims.
/// Maps to 'fraud_checks' table in insurance_schema.
/// Proto: insuretech.claims.entity.v1.FraudCheckResult
/// </summary>
public class FraudCheckResult
{
    public Guid Id { get; set; }
    public Guid ClaimId { get; set; }

    /// <summary>
    /// Fraud risk score (0-100, higher is more risky). Proto: fraud_score DECIMAL(5,2)
    /// </summary>
    public double FraudScore { get; set; }

    /// <summary>
    /// Array of detected risk factors. Proto: risk_factors TEXT[]
    /// </summary>
    public List<string> RiskFactors { get; set; } = new();

    /// <summary>
    /// Whether claim is flagged for manual review. Proto: flagged BOOLEAN default false
    /// </summary>
    public bool Flagged { get; set; }

    /// <summary>
    /// User ID who reviewed the flagged claim. Proto: reviewed_by UUID FK→users
    /// </summary>
    public Guid? ReviewedBy { get; set; }

    /// <summary>
    /// Manual review timestamp. Proto: reviewed_at TIMESTAMP WITH TIME ZONE
    /// </summary>
    public DateTime? ReviewedAt { get; set; }

    public DateTime CreatedAt { get; set; }
}
