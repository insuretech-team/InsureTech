using MediatR;

namespace PoliSync.SharedKernel.Domain;

/// <summary>
/// Base class for all aggregate roots. Holds domain events for dispatch after persistence.
/// </summary>
public abstract class Entity
{
    private readonly List<INotification> _domainEvents = [];

    public IReadOnlyList<INotification> DomainEvents => _domainEvents.AsReadOnly();

    protected void RaiseDomainEvent(INotification domainEvent)
        => _domainEvents.Add(domainEvent);

    public void ClearDomainEvents() => _domainEvents.Clear();
}
