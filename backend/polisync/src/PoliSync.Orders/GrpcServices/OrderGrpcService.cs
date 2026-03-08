using Google.Protobuf.WellKnownTypes;
using Grpc.Core;
using Insuretech.Common.V1;
using Insuretech.Orders.Entity.V1;
using Insuretech.Orders.Services.V1;
using Insuretech.Policy.Entity.V1;
using MediatR;
using Microsoft.Extensions.Logging;
using PoliSync.Orders.Infrastructure;
using PoliSync.Quotes.Infrastructure;
using QuotationEntity = Insuretech.Policy.Entity.V1.Quotation;

namespace PoliSync.Orders.GrpcServices;

public sealed class OrderGrpcService : OrderService.OrderServiceBase
{
    private readonly IMediator _mediator;
    private readonly ILogger<OrderGrpcService> _logger;
    private readonly IQuotationDataGateway _quotationDataGateway;
    private readonly IOrderDataGateway _orderDataGateway;

    public OrderGrpcService(
        IMediator mediator,
        ILogger<OrderGrpcService> logger,
        IQuotationDataGateway quotationDataGateway,
        IOrderDataGateway orderDataGateway)
    {
        _mediator = mediator;
        _logger = logger;
        _quotationDataGateway = quotationDataGateway;
        _orderDataGateway = orderDataGateway;
    }

    public override async Task<CreateOrderResponse> CreateOrder(CreateOrderRequest request, ServerCallContext context)
    {
        if (string.IsNullOrWhiteSpace(request.QuotationId) || string.IsNullOrWhiteSpace(request.CustomerId))
        {
            return new CreateOrderResponse
            {
                Error = BuildError("VALIDATION_ERROR", "QuotationId and CustomerId are required")
            };
        }

        QuotationEntity? quotation;
        try
        {
            quotation = await _quotationDataGateway.GetQuotationAsync(request.QuotationId, GetCancellationToken(context));
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "Failed to load quotation {QuotationId}", request.QuotationId);
            return new CreateOrderResponse
            {
                Error = BuildError("UPSTREAM_ERROR", ex.Status.Detail)
            };
        }

        if (quotation is null)
        {
            return new CreateOrderResponse
            {
                Error = BuildError("NOT_FOUND", "Quotation not found")
            };
        }

        if (quotation.Status != QuotationStatus.Approved)
        {
            return new CreateOrderResponse
            {
                Error = BuildError("INVALID_STATE", "Only approved quotations can be converted to orders")
            };
        }

        if (HasValue(quotation.ValidUntil) && quotation.ValidUntil.ToDateTime() < DateTime.UtcNow)
        {
            return new CreateOrderResponse
            {
                Error = BuildError("EXPIRED", "Quotation has expired")
            };
        }

        try
        {
            var response = await _orderDataGateway.CreateOrderAsync(
                request.QuotationId,
                request.CustomerId,
                request.PaymentMethod,
                GetCancellationToken(context));

            if (response.Order is not null)
            {
                EnrichOrderView(response.Order, quotation, request.CustomerId);
            }

            return response;
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "Failed to create order for quotation {QuotationId}", request.QuotationId);
            return new CreateOrderResponse
            {
                Error = BuildError(MapErrorCode(ex.StatusCode), ex.Status.Detail)
            };
        }
    }

    public override async Task<GetOrderResponse> GetOrder(GetOrderRequest request, ServerCallContext context)
    {
        try
        {
            var order = await _orderDataGateway.GetOrderAsync(request.OrderId, GetCancellationToken(context));
            if (order is null)
            {
                return new GetOrderResponse
                {
                    Error = BuildError("NOT_FOUND", "Order not found")
                };
            }

            await TryEnrichOrderViewAsync(order, GetCancellationToken(context));
            return new GetOrderResponse { Order = order };
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "Failed to get order {OrderId}", request.OrderId);
            return new GetOrderResponse
            {
                Error = BuildError(MapErrorCode(ex.StatusCode), ex.Status.Detail)
            };
        }
    }

    public override async Task<ListOrdersResponse> ListOrders(ListOrdersRequest request, ServerCallContext context)
    {
        try
        {
            var response = await _orderDataGateway.ListOrdersAsync(new ListOrdersRequest
            {
                PageSize = request.PageSize <= 0 ? 20 : request.PageSize,
                PageToken = request.PageToken,
                CustomerId = request.CustomerId,
                Status = request.Status,
                StartDate = request.StartDate,
                EndDate = request.EndDate
            }, GetCancellationToken(context));

            foreach (var order in response.Orders)
            {
                await TryEnrichOrderViewAsync(order, GetCancellationToken(context));
            }

            return response;
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "Failed to list orders for customer {CustomerId}", request.CustomerId);
            return new ListOrdersResponse
            {
                Error = BuildError(MapErrorCode(ex.StatusCode), ex.Status.Detail)
            };
        }
    }

    public override async Task<InitiatePaymentResponse> InitiatePayment(InitiatePaymentRequest request, ServerCallContext context)
    {
        if (string.IsNullOrWhiteSpace(request.PaymentMethod) || string.IsNullOrWhiteSpace(request.IdempotencyKey))
        {
            return new InitiatePaymentResponse
            {
                Error = BuildError("VALIDATION_ERROR", "PaymentMethod and IdempotencyKey are required")
            };
        }

        try
        {
            var order = await _orderDataGateway.GetOrderAsync(request.OrderId, GetCancellationToken(context));
            if (order is null)
            {
                return new InitiatePaymentResponse
                {
                    Error = BuildError("NOT_FOUND", "Order not found")
                };
            }

            if (order.Order?.Status != OrderStatus.Pending)
            {
                return new InitiatePaymentResponse
                {
                    Error = BuildError("INVALID_STATE", "Payment can only be initiated for pending orders")
                };
            }

            return await _orderDataGateway.InitiatePaymentAsync(
                request.OrderId,
                request.PaymentMethod,
                request.CallbackUrl,
                request.IdempotencyKey,
                GetCancellationToken(context));
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "Failed to initiate payment for order {OrderId}", request.OrderId);
            return new InitiatePaymentResponse
            {
                Error = BuildError(MapErrorCode(ex.StatusCode), ex.Status.Detail)
            };
        }
    }

    public override async Task<ConfirmPaymentResponse> ConfirmPayment(ConfirmPaymentRequest request, ServerCallContext context)
    {
        try
        {
            var order = await _orderDataGateway.GetOrderAsync(request.OrderId, GetCancellationToken(context));
            if (order is null)
            {
                return new ConfirmPaymentResponse
                {
                    Error = BuildError("NOT_FOUND", "Order not found")
                };
            }

            if (order.Order?.Status != OrderStatus.PaymentInitiated)
            {
                return new ConfirmPaymentResponse
                {
                    Error = BuildError("INVALID_STATE", "Only payment-initiated orders can be confirmed")
                };
            }

            if (!string.IsNullOrWhiteSpace(order.Order.PaymentId) && order.Order.PaymentId != request.PaymentId)
            {
                return new ConfirmPaymentResponse
                {
                    Error = BuildError("PAYMENT_MISMATCH", "PaymentId does not match the order")
                };
            }

            return await _orderDataGateway.ConfirmPaymentAsync(
                request.OrderId,
                request.PaymentId,
                request.TransactionId,
                GetCancellationToken(context));
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "Failed to confirm payment for order {OrderId}", request.OrderId);
            return new ConfirmPaymentResponse
            {
                Error = BuildError(MapErrorCode(ex.StatusCode), ex.Status.Detail)
            };
        }
    }

    public override async Task<CancelOrderResponse> CancelOrder(CancelOrderRequest request, ServerCallContext context)
    {
        if (string.IsNullOrWhiteSpace(request.Reason))
        {
            return new CancelOrderResponse
            {
                Error = BuildError("VALIDATION_ERROR", "Reason is required")
            };
        }

        try
        {
            var order = await _orderDataGateway.GetOrderAsync(request.OrderId, GetCancellationToken(context));
            if (order is null)
            {
                return new CancelOrderResponse
                {
                    Error = BuildError("NOT_FOUND", "Order not found")
                };
            }

            if (order.Order?.Status is OrderStatus.PolicyIssued or OrderStatus.Cancelled)
            {
                return new CancelOrderResponse
                {
                    Error = BuildError("INVALID_STATE", $"Cannot cancel order from {order.Order.Status}")
                };
            }

            return await _orderDataGateway.CancelOrderAsync(
                request.OrderId,
                request.Reason,
                GetCancellationToken(context));
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "Failed to cancel order {OrderId}", request.OrderId);
            return new CancelOrderResponse
            {
                Error = BuildError(MapErrorCode(ex.StatusCode), ex.Status.Detail)
            };
        }
    }

    public override async Task<GetOrderStatusResponse> GetOrderStatus(GetOrderStatusRequest request, ServerCallContext context)
    {
        try
        {
            var response = await _orderDataGateway.GetOrderStatusAsync(request.OrderId, GetCancellationToken(context));
            if (response is null)
            {
                return new GetOrderStatusResponse
                {
                    Error = BuildError("NOT_FOUND", "Order not found")
                };
            }

            return response;
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "Failed to get order status {OrderId}", request.OrderId);
            return new GetOrderStatusResponse
            {
                Error = BuildError(MapErrorCode(ex.StatusCode), ex.Status.Detail)
            };
        }
    }

    private async Task TryEnrichOrderViewAsync(OrderView order, CancellationToken cancellationToken)
    {
        var quotationId = order.Order?.QuotationId;
        if (string.IsNullOrWhiteSpace(quotationId))
        {
            return;
        }

        try
        {
            var quotation = await _quotationDataGateway.GetQuotationAsync(quotationId, cancellationToken);
            if (quotation is not null)
            {
                EnrichOrderView(order, quotation, order.Order?.CustomerId ?? string.Empty);
            }
        }
        catch (RpcException ex)
        {
            _logger.LogDebug(ex, "Failed to enrich order {OrderId} with quotation data", order.Order?.OrderId);
        }
    }

    private static void EnrichOrderView(OrderView orderView, QuotationEntity quotation, string customerId)
    {
        orderView.QuotationNumber = string.IsNullOrWhiteSpace(orderView.QuotationNumber)
            ? quotation.QuotationNumber
            : orderView.QuotationNumber;
        orderView.PlanName = string.IsNullOrWhiteSpace(orderView.PlanName)
            ? quotation.PlanName
            : orderView.PlanName;
        orderView.ProductName = string.IsNullOrWhiteSpace(orderView.ProductName)
            ? quotation.InsuranceCategory.ToString()
            : orderView.ProductName;
        orderView.CustomerName = string.IsNullOrWhiteSpace(orderView.CustomerName)
            ? customerId
            : orderView.CustomerName;
    }

    private static bool HasValue(Timestamp? timestamp)
        => timestamp is not null && timestamp.Seconds > 0;

    private static Error BuildError(string code, string message)
        => new() { Code = code, Message = message };

    private static CancellationToken GetCancellationToken(ServerCallContext? context)
        => context?.CancellationToken ?? CancellationToken.None;

    private static string MapErrorCode(StatusCode statusCode)
    {
        return statusCode switch
        {
            StatusCode.InvalidArgument => "VALIDATION_ERROR",
            StatusCode.NotFound => "NOT_FOUND",
            StatusCode.FailedPrecondition => "INVALID_STATE",
            StatusCode.Aborted => "PAYMENT_FAILED",
            StatusCode.Unimplemented => "NOT_IMPLEMENTED",
            _ => "UPSTREAM_ERROR"
        };
    }
}
