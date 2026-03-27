using Microsoft.AspNetCore.Mvc;
using MediatR;
using InsuranceEngine.Policy.Application.Commands;
using InsuranceEngine.Policy.Application.Queries;
using InsuranceEngine.SharedKernel.DTOs;

namespace InsuranceEngine.ApiHost.Controllers;

// 抽离出了 Claims 模块的具体文件路径
/*
[ApiController]
[Route("v1/policies")]
public class PoliciesController : ControllerBase
{
    private readonly IMediator _mediator;
    private readonly ILogger<PoliciesController> _logger;

    public PoliciesController(IMediator mediator, ILogger<PoliciesController> logger)
    {
        _mediator = mediator;
        _logger = logger;
    }

    [HttpPost]
    public async Task<IActionResult> Create([FromBody] CreatePolicyRequest request)
    {
        var command = new CreatePolicyCommand(
            request.ProductId,
            request.CustomerId,
            request.PartnerId,
            request.AgentId,
            request.QuoteId,
            request.PremiumAmount,
            request.SumInsured,
            request.TenureMonths,
            request.StartDate,
            request.ProposerDetails,
            request.Nominees?.Select(n => new InsuranceEngine.Policy.Application.Commands.NomineeDto(
                n.BeneficiaryId, n.FullName, n.Relationship, n.SharePercentage,
                n.DateOfBirth, n.NidNumber, n.PhoneNumber, n.NomineeDobText)).ToList());

        var result = await _mediator.Send(command);

        if (result.IsSuccess)
            return CreatedAtAction(nameof(Get), new { id = result.Value }, new { id = result.Value, message = "Policy created successfully" });

        return BadRequest(MapError(result.Error!));
    }

    [HttpGet("{id}")]
    public async Task<IActionResult> Get(string id)
    {
        var result = await _mediator.Send(new GetPolicyQuery(id));

        if (result.IsNotFound)
            return NotFound(MapError(result.Error!));

        if (result.IsFailure)
            return BadRequest(MapError(result.Error!));

        return Ok(MapToResponse(result.Value!));
    }

    [HttpGet]
    public async Task<IActionResult> List(
        [FromQuery] string? customerId,
        [FromQuery] string? status,
        [FromQuery] string? productId,
        [FromQuery] int page = 1,
        [FromQuery] int pageSize = 20)
    {
        var result = await _mediator.Send(new ListPoliciesQuery(customerId, status, productId, page, pageSize));

        if (result.IsFailure)
            return BadRequest(MapError(result.Error!));

        var items = result.Value.Items.Select(MapToResponse).ToList();
        return Ok(new PolicyListResponse(items, result.Value.TotalCount));
    }

    [HttpPost("{id}/issue")]
    public async Task<IActionResult> Issue(string id)
    {
        var result = await _mediator.Send(new IssuePolicyCommand(id));

        if (result.IsNotFound)
            return NotFound(MapError(result.Error!));

        if (result.IsFailure)
            return BadRequest(MapError(result.Error!));

        return Ok(new { message = "Policy issued successfully" });
    }

    [HttpPost("{id}/cancel")]
    public async Task<IActionResult> Cancel(string id, [FromBody] CancelPolicyRequest request)
    {
        var result = await _mediator.Send(new CancelPolicyCommand(id, request.Reason));

        if (result.IsNotFound)
            return NotFound(MapError(result.Error!));

        if (result.IsFailure)
            return BadRequest(MapError(result.Error!));

        return Ok(new { message = "Policy cancelled successfully" });
    }

    [HttpPost("{id}/renew")]
    public async Task<IActionResult> Renew(string id, [FromBody] RenewPolicyRequest request)
    {
        var result = await _mediator.Send(new RenewPolicyCommand(id, request.TenureMonths));

        if (result.IsNotFound)
            return NotFound(MapError(result.Error!));

        if (result.IsFailure)
            return BadRequest(MapError(result.Error!));

        return Ok(new { id = result.Value.NewPolicyId, policyNumber = result.Value.NewPolicyNumber, message = "Policy renewed successfully" });
    }

    [HttpGet("{policyId}/nominees")]
    public async Task<IActionResult> ListNominees(string policyId)
    {
        var result = await _mediator.Send(new ListNomineesQuery(policyId));

        if (result.IsFailure)
            return BadRequest(MapError(result.Error!));

        var items = result.Value.Select(n => new NomineeResponse(
            n.NomineeId, n.FullName, n.Relationship, n.SharePercentage, n.DateOfBirth, n.PhoneNumber)).ToList();

        return Ok(new NomineeListResponse(items, items.Count));
    }

    [HttpPost("{policyId}/nominees")]
    public async Task<IActionResult> AddNominee(string policyId, [FromBody] AddNomineeRequest request)
    {
        var command = new AddNomineeCommand(
            policyId,
            request.FullName,
            request.Relationship,
            request.SharePercentage,
            request.DateOfBirth,
            request.NidNumber,
            request.PhoneNumber,
            request.NomineeDobText);

        var result = await _mediator.Send(command);

        if (result.IsFailure)
            return BadRequest(MapError(result.Error!));

        return Created($"v1/policies/{policyId}/nominees/{result.Value}", new { id = result.Value, message = "Nominee added successfully" });
    }

    [HttpPut("{policyId}/nominees/{nomineeId}")]
    public async Task<IActionResult> UpdateNominee(string policyId, string nomineeId, [FromBody] UpdateNomineeRequest request)
    {
        var command = new UpdateNomineeCommand(
            policyId,
            nomineeId,
            request.FullName,
            request.Relationship,
            request.SharePercentage,
            request.DateOfBirth,
            request.NidNumber,
            request.PhoneNumber);

        var result = await _mediator.Send(command);

        if (result.IsNotFound)
            return NotFound(MapError(result.Error!));

        if (result.IsFailure)
            return BadRequest(MapError(result.Error!));

        return Ok(new { message = "Nominee updated successfully" });
    }

    [HttpDelete("{policyId}/nominees/{nomineeId}")]
    public async Task<IActionResult> DeleteNominee(string policyId, string nomineeId)
    {
        var result = await _mediator.Send(new DeleteNomineeCommand(policyId, nomineeId));

        if (result.IsNotFound)
            return NotFound(MapError(result.Error!));

        if (result.IsFailure)
            return BadRequest(MapError(result.Error!));

        return Ok(new { message = "Nominee deleted successfully" });
    }

    private static ErrorDto MapError(SharedKernel.CQRS.Error error) => new(error.Code, error.Message);

    private static PolicyResponse MapToResponse(PolicyDto dto) => new(
        dto.PolicyId,
        dto.PolicyNumber,
        dto.ProductId,
        dto.CustomerId,
        dto.PartnerId,
        dto.AgentId,
        dto.Status,
        dto.PremiumAmount,
        dto.SumInsured,
        dto.TenureMonths,
        dto.StartDate,
        dto.EndDate,
        dto.IssuedAt,
        dto.CreatedAt);
}

public record CreatePolicyRequest(
    string ProductId,
    string CustomerId,
    string? PartnerId,
    string? AgentId,
    string? QuoteId,
    decimal PremiumAmount,
    decimal SumInsured,
    int TenureMonths,
    DateTime StartDate,
    string? ProposerDetails,
    List<RequestNomineeDto>? Nominees);

public record RequestNomineeDto(
    string? BeneficiaryId,
    string FullName,
    string Relationship,
    double SharePercentage,
    DateTime? DateOfBirth,
    string? NidNumber,
    string? PhoneNumber,
    string? NomineeDobText);

public record CancelPolicyRequest(string Reason);
public record RenewPolicyRequest(int TenureMonths);

public record AddNomineeRequest(
    string FullName,
    string Relationship,
    double SharePercentage,
    DateTime? DateOfBirth = null,
    string? NidNumber = null,
    string? PhoneNumber = null,
    string? NomineeDobText = null);

public record UpdateNomineeRequest(
    string? FullName = null,
    string? Relationship = null,
    double? SharePercentage = null,
    DateTime? DateOfBirth = null,
    string? NidNumber = null,
    string? PhoneNumber = null);

public record PolicyResponse(
    string PolicyId,
    string PolicyNumber,
    string ProductId,
    string CustomerId,
    string? PartnerId,
    string? AgentId,
    string Status,
    decimal PremiumAmount,
    decimal SumInsured,
    int TenureMonths,
    DateTime StartDate,
    DateTime EndDate,
    DateTime? IssuedAt,
    DateTime? CreatedAt);

public record PolicyListResponse(IReadOnlyList<PolicyResponse> Items, int TotalCount);
*/

public record NomineeResponse(
    string NomineeId,
    string FullName,
    string Relationship,
    double SharePercentage,
    DateTime? DateOfBirth,
    string? PhoneNumber);

public record NomineeListResponse(IReadOnlyList<NomineeResponse> Items, int TotalCount);
