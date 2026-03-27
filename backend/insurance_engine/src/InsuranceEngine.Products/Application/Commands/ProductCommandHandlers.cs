using MediatR;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using InsuranceEngine.SharedKernel.CQRS;
using Dapper;

namespace InsuranceEngine.Products.Application.Commands;

public sealed class UpdateProductCommandHandler : IRequestHandler<UpdateProductCommand, Result<bool>>
{
    private readonly DbContext _dbContext;
    private readonly ILogger<UpdateProductCommandHandler> _logger;

    public UpdateProductCommandHandler(DbContext dbContext, ILogger<UpdateProductCommandHandler> logger)
    {
        _dbContext = dbContext;
        _logger = logger;
    }

    public async Task<Result<bool>> Handle(UpdateProductCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var sql = @"
                UPDATE insurance_schema.products
                SET product_name = COALESCE(@ProductName, product_name),
                    description = COALESCE(@Description, description),
                    base_premium = COALESCE(@BasePremium, base_premium),
                    min_sum_insured = COALESCE(@MinSumInsured, min_sum_insured),
                    max_sum_insured = COALESCE(@MaxSumInsured, max_sum_insured),
                    min_tenure_months = COALESCE(@MinTenureMonths, min_tenure_months),
                    max_tenure_months = COALESCE(@MaxTenureMonths, max_tenure_months),
                    updated_at = @UpdatedAt
                WHERE product_id = @ProductId::uuid AND deleted_at IS NULL";

            using var connection = _dbContext.Database.GetDbConnection();
            await connection.OpenAsync(cancellationToken);
            
            var rows = await connection.ExecuteAsync(sql, new
            {
                ProductId = request.ProductId,
                request.ProductName,
                request.Description,
                request.BasePremium,
                request.MinSumInsured,
                request.MaxSumInsured,
                request.MinTenureMonths,
                request.MaxTenureMonths,
                UpdatedAt = DateTime.UtcNow
            });

            if (rows == 0)
                return Result<bool>.NotFound("PRODUCT_NOT_FOUND", "Product not found");

            _logger.LogInformation("Product updated: {ProductId}", request.ProductId);
            return Result<bool>.Ok(true);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to update product {ProductId}", request.ProductId);
            return Result<bool>.Fail("PRODUCT_UPDATE_FAILED", ex.Message);
        }
    }
}

public sealed class ActivateProductCommandHandler : IRequestHandler<ActivateProductCommand, Result<bool>>
{
    private readonly DbContext _dbContext;
    private readonly ILogger<ActivateProductCommandHandler> _logger;

    public ActivateProductCommandHandler(DbContext dbContext, ILogger<ActivateProductCommandHandler> logger)
    {
        _dbContext = dbContext;
        _logger = logger;
    }

    public async Task<Result<bool>> Handle(ActivateProductCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var sql = @"
                UPDATE insurance_schema.products
                SET status = 'ACTIVE', updated_at = @UpdatedAt
                WHERE product_id = @ProductId::uuid AND deleted_at IS NULL";

            using var connection = _dbContext.Database.GetDbConnection();
            await connection.OpenAsync(cancellationToken);
            
            var rows = await connection.ExecuteAsync(sql, new
            {
                ProductId = request.ProductId,
                UpdatedAt = DateTime.UtcNow
            });

            if (rows == 0)
                return Result<bool>.NotFound("PRODUCT_NOT_FOUND", "Product not found");

            _logger.LogInformation("Product activated: {ProductId}", request.ProductId);
            return Result<bool>.Ok(true);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to activate product {ProductId}", request.ProductId);
            return Result<bool>.Fail("PRODUCT_STATUS_UPDATE_FAILED", ex.Message);
        }
    }
}

public sealed class DeactivateProductCommandHandler : IRequestHandler<DeactivateProductCommand, Result<bool>>
{
    private readonly DbContext _dbContext;
    private readonly ILogger<DeactivateProductCommandHandler> _logger;

    public DeactivateProductCommandHandler(DbContext dbContext, ILogger<DeactivateProductCommandHandler> logger)
    {
        _dbContext = dbContext;
        _logger = logger;
    }

    public async Task<Result<bool>> Handle(DeactivateProductCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var sql = @"
                UPDATE insurance_schema.products
                SET status = 'INACTIVE', updated_at = @UpdatedAt
                WHERE product_id = @ProductId::uuid AND deleted_at IS NULL";

            using var connection = _dbContext.Database.GetDbConnection();
            await connection.OpenAsync(cancellationToken);
            
            var rows = await connection.ExecuteAsync(sql, new
            {
                ProductId = request.ProductId,
                UpdatedAt = DateTime.UtcNow
            });

            if (rows == 0)
                return Result<bool>.NotFound("PRODUCT_NOT_FOUND", "Product not found");

            _logger.LogInformation("Product deactivated: {ProductId}", request.ProductId);
            return Result<bool>.Ok(true);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to deactivate product {ProductId}", request.ProductId);
            return Result<bool>.Fail("PRODUCT_STATUS_UPDATE_FAILED", ex.Message);
        }
    }
}

public sealed class DiscontinueProductCommandHandler : IRequestHandler<DiscontinueProductCommand, Result<bool>>
{
    private readonly DbContext _dbContext;
    private readonly ILogger<DiscontinueProductCommandHandler> _logger;

    public DiscontinueProductCommandHandler(DbContext dbContext, ILogger<DiscontinueProductCommandHandler> logger)
    {
        _dbContext = dbContext;
        _logger = logger;
    }

    public async Task<Result<bool>> Handle(DiscontinueProductCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var sql = @"
                UPDATE insurance_schema.products
                SET status = 'DISCONTINUED', updated_at = @UpdatedAt
                WHERE product_id = @ProductId::uuid AND deleted_at IS NULL";

            using var connection = _dbContext.Database.GetDbConnection();
            await connection.OpenAsync(cancellationToken);
            
            var rows = await connection.ExecuteAsync(sql, new
            {
                ProductId = request.ProductId,
                UpdatedAt = DateTime.UtcNow
            });

            if (rows == 0)
                return Result<bool>.NotFound("PRODUCT_NOT_FOUND", "Product not found");

            _logger.LogInformation("Product discontinued: {ProductId}", request.ProductId);
            return Result<bool>.Ok(true);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to discontinue product {ProductId}", request.ProductId);
            return Result<bool>.Fail("PRODUCT_STATUS_UPDATE_FAILED", ex.Message);
        }
    }
}
