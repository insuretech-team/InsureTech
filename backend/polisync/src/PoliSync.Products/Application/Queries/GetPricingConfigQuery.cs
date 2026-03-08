using Microsoft.Extensions.Logging;
using PoliSync.Products.Application.Mappers;
using PoliSync.Products.Domain;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Products.Application.Queries;

public record GetPricingConfigQuery(
    Guid ProductId
) : IQuery<Insuretech.Products.Entity.V1.PricingConfig>;

public sealed class GetPricingConfigQueryHandler : IQueryHandler<GetPricingConfigQuery, Insuretech.Products.Entity.V1.PricingConfig>
{
    private readonly IProductRepository _productRepository;
    private readonly ILogger<GetPricingConfigQueryHandler> _logger;

    public GetPricingConfigQueryHandler(
        IProductRepository productRepository,
        ILogger<GetPricingConfigQueryHandler> logger)
    {
        _productRepository = productRepository;
        _logger = logger;
    }

    public async Task<Result<Insuretech.Products.Entity.V1.PricingConfig>> Handle(GetPricingConfigQuery request, CancellationToken ct)
    {
        // Load product
        var product = await _productRepository.GetByIdAsync(request.ProductId, ct);
        if (product is null)
        {
            _logger.LogWarning("Product not found: {ProductId}", request.ProductId);
            return Result<Insuretech.Products.Entity.V1.PricingConfig>.NotFound($"Product '{request.ProductId}' not found");
        }

        // Get pricing config
        if (product.PricingConfig is null)
        {
            _logger.LogWarning("Pricing config not found for product {ProductId}", request.ProductId);
            return Result<Insuretech.Products.Entity.V1.PricingConfig>.NotFound($"Pricing config not found for product '{request.ProductId}'");
        }

        // Map to proto
        var protoConfig = product.PricingConfig.ToProto();

        _logger.LogInformation("Pricing config retrieved for product {ProductId}", request.ProductId);

        return Result<Insuretech.Products.Entity.V1.PricingConfig>.Ok(protoConfig);
    }
}
