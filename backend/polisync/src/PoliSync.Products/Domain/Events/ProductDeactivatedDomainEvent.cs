using PoliSync.SharedKernel.Domain;

namespace PoliSync.Products.Domain.Events;

/// <summary>
/// Raised when a product is deactivated.
/// </summary>
public record ProductDeactivatedDomainEvent(
    Guid ProductId,
    Guid TenantId,
    string ProductCode,
    string Reason,
    string DeactivatedBy) : DomainEvent
{
    public override string EventType => nameof(ProductDeactivatedDomainEvent);
}
