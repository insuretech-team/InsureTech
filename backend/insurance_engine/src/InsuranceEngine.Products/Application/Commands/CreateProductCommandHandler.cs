using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Products.Services.V1;
using InsuranceEngine.Products.Domain;
using InsuranceEngine.SharedKernel.Domain;
using InsuranceEngine.SharedKernel.Persistence;
using InsuranceEngine.SharedKernel.Persistence.Entities;

namespace InsuranceEngine.Products.Application.Commands;

public sealed class CreateProductCommandHandler : IRequestHandler<CreateProductCommand, CreateProductResponse>
{
    private readonly IRepository<ProductEntity> _productRepository;
    private readonly ILogger<CreateProductCommandHandler> _logger;

    public CreateProductCommandHandler(
        IRepository<ProductEntity> productRepository,
        ILogger<CreateProductCommandHandler> logger)
    {
        _productRepository = productRepository;
        _logger = logger;
    }

    public async Task<CreateProductResponse> Handle(CreateProductCommand request, CancellationToken cancellationToken)
    {
        try
        {
            // 1. Create Domain Aggregate (DDD)
            // Note: In a full VSA with isolated storage, the aggregate might be used.
            // But since we are using EF Core (Option A), we can map directly to the entity 
            // after domain validation logic.
            
            var productCode = request.ProductCode;
            if (string.IsNullOrWhiteSpace(productCode))
                throw new ArgumentException("Product code is required");

            // Check uniqueness
            if (await _productRepository.ExistsAsync(p => p.ProductCode == productCode, cancellationToken))
            {
                return new CreateProductResponse
                {
                    Error = new Insuretech.Common.V1.Error
                    {
                        Code = "DUPLICATE_PRODUCT",
                        Message = $"Product with code {productCode} already exists"
                    }
                };
            }

            // 2. Create Entity and Persist
            var entity = new ProductEntity
            {
                ProductId = Guid.NewGuid(),
                ProductCode = productCode,
                ProductName = request.ProductName,
                Category = request.Category,
                Description = request.Description,
                BasePremium = (long)(request.BasePremium * 100), // Store in paisa
                BasePremiumCurrency = "BDT",
                MinSumInsured = (long)(request.MinSumInsured * 100),
                MaxSumInsured = (long)(request.MaxSumInsured * 100),
                MinTenureMonths = request.MinTenureMonths,
                MaxTenureMonths = request.MaxTenureMonths,
                MinAge = request.MinAge,
                MaxAge = request.MaxAge,
                Status = "ACTIVE", // Start as ACTIVE for now or follow lifecycle
                CreatedAt = DateTime.UtcNow,
                UpdatedAt = DateTime.UtcNow,
                CreatedBy = request.CreatedBy ?? "SYSTEM",
                Version = 1
            };

            await _productRepository.AddAsync(entity, cancellationToken);

            _logger.LogInformation("Product created successfully: {ProductId} ({ProductCode})", entity.ProductId, entity.ProductCode);

            return new CreateProductResponse
            {
                ProductId = entity.ProductId.ToString(),
                Message = "Product created successfully"
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to create product {ProductCode}", request.ProductCode);
            return new CreateProductResponse
            {
                Error = new Insuretech.Common.V1.Error
                {
                    Code = "PRODUCT_CREATION_FAILED",
                    Message = ex.Message
                }
            };
        }
    }
}
