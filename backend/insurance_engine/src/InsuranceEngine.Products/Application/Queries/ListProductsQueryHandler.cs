using System;
using System.Collections.Generic;
using System.Threading;
using System.Threading.Tasks;
using MediatR;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using Dapper;
using Insuretech.Products.Services.V1;
using Insuretech.Products.Entity.V1;
using Insuretech.Common.V1;
using Microsoft.Extensions.Caching.Distributed;
using System.Text.Json;

namespace InsuranceEngine.Products.Application.Queries;

public sealed class ListProductsQueryHandler : IRequestHandler<ListProductsQuery, ListProductsResponse>
{
    private readonly DbContext _dbContext;
    private readonly ILogger<ListProductsQueryHandler> _logger;
    private readonly IDistributedCache _cache;

    public ListProductsQueryHandler(
        DbContext dbContext, 
        ILogger<ListProductsQueryHandler> logger,
        IDistributedCache cache)
    {
        _dbContext = dbContext;
        _logger = logger;
        _cache = cache;
    }

    public async Task<ListProductsResponse> Handle(ListProductsQuery request, CancellationToken cancellationToken)
    {
        string cacheKey = $"products_list_{request.Category}_{request.Page}_{request.PageSize}";
        var cachedData = await _cache.GetStringAsync(cacheKey, cancellationToken);
        
        if (!string.IsNullOrEmpty(cachedData))
        {
            return JsonSerializer.Deserialize<ListProductsResponse>(cachedData)!;
        }

        try
        {
            var sql = @"
                SELECT product_id, product_code, product_name, description, category,
                       base_premium, min_sum_insured, max_sum_insured,
                       min_tenure_months, max_tenure_months, status, created_at
                FROM insurance_schema.products
                WHERE (@Category IS NULL OR category = @Category)
                  AND (@SearchTerm IS NULL OR product_name ILIKE @SearchTerm OR product_code ILIKE @SearchTerm)
                  AND deleted_at IS NULL
                ORDER BY created_at DESC
                LIMIT @PageSize OFFSET @Offset";

            var offset = (request.Page - 1) * request.PageSize;

            using var connection = _dbContext.Database.GetDbConnection();
            if (connection.State != System.Data.ConnectionState.Open) await connection.OpenAsync(cancellationToken);

            var items = await connection.QueryAsync<dynamic>(sql, new
            {
                Category = request.Category?.ToString(),
                SearchTerm = string.IsNullOrEmpty(request.SearchTerm) ? null : $"%{request.SearchTerm}%",
                PageSize = request.PageSize,
                Offset = offset
            });

            var countSql = @"
                SELECT COUNT(*) FROM insurance_schema.products
                WHERE (@Category IS NULL OR category = @Category)
                  AND deleted_at IS NULL";

            var totalCount = await connection.ExecuteScalarAsync<int>(countSql, new
            {
                Category = request.Category?.ToString(),
                SearchTerm = string.IsNullOrEmpty(request.SearchTerm) ? null : $"%{request.SearchTerm}%"
            });

            var response = new ListProductsResponse
            {
                TotalCount = totalCount,
                Page = request.Page,
                PageSize = request.PageSize
            };

            foreach (var item in items)
            {
                var product = new Insuretech.Products.Entity.V1.Product
                {
                    ProductId = item.product_id?.ToString() ?? "",
                    ProductCode = item.product_code?.ToString() ?? "",
                    ProductName = item.product_name?.ToString() ?? "",
                    Description = item.description?.ToString() ?? ""
                };

                string categoryStr = item.category?.ToString() ?? "";
                if (System.Enum.TryParse<ProductCategory>(categoryStr, true, out var cat)) product.Category = cat;

                string statusStr = item.status?.ToString() ?? "";
                if (System.Enum.TryParse<ProductStatus>(statusStr, true, out var stat)) product.Status = stat;

                product.BasePremium = new Money { Amount = (long)((decimal)(item.base_premium ?? 0) * 100), Currency = "BDT" };
                product.MinSumInsured = new Money { Amount = (long)((decimal)(item.min_sum_insured ?? 0) * 100), Currency = "BDT" };
                product.MaxSumInsured = new Money { Amount = (long)((decimal)(item.max_sum_insured ?? 0) * 100), Currency = "BDT" };
                product.MinTenureMonths = item.min_tenure_months ?? 0;
                product.MaxTenureMonths = item.max_tenure_months ?? 0;

                if (item.created_at != null)
                {
                    product.CreatedAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(DateTime.SpecifyKind((DateTime)item.created_at, DateTimeKind.Utc));
                }

                response.Products.Add(product);
            }

            // Cache the result for 5 minutes (FR-028)
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
