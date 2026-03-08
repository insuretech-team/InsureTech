using Insuretech.Orders.Services.V1;

namespace PoliSync.Orders.Infrastructure;

public interface IOrderDataGateway
{
    Task<CreateOrderResponse> CreateOrderAsync(string quotationId, string customerId, string paymentMethod, CancellationToken cancellationToken = default);
    Task<OrderView?> GetOrderAsync(string orderId, CancellationToken cancellationToken = default);
    Task<ListOrdersResponse> ListOrdersAsync(ListOrdersRequest request, CancellationToken cancellationToken = default);
    Task<InitiatePaymentResponse> InitiatePaymentAsync(string orderId, string paymentMethod, string callbackUrl, string idempotencyKey, CancellationToken cancellationToken = default);
    Task<ConfirmPaymentResponse> ConfirmPaymentAsync(string orderId, string paymentId, string transactionId, CancellationToken cancellationToken = default);
    Task<CancelOrderResponse> CancelOrderAsync(string orderId, string reason, CancellationToken cancellationToken = default);
    Task<GetOrderStatusResponse?> GetOrderStatusAsync(string orderId, CancellationToken cancellationToken = default);
}
