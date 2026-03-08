using System;
using System.Threading;
using System.Threading.Tasks;
using MediatR;
using InsuranceEngine.Application.Interfaces;
using InsuranceEngine.Domain.Entities;
using InsuranceEngine.Domain.Enums;

namespace InsuranceEngine.Application.Features.Products.Commands.CreateProduct;

public class CreateProductCommandHandler : IRequestHandler<CreateProductCommand, Guid>
{
    private readonly IProductRepository _productRepository;
    private readonly IEventBus _eventBus;

    public CreateProductCommandHandler(IProductRepository productRepository, IEventBus eventBus)
    {
        _productRepository = productRepository;
        _eventBus = eventBus;
    }

    public async Task<Guid> Handle(CreateProductCommand request, CancellationToken cancellationToken)
    {
        var product = new Product
        {
            Id = Guid.NewGuid(),
            ProductCode = request.ProductCode,
            ProductName = request.ProductName,
            ProductNameBn = request.ProductNameBn,
            Description = request.Description,
            Category = request.Category,
            Status = ProductStatus.Draft,
            MinSumInsured = request.MinSumInsured,
            MaxSumInsured = request.MaxSumInsured,
            MinAge = request.MinAge,
            MaxAge = request.MaxAge,
            InsurerId = request.InsurerId,
            CreatedAt = DateTime.UtcNow,
            UpdatedAt = DateTime.UtcNow
        };

        var id = await _productRepository.AddAsync(product);
        
        await _eventBus.PublishAsync("product-created", new { ProductId = id, request.ProductCode });

        return id;
    }
}
