using Microsoft.Extensions.Logging;
using PoliSync.Products.Application.Mappers;
using PoliSync.Products.Domain;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Products.Application.Queries;

public record ListRidersQuery(
    Guid ProductId
) : IQuery<List<Insuretech.Products.Entity.V1.Rider>>;

public sealed class ListRidersQueryHandler : IQueryHandler<ListRidersQuery, List<Insuretech.Products.Entity.V1.Rider>>
{
    private readonly IProductRepository _productRepository;
    private readonly ILogger<ListRidersQueryHandler> _logger;

    public ListRidersQueryHandler(
        IProductRepository productRepository,
        ILogger<ListRidersQueryHandler> logger)
    {
        _productRepository = productRepository;
        _logger = logger;
    }

    public async Task<Result<List<Insuretech.Products.Entity.V1.Rider>>> Handle(ListRidersQuery request, CancellationToken ct)
    {
        // Load product
        var product = await _productRepository.GetByIdAsync(request.ProductId, ct);
        if (product is null)
        {
            _logger.LogWarning("Product not found: {ProductId}", request.ProductId);
            return Result<List<Insuretech.Products.Entity.V1.Rider>>.NotFound($"Product '{request.ProductId}' not found");
        }

        // Map riders to proto
        var protoRiders = product.Riders.Select(r => r.ToProto()).ToList();

        _logger.LogInformation("Listed riders for product {ProductId}: {Count}", request.ProductId, protoRiders.Count);

        return Result<List<Insuretech.Products.Entity.V1.Rider>>.Ok(protoRiders);
    }
}
