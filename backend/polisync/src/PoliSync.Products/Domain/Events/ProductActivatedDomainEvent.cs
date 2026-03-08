using PoliSync.SharedKernel.Domain;

namespace PoliSync.Products.Domain.Events;

/// <summary>
/// Raised when a product is activated.
/// </summary>
public record ProductActivatedDomainEvent(
    Guid ProductId,
    Guid TenantId,
    string ProductCode,
    string ActivatedBy) : DomainEvent
{
    public override string EventType => nameof(ProductActivatedDomainEvent);
}
