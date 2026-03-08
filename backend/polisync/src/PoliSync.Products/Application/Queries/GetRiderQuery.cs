using Microsoft.Extensions.Logging;
using PoliSync.Products.Application.Mappers;
using PoliSync.Products.Domain;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Products.Application.Queries;

public record GetRiderQuery(
    Guid RiderId,
    Guid ProductId
) : IQuery<Insuretech.Products.Entity.V1.Rider>;

public sealed class GetRiderQueryHandler : IQueryHandler<GetRiderQuery, Insuretech.Products.Entity.V1.Rider>
{
    private readonly IProductRepository _productRepository;
    private readonly ILogger<GetRiderQueryHandler> _logger;

    public GetRiderQueryHandler(
        IProductRepository productRepository,
        ILogger<GetRiderQueryHandler> logger)
    {
        _productRepository = productRepository;
        _logger = logger;
    }

    public async Task<Result<Insuretech.Products.Entity.V1.Rider>> Handle(GetRiderQuery request, CancellationToken ct)
    {
        // Load product
        var product = await _productRepository.GetByIdAsync(request.ProductId, ct);
        if (product is null)
        {
            _logger.LogWarning("Product not found: {ProductId}", request.ProductId);
            return Result<Insuretech.Products.Entity.V1.Rider>.NotFound($"Product '{request.ProductId}' not found");
        }

        // Find rider
        var rider = product.Riders.FirstOrDefault(r => r.Id == request.RiderId);
        if (rider is null)
        {
            _logger.LogWarning("Rider not found: {RiderId} in product {ProductId}", request.RiderId, request.ProductId);
            return Result<Insuretech.Products.Entity.V1.Rider>.NotFound($"Rider '{request.RiderId}' not found");
        }

        // Map to proto
        var protoRider = rider.ToProto();

        _logger.LogInformation("Rider retrieved: {RiderId}", request.RiderId);

        return Result<Insuretech.Products.Entity.V1.Rider>.Ok(protoRider);
    }
}
