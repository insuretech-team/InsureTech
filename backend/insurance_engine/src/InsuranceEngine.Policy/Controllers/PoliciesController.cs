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

namespace InsuranceEngine.Policy.Controllers;

[ApiController]
[Route("api/[controller]")]
public class PoliciesController : ControllerBase
{
    private readonly IMediator _mediator;

    public PoliciesController(IMediator mediator) => _mediator = mediator;

    // ===================== Policy CRUD =====================

    [HttpGet]
    public async Task<ActionResult<PaginatedResponse<PolicyListDto>>> List(
        [FromQuery] Guid? customerId, [FromQuery] PolicyStatus? status,
        [FromQuery] Guid? productId, [FromQuery] int page = 1, [FromQuery] int pageSize = 20)
    {
        return Ok(await _mediator.Send(new ListPoliciesQuery(customerId, status, productId, page, pageSize)));
    }

    [HttpGet("{id}")]
    public async Task<ActionResult<PolicyDto>> Get(Guid id)
    {
        var result = await _mediator.Send(new GetPolicyQuery(id));
        return result != null ? Ok(result) : NotFound();
    }

    [HttpPost]
    public async Task<IActionResult> Create([FromBody] CreatePolicyCommand command)
    {
        var result = await _mediator.Send(command);
        return result.Match<IActionResult>(
            onSuccess: r => CreatedAtAction(nameof(Get), new { id = r.PolicyId },
                new { r.PolicyId, r.PolicyNumber, message = "Policy created successfully." }),
            onFailure: error => error.Code switch
            {
                "NOT_FOUND" => NotFound(new { error.Code, error.Message }),
                "VALIDATION_ERROR" => BadRequest(new { error.Code, error.Message }),
                _ => BadRequest(new { error.Code, error.Message })
            }
        );
    }

    // ===================== Lifecycle =====================

    [HttpPost("{id}/issue")]
    public async Task<IActionResult> Issue(Guid id)
    {
        var result = await _mediator.Send(new IssuePolicyCommand(id));
        return HandleResult(result, "Policy issued successfully.");
    }

    [HttpPost("{id}/cancel")]
    public async Task<IActionResult> Cancel(Guid id, [FromBody] CancelPolicyRequest request)
    {
        var result = await _mediator.Send(new CancelPolicyCommand(id, request.Reason));
        return HandleResult(result, "Policy cancelled successfully.");
    }

    [HttpPost("{id}/renew")]
    public async Task<IActionResult> Renew(Guid id, [FromBody] RenewPolicyRequest request)
    {
        var result = await _mediator.Send(new RenewPolicyCommand(id, request.TenureMonths));
        return result.Match<IActionResult>(
            onSuccess: r => Ok(new { r.NewPolicyId, r.NewPolicyNumber, message = "Policy renewed successfully." }),
            onFailure: error => HandleErrorResult(error)
        );
    }

    // ===================== Grace Period & Renewal =====================

    [HttpGet("{id}/grace-period")]
    public async Task<IActionResult> GetGracePeriod(Guid id)
    {
        var result = await _mediator.Send(new GetGracePeriodQuery(id));
        return result != null ? Ok(result) : NotFound();
    }

    [HttpGet("{id}/renewal-schedule")]
    public async Task<IActionResult> GetRenewalSchedule(Guid id)
    {
        var result = await _mediator.Send(new GetRenewalScheduleQuery(id));
        return result != null ? Ok(result) : NotFound();
    }

    // ===================== Nominees =====================

    [HttpGet("{policyId}/nominees")]
    public async Task<ActionResult<List<NomineeDto>>> ListNominees(Guid policyId)
    {
        return Ok(await _mediator.Send(new ListNomineesQuery(policyId)));
    }

    [HttpPost("{policyId}/nominees")]
    public async Task<IActionResult> AddNominee(Guid policyId, [FromBody] AddNomineeRequest request)
    {
        var command = new AddNomineeCommand(policyId, request.BeneficiaryId, request.FullName, request.Relationship, request.SharePercentage,
            request.DateOfBirth, request.NidNumber, request.PhoneNumber, request.NomineeDobText);
        var result = await _mediator.Send(command);
        return result.Match<IActionResult>(
            onSuccess: id => Created($"api/policies/{policyId}/nominees/{id}",
                new { nomineeId = id, message = "Nominee added successfully." }),
            onFailure: error => HandleErrorResult(error)
        );
    }

    [HttpPut("{policyId}/nominees/{nomineeId}")]
    public async Task<IActionResult> UpdateNominee(Guid policyId, Guid nomineeId, [FromBody] UpdateNomineeRequest request)
    {
        var command = new UpdateNomineeCommand(policyId, nomineeId, request.FullName, request.Relationship, request.SharePercentage,
            request.DateOfBirth, request.NidNumber, request.PhoneNumber, request.NomineeDobText);
        var result = await _mediator.Send(command);
        return HandleResult(result, "Nominee updated successfully.");
    }

    [HttpDelete("{policyId}/nominees/{nomineeId}")]
    public async Task<IActionResult> DeleteNominee(Guid policyId, Guid nomineeId)
    {
        var result = await _mediator.Send(new DeleteNomineeCommand(policyId, nomineeId));
        return HandleResult(result, "Nominee deleted successfully.");
    }

    // ===================== Helpers =====================

    private IActionResult HandleResult(SharedKernel.CQRS.Result result, string successMessage)
    {
        if (result.IsSuccess) return Ok(new { message = successMessage });
        return HandleErrorResult(result.Error!);
    }

    private IActionResult HandleErrorResult(SharedKernel.CQRS.Error error)
    {
        return error.Code switch
        {
            "NOT_FOUND" => NotFound(new { error.Code, error.Message }),
            "INVALID_STATE_TRANSITION" => Conflict(new { error.Code, error.Message }),
            "VALIDATION_ERROR" => BadRequest(new { error.Code, error.Message }),
            _ => BadRequest(new { error.Code, error.Message })
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
