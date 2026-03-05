using MediatR;
using Microsoft.AspNetCore.Mvc;
using PoliSync.Products.Application.Commands;
using PoliSync.Products.Application.Queries;
using PoliSync.Products.Domain;

namespace PoliSync.ApiHost.Controllers;

/// <summary>
/// REST API controller for Product operations.
/// Maps to proto ProductService RPC definitions.
/// </summary>
[ApiController]
[Route("v1/products")]
public class ProductsController : ControllerBase
{
    private readonly IMediator _mediator;

    public ProductsController(IMediator mediator) => _mediator = mediator;

    /// <summary>List all products with optional category filter and pagination.</summary>
    [HttpGet]
    public async Task<IActionResult> ListProducts(
        [FromQuery] ProductCategory? category,
        [FromQuery] int page = 1,
        [FromQuery] int pageSize = 20,
        CancellationToken ct = default)
    {
        var result = await _mediator.Send(new ListProductsQuery(category, page, pageSize), ct);
        if (result.IsFailure) return MapError(result.Error!);

        var data = result.Value!;
        return Ok(new
        {
            products = data.Products.Select(MapProductDto),
            total_count = data.TotalCount,
            page = data.Page,
            page_size = data.PageSize
        });
    }

    /// <summary>Get product details by ID.</summary>
    [HttpGet("{productId:guid}")]
    public async Task<IActionResult> GetProduct(Guid productId, CancellationToken ct = default)
    {
        var result = await _mediator.Send(new GetProductQuery(productId), ct);
        if (result.IsFailure) return MapError(result.Error!);

        return Ok(new { product = MapProductDto(result.Value!) });
    }

    /// <summary>Create a new product.</summary>
    [HttpPost]
    public async Task<IActionResult> CreateProduct([FromBody] CreateProductRequest req, CancellationToken ct = default)
    {
        var result = await _mediator.Send(new CreateProductCommand(
            ProductCode: req.ProductCode,
            ProductName: req.ProductName,
            Category: req.Category,
            BasePremium: req.BasePremium,
            MinSumInsured: req.MinSumInsured,
            MaxSumInsured: req.MaxSumInsured,
            MinTenureMonths: req.MinTenureMonths,
            MaxTenureMonths: req.MaxTenureMonths,
            CreatedBy: req.CreatedBy,
            Description: req.Description,
            Exclusions: req.Exclusions,
            ProductAttributes: req.ProductAttributes
        ), ct);

        if (result.IsFailure) return MapError(result.Error!);

        return Created($"/v1/products/{result.Value}", new
        {
            product_id = result.Value,
            message = "Product created successfully"
        });
    }

    /// <summary>Update an existing product.</summary>
    [HttpPatch("{productId:guid}")]
    public async Task<IActionResult> UpdateProduct(Guid productId, [FromBody] UpdateProductRequest req, CancellationToken ct = default)
    {
        var result = await _mediator.Send(new UpdateProductCommand(
            ProductId: productId,
            ProductName: req.ProductName,
            Description: req.Description,
            BasePremium: req.BasePremium,
            MinSumInsured: req.MinSumInsured,
            MaxSumInsured: req.MaxSumInsured,
            MinTenureMonths: req.MinTenureMonths,
            MaxTenureMonths: req.MaxTenureMonths,
            Exclusions: req.Exclusions,
            ProductAttributes: req.ProductAttributes
        ), ct);

        if (result.IsFailure) return MapError(result.Error!);
        return Ok(new { message = "Product updated successfully" });
    }

    /// <summary>Search products by name, category, or premium range.</summary>
    [HttpGet("search")]
    public async Task<IActionResult> SearchProducts(
        [FromQuery] string? query,
        [FromQuery] ProductCategory? category,
        [FromQuery] long? minPremium,
        [FromQuery] long? maxPremium,
        CancellationToken ct = default)
    {
        var result = await _mediator.Send(new SearchProductsQuery(query, category, minPremium, maxPremium), ct);
        if (result.IsFailure) return MapError(result.Error!);

        var data = result.Value!;
        return Ok(new
        {
            products = data.Products.Select(MapProductDto),
            total_count = data.TotalCount
        });
    }

    /// <summary>Activate a product (make available for sale).</summary>
    [HttpPost("{productId:guid}:activate")]
    public async Task<IActionResult> ActivateProduct(Guid productId, CancellationToken ct = default)
    {
        var result = await _mediator.Send(new ActivateProductCommand(productId), ct);
        if (result.IsFailure) return MapError(result.Error!);
        return Ok(new { message = "Product activated successfully" });
    }

    /// <summary>Deactivate a product (temporarily disable).</summary>
    [HttpPost("{productId:guid}:deactivate")]
    public async Task<IActionResult> DeactivateProduct(Guid productId, [FromBody] ReasonRequest? req = null, CancellationToken ct = default)
    {
        var result = await _mediator.Send(new DeactivateProductCommand(productId, req?.Reason), ct);
        if (result.IsFailure) return MapError(result.Error!);
        return Ok(new { message = "Product deactivated successfully" });
    }

    /// <summary>Discontinue a product (permanently remove).</summary>
    [HttpPost("{productId:guid}:discontinue")]
    public async Task<IActionResult> DiscontinueProduct(Guid productId, [FromBody] ReasonRequest? req = null, CancellationToken ct = default)
    {
        var result = await _mediator.Send(new DiscontinueProductCommand(productId, req?.Reason), ct);
        if (result.IsFailure) return MapError(result.Error!);
        return Ok(new { message = "Product discontinued successfully" });
    }

    // ── Private helpers ──────────────────────────────────────────────

    private static object MapProductDto(Product p) => new
    {
        product_id = p.ProductId,
        product_code = p.ProductCode,
        product_name = p.ProductName,
        category = p.Category.ToString(),
        description = p.Description,
        base_premium = p.BasePremium,
        base_premium_currency = p.BasePremiumCurrency,
        min_sum_insured = p.MinSumInsured,
        max_sum_insured = p.MaxSumInsured,
        min_tenure_months = p.MinTenureMonths,
        max_tenure_months = p.MaxTenureMonths,
        exclusions = p.Exclusions,
        status = p.Status.ToString(),
        plans = p.Plans.Select(plan => new
        {
            plan_id = plan.PlanId,
            plan_name = plan.PlanName,
            plan_description = plan.PlanDescription,
            premium_amount = plan.PremiumAmount,
            min_sum_insured = plan.MinSumInsured,
            max_sum_insured = plan.MaxSumInsured
        }),
        riders = p.AvailableRiders.Select(r => new
        {
            rider_id = r.RiderId,
            rider_name = r.RiderName,
            description = r.Description,
            premium_amount = r.PremiumAmount,
            coverage_amount = r.CoverageAmount,
            is_mandatory = r.IsMandatory
        }),
        created_by = p.CreatedBy,
        created_at = p.CreatedAt,
        updated_at = p.UpdatedAt
    };

    private IActionResult MapError(SharedKernel.CQRS.ResultError error) => error.Kind switch
    {
        SharedKernel.CQRS.ResultErrorKind.NotFound => NotFound(new { error = error.Message }),
        SharedKernel.CQRS.ResultErrorKind.Conflict => Conflict(new { error = error.Message }),
        SharedKernel.CQRS.ResultErrorKind.Unauthorized => Unauthorized(new { error = error.Message }),
        SharedKernel.CQRS.ResultErrorKind.Validation => BadRequest(new { error = error.Message }),
        _ => StatusCode(500, new { error = error.Message })
    };
}

// ── Request DTOs ─────────────────────────────────────────────────────

public record CreateProductRequest
{
    public string ProductCode { get; init; } = string.Empty;
    public string ProductName { get; init; } = string.Empty;
    public ProductCategory Category { get; init; }
    public long BasePremium { get; init; }
    public long MinSumInsured { get; init; }
    public long MaxSumInsured { get; init; }
    public int MinTenureMonths { get; init; }
    public int MaxTenureMonths { get; init; }
    public string CreatedBy { get; init; } = string.Empty;
    public string? Description { get; init; }
    public List<string>? Exclusions { get; init; }
    public string? ProductAttributes { get; init; }
}

public record UpdateProductRequest
{
    public string? ProductName { get; init; }
    public string? Description { get; init; }
    public long? BasePremium { get; init; }
    public long? MinSumInsured { get; init; }
    public long? MaxSumInsured { get; init; }
    public int? MinTenureMonths { get; init; }
    public int? MaxTenureMonths { get; init; }
    public List<string>? Exclusions { get; init; }
    public string? ProductAttributes { get; init; }
}

public record ReasonRequest
{
    public string? Reason { get; init; }
}
