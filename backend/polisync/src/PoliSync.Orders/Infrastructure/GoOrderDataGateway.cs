using Grpc.Core;
using Insuretech.Orders.Services.V1;
using Microsoft.AspNetCore.Http;
using PoliSync.Infrastructure.Clients;
using PoliSync.SharedKernel.Auth;

namespace PoliSync.Orders.Infrastructure;

public sealed class GoOrderDataGateway : IOrderDataGateway
{
    private readonly OrderServiceGrpcClient _orderClient;
    private readonly ICurrentUser _currentUser;
    private readonly IHttpContextAccessor _httpContextAccessor;

    public GoOrderDataGateway(OrderServiceGrpcClient orderClient, ICurrentUser currentUser, IHttpContextAccessor httpContextAccessor)
    {
        _orderClient = orderClient;
        _currentUser = currentUser;
        _httpContextAccessor = httpContextAccessor;
    }

    /// <summary>
    /// Builds gRPC CallOptions with identity headers forwarded from the incoming HTTP request.
    /// The Go orders service requires x-user-id, x-tenant-id etc. in gRPC metadata.
    /// </summary>
    private CallOptions BuildCallOptions(CancellationToken ct = default)
    {
        var headers = new Metadata();
        var httpCtx = _httpContextAccessor.HttpContext;
        if (httpCtx != null)
        {
            // Forward all X-* identity headers injected by the Go gateway
            foreach (var key in new[] { "x-user-id", "x-tenant-id", "x-partner-id", "x-token-id",
                                         "x-user-type", "x-portal", "x-roles", "x-request-id", "x-session-id" })
            {
                var val = httpCtx.Request.Headers[key].FirstOrDefault();
                if (!string.IsNullOrEmpty(val)) headers.Add(key, val);
            }
        }
        else if (_currentUser.UserId != Guid.Empty)
        {
            // Fallback: populate from ICurrentUser when HttpContext is unavailable
            headers.Add("x-user-id", _currentUser.UserId.ToString());
            if (_currentUser.TenantId != Guid.Empty)
                headers.Add("x-tenant-id", _currentUser.TenantId.ToString());
        }
        return new CallOptions(headers: headers, cancellationToken: ct);
    }

    public async Task<CreateOrderResponse> CreateOrderAsync(
        string quotationId,
        string customerId,
        string paymentMethod,
        CancellationToken cancellationToken = default,
        string? productId = null,
        string? planId = null,
        long totalPayable = 0,
        string currency = "BDT")
    {
        var req = new CreateOrderRequest
        {
            QuotationId = quotationId,
            CustomerId = customerId,
            PaymentMethod = paymentMethod,
        };
        if (!string.IsNullOrEmpty(productId)) req.ProductId = productId;
        if (!string.IsNullOrEmpty(planId)) req.PlanId = planId;
        if (totalPayable > 0)
            req.TotalPayable = new Insuretech.Common.V1.Money { Amount = totalPayable, Currency = currency };
        return await _orderClient.Client.CreateOrderAsync(req, BuildCallOptions(cancellationToken));
    }

    public async Task<OrderView?> GetOrderAsync(string orderId, CancellationToken cancellationToken = default)
    {
        try
        {
            var response = await _orderClient.Client.GetOrderAsync(new GetOrderRequest
            {
                OrderId = orderId
            }, BuildCallOptions(cancellationToken));
            return response.Order;
        }
        catch (RpcException ex) when (ex.StatusCode == StatusCode.NotFound)
        {
            return null;
        }
    }

    public async Task<ListOrdersResponse> ListOrdersAsync(ListOrdersRequest request, CancellationToken cancellationToken = default)
    {
        // Re-create request with same fields but use BuildCallOptions for auth headers
        return await _orderClient.Client.ListOrdersAsync(request, BuildCallOptions(cancellationToken));
    }

    public async Task<InitiatePaymentResponse> InitiatePaymentAsync(string orderId, string paymentMethod, string callbackUrl, string idempotencyKey, CancellationToken cancellationToken = default)
    {
        return await _orderClient.Client.InitiatePaymentAsync(new InitiatePaymentRequest
        {
            OrderId = orderId,
            PaymentMethod = paymentMethod,
            CallbackUrl = callbackUrl,
            IdempotencyKey = idempotencyKey
        }, BuildCallOptions(cancellationToken));
    }

    public async Task<ConfirmPaymentResponse> ConfirmPaymentAsync(string orderId, string paymentId, string transactionId, CancellationToken cancellationToken = default)
    {
        return await _orderClient.Client.ConfirmPaymentAsync(new ConfirmPaymentRequest
        {
            OrderId = orderId,
            PaymentId = paymentId,
            TransactionId = transactionId
        }, BuildCallOptions(cancellationToken));
    }

    public async Task<CancelOrderResponse> CancelOrderAsync(string orderId, string reason, CancellationToken cancellationToken = default)
    {
        return await _orderClient.Client.CancelOrderAsync(new CancelOrderRequest
        {
            OrderId = orderId,
            Reason = reason
        }, BuildCallOptions(cancellationToken));
    }

    public async Task<GetOrderStatusResponse?> GetOrderStatusAsync(string orderId, CancellationToken cancellationToken = default)
    {
        try
        {
            return await _orderClient.Client.GetOrderStatusAsync(new GetOrderStatusRequest
            {
                OrderId = orderId
            }, BuildCallOptions(cancellationToken));
        }
        catch (RpcException ex) when (ex.StatusCode == StatusCode.NotFound)
        {
            return null;
        }
    }
}
