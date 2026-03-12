using System;
using System.Collections.Generic;
using System.Threading.Tasks;
using MediatR;
using Microsoft.AspNetCore.Mvc;
using InsuranceEngine.Products.Application.Features.Queries.ListProducts;
using InsuranceEngine.Products.Application.Features.Queries.GetProduct;
using InsuranceEngine.Products.Application.Features.Queries.SearchProducts;
using InsuranceEngine.Products.Application.Features.Commands.CreateProduct;
using InsuranceEngine.Products.Application.Features.Commands.UpdateProduct;
using InsuranceEngine.Products.Application.Features.Commands.ActivateProduct;
using InsuranceEngine.Products.Application.Features.Commands.DeactivateProduct;
using InsuranceEngine.Products.Application.Features.Commands.DiscontinueProduct;
using InsuranceEngine.Products.Application.Features.Commands.CalculatePremium;
using InsuranceEngine.Products.Application.DTOs;
using InsuranceEngine.Products.Domain.Enums;

namespace InsuranceEngine.Products.Controllers;

[ApiController]
[Route("api/[controller]")]
public class ProductsController : ControllerBase
{
    private readonly IMediator _mediator;

    public ProductsController(IMediator mediator)
    {
        _mediator = mediator;
    }

    /// <summary>
    /// List active products with optional category filter and pagination
    /// </summary>
    [HttpGet]
    public async Task<ActionResult<PaginatedResponse<ProductListDto>>> List(
        [FromQuery] ProductCategory? category = null,
        [FromQuery] int page = 1,
        [FromQuery] int pageSize = 20)
    {
        var result = await _mediator.Send(new ListProductsQuery(category, page, pageSize));
        return Ok(result);
    }

    /// <summary>
    /// Get product by UUID
    /// </summary>
    [HttpGet("{id}")]
    public async Task<ActionResult<ProductDto>> Get(Guid id)
    {
        var result = await _mediator.Send(new GetProductQuery(id));
        return result != null ? Ok(result) : NotFound();
    }

    /// <summary>
    /// Full-text search products
    /// </summary>
    [HttpGet("search")]
    public async Task<ActionResult<List<ProductListDto>>> Search(
        [FromQuery] string? q,
        [FromQuery] decimal? minPremium,
        [FromQuery] decimal? maxPremium)
    {
        return Ok(await _mediator.Send(new SearchProductsQuery(q, minPremium, maxPremium)));
    }

    /// <summary>
    /// Create product (Admin only). Returns 201 with location header.
    /// </summary>
    [HttpPost]
    public async Task<IActionResult> Create([FromBody] CreateProductCommand command)
    {
        var result = await _mediator.Send(command);
        return result.Match<IActionResult>(
            onSuccess: id => CreatedAtAction(nameof(Get), new { id }, new { productId = id, message = "Product created successfully." }),
            onFailure: error => error.Code switch
            {
                "CONFLICT" => Conflict(new { error.Code, error.Message }),
                _ => BadRequest(new { error.Code, error.Message })
            }
        );
    }

    /// <summary>
    /// Full update of product in DRAFT status
    /// </summary>
    [HttpPut("{id}")]
    public async Task<IActionResult> Update(Guid id, [FromBody] UpdateProductCommand command)
    {
        if (id != command.Id) return BadRequest(new { Code = "VALIDATION_ERROR", Message = "Route ID does not match body ID." });

        var result = await _mediator.Send(command);
        if (result.IsSuccess) return NoContent();

        return result.Error!.Code switch
        {
            "NOT_FOUND" => NotFound(new { result.Error.Code, result.Error.Message }),
            "INVALID_STATE_TRANSITION" => Conflict(new { result.Error.Code, result.Error.Message }),
            _ => BadRequest(new { result.Error!.Code, result.Error.Message })
        };
    }

    /// <summary>
    /// Transition DRAFT → ACTIVE
    /// </summary>
    [HttpPost("{id}/activate")]
    public async Task<IActionResult> Activate(Guid id)
    {
        var result = await _mediator.Send(new ActivateProductCommand(id));
        if (result.IsSuccess) return Ok(new { message = "Product activated successfully." });

        return result.Error!.Code switch
        {
            "NOT_FOUND" => NotFound(new { result.Error.Code, result.Error.Message }),
            "INVALID_STATE_TRANSITION" => Conflict(new { result.Error.Code, result.Error.Message }),
            _ => BadRequest(new { result.Error!.Code, result.Error.Message })
        };
    }

    /// <summary>
    /// Transition ACTIVE → INACTIVE
    /// </summary>
    [HttpPost("{id}/deactivate")]
    public async Task<IActionResult> Deactivate(Guid id, [FromBody] ReasonRequest? request = null)
    {
        var result = await _mediator.Send(new DeactivateProductCommand(id, request?.Reason));
        if (result.IsSuccess) return Ok(new { message = "Product deactivated successfully." });

        return result.Error!.Code switch
        {
            "NOT_FOUND" => NotFound(new { result.Error.Code, result.Error.Message }),
            "INVALID_STATE_TRANSITION" => Conflict(new { result.Error.Code, result.Error.Message }),
            _ => BadRequest(new { result.Error!.Code, result.Error.Message })
        };
    }

    /// <summary>
    /// Transition any → DISCONTINUED
    /// </summary>
    [HttpPost("{id}/discontinue")]
    public async Task<IActionResult> Discontinue(Guid id, [FromBody] ReasonRequest? request = null)
    {
        var result = await _mediator.Send(new DiscontinueProductCommand(id, request?.Reason));
        if (result.IsSuccess) return Ok(new { message = "Product discontinued successfully." });

        return result.Error!.Code switch
        {
            "NOT_FOUND" => NotFound(new { result.Error.Code, result.Error.Message }),
            "INVALID_STATE_TRANSITION" => Conflict(new { result.Error.Code, result.Error.Message }),
            _ => BadRequest(new { result.Error!.Code, result.Error.Message })
        };
    }

    /// <summary>
    /// Calculate premium for a product
    /// </summary>
    [HttpPost("{id}/calculate-premium")]
    public async Task<IActionResult> CalculatePremium(Guid id, [FromBody] CalculatePremiumRequest request)
    {
        var result = await _mediator.Send(new CalculatePremiumCommand(
            id, request.SumInsuredAmount, request.TenureMonths, request.RiderIds, request.ApplicantData));

        return result.Match<IActionResult>(
            onSuccess: Ok,
            onFailure: error => error.Code switch
            {
                "NOT_FOUND" => NotFound(new { error.Code, error.Message }),
                _ => BadRequest(new { error.Code, error.Message })
            }
        );
    }
}

public record ReasonRequest(string? Reason);
