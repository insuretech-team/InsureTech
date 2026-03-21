using System;
using InsuranceEngine.Policy.Domain.Enums;
using InsuranceEngine.SharedKernel.Interfaces;

namespace InsuranceEngine.Policy.Domain.Entities;

public class Endorsement
{
    public Guid Id { get; set; }

    public string EndorsementNumber { get; set; } = string.Empty;
    public Guid PolicyId { get; set; }
    public EndorsementType Type { get; set; }
    public string Reason { get; set; } = string.Empty;
    public string ChangesJson { get; set; } = "{}"; // Store before/after as JSON
    
    public decimal PremiumAdjustmentAmount { get; set; }
    public string PremiumAdjustmentCurrency { get; set; } = "BDT";
    public bool PremiumRefundRequired { get; set; }
    
    public EndorsementStatus Status { get; set; }
    public Guid RequestedBy { get; set; }
    public Guid? ApprovedBy { get; set; }
    
    public DateTime EffectiveDate { get; set; }
    public DateTime At { get; set; } // ApprovedAt
    public string AuditInfoJson { get; set; } = "{}";
    
    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }

    // Navigation (Optional)
    // public virtual PolicyEntity Policy { get; set; } = null!;
}
