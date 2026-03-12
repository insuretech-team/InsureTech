using InsuranceEngine.SharedKernel.Domain.Events;

namespace InsuranceEngine.Policy.Domain.Events;

public record PolicyCreatedEvent(
    Guid PolicyId,
    string PolicyNumber,
    Guid CustomerId,
    Guid ProductId,
    long PremiumAmount,
    DateTime StartDate,
    DateTime EndDate
) : DomainEvent;

public record PolicyIssuedEvent(
    Guid PolicyId,
    string PolicyNumber,
    DateTime IssuedAt
) : DomainEvent;

public record PolicyRenewedEvent(
    Guid OldPolicyId,
    Guid NewPolicyId,
    string NewPolicyNumber,
    DateTime RenewalDate
) : DomainEvent;

public record PolicyCancelledEvent(
    Guid PolicyId,
    string PolicyNumber,
    DateTime CancelledAt,
    string Reason
) : DomainEvent;

public record PolicyEndorsedEvent(
    Guid PolicyId,
    string PolicyNumber,
    List<string> ChangedFields,
    DateTime EndorsedAt
) : DomainEvent;
