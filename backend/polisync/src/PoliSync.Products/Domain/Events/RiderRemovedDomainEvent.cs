using PoliSync.SharedKernel.Domain;

namespace PoliSync.Products.Domain.Events;

/// <summary>
/// Raised when a rider is removed from a product.
/// </summary>
public record RiderRemovedDomainEvent(
    Guid RiderId,
    Guid ProductId,
    Guid TenantId,
    string Reason,
    string RemovedBy) : DomainEvent
{
    public override string EventType => nameof(RiderRemovedDomainEvent);
}
