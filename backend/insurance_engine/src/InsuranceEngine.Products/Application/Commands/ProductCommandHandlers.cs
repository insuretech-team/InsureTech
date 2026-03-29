using MediatR;
using Microsoft.Extensions.Logging;
using InsuranceEngine.SharedKernel.CQRS;
using InsuranceEngine.SharedKernel.Persistence;
using InsuranceEngine.SharedKernel.Persistence.Entities;

namespace InsuranceEngine.Products.Application.Commands;

public sealed class UpdateProductCommandHandler : IRequestHandler<UpdateProductCommand, Result<bool>>
{
    private readonly IRepository<ProductEntity> _repository;
    private readonly ILogger<UpdateProductCommandHandler> _logger;

    public UpdateProductCommandHandler(IRepository<ProductEntity> repository, ILogger<UpdateProductCommandHandler> logger)
    {
        _repository = repository;
        _logger = logger;
    }

    public async Task<Result<bool>> Handle(UpdateProductCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var product = await _repository.GetByIdAsync(Guid.Parse(request.ProductId), cancellationToken);
            if (product == null)
                return Result<bool>.NotFound("PRODUCT_NOT_FOUND", "Product not found");

            // Update fields
            if (request.ProductName != null) product.ProductName = request.ProductName;
            if (request.Description != null) product.Description = request.Description;
            if (request.BasePremium.HasValue) product.BasePremium = (long)(request.BasePremium.Value * 100);
            if (request.MinSumInsured.HasValue) product.MinSumInsured = (long)(request.MinSumInsured.Value * 100);
            if (request.MaxSumInsured.HasValue) product.MaxSumInsured = (long)(request.MaxSumInsured.Value * 100);
            if (request.MinTenureMonths.HasValue) product.MinTenureMonths = request.MinTenureMonths.Value;
            if (request.MaxTenureMonths.HasValue) product.MaxTenureMonths = request.MaxTenureMonths.Value;
            
            product.UpdatedAt = DateTime.UtcNow;
            await _repository.UpdateAsync(product, cancellationToken);

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
    private readonly IRepository<ProductEntity> _repository;
    private readonly ILogger<ActivateProductCommandHandler> _logger;

    public ActivateProductCommandHandler(IRepository<ProductEntity> repository, ILogger<ActivateProductCommandHandler> logger)
    {
        _repository = repository;
        _logger = logger;
    }

    public async Task<Result<bool>> Handle(ActivateProductCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var product = await _repository.GetByIdAsync(Guid.Parse(request.ProductId), cancellationToken);
            if (product == null)
                return Result<bool>.NotFound("PRODUCT_NOT_FOUND", "Product not found");

            product.Status = "ACTIVE";
            product.UpdatedAt = DateTime.UtcNow;
            await _repository.UpdateAsync(product, cancellationToken);

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
    private readonly IRepository<ProductEntity> _repository;
    private readonly ILogger<DeactivateProductCommandHandler> _logger;

    public DeactivateProductCommandHandler(IRepository<ProductEntity> repository, ILogger<DeactivateProductCommandHandler> logger)
    {
        _repository = repository;
        _logger = logger;
    }

    public async Task<Result<bool>> Handle(DeactivateProductCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var product = await _repository.GetByIdAsync(Guid.Parse(request.ProductId), cancellationToken);
            if (product == null)
                return Result<bool>.NotFound("PRODUCT_NOT_FOUND", "Product not found");

            product.Status = "INACTIVE";
            product.UpdatedAt = DateTime.UtcNow;
            await _repository.UpdateAsync(product, cancellationToken);

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
    private readonly IRepository<ProductEntity> _repository;
    private readonly ILogger<DiscontinueProductCommandHandler> _logger;

    public DiscontinueProductCommandHandler(IRepository<ProductEntity> repository, ILogger<DiscontinueProductCommandHandler> logger)
    {
        _repository = repository;
        _logger = logger;
    }

    public async Task<Result<bool>> Handle(DiscontinueProductCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var product = await _repository.GetByIdAsync(Guid.Parse(request.ProductId), cancellationToken);
            if (product == null)
                return Result<bool>.NotFound("PRODUCT_NOT_FOUND", "Product not found");

            product.Status = "DISCONTINUED";
            product.UpdatedAt = DateTime.UtcNow;
            await _repository.UpdateAsync(product, cancellationToken);

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
