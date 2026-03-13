using System;
using InsuranceEngine.Policy.Domain.Enums;

namespace InsuranceEngine.Policy.Domain.Entities;

public class ClaimApproval
{
    public Guid Id { get; set; }
    public Guid ClaimId { get; set; }
    public Guid ApproverId { get; set; }
    public string ApproverRole { get; set; } = string.Empty;
    public int ApprovalLevel { get; set; }
    public ApprovalDecision Decision { get; set; }
    public long ApprovedAmount { get; set; }
    public string ApprovedCurrency { get; set; } = "BDT";
    public string? Notes { get; set; }
    public DateTime ApprovedAt { get; set; }
    public DateTime CreatedAt { get; set; }
}
