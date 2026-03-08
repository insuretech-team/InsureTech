using PoliSync.SharedKernel.Domain;

namespace PoliSync.SharedKernel.Messaging;

/// <summary>
/// Event bus abstraction for publishing domain events to Kafka
/// </summary>
public interface IEventBus
{
    Task PublishAsync<TEvent>(
        TEvent @event, 
        string topic, 
        CancellationToken cancellationToken = default) 
        where TEvent : DomainEvent;
    
    Task PublishAsync<TEvent>(
        TEvent @event, 
        CancellationToken cancellationToken = default) 
        where TEvent : DomainEvent;
    
    Task PublishBatchAsync<TEvent>(
        IEnumerable<TEvent> events, 
        string topic, 
        CancellationToken cancellationToken = default) 
        where TEvent : DomainEvent;
}
