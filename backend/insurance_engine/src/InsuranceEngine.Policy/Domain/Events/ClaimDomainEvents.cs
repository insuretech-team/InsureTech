using System;
using InsuranceEngine.Policy.Domain.Enums;
using InsuranceEngine.SharedKernel.Domain.Events;

namespace InsuranceEngine.Policy.Domain.Events;

public record ClaimSubmittedEvent(
    Guid ClaimId,
    string ClaimNumber,
    Guid PolicyId,
    Guid CustomerId,
    long Amount,
    string Currency,
    DateTime IncidentDate
) : DomainEvent;

public record ClaimProcessedEvent(
    Guid ClaimId,
    string ClaimNumber,
    ClaimStatus NewStatus,
    long? ApprovedAmount,
    string? Notes
) : DomainEvent;
