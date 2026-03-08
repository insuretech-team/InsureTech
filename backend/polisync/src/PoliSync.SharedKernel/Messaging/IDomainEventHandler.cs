using MediatR;
using PoliSync.SharedKernel.Domain;

namespace PoliSync.SharedKernel.Messaging;

/// <summary>
/// Handler for domain events (in-process)
/// </summary>
public interface IDomainEventHandler<in TEvent> : INotificationHandler<TEvent>
    where TEvent : DomainEvent, INotification
{
}
