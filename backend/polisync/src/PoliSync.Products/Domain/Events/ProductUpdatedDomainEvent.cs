using PoliSync.SharedKernel.Domain;

namespace PoliSync.Products.Domain.Events;

/// <summary>
/// Raised when a product is updated.
/// </summary>
public record ProductUpdatedDomainEvent(
    Guid ProductId,
    Guid TenantId,
    string ProductCode,
    string ProductName,
    int Version,
    string UpdatedBy) : DomainEvent
{
    public override string EventType => nameof(ProductUpdatedDomainEvent);
}
