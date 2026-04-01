using MediatR;
using Microsoft.AspNetCore.Mvc;
using PoliSync.Orders.Application.Commands;
using PoliSync.Orders.Application.Queries;
using PoliSync.Orders.Domain;

namespace PoliSync.ApiHost.Controllers;

// BUG-002 FIX: Changed from [Route("api/[controller]")] to explicit /v1/orders routes.
// The InScore gateway proxies /v1/orders/* to this controller on port 50141.
[ApiController]
public class OrdersController : ControllerBase
{
    private readonly IMediator _mediator;
    private readonly ILogger<OrdersController> _logger;

    public OrdersController(IMediator mediator, ILogger<OrdersController> logger)
    {
        _mediator = mediator;
        _logger = logger;
    }

    /// <summary>
    /// Create a new order from an approved quotation
    /// </summary>
    [HttpPost("/v1/orders")]
    public async Task<IActionResult> CreateOrder([FromBody] CreateOrderRequest request)
    {
        var command = new CreateOrderCommand(
            request.QuotationId,
            request.CustomerId,
            request.ProductId,
            request.PlanId,
            request.TotalPayable,
            request.Currency,
            request.PaymentDueAt,
            request.CoverageStartAt,
            request.CoverageEndAt);

        var result = await _mediator.Send(command);

        if (!result.IsSuccess)
            return BadRequest(new { error = result.Error });

        var order = await _mediator.Send(new GetOrderByIdQuery(result.Value));

        return CreatedAtAction(
            nameof(GetOrder),
            new { id = result.Value },
            order is null ? new { orderId = result.Value } : order);
    }

    /// <summary>
    /// Get order by ID
    /// </summary>
    [HttpGet("/v1/orders/{id}")]
    public async Task<IActionResult> GetOrder(Guid id)
    {
        var query = new GetOrderByIdQuery(id);
        var order = await _mediator.Send(query);

        if (order == null)
            return NotFound(new { error = "Order not found" });

        return Ok(order);
    }

    /// <summary>
    /// List orders with filtering
    /// </summary>
    [HttpGet("/v1/orders")]
    public async Task<IActionResult> ListOrders(
        [FromQuery] Guid? customerId = null,
        [FromQuery] OrderStatus? status = null,
        [FromQuery] int pageNumber = 1,
        [FromQuery] int pageSize = 20)
    {
        var query = new ListOrdersQuery(customerId, status, pageNumber, pageSize);
        var result = await _mediator.Send(query);

        if (!result.IsSuccess)
            return BadRequest(new { error = result.Error });

        return Ok(result.Value);
    }

    /// <summary>
    /// Initiate payment for an order
    /// </summary>
    [HttpPost("/v1/orders/{id}/initiate-payment")]
    public async Task<IActionResult> InitiatePayment(Guid id, [FromBody] InitiatePaymentRequest request)
    {
        var command = new InitiatePaymentCommand(
            id,
            request.PaymentMethod,
            request.CallbackUrl,
            request.IdempotencyKey);

        var result = await _mediator.Send(command);

        if (!result.IsSuccess)
            return BadRequest(new { error = result.Error });

        return Ok(result.Value);
    }

    /// <summary>
    /// Confirm payment for an order
    /// </summary>
    [HttpPost("/v1/orders/{id}/confirm-payment")]
    public async Task<IActionResult> ConfirmPayment(Guid id, [FromBody] ConfirmPaymentRequest request)
    {
        var command = new ConfirmPaymentCommand(id, request.PaymentId, request.TransactionId);
        var result = await _mediator.Send(command);

        if (!result.IsSuccess)
            return result.Error?.Code == "NOT_FOUND"
                ? NotFound(new { error = result.Error })
                : BadRequest(new { error = result.Error });

        var order = await _mediator.Send(new GetOrderByIdQuery(id));
        return order is null
            ? NotFound(new { error = "Order not found after payment confirmation" })
            : Ok(order);
    }

    /// <summary>
    /// Cancel an order
    /// </summary>
    [HttpPost("/v1/orders/{id}/cancel")]
    public async Task<IActionResult> CancelOrder(Guid id, [FromBody] CancelOrderRequest request)
    {
        var command = new CancelOrderCommand(id, request.Reason);
        var result = await _mediator.Send(command);

        if (!result.IsSuccess)
            return result.Error?.Code == "NOT_FOUND"
                ? NotFound(new { error = result.Error })
                : BadRequest(new { error = result.Error });

        var order = await _mediator.Send(new GetOrderByIdQuery(id));
        return order is null
            ? NotFound(new { error = "Order not found after cancellation" })
            : Ok(order);
    }
}

// Request DTOs
public record CreateOrderRequest(
    Guid QuotationId,
    Guid CustomerId,
    Guid ProductId,
    Guid PlanId,
    long TotalPayable,
    string Currency = "BDT",
    DateTime? PaymentDueAt = null,
    DateTime? CoverageStartAt = null,
    DateTime? CoverageEndAt = null);

public record InitiatePaymentRequest(
    string PaymentMethod,
    string? CallbackUrl = null,
    string? IdempotencyKey = null);

public record ConfirmPaymentRequest(
    string PaymentId,
    string TransactionId);

public record CancelOrderRequest(
    string Reason);
