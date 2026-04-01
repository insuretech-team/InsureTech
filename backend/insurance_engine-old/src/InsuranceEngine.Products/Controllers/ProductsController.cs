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
using InsuranceEngine.SharedKernel.DTOs;

namespace InsuranceEngine.Products.Controllers;

[ApiController]
[Route("v1/products")]
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
    public async Task<IActionResult> List(
        [FromQuery] ProductCategory? category = null,
        [FromQuery] int page = 1,
        [FromQuery] int pageSize = 20)
    {
        var result = await _mediator.Send(new ListProductsQuery(category, page, pageSize));
        return Ok(new ProductsListingResponse(result.Items, result.TotalCount, result.Page, result.PageSize));
    }

    /// <summary>
    /// Get product by UUID
    /// </summary>
    [HttpGet("{id}")]
    public async Task<IActionResult> Get(Guid id)
    {
        var result = await _mediator.Send(new GetProductQuery(id));
        if (result == null) return NotFound(new ErrorDto("PRODUCT_NOT_FOUND", "Product not found."));
        return Ok(new ProductRetrievalResponse(result));
    }

    /// <summary>
    /// Full-text search products
    /// </summary>
    [HttpGet("search")]
    public async Task<IActionResult> Search(
        [FromQuery] string? q,
        [FromQuery] decimal? minPremium,
        [FromQuery] decimal? maxPremium)
    {
        var items = await _mediator.Send(new SearchProductsQuery(q, minPremium, maxPremium));
        return Ok(new ProductsListingResponse(items, items.Count, 1, items.Count));
    }

    /// <summary>
    /// Create product (Admin only). Returns 201 with location header.
    /// </summary>
    [HttpPost]
    public async Task<IActionResult> Create([FromBody] CreateProductCommand command)
    {
        var result = await _mediator.Send(command);
        if (result.IsSuccess)
        {
            var id = result.Value!;
            return CreatedAtAction(nameof(Get), new { id }, new ProductCreationResponse(id, "Product created successfully."));
        }
        return BadRequest(MapError(result.Error!));
    }

    /// <summary>
    /// Full update of product in DRAFT status
    /// </summary>
    [HttpPut("{id}")]
    public async Task<IActionResult> Update(Guid id, [FromBody] UpdateProductCommand command)
    {
        if (id != command.Id) return BadRequest(new ErrorDto("VALIDATION_ERROR", "Route ID does not match body ID."));

        var result = await _mediator.Send(command);
        if (result.IsSuccess) return NoContent();

        return HandleErrorResult(result.Error!);
    }

    /// <summary>
    /// Transition DRAFT → ACTIVE
    /// </summary>
    [HttpPost("{id}/activate")]
    public async Task<IActionResult> Activate(Guid id)
    {
        var result = await _mediator.Send(new ActivateProductCommand(id));
        if (result.IsSuccess) return Ok(new ProductUpdateResponse("Product activated successfully."));

        return HandleErrorResult(result.Error!);
    }

    /// <summary>
    /// Transition ACTIVE → INACTIVE
    /// </summary>
    [HttpPost("{id}/deactivate")]
    public async Task<IActionResult> Deactivate(Guid id, [FromBody] ReasonRequest? request = null)
    {
        var result = await _mediator.Send(new DeactivateProductCommand(id, request?.Reason));
        if (result.IsSuccess) return Ok(new ProductUpdateResponse("Product deactivated successfully."));

        return HandleErrorResult(result.Error!);
    }

    /// <summary>
    /// Transition any → DISCONTINUED
    /// </summary>
    [HttpPost("{id}/discontinue")]
    public async Task<IActionResult> Discontinue(Guid id, [FromBody] ReasonRequest? request = null)
    {
        var result = await _mediator.Send(new DiscontinueProductCommand(id, request?.Reason));
        if (result.IsSuccess) return Ok(new ProductUpdateResponse("Product discontinued successfully."));

        return HandleErrorResult(result.Error!);
    }


    [HttpPost("{id}/calculate-premium")]
    public async Task<IActionResult> CalculatePremium(Guid id, [FromBody] CalculatePremiumRequest request)
    {
        var result = await _mediator.Send(new CalculatePremiumCommand(
            id, request.SumInsuredAmount, request.TenureMonths, request.RiderIds, request.ApplicantData));

        if (result.IsSuccess)
        {
            var r = result.Value!;
            return Ok(new PremiumCalculationResponse(r.BasePremium, r.RiderPremium, r.TotalPremium, r.Breakdown));
        }
        return HandleErrorResult(result.Error!);
    }

    // ===================== Helpers =====================

    private ErrorDto MapError(InsuranceEngine.SharedKernel.CQRS.Error error)
    {
        return new ErrorDto(error.Code, error.Message);
    }

    private IActionResult HandleErrorResult(InsuranceEngine.SharedKernel.CQRS.Error error)
    {
        var errorDto = MapError(error);
        return error.Code switch
        {
            "NOT_FOUND" => NotFound(errorDto),
            "INVALID_STATE_TRANSITION" => Conflict(errorDto),
            "CONFLICT" => Conflict(errorDto),
            "VALIDATION_ERROR" => BadRequest(errorDto),
            _ => BadRequest(errorDto)
        };
    }
}

public record ReasonRequest(string? Reason);
