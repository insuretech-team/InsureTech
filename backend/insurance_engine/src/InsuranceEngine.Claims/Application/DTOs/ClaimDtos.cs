using System;
using System.Collections.Generic;
using InsuranceEngine.Claims.Domain.Enums;

namespace InsuranceEngine.Claims.Application.DTOs;

public record ClaimResponseDto
{
    public Guid Id { get; init; }
    public string ClaimNumber { get; init; } = string.Empty;
    public Guid PolicyId { get; init; }
    public Guid CustomerId { get; init; }
    public string Status { get; init; } = string.Empty;
    public string ClaimType { get; init; } = string.Empty;
    public decimal ClaimedAmount { get; init; }
    public string Currency { get; init; } = "BDT";
    public DateTime IncidentDate { get; init; }
    public string IncidentDescription { get; init; } = string.Empty;
    public string PlaceOfIncident { get; init; } = string.Empty;
    public DateTime SubmittedAt { get; init; }
    public string? RejectionReason { get; init; }
    public List<ClaimApprovalDto> Approvals { get; init; } = new();
    public List<ClaimDocumentDto> Documents { get; init; } = new();
}

public record ClaimApprovalDto
{
    public Guid Id { get; init; }
    public string ApproverName { get; init; } = string.Empty;
    public string Role { get; init; } = string.Empty;
    public int Level { get; init; }
    public string Decision { get; init; } = string.Empty;
    public string? Notes { get; init; }
    public DateTime? DecidedAt { get; init; }
}

public record ClaimDocumentDto
{
    public Guid Id { get; init; }
    public string DocumentType { get; init; } = string.Empty;
    public string FileUrl { get; init; } = string.Empty;
    public bool IsVerified { get; init; }
    public DateTime UploadedAt { get; init; }
}

public record SubmitClaimRestRequest
{
    public Guid PolicyId { get; init; }
    public Guid CustomerId { get; init; }
    public ClaimType Type { get; init; }
    public decimal ClaimedAmount { get; init; }
    public DateTime IncidentDate { get; init; }
    public string IncidentDescription { get; init; } = string.Empty;
    public string PlaceOfIncident { get; init; } = string.Empty;
}
