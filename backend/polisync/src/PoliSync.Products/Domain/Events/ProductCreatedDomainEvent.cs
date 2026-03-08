using PoliSync.SharedKernel.Domain;

namespace PoliSync.Products.Domain.Events;

/// <summary>
/// Raised when a new product is created.
/// </summary>
public record ProductCreatedDomainEvent(
    Guid ProductId,
    Guid TenantId,
    Guid PartnerId,
    string ProductCode,
    string ProductName,
    string Category,
    long BasePremiumPaisa,
    string CreatedBy) : DomainEvent
{
    public override string EventType => nameof(ProductCreatedDomainEvent);
}
