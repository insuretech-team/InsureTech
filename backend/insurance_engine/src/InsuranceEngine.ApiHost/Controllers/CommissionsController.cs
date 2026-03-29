using Microsoft.AspNetCore.Mvc;
using MediatR;
using InsuranceEngine.Commission.Application.Commands;
using InsuranceEngine.Commission.Application.Queries;
using InsuranceEngine.SharedKernel.DTOs;

namespace InsuranceEngine.ApiHost.Controllers;

[ApiController]
[Route("v1/commissions")]
public class CommissionsController : ControllerBase
{
    private readonly IMediator _mediator;
    private readonly ILogger<CommissionsController> _logger;

    public CommissionsController(IMediator mediator, ILogger<CommissionsController> logger)
    {
        _mediator = mediator;
        _logger = logger;
    }

    [HttpPost("calculate")]
    public async Task<IActionResult> Calculate([FromBody] CalculateCommissionRequest request)
    {
        var command = new CalculateCommissionCommand(
            PolicyId: request.PolicyId,
            CommissionType: "ACQUISITION", // Default or from request
            RecipientType: "AGENT",        // Default or from request
            RecipientId: request.AgentId);

        var result = await _mediator.Send(command);

        if (string.IsNullOrEmpty(result.Error?.Code))
            return Ok(new { id = result.CommissionId, message = "Commission calculated successfully" });

        return BadRequest(new { code = result.Error.Code, message = result.Error.Message });
    }

    [HttpGet("{id}")]
    public async Task<IActionResult> Get(string id)
    {
        var response = await _mediator.Send(new GetCommissionQuery(id));

        if (!string.IsNullOrEmpty(response.Error?.Code))
        {
            if (response.Error.Code == "NOT_FOUND") return NotFound(new { code = response.Error.Code, message = response.Error.Message });
            return BadRequest(new { code = response.Error.Code, message = response.Error.Message });
        }

        return Ok(MapToResponse(response.Commission));
    }

    [HttpGet]
    public async Task<IActionResult> List(
        [FromQuery] string? agentId,
        [FromQuery] string? status,
        [FromQuery] int page = 1,
        [FromQuery] int pageSize = 20)
    {
        var query = new ListCommissionsQuery(
            RecipientType: "AGENT",
            RecipientId: agentId ?? "",
            Status: status,
            StartDate: null,
            EndDate: null,
            Page: page,
            PageSize: pageSize);

        var response = await _mediator.Send(query);

        var items = response.Commissions.Select(MapToResponse).ToList();
        return Ok(new CommissionListResponse(items, response.TotalCount));
    }

    [HttpPost("{id}/payout")]
    public async Task<IActionResult> ProcessPayout(string id, [FromQuery] string paymentMethod = "BANK_TRANSFER")
    {
        var command = new ProcessPayoutCommand(
            PayoutId: id,
            PaymentMethod: paymentMethod,
            PaymentReference: null);

        var result = await _mediator.Send(command);
        return Ok(new { message = "Commission payout processed" });
    }

    private static ErrorDto MapError(SharedKernel.CQRS.Error error) => new(error.Code, error.Message);

    private static CommissionResponse MapToResponse(Insuretech.Partner.Entity.V1.Commission c) => new(
        c.CommissionId,
        c.PolicyId,
        c.AgentId ?? "",
        (decimal)c.CommissionAmount.Amount / 100m,
        (decimal)c.CommissionRate,
        (decimal)c.CommissionAmount.Amount / 100m, // Simplification for API
        c.Status.ToString(),
        c.PaidAt?.ToDateTime(),
        c.CreatedAt?.ToDateTime());
}

public record CalculateCommissionRequest(
    string PolicyId,
    string AgentId,
    decimal PremiumAmount);

public record CommissionResponse(
    string CommissionId,
    string PolicyId,
    string AgentId,
    decimal PremiumAmount,
    decimal CommissionRate,
    decimal CommissionAmount,
    string Status,
    DateTime? PaidAt,
    DateTime? CreatedAt);

public record CommissionListResponse(IReadOnlyList<CommissionResponse> Items, int TotalCount);
