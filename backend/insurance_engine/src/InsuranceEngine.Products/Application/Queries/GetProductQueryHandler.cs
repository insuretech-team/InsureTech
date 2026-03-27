using System;
using System.Threading;
using System.Threading.Tasks;
using MediatR;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using Dapper;
using Insuretech.Products.Services.V1;
using Insuretech.Products.Entity.V1;
using Insuretech.Common.V1;

namespace InsuranceEngine.Products.Application.Queries;

public sealed class GetProductQueryHandler : IRequestHandler<GetProductQuery, GetProductResponse>
{
    private readonly DbContext _dbContext;
    private readonly ILogger<GetProductQueryHandler> _logger;

    public GetProductQueryHandler(DbContext dbContext, ILogger<GetProductQueryHandler> logger)
    {
        _dbContext = dbContext;
        _logger = logger;
    }

    public async Task<GetProductResponse> Handle(GetProductQuery request, CancellationToken cancellationToken)
    {
        try
        {
            var sql = @"
                SELECT product_id, product_code, product_name, description, category,
                       base_premium, min_sum_insured, max_sum_insured,
                       min_tenure_months, max_tenure_months, status, created_at
                FROM insurance_schema.products
                WHERE product_id = @ProductId AND deleted_at IS NULL";

            using var connection = _dbContext.Database.GetDbConnection();
            if (connection.State != System.Data.ConnectionState.Open) await connection.OpenAsync(cancellationToken);

            var item = await connection.QueryFirstOrDefaultAsync<dynamic>(sql, new
            {
                ProductId = Guid.Parse(request.ProductId)
            });

            if (item == null) throw new Exception("Product not found");

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

            return new GetProductResponse { Product = product };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to get product {ProductId}", request.ProductId);
            throw;
        }
    }
}
