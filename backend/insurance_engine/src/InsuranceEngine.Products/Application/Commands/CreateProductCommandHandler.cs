using MediatR;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using Insuretech.Products.Services.V1;
using Dapper;
using System.Data;
using InsuranceEngine.Products.Domain;
using InsuranceEngine.SharedKernel.Domain;

namespace InsuranceEngine.Products.Application.Commands;

public sealed class CreateProductCommandHandler : IRequestHandler<CreateProductCommand, CreateProductResponse>
{
    private readonly DbContext _dbContext;
    private readonly ILogger<CreateProductCommandHandler> _logger;

    public CreateProductCommandHandler(
        DbContext dbContext,
        ILogger<CreateProductCommandHandler> logger)
    {
        _dbContext = dbContext;
        _logger = logger;
    }

    public async Task<CreateProductResponse> Handle(CreateProductCommand request, CancellationToken cancellationToken)
    {
        var connection = _dbContext.Database.GetDbConnection();
        if (connection.State != ConnectionState.Open) await connection.OpenAsync(cancellationToken);

        using var transaction = await connection.BeginTransactionAsync(cancellationToken);
        try
        {
            // 1. Create Domain Aggregate (DDD)
            var product = Product.Create(
                code: request.ProductCode,
                name: request.ProductName,
                category: Enum.TryParse<Insuretech.Products.Entity.V1.ProductCategory>(request.Category, true, out var cat) ? cat : Insuretech.Products.Entity.V1.ProductCategory.Unspecified,
                enDesc: request.Description ?? "",
                bnDesc: "", // Default empty for now or handle from request if available
                basePremium: Money.FromDecimal(request.BasePremium),
                minSum: Money.FromDecimal(request.MinSumInsured),
                maxSum: Money.FromDecimal(request.MaxSumInsured),
                minTenure: request.MinTenureMonths,
                maxTenure: request.MaxTenureMonths
            );

            // 2. Persist Aggregate State using Dapper
            var insertProductSql = @"
                INSERT INTO insurance_schema.products (
                    product_id, product_code, product_name, category, description,
                    base_premium, min_sum_insured, max_sum_insured,
                    min_tenure_months, max_tenure_months,
                    status, created_at, created_by
                ) VALUES (
                    @ProductId, @ProductCode, @ProductName, @Category, @Description,
                    @BasePremium, @MinSumInsured, @MaxSumInsured,
                    @MinTenureMonths, @MaxTenureMonths,
                    @Status, @CreatedAt, @CreatedBy
                )";

            await connection.ExecuteAsync(insertProductSql, new
            {
                ProductId = product.Id,
                ProductCode = product.ProductCode,
                ProductName = product.ProductName,
                Category = product.Category.ToString(),
                Description = product.Description.English,
                BasePremium = product.BasePremium.Amount,
                MinSumInsured = product.MinSumInsured.Amount,
                MaxSumInsured = product.MaxSumInsured.Amount,
                MinTenureMonths = product.MinTenureMonths,
                MaxTenureMonths = product.MaxTenureMonths,
                Status = product.Status.ToString(),
                CreatedAt = product.CreatedAt,
                CreatedBy = string.IsNullOrEmpty(request.CreatedBy) || request.CreatedBy.ToUpper() == "SYSTEM" 
                    ? Guid.Parse("00000000-0000-0000-0000-000000000001")  // System user UUID
                    : Guid.Parse(request.CreatedBy)
            }, transaction);

            await transaction.CommitAsync(cancellationToken);

            _logger.LogInformation("Product created successfully: {ProductId} ({ProductCode})", product.Id, product.ProductCode);

            return new CreateProductResponse
            {
                ProductId = product.Id.ToString(),
                Message = "Product created successfully"
            };
        }
        catch (Exception ex)
        {
            await transaction.RollbackAsync(cancellationToken);
            _logger.LogError(ex, "Failed to create product {ProductCode}", request.ProductCode);
            throw;
        }
    }
}
