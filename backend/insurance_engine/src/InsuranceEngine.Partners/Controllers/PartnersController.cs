using System;
using System.Collections.Generic;
using System.Threading.Tasks;
using InsuranceEngine.Partners.Application.Features.Commands.Agents;
using InsuranceEngine.Partners.Application.Features.Commands.Partners;
using InsuranceEngine.Partners.Application.Features.Queries;
using InsuranceEngine.Partners.Application.Features.Queries.Partners;
using InsuranceEngine.SharedKernel.CQRS;
using MediatR;
using Microsoft.AspNetCore.Mvc;

namespace InsuranceEngine.Partners.Controllers;

[ApiController]
[Route("api/partners")]
public class PartnersController : ControllerBase
{
    private readonly IMediator _mediator;

    public PartnersController(IMediator mediator)
    {
        _mediator = mediator;
    }

    [HttpPost]
    public async Task<ActionResult<Guid>> CreatePartner(CreatePartnerCommand command)
    {
        var result = await _mediator.Send(command);
        return result.IsSuccess ? Ok(result.Value) : BadRequest(result.Error);
    }

    [HttpGet]
    public async Task<ActionResult<List<PartnerDto>>> ListPartners()
    {
        var result = await _mediator.Send(new ListPartnersQuery());
        return result.IsSuccess ? Ok(result.Value) : BadRequest(result.Error);
    }

    [HttpGet("{id}")]
    public async Task<ActionResult<PartnerDto>> GetPartner(Guid id)
    {
        var result = await _mediator.Send(new GetPartnerQuery(id));
        return result.IsSuccess ? Ok(result.Value) : NotFound(result.Error);
    }

    [HttpPost("{id}/agents")]
    public async Task<ActionResult> CreateAgent(Guid id, CreateAgentCommand command)
    {
        if (id != command.PartnerId) return BadRequest("Partner ID mismatch");
        var result = await _mediator.Send(command);
        return result.IsSuccess ? Ok() : BadRequest(result.Error);
    }
}
