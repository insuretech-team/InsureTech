using MediatR;
using Microsoft.Extensions.Logging;
using InsuranceEngine.SharedKernel.CQRS;
using InsuranceEngine.Grpc.Gateways;
using Insuretech.Products.Entity.V1;
using Insuretech.Common.V1;

namespace InsuranceEngine.Products.Application.Commands;

public sealed class UpdateProductCommandHandler : IRequestHandler<UpdateProductCommand, Result<bool>>
{
    private readonly IProductDataGateway _gateway;
    private readonly ILogger<UpdateProductCommandHandler> _logger;

    public UpdateProductCommandHandler(IProductDataGateway gateway, ILogger<UpdateProductCommandHandler> logger)
    {
        _gateway = gateway;
        _logger = logger;
    }

    public async Task<Result<bool>> Handle(UpdateProductCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var product = await _gateway.GetProductAsync(request.ProductId, cancellationToken);
            if (product == null)
                return Result<bool>.NotFound("PRODUCT_NOT_FOUND", "Product not found");

            // Update fields in the Proto object
            if (request.ProductName != null) product.ProductName = request.ProductName;
            if (request.Description != null) product.Description = request.Description;
            
            if (request.BasePremium.HasValue) 
                product.BasePremium = new Money { Amount = (long)(request.BasePremium.Value * 100), Currency = "BDT" };
            
            if (request.MinSumInsured.HasValue) 
                product.MinSumInsured = new Money { Amount = (long)(request.MinSumInsured.Value * 100), Currency = "BDT" };
            
            if (request.MaxSumInsured.HasValue) 
                product.MaxSumInsured = new Money { Amount = (long)(request.MaxSumInsured.Value * 100), Currency = "BDT" };
            
            if (request.MinTenureMonths.HasValue) product.MinTenureMonths = request.MinTenureMonths.Value;
            if (request.MaxTenureMonths.HasValue) product.MaxTenureMonths = request.MaxTenureMonths.Value;
            
            await _gateway.UpdateProductAsync(product, cancellationToken);

            _logger.LogInformation("Product updated via Go Gateway: {ProductId}", request.ProductId);
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
    private readonly IProductDataGateway _gateway;
    private readonly ILogger<ActivateProductCommandHandler> _logger;

    public ActivateProductCommandHandler(IProductDataGateway gateway, ILogger<ActivateProductCommandHandler> logger)
    {
        _gateway = gateway;
        _logger = logger;
    }

    public async Task<Result<bool>> Handle(ActivateProductCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var product = await _gateway.GetProductAsync(request.ProductId, cancellationToken);
            if (product == null)
                return Result<bool>.NotFound("PRODUCT_NOT_FOUND", "Product not found");

            product.Status = ProductStatus.Active;
            await _gateway.UpdateProductAsync(product, cancellationToken);

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
    private readonly IProductDataGateway _gateway;
    private readonly ILogger<DeactivateProductCommandHandler> _logger;

    public DeactivateProductCommandHandler(IProductDataGateway gateway, ILogger<DeactivateProductCommandHandler> logger)
    {
        _gateway = gateway;
        _logger = logger;
    }

    public async Task<Result<bool>> Handle(DeactivateProductCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var product = await _gateway.GetProductAsync(request.ProductId, cancellationToken);
            if (product == null)
                return Result<bool>.NotFound("PRODUCT_NOT_FOUND", "Product not found");

            product.Status = ProductStatus.Inactive;
            await _gateway.UpdateProductAsync(product, cancellationToken);

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
    private readonly IProductDataGateway _gateway;
    private readonly ILogger<DiscontinueProductCommandHandler> _logger;

    public DiscontinueProductCommandHandler(IProductDataGateway gateway, ILogger<DiscontinueProductCommandHandler> logger)
    {
        _gateway = gateway;
        _logger = logger;
    }

    public async Task<Result<bool>> Handle(DiscontinueProductCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var product = await _gateway.GetProductAsync(request.ProductId, cancellationToken);
            if (product == null)
                return Result<bool>.NotFound("PRODUCT_NOT_FOUND", "Product not found");

            product.Status = ProductStatus.Discontinued;
            await _gateway.UpdateProductAsync(product, cancellationToken);

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
