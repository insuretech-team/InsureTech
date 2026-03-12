using System;
using System.Threading;
using System.Threading.Tasks;
using MediatR;
using InsuranceEngine.Products.Application.Interfaces;
using InsuranceEngine.Products.Domain;
using InsuranceEngine.Products.Domain.Enums;
using InsuranceEngine.Products.Domain.Events;
using InsuranceEngine.SharedKernel.CQRS;
using InsuranceEngine.SharedKernel.Interfaces;

namespace InsuranceEngine.Products.Application.Features.Commands.CreateProduct;

public class CreateProductCommandHandler : IRequestHandler<CreateProductCommand, Result<Guid>>
{
    private readonly IProductRepository _productRepository;
    private readonly IEventBus _eventBus;

    public CreateProductCommandHandler(IProductRepository productRepository, IEventBus eventBus)
    {
        _productRepository = productRepository;
        _eventBus = eventBus;
    }

    public async Task<Result<Guid>> Handle(CreateProductCommand request, CancellationToken cancellationToken)
    {
        // Check for duplicate product code
        var existing = await _productRepository.GetByCodeAsync(request.ProductCode);
        if (existing != null)
            return Result<Guid>.Fail(Error.Conflict($"Product with code '{request.ProductCode}' already exists."));

        var product = new Product
        {
            Id = Guid.NewGuid(),
            ProductCode = request.ProductCode,
            ProductName = request.ProductName,
            ProductNameBn = request.ProductNameBn,
            Description = request.Description,
            Category = request.Category,
            Status = ProductStatus.Draft,
            BasePremiumAmount = request.BasePremiumAmount,
            BasePremiumCurrency = "BDT",
            MinSumInsuredAmount = request.MinSumInsuredAmount,
            MinSumInsuredCurrency = "BDT",
            MaxSumInsuredAmount = request.MaxSumInsuredAmount,
            MaxSumInsuredCurrency = "BDT",
            MinAge = request.MinAge,
            MaxAge = request.MaxAge,
            MinTenureMonths = request.MinTenureMonths,
            MaxTenureMonths = request.MaxTenureMonths,
            Exclusions = request.Exclusions ?? new(),
            CreatedBy = request.CreatedBy,
            CreatedAt = DateTime.UtcNow,
            UpdatedAt = DateTime.UtcNow
        };

        var id = await _productRepository.AddAsync(product);

        await _eventBus.PublishAsync("product.events", new ProductCreatedDomainEvent(
            ProductId: id,
            ProductCode: product.ProductCode,
            ProductName: product.ProductName,
            Category: product.Category.ToString(),
            BasePremium: product.BasePremiumAmount,
            CreatedBy: product.CreatedBy
        ));

        return Result<Guid>.Ok(id);
    }
}
