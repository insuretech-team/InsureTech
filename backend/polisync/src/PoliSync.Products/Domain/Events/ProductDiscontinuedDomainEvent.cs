using PoliSync.SharedKernel.Domain;

namespace PoliSync.Products.Domain.Events;

/// <summary>
/// Raised when a product is discontinued.
/// </summary>
public record ProductDiscontinuedDomainEvent(
    Guid ProductId,
    Guid TenantId,
    string ProductCode,
    string Reason,
    string DiscontinuedBy) : DomainEvent
{
    public override string EventType => nameof(ProductDiscontinuedDomainEvent);
}
