using System;
using System.Collections.Generic;
using System.Threading.Tasks;
using InsuranceEngine.Commission.Application.Features.Queries;
using InsuranceEngine.Commission.Application.Features.Queries.Commissions;
using MediatR;
using Microsoft.AspNetCore.Mvc;

namespace InsuranceEngine.Commission.Controllers;

[ApiController]
[Route("api/commissions")]
public class CommissionsController : ControllerBase
{
    private readonly IMediator _mediator;

    public CommissionsController(IMediator mediator)
    {
        _mediator = mediator;
    }

    [HttpGet("recipient/{id}")]
    public async Task<ActionResult<List<CommissionDto>>> GetRecipientCommissions(Guid id)
    {
        var result = await _mediator.Send(new GetRecipientCommissionsQuery(id));
        return result.IsSuccess ? Ok(result.Value) : BadRequest(result.Error);
    }
}
