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
            request.PolicyId,
            request.AgentId,
            request.PremiumAmount);

        var result = await _mediator.Send(command);

        if (result.IsSuccess)
            return Ok(new { id = result.Value, message = "Commission calculated successfully" });

        return BadRequest(MapError(result.Error!));
    }

    [HttpGet("{id}")]
    public async Task<IActionResult> Get(string id)
    {
        var result = await _mediator.Send(new GetCommissionQuery(id));

        if (result.IsNotFound)
            return NotFound(MapError(result.Error!));

        if (result.IsFailure)
            return BadRequest(MapError(result.Error!));

        return Ok(MapToResponse(result.Value!));
    }

    [HttpGet]
    public async Task<IActionResult> List(
        [FromQuery] string? agentId,
        [FromQuery] string? status,
        [FromQuery] int page = 1,
        [FromQuery] int pageSize = 20)
    {
        var result = await _mediator.Send(new ListCommissionsQuery(agentId, status, page, pageSize));

        if (result.IsFailure)
            return BadRequest(MapError(result.Error!));

        var items = result.Value.Items.Select(MapToResponse).ToList();
        return Ok(new CommissionListResponse(items, result.Value.TotalCount));
    }

    [HttpPost("{id}/payout")]
    public async Task<IActionResult> ProcessPayout(string id)
    {
        var result = await _mediator.Send(new ProcessPayoutCommand(id));

        if (result.IsFailure)
            return BadRequest(MapError(result.Error!));

        return Ok(new { message = "Commission payout processed" });
    }

    private static ErrorDto MapError(SharedKernel.CQRS.Error error) => new(error.Code, error.Message);

    private static CommissionResponse MapToResponse(CommissionDto dto) => new(
        dto.CommissionId,
        dto.PolicyId,
        dto.AgentId,
        dto.PremiumAmount,
        dto.CommissionRate,
        dto.CommissionAmount,
        dto.Status,
        dto.PaidAt,
        dto.CreatedAt);
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
