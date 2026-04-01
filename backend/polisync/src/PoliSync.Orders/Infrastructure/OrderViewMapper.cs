using Google.Protobuf.WellKnownTypes;
using Insuretech.Orders.Entity.V1;
using Insuretech.Orders.Services.V1;
using PoliSync.Orders.Application.Queries;
using DomainOrderPaymentStatus = PoliSync.Orders.Domain.OrderPaymentStatus;
using DomainOrderStatus = PoliSync.Orders.Domain.OrderStatus;

namespace PoliSync.Orders.Infrastructure;

internal static class OrderViewMapper
{
    public static OrderDto ToDto(OrderView orderView)
    {
        var order = orderView.Order;

        return new OrderDto(
            ParseGuid(order?.OrderId),
            order?.OrderNumber ?? string.Empty,
            ParseGuid(order?.QuotationId),
            ParseGuid(order?.CustomerId),
            ParseGuid(order?.ProductId),
            ParseGuid(order?.PlanId),
            ToDomainStatus(order?.Status ?? OrderStatus.Unspecified),
            order?.TotalPayable?.Amount ?? 0,
            string.IsNullOrWhiteSpace(order?.Currency) ? "BDT" : order.Currency,
            NullIfEmpty(order?.PaymentId),
            NullIfEmpty(order?.PaymentGatewayRef),
            ToDomainPaymentStatus(order?.PaymentStatus ?? OrderPaymentStatus.Unspecified),
            NullIfEmpty(order?.PolicyId),
            NullIfEmpty(order?.CancellationReason),
            NullIfEmpty(order?.FailureReason),
            ToNullableDateTime(order?.PaymentDueAt),
            ToNullableDateTime(order?.CoverageStartAt),
            ToNullableDateTime(order?.CoverageEndAt),
            ToDateTime(order?.CreatedAt),
            ToNullableDateTime(order?.UpdatedAt)
        );
    }

    public static OrderStatus ToProtoStatus(DomainOrderStatus status)
        => status switch
        {
            DomainOrderStatus.Pending => OrderStatus.Pending,
            DomainOrderStatus.PaymentInitiated => OrderStatus.PaymentInitiated,
            DomainOrderStatus.Paid => OrderStatus.Paid,
            DomainOrderStatus.PolicyIssued => OrderStatus.PolicyIssued,
            DomainOrderStatus.Cancelled => OrderStatus.Cancelled,
            DomainOrderStatus.Failed => OrderStatus.Failed,
            _ => OrderStatus.Unspecified
        };

    private static DomainOrderStatus ToDomainStatus(OrderStatus status)
        => status switch
        {
            OrderStatus.PaymentInitiated => DomainOrderStatus.PaymentInitiated,
            OrderStatus.Paid => DomainOrderStatus.Paid,
            OrderStatus.PolicyIssued => DomainOrderStatus.PolicyIssued,
            OrderStatus.Cancelled => DomainOrderStatus.Cancelled,
            OrderStatus.Failed => DomainOrderStatus.Failed,
            _ => DomainOrderStatus.Pending
        };

    private static DomainOrderPaymentStatus ToDomainPaymentStatus(OrderPaymentStatus status)
        => status switch
        {
            OrderPaymentStatus.PaymentInProgress => DomainOrderPaymentStatus.PaymentInProgress,
            OrderPaymentStatus.Paid => DomainOrderPaymentStatus.Paid,
            OrderPaymentStatus.PaymentFailed => DomainOrderPaymentStatus.Failed,
            _ => DomainOrderPaymentStatus.Unpaid
        };

    private static Guid ParseGuid(string? value)
        => Guid.TryParse(value, out var guid) ? guid : Guid.Empty;

    private static string? NullIfEmpty(string? value)
        => string.IsNullOrWhiteSpace(value) ? null : value;

    private static DateTime ToDateTime(Timestamp? timestamp)
        => timestamp?.ToDateTime() ?? DateTime.UtcNow;

    private static DateTime? ToNullableDateTime(Timestamp? timestamp)
        => timestamp is null || timestamp.Seconds <= 0 ? null : timestamp.ToDateTime();
}
