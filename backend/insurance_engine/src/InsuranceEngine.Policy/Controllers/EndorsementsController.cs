using System;
using System.Threading.Tasks;
using MediatR;
using Microsoft.AspNetCore.Mvc;
using InsuranceEngine.Policy.Application.Features.Commands.Endorsements;
using InsuranceEngine.Policy.Application.Features.Queries.Endorsements;
using InsuranceEngine.Policy.Domain.Enums;

namespace InsuranceEngine.Policy.Controllers;

[ApiController]
[Route("api/[controller]")]
public class EndorsementsController : ControllerBase
{
    private readonly IMediator _mediator;

    public EndorsementsController(IMediator mediator)
    {
        _mediator = mediator;
    }

    [HttpPost]
    public async Task<IActionResult> Submit([FromBody] SubmitEndorsementCommand command)
    {
        var result = await _mediator.Send(command);
        return result.IsSuccess ? Ok(result.Value) : BadRequest(result.Error);
    }

    [HttpPost("{id}/approve")]
    public async Task<IActionResult> Approve(Guid id, [FromBody] Guid approvedBy)
    {
        var result = await _mediator.Send(new ApproveEndorsementCommand(id, approvedBy));
        return result.IsSuccess ? Ok() : BadRequest(result.Error);
    }

    [HttpGet("{id}")]
    public async Task<IActionResult> Get(Guid id)
    {
        var result = await _mediator.Send(new GetEndorsementQuery(id));
        return result.IsSuccess ? Ok(result.Value) : NotFound(result.Error);
    }

    [HttpGet]
    public async Task<IActionResult> List([FromQuery] Guid? policyId, [FromQuery] EndorsementStatus? status, [FromQuery] int page = 1, [FromQuery] int pageSize = 10)
    {
        var result = await _mediator.Send(new ListEndorsementsQuery(policyId, status, page, pageSize));
        return result.IsSuccess ? Ok(result.Value) : BadRequest(result.Error);
    }
}
