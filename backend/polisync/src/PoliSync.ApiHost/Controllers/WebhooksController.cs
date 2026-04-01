using MediatR;
using Microsoft.AspNetCore.Mvc;
using PoliSync.Orders.Application.Commands;
using System.Text.Json;

namespace PoliSync.ApiHost.Controllers;

/// <summary>
/// Controller for handling payment gateway webhooks
/// </summary>
[ApiController]
[Route("api/[controller]")]
public class WebhooksController : ControllerBase
{
    private readonly IMediator _mediator;
    private readonly ILogger<WebhooksController> _logger;

    public WebhooksController(IMediator mediator, ILogger<WebhooksController> logger)
    {
        _mediator = mediator;
        _logger = logger;
    }

    /// <summary>
    /// Handle SSLCommerz payment webhook callback
    /// </summary>
    [HttpPost("sslcommerz")]
    public async Task<IActionResult> HandleSSLCommerzWebhook([FromForm] SSLCommerzWebhookPayload payload)
    {
        _logger.LogInformation("Received SSLCommerz webhook for transaction {TransactionId}", 
            payload.TranId);

        try
        {
            // Extract order ID from tran_id (format: ORDER-{orderId})
            if (!payload.TranId.StartsWith("ORDER-"))
            {
                _logger.LogWarning("Invalid transaction ID format: {TranId}", payload.TranId);
                return BadRequest(new { error = "Invalid transaction ID format" });
            }

            var orderIdStr = payload.TranId.Substring(6); // Remove "ORDER-" prefix
            if (!Guid.TryParse(orderIdStr, out var orderId))
            {
                _logger.LogWarning("Invalid order ID in transaction: {TranId}", payload.TranId);
                return BadRequest(new { error = "Invalid order ID" });
            }

            // Parse amount
            if (!decimal.TryParse(payload.Amount, out var amount))
            {
                _logger.LogWarning("Invalid amount in webhook: {Amount}", payload.Amount);
                return BadRequest(new { error = "Invalid amount" });
            }

            // Verify payment
            var command = new VerifyPaymentCommand(
                orderId,
                payload.TranId,
                payload.BankTranId ?? payload.TranId,
                payload.Status,
                amount,
                payload.Currency);

            var result = await _mediator.Send(command);

            if (!result.IsSuccess)
            {
                _logger.LogError("Payment verification failed: {Error}", result.Error?.Message);
                return BadRequest(new { error = result.Error?.Message });
            }

            return Ok(new { message = "Webhook processed successfully", verified = result.Value });
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error processing SSLCommerz webhook");
            return StatusCode(500, new { error = "Internal server error" });
        }
    }
}

/// <summary>
/// SSLCommerz webhook payload structure
/// </summary>
public class SSLCommerzWebhookPayload
{
    public string Status { get; set; } = string.Empty;
    public string TranId { get; set; } = string.Empty;
    public string? BankTranId { get; set; }
    public string Amount { get; set; } = string.Empty;
    public string Currency { get; set; } = "BDT";
    public string? CardType { get; set; }
    public string? CardNo { get; set; }
    public string? CardIssuer { get; set; }
    public string? CardBrand { get; set; }
    public string? CardIssuerCountry { get; set; }
    public string? StoreAmount { get; set; }
    public string? VerifySign { get; set; }
    public string? VerifyKey { get; set; }
    public string? RiskLevel { get; set; }
    public string? RiskTitle { get; set; }
}
