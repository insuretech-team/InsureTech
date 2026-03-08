using PoliSync.SharedKernel.Domain;

namespace PoliSync.SharedKernel.Messaging;

/// <summary>
/// Dispatches in-process domain events (MediatR INotification) after SaveChanges.
/// </summary>
public interface IDomainEventDispatcher
{
    Task DispatchAsync(IEnumerable<Entity> entities, CancellationToken ct = default);
}
