using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Products.Services.V1;
using Insuretech.Products.Entity.V1;
using Insuretech.Common.V1;
using InsuranceEngine.SharedKernel.Persistence;
using InsuranceEngine.SharedKernel.Persistence.Entities;
using Microsoft.Extensions.Caching.Distributed;
using System.Text.Json;
using Google.Protobuf.WellKnownTypes;
using System.Linq.Expressions;

namespace InsuranceEngine.Products.Application.Queries;

public sealed class ListProductsQueryHandler : IRequestHandler<ListProductsQuery, ListProductsResponse>
{
    private readonly IRepository<ProductEntity> _repository;
    private readonly ILogger<ListProductsQueryHandler> _logger;
    private readonly IDistributedCache _cache;

    public ListProductsQueryHandler(
        IRepository<ProductEntity> repository, 
        ILogger<ListProductsQueryHandler> logger,
        IDistributedCache cache)
    {
        _repository = repository;
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
            Expression<Func<ProductEntity, bool>>? predicate = null;
            
            if (request.Category != ProductCategory.Unspecified)
            {
                var categoryStr = request.Category.ToString();
                predicate = p => p.Category == categoryStr;
            }


            var (items, totalCount) = await _repository.GetPagedAsync(
                page: request.Page,
                pageSize: request.PageSize,
                predicate: predicate,
                orderBy: p => p.CreatedAt,
                descending: true,
                cancellationToken: cancellationToken
            );

            var response = new ListProductsResponse
            {
                TotalCount = totalCount,
                Page = request.Page,
                PageSize = request.PageSize
            };

            foreach (var entity in items)
            {
                response.Products.Add(MapToProto(entity));
            }

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

    private static Expression<Func<T, bool>> Combine<T>(Expression<Func<T, bool>> expr1, Expression<Func<T, bool>> expr2)
    {
        var parameter = Expression.Parameter(typeof(T));
        var leftVisitor = new ReplaceExpressionVisitor(expr1.Parameters[0], parameter);
        var left = leftVisitor.Visit(expr1.Body);
        var rightVisitor = new ReplaceExpressionVisitor(expr2.Parameters[0], parameter);
        var right = rightVisitor.Visit(expr2.Body);
        return Expression.Lambda<Func<T, bool>>(Expression.AndAlso(left, right), parameter);
    }

    private class ReplaceExpressionVisitor : ExpressionVisitor
    {
        private readonly Expression _oldValue;
        private readonly Expression _newValue;
        public ReplaceExpressionVisitor(Expression oldValue, Expression newValue)
        {
            _oldValue = oldValue;
            _newValue = newValue;
        }
        public override Expression Visit(Expression node) => node == _oldValue ? _newValue : base.Visit(node);
    }

    private static Insuretech.Products.Entity.V1.Product MapToProto(ProductEntity entity)
    {
        var product = new Insuretech.Products.Entity.V1.Product
        {
            ProductId = entity.ProductId.ToString(),
            ProductCode = entity.ProductCode,
            ProductName = entity.ProductName,
            Description = entity.Description ?? ""
        };

        if (System.Enum.TryParse<ProductCategory>(entity.Category, true, out var cat)) product.Category = cat;
        if (System.Enum.TryParse<ProductStatus>(entity.Status, true, out var stat)) product.Status = stat;

        product.BasePremium = new Money { Amount = entity.BasePremium, Currency = entity.BasePremiumCurrency };
        product.MinSumInsured = new Money { Amount = entity.MinSumInsured, Currency = entity.MinSumInsuredCurrency };
        product.MaxSumInsured = new Money { Amount = entity.MaxSumInsured, Currency = entity.MaxSumInsuredCurrency };
        product.MinTenureMonths = entity.MinTenureMonths;
        product.MaxTenureMonths = entity.MaxTenureMonths;
        // MinAge, MaxAge, TermsUrl removed as they are missing in Proto definition
        
        product.CreatedAt = Timestamp.FromDateTime(DateTime.SpecifyKind(entity.CreatedAt, DateTimeKind.Utc));
        product.UpdatedAt = Timestamp.FromDateTime(DateTime.SpecifyKind(entity.UpdatedAt, DateTimeKind.Utc));

        return product;
    }
}
