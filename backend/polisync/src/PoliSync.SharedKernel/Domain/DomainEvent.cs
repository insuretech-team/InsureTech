namespace PoliSync.SharedKernel.Domain;

/// <summary>
/// Base class for domain events - something that happened in the domain
/// </summary>
public abstract record DomainEvent
{
    public Guid EventId { get; init; } = Guid.NewGuid();
    public DateTime OccurredAt { get; init; } = DateTime.UtcNow;
    public string EventType => GetType().Name;
}
