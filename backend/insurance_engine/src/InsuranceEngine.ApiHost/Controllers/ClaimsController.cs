using Microsoft.AspNetCore.Mvc;
using MediatR;
using InsuranceEngine.Claims.Application.Commands;
using InsuranceEngine.Claims.Application.Queries;
using InsuranceEngine.SharedKernel.DTOs;

namespace InsuranceEngine.ApiHost.Controllers;

/*
[ApiController]
[Route("v1/claims")]
public class ClaimsController : ControllerBase
{
    private readonly IMediator _mediator;
    private readonly ILogger<ClaimsController> _logger;

    public ClaimsController(IMediator mediator, ILogger<ClaimsController> logger)
    {
        _mediator = mediator;
        _logger = logger;
    }

    [HttpPost]
    public async Task<IActionResult> Submit([FromBody] SubmitClaimRequest request)
    {
        var command = new SubmitClaimCommand(
            request.PolicyId,
            request.ClaimType,
            request.ClaimAmount,
            request.Description,
            request.BeneficiaryId);

        var result = await _mediator.Send(command);

        if (result.IsSuccess)
            return CreatedAtAction(nameof(Get), new { id = result.Value }, new { id = result.Value, message = "Claim submitted successfully" });

        return BadRequest(MapError(result.Error!));
    }

    [HttpGet("{id}")]
    public async Task<IActionResult> Get(string id)
    {
        var result = await _mediator.Send(new GetClaimQuery(id));

        if (result.IsNotFound)
            return NotFound(MapError(result.Error!));

        if (result.IsFailure)
            return BadRequest(MapError(result.Error!));

        return Ok(MapToResponse(result.Value!));
    }

    [HttpGet]
    public async Task<IActionResult> List(
        [FromQuery] string? policyId,
        [FromQuery] string? status,
        [FromQuery] int page = 1,
        [FromQuery] int pageSize = 20)
    {
        var result = await _mediator.Send(new ListClaimsQuery(policyId, status, page, pageSize));

        if (result.IsFailure)
            return BadRequest(MapError(result.Error!));

        var items = result.Value.Items.Select(MapToResponse).ToList();
        return Ok(new ClaimListResponse(items, result.Value.TotalCount));
    }

    [HttpPost("{id}/approve")]
    public async Task<IActionResult> Approve(string id, [FromBody] ApproveClaimRequest request)
    {
        var result = await _mediator.Send(new ApproveClaimCommand(id, request.ApprovedAmount));

        if (result.IsFailure)
            return BadRequest(MapError(result.Error!));

        return Ok(new { message = "Claim approved" });
    }

    [HttpPost("{id}/reject")]
    public async Task<IActionResult> Reject(string id, [FromBody] RejectClaimRequest request)
    {
        var result = await _mediator.Send(new RejectClaimCommand(id, request.Reason));

        if (result.IsFailure)
            return BadRequest(MapError(result.Error!));

        return Ok(new { message = "Claim rejected" });
    }

    [HttpPost("{id}/settle")]
    public async Task<IActionResult> Settle(string id, [FromBody] SettleClaimRequest request)
    {
        var result = await _mediator.Send(new SettleClaimCommand(id, request.SettlementAmount));

        if (result.IsFailure)
            return BadRequest(MapError(result.Error!));

        return Ok(new { message = "Claim settled" });
    }

    private static ErrorDto MapError(SharedKernel.CQRS.Error error) => new(error.Code, error.Message);

    private static ClaimResponse MapToResponse(ClaimDto dto) => new(
        dto.ClaimId,
        dto.ClaimNumber,
        dto.PolicyId,
        dto.ClaimType,
        dto.ClaimAmount,
        dto.ApprovedAmount,
        dto.SettlementAmount,
        dto.Description,
        dto.Status,
        dto.RejectionReason,
        dto.SettledAt,
        dto.CreatedAt);
}

public record SubmitClaimRequest(
    string PolicyId,
    string ClaimType,
    decimal ClaimAmount,
    string Description,
    string? BeneficiaryId);

public record ApproveClaimRequest(decimal ApprovedAmount);
public record RejectClaimRequest(string Reason);
public record SettleClaimRequest(decimal SettlementAmount);

public record ClaimResponse(
    string ClaimId,
    string ClaimNumber,
    string PolicyId,
    string ClaimType,
    decimal ClaimAmount,
    decimal? ApprovedAmount,
    decimal? SettlementAmount,
    string? Description,
    string Status,
    string? RejectionReason,
    DateTime? SettledAt,
    DateTime? CreatedAt);

public record ClaimListResponse(IReadOnlyList<ClaimResponse> Items, int TotalCount);
*/
