using PoliSync.SharedKernel.Domain;

namespace PoliSync.Products.Domain.Events;

/// <summary>
/// Raised when pricing configuration is updated.
/// </summary>
public record PricingUpdatedDomainEvent(
    Guid PricingConfigId,
    Guid ProductId,
    Guid TenantId,
    int Version,
    string UpdatedBy) : DomainEvent
{
    public override string EventType => nameof(PricingUpdatedDomainEvent);
}
