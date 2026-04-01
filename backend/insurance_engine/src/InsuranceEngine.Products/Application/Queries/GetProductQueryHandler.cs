using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Products.Services.V1;
using Insuretech.Products.Entity.V1;
using Insuretech.Common.V1;
using InsuranceEngine.SharedKernel.Persistence;
using InsuranceEngine.SharedKernel.Persistence.Entities;
using Google.Protobuf.WellKnownTypes;

namespace InsuranceEngine.Products.Application.Queries;

public sealed class GetProductQueryHandler : IRequestHandler<GetProductQuery, GetProductResponse>
{
    private readonly IRepository<ProductEntity> _repository;
    private readonly ILogger<GetProductQueryHandler> _logger;

    public GetProductQueryHandler(IRepository<ProductEntity> repository, ILogger<GetProductQueryHandler> logger)
    {
        _repository = repository;
        _logger = logger;
    }

    public async Task<GetProductResponse> Handle(GetProductQuery request, CancellationToken cancellationToken)
    {
        try
        {
            var entity = await _repository.GetByIdAsync(Guid.Parse(request.ProductId), cancellationToken);

            if (entity == null)
            {
                return new GetProductResponse
                {
                    Error = new Insuretech.Common.V1.Error
                    {
                        Code = "PRODUCT_NOT_FOUND",
                        Message = "Product not found"
                    }
                };
            }

            var product = MapToProto(entity);

            return new GetProductResponse { Product = product };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to get product {ProductId}", request.ProductId);
            throw;
        }
    }

    private static Insuretech.Products.Entity.V1.Product MapToProto(ProductEntity entity)
    {
        var product = new Insuretech.Products.Entity.V1.Product
        {
            ProductId = entity.ProductId.ToString(),
            ProductCode = entity.ProductCode,
            ProductName = entity.ProductName,
            Description = entity.Description ?? ""
        };

        if (System.Enum.TryParse<ProductCategory>(entity.Category, true, out var cat)) product.Category = cat;
        if (System.Enum.TryParse<ProductStatus>(entity.Status, true, out var stat)) product.Status = stat;

        product.BasePremium = new Money { Amount = entity.BasePremium, Currency = entity.BasePremiumCurrency };
        product.MinSumInsured = new Money { Amount = entity.MinSumInsured, Currency = entity.MinSumInsuredCurrency };
        product.MaxSumInsured = new Money { Amount = entity.MaxSumInsured, Currency = entity.MaxSumInsuredCurrency };
        product.MinTenureMonths = entity.MinTenureMonths;
        product.MaxTenureMonths = entity.MaxTenureMonths;
        // MinAge, MaxAge, TermsUrl removed as they are not in the Proto definition
        
        product.CreatedAt = Timestamp.FromDateTime(DateTime.SpecifyKind(entity.CreatedAt, DateTimeKind.Utc));
        product.UpdatedAt = Timestamp.FromDateTime(DateTime.SpecifyKind(entity.UpdatedAt, DateTimeKind.Utc));

        return product;
    }
}
