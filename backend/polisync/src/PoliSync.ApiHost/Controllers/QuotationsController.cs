using MediatR;
using Microsoft.AspNetCore.Mvc;
using PoliSync.Quotes.Application.Commands;
using PoliSync.Quotes.Application.Queries;
using PoliSync.Quotes.Domain;

namespace PoliSync.ApiHost.Controllers;

// BUG-003 FIX: Changed from [Route("api/[controller]")] to explicit /v1/quotes routes.
// The InScore gateway proxies /v1/quotes/* to this controller on port 50131.
[ApiController]
public class QuotationsController : ControllerBase
{
    private readonly IMediator _mediator;
    private readonly ILogger<QuotationsController> _logger;

    public QuotationsController(IMediator mediator, ILogger<QuotationsController> logger)
    {
        _mediator = mediator;
        _logger = logger;
    }

    /// <summary>
    /// Create a new quotation
    /// </summary>
    [HttpPost("/v1/quotations")]
    public async Task<IActionResult> CreateQuotation([FromBody] CreateQuotationRequest request)
    {
        var tenantId = request.TenantId;
        if (tenantId == Guid.Empty)
        {
            var headerTenant = HttpContext.Request.Headers["X-Tenant-ID"].FirstOrDefault();
            if (!string.IsNullOrEmpty(headerTenant) && Guid.TryParse(headerTenant, out var parsedTenant))
            {
                tenantId = parsedTenant;
            }
        }

        var command = new CreateQuotationCommand(
            tenantId,
            request.ProductId,
            request.PlanId,
            request.CustomerId,
            request.BasePremium,
            request.RiderPremium,
            request.ExpiryDays);

        var result = await _mediator.Send(command);

        if (!result.IsSuccess)
            return BadRequest(new { error = result.Error });

        var quotation = await _mediator.Send(new GetQuotationByIdQuery(result.Value));

        return CreatedAtAction(
            nameof(GetQuotation),
            new { id = result.Value },
            quotation is null ? new { quotationId = result.Value } : quotation);
    }

    /// <summary>
    /// Get quotation by ID
    /// </summary>
    [HttpGet("/v1/quotations/{id}")]
    public async Task<IActionResult> GetQuotation(Guid id)
    {
        var query = new GetQuotationByIdQuery(id);
        var quotation = await _mediator.Send(query);

        if (quotation == null)
            return NotFound(new { error = "Quotation not found" });

        return Ok(quotation);
    }

    /// <summary>
    /// List quotations with filtering
    /// </summary>
    [HttpGet("/v1/quotations")]
    public async Task<IActionResult> ListQuotations(
        [FromQuery] Guid? customerId = null,
        [FromQuery] QuotationStatus? status = null,
        [FromQuery] int pageNumber = 1,
        [FromQuery] int pageSize = 20)
    {
        var query = new ListQuotationsQuery(customerId, status, pageNumber, pageSize);
        var result = await _mediator.Send(query);

        if (!result.IsSuccess)
            return BadRequest(new { error = result.Error });

        return Ok(result.Value);
    }

    /// <summary>
    /// Submit quotation for underwriting review
    /// </summary>
    [HttpPost("/v1/quotations/{id}/submit")]
    public async Task<IActionResult> SubmitQuotation(Guid id)
    {
        var command = new SubmitQuotationCommand(id);
        var result = await _mediator.Send(command);

        if (!result.IsSuccess)
            return result.Error?.Code == "NOT_FOUND"
                ? NotFound(new { error = result.Error })
                : BadRequest(new { error = result.Error });

        var quotation = await _mediator.Send(new GetQuotationByIdQuery(id));
        return quotation is null
            ? NotFound(new { error = "Quotation not found after submit" })
            : Ok(quotation);
    }

    /// <summary>
    /// Mark quotation as received by underwriting
    /// </summary>
    [HttpPost("/v1/quotations/{id}/mark-received")]
    public async Task<IActionResult> MarkQuotationReceived(Guid id)
    {
        var command = new MarkQuotationReceivedCommand(id);
        var result = await _mediator.Send(command);

        if (!result.IsSuccess)
            return result.Error?.Code == "NOT_FOUND"
                ? NotFound(new { error = result.Error })
                : BadRequest(new { error = result.Error });

        var quotation = await _mediator.Send(new GetQuotationByIdQuery(id));
        return quotation is null
            ? NotFound(new { error = "Quotation not found after marking as received" })
            : Ok(quotation);
    }

    /// <summary>
    /// Apply loading to quotation premium
    /// </summary>
    [HttpPost("/v1/quotations/{id}/apply-loading")]
    public async Task<IActionResult> ApplyLoading(Guid id, [FromBody] ApplyLoadingRequest request)
    {
        var command = new ApplyLoadingCommand(id, request.LoadingAmount);
        var result = await _mediator.Send(command);

        if (!result.IsSuccess)
            return result.Error?.Code == "NOT_FOUND"
                ? NotFound(new { error = result.Error })
                : BadRequest(new { error = result.Error });

        var quotation = await _mediator.Send(new GetQuotationByIdQuery(id));
        return quotation is null
            ? NotFound(new { error = "Quotation not found after applying loading" })
            : Ok(quotation);
    }

    /// <summary>
    /// Apply discount to quotation premium
    /// </summary>
    [HttpPost("/v1/quotations/{id}/apply-discount")]
    public async Task<IActionResult> ApplyDiscount(Guid id, [FromBody] ApplyDiscountRequest request)
    {
        var command = new ApplyDiscountCommand(id, request.DiscountAmount);
        var result = await _mediator.Send(command);

        if (!result.IsSuccess)
            return result.Error?.Code == "NOT_FOUND"
                ? NotFound(new { error = result.Error })
                : BadRequest(new { error = result.Error });

        var quotation = await _mediator.Send(new GetQuotationByIdQuery(id));
        return quotation is null
            ? NotFound(new { error = "Quotation not found after applying discount" })
            : Ok(quotation);
    }

    /// <summary>
    /// Approve quotation
    /// </summary>
    [HttpPost("/v1/quotations/{id}/approve")]
    public async Task<IActionResult> ApproveQuotation(Guid id)
    {
        var command = new ApproveQuotationCommand(id);
        var result = await _mediator.Send(command);

        if (!result.IsSuccess)
            return result.Error?.Code == "NOT_FOUND"
                ? NotFound(new { error = result.Error })
                : BadRequest(new { error = result.Error });

        var quotation = await _mediator.Send(new GetQuotationByIdQuery(id));
        return quotation is null
            ? NotFound(new { error = "Quotation not found after approval" })
            : Ok(quotation);
    }

    /// <summary>
    /// Reject quotation
    /// </summary>
    [HttpPost("/v1/quotations/{id}/reject")]
    public async Task<IActionResult> RejectQuotation(Guid id, [FromBody] RejectQuotationRequest request)
    {
        var command = new RejectQuotationCommand(id, request.Reason);
        var result = await _mediator.Send(command);

        if (!result.IsSuccess)
            return result.Error?.Code == "NOT_FOUND"
                ? NotFound(new { error = result.Error })
                : BadRequest(new { error = result.Error });

        var quotation = await _mediator.Send(new GetQuotationByIdQuery(id));
        return quotation is null
            ? NotFound(new { error = "Quotation not found after rejection" })
            : Ok(quotation);
    }
}

// Request DTOs
public record CreateQuotationRequest(
    Guid TenantId,
    Guid ProductId,
    Guid PlanId,
    Guid CustomerId,
    long BasePremium,
    long RiderPremium,
    int ExpiryDays = 30);

public record ApplyLoadingRequest(
    long LoadingAmount);

public record ApplyDiscountRequest(
    long DiscountAmount);

public record RejectQuotationRequest(
    string Reason);

