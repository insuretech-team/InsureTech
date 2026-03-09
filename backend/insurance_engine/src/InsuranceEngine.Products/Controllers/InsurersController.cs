using System;
using System.Collections.Generic;
using System.Threading.Tasks;
using MediatR;
using Microsoft.AspNetCore.Mvc;
using InsuranceEngine.Products.Application.Features.Queries.ListInsurers;
using InsuranceEngine.Products.Domain;

namespace InsuranceEngine.Products.Controllers;

[ApiController]
[Route("api/[controller]")]
public class InsurersController : ControllerBase
{
    private readonly IMediator _mediator;

    public InsurersController(IMediator mediator)
    {
        _mediator = mediator;
    }

    [HttpGet]
    public async Task<ActionResult<List<Insurer>>> List()
    {
        return Ok(await _mediator.Send(new ListInsurersQuery()));
    }
}
