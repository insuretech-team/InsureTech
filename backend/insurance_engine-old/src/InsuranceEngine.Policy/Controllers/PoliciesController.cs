using System;
using System.Collections.Generic;
using System.Threading.Tasks;
using MediatR;
using Microsoft.AspNetCore.Mvc;
using InsuranceEngine.Policy.Application.Features.Commands.CreatePolicy;
using InsuranceEngine.Policy.Application.Features.Commands.IssuePolicy;
using InsuranceEngine.Policy.Application.Features.Commands.CancelPolicy;
using InsuranceEngine.Policy.Application.Features.Commands.RenewPolicy;
using InsuranceEngine.Policy.Application.Features.Commands.Nominees;
using InsuranceEngine.Policy.Application.Features.Queries;
using InsuranceEngine.Policy.Application.DTOs;
using InsuranceEngine.Policy.Domain.Enums;
using InsuranceEngine.SharedKernel.DTOs;

namespace InsuranceEngine.Policy.Controllers;

[ApiController]
[Route("v1/policies")]
public class PoliciesController : ControllerBase
{
    private readonly IMediator _mediator;

    public PoliciesController(IMediator mediator) => _mediator = mediator;

    // ===================== Policy CRUD =====================

    [HttpGet]
    public async Task<IActionResult> List(
        [FromQuery] Guid? customerId, [FromQuery] PolicyStatus? status,
        [FromQuery] Guid? productId, [FromQuery] int page = 1, [FromQuery] int pageSize = 20)
    {
        var result = await _mediator.Send(new ListPoliciesQuery(customerId, status, productId, page, pageSize));
        return Ok(new PolicyListingResponse(result.Items, result.TotalCount));
    }

    [HttpGet("{id}")]
    public async Task<IActionResult> Get(Guid id)
    {
        var result = await _mediator.Send(new GetPolicyQuery(id));
        if (result == null) return NotFound(new ErrorDto("POLICY_NOT_FOUND", "Policy not found."));
        return Ok(new PolicyRetrievalResponse(result));
    }

    [HttpPost]
    public async Task<IActionResult> Create([FromBody] CreatePolicyCommand command)
    {
        var result = await _mediator.Send(command);
        if (result.IsSuccess)
        {
            var r = result.Value!;
            var response = new PolicyCreationResponse(r.PolicyId, r.PolicyNumber, "Policy created successfully.");
            return CreatedAtAction(nameof(Get), new { id = r.PolicyId }, response);
        }
        return BadRequest(MapError(result.Error!));
    }

    // ===================== Lifecycle =====================

    [HttpPost("{id}/issue")]
    public async Task<IActionResult> Issue(Guid id)
    {
        var result = await _mediator.Send(new IssuePolicyCommand(id));
        if (result.IsSuccess) return Ok(new PolicyIssueResponse("Policy issued successfully."));
        return HandleErrorResult(result.Error!);
    }

    [HttpPost("{id}/cancel")]
    public async Task<IActionResult> Cancel(Guid id, [FromBody] CancelPolicyRequest request)
    {
        var result = await _mediator.Send(new CancelPolicyCommand(id, request.Reason));
        if (result.IsSuccess) return Ok(new PolicyCancelResponse("Policy cancelled successfully."));
        return HandleErrorResult(result.Error!);
    }

    [HttpPost("{id}/renew")]
    public async Task<IActionResult> Renew(Guid id, [FromBody] RenewPolicyRequest request)
    {
        var result = await _mediator.Send(new RenewPolicyCommand(id, request.TenureMonths));
        if (result.IsSuccess)
        {
            var r = result.Value!;
            return Ok(new PolicyRenewResponse(r.NewPolicyId, r.NewPolicyNumber, "Policy renewed successfully."));
        }
        return HandleErrorResult(result.Error!);
    }

    //// ===================== Grace Period & Renewal =====================

    [HttpGet("{id}/grace-period")]
    public async Task<IActionResult> GetGracePeriod(Guid id)
    {
        var result = await _mediator.Send(new GetGracePeriodQuery(id));
        if (result == null) return NotFound(new ErrorDto("POLICY_NOT_FOUND", "Policy not found."));
        return Ok(new GracePeriodResponse(result));
    }

    [HttpGet("{id}/renewal-schedule")]
    public async Task<IActionResult> GetRenewalSchedule(Guid id)
    {
        var result = await _mediator.Send(new GetRenewalScheduleQuery(id));
        if (result == null) return NotFound(new ErrorDto("POLICY_NOT_FOUND", "Policy not found."));
        return Ok(new RenewalScheduleResponse(result));
    }

    //// ===================== Nominees =====================

    [HttpGet("{policyId}/nominees")]
    public async Task<IActionResult> ListNominees(Guid policyId)
    {
        var items = await _mediator.Send(new ListNomineesQuery(policyId));
        return Ok(new NomineeListingResponse(items, items.Count));
    }

    [HttpPost("{policyId}/nominees")]
    public async Task<IActionResult> AddNominee(Guid policyId, [FromBody] AddNomineeRequest request)
    {
        var command = new AddNomineeCommand(policyId, request.BeneficiaryId, request.FullName, request.Relationship, request.SharePercentage,
            request.DateOfBirth, request.NidNumber, request.PhoneNumber, request.NomineeDobText);
        var result = await _mediator.Send(command);
        if (result.IsSuccess)
        {
            return Created($"api/policies/{policyId}/nominees/{result.Value}",
                new NomineeResponse(result.Value, "Nominee added successfully."));
        }
        return HandleErrorResult(result.Error!);
    }

    [HttpPut("{policyId}/nominees/{nomineeId}")]
    public async Task<IActionResult> UpdateNominee(Guid policyId, Guid nomineeId, [FromBody] UpdateNomineeRequest request)
    {
        var command = new UpdateNomineeCommand(policyId, nomineeId, request.FullName, request.Relationship, request.SharePercentage,
            request.DateOfBirth, request.NidNumber, request.PhoneNumber, request.NomineeDobText);
        var result = await _mediator.Send(command);
        if (result.IsSuccess) return Ok(new NomineeResponse(nomineeId, "Nominee updated successfully."));
        return HandleErrorResult(result.Error!);
    }

    [HttpDelete("{policyId}/nominees/{nomineeId}")]
    public async Task<IActionResult> DeleteNominee(Guid policyId, Guid nomineeId)
    {
        var result = await _mediator.Send(new DeleteNomineeCommand(policyId, nomineeId));
        if (result.IsSuccess) return Ok(new { message = "Nominee deleted successfully." });
        return HandleErrorResult(result.Error!);
    }

    // ===================== Helpers =====================

    private ErrorDto MapError(SharedKernel.CQRS.Error error)
    {
        return new ErrorDto(error.Code, error.Message);
    }

    private IActionResult HandleErrorResult(SharedKernel.CQRS.Error error)
    {
        var errorDto = MapError(error);
        return error.Code switch
        {
            "NOT_FOUND" => NotFound(errorDto),
            "INVALID_STATE_TRANSITION" => Conflict(errorDto),
            "VALIDATION_ERROR" => BadRequest(errorDto),
            _ => BadRequest(errorDto)
        };
    }
}

// --- Request DTOs ---
public record CancelPolicyRequest(string Reason);
public record RenewPolicyRequest(int TenureMonths);
public record AddNomineeRequest(
    Guid? BeneficiaryId, string FullName, string Relationship, double SharePercentage,
    DateTime? DateOfBirth = null, string? NidNumber = null, string? PhoneNumber = null, string? NomineeDobText = null);
public record UpdateNomineeRequest(
    string? FullName = null, string? Relationship = null, double? SharePercentage = null,
    DateTime? DateOfBirth = null, string? NidNumber = null, string? PhoneNumber = null, string? NomineeDobText = null);
