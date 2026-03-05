namespace PoliSync.SharedKernel.Messaging;

/// <summary>
/// Abstraction for publishing domain events to Kafka (or in-memory for dev).
/// </summary>
public interface IEventBus
{
    Task PublishAsync(string topic, object @event, CancellationToken ct = default);
}

/// <summary>
/// Dispatches domain events collected from entities after SaveChanges.
/// </summary>
public interface IDomainEventDispatcher
{
    Task DispatchAsync(IReadOnlyList<MediatR.INotification> events, CancellationToken ct = default);
}
