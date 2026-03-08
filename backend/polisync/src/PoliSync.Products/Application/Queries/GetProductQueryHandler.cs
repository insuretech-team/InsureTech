using Insuretech.Products.Entity.V1;
using MediatR;
using Microsoft.Extensions.Logging;
using PoliSync.Products.Infrastructure;
using PoliSync.SharedKernel.CQRS;
using System;
using System.Threading;
using System.Threading.Tasks;

namespace PoliSync.Products.Application.Queries;

public class GetProductQueryHandler : IRequestHandler<GetProductQuery, Result<Product?>>
{
    private readonly IProductRepository _repository;
    private readonly ILogger<GetProductQueryHandler> _logger;

    public GetProductQueryHandler(
        IProductRepository repository,
        ILogger<GetProductQueryHandler> logger)
    {
        _repository = repository;
        _logger = logger;
    }

    public async Task<Result<Product?>> Handle(GetProductQuery request, CancellationToken cancellationToken)
    {
        try
        {
            var product = await _repository.GetByIdAsync(request.ProductId, cancellationToken);
            
            if (product == null)
            {
                _logger.LogWarning("Product not found: {ProductId}", request.ProductId);
                return Result<Product?>.Failure($"Product not found: {request.ProductId}");
            }

            return Result<Product?>.Success(product);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to get product: {ProductId}", request.ProductId);
            return Result<Product?>.Failure($"Failed to get product: {ex.Message}");
        }
    }
}
