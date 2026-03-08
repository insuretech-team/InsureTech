using MediatR;
using PoliSync.SharedKernel.Domain;
using PoliSync.SharedKernel.Messaging;

namespace PoliSync.Infrastructure.Messaging;

/// <summary>
/// Dispatches in-memory domain events via MediatR after EF SaveChanges.
/// </summary>
public sealed class DomainEventDispatcher : IDomainEventDispatcher
{
    private readonly IMediator _mediator;

    public DomainEventDispatcher(IMediator mediator) => _mediator = mediator;

    public async Task DispatchAsync(IEnumerable<Entity> entities, CancellationToken ct = default)
    {
        var events = entities
            .SelectMany(e => e.DomainEvents)
            .ToList();

        foreach (var entity in entities)
            entity.ClearDomainEvents();

        foreach (var @event in events)
            await _mediator.Publish(@event, ct);
    }
}
