using PoliSync.SharedKernel.Domain;

namespace PoliSync.Orders.Domain;

public sealed record OrderCreatedEvent(
    Guid OrderId,
    string OrderNumber,
    Guid QuotationId,
    Guid CustomerId,
    long TotalPayable
) : DomainEvent;

public sealed record OrderPaymentInitiatedEvent(
    Guid OrderId,
    string OrderNumber,
    string PaymentId
) : DomainEvent;

public sealed record OrderPaymentConfirmedEvent(
    Guid OrderId,
    string OrderNumber,
    string PaymentId
) : DomainEvent;

public sealed record OrderPolicyIssuedEvent(
    Guid OrderId,
    string OrderNumber,
    string PolicyId
) : DomainEvent;

public sealed record OrderCancelledEvent(
    Guid OrderId,
    string OrderNumber,
    string Reason
) : DomainEvent;

public sealed record OrderFailedEvent(
    Guid OrderId,
    string OrderNumber,
    string Reason
) : DomainEvent;
