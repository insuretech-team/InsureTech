using Insuretech.Products.Entity.V1;
using MediatR;
using Microsoft.Extensions.Logging;
using PoliSync.Products.Infrastructure;
using PoliSync.SharedKernel.CQRS;
using System;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;

namespace PoliSync.Products.Application.Queries;

public class ListProductsQueryHandler : IRequestHandler<ListProductsQuery, Result<ListProductsResult>>
{
    private readonly IProductRepository _repository;
    private readonly ILogger<ListProductsQueryHandler> _logger;

    public ListProductsQueryHandler(
        IProductRepository repository,
        ILogger<ListProductsQueryHandler> logger)
    {
        _repository = repository;
        _logger = logger;
    }

    public async Task<Result<ListProductsResult>> Handle(ListProductsQuery request, CancellationToken cancellationToken)
    {
        try
        {
            var products = request.Category.HasValue
                ? await _repository.GetByCategoryAsync(request.Category.Value, cancellationToken)
                : await _repository.GetAllAsync(request.Page, request.PageSize, cancellationToken);

            // Filter by status if specified
            if (request.Status.HasValue)
            {
                products = products.Where(p => p.Status == request.Status.Value).ToList();
            }

            var result = new ListProductsResult
            {
                Products = products,
                TotalCount = products.Count
            };

            _logger.LogInformation("Listed {Count} products", products.Count);

            return Result<ListProductsResult>.Success(result);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to list products");
            return Result<ListProductsResult>.Failure($"Failed to list products: {ex.Message}");
        }
    }
}
