using Grpc.Core;
using Insuretech.Orders.Services.V1;
using PoliSync.Infrastructure.Clients;

namespace PoliSync.Orders.Infrastructure;

public sealed class GoOrderDataGateway : IOrderDataGateway
{
    private readonly OrderServiceGrpcClient _orderClient;

    public GoOrderDataGateway(OrderServiceGrpcClient orderClient)
    {
        _orderClient = orderClient;
    }

    public async Task<CreateOrderResponse> CreateOrderAsync(string quotationId, string customerId, string paymentMethod, CancellationToken cancellationToken = default)
    {
        return await _orderClient.Client.CreateOrderAsync(new CreateOrderRequest
        {
            QuotationId = quotationId,
            CustomerId = customerId,
            PaymentMethod = paymentMethod
        }, cancellationToken: cancellationToken);
    }

    public async Task<OrderView?> GetOrderAsync(string orderId, CancellationToken cancellationToken = default)
    {
        try
        {
            var response = await _orderClient.Client.GetOrderAsync(new GetOrderRequest
            {
                OrderId = orderId
            }, cancellationToken: cancellationToken);

            return response.Order;
        }
        catch (RpcException ex) when (ex.StatusCode == StatusCode.NotFound)
        {
            return null;
        }
    }

    public async Task<ListOrdersResponse> ListOrdersAsync(ListOrdersRequest request, CancellationToken cancellationToken = default)
    {
        return await _orderClient.Client.ListOrdersAsync(request, cancellationToken: cancellationToken);
    }

    public async Task<InitiatePaymentResponse> InitiatePaymentAsync(string orderId, string paymentMethod, string callbackUrl, string idempotencyKey, CancellationToken cancellationToken = default)
    {
        return await _orderClient.Client.InitiatePaymentAsync(new InitiatePaymentRequest
        {
            OrderId = orderId,
            PaymentMethod = paymentMethod,
            CallbackUrl = callbackUrl,
            IdempotencyKey = idempotencyKey
        }, cancellationToken: cancellationToken);
    }

    public async Task<ConfirmPaymentResponse> ConfirmPaymentAsync(string orderId, string paymentId, string transactionId, CancellationToken cancellationToken = default)
    {
        return await _orderClient.Client.ConfirmPaymentAsync(new ConfirmPaymentRequest
        {
            OrderId = orderId,
            PaymentId = paymentId,
            TransactionId = transactionId
        }, cancellationToken: cancellationToken);
    }

    public async Task<CancelOrderResponse> CancelOrderAsync(string orderId, string reason, CancellationToken cancellationToken = default)
    {
        return await _orderClient.Client.CancelOrderAsync(new CancelOrderRequest
        {
            OrderId = orderId,
            Reason = reason
        }, cancellationToken: cancellationToken);
    }

    public async Task<GetOrderStatusResponse?> GetOrderStatusAsync(string orderId, CancellationToken cancellationToken = default)
    {
        try
        {
            return await _orderClient.Client.GetOrderStatusAsync(new GetOrderStatusRequest
            {
                OrderId = orderId
            }, cancellationToken: cancellationToken);
        }
        catch (RpcException ex) when (ex.StatusCode == StatusCode.NotFound)
        {
            return null;
        }
    }
}
