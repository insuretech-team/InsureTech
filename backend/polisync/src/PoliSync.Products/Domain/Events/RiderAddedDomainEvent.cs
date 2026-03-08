using PoliSync.SharedKernel.Domain;

namespace PoliSync.Products.Domain.Events;

/// <summary>
/// Raised when a rider is added to a product.
/// </summary>
public record RiderAddedDomainEvent(
    Guid RiderId,
    Guid ProductId,
    Guid TenantId,
    string RiderCode,
    string RiderName,
    string AddedBy) : DomainEvent
{
    public override string EventType => nameof(RiderAddedDomainEvent);
}
