using Insuretech.Common.V1;
using Insuretech.Products.Entity.V1;
using MediatR;
using Microsoft.Extensions.Logging;
using PoliSync.Products.Infrastructure;
using PoliSync.SharedKernel.CQRS;
using System;
using System.Threading;
using System.Threading.Tasks;

namespace PoliSync.Products.Application.Commands;

public class UpdateProductCommandHandler : IRequestHandler<UpdateProductCommand, Result<Product>>
{
    private readonly IProductRepository _repository;
    private readonly ILogger<UpdateProductCommandHandler> _logger;

    public UpdateProductCommandHandler(
        IProductRepository repository,
        ILogger<UpdateProductCommandHandler> logger)
    {
        _repository = repository;
        _logger = logger;
    }

    public async Task<Result<Product>> Handle(UpdateProductCommand request, CancellationToken cancellationToken)
    {
        try
        {
            // Fetch existing product
            var product = await _repository.GetByIdAsync(request.ProductId, cancellationToken);
            if (product == null)
            {
                return Result<Product>.Failure($"Product not found: {request.ProductId}");
            }

            // Update fields if provided
            if (!string.IsNullOrEmpty(request.ProductName))
            {
                product.ProductName = request.ProductName;
            }

            if (!string.IsNullOrEmpty(request.Description))
            {
                product.Description = request.Description;
            }

            if (request.BasePremiumAmount.HasValue)
            {
                if (request.BasePremiumAmount.Value <= 0)
                {
                    return Result<Product>.Failure("Base premium must be greater than zero");
                }
                product.BasePremium = new Money { Amount = request.BasePremiumAmount.Value };
            }

            if (request.MinSumInsuredAmount.HasValue)
            {
                product.MinSumInsured = new Money { Amount = request.MinSumInsuredAmount.Value };
            }

            if (request.MaxSumInsuredAmount.HasValue)
            {
                product.MaxSumInsured = new Money { Amount = request.MaxSumInsuredAmount.Value };
            }

            // Validate sum insured range
            if (product.MinSumInsured.Amount >= product.MaxSumInsured.Amount)
            {
                return Result<Product>.Failure("Min sum insured must be less than max sum insured");
            }

            if (request.MinTenureMonths.HasValue)
            {
                product.MinTenureMonths = request.MinTenureMonths.Value;
            }

            if (request.MaxTenureMonths.HasValue)
            {
                product.MaxTenureMonths = request.MaxTenureMonths.Value;
            }

            // Validate tenure range
            if (product.MinTenureMonths >= product.MaxTenureMonths)
            {
                return Result<Product>.Failure("Min tenure must be less than max tenure");
            }

            if (request.Exclusions != null)
            {
                product.Exclusions.Clear();
                product.Exclusions.AddRange(request.Exclusions);
            }

            if (request.Status.HasValue)
            {
                product.Status = request.Status.Value;
            }

            // Save changes
            var updated = await _repository.UpdateAsync(product, cancellationToken);

            _logger.LogInformation("Product updated: {ProductId}", request.ProductId);

            return Result<Product>.Success(updated);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to update product: {ProductId}", request.ProductId);
            return Result<Product>.Failure($"Failed to update product: {ex.Message}");
        }
    }
}
