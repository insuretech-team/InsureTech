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

public class CreateProductCommandHandler : IRequestHandler<CreateProductCommand, Result<Product>>
{
    private readonly IProductRepository _repository;
    private readonly ILogger<CreateProductCommandHandler> _logger;

    public CreateProductCommandHandler(
        IProductRepository repository,
        ILogger<CreateProductCommandHandler> logger)
    {
        _repository = repository;
        _logger = logger;
    }

    public async Task<Result<Product>> Handle(CreateProductCommand request, CancellationToken cancellationToken)
    {
        try
        {
            // Validate product code uniqueness
            var existing = await _repository.GetByCodeAsync(request.ProductCode, cancellationToken);
            if (existing != null)
            {
                return Result<Product>.Failure($"Product code '{request.ProductCode}' already exists");
            }

            // Validate business rules
            if (request.MinSumInsuredAmount >= request.MaxSumInsuredAmount)
            {
                return Result<Product>.Failure("Min sum insured must be less than max sum insured");
            }

            if (request.MinTenureMonths >= request.MaxTenureMonths)
            {
                return Result<Product>.Failure("Min tenure must be less than max tenure");
            }

            if (request.BasePremiumAmount <= 0)
            {
                return Result<Product>.Failure("Base premium must be greater than zero");
            }

            // Create product entity
            var product = new Product
            {
                ProductId = Guid.NewGuid().ToString(),
                ProductCode = request.ProductCode,
                ProductName = request.ProductName,
                Category = request.Category,
                Description = request.Description,
                BasePremium = new Money { Amount = request.BasePremiumAmount },
                MinSumInsured = new Money { Amount = request.MinSumInsuredAmount },
                MaxSumInsured = new Money { Amount = request.MaxSumInsuredAmount },
                MinTenureMonths = request.MinTenureMonths,
                MaxTenureMonths = request.MaxTenureMonths,
                Status = ProductStatus.Draft,
                CreatedBy = request.CreatedBy,
                BasePremiumCurrency = "BDT",
                MinSumInsuredCurrency = "BDT",
                MaxSumInsuredCurrency = "BDT"
            };

            product.Exclusions.AddRange(request.Exclusions);

            // Save to database via Go proxy
            var created = await _repository.CreateAsync(product, cancellationToken);

            _logger.LogInformation(
                "Product created: {ProductId} - {ProductCode} - {ProductName}",
                created.ProductId,
                created.ProductCode,
                created.ProductName
            );

            return Result<Product>.Success(created);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to create product: {ProductCode}", request.ProductCode);
            return Result<Product>.Failure($"Failed to create product: {ex.Message}");
        }
    }
}
