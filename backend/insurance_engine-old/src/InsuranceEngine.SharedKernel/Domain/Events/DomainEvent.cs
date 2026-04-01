using System;
using System.ComponentModel.DataAnnotations.Schema;
using MediatR;

namespace InsuranceEngine.SharedKernel.Domain.Events;

/// <summary>
/// Base class for all domain events published to Kafka
/// </summary>
[NotMapped]
public abstract record DomainEvent : INotification
{
    public Guid EventId { get; init; } = Guid.NewGuid();
    public DateTime OccurredAt { get; init; } = DateTime.UtcNow;
}
