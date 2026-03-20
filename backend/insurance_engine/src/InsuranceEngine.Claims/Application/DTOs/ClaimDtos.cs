using System;
using System.Collections.Generic;
using InsuranceEngine.Claims.Domain.Enums;

namespace InsuranceEngine.Claims.Application.DTOs;

public record MoneyDto(long Amount, string CurrencyCode = "BDT");

public record ClaimResponseDto
{
    public Guid Id { get; init; }
    public string ClaimNumber { get; init; } = string.Empty;
    public Guid PolicyId { get; init; }
    public Guid CustomerId { get; init; }
    public ClaimType Type { get; init; }
    public ClaimStatus Status { get; init; }
    public ClaimProcessingType ProcessingType { get; init; }
    public MoneyDto ClaimedAmount { get; init; } = new(0);
    public MoneyDto ApprovedAmount { get; init; } = new(0);
    public MoneyDto SettledAmount { get; init; } = new(0);
    public MoneyDto DeductibleAmount { get; init; } = new(0);
    public MoneyDto CoPayAmount { get; init; } = new(0);
    public DateTime IncidentDate { get; init; }
    public string IncidentDescription { get; init; } = string.Empty;
    public string? PlaceOfIncident { get; init; }
    public DateTime SubmittedAt { get; init; }
    public DateTime? ApprovedAt { get; init; }
    public DateTime? SettledAt { get; init; }
    public string? RejectionReason { get; init; }
    public bool AppealOptionAvailable { get; init; }
    public FraudCheckResultDto? FraudCheck { get; init; }
    public List<ClaimApprovalDto> Approvals { get; init; } = new();
    public List<ClaimDocumentDto> Documents { get; init; } = new();
    public DateTime CreatedAt { get; init; }
    public DateTime UpdatedAt { get; init; }
}

public record ClaimListDto
{
    public Guid Id { get; init; }
    public string ClaimNumber { get; init; } = string.Empty;
    public Guid PolicyId { get; init; }
    public ClaimType Type { get; init; }
    public ClaimStatus Status { get; init; }
    public MoneyDto ClaimedAmount { get; init; } = new(0);
    public MoneyDto ApprovedAmount { get; init; } = new(0);
    public DateTime SubmittedAt { get; init; }
}

public record ClaimApprovalDto
{
    public Guid Id { get; init; }
    public Guid ApproverId { get; init; }
    public string ApproverRole { get; init; } = string.Empty;
    public int ApprovalLevel { get; init; }
    public ApprovalDecision Decision { get; init; }
    public MoneyDto ApprovedAmount { get; init; } = new(0);
    public string? Notes { get; init; }
    public DateTime ApprovedAt { get; init; }
}

public record ClaimDocumentDto
{
    public Guid Id { get; init; }
    public string DocumentType { get; init; } = string.Empty;
    public string FileUrl { get; init; } = string.Empty;
    public string FileHash { get; init; } = string.Empty;
    public bool Verified { get; init; }
    public DateTime UploadedAt { get; init; }
}

public record FraudCheckResultDto
{
    public Guid Id { get; init; }
    public double FraudScore { get; init; }
    public List<string> RiskFactors { get; init; } = new();
    public bool Flagged { get; init; }
    public Guid? ReviewedBy { get; init; }
    public DateTime? ReviewedAt { get; init; }
}

public record SubmitClaimRestRequest
{
    public Guid PolicyId { get; init; }
    public Guid CustomerId { get; init; }
    public ClaimType Type { get; init; }
    public long ClaimedAmount { get; init; }
    public DateTime IncidentDate { get; init; }
    public string IncidentDescription { get; init; } = string.Empty;
    public string? PlaceOfIncident { get; init; }
    public string? BankDetailsForPayout { get; init; }
}
