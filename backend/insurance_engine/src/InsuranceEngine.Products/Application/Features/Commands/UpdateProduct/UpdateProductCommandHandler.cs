using System;
using System.Threading;
using System.Threading.Tasks;
using MediatR;
using InsuranceEngine.Products.Application.Interfaces;
using InsuranceEngine.Products.Domain.Enums;
using InsuranceEngine.Products.Domain.Events;
using InsuranceEngine.SharedKernel.CQRS;
using InsuranceEngine.SharedKernel.Interfaces;

namespace InsuranceEngine.Products.Application.Features.Commands.UpdateProduct;

public class UpdateProductCommandHandler : IRequestHandler<UpdateProductCommand, Result>
{
    private readonly IProductRepository _productRepository;
    private readonly IEventBus _eventBus;

    public UpdateProductCommandHandler(IProductRepository productRepository, IEventBus eventBus)
    {
        _productRepository = productRepository;
        _eventBus = eventBus;
    }

    public async Task<Result> Handle(UpdateProductCommand request, CancellationToken cancellationToken)
    {
        var product = await _productRepository.GetByIdAsync(request.Id);
        if (product == null)
            return Result.Fail(Error.NotFound("Product", request.Id.ToString()));

        if (product.Status != ProductStatus.Draft)
            return Result.Fail(Error.InvalidStateTransition(
                "Product can only be updated in DRAFT status."));

        var changedFields = new List<string>();

        if (product.ProductName != request.ProductName) { product.ProductName = request.ProductName; changedFields.Add("ProductName"); }
        if (product.Description != request.Description) { product.Description = request.Description; changedFields.Add("Description"); }
        if (product.Category != request.Category) { product.Category = request.Category; changedFields.Add("Category"); }
        if (product.BasePremium != request.BasePremiumAmount) { product.BasePremium = request.BasePremiumAmount; changedFields.Add("BasePremium"); }
        if (product.MinSumInsured != request.MinSumInsuredAmount) { product.MinSumInsured = request.MinSumInsuredAmount; changedFields.Add("MinSumInsured"); }
        if (product.MaxSumInsured != request.MaxSumInsuredAmount) { product.MaxSumInsured = request.MaxSumInsuredAmount; changedFields.Add("MaxSumInsured"); }
        if (product.MinTenureMonths != request.MinTenureMonths) { product.MinTenureMonths = request.MinTenureMonths; changedFields.Add("MinTenureMonths"); }
        if (product.MaxTenureMonths != request.MaxTenureMonths) { product.MaxTenureMonths = request.MaxTenureMonths; changedFields.Add("MaxTenureMonths"); }
        if (request.Exclusions != null) { product.Exclusions = request.Exclusions; changedFields.Add("Exclusions"); }

        product.UpdatedAt = DateTime.UtcNow;

        await _productRepository.UpdateAsync(product);

        if (changedFields.Count > 0)
        {
            await _eventBus.PublishAsync("product.events", new ProductUpdatedDomainEvent(
                ProductId: product.Id,
                ProductCode: product.ProductCode,
                UpdatedFields: changedFields
            ));
        }

        return Result.Ok();
    }
}
