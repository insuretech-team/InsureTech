using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Products.Services.V1;
using InsuranceEngine.Grpc.Gateways;
using Microsoft.Extensions.Caching.Distributed;
using System.Text.Json;

namespace InsuranceEngine.Products.Application.Queries;

public sealed class ListProductsQueryHandler : IRequestHandler<ListProductsQuery, ListProductsResponse>
{
    private readonly IProductDataGateway _gateway;
    private readonly ILogger<ListProductsQueryHandler> _logger;
    private readonly IDistributedCache _cache;

    public ListProductsQueryHandler(
        IProductDataGateway gateway, 
        ILogger<ListProductsQueryHandler> logger,
        IDistributedCache cache)
    {
        _gateway = gateway;
        _logger = logger;
        _cache = cache;
    }

    public async Task<ListProductsResponse> Handle(ListProductsQuery request, CancellationToken cancellationToken)
    {
        string cacheKey = $"products_list_{request.Category}_{request.Page}_{request.PageSize}";
        var cachedData = await _cache.GetStringAsync(cacheKey, cancellationToken);
        
        if (!string.IsNullOrEmpty(cachedData))
        {
            try 
            {
                return JsonSerializer.Deserialize<ListProductsResponse>(cachedData)!;
            }
            catch (JsonException)
            {
                _logger.LogWarning("Failed to deserialize product list cache for key {CacheKey}", cacheKey);
            }
        }

        try
        {
            // Note: The current gRPC contract does not support category filtering.
            // We list all products from the gateway and handle pagination.
            var items = await _gateway.ListProductsAsync(request.Page, request.PageSize, cancellationToken);

            var response = new ListProductsResponse
            {
                TotalCount = items.Count, // Temporary: ideally the gateway returns the total
                Page = request.Page,
                PageSize = request.PageSize
            };

            response.Products.AddRange(items);

            // Cache for 5 minutes (FR-028)
            var cacheOptions = new DistributedCacheEntryOptions { AbsoluteExpirationRelativeToNow = TimeSpan.FromMinutes(5) };
            await _cache.SetStringAsync(cacheKey, JsonSerializer.Serialize(response), cacheOptions, cancellationToken);

            return response;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to list products");
            throw;
        }
    }
}
