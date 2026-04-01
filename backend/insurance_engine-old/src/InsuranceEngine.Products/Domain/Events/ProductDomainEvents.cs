using InsuranceEngine.SharedKernel.Domain.Events;

namespace InsuranceEngine.Products.Domain.Events;

public record ProductCreatedDomainEvent(
    Guid ProductId,
    string ProductCode,
    string ProductName,
    string Category,
    long BasePremium,
    Guid CreatedBy
) : DomainEvent;

public record ProductUpdatedDomainEvent(
    Guid ProductId,
    string ProductCode,
    List<string> UpdatedFields
) : DomainEvent;

public record ProductActivatedDomainEvent(
    Guid ProductId,
    string ProductCode,
    DateTime ActivatedAt
) : DomainEvent;

public record ProductDeactivatedDomainEvent(
    Guid ProductId,
    string ProductCode,
    string? Reason,
    DateTime DeactivatedAt
) : DomainEvent;

public record ProductDiscontinuedDomainEvent(
    Guid ProductId,
    string ProductCode,
    string? Reason,
    DateTime DiscontinuedAt
) : DomainEvent;
