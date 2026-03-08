using System;
using System.Collections.Generic;
using System.Threading.Tasks;
using MediatR;
using Microsoft.AspNetCore.Mvc;
using InsuranceEngine.Application.Features.Products.Queries.ListProducts;
using InsuranceEngine.Application.Features.Products.Queries.GetProduct;
using InsuranceEngine.Application.Features.Products.Queries.SearchProducts;
using InsuranceEngine.Application.Features.Products.Commands.CreateProduct;
using InsuranceEngine.Application.Features.Products.Commands.UpdateProduct;
using InsuranceEngine.Application.DTOs;

namespace InsuranceEngine.Api.Controllers;

[ApiController]
[Route("api/[controller]")]
public class ProductsController : ControllerBase
{
    private readonly IMediator _mediator;

    public ProductsController(IMediator mediator)
    {
        _mediator = mediator;
    }

    [HttpGet]
    public async Task<ActionResult<List<ProductDto>>> List()
    {
        return Ok(await _mediator.Send(new ListProductsQuery()));
    }

    [HttpGet("{id}")]
    public async Task<ActionResult<ProductDto>> Get(Guid id)
    {
        var result = await _mediator.Send(new GetProductQuery(id));
        return result != null ? Ok(result) : NotFound();
    }

    [HttpGet("search")]
    public async Task<ActionResult<List<ProductDto>>> Search([FromQuery] string? query, [FromQuery] decimal? minPremium, [FromQuery] decimal? maxPremium)
    {
        return Ok(await _mediator.Send(new SearchProductsQuery(query, minPremium, maxPremium)));
    }

    [HttpPost]
    public async Task<ActionResult<Guid>> Create(CreateProductCommand command)
    {
        var id = await _mediator.Send(command);
        return Ok(id);
    }

    [HttpPut("{id}")]
    public async Task<IActionResult> Update(Guid id, UpdateProductCommand command)
    {
        if (id != command.Id) return BadRequest();
        var success = await _mediator.Send(command);
        return success ? NoContent() : NotFound();
    }
}
