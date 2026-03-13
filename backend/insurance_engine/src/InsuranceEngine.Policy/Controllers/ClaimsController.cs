using System;
using System.Collections.Generic;
using System.Threading.Tasks;
using MediatR;
using Microsoft.AspNetCore.Mvc;
using InsuranceEngine.Policy.Application.Features.Commands.Claims;
using InsuranceEngine.Policy.Application.Features.Queries.Claims;
using InsuranceEngine.Policy.Application.DTOs;

namespace InsuranceEngine.Policy.Controllers;

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
            request.PlaceOfIncident
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
}
