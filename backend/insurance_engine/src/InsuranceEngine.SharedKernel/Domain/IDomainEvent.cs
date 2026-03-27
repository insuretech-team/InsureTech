using MediatR;

namespace InsuranceEngine.SharedKernel.Domain;

public interface IDomainEvent : INotification
{
    DateTime OccurredOn { get; }
}

public abstract record DomainEvent : IDomainEvent
{
    public DateTime OccurredOn { get; } = DateTime.UtcNow;
}
