using Microsoft.Extensions.Logging;
using PoliSync.Products.Application.Mappers;
using PoliSync.Products.Domain;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Products.Application.Queries;

public record ListProductPlansQuery(
    Guid ProductId
) : IQuery<List<Insuretech.Products.Entity.V1.ProductPlan>>;

public sealed class ListProductPlansQueryHandler : IQueryHandler<ListProductPlansQuery, List<Insuretech.Products.Entity.V1.ProductPlan>>
{
    private readonly IProductRepository _productRepository;
    private readonly ILogger<ListProductPlansQueryHandler> _logger;

    public ListProductPlansQueryHandler(
        IProductRepository productRepository,
        ILogger<ListProductPlansQueryHandler> logger)
    {
        _productRepository = productRepository;
        _logger = logger;
    }

    public async Task<Result<List<Insuretech.Products.Entity.V1.ProductPlan>>> Handle(ListProductPlansQuery request, CancellationToken ct)
    {
        // Load product
        var product = await _productRepository.GetByIdAsync(request.ProductId, ct);
        if (product is null)
        {
            _logger.LogWarning("Product not found: {ProductId}", request.ProductId);
            return Result<List<Insuretech.Products.Entity.V1.ProductPlan>>.NotFound($"Product '{request.ProductId}' not found");
        }

        // Map plans to proto
        var protoPlans = product.Plans.Select(p => p.ToProto()).ToList();

        _logger.LogInformation("Listed product plans for product {ProductId}: {Count}", request.ProductId, protoPlans.Count);

        return Result<List<Insuretech.Products.Entity.V1.ProductPlan>>.Ok(protoPlans);
    }
}
