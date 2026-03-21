using System;
using System.Collections.Generic;
using System.Threading.Tasks;
using MediatR;
using Microsoft.AspNetCore.Mvc;
using InsuranceEngine.Claims.Application.Features.Commands.Claims;
using InsuranceEngine.Claims.Application.Features.Queries.Claims;
using InsuranceEngine.Claims.Application.DTOs;

namespace InsuranceEngine.Claims.Controllers;

[ApiController]
[Route("api/[controller]")]
public class ClaimsController : ControllerBase
{
    private readonly IMediator _mediator;

    public ClaimsController(IMediator mediator)
    {
        _mediator = mediator;
    }

    [HttpPost]
    public async Task<IActionResult> SubmitClaim([FromBody] SubmitClaimRestRequest request)
    {
        var command = new SubmitClaimCommand(
            request.PolicyId,
            request.CustomerId,
            request.Type,
            request.ClaimedAmount,
            request.IncidentDate,
            request.IncidentDescription,
            request.PlaceOfIncident,
            request.BankDetailsForPayout
        );

        var result = await _mediator.Send(command);

        if (result.IsSuccess)
        {
            return CreatedAtAction(nameof(GetClaim), new { id = result.Value }, result.Value);
        }

        return BadRequest(result.Error);
    }

    [HttpGet("{id}")]
    public async Task<IActionResult> GetClaim(Guid id)
    {
        var result = await _mediator.Send(new GetClaimByIdQuery(id));

        if (result.IsSuccess)
        {
            return Ok(result.Value);
        }

        return NotFound(result.Error);
    }

    [HttpGet("customer/{customerId}")]
    public async Task<IActionResult> ListByCustomer(Guid customerId, [FromQuery] int page = 1, [FromQuery] int pageSize = 10)
    {
        var result = await _mediator.Send(new ListClaimsByCustomerQuery(customerId, page, pageSize));

        if (result.IsSuccess)
        {
            return Ok(result.Value);
        }

        return BadRequest(result.Error);
    }

    [HttpPost("{id}/approve")]
    public async Task<IActionResult> ApproveClaim(Guid id, [FromBody] ApproveClaimRestRequest request)
    {
        var command = new ApproveClaimCommand(
            id,
            request.ApproverId,
            request.ApproverRole,
            request.ApprovalLevel,
            request.Decision,
            request.ApprovedAmount,
            request.Notes
        );

        var result = await _mediator.Send(command);
        return result.IsSuccess ? Ok() : BadRequest(result.Error);
    }
}

public record ApproveClaimRestRequest(
    Guid ApproverId,
    string ApproverRole,
    int ApprovalLevel,
    Domain.Enums.ApprovalDecision Decision,
    long ApprovedAmount,
    string Notes
);
