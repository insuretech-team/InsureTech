using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Products.Services.V1;
using Insuretech.Products.Entity.V1;
using Insuretech.Common.V1;
using InsuranceEngine.Grpc.Gateways;

namespace InsuranceEngine.Products.Application.Commands;

public sealed class CreateProductCommandHandler : IRequestHandler<CreateProductCommand, CreateProductResponse>
{
    private readonly IProductDataGateway _gateway;
    private readonly ILogger<CreateProductCommandHandler> _logger;

    public CreateProductCommandHandler(
        IProductDataGateway gateway,
        ILogger<CreateProductCommandHandler> logger)
    {
        _gateway = gateway;
        _logger = logger;
    }

    public async Task<CreateProductResponse> Handle(CreateProductCommand request, CancellationToken cancellationToken)
    {
        try
        {
            // Note: Validation logic (uniqueness etc.) is now handled by the Go backend (SSOT).
            // We map directly from the C# command to the Proto Product message.

            var product = new Product
            {
                ProductCode = request.ProductCode,
                ProductName = request.ProductName,
                Category = Enum.TryParse<ProductCategory>(request.Category, true, out var category) ? category : ProductCategory.Unspecified,
                Description = request.Description,
                BasePremium = new Money { Amount = (long)(request.BasePremium * 100), Currency = "BDT" },
                MinSumInsured = new Money { Amount = (long)(request.MinSumInsured * 100), Currency = "BDT" },
                MaxSumInsured = new Money { Amount = (long)(request.MaxSumInsured * 100), Currency = "BDT" },
                MinTenureMonths = request.MinTenureMonths,
                MaxTenureMonths = request.MaxTenureMonths,
                Status = ProductStatus.Active
            };

            var createdProduct = await _gateway.CreateProductAsync(product, cancellationToken);

            _logger.LogInformation("Product created successfully via Go SSOT: {ProductId} ({ProductCode})", 
                createdProduct.ProductId, createdProduct.ProductCode);

            return new CreateProductResponse
            {
                ProductId = createdProduct.ProductId,
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
